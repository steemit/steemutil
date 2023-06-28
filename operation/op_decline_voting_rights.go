package operation

type DeclineVotingRightsOperation struct {
}

func (op *DeclineVotingRightsOperation) Type() OpType {
	return TypeDeclineVotingRights
}
func (op *DeclineVotingRightsOperation) Data() interface{} {
	return op
}
