package operation

type ClaimAccountOperation struct {
}

func (op *ClaimAccountOperation) Type() OpType {
	return TypeClaimAccount
}
func (op *ClaimAccountOperation) Data() interface{} {
	return op
}