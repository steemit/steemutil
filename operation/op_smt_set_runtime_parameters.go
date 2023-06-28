package operation

type SmtSetRuntimeParaetersOperation struct {
}

func (op *SmtSetRuntimeParaetersOperation) Type() OpType {
	return TypeSmtSetRuntimeParameters
}
func (op *SmtSetRuntimeParaetersOperation) Data() interface{} {
	return op
}
