package rpc

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/steemit/steemutil/wif"
)

// K is the signing constant used to reserve opcode space and prevent cross-protocol attacks.
// This is the output of sha256('steem_jsonrpc_auth').
var K = []byte{
	0x3b, 0x3b, 0x08, 0x1e, 0x46, 0xea, 0x80, 0x8d,
	0x5a, 0x96, 0xb0, 0x8c, 0x4b, 0xc5, 0x00, 0x3f,
	0x5e, 0x15, 0x76, 0x70, 0x90, 0xf3, 0x44, 0xfa,
	0xab, 0x53, 0x1e, 0xc5, 0x75, 0x65, 0x13, 0x6b,
}

// RpcRequest represents a JSON-RPC request to be signed.
type RpcRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     int           `json:"id"`
}

// SignedRequest represents a signed JSON-RPC request.
type SignedRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`
	Params  struct {
		Signed SignedParams `json:"__signed"`
	} `json:"params"`
}

// SignedParams contains the signed payload data.
type SignedParams struct {
	Account    string   `json:"account"`
	Nonce      string   `json:"nonce"`
	Params     string   `json:"params"`
	Signatures []string `json:"signatures"`
	Timestamp  string   `json:"timestamp"`
}

// Sign creates a signed JSON-RPC request.
// The request is signed using the provided account and private keys.
func Sign(request *RpcRequest, account string, privateKeys []string) (*SignedRequest, error) {
	if request.Params == nil {
		return nil, errors.New("unable to sign a request without params")
	}

	// Encode params as base64
	paramsJSON, err := json.Marshal(request.Params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal params to JSON")
	}
	params := base64.StdEncoding.EncodeToString(paramsJSON)

	// Generate 8-byte random nonce
	nonceBytes := make([]byte, 8)
	if _, err := rand.Read(nonceBytes); err != nil {
		return nil, errors.Wrap(err, "failed to generate nonce")
	}
	nonce := hex.EncodeToString(nonceBytes)

	// Create ISO8601 timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)

	// Create message hash
	message := hashMessage(timestamp, account, request.Method, params, nonceBytes)

	// Sign with each private key
	signatures := make([]string, 0, len(privateKeys))
	for _, keyWif := range privateKeys {
		privateKey := &wif.PrivateKey{}
		if err := privateKey.FromWif(keyWif); err != nil {
			return nil, errors.Wrapf(err, "failed to decode private key: %s", keyWif)
		}

		signature, err := privateKey.SignSha256(message)
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign message")
		}

		signatures = append(signatures, hex.EncodeToString(signature))
	}

	// Create signed request
	signedRequest := &SignedRequest{
		JsonRpc: "2.0",
		Method:  request.Method,
		ID:      request.ID,
	}
	signedRequest.Params.Signed = SignedParams{
		Account:    account,
		Nonce:      nonce,
		Params:     params,
		Signatures: signatures,
		Timestamp:  timestamp,
	}

	return signedRequest, nil
}

// Validate validates a signed JSON-RPC request.
// The verifyFunc should verify that the signatures are valid for the given account.
func Validate(request *SignedRequest, verifyFunc func(message []byte, signatures []string, account string) error) ([]interface{}, error) {
	if request.JsonRpc != "2.0" || request.Method == "" {
		return nil, errors.New("invalid JSON RPC request")
	}

	signed := request.Params.Signed

	if signed.Account == "" {
		return nil, errors.New("missing account")
	}

	// Decode and validate params
	var params []interface{}
	paramsJSON, err := base64.StdEncoding.DecodeString(signed.Params)
	if err != nil {
		return nil, errors.Wrap(err, "invalid encoded params")
	}

	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal params")
	}

	// Validate nonce
	if signed.Nonce == "" {
		return nil, errors.New("invalid nonce")
	}

	nonceBytes, err := hex.DecodeString(signed.Nonce)
	if err != nil || len(nonceBytes) != 8 {
		return nil, errors.New("invalid nonce format")
	}

	// Validate timestamp
	timestamp, err := time.Parse(time.RFC3339Nano, signed.Timestamp)
	if err != nil {
		return nil, errors.Wrap(err, "invalid timestamp")
	}

	// Check if signature has expired (60 seconds)
	if time.Since(timestamp) > 60*time.Second {
		return nil, errors.New("signature expired")
	}

	// Recreate message hash
	message := hashMessage(signed.Timestamp, signed.Account, request.Method, signed.Params, nonceBytes)

	// Verify signatures
	if err := verifyFunc(message, signed.Signatures, signed.Account); err != nil {
		return nil, errors.Wrap(err, "verification failed")
	}

	return params, nil
}

// hashMessage creates the message hash to be signed.
// This follows the same algorithm as steem-js:
// message = sha256(K + sha256(timestamp + account + method + params + nonce))
func hashMessage(timestamp, account, method, params string, nonce []byte) []byte {
	// First hash: sha256(timestamp + account + method + params + nonce)
	first := sha256.New()
	first.Write([]byte(timestamp))
	first.Write([]byte(account))
	first.Write([]byte(method))
	first.Write([]byte(params))
	first.Write(nonce)
	firstHash := first.Sum(nil)

	// Second hash: sha256(K + firstHash)
	second := sha256.New()
	second.Write(K)
	second.Write(firstHash)

	return second.Sum(nil)
}

// DefaultVerifyFunc provides a default verification function that uses public key recovery.
// This function attempts to recover the public key from the signature and verify it matches
// the expected account's public keys.
func DefaultVerifyFunc(message []byte, signatures []string, account string) error {
	// This is a placeholder implementation.
	// In a real implementation, you would:
	// 1. Get the account's public keys from the blockchain
	// 2. For each signature, recover the public key
	// 3. Verify that the recovered public key matches one of the account's keys

	if len(signatures) == 0 {
		return errors.New("no signatures provided")
	}

	// For now, just validate that signatures are properly formatted hex strings
	for i, sig := range signatures {
		if _, err := hex.DecodeString(sig); err != nil {
			return errors.Wrapf(err, "invalid signature format at index %d", i)
		}
	}

	return nil
}

// SignRequest is a convenience function that signs a request with a single private key.
func SignRequest(method string, params []interface{}, id int, account string, privateKey string) (*SignedRequest, error) {
	request := &RpcRequest{
		Method: method,
		Params: params,
		ID:     id,
	}

	return Sign(request, account, []string{privateKey})
}
