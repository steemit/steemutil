package operation

type ProducerRewardOperation struct {
}

func (op *ProducerRewardOperation) Type() OpType {
	return TypeProducerReward
}
func (op *ProducerRewardOperation) Data() interface{} {
	return op
}
