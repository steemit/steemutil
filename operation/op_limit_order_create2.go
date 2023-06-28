package operation

type LimitOrderCreate2Operation struct {
}

func (op *LimitOrderCreate2Operation) Type() OpType {
	return TypeLimitOrderCreate2
}
func (op *LimitOrderCreate2Operation) Data() interface{} {
	return op
}