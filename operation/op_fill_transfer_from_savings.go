package operation

type FillTransferFromSavingsOperation struct {
}

func (op *FillTransferFromSavingsOperation) Type() OpType {
	return TypeFillTransferFromSavings
}
func (op *FillTransferFromSavingsOperation) Data() interface{} {
	return op
}
