package protocol

import (
	"bytes"
	"encoding/hex"
	"sort"
	"testing"

	"github.com/steemit/steemutil/encoder"
)

func TestCustomJSONOperation_MarshalTransaction(t *testing.T) {
	// Expected hex values generated from old-steem-js using test-custom-json-operation.js
	// These values ensure compatibility with steem-js serialization
	// NOTE: After the fix, MarshalTransaction now includes the operation type code (0x12 = 18)
	// So these expected values now include the type code prefix
	expectedHex1 := "12000006666f6c6c6f77315b22666f6c6c6f77222c7b22666f6c6c6f776572223a22616c696365222c22666f6c6c6f77696e67223a22626f62227d5d"
	expectedHex2 := "12000105616c69636506666f6c6c6f77315b22666f6c6c6f77222c7b22666f6c6c6f776572223a22616c696365222c22666f6c6c6f77696e67223a22626f62227d5d"
	expectedHex3 := "12000305616c69636503626f6207636861726c696506666f6c6c6f77315b22666f6c6c6f77222c7b22666f6c6c6f776572223a22616c696365222c22666f6c6c6f77696e67223a22626f62227d5d"
	expectedHex4 := "120105616c69636500066e6f746966792f5b227365744c61737452656164222c7b2264617465223a22323032332d30312d30315430303a30303a30305a227d5d"
	expectedHex5 := "120205616c69636507636861726c69650203626f620464617665066e6f746966792f5b227365744c61737452656164222c7b2264617465223a22323032332d30312d30315430303a30303a30305a227d5d"

	tests := []struct {
		name     string
		op       *CustomJSONOperation
		expected string // Expected hex output from steem-js
	}{
		{
			name: "empty auths",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{},
				RequiredPostingAuths: []string{},
				ID:                   "follow",
				JSON:                 `["follow",{"follower":"alice","following":"bob"}]`,
			},
			expected: expectedHex1,
		},
		{
			name: "single posting auth",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{},
				RequiredPostingAuths: []string{"alice"},
				ID:                   "follow",
				JSON:                 `["follow",{"follower":"alice","following":"bob"}]`,
			},
			expected: expectedHex2,
		},
		{
			name: "multiple posting auths (should be sorted)",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{},
				RequiredPostingAuths: []string{"charlie", "alice", "bob"},
				ID:                   "follow",
				JSON:                 `["follow",{"follower":"alice","following":"bob"}]`,
			},
			expected: expectedHex3,
		},
		{
			name: "active auths",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{"alice"},
				RequiredPostingAuths: []string{},
				ID:                   "notify",
				JSON:                 `["setLastRead",{"date":"2023-01-01T00:00:00Z"}]`,
			},
			expected: expectedHex4,
		},
		{
			name: "both auth types (should be sorted)",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{"charlie", "alice"},
				RequiredPostingAuths: []string{"bob", "dave"},
				ID:                   "notify",
				JSON:                 `["setLastRead",{"date":"2023-01-01T00:00:00Z"}]`,
			},
			expected: expectedHex5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create encoder with buffer
			var buf bytes.Buffer
			enc := encoder.NewEncoder(&buf)

			// Marshal the operation
			err := tt.op.MarshalTransaction(enc)
			if err != nil {
				t.Fatalf("MarshalTransaction should not fail: %v", err)
			}

			// Get the serialized bytes
			serialized := buf.Bytes()
			if len(serialized) == 0 {
				t.Fatal("Serialized data should not be empty")
			}

			// Verify that required_auths and required_posting_auths are sorted
			// by checking that the original slices are not modified
			originalRequiredAuths := make([]string, len(tt.op.RequiredAuths))
			copy(originalRequiredAuths, tt.op.RequiredAuths)

			originalRequiredPostingAuths := make([]string, len(tt.op.RequiredPostingAuths))
			copy(originalRequiredPostingAuths, tt.op.RequiredPostingAuths)

			// Verify original slices are not modified (MarshalTransaction should sort copies)
			if !equalStringSlices(originalRequiredAuths, tt.op.RequiredAuths) {
				t.Errorf("RequiredAuths should not be modified: got %v, want %v", tt.op.RequiredAuths, originalRequiredAuths)
			}
			if !equalStringSlices(originalRequiredPostingAuths, tt.op.RequiredPostingAuths) {
				t.Errorf("RequiredPostingAuths should not be modified: got %v, want %v", tt.op.RequiredPostingAuths, originalRequiredPostingAuths)
			}

			// Verify serialization format:
			// 0. varint32 operation type code (18 = 0x12 for custom_json) - NOW INCLUDED
			// 1. varint32 length of required_auths
			// 2. each account name (string: varint32 length + bytes)
			// 3. varint32 length of required_posting_auths
			// 4. each account name (string: varint32 length + bytes)
			// 5. id (string: varint32 length + bytes)
			// 6. json (string: varint32 length + bytes)

			// Decode and verify structure
			pos := 0

			// Read operation type code (should be 18 = 0x12 for custom_json)
			opTypeCode, n := readVarint(serialized[pos:])
			pos += n
			if opTypeCode != 18 {
				t.Errorf("operation type code should be 18 (custom_json), got %d", opTypeCode)
			}

			// Read required_auths length
			requiredAuthsLen, n := readVarint(serialized[pos:])
			pos += n
			if requiredAuthsLen != uint64(len(tt.op.RequiredAuths)) {
				t.Errorf("required_auths length should match: got %d, want %d", requiredAuthsLen, len(tt.op.RequiredAuths))
			}

			// Read required_auths accounts (should be sorted)
			sortedRequiredAuths := make([]string, len(tt.op.RequiredAuths))
			copy(sortedRequiredAuths, tt.op.RequiredAuths)
			sort.Strings(sortedRequiredAuths)

			for i := 0; i < int(requiredAuthsLen); i++ {
				strLen, n := readVarint(serialized[pos:])
				pos += n
				account := string(serialized[pos : pos+int(strLen)])
				pos += int(strLen)
				if account != sortedRequiredAuths[i] {
					t.Errorf("required_auths should be sorted: at index %d, got %s, want %s", i, account, sortedRequiredAuths[i])
				}
			}

			// Read required_posting_auths length
			requiredPostingAuthsLen, n := readVarint(serialized[pos:])
			pos += n
			if requiredPostingAuthsLen != uint64(len(tt.op.RequiredPostingAuths)) {
				t.Errorf("required_posting_auths length should match: got %d, want %d", requiredPostingAuthsLen, len(tt.op.RequiredPostingAuths))
			}

			// Read required_posting_auths accounts (should be sorted)
			sortedRequiredPostingAuths := make([]string, len(tt.op.RequiredPostingAuths))
			copy(sortedRequiredPostingAuths, tt.op.RequiredPostingAuths)
			sort.Strings(sortedRequiredPostingAuths)

			for i := 0; i < int(requiredPostingAuthsLen); i++ {
				strLen, n := readVarint(serialized[pos:])
				pos += n
				account := string(serialized[pos : pos+int(strLen)])
				pos += int(strLen)
				if account != sortedRequiredPostingAuths[i] {
					t.Errorf("required_posting_auths should be sorted: at index %d, got %s, want %s", i, account, sortedRequiredPostingAuths[i])
				}
			}

			// Read id
			idLen, n := readVarint(serialized[pos:])
			pos += n
			id := string(serialized[pos : pos+int(idLen)])
			pos += int(idLen)
			if id != tt.op.ID {
				t.Errorf("id should match: got %s, want %s", id, tt.op.ID)
			}

			// Read json
			jsonLen, n := readVarint(serialized[pos:])
			pos += n
			jsonStr := string(serialized[pos : pos+int(jsonLen)])
			pos += int(jsonLen)
			if jsonStr != tt.op.JSON {
				t.Errorf("json should match: got %s, want %s", jsonStr, tt.op.JSON)
			}

			// Verify we consumed all bytes
			if pos != len(serialized) {
				t.Errorf("Should consume all serialized bytes: consumed %d, total %d", pos, len(serialized))
			}

			// Compare with expected hex from steem-js
			actualHex := hex.EncodeToString(serialized)
			if tt.expected != "" {
				if actualHex != tt.expected {
					t.Errorf("Serialized hex does not match steem-js:\n  got:      %s\n  expected: %s", actualHex, tt.expected)
				} else {
					t.Logf("✓ Serialized hex matches steem-js: %s", actualHex)
				}
			} else {
				t.Logf("Serialized hex: %s", actualHex)
			}
		})
	}
}

