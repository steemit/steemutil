package protocol

import (
	"encoding/json"
)

type OperationObject struct {
	BlockNumber            uint32    `json:"block"`
	TransactionID          string    `json:"trx_id"`
	TransactionInBlock     uint32    `json:"trx_in_block"`
	Operation              Operation `json:"op"`
	OperationInTransaction uint16    `json:"op_in_trx"`
	VirtualOperation       uint64    `json:"virtual_op"`
	Timestamp              *Time     `json:"timestamp"`
}

type rawOperationObject struct {
	BlockNumber            uint32          `json:"block"`
	TransactionID          string          `json:"trx_id"`
	TransactionInBlock     uint32          `json:"trx_in_block"`
	Operation              *operationTuple `json:"op"`
	OperationInTransaction uint16          `json:"op_in_trx"`
	VirtualOperation       uint64          `json:"virtual_op"`
	Timestamp              *Time           `json:"timestamp"`
}

func (op *OperationObject) UnmarshalJSON(p []byte) error {
	var raw rawOperationObject
	if err := json.Unmarshal(p, &raw); err != nil {
		return err
	}

	op.BlockNumber = raw.BlockNumber
	op.TransactionID = raw.TransactionID
	op.TransactionInBlock = raw.TransactionInBlock
	op.Operation = raw.Operation.Data
	op.OperationInTransaction = raw.OperationInTransaction
	op.VirtualOperation = raw.VirtualOperation
	op.Timestamp = raw.Timestamp
	return nil
}

func (op *OperationObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&rawOperationObject{
		BlockNumber:            op.BlockNumber,
		TransactionID:          op.TransactionID,
		TransactionInBlock:     op.TransactionInBlock,
		Operation:              &operationTuple{op.Operation.Type(), op.Operation},
		OperationInTransaction: op.OperationInTransaction,
		VirtualOperation:       op.VirtualOperation,
		Timestamp:              op.Timestamp,
	})
}
