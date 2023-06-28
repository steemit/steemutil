package operation

type LiquidityRewardOperation struct {
}

func (op *LiquidityRewardOperation) Type() OpType {
	return TypeLiquidityReward
}
func (op *LiquidityRewardOperation) Data() interface{} {
	return op
}
