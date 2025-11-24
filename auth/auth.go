package auth

import (
	"crypto/sha256"

	"github.com/pkg/errors"
	"github.com/steemit/steemutil/transaction"
	"github.com/steemit/steemutil/wif"
)

// Verify verifies if the account name and password match the given authorities.
// It returns true if at least one role's public key matches.
func Verify(name, password string, auths map[string]interface{}) (bool, error) {
	// Extract roles from auths
	roles := make([]string, 0, len(auths))
	for role := range auths {
		roles = append(roles, role)
	}

	// Generate public keys for all roles
	pubKeys, err := GenerateKeys(name, password, roles)
	if err != nil {
		return false, err
	}

	// Check if any role's public key matches
	for role, pubKey := range pubKeys {
		if authData, ok := auths[role].([]interface{}); ok && len(authData) > 0 {
			if keyAuths, ok := authData[0].([]interface{}); ok && len(keyAuths) > 0 {
				if expectedPubKey, ok := keyAuths[0].(string); ok {
					if expectedPubKey == pubKey {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

// GenerateKeys generates public keys for the given roles from account name and password.
func GenerateKeys(name, password string, roles []string) (map[string]string, error) {
	pubKeys := make(map[string]string, len(roles))

	for _, role := range roles {
		// Generate WIF from brain key
		privWif, err := ToWif(name, password, role)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate WIF for role %s", role)
		}

		// Convert WIF to public key
		pubKey, err := WifToPublic(privWif)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert WIF to public key for role %s", role)
		}

		pubKeys[role] = pubKey
	}

	return pubKeys, nil
}

// GetPrivateKeys returns private keys and public keys for the given roles.
// Default roles are: "owner", "active", "posting", "memo"
func GetPrivateKeys(name, password string, roles []string) (map[string]string, error) {
	if len(roles) == 0 {
		roles = []string{"owner", "active", "posting", "memo"}
	}

	keys := make(map[string]string, len(roles)*2)

	for _, role := range roles {
		// Generate WIF
		privWif, err := ToWif(name, password, role)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate WIF for role %s", role)
		}
		keys[role] = privWif

		// Generate public key
		pubKey, err := WifToPublic(privWif)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert WIF to public key for role %s", role)
		}
		keys[role+"Pubkey"] = pubKey
	}

	return keys, nil
}

// IsWif checks if the given string is a valid WIF format.
func IsWif(privWif string) bool {
	privKey := &wif.PrivateKey{}
	err := privKey.FromWif(privWif)
	return err == nil
}

// WifIsValid checks if the given WIF corresponds to the given public key.
func WifIsValid(privWif, pubKey string) bool {
	derivedPubKey, err := WifToPublic(privWif)
	if err != nil {
		return false
	}
	return derivedPubKey == pubKey
}

// WifToPublic converts a WIF to a public key string.
func WifToPublic(privWif string) (string, error) {
	privKey := &wif.PrivateKey{}
	if err := privKey.FromWif(privWif); err != nil {
		return "", errors.Wrap(err, "failed to decode WIF")
	}
	return privKey.ToPubKeyStr(), nil
}

// IsPubkey checks if the given string is a valid public key format.
func IsPubkey(pubkey string) bool {
	pubKey := &wif.PublicKey{}
	err := pubKey.FromStr(pubkey)
	return err == nil
}

// SignTransaction signs a transaction with the given private keys.
// The keys parameter should be a map of role -> WIF string.
func SignTransaction(tx *transaction.SignedTransaction, keys map[string]string, chain *transaction.Chain) error {
	// Collect private keys
	privKeys := make([]*wif.PrivateKey, 0, len(keys))
	for _, wifStr := range keys {
		privKey := &wif.PrivateKey{}
		if err := privKey.FromWif(wifStr); err != nil {
			return errors.Wrap(err, "failed to decode WIF")
		}
		privKeys = append(privKeys, privKey)
	}

	// Sign the transaction
	return tx.Sign(privKeys, chain)
}

// VerifyChecksum verifies the checksum of a WIF.
func VerifyChecksum(wifBytes []byte) bool {
	if len(wifBytes) < 5 {
		return false
	}

	privKey := wifBytes[:len(wifBytes)-4]
	checksum := wifBytes[len(wifBytes)-4:]

	// Double SHA256
	hash1 := sha256.Sum256(privKey)
	hash2 := sha256.Sum256(hash1[:])
	newChecksum := hash2[:4]

	// Compare checksums
	for i := 0; i < 4; i++ {
		if checksum[i] != newChecksum[i] {
			return false
		}
	}

	return true
}
