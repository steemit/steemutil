package auth

import (
	"testing"

	"github.com/steemit/steemutil/wif"
)

func TestToWif(t *testing.T) {
	name := "testuser"
	password := "testpassword"
	role := "posting"

	wifStr, err := ToWif(name, password, role)
	if err != nil {
		t.Fatalf("ToWif failed: %v", err)
	}

	if wifStr == "" {
		t.Error("ToWif returned empty string")
	}

	// Verify it's a valid WIF
	privKey := &wif.PrivateKey{}
	if err := privKey.FromWif(wifStr); err != nil {
		t.Errorf("Generated WIF is invalid: %v", err)
	}
}

func TestWifToPublic(t *testing.T) {
	// Use a known WIF for testing
	testWif := "5JRaypasxMx1L97ZUX7YuC5Psb5EAbF821kkAGtBj7xCJFQcbLg"
	
	pubKey, err := WifToPublic(testWif)
	if err != nil {
		t.Fatalf("WifToPublic failed: %v", err)
	}

	if pubKey == "" {
		t.Error("WifToPublic returned empty string")
	}

	// Verify it's a valid public key
	pubKeyObj := &wif.PublicKey{}
	if err := pubKeyObj.FromStr(pubKey); err != nil {
		t.Errorf("Generated public key is invalid: %v", err)
	}
}

func TestIsWif(t *testing.T) {
	validWif := "5JRaypasxMx1L97ZUX7YuC5Psb5EAbF821kkAGtBj7xCJFQcbLg"
	invalidWif := "invalidwif"

	if !IsWif(validWif) {
		t.Error("IsWif should return true for valid WIF")
	}

	if IsWif(invalidWif) {
		t.Error("IsWif should return false for invalid WIF")
	}
}

func TestWifIsValid(t *testing.T) {
	testWif := "5JRaypasxMx1L97ZUX7YuC5Psb5EAbF821kkAGtBj7xCJFQcbLg"
	
	pubKey, err := WifToPublic(testWif)
	if err != nil {
		t.Fatalf("WifToPublic failed: %v", err)
	}

	if !WifIsValid(testWif, pubKey) {
		t.Error("WifIsValid should return true for matching WIF and public key")
	}

	if WifIsValid(testWif, "STMinvalid") {
		t.Error("WifIsValid should return false for non-matching public key")
	}
}

func TestIsPubkey(t *testing.T) {
	validPubKey := "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA"
	invalidPubKey := "invalidpubkey"

	if !IsPubkey(validPubKey) {
		t.Error("IsPubkey should return true for valid public key")
	}

	if IsPubkey(invalidPubKey) {
		t.Error("IsPubkey should return false for invalid public key")
	}
}

func TestGenerateKeys(t *testing.T) {
	name := "testuser"
	password := "testpassword"
	roles := []string{"owner", "active", "posting", "memo"}

	pubKeys, err := GenerateKeys(name, password, roles)
	if err != nil {
		t.Fatalf("GenerateKeys failed: %v", err)
	}

	if len(pubKeys) != len(roles) {
		t.Errorf("GenerateKeys returned %d keys, expected %d", len(pubKeys), len(roles))
	}

	// Verify all roles have keys
	for _, role := range roles {
		if _, ok := pubKeys[role]; !ok {
			t.Errorf("GenerateKeys missing key for role: %s", role)
		}
	}

	// Verify keys are valid public keys
	for role, pubKey := range pubKeys {
		if !IsPubkey(pubKey) {
			t.Errorf("GenerateKeys returned invalid public key for role %s: %s", role, pubKey)
		}
	}
}

func TestGetPrivateKeys(t *testing.T) {
	name := "testuser"
	password := "testpassword"
	roles := []string{"owner", "active", "posting", "memo"}

	keys, err := GetPrivateKeys(name, password, roles)
	if err != nil {
		t.Fatalf("GetPrivateKeys failed: %v", err)
	}

	// Should have both private keys and public keys
	expectedCount := len(roles) * 2
	if len(keys) != expectedCount {
		t.Errorf("GetPrivateKeys returned %d keys, expected %d", len(keys), expectedCount)
	}

	// Verify private keys are valid WIFs
	for _, role := range roles {
		wifStr, ok := keys[role]
		if !ok {
			t.Errorf("GetPrivateKeys missing private key for role: %s", role)
			continue
		}
		if !IsWif(wifStr) {
			t.Errorf("GetPrivateKeys returned invalid WIF for role %s: %s", role, wifStr)
		}

		// Verify corresponding public key
		pubKeyKey := role + "Pubkey"
		pubKey, ok := keys[pubKeyKey]
		if !ok {
			t.Errorf("GetPrivateKeys missing public key for role: %s", role)
			continue
		}

		// Verify WIF and public key match
		if !WifIsValid(wifStr, pubKey) {
			t.Errorf("GetPrivateKeys: WIF and public key don't match for role %s", role)
		}
	}
}

func TestGetPrivateKeysDefaultRoles(t *testing.T) {
	name := "testuser"
	password := "testpassword"

	keys, err := GetPrivateKeys(name, password, nil)
	if err != nil {
		t.Fatalf("GetPrivateKeys failed: %v", err)
	}

	// Default roles: owner, active, posting, memo
	defaultRoles := []string{"owner", "active", "posting", "memo"}
	expectedCount := len(defaultRoles) * 2
	if len(keys) != expectedCount {
		t.Errorf("GetPrivateKeys with nil roles returned %d keys, expected %d", len(keys), expectedCount)
	}
}

