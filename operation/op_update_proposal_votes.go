package operation

type UpdateProposalVotesOperation struct {
}

func (op *UpdateProposalVotesOperation) Type() OpType {
	return TypeUpdateProposalVotes
}
func (op *UpdateProposalVotesOperation) Data() interface{} {
	return op
}
