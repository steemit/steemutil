package operation

type EescrowReleaseOperation struct {
}

func (op *EescrowReleaseOperation) Type() OpType {
	return TypeEscrowRelease
}
func (op *EescrowReleaseOperation) Data() interface{} {
	return op
}
