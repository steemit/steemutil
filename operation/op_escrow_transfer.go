package operation

type EscrowTransferOperation struct {
}

func (op *EscrowTransferOperation) Type() OpType {
	return TypeEscrowTransfer
}
func (op *EscrowTransferOperation) Data() interface{} {
	return op
}
