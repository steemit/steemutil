package protocol

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/steemit/steemutil/encoder"
)

// PublicKey used across all tests:
//   STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA
//   Binary (33 bytes): 03fdf4907810a9f5d9462a1ae09feee5ab205d32798b0ffcc379442021f84c5bbf
const testPubKey = "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA"

func encodeOp(t *testing.T, op interface{}) []byte {
	t.Helper()
	var buf bytes.Buffer
	enc := encoder.NewEncoder(&buf)
	if err := enc.Encode(op); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	return buf.Bytes()
}

func assertHex(t *testing.T, name string, got []byte, expectedHex string) {
	t.Helper()
	gotHex := hex.EncodeToString(got)
	if gotHex != expectedHex {
		t.Errorf("%s hex mismatch:\n  expected (%d bytes): %s\n  got      (%d bytes): %s",
			name, len(expectedHex)/2, expectedHex, len(got)/2, gotHex)
	}
}

// TestWitnessUpdateSerialization verifies that block_signing_key is encoded
// as a 33-byte public_key_type binary, NOT a varint-prefixed string.
func TestWitnessUpdateSerialization(t *testing.T) {
	op := &WitnessUpdateOperation{
		Owner:           "owner",
		URL:             "https://example.com",
		BlockSigningKey: testPubKey,
		Props: &ChainProperties{
			AccountCreationFee: "0.100 STEEM",
			MaximumBlockSize:   131072,
			SBDInterestRate:    1000,
		},
		Fee: "0.000 STEEM",
	}
	// op_type(0b) + owner + url + 33-byte pubkey + props(3 assets) + fee asset
	assertHex(t, "WitnessUpdate", encodeOp(t, op),
		"0b056f776e65721368747470733a2f2f6578616d706c652e636f6d03fdf4907810a9f5d9462a1ae09feee5ab205d32798b0ffcc379442021f84c5bbf640000000000000003535445454d000000000200e803000000000000000003535445454d0000")
}

// TestAccountCreateSerialization verifies memo_key is encoded as 33-byte pubkey
// and the plain (non-optional) authority pointers are serialized without presence bytes.
func TestAccountCreateSerialization(t *testing.T) {
	op := &AccountCreateOperation{
		Fee:            "0.000 STEEM",
		Creator:        "creator",
		NewAccountName: "newaccount",
		Owner:          &Authority{WeightThreshold: 1},
		Active:         &Authority{WeightThreshold: 1},
		Posting:        &Authority{WeightThreshold: 1},
		MemoKey:        testPubKey,
		JsonMetadata:   "{}",
	}
	assertHex(t, "AccountCreate", encodeOp(t, op),
		"09000000000000000003535445454d00000763726561746f720a6e65776163636f756e7400000100000000000100000000000100000003fdf4907810a9f5d9462a1ae09feee5ab205d32798b0ffcc379442021f84c5bbf027b7d")
}

// TestAccountUpdateSerialization verifies optional authority fields (nil → 0x00)
// alongside pubkey encoding for memo_key.
func TestAccountUpdateSerialization(t *testing.T) {
	op := &AccountUpdateOperation{
		Account:      "account",
		Owner:        nil, // optional → 0x00
		Active:       nil, // optional → 0x00
		Posting:      nil, // optional → 0x00
		MemoKey:      testPubKey,
		JsonMetadata: "",
	}
	// op_type(0a) + account + 3×0x00 (optional nil) + 33-byte pubkey + empty metadata + (no extensions field)
	assertHex(t, "AccountUpdate", encodeOp(t, op),
		"0a076163636f756e7400000003fdf4907810a9f5d9462a1ae09feee5ab205d32798b0ffcc379442021f84c5bbf00")
}

// TestAccountUpdate2Serialization verifies memo_key pubkey encoding with
// optional authorities and extensions field.
func TestAccountUpdate2Serialization(t *testing.T) {
	op := &AccountUpdate2Operation{
		Account:             "account",
		Owner:               nil,
		Active:              nil,
		Posting:             nil,
		MemoKey:             testPubKey,
		JsonMetadata:        "",
		PostingJsonMetadata: "",
		Extensions:          []interface{}{},
	}
	assertHex(t, "AccountUpdate2", encodeOp(t, op),
		"2b076163636f756e7400000003fdf4907810a9f5d9462a1ae09feee5ab205d32798b0ffcc379442021f84c5bbf000000")
}

// TestCreateClaimedAccountSerialization verifies memo_key pubkey encoding
// with plain authorities and extensions.
func TestCreateClaimedAccountSerialization(t *testing.T) {
	op := &CreateClaimedAccountOperation{
		Creator:        "creator",
		NewAccountName: "newaccount",
		Owner:          &Authority{WeightThreshold: 1},
		Active:         &Authority{WeightThreshold: 1},
		Posting:        &Authority{WeightThreshold: 1},
		MemoKey:        testPubKey,
		JsonMetadata:   "{}",
		Extensions:     []interface{}{},
	}
	assertHex(t, "CreateClaimedAccount", encodeOp(t, op),
		"170763726561746f720a6e65776163636f756e7400000100000000000100000000000100000003fdf4907810a9f5d9462a1ae09feee5ab205d32798b0ffcc379442021f84c5bbf027b7d00")
}

// TestAccountCreateWithDelegationSerialization verifies memo_key pubkey encoding
// with delegation asset and extensions.
func TestAccountCreateWithDelegationSerialization(t *testing.T) {
	op := &AccountCreateWithDelegationOperation{
		Fee:            "0.000 STEEM",
		Delegation:     "1.000000 VESTS",
		Creator:        "creator",
		NewAccountName: "newaccount",
		Owner:          &Authority{WeightThreshold: 1},
		Active:         &Authority{WeightThreshold: 1},
		Posting:        &Authority{WeightThreshold: 1},
		MemoKey:        testPubKey,
		JsonMetadata:   "{}",
		Extensions:     []interface{}{},
	}
	assertHex(t, "AccountCreateWithDelegation", encodeOp(t, op),
		"29000000000000000003535445454d000040420f000000000006564553545300000763726561746f720a6e65776163636f756e7400000100000000000100000000000100000003fdf4907810a9f5d9462a1ae09feee5ab205d32798b0ffcc379442021f84c5bbf027b7d00")
}
