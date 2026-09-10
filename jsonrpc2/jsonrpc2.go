package jsonrpc2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/steemit/steemutil/protocol/api"
)

type IJsonRpc interface {
	Send() (*api.RpcResultData, error)
}

type JsonRpc struct {
	Url      string
	SendData []byte

	// client is the http.Client used by SendWithContext. A nil client means
	// "legacy behavior": a client with a 30s timeout per call. Set it via
	// WithHTTPClient / NewClientWithOptions.
	client *http.Client
}

func (j *JsonRpc) BuildSendData(method string, params []any) (err error) {
	data := &api.RpcSendData{
		Id:      1,
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
	tmp, err := json.Marshal(data)
	if err != nil {
		return
	}
	j.SendData = tmp

	// Debug: Print JSON request if DEBUG is set
	if os.Getenv("DEBUG") != "" && method == "condenser_api.broadcast_transaction_synchronous" {
		fmt.Printf("=== JSON-RPC Request ===\n%s\n", string(tmp))
	}

	return
}

func NewClient(url string) *JsonRpc {
	return NewClientWithOptions(url)
}
