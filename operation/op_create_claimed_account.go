package operation

type CreateClaimedAccountOperation struct {
}

func (op *CreateClaimedAccountOperation) Type() OpType {
	return TypeCreateClaimedAccount
}
func (op *CreateClaimedAccountOperation) Data() interface{} {
	return op
}
