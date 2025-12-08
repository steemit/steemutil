package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/steemit/steemutil/protocol/api"
)

type IJsonRpc interface {
	Send() (*api.RpcResultData, error)
}

type JsonRpc struct {
	Url      string
	SendData []byte
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

func (j *JsonRpc) Send() (result *api.RpcResultData, err error) {
	bodyReader := bytes.NewReader(j.SendData)
	req, err := http.NewRequest(http.MethodPost, j.Url, bodyReader)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return result, errors.Errorf("failed to response(http code): %v", res.StatusCode)
	}
	result = &api.RpcResultData{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

func NewClient(url string) *JsonRpc {
	client := &JsonRpc{
		Url: url,
	}
	return client
}
