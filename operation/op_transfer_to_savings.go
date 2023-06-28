package operation

type TransferToSavingsOperation struct {
}

func (op *TransferToSavingsOperation) Type() OpType {
	return TypeTransferToSavings
}
func (op *TransferToSavingsOperation) Data() interface{} {
	return op
}
