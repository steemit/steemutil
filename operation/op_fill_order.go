package operation

type FillOrderOperation struct {
}

func (op *FillOrderOperation) Type() OpType {
	return TypeFillOrder
}
func (op *FillOrderOperation) Data() interface{} {
	return op
}
