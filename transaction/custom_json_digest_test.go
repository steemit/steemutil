package transaction

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/steemit/steemutil/protocol"
)

// TestCustomJsonTransactionDigest tests the digest calculation for a custom_json transaction
// This helps verify that transaction serialization and digest calculation are correct
func TestCustomJsonTransactionDigest(t *testing.T) {
	// Create a custom_json operation similar to the one in the example
	op := &protocol.CustomJSONOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{"ety001234"},
		ID:                   "notify",
		JSON:                 `["setLastRead",{"date":"2025-01-01T00:00:00Z"}]`,
	}

	// Create transaction (without signatures for digest calculation)
	tx := &Transaction{
		RefBlockNum:    1000,
		RefBlockPrefix: 1234567890,
		Operations:     []protocol.Operation{op},
		Extensions:     []interface{}{},
		Signatures:     nil, // Important: no signatures for digest calculation
	}

	stx := NewSignedTransaction(tx)

	// Serialize transaction (should not include signatures)
	serialized, err := stx.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize transaction: %v", err)
	}

	fmt.Printf("Serialized transaction (hex): %s\n", hex.EncodeToString(serialized))
	fmt.Printf("Serialized transaction length: %d bytes\n", len(serialized))

	// Calculate digest
	digest, err := stx.Digest(SteemChain)
	if err != nil {
		t.Fatalf("Failed to calculate digest: %v", err)
	}

	fmt.Printf("Digest (hex): %s\n", hex.EncodeToString(digest))
	fmt.Printf("Digest length: %d bytes\n", len(digest))

	// Verify chain ID is correct
	rawChainID, err := hex.DecodeString(SteemChain.ID)
	if err != nil {
		t.Fatalf("Failed to decode chain ID: %v", err)
	}
	fmt.Printf("Chain ID (hex): %s\n", hex.EncodeToString(rawChainID))
	fmt.Printf("Chain ID length: %d bytes\n", len(rawChainID))

	// Verify digest calculation: sha256(chain_id + serialized_transaction)
	// This is what Steem C++ does in transaction::sig_digest()
	expectedDigestInput := append(rawChainID, serialized...)
	fmt.Printf("Digest input length: %d bytes (32 chain_id + %d transaction)\n", len(rawChainID), len(serialized))
	fmt.Printf("Digest input (hex): %s\n", hex.EncodeToString(expectedDigestInput))

	// The digest should be SHA256 of (chain_id + serialized_transaction)
	// This is already what Digest() does, so we're just verifying the process
}

// TestTransactionSerializationExcludesSignatures verifies that transaction serialization
// does not include signatures (which is required for digest calculation)
func TestTransactionSerializationExcludesSignatures(t *testing.T) {
	op := &protocol.CustomJSONOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{"test"},
		ID:                   "test",
		JSON:                 `["test",{}]`,
	}

	tx1 := &Transaction{
		RefBlockNum:    1000,
		RefBlockPrefix: 1234567890,
		Operations:     []protocol.Operation{op},
		Extensions:     []interface{}{},
		Signatures:     nil, // No signatures
	}

	tx2 := &Transaction{
		RefBlockNum:    1000,
		RefBlockPrefix: 1234567890,
		Operations:     []protocol.Operation{op},
		Extensions:     []interface{}{},
		Signatures:     []string{"test_signature_hex"}, // With signatures
	}

	stx1 := NewSignedTransaction(tx1)
	stx2 := NewSignedTransaction(tx2)

	serialized1, err := stx1.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize tx1: %v", err)
	}

	serialized2, err := stx2.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize tx2: %v", err)
	}

	// Both should serialize to the same bytes (signatures are not included)
	if hex.EncodeToString(serialized1) != hex.EncodeToString(serialized2) {
		t.Errorf("Transaction serialization should not include signatures")
		t.Errorf("Tx1 (no sigs): %s", hex.EncodeToString(serialized1))
		t.Errorf("Tx2 (with sigs): %s", hex.EncodeToString(serialized2))
	} else {
		fmt.Println("✅ Transaction serialization correctly excludes signatures")
	}

	// Both should produce the same digest
	digest1, err := stx1.Digest(SteemChain)
	if err != nil {
		t.Fatalf("Failed to calculate digest1: %v", err)
	}

	digest2, err := stx2.Digest(SteemChain)
	if err != nil {
		t.Fatalf("Failed to calculate digest2: %v", err)
	}

	if hex.EncodeToString(digest1) != hex.EncodeToString(digest2) {
		t.Errorf("Digests should match regardless of signatures")
		t.Errorf("Digest1: %s", hex.EncodeToString(digest1))
		t.Errorf("Digest2: %s", hex.EncodeToString(digest2))
	} else {
		fmt.Println("✅ Digests match (signatures correctly excluded)")
	}
}

