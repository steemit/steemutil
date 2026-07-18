package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/protocol"
	"github.com/steemit/steemutil/wif"

	"github.com/pkg/errors"
)

// errInvalidSignature is returned by Verify when a signature fails any
// non-parse validation check (wrong key, tampered digest, bad format).
// Recoverable parse failures are surfaced as errors too; callers that only
// care about the boolean result can ignore the error.
var errInvalidSignature = errors.New("invalid transaction signature")

type SignedTransaction struct {
	*Transaction
}

func NewSignedTransaction(tx *Transaction) *SignedTransaction {
	if tx.Expiration == nil {
		// Use UTC time to match steemjs behavior
		expiration := time.Now().UTC().Add(600 * time.Second)
		tx.Expiration = &protocol.Time{Time: &expiration}
	}

	return &SignedTransaction{tx}
}

func (tx *SignedTransaction) Serialize() ([]byte, error) {
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)

	if err := encoderObj.Encode(tx.Transaction); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (tx *SignedTransaction) Digest(chain *Chain) ([]byte, error) {
	var msgBuffer bytes.Buffer

	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chain.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode chain ID: %v", chain.ID)
	}

	if _, err := msgBuffer.Write(rawChainID); err != nil {
		return nil, errors.Wrap(err, "failed to write chain ID")
	}

	// Write the serialized transaction.
	rawTx, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return nil, errors.Wrap(err, "failed to write serialized transaction")
	}

	// Compute the digest.
	digest := sha256.Sum256(msgBuffer.Bytes())
	return digest[:], nil
}

func (tx *SignedTransaction) Sign(privKeys []*wif.PrivateKey, chain *Chain) error {
	// Compute digest
	digest, err := tx.Digest(chain)
	if err != nil {
		return err
	}

	// Sign digest
	sigs := make([][]byte, 0, len(privKeys))
	for _, v := range privKeys {
		sig, err := ecdsa.SignCompact(v.Raw.PrivKey, digest, true)
		if err != nil {
			return err
		}
		sigs = append(sigs, sig)
	}

	// Set the signature array in the transaction.
	sigsHex := make([]string, 0, len(sigs))
	for _, sig := range sigs {
		sigsHex = append(sigsHex, hex.EncodeToString(sig))
	}

	tx.Transaction.Signatures = sigsHex
	return nil
}

// Verify validates every signature on the transaction against pubKeys by
// recovering each signer's public key from the compact (recoverable)
// signature and comparing it to the expected key.
//
// Signatures are produced by Sign via ecdsa.SignCompact, i.e. they are
// 65-byte compact signatures whose first byte encodes the recovery id
// (range 31..34 for compressed keys). Verification therefore must use
// recover-and-compare (ecdsa.RecoverCompact), not DER parsing.
//
// Verify is fail-closed: an empty signature slice, a count mismatch, a
// malformed signature, or any key mismatch yields (false, <non-nil error>).
func (tx *SignedTransaction) Verify(pubKeys []*wif.PublicKey, chain *Chain) (bool, error) {
	// Compute digest
	digest, err := tx.Digest(chain)
	if err != nil {
		return false, err
	}

	// Fail-closed: no signatures means nothing to verify.
	if len(tx.Signatures) == 0 {
		return false, errors.Wrap(errInvalidSignature, "no signatures to verify")
	}
	// Each signature corresponds positionally to one pubKey.
	if len(pubKeys) < len(tx.Signatures) {
		return false, errors.Wrapf(errInvalidSignature, "need %d pubkeys, got %d", len(tx.Signatures), len(pubKeys))
	}

	for i, sigHex := range tx.Signatures {
		sig, err := hex.DecodeString(sigHex)
		if err != nil {
			return false, errors.Wrapf(errInvalidSignature, "signature %d is not valid hex", i)
		}

		// Reject anything that is not a well-formed compact signature in the
		// compressed-key recovery range. This mirrors wif.ValidateCompactSignature
		// without the package cycle of importing wif here.
		if _, err := wif.ValidateCompactSignature(sig); err != nil {
			return false, errors.Wrapf(errInvalidSignature, "signature %d: %v", i, err)
		}

		// Recover the public key that produced this signature over the digest.
		recovered, wasCompressed, err := ecdsa.RecoverCompact(sig, digest)
		if err != nil {
			return false, errors.Wrapf(errInvalidSignature, "signature %d recovery failed: %v", i, err)
		}
		if !wasCompressed {
			return false, errors.Wrapf(errInvalidSignature, "signature %d was not compressed", i)
		}

		// Constant-time-independent byte compare of the recovered compressed
		// key against the expected signer's compressed key.
		if !bytes.Equal(recovered.SerializeCompressed(), pubKeys[i].Raw.SerializeCompressed()) {
			return false, errors.Wrapf(errInvalidSignature, "signature %d was not produced by pubkey %d", i, i)
		}
	}
	return true, nil
}

// ID calculates and returns the transaction ID.
// The transaction ID is the SHA256 hash of the serialized transaction (without signatures),
// represented as a 40-character hex string (first 20 bytes of the hash).
func (tx *SignedTransaction) ID() string {
	// Serialize the transaction (without signatures)
	rawTx, err := tx.Serialize()
	if err != nil {
		return ""
	}

	// Compute SHA256 hash
	hash := sha256.Sum256(rawTx)

	// Return first 20 bytes as hex string (40 characters)
	// This matches Steem's transaction ID format
	return hex.EncodeToString(hash[:20])
}
