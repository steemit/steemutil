package operation

type SmtSetupOperation struct {
}

func (op *SmtSetupOperation) Type() OpType {
	return TypeSmtSetup
}
func (op *SmtSetupOperation) Data() interface{} {
	return op
}
