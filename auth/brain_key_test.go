package auth

import (
	"crypto/sha256"
	"strings"
	"testing"
)

func TestGenerateBrainKey(t *testing.T) {
	name := "testuser"
	password := "testpassword"
	role := "posting"

	brainKey, err := GenerateBrainKey(name, password, role)
	if err != nil {
		t.Fatalf("GenerateBrainKey failed: %v", err)
	}

	if len(brainKey) != 32 {
		t.Errorf("GenerateBrainKey returned key of length %d, expected 32", len(brainKey))
	}

	// Test that same input produces same output
	brainKey2, err := GenerateBrainKey(name, password, role)
	if err != nil {
		t.Fatalf("GenerateBrainKey failed on second call: %v", err)
	}

	if string(brainKey) != string(brainKey2) {
		t.Error("GenerateBrainKey should produce deterministic output")
	}
}

func TestGenerateBrainKeyNormalization(t *testing.T) {
	// Test that whitespace normalization works
	name1 := "testuser"
	password1 := "test  password"  // Multiple spaces
	role1 := "posting"

	name2 := "testuser"
	password2 := "test password"   // Single space
	role2 := "posting"

	brainKey1, err1 := GenerateBrainKey(name1, password1, role1)
	if err1 != nil {
		t.Fatalf("GenerateBrainKey failed: %v", err1)
	}

	brainKey2, err2 := GenerateBrainKey(name2, password2, role2)
	if err2 != nil {
		t.Fatalf("GenerateBrainKey failed: %v", err2)
	}

	// Should produce same result after normalization
	if string(brainKey1) != string(brainKey2) {
		t.Error("GenerateBrainKey should normalize whitespace")
	}
}

func TestGenerateBrainKeyDifferentRoles(t *testing.T) {
	name := "testuser"
	password := "testpassword"

	roles := []string{"owner", "active", "posting", "memo"}
	keys := make(map[string][]byte)

	for _, role := range roles {
		key, err := GenerateBrainKey(name, password, role)
		if err != nil {
			t.Fatalf("GenerateBrainKey failed for role %s: %v", role, err)
		}
		keys[role] = key
	}

	// All keys should be different
	for i, role1 := range roles {
		for j, role2 := range roles {
			if i != j && string(keys[role1]) == string(keys[role2]) {
				t.Errorf("GenerateBrainKey produced same key for different roles: %s and %s", role1, role2)
			}
		}
	}
}

func TestBrainKeyAlgorithm(t *testing.T) {
	// Test the brain key algorithm matches expected behavior
	name := "test"
	password := "password"
	role := "posting"

	seed := name + role + password
	brainKey := strings.TrimSpace(seed)
	fields := strings.Fields(brainKey)
	brainKey = strings.Join(fields, " ")

	hash := sha256.Sum256([]byte(brainKey))
	expectedKey := hash[:]

	actualKey, err := GenerateBrainKey(name, password, role)
	if err != nil {
		t.Fatalf("GenerateBrainKey failed: %v", err)
	}

	if string(actualKey) != string(expectedKey) {
		t.Error("GenerateBrainKey algorithm doesn't match expected implementation")
	}
}

