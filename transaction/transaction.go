package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/steemit/steemutil/chain"
	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/operation"
	"github.com/steemit/steemutil/util"
	"github.com/steemit/steemutil/wif"

	"github.com/pkg/errors"
)

// Transaction represents a blockchain transaction.
type Transaction struct {
	RefBlockNum    util.UInt16          `json:"ref_block_num"`
	RefBlockPrefix util.UInt32          `json:"ref_block_prefix"`
	Expiration     *util.Time           `json:"expiration"`
	Operations     operation.Operations `json:"operations"`
	Signatures     []string             `json:"signatures"`
}

// MarshalTransaction implements transaction.Marshaller interface.
func (tx *Transaction) MarshalTransaction(encoderObj *encoder.Encoder) (err error) {
	if len(tx.Operations) == 0 {
		return errors.New("no operation specified")
	}

	encoderObj.Encode(tx.RefBlockNum)
	encoderObj.Encode(tx.RefBlockPrefix)
	encoderObj.Encode(tx.Expiration)

	encoderObj.EncodeUVarint(uint64(len(tx.Operations)))
	for _, op := range tx.Operations {
		encoderObj.Encode(op)
	}

	// TODO: Extensions are not supported yet.
	encoderObj.EncodeUVarint(0)

	return err
}

// PushOperation can be used to add an operation into the transaction.
func (tx *Transaction) PushOperation(op operation.IOperation) {
	tx.Operations = append(tx.Operations, op)
}

func RefBlockNum(blockNumber util.UInt32) util.UInt16 {
	return util.UInt16(blockNumber)
}

func RefBlockPrefix(blockID string) (util.UInt32, error) {
	// Block ID is hex-encoded.
	rawBlockID, err := hex.DecodeString(blockID)
	if err != nil {
		return 0, errors.Wrapf(err, "networkbroadcast: failed to decode block ID: %v", blockID)
	}

	// Raw prefix = raw block ID [4:8].
	// Make sure we don't trigger a slice bounds out of range panic.
	if len(rawBlockID) < 8 {
		return 0, errors.Errorf("networkbroadcast: invalid block ID: %v", blockID)
	}
	rawPrefix := rawBlockID[4:8]

	// Decode the prefix.
	var prefix uint32
	if err := binary.Read(bytes.NewReader(rawPrefix), binary.LittleEndian, &prefix); err != nil {
		return 0, errors.Wrapf(err, "networkbroadcast: failed to read block prefix: %v", rawPrefix)
	}

	// Done, return the prefix.
	return util.UInt32(prefix), nil
}

func (tx *Transaction) Serialize() ([]byte, error) {
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)

	if err := encoderObj.Encode(tx); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (tx *Transaction) Digest(chain *chain.Chain) ([]byte, error) {
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

func (tx *Transaction) Sign(privKeys []*wif.PrivateKey, chain *chain.Chain) error {
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

	tx.Signatures = sigsHex
	return nil
}

func (tx *Transaction) Verify(pubKeys []*wif.PublicKey, chain *chain.Chain) (bool, error) {
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
