package auth

import (
	"testing"

	"github.com/steemit/steemutil/wif"
)

func TestEncodeDecodePlainText(t *testing.T) {
	// Plain text should be returned as-is
	memo := "plain text memo"
	
	encoded, err := Encode(nil, nil, memo)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if encoded != memo {
		t.Errorf("Encode should return plain text as-is, got: %s", encoded)
	}

	decoded, err := Decode(nil, memo)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded != memo {
		t.Errorf("Decode should return plain text as-is, got: %s", decoded)
	}
}

func TestEncodeDecodeWithHashPrefix(t *testing.T) {
	// Create test keys
	senderWif := "5JRaypasxMx1L97ZUX7YuC5Psb5EAbF821kkAGtBj7xCJFQcbLg"
	recipientPubKey := "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA"

	memo := "#test memo"

	// Encode
	encoded, err := Encode(senderWif, recipientPubKey, memo)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if encoded == memo {
		t.Error("Encode should encrypt memo with # prefix")
	}

	if len(encoded) == 0 {
		t.Error("Encode returned empty string")
	}

	// Decode
	decoded, err := Decode(senderWif, encoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded != memo {
		t.Errorf("Decode failed: expected %s, got %s", memo, decoded)
	}
}

func TestEncodeDecodeWithPrivateKeyObject(t *testing.T) {
	// Test with PrivateKey object instead of WIF string
	senderWif := "5JRaypasxMx1L97ZUX7YuC5Psb5EAbF821kkAGtBj7xCJFQcbLg"
	recipientPubKey := "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA"

	senderPrivKey := &wif.PrivateKey{}
	if err := senderPrivKey.FromWif(senderWif); err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	recipientPubKeyObj := &wif.PublicKey{}
	if err := recipientPubKeyObj.FromStr(recipientPubKey); err != nil {
		t.Fatalf("Failed to create public key: %v", err)
	}

	memo := "#test memo with objects"

	// Encode with objects
	encoded, err := Encode(senderPrivKey, recipientPubKeyObj, memo)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode with WIF string (should work both ways)
	decoded, err := Decode(senderWif, encoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded != memo {
		t.Errorf("Decode failed: expected %s, got %s", memo, decoded)
	}
}

func TestEncodeDecodeEmptyMemo(t *testing.T) {
	memo := ""
	
	encoded, err := Encode(nil, nil, memo)
	if err == nil {
		t.Error("Encode should fail for empty memo")
	}

	decoded, err := Decode(nil, memo)
	if err == nil {
		t.Error("Decode should fail for empty memo")
	}

	_ = encoded
	_ = decoded
}

