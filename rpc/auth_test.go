package rpc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"
)

// Test constants
const (
	testAccount    = "testuser"
	testPrivateKey = "5JLw5dgQAx6rhZEgNN5C2ds1V47RweGshynFSWFbaMohsYsBvE8"
	testMethod     = "condenser_api.get_accounts"
)

var testParams = []interface{}{[]string{"testuser"}}

func TestSign(t *testing.T) {
	request := &RpcRequest{
		Method: testMethod,
		Params: testParams,
		ID:     1,
	}

	signedRequest, err := Sign(request, testAccount, []string{testPrivateKey})
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	// Verify basic structure
	if signedRequest.JsonRpc != "2.0" {
		t.Errorf("Expected jsonrpc '2.0', got '%s'", signedRequest.JsonRpc)
	}

	if signedRequest.Method != testMethod {
		t.Errorf("Expected method '%s', got '%s'", testMethod, signedRequest.Method)
	}

	if signedRequest.ID != 1 {
		t.Errorf("Expected ID 1, got %d", signedRequest.ID)
	}

	signed := signedRequest.Params.Signed

	// Verify signed params structure
	if signed.Account != testAccount {
		t.Errorf("Expected account '%s', got '%s'", testAccount, signed.Account)
	}

	if len(signed.Nonce) != 16 { // 8 bytes = 16 hex chars
		t.Errorf("Expected nonce length 16, got %d", len(signed.Nonce))
	}

	if len(signed.Signatures) != 1 {
		t.Errorf("Expected 1 signature, got %d", len(signed.Signatures))
	}

	// Verify timestamp format
	if _, err := time.Parse(time.RFC3339Nano, signed.Timestamp); err != nil {
		t.Errorf("Invalid timestamp format: %v", err)
	}

	// Verify params encoding
	if signed.Params == "" {
		t.Error("Params should not be empty")
	}
}

func TestSignWithoutParams(t *testing.T) {
	request := &RpcRequest{
		Method: testMethod,
		Params: nil,
		ID:     1,
	}

	_, err := Sign(request, testAccount, []string{testPrivateKey})
	if err == nil {
		t.Error("Expected error when signing request without params")
	}
}

func TestSignWithInvalidPrivateKey(t *testing.T) {
	request := &RpcRequest{
		Method: testMethod,
		Params: testParams,
		ID:     1,
	}

	_, err := Sign(request, testAccount, []string{"invalid-key"})
	if err == nil {
		t.Error("Expected error when signing with invalid private key")
	}
}

func TestValidate(t *testing.T) {
	// First create a signed request
	request := &RpcRequest{
		Method: testMethod,
		Params: testParams,
		ID:     1,
	}

	signedRequest, err := Sign(request, testAccount, []string{testPrivateKey})
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	// Now validate it
	params, err := Validate(signedRequest, DefaultVerifyFunc)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	// Verify that we got back the original params
	if len(params) != len(testParams) {
		t.Errorf("Expected %d params, got %d", len(testParams), len(params))
	}

	// Compare the first param (should be a slice of strings)
	originalParam := testParams[0].([]string)
	recoveredParam := params[0].([]interface{})

	if len(recoveredParam) != len(originalParam) {
		t.Errorf("Expected param length %d, got %d", len(originalParam), len(recoveredParam))
	}

	if recoveredParam[0].(string) != originalParam[0] {
		t.Errorf("Expected param value '%s', got '%s'", originalParam[0], recoveredParam[0])
	}
}

func TestValidateInvalidJsonRpc(t *testing.T) {
	signedRequest := &SignedRequest{
		JsonRpc: "1.0", // Invalid version
		Method:  testMethod,
		ID:      1,
	}

	_, err := Validate(signedRequest, DefaultVerifyFunc)
	if err == nil {
		t.Error("Expected error for invalid JSON-RPC version")
	}
}

func TestValidateExpiredSignature(t *testing.T) {
	// Create a signed request with an old timestamp
	request := &RpcRequest{
		Method: testMethod,
		Params: testParams,
		ID:     1,
	}

	signedRequest, err := Sign(request, testAccount, []string{testPrivateKey})
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	// Modify the timestamp to be older than 60 seconds
	oldTime := time.Now().UTC().Add(-70 * time.Second)
	signedRequest.Params.Signed.Timestamp = oldTime.Format(time.RFC3339Nano)

	_, err = Validate(signedRequest, DefaultVerifyFunc)
	if err == nil {
		t.Error("Expected error for expired signature")
	}
}

