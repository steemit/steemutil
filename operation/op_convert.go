package operation

// FC_REFLECT( steemit::chain::convert_operation,
//             (owner)
//             (requestid)
//             (amount) )

type ConvertOperation struct {
	Owner     string `json:"owner"`
	RequestID uint32 `json:"requestid"`
	Amount    string `json:"amount"`
}

func (op *ConvertOperation) Type() OpType {
	return TypeConvert
}

func (op *ConvertOperation) Data() interface{} {
	return op
}
