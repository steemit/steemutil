package operation

type HardforkOperation struct {
}

func (op *HardforkOperation) Type() OpType {
	return TypeHardfork
}
func (op *HardforkOperation) Data() interface{} {
	return op
}
