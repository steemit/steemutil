package operation

type WitnessSetPropertiesOperation struct {
}

func (op *WitnessSetPropertiesOperation) Type() OpType {
	return TypeWitnessSetProperties
}
func (op *WitnessSetPropertiesOperation) Data() interface{} {
	return op
}