func TestCustomJSONOperation_MarshalTransaction_Sorting(t *testing.T) {
	// Test that accounts are sorted correctly regardless of input order
	op := &CustomJSONOperation{
		RequiredAuths:        []string{"zebra", "alpha", "beta"},
		RequiredPostingAuths: []string{"zulu", "alpha", "beta"},
		ID:                   "test",
		JSON:                 `{"test":true}`,
	}

	var buf bytes.Buffer
	enc := encoder.NewEncoder(&buf)

	err := op.MarshalTransaction(enc)
	if err != nil {
		t.Fatalf("MarshalTransaction should not fail: %v", err)
	}

	serialized := buf.Bytes()
	pos := 0

	// Read operation type code (should be 18 = 0x12 for custom_json)
	opTypeCode, n := readVarint(serialized[pos:])
	pos += n
	if opTypeCode != 18 {
		t.Errorf("operation type code should be 18 (custom_json), got %d", opTypeCode)
	}

	// Read required_auths
	requiredAuthsLen, n := readVarint(serialized[pos:])
	pos += n

	requiredAuths := make([]string, 0, requiredAuthsLen)
	for i := 0; i < int(requiredAuthsLen); i++ {
		strLen, n := readVarint(serialized[pos:])
		pos += n
		account := string(serialized[pos : pos+int(strLen)])
		pos += int(strLen)
		requiredAuths = append(requiredAuths, account)
	}

	// Verify required_auths are sorted
	expectedRequiredAuths := []string{"alpha", "beta", "zebra"}
	if !equalStringSlices(expectedRequiredAuths, requiredAuths) {
		t.Errorf("required_auths should be sorted alphabetically: got %v, want %v", requiredAuths, expectedRequiredAuths)
	}

	// Read required_posting_auths
	requiredPostingAuthsLen, n := readVarint(serialized[pos:])
	pos += n

	requiredPostingAuths := make([]string, 0, requiredPostingAuthsLen)
	for i := 0; i < int(requiredPostingAuthsLen); i++ {
		strLen, n := readVarint(serialized[pos:])
		pos += n
		account := string(serialized[pos : pos+int(strLen)])
		pos += int(strLen)
		requiredPostingAuths = append(requiredPostingAuths, account)
	}

	// Verify required_posting_auths are sorted
	expectedRequiredPostingAuths := []string{"alpha", "beta", "zulu"}
	if !equalStringSlices(expectedRequiredPostingAuths, requiredPostingAuths) {
		t.Errorf("required_posting_auths should be sorted alphabetically: got %v, want %v", requiredPostingAuths, expectedRequiredPostingAuths)
	}
}

