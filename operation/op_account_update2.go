package operation

type AccountUpdate2Operation struct {
}

func (op *AccountUpdate2Operation) Type() OpType {
	return TypeAccountUpdate2
}
func (op *AccountUpdate2Operation) Data() interface{} {
	return op
}
