package operation

type CustomBinaryOperation struct {
}

func (op *CustomBinaryOperation) Type() OpType {
	return TypeCustomBinary
}
func (op *CustomBinaryOperation) Data() interface{} {
	return op
}
