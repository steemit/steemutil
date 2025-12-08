package transaction

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/steemit/steemutil/wif"
)

// TestSignatureFormat tests the signature format to understand how btcec.SignCompact encodes recovery parameter
func TestSignatureFormat(t *testing.T) {
	// Use a test private key
	testWif := "5JRaypasxMx1L97ZUX7YuC5Psb5EAbF821kkAGtBj7xCJFQcbLg"
	privKey := &wif.PrivateKey{}
	if err := privKey.FromWif(testWif); err != nil {
		t.Fatalf("Failed to decode WIF: %v", err)
	}

	// Get public key
	pubKeyStr := privKey.ToPubKeyStr()
	fmt.Printf("Public Key: %s\n", pubKeyStr)

	// Create a test digest (32 bytes)
	digest := make([]byte, 32)
	for i := range digest {
		digest[i] = byte(i)
	}

	// Sign with SignCompact
	sig, err := ecdsa.SignCompact(privKey.Raw.PrivKey, digest, true)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	fmt.Printf("Signature length: %d bytes\n", len(sig))
	fmt.Printf("Signature (hex): %s\n", hex.EncodeToString(sig))

	// Check first byte (recovery parameter)
	recoveryID := sig[0]
	fmt.Printf("Recovery ID (first byte): %d (0x%02x)\n", recoveryID, recoveryID)

	// Extract recovery parameter according to Steem format
	// Steem C++ uses: (recoveryID - 27) & 3
	recoveryParam := (recoveryID - 27) & 3
	fmt.Printf("Recovery parameter (Steem format): %d\n", recoveryParam)

	// Check if compressed flag is set
	compressedFlag := (recoveryID - 27) & 4
	fmt.Printf("Compressed flag: %d (should be 4 for compressed)\n", compressedFlag)

	// Try to recover public key
	recoveredPubKey, wasCompressed, err := ecdsa.RecoverCompact(sig, digest)
	if err != nil {
		t.Fatalf("Failed to recover public key: %v", err)
	}

	fmt.Printf("Recovered compressed: %v\n", wasCompressed)

	// Convert recovered public key to string format
	recoveredPubKeyBytes := recoveredPubKey.SerializeCompressed()
	fmt.Printf("Recovered public key (compressed): %s\n", hex.EncodeToString(recoveredPubKeyBytes))

	// Verify it matches
	expectedPubKey := &wif.PublicKey{}
	if err := expectedPubKey.FromStr(pubKeyStr); err != nil {
		t.Fatalf("Failed to parse expected public key: %v", err)
	}

	expectedPubKeyBytes := expectedPubKey.Raw.SerializeCompressed()
	fmt.Printf("Expected public key (compressed): %s\n", hex.EncodeToString(expectedPubKeyBytes))

	if hex.EncodeToString(recoveredPubKeyBytes) == hex.EncodeToString(expectedPubKeyBytes) {
		fmt.Println("✅ Recovered public key matches!")
	} else {
		fmt.Println("❌ Recovered public key does NOT match!")
		t.Errorf("Public key mismatch")
	}

	// Check if recovery ID is in expected range for compressed keys
	// For compressed: 27 + 4 + (0-3) = 31-34
	// For uncompressed: 27 + (0-3) = 27-30
	if recoveryID >= 31 && recoveryID <= 34 {
		fmt.Printf("✅ Recovery ID %d is in expected range for compressed keys (31-34)\n", recoveryID)
	} else if recoveryID >= 27 && recoveryID <= 30 {
		fmt.Printf("⚠️  Recovery ID %d is in range for uncompressed keys (27-30), but we used compressed=true\n", recoveryID)
	} else {
		fmt.Printf("❌ Recovery ID %d is outside expected range\n", recoveryID)
	}

	// Check if signature is canonical (Steem fc_canonical format)
	// is_fc_canonical checks:
	//   !(c.data[1] & 0x80) && !(c.data[1] == 0 && !(c.data[2] & 0x80))
	//   && !(c.data[33] & 0x80) && !(c.data[33] == 0 && !(c.data[34] & 0x80))
	isCanonical := !(sig[1]&0x80 != 0) &&
		!(sig[1] == 0 && sig[2]&0x80 == 0) &&
		!(sig[33]&0x80 != 0) &&
		!(sig[33] == 0 && sig[34]&0x80 == 0)

	if isCanonical {
		fmt.Println("✅ Signature is canonical (fc_canonical format)")
	} else {
		fmt.Println("❌ Signature is NOT canonical (fc_canonical format)")
		fmt.Printf("  r[0] = 0x%02x, r[1] = 0x%02x\n", sig[1], sig[2])
		fmt.Printf("  s[0] = 0x%02x, s[1] = 0x%02x\n", sig[33], sig[34])
		t.Errorf("Signature is not canonical")
	}
}

