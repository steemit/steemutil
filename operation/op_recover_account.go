package operation

type RecoverAccountOperation struct {
}

func (op *RecoverAccountOperation) Type() OpType {
	return TypeRecoverAccount
}
func (op *RecoverAccountOperation) Data() interface{} {
	return op
}
