package operation

type CommentBenefactorRewardOperation struct {
}

func (op *CommentBenefactorRewardOperation) Type() OpType {
	return TypeCommentBenefactorReward
}
func (op *CommentBenefactorRewardOperation) Data() interface{} {
	return op
}
