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

type SignedTransaction struct {
	*Transaction
}

func NewSignedTransaction(tx *Transaction) *SignedTransaction {
	if tx.Expiration == nil {
		expiration := time.Now().Add(600 * time.Second)
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

func (tx *SignedTransaction) Verify(pubKeys []*wif.PublicKey, chain *Chain) (bool, error) {
	// Compute digest
	digest, err := tx.Digest(chain)
	if err != nil {
		return false, err
	}

	// Parse signatures
	sigs := make([][]byte, 0, len(tx.Signatures))
	for i, sig := range sigs {
		tmpSig, err := ecdsa.ParseSignature(sig)
		if err != nil {
			return false, err
		}
		verified := tmpSig.Verify(digest, pubKeys[i].Raw)
		if !verified {
			return false, nil
		}
	}
	return true, nil
}
