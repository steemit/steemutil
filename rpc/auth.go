package rpc

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	api "github.com/steemit/steemutil/protocol/api"
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
//
// The decoded params are returned as interface{} (not []interface{}) to match
// the shape-agnostic semantics of the JS reference (@steemit/rpc-auth validate),
// which returns JSON.parse(jsonString). Signed params can be a JSON object,
// array, or scalar — conveyor and koa-jsonrpc clients sign object params
// (e.g. {"account":"foo"}), so constraining to []interface{} would reject
// every real-world authenticated request.
func Validate(request *SignedRequest, verifyFunc func(message []byte, signatures []string, account string) error) (interface{}, error) {
	if request.JsonRpc != "2.0" || request.Method == "" {
		return nil, errors.New("invalid JSON RPC request")
	}

	signed := request.Params.Signed

	if signed.Account == "" {
		return nil, errors.New("missing account")
	}

	// Decode and validate params. The target is interface{} so that any valid
	// JSON shape (object, array, scalar) is accepted, matching JSON.parse in JS.
	var params interface{}
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
// message = sha256( K ‖ sha256(timestamp ‖ account ‖ method ‖ params) ‖ nonce )
// The nonce (8 raw bytes) is included in the OUTER (second) hash, matching
// @steemit/rpc-auth (see src/index.ts hashMessage). This is the authoritative
// layout; do not move the nonce into the inner hash or signatures will not
// verify against JS-signed requests.
func hashMessage(timestamp, account, method, params string, nonce []byte) []byte {
	// First (inner) hash: sha256(timestamp + account + method + params)
	first := sha256.New()
	first.Write([]byte(timestamp))
	first.Write([]byte(account))
	first.Write([]byte(method))
	first.Write([]byte(params))
	firstHash := first.Sum(nil)

	// Second (outer) hash: sha256(K + firstHash + nonce)
	second := sha256.New()
	second.Write(K)
	second.Write(firstHash)
	second.Write(nonce)

	return second.Sum(nil)
}

// DefaultVerifyFunc is a placeholder verification function: it only checks
// that signatures are well-formed hex strings and performs NO cryptographic
// verification. Use VerifySignedRpc (bound to an AccountFetcher) for real
// signature verification against an account's posting authority.
func DefaultVerifyFunc(message []byte, signatures []string, account string) error {
	_ = message
	_ = account

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

// AccountFetcher returns the posting authority for the given account. It is
// injected by the caller (e.g. conveyor uses steemgosdk's get_accounts) so the
// rpc package does not depend on a specific chain-access mechanism. Returning
// the (api.Authority) of the account's posting key_auths/weight_threshold lets
// VerifySignedRpc consume the same type conveyor already deserializes into,
// with no conversion glue. An error from the fetcher (account not found) is
// surfaced to the caller as "No such account".
type AccountFetcher func(account string) (api.Authority, error)

// VerifySignedRpc verifies that one of the signatures was produced by the
// account's single posting key over the given 32-byte message digest. It
// mirrors the @steemit/koa-jsonrpc verifier (src/auth.ts:65-98) byte-for-byte,
// including the order and exact wording of validation errors, so that Go and
// JS servers reject/accept the same requests.
//
// Only the posting authority is checked (active/owner/memo are ignored). The
// JS verifier only supports accounts with a single posting key and a single
// signature, so this function rejects multisig configurations.
//
// accountFetcher must return the account's posting authority (key_auths +
// weight_threshold); a not-found error from the fetcher is reported as
// "No such account".
func VerifySignedRpc(message []byte, signatures []string, account string, accountFetcher AccountFetcher) error {
	// 1. message must be a 32-byte digest
	if len(message) != 32 {
		return errors.New("Invalid message")
	}

	// 2. account name length bounds
	if len(account) < 3 || len(account) > 16 {
		return errors.New("Invalid account name")
	}

	// 3. fetch the account's posting authority
	auth, err := accountFetcher(account)
	if err != nil {
		return errors.New("No such account")
	}

	// 4. only single posting key is supported
	if len(auth.KeyAuths) != 1 {
		return errors.New("Unsupported posting key configuration for account")
	}

	// 5. the posting key must clear the weight threshold
	if uint32(auth.WeightThreshold) > uint32(auth.KeyAuths[0].Weight) {
		return errors.New("Signing key not above weight threshold")
	}

	// 6. only one signature is supported (no multisig)
	if len(signatures) != 1 {
		return errors.New("Multisig not supported")
	}

	// 7. recover the public key from the signature and compare it to the
	//    account's posting key.
	signature, err := hex.DecodeString(signatures[0])
	if err != nil {
		return errors.New("Invalid signature")
	}

	// The signature arrives in dsteem/libcrypto compact format. Since the
	// btcec/decred recovery layout is byte-for-byte identical for compressed
	// keys, ValidateCompactSignature just validates the format (65 bytes,
	// byte0 in [31, 35)) and returns it unchanged before we hand it to
	// RecoverPublicKeyFromSignature.
	btcecSig, err := wif.ValidateCompactSignature(signature)
	if err != nil {
		return errors.New("Invalid signature")
	}

	recovered, err := wif.RecoverPublicKeyFromSignature(message, btcecSig)
	if err != nil {
		return errors.New("Invalid signature")
	}

	expected := &wif.PublicKey{}
	if err := expected.FromStr(auth.KeyAuths[0].PubKey); err != nil {
		return errors.New("Invalid signature")
	}

	if !bytes.Equal(expected.ToByte(), recovered.ToByte()) {
		return errors.New("Invalid signature")
	}

	return nil
}
