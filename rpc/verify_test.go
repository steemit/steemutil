package rpc

import (
	"encoding/hex"
	"strings"
	"testing"
)

// errNotFound is a sentinel returned by test fetchers to mean "no such account".
var errNotFound = &notFoundErr{}

type notFoundErr struct{}

func (*notFoundErr) Error() string { return "no such account" }

// jsVector bundles the components of a request signed by the real JS stack
// (@steemit/libcrypto signRecoverably + the rpc-auth hashMessage algorithm).
// The signature is non-deterministic, so it was captured once and fixed here.
// The digest is reproducible from the other fields via hashMessage.
var jsVector = struct {
	wif, account, method, timestamp, nonceHex, paramsBase64, sigHex, pubKey string
}{
	wif:          "5JWHY5DxTF6qN5grTtChDCYBmWHfY9zaSsw4CxEKN5eZpH9iBma",
	account:      "testaccount",
	method:       "condenser_api.get_accounts",
	timestamp:    "2024-01-01T00:00:00.000Z",
	nonceHex:     "0102030405060708",
	paramsBase64: "W1sidGVzdGFjY291bnQiXV0=",
	sigHex:       "201cba8696cdecef383702be9b4349e6418adf0973c018534e11b30012f5d893d32e083c38fc719c7470c64bba949b2d4c94bd74bafb8df80a6c4495a51a4ed054",
	pubKey:       "STM7jNh5ejQoqHqWcGWFJ1v4F5CzsG3EiBuz1VooCng1cH5QpJD27",
}

// TestHashMessage_MatchesJS pins the Go hashMessage digest to the value
// computed by the reference JS implementation (@steemit/rpc-auth hashMessage).
// This is the regression guard for U3: the nonce must live in the OUTER hash.
// The digest is derived from the fixed jsVector inputs and was independently
// confirmed against Node's crypto by running the JS hashMessage algorithm.
func TestHashMessage_MatchesJS(t *testing.T) {
	nonce, err := hex.DecodeString(jsVector.nonceHex)
	if err != nil {
		t.Fatalf("bad nonce hex: %v", err)
	}
	got := hashMessage(jsVector.timestamp, jsVector.account, jsVector.method, jsVector.paramsBase64, nonce)
	// Computed via the JS @steemit/rpc-auth hashMessage (nonce in second hash).
	const wantHex = "4b1eecc536155df76ce97c8879c4429154f856e9dee2d1fb5fd942a9f1a7ebf4"
	if hex.EncodeToString(got) != wantHex {
		t.Fatalf("hashMessage mismatch\nGo  : %s\nWant: %s", hex.EncodeToString(got), wantHex)
	}
}

// singleKeyFetcher builds an AccountFetcher that returns the given posting key
// (weight 1, threshold 1) for any account.
func singleKeyFetcher(pubKey string) AccountFetcher {
	return func(string) (AccountPostingAuth, error) {
		return AccountPostingAuth{
			KeyAuths:        []KeyAuth{{PubKey: pubKey, Weight: 1}},
			WeightThreshold: 1,
		}, nil
	}
}

// TestVerifySignedRpc_JSVector verifies a signature produced by the real JS
// signing stack (@steemit/libcrypto). This is acceptance criterion 1: the Go
// verifier accepts a genuine JS-signed request.
func TestVerifySignedRpc_JSVector(t *testing.T) {
	msg, _ := hex.DecodeString("4b1eecc536155df76ce97c8879c4429154f856e9dee2d1fb5fd942a9f1a7ebf4")
	err := VerifySignedRpc(msg, []string{jsVector.sigHex}, jsVector.account, singleKeyFetcher(jsVector.pubKey))
	if err != nil {
		t.Fatalf("expected valid signature to verify, got: %v", err)
	}
}

// TestVerifySignedRpc_Tampered verifies that altering any byte of the signature
// fails verification. Acceptance criterion 2.
func TestVerifySignedRpc_Tampered(t *testing.T) {
	msg, _ := hex.DecodeString("4b1eecc536155df76ce97c8879c4429154f856e9dee2d1fb5fd942a9f1a7ebf4")

	sigBytes, _ := hex.DecodeString(jsVector.sigHex)

	// Tamper each of bytes [1, 65) (the r‖s payload) and expect failure.
	for i := 1; i < len(sigBytes); i++ {
		flipped := make([]byte, len(sigBytes))
		copy(flipped, sigBytes)
		flipped[i] ^= 0x40
		err := VerifySignedRpc(msg, []string{hex.EncodeToString(flipped)}, jsVector.account, singleKeyFetcher(jsVector.pubKey))
		if err == nil {
			t.Errorf("tampering payload byte %d did not fail verification", i)
		}
	}

	// Tamper byte0 (the recovery byte). Values in [31, 35) pass DsteemSigToBtcec
	// validation but recover a different (or invalid) public key; values outside
	// [31, 35) are rejected by DsteemSigToBtcec. All must fail verification.
	origByte0 := sigBytes[0] // 32
	for _, b0 := range []byte{27, 30, 31, 33, 34, 35} {
		flipped := make([]byte, len(sigBytes))
		copy(flipped, sigBytes)
		flipped[0] = b0
		err := VerifySignedRpc(msg, []string{hex.EncodeToString(flipped)}, jsVector.account, singleKeyFetcher(jsVector.pubKey))
		if err == nil {
			t.Errorf("tampering recovery byte to %d (orig %d) did not fail verification", b0, origByte0)
		}
	}
}

