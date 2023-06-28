package api

import (
	"github.com/steemit/steemutil/util"
)

type RpcSendData struct {
	Id      util.UInt `json:"id"`
	JsonRpc string    `json:"jsonrpc"`
	Method  string    `json:"method"`
	Params  []any     `json:"params"`
}

type RpcResultData struct {
	Id      util.UInt `json:"id"`
	JsonRpc string    `json:"jsonrpc"`
	Result  any       `json:"result,omitempty"`
	Error   any       `json:"error,omitempty"`
}
