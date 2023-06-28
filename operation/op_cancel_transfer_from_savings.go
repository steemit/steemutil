package operation

type CancelTransferFromSavingsOperation struct {
}

func (op *CancelTransferFromSavingsOperation) Type() OpType {
	return TypeCancelTransferFromSavings
}
func (op *CancelTransferFromSavingsOperation) Data() interface{} {
	return op
}
