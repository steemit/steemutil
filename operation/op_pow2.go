package operation

type POW2Operation struct {
}

func (op *POW2Operation) Type() OpType {
	return TypePOW2
}
func (op *POW2Operation) Data() interface{} {
	return op
}
