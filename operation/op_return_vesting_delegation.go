package operation

type ReturnVestingDelegationOperation struct {
}

func (op *ReturnVestingDelegationOperation) Type() OpType {
	return TypeReturnVestingDelegation
}
func (op *ReturnVestingDelegationOperation) Data() interface{} {
	return op
}
