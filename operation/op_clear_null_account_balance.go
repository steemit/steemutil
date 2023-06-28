package operation

type ClearNullAccountBalanceOperation struct {
}

func (op *ClearNullAccountBalanceOperation) Type() OpType {
	return TypeClearNullAccountBalance
}
func (op *ClearNullAccountBalanceOperation) Data() interface{} {
	return op
}
