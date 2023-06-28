package operation

// FC_REFLECT( steemit::chain::account_create_operation,
//             (fee)
//             (creator)
//             (new_account_name)
//             (owner)
//             (active)
//             (posting)
//             (memo_key)
//             (json_metadata) )

type AccountCreateOperation struct {
	Fee            string     `json:"fee"`
	Creator        string     `json:"creator"`
	NewAccountName string     `json:"new_account_name"`
	Owner          *Authority `json:"owner"`
	Active         *Authority `json:"active"`
	Posting        *Authority `json:"posting"`
	MemoKey        string     `json:"memo_key"`
	JsonMetadata   string     `json:"json_metadata"`
}

func (op *AccountCreateOperation) Type() OpType {
	return TypeAccountCreate
}

func (op *AccountCreateOperation) Data() interface{} {
	return op
}
