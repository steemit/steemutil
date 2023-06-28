package operation

type Vote2Operation struct {
}

func (op *Vote2Operation) Type() OpType {
	return TypeVote2
}
func (op *Vote2Operation) Data() interface{} {
	return op
}
