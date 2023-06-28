package operation

type AccountCreateWithDelegationOperation struct {
}

func (op *AccountCreateWithDelegationOperation) Type() OpType {
	return TypeAccountCreateWithDelegation
}
func (op *AccountCreateWithDelegationOperation) Data() interface{} {
	return op
}
