package operation

type EscrowApproveOperation struct {
}

func (op *EscrowApproveOperation) Type() OpType {
	return TypeEscrowApprove
}
func (op *EscrowApproveOperation) Data() interface{} {
	return op
}
