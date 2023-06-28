package operation

type SpsFundOperation struct {
}

func (op *SpsFundOperation) Type() OpType {
	return TypeSpsFund
}
func (op *SpsFundOperation) Data() interface{} {
	return op
}