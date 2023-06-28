package operation

type RemoveProposalOperation struct {
}

func (op *RemoveProposalOperation) Type() OpType {
	return TypeRemoveProposal
}
func (op *RemoveProposalOperation) Data() interface{} {
	return op
}
