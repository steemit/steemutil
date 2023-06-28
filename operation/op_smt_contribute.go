package operation

type SmtContributeOperation struct {
}

func (op *SmtContributeOperation) Type() OpType {
	return TypeSmtContribute
}
func (op *SmtContributeOperation) Data() interface{} {
	return op
}
