package operation

type FillVestingWithdrawOperation struct {
}

func (op *FillVestingWithdrawOperation) Type() OpType {
	return TypeFillVestingWithdraw
}
func (op *FillVestingWithdrawOperation) Data() interface{} {
	return op
}
