package operation

type ShutdownWitnessOperation struct {
}

func (op *ShutdownWitnessOperation) Type() OpType {
	return TypeShutdownWitness
}
func (op *ShutdownWitnessOperation) Data() interface{} {
	return op
}
