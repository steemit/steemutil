package operation

type CommentRewardOperation struct {
}

func (op *CommentRewardOperation) Type() OpType {
	return TypeCommentReward
}
func (op *CommentRewardOperation) Data() interface{} {
	return op
}
