package operation

type EscrowDisputeOperation struct {
}

func (op *EscrowDisputeOperation) Type() OpType {
	return TypeEscrowDispute
}
func (op *EscrowDisputeOperation) Data() interface{} {
	return op
}
