package operation

type ClaimRewardBalanceOperation struct {
}

func (op *ClaimRewardBalanceOperation) Type() OpType {
	return TypeClaimRewardBalance
}
func (op *ClaimRewardBalanceOperation) Data() interface{} {
	return op
}
