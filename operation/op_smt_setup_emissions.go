package operation

type SmtSetupEmissionsOperation struct {
}

func (op *SmtSetupEmissionsOperation) Type() OpType {
	return TypeSmtSetupEmissions
}
func (op *SmtSetupEmissionsOperation) Data() interface{} {
	return op
}
