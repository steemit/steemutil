package operation

type ClaimRewardBalance2Operation struct {
}

func (op *ClaimRewardBalance2Operation) Type() OpType {
	return TypeClaimRewardBalance2
}
func (op *ClaimRewardBalance2Operation) Data() interface{} {
	return op
}
