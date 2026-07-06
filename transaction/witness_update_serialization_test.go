package transaction

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/steemit/steemutil/protocol"
)

// TestWitnessUpdateSerialization verifies that WitnessUpdateOperation serializes
// the props field as a plain struct WITHOUT an erroneous 0x01 "pointer present"
// flag. The bug was that the encoder treated all pointer fields as optional<T>.
func TestWitnessUpdateSerialization(t *testing.T) {
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

	// The bug produced 132 bytes (extra 0x01 before props). Correct is 131.
	if len(serialized) != 131 {
		t.Errorf("expected 131 bytes, got %d bytes (hex: %s)", len(serialized), hex.EncodeToString(serialized))
	}

	// The byte right after block_signing_key should NOT be 0x01 (pointer-present flag).
	// It should be the first byte of the ChainProperties.account_creation_fee asset (0xb8).
	keyStr := "STM7cKkyBCxTRkdxoGSP5ypagNeaVYeR7oWcBsNfKurXg5jnAtNqm"
	keyEnd := -1
	for i := 0; i < len(serialized)-len(keyStr); i++ {
		if string(serialized[i:i+len(keyStr)]) == keyStr {
			keyEnd = i + len(keyStr)
			break
		}
	}
	if keyEnd == -1 {
		t.Fatalf("could not find block_signing_key in serialized output")
	}

	if serialized[keyEnd] == 0x01 {
		t.Errorf("found erroneous 0x01 pointer-present flag at offset %d; "+
			"witness_update props must serialize as plain struct, not optional", keyEnd)
	}
}
