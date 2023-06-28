package operation

type ProposalPayOperation struct {
}

func (op *ProposalPayOperation) Type() OpType {
	return TypeProposalPay
}
func (op *ProposalPayOperation) Data() interface{} {
	return op
}
