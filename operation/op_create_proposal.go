package operation

type CreateProposalOperation struct {
}

func (op *CreateProposalOperation) Type() OpType {
	return TypeCreateProposal
}
func (op *CreateProposalOperation) Data() interface{} {
	return op
}
