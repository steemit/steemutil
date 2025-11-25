package protocol

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/steemit/steemutil/encoder"
)

func TestParseAsset_STEEM(t *testing.T) {
	asset, err := ParseAsset("0.001 STEEM")
	if err != nil {
		t.Fatalf("ParseAsset failed: %v", err)
	}

	if asset.Amount != 1 {
		t.Errorf("expected amount 1, got %d", asset.Amount)
	}
	if asset.Precision != 3 {
		t.Errorf("expected precision 3, got %d", asset.Precision)
	}
	if asset.Symbol != "STEEM" {
		t.Errorf("expected symbol STEEM, got %s", asset.Symbol)
	}
}

func TestParseAsset_SBD(t *testing.T) {
	asset, err := ParseAsset("1.000 SBD")
	if err != nil {
		t.Fatalf("ParseAsset failed: %v", err)
	}

	if asset.Amount != 1000 {
		t.Errorf("expected amount 1000, got %d", asset.Amount)
	}
	if asset.Precision != 3 {
		t.Errorf("expected precision 3, got %d", asset.Precision)
	}
	if asset.Symbol != "SBD" {
		t.Errorf("expected symbol SBD, got %s", asset.Symbol)
	}
}

func TestParseAsset_VESTS(t *testing.T) {
	asset, err := ParseAsset("0.000001 VESTS")
	if err != nil {
		t.Fatalf("ParseAsset failed: %v", err)
	}

	if asset.Amount != 1 {
		t.Errorf("expected amount 1, got %d", asset.Amount)
	}
	if asset.Precision != 6 {
		t.Errorf("expected precision 6, got %d", asset.Precision)
	}
	if asset.Symbol != "VESTS" {
		t.Errorf("expected symbol VESTS, got %s", asset.Symbol)
	}
}

func TestAsset_MarshalTransaction_STEEM(t *testing.T) {
	asset, err := ParseAsset("1.000 STEEM")
	if err != nil {
		t.Fatalf("ParseAsset failed: %v", err)
	}

	var b bytes.Buffer
	enc := encoder.NewEncoder(&b)

	if err := asset.MarshalTransaction(enc); err != nil {
		t.Fatalf("MarshalTransaction failed: %v", err)
	}

	// Expected format: int64(1000) + uint8(3) + "STEEM" + 2 null bytes
	// int64(1000) little-endian = e803000000000000
	// uint8(3) = 03
	// "STEEM" + 2 null = 535445454d0000
	expectedHex := "e80300000000000003535445454d0000"
	actualHex := hex.EncodeToString(b.Bytes())

	if actualHex != expectedHex {
		t.Errorf("expected %s, got %s", expectedHex, actualHex)
	}
}

func TestAsset_MarshalTransaction_SBD(t *testing.T) {
	asset, err := ParseAsset("1.000 SBD")
	if err != nil {
		t.Fatalf("ParseAsset failed: %v", err)
	}

	var b bytes.Buffer
	enc := encoder.NewEncoder(&b)

	if err := asset.MarshalTransaction(enc); err != nil {
		t.Fatalf("MarshalTransaction failed: %v", err)
	}

	// Expected format: int64(1000) + uint8(3) + "SBD" + 4 null bytes
	// int64(1000) little-endian = e803000000000000
	// uint8(3) = 03
	// "SBD" + 4 null = 53424400000000
	expectedHex := "e8030000000000000353424400000000"
	actualHex := hex.EncodeToString(b.Bytes())

	if actualHex != expectedHex {
		t.Errorf("expected %s, got %s", expectedHex, actualHex)
	}
}

func TestAsset_String(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"0.001 STEEM", "0.001 STEEM"},
		{"1.000 SBD", "1.000 SBD"},
		{"100.500 SBD", "100.500 SBD"},
		{"0.000001 VESTS", "0.000001 VESTS"},
	}

	for _, tc := range testCases {
		asset, err := ParseAsset(tc.input)
		if err != nil {
			t.Fatalf("ParseAsset failed for %s: %v", tc.input, err)
		}

		result := asset.String()
		if result != tc.expected {
			t.Errorf("for input %s, expected %s, got %s", tc.input, tc.expected, result)
		}
	}
}