// TestVerifySignedRpc_ErrorMessages checks each failure path produces the exact
// error message the JS verifier (@steemit/koa-jsonrpc auth.ts) emits, so that
// Go and JS servers are indistinguishable to clients.
func TestVerifySignedRpc_ErrorMessages(t *testing.T) {
	msg := make([]byte, 32) // valid-length digest placeholder
	fetcher := singleKeyFetcher(jsVector.pubKey)

	cases := []struct {
		name       string
		message    []byte
		signatures []string
		account    string
		fetcher    AccountFetcher
		want       string
	}{
		{
			name:       "wrong message length",
			message:    make([]byte, 31),
			signatures: []string{jsVector.sigHex},
			account:    jsVector.account,
			fetcher:    fetcher,
			want:       "Invalid message",
		},
		{
			name:       "account too short",
			message:    msg,
			signatures: []string{jsVector.sigHex},
			account:    "ab",
			fetcher:    fetcher,
			want:       "Invalid account name",
		},
		{
			name:       "account too long",
			message:    msg,
			signatures: []string{jsVector.sigHex},
			account:    strings.Repeat("a", 17),
			fetcher:    fetcher,
			want:       "Invalid account name",
		},
		{
			name:       "no such account",
			message:    msg,
			signatures: []string{jsVector.sigHex},
			account:    "noexistaccount",
			fetcher: func(string) (AccountPostingAuth, error) {
				return AccountPostingAuth{}, errNotFound
			},
			want: "No such account",
		},
		{
			name:       "multi posting key",
			message:    msg,
			signatures: []string{jsVector.sigHex},
			account:    jsVector.account,
			fetcher: func(string) (AccountPostingAuth, error) {
				return AccountPostingAuth{
					KeyAuths: []KeyAuth{
						{PubKey: jsVector.pubKey, Weight: 1},
						{PubKey: "STM7W7ACQDZJZ6rZGKeT9auipnSiSxFxJ4k71QXmrhY9HbvYsNnQ2", Weight: 1},
					},
					WeightThreshold: 1,
				}, nil
			},
			want: "Unsupported posting key configuration for account",
		},
		{
			name:       "key below weight threshold",
			message:    msg,
			signatures: []string{jsVector.sigHex},
			account:    jsVector.account,
			fetcher: func(string) (AccountPostingAuth, error) {
				return AccountPostingAuth{
					KeyAuths:        []KeyAuth{{PubKey: jsVector.pubKey, Weight: 1}},
					WeightThreshold: 2,
				}, nil
			},
			want: "Signing key not above weight threshold",
		},
		{
			name:       "multisig (two signatures)",
			message:    msg,
			signatures: []string{jsVector.sigHex, jsVector.sigHex},
			account:    jsVector.account,
			fetcher:    fetcher,
			want:       "Multisig not supported",
		},
		{
			name:       "bad signature hex",
			message:    msg,
			signatures: []string{"nothex"},
			account:    jsVector.account,
			fetcher:    fetcher,
			want:       "Invalid signature",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifySignedRpc(tc.message, tc.signatures, tc.account, tc.fetcher)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tc.want)
			}
			if err.Error() != tc.want {
				t.Errorf("error mismatch\nwant: %q\ngot:  %q", tc.want, err.Error())
			}
		})
	}
}

// TestVerifySignedRpc_WrongPubKey verifies that a valid signature which does
// not correspond to the account's posted key fails.
func TestVerifySignedRpc_WrongPubKey(t *testing.T) {
	msg, _ := hex.DecodeString("4b1eecc536155df76ce97c8879c4429154f856e9dee2d1fb5fd942a9f1a7ebf4")
	// account claims a DIFFERENT key than the one that signed
	other := singleKeyFetcher("STM7W7ACQDZJZ6rZGKeT9auipnSiSxFxJ4k71QXmrhY9HbvYsNnQ2")
	err := VerifySignedRpc(msg, []string{jsVector.sigHex}, jsVector.account, other)
	if err == nil || err.Error() != "Invalid signature" {
		t.Errorf("expected Invalid signature, got: %v", err)
	}
}

// TestValidate_WithVerifySignedRpc confirms the end-to-end Validate path works
// when wired up with a VerifySignedRpc-bound closure, matching how conveyor
// will use it: Validate(req, func(...){ return VerifySignedRpc(..., fetcher) }).
//
// It signs fresh (so the timestamp is within the 60s window), reconstructs the
// same digest inside Validate, and the bound VerifySignedRpc recovers the key.
// JS compatibility of the digest/signature is covered separately by
// TestHashMessage_MatchesJS and TestVerifySignedRpc_JSVector.
func TestValidate_WithVerifySignedRpc(t *testing.T) {
	account := jsVector.account
	request := &RpcRequest{
		Method: jsVector.method,
		Params: []interface{}{[]string{account}},
		ID:     1,
	}
	signed, err := Sign(request, account, []string{jsVector.wif})
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	fetcher := singleKeyFetcher(jsVector.pubKey)
	params, err := Validate(signed, func(message []byte, signatures []string, acct string) error {
		return VerifySignedRpc(message, signatures, acct, fetcher)
	})
	if err != nil {
		t.Fatalf("Validate with VerifySignedRpc failed: %v", err)
	}
	if len(params) != 1 {
		t.Errorf("expected 1 param, got %d", len(params))
	}
}

// sanity guard that a recovered key matches expectations, reused above.
var _ = errNotFound
