package operation

type FillConvertRequestOperation struct {
}

func (op *FillConvertRequestOperation) Type() OpType {
	return TypeFillConvertRequest
}
func (op *FillConvertRequestOperation) Data() interface{} {
	return op
}
