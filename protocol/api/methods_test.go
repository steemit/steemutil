package api

import (
	"testing"
)

func TestMethodsDataNotEmpty(t *testing.T) {
	if len(MethodsData) == 0 {
		t.Error("MethodsData should not be empty")
	}
}

func TestMethodsDataStructure(t *testing.T) {
	for i, method := range MethodsData {
		if method.API == "" {
			t.Errorf("Method %d: API field is empty", i)
		}
		if method.Method == "" {
			t.Errorf("Method %d: Method field is empty", i)
		}
	}
}

func TestMethodsDataHasCommonAPIs(t *testing.T) {
	apis := make(map[string]bool)
	for _, method := range MethodsData {
		apis[method.API] = true
	}

	expectedAPIs := []string{
		"database_api",
		"network_broadcast_api",
		"follow_api",
		"market_history_api",
		"condenser_api",
	}

	for _, api := range expectedAPIs {
		if !apis[api] {
			t.Errorf("MethodsData missing API: %s", api)
		}
	}
}

func TestMethodsDataHasCommonMethods(t *testing.T) {
	methods := make(map[string]bool)
	for _, method := range MethodsData {
		methods[method.Method] = true
	}

	expectedMethods := []string{
		"get_block",
		"get_dynamic_global_properties",
		"broadcast_transaction_synchronous",
		"get_content",
		"get_followers",
	}

	for _, method := range expectedMethods {
		if !methods[method] {
			t.Errorf("MethodsData missing method: %s", method)
		}
	}
}

func TestAPIMethodStructure(t *testing.T) {
	method := APIMethod{
		API:    "database_api",
		Method: "get_block",
		Params: []string{"blockNum"},
	}

	if method.API != "database_api" {
		t.Error("APIMethod.API not set correctly")
	}
	if method.Method != "get_block" {
		t.Error("APIMethod.Method not set correctly")
	}
	if len(method.Params) != 1 || method.Params[0] != "blockNum" {
		t.Error("APIMethod.Params not set correctly")
	}
}

