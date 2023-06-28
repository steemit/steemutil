package operation

import (
	"github.com/steemit/steemutil/util"
)

type Authority struct {
	AccountAuths    util.StringInt64Map `json:"account_auths"`
	KeyAuths        util.StringInt64Map `json:"key_auths"`
	WeightThreshold uint32              `json:"weight_threshold"`
}