func TestHashMessage(t *testing.T) {
	timestamp := "2023-11-27T10:30:00.000Z"
	account := "testuser"
	method := "condenser_api.get_accounts"
	params := "W1sidGVzdHVzZXIiXV0=" // base64 encoded [["testuser"]]
	nonce := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	hash := hashMessage(timestamp, account, method, params, nonce)

	// Verify hash length
	if len(hash) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hash))
	}

	// Verify hash is deterministic
	hash2 := hashMessage(timestamp, account, method, params, nonce)
	if hex.EncodeToString(hash) != hex.EncodeToString(hash2) {
		t.Error("Hash should be deterministic")
	}

	// Verify hash changes with different inputs
	hash3 := hashMessage(timestamp, "different", method, params, nonce)
	if hex.EncodeToString(hash) == hex.EncodeToString(hash3) {
		t.Error("Hash should change with different account")
	}
}

func TestSigningConstantK(t *testing.T) {
	// Verify that K matches the expected value from steem-js
	// K = sha256('steem_jsonrpc_auth')
	expected := sha256.Sum256([]byte("steem_jsonrpc_auth"))
	
	if hex.EncodeToString(K) != hex.EncodeToString(expected[:]) {
		t.Errorf("Signing constant K mismatch.\nExpected: %s\nGot: %s", 
			hex.EncodeToString(expected[:]), hex.EncodeToString(K))
	}
}

func TestSignRequest(t *testing.T) {
	// Test the convenience function
	signedRequest, err := SignRequest(testMethod, testParams, 1, testAccount, testPrivateKey)
	if err != nil {
		t.Fatalf("SignRequest failed: %v", err)
	}

	if signedRequest.Method != testMethod {
		t.Errorf("Expected method '%s', got '%s'", testMethod, signedRequest.Method)
	}

	if signedRequest.ID != 1 {
		t.Errorf("Expected ID 1, got %d", signedRequest.ID)
	}

	if signedRequest.Params.Signed.Account != testAccount {
		t.Errorf("Expected account '%s', got '%s'", testAccount, signedRequest.Params.Signed.Account)
	}
}

func TestMultipleSignatures(t *testing.T) {
	// Test signing with multiple private keys
	privateKeys := []string{
		testPrivateKey,
		"5JRaypasxMx1L97ZUX7YuC5Psb5EAbF821kkAGtBj7xCJFQcbLg", // Another test key
	}

	request := &RpcRequest{
		Method: testMethod,
		Params: testParams,
		ID:     1,
	}

	signedRequest, err := Sign(request, testAccount, privateKeys)
	if err != nil {
		t.Fatalf("Sign with multiple keys failed: %v", err)
	}

	if len(signedRequest.Params.Signed.Signatures) != 2 {
		t.Errorf("Expected 2 signatures, got %d", len(signedRequest.Params.Signed.Signatures))
	}

	// Verify both signatures are different
	sig1 := signedRequest.Params.Signed.Signatures[0]
	sig2 := signedRequest.Params.Signed.Signatures[1]
	if sig1 == sig2 {
		t.Error("Multiple signatures should be different")
	}
}

func TestJsonSerialization(t *testing.T) {
	request := &RpcRequest{
		Method: testMethod,
		Params: testParams,
		ID:     1,
	}

	signedRequest, err := Sign(request, testAccount, []string{testPrivateKey})
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(signedRequest)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled SignedRequest
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify the unmarshaled data matches
	if unmarshaled.Method != signedRequest.Method {
		t.Errorf("Method mismatch after JSON round-trip")
	}

	if unmarshaled.Params.Signed.Account != signedRequest.Params.Signed.Account {
		t.Errorf("Account mismatch after JSON round-trip")
	}
}

// Benchmark tests
func BenchmarkSign(b *testing.B) {
	request := &RpcRequest{
		Method: testMethod,
		Params: testParams,
		ID:     1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sign(request, testAccount, []string{testPrivateKey})
		if err != nil {
			b.Fatalf("Sign failed: %v", err)
		}
	}
}

func BenchmarkHashMessage(b *testing.B) {
	timestamp := "2023-11-27T10:30:00.000Z"
	account := "testuser"
	method := "condenser_api.get_accounts"
	params := "W1sidGVzdHVzZXIiXV0="
	nonce := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashMessage(timestamp, account, method, params, nonce)
	}
}
