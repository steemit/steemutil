package api

type RpcSendData struct {
	Id      uint   `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type RpcResultData struct {
	Id      uint   `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}
