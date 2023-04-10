package api

import "github.com/steemit/steemutil/protocol"

type RpcSendData struct {
	Id      protocol.UInt `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []any         `json:"params"`
}

type RpcResultData struct {
	Id      protocol.UInt `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Result  any           `json:"result,omitempty"`
	Error   any           `json:"error,omitempty"`
}
