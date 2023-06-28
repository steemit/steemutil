package operation

type WitnessUpdateOperation struct {
}

func (op *WitnessUpdateOperation) Type() OpType {
	return TypeWitnessUpdate
}

func (op *WitnessUpdateOperation) Data() interface{} {
	return op
}
