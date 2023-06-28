package operation

type SetWithdrawVestingRouteOperation struct {
}

func (op *SetWithdrawVestingRouteOperation) Type() OpType {
	return TypeSetWithdrawVestingRoute
}
func (op *SetWithdrawVestingRouteOperation) Data() interface{} {
	return op
}
