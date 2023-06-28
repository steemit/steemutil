package operation

type SmtSetSetupParametersOperation struct {
}

func (op *SmtSetSetupParametersOperation) Type() OpType {
	return TypeSmtSetSetupParameters
}
func (op *SmtSetSetupParametersOperation) Data() interface{} {
	return op
}
