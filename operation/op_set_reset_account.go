package operation

type SetResetAccountOperation struct {
}

func (op *SetResetAccountOperation) Type() OpType {
	return TypeSetResetAccount
}
func (op *SetResetAccountOperation) Data() interface{} {
	return op
}
