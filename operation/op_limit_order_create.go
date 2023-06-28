package operation

import (
	"github.com/steemit/steemutil/util"
)

// FC_REFLECT( steemit::chain::limit_order_create_operation,
//             (owner)
//             (orderid)
//             (amount_to_sell)
//             (min_to_receive)
//             (fill_or_kill)
//             (expiration) )

type LimitOrderCreateOperation struct {
	Owner        string     `json:"owner"`
	OrderID      uint32     `json:"orderid"`
	AmountToSell string     `json:"amount_to_sell"`
	MinToReceive string     `json:"min_to_receive"`
	FillOrKill   bool       `json:"fill_or_kill"`
	Expiration   *util.Time `json:"expiration"`
}

func (op *LimitOrderCreateOperation) Type() OpType {
	return TypeLimitOrderCreate
}

func (op *LimitOrderCreateOperation) Data() interface{} {
	return op
}