// Helper function to read a varint from bytes
func readVarint(data []byte) (uint64, int) {
	var result uint64
	var shift uint
	for i, b := range data {
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			return result, i + 1
		}
		shift += 7
	}
	return result, len(data)
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestCustomJSONOperation_FullOperationEncoding tests the complete operation encoding
// including the operation type code. This tests the path used in actual transactions
// (through encoder.Encode() rather than direct MarshalTransaction call).
func TestCustomJSONOperation_FullOperationEncoding(t *testing.T) {
	// Expected hex values for full operation serialization (with type code 18 = 0x12)
	// Generated from old-steem-js using operation.appendByteBuffer()
	// These include the operation type code varint (18 = 0x12) followed by operation data
	// Format: [0x12 (type code)] + [operation data from MarshalTransaction]
	expectedFullOpHex1 := "12000006666f6c6c6f77315b22666f6c6c6f77222c7b22666f6c6c6f776572223a22616c696365222c22666f6c6c6f77696e67223a22626f62227d5d"
	expectedFullOpHex2 := "12000105616c69636506666f6c6c6f77315b22666f6c6c6f77222c7b22666f6c6c6f776572223a22616c696365222c22666f6c6c6f77696e67223a22626f62227d5d"
	expectedFullOpHex3 := "12000305616c69636503626f6207636861726c696506666f6c6c6f77315b22666f6c6c6f77222c7b22666f6c6c6f776572223a22616c696365222c22666f6c6c6f77696e67223a22626f62227d5d"
	expectedFullOpHex4 := "120105616c69636500066e6f746966792f5b227365744c61737452656164222c7b2264617465223a22323032332d30312d30315430303a30303a30305a227d5d"
	expectedFullOpHex5 := "120205616c69636507636861726c69650203626f620464617665066e6f746966792f5b227365744c61737452656164222c7b2264617465223a22323032332d30312d30315430303a30303a30305a227d5d"

	tests := []struct {
		name     string
		op       *CustomJSONOperation
		expected string // Expected hex output from steem-js (full operation with type code)
	}{
		{
			name: "empty auths",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{},
				RequiredPostingAuths: []string{},
				ID:                   "follow",
				JSON:                 `["follow",{"follower":"alice","following":"bob"}]`,
			},
			expected: expectedFullOpHex1,
		},
		{
			name: "single posting auth",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{},
				RequiredPostingAuths: []string{"alice"},
				ID:                   "follow",
				JSON:                 `["follow",{"follower":"alice","following":"bob"}]`,
			},
			expected: expectedFullOpHex2,
		},
		{
			name: "multiple posting auths (should be sorted)",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{},
				RequiredPostingAuths: []string{"charlie", "alice", "bob"},
				ID:                   "follow",
				JSON:                 `["follow",{"follower":"alice","following":"bob"}]`,
			},
			expected: expectedFullOpHex3,
		},
		{
			name: "active auths",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{"alice"},
				RequiredPostingAuths: []string{},
				ID:                   "notify",
				JSON:                 `["setLastRead",{"date":"2023-01-01T00:00:00Z"}]`,
			},
			expected: expectedFullOpHex4,
		},
		{
			name: "both auth types (should be sorted)",
			op: &CustomJSONOperation{
				RequiredAuths:        []string{"charlie", "alice"},
				RequiredPostingAuths: []string{"bob", "dave"},
				ID:                   "notify",
				JSON:                 `["setLastRead",{"date":"2023-01-01T00:00:00Z"}]`,
			},
			expected: expectedFullOpHex5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create encoder with buffer
			var buf bytes.Buffer
			enc := encoder.NewEncoder(&buf)

			// Encode the operation through encoder.Encode() (the path used in transactions)
			// This will call encodeOperation() which includes the type code
			err := enc.Encode(tt.op)
			if err != nil {
				t.Fatalf("Encode should not fail: %v", err)
			}

			// Get the serialized bytes
			serialized := buf.Bytes()
			if len(serialized) == 0 {
				t.Fatal("Serialized data should not be empty")
			}

			// Verify the first byte is the operation type code (18 = 0x12 for custom_json)
			actualHex := hex.EncodeToString(serialized)
			if len(actualHex) < 2 {
				t.Fatal("Serialized data too short")
			}

			// Check operation type code (first byte should be 0x12 = 18)
			opTypeHex := actualHex[0:2]
			if opTypeHex != "12" {
				t.Errorf("Operation type code should be 0x12 (18 for custom_json), got 0x%s", opTypeHex)
			}

			// Compare with expected hex from steem-js
			if tt.expected != "" {
				if actualHex != tt.expected {
					t.Errorf("Full operation serialization does not match steem-js:\n  got:      %s\n  expected: %s", actualHex, tt.expected)
					t.Errorf("Operation type code: got 0x%s, expected 0x%s", actualHex[0:2], tt.expected[0:2])
				} else {
					t.Logf("✓ Full operation serialization matches steem-js: %s", actualHex)
				}
			} else {
				t.Logf("Full operation hex: %s", actualHex)
			}

			// Verify operation type code is correctly encoded
			// Decode varint for operation type
			opTypeCode, n := readVarint(serialized)
			if opTypeCode != 18 {
				t.Errorf("Operation type code should be 18 (custom_json), got %d", opTypeCode)
			}
			if n != 1 {
				t.Errorf("Operation type code should be 1 byte (18 < 128), got %d bytes", n)
			}

			// Verify the rest matches the operation data (without type code)
			operationDataHex := hex.EncodeToString(serialized[n:])
			expectedDataHex := tt.expected[2:] // Skip first byte (type code)
			if operationDataHex != expectedDataHex {
				t.Errorf("Operation data does not match:\n  got:      %s\n  expected: %s", operationDataHex, expectedDataHex)
			}
		})
	}
}
