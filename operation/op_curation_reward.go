package operation

type CurationRewardOperation struct {
}

func (op *CurationRewardOperation) Type() OpType {
	return TypeCurationReward
}
func (op *CurationRewardOperation) Data() interface{} {
	return op
}
