package jsonrpc2

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/steemit/steemutil/protocol/api"
)

func TestBuildSendData(t *testing.T) {
	client := NewClient("https://api.steemit.com")
	client.BuildSendData("condenser_api.get_block", []any{1})
	gotData := client.SendData
	expectedData := `{"id":1,"jsonrpc":"2.0","method":"condenser_api.get_block","params":[1]}`

	if string(gotData) != expectedData {
		t.Errorf("expected %v, got %v", expectedData, gotData)
	}
}

func TestSend(t *testing.T) {
	testData := `{"jsonrpc":"2.0","result":{"previous":"0000000000000000000000000000000000000000","timestamp":"2016-03-24T16:05:00","witness":"initminer","transaction_merkle_root":"0000000000000000000000000000000000000000","extensions":[],"witness_signature":"204f8ad56a8f5cf722a02b035a61b500aa59b9519b2c33c77a80c0a714680a5a5a7a340d909d19996613c5e4ae92146b9add8a7a663eef37d837ef881477313043","transactions":[],"block_id":"0000000109833ce528d5bbfb3f6225b39ee10086","signing_key":"STM8GC13uCZbP44HzMLV6zPZGwVQ8Nt4Kji8PapsPiNq1BK153XTX","transaction_ids":[]},"id":1}`
	testDataBody := &api.RpcResultData{}
	err := json.Unmarshal([]byte(testData), testDataBody)
	if err != nil {
		t.Errorf("json unmarshal error: %v", err)
		return
	}
	expectedData, err := json.Marshal(testDataBody.Result)
	if err != nil {
		t.Errorf("json marshal error: %v", err)
		return
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.steemit.com",
		func(req *http.Request) (res *http.Response, err error) {
			tmpReq := &api.RpcSendData{}
			err = json.NewDecoder(req.Body).Decode(tmpReq)
			if err != nil {
				return
			}
			if tmpReq.Method != "condenser_api.get_block" {
				t.Errorf("unexpected method")
				return
			}
			return httpmock.NewJsonResponse(200, testDataBody)
		},
	)
	client := NewClient("https://api.steemit.com")
	err = client.BuildSendData("condenser_api.get_block", []any{20000000})
	if err != nil {
		t.Errorf("build send data error: %v", err)
		return
	}
	result, err := client.Send()
	if err != nil {
		t.Errorf("send error: %v", err)
		return
	}

	gotData, err := json.Marshal(result.Result)
	if err != nil {
		t.Errorf("json marshal error: %v", err)
		return
	}

	if string(expectedData) != string(gotData) {
		t.Errorf("expected data: %v, got: %v", string(expectedData), string(gotData))
	}
}
