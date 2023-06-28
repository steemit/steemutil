package operation

type RequestAccountRecoveryOperation struct {
}

func (op *RequestAccountRecoveryOperation) Type() OpType {
	return TypeRequestAccountRecovery
}
func (op *RequestAccountRecoveryOperation) Data() interface{} {
	return op
}
