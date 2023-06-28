package operation

type AuthorRewardOperation struct {
}

func (op *AuthorRewardOperation) Type() OpType {
	return TypeAuthorReward
}
func (op *AuthorRewardOperation) Data() interface{} {
	return op
}
