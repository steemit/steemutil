package operation

type ChangeRecoverAccountOperation struct {
}

func (op *ChangeRecoverAccountOperation) Type() OpType {
	return TypeChangeRecoveryAccount
}
func (op *ChangeRecoverAccountOperation) Data() interface{} {
	return op
}
