package operation

type ResetAccountOperation struct {
}

func (op *ResetAccountOperation) Type() OpType {
	return TypeResetAccount
}
func (op *ResetAccountOperation) Data() interface{} {
	return op
}
