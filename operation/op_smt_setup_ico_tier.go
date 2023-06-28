package operation

type SmtSetupIcoTierOperation struct {
}

func (op *SmtSetupIcoTierOperation) Type() OpType {
	return TypeSmtSetupIcoTier
}
func (op *SmtSetupIcoTierOperation) Data() interface{} {
	return op
}
