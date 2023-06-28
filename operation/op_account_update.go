package operation

import (
	"github.com/steemit/steemutil/encoder"
)

// FC_REFLECT( steemit::chain::account_update_operation,
//             (account)
//             (owner)
//             (active)
//             (posting)
//             (memo_key)
//             (json_metadata) )

type AccountUpdateOperation struct {
	Account      string     `json:"account"`
	Owner        *Authority `json:"owner"`
	Active       *Authority `json:"active"`
	Posting      *Authority `json:"posting"`
	MemoKey      string     `json:"memo_key"`
	JsonMetadata string     `json:"json_metadata"`
}

func (op *AccountUpdateOperation) Type() OpType {
	return TypeAccountUpdate
}

func (op *AccountUpdateOperation) Data() interface{} {
	return op
}

func (op *AccountUpdateOperation) MarshalTransaction(encoderObj *encoder.Encoder) (err error) {
	if err = encoderObj.EncodeUVarint(uint64(TypeAccountUpdate.Code())); err != nil {
		return err
	}
	if err = encoderObj.Encode(op.Account); err != nil {
		return err
	}
	if err = encoderObj.Encode(op.Owner); err != nil {
		return err
	}
	if err = encoderObj.Encode(op.Active); err != nil {
		return err
	}
	if err = encoderObj.Encode(op.Posting); err != nil {
		return err
	}
	if err = encoderObj.Encode(op.MemoKey); err != nil {
		return err
	}
	if err = encoderObj.Encode(op.JsonMetadata); err != nil {
		return err
	}
	return err
}
