package steemutil

import (
	"testing"
)

func TestRequestDataStructure(t *testing.T) {
	req := RequestData{
		Id:      1,
		JsonRPC: "2.0",
		Method:  "test_method",
		Params:  []any{"param1", "param2"},
	}

	if req.Id != 1 {
		t.Error("RequestData.Id not set correctly")
	}
	if req.JsonRPC != "2.0" {
		t.Error("RequestData.JsonRPC not set correctly")
	}
	if req.Method != "test_method" {
		t.Error("RequestData.Method not set correctly")
	}
	if len(req.Params) != 2 {
		t.Error("RequestData.Params not set correctly")
	}
}

func TestResponseDataStructure(t *testing.T) {
	res := ResponseData{
		Id:      1,
		JsonRPC: "2.0",
		Result:  "test_result",
	}

	if res.Id != 1 {
		t.Error("ResponseData.Id not set correctly")
	}
	if res.JsonRPC != "2.0" {
		t.Error("ResponseData.JsonRPC not set correctly")
	}
	if res.Result != "test_result" {
		t.Error("ResponseData.Result not set correctly")
	}
}

func TestGetClient(t *testing.T) {
	api := "https://api.steemit.com"
	timeout := uint(30)

	client := GetClient(api, timeout)
	if client == nil {
		t.Error("GetClient returned nil")
	}

	// Test that client implements IClient interface
	var _ IClient = client
}
