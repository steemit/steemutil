package transaction

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/steemit/steemutil/protocol"
)

// TestWitnessUpdateGoldenHex verifies the exact byte-level serialization of a
// witness_update transaction. The golden hex was captured from the fixed encoder
// and validated against steemd's expected C++ FC_REFLECT serialization.
//
// The regression this guards against: the encoder used to insert an erroneous
// 0x01 "pointer-present" byte before ChainProperties (a plain, non-optional
// struct), producing 132 bytes instead of 131 and corrupting the digest.
func TestWitnessUpdateGoldenHex(t *testing.T) {
	expiration := time.Date(2026, 7, 6, 11, 31, 28, 0, time.UTC)

	op := &protocol.WitnessUpdateOperation{
		Owner:           "ety001",
		URL:             "https://steem.fans",
		BlockSigningKey: "STM7cKkyBCxTRkdxoGSP5ypagNeaVYeR7oWcBsNfKurXg5jnAtNqm",
		Props: &protocol.ChainProperties{
			AccountCreationFee: "3.000 STEEM",
			MaximumBlockSize:   65536,
			SBDInterestRate:    0,
		},
		Fee: "0.000 STEEM",
	}

	stx := NewSignedTransaction(&Transaction{
		RefBlockNum:    18668,
		RefBlockPrefix: 2441809904,
		Expiration:     &protocol.Time{Time: &expiration},
		Extensions:     []interface{}{},
	})
	stx.PushOperation(op)

	serialized, err := stx.Serialize()
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	// Expected golden hex (131 bytes). Breakdown for reference:
	//   ec48            ref_block_num (18668 LE)
	//   f00f8b91        ref_block_prefix (2441809904 LE)
	//   0924b6a0        expiration (unix 1783337488 LE)
	//   01              op_count (1)
	//   0b              op_type (11 = witness_update)
	//   06657479303031  owner "ety001" (varint-len 6)
	//   1268...66616e73 url "https://steem.fans" (varint-len 18)
	//   3553...4e716d   block_signing_key (varint-len 53)
	//   b80b...000003   props.account_creation_fee (asset: uint64 3000, prec 3, "STEEM\0\0")
	//   00010000        props.maximum_block_size (65536 LE)
	//   0000            props.sbd_interest_rate (0 LE)
	//   0000...0003     fee (asset: uint64 0, prec 3, "STEEM\0\0")
	//   00              extensions count (0)
	//
	// CRITICAL: the byte after block_signing_key is 0xb8 (first byte of the
	// asset amount), NOT 0x01. The old encoder inserted an erroneous 0x01 here.
	const expectedHex = "ec48f00f8b9110924b6a010b066574793030311268747470733a2f2f737465656d2e66616e733553544d37634b6b7942437854526b64786f47535035797061674e656156596552376f576342734e664b75725867356a6e41744e716db80b00000000000003535445454d0000000001000000000000000000000003535445454d000000"

	got := hex.EncodeToString(serialized)
	if got != expectedHex {
		t.Errorf("golden hex mismatch:\n  expected (%d bytes): %s\n  got      (%d bytes): %s",
			len(expectedHex)/2, expectedHex, len(got)/2, got)
	}
}

// TestAccountUpdateOptionalFields verifies that truly optional pointer fields
// (account_update owner/active/posting) still emit the 0x01/0x00 presence byte,
// while non-optional pointers (witness_update props) do not.
func TestAccountUpdateOptionalFields(t *testing.T) {
	expiration := time.Date(2026, 7, 6, 11, 31, 28, 0, time.UTC)

	activeAuth := &protocol.Authority{
		WeightThreshold: 1,
		KeyAuths:        protocol.StringInt64Map{"STM6z4ueBGtT96KczqAEzQQTxPUkX5FTVKMEjRb1viFeS2eBSGgKu": 1},
	}

	// owner=nil (optional absent), active=set (optional present), posting=nil
	op := &protocol.AccountUpdateOperation{
		Account:      "ety001",
		Owner:        nil, // optional → should serialize as 0x00
		Active:       activeAuth,
		Posting:      nil, // optional → should serialize as 0x00
		MemoKey:      "STM6z4ueBGtT96KczqAEzQQTxPUkX5FTVKMEjRb1viFeS2eBSGgKu",
		JsonMetadata: "",
	}

	stx := NewSignedTransaction(&Transaction{
		RefBlockNum:    18668,
		RefBlockPrefix: 2441809904,
		Expiration:     &protocol.Time{Time: &expiration},
		Extensions:     []interface{}{},
	})
	stx.PushOperation(op)

	serialized, err := stx.Serialize()
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	// Locate the account name "ety001" (varint-len 6 + "ety001") = 7 bytes after op_type.
	// Layout: header(11) + op_count(1) + op_type(1) + acct_len(1) + "ety001"(6)
	//         = offset 20 → owner flag at offset 20.
	acctEnd := -1
	acct := []byte("ety001")
	for i := 0; i < len(serialized)-len(acct); i++ {
		if string(serialized[i:i+len(acct)]) == "ety001" {
			acctEnd = i + len(acct)
			break
		}
	}
	if acctEnd == -1 {
		t.Fatalf("could not find account name in serialized output")
	}

	// owner (optional nil) → must be 0x00
	if serialized[acctEnd] != 0x00 {
		t.Errorf("owner (optional, nil) should encode as 0x00, got 0x%02x", serialized[acctEnd])
	}

	// active (optional, set) → must be 0x01 followed by the Authority struct
	activeOffset := acctEnd + 1
	if serialized[activeOffset] != 0x01 {
		t.Errorf("active (optional, set) should encode as 0x01, got 0x%02x", serialized[activeOffset])
	}
}
