package operation

type TransferFromSavingsOperation struct {
}

func (op *TransferFromSavingsOperation) Type() OpType {
	return TypeTransferFromSavings
}
func (op *TransferFromSavingsOperation) Data() interface{} {
	return op
}
