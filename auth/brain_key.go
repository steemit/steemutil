package auth

import (
	"crypto/sha256"
	"strings"

	"github.com/pkg/errors"
	"github.com/steemit/steemutil/wif"
)

// GenerateBrainKey generates a private key from account name, password, and role.
// This implements the brain key generation algorithm used by Steem.
func GenerateBrainKey(name, password, role string) ([]byte, error) {
	// Create seed: name + role + password
	seed := name + role + password

	// Normalize brain key: trim and normalize whitespace
	brainKey := strings.TrimSpace(seed)
	// Replace multiple whitespace with single space
	fields := strings.Fields(brainKey)
	brainKey = strings.Join(fields, " ")

	// Hash the brain key
	hash := sha256.Sum256([]byte(brainKey))
	return hash[:], nil
}

// ToWif generates a WIF from account name, password, and role.
func ToWif(name, password, role string) (string, error) {
	// Generate brain key
	privKeyBytes, err := GenerateBrainKey(name, password, role)
	if err != nil {
		return "", err
	}

	// Create private key from bytes
	privKey := &wif.PrivateKey{}
	if err := privKey.FromByte(privKeyBytes); err != nil {
		return "", errors.Wrap(err, "failed to create private key from brain key")
	}

	return privKey.ToWif(), nil
}
