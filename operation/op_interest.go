package operation

type InterestOperation struct {
}

func (op *InterestOperation) Type() OpType {
	return TypeInterest
}
func (op *InterestOperation) Data() interface{} {
	return op
}
