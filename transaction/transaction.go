package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/protocol"

	"github.com/pkg/errors"
)

// Transaction represents a blockchain transaction.
type Transaction struct {
	RefBlockNum    protocol.UInt16     `json:"ref_block_num"`
	RefBlockPrefix protocol.UInt32     `json:"ref_block_prefix"`
	Expiration     *protocol.Time      `json:"expiration"`
	Operations     protocol.Operations `json:"operations"`
	Signatures     []string            `json:"signatures"`
}

// MarshalTransaction implements transaction.Marshaller interface.
func (tx *Transaction) MarshalTransaction(encoderObj *encoder.Encoder) error {
	if len(tx.Operations) == 0 {
		return errors.New("no operation specified")
	}

	enc := encoder.NewRollingEncoder(encoderObj)

	enc.Encode(tx.RefBlockNum)
	enc.Encode(tx.RefBlockPrefix)
	enc.Encode(tx.Expiration)

	enc.EncodeUVarint(uint64(len(tx.Operations)))
	for _, op := range tx.Operations {
		enc.Encode(op)
	}

	// Extensions are not supported yet.
	enc.EncodeUVarint(0)

	return enc.Err()
}

// PushOperation can be used to add an operation into the transaction.
func (tx *Transaction) PushOperation(op protocol.Operation) {
	tx.Operations = append(tx.Operations, op)
}

func RefBlockNum(blockNumber protocol.UInt32) protocol.UInt16 {
	return protocol.UInt16(blockNumber)
}

func RefBlockPrefix(blockID string) (protocol.UInt32, error) {
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
	return protocol.UInt32(prefix), nil
}
