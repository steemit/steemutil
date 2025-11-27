package wif

import (
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/pkg/errors"
)

// SignSha256 signs a SHA256 message hash using the private key.
// This method is specifically designed for signing arbitrary message hashes,
// such as those used in RPC authentication.
func (pk *PrivateKey) SignSha256(message []byte) ([]byte, error) {
	if pk.Raw == nil || pk.Raw.PrivKey == nil {
		return nil, errors.New("private key not initialized")
	}

	// Sign the message hash directly
	signature, err := ecdsa.SignCompact(pk.Raw.PrivKey, message, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compact signature")
	}

	return signature, nil
}

// VerifySha256 verifies a SHA256 message signature using the public key.
// This method is used to verify signatures created with SignSha256.
func (pk *PublicKey) VerifySha256(message, signature []byte) bool {
	if pk.Raw == nil {
		return false
	}

	// Parse the compact signature
	parsedSig, wasCompressed, err := ecdsa.RecoverCompact(signature, message)
	if err != nil {
		return false
	}

	// Verify that the recovered public key matches this public key
	if !wasCompressed {
		return false // We expect compressed signatures
	}

	// Compare the recovered public key with our public key
	// Convert both to the same format for comparison
	recoveredBytes := parsedSig.SerializeCompressed()
	ourBytes := pk.Raw.SerializeCompressed()
	
	if len(recoveredBytes) != len(ourBytes) {
		return false
	}
	
	for i := range recoveredBytes {
		if recoveredBytes[i] != ourBytes[i] {
			return false
		}
	}
	
	return true
}

// RecoverPublicKeyFromSignature recovers the public key from a signature and message hash.
// This is useful for verifying signatures when you only have the signature and message.
func RecoverPublicKeyFromSignature(message, signature []byte) (*PublicKey, error) {
	// Recover the public key from the compact signature
	recoveredKey, wasCompressed, err := ecdsa.RecoverCompact(signature, message)
	if err != nil {
		return nil, errors.Wrap(err, "failed to recover public key from signature")
	}

	if !wasCompressed {
		return nil, errors.New("signature was not created with compressed public key")
	}

	// Convert the recovered btcec public key to secp256k1 format
	recoveredBytes := recoveredKey.SerializeCompressed()
	secp256k1PubKey, err := secp256k1.ParsePubKey(recoveredBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert recovered public key")
	}

	// Create a PublicKey instance
	pubKey := &PublicKey{
		Raw: secp256k1PubKey,
	}

	return pubKey, nil
}

// SignMessage signs an arbitrary message (not a hash) using the private key.
// The message will be hashed with SHA256 before signing.
func (pk *PrivateKey) SignMessage(message []byte) ([]byte, error) {
	// Hash the message with SHA256
	hash := sha256.Sum256(message)
	
	// Sign the hash
	return pk.SignSha256(hash[:])
}

// VerifyMessage verifies a message signature using the public key.
// The message will be hashed with SHA256 before verification.
func (pk *PublicKey) VerifyMessage(message, signature []byte) bool {
	// Hash the message with SHA256
	hash := sha256.Sum256(message)
	
	// Verify the signature
	return pk.VerifySha256(hash[:], signature)
}
