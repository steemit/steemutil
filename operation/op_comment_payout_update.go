package operation

type CommentPayoutUpdateOperation struct {
}

func (op *CommentPayoutUpdateOperation) Type() OpType {
	return TypeCommentPayoutUpdate
}
func (op *CommentPayoutUpdateOperation) Data() interface{} {
	return op
}
