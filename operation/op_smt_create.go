package operation

type SmtCreateOperation struct {
}

func (op *SmtCreateOperation) Type() OpType {
	return TypeSmtCreate
}
func (op *SmtCreateOperation) Data() interface{} {
	return op
}
