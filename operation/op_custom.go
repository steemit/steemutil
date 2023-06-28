package operation

type CustomOperation struct {
}

func (op *CustomOperation) Type() OpType {
	return TypeCustom
}

func (op *CustomOperation) Data() interface{} {
	return op
}
