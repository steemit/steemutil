package operation

import (
	"github.com/steemit/steemutil/encoder"
)

// Delegate vesting shares from one account to the other. The vesting shares are still owned
// by the original account, but content voting rights and bandwidth allocation are transferred
// to the receiving account. This sets the delegation to `vesting_shares`, increasing it or
// decreasing it as needed. (i.e. a delegation of 0 removes the delegation)
//
// When a delegation is removed the shares are placed in limbo for a week to prevent a satoshi
// of VESTS from voting on the same content twice.
//
// FC_REFLECT(
//		steem::protocol::delegate_vesting_shares_operation,
//		(delegator)(delegatee)(vesting_shares) );
type DelegateVestingSharesOperation struct {
	Delegator     string `json:"delegator"`
	Delegatee     string `json:"delegatee"`
	VestingShares Asset  `json:"vesting_shares"`
}

func (op *DelegateVestingSharesOperation) Type() OpType {
	return TypeDelegateVestingShares
}
func (op *DelegateVestingSharesOperation) Data() interface{} {
	return op
}
func (op *DelegateVestingSharesOperation) MarshalTransaction(encoderObj *encoder.Encoder) (err error) {
	encoderObj.Encode(op.Delegator)
	encoderObj.Encode(op.Delegatee)
	encoderObj.Encode(op.VestingShares)
	return
}
