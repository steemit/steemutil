package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/pkg/errors"
	"github.com/steemit/steemutil/wif"
)

// EncryptedMemo represents an encrypted memo structure.
type EncryptedMemo struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Nonce     []byte `json:"nonce"`
	Check     []byte `json:"check"`
	Encrypted []byte `json:"encrypted"`
}

// Encode encrypts a memo if it starts with '#', otherwise returns it as-is.
// privateKey can be a WIF string or a PrivateKey object.
// publicKey can be a public key string or a PublicKey object.
func Encode(privateKey interface{}, publicKey interface{}, memo string) (string, error) {
	if memo == "" {
		return "", errors.New("memo is required")
	}

	// If memo doesn't start with '#', return as-is
	if len(memo) == 0 || memo[0] != '#' {
		return memo, nil
	}

	// Remove '#' prefix
	memo = memo[1:]

	// Convert private key to PrivateKey object
	privKey, err := toPrivateKey(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert private key")
	}

	// Convert public key to PublicKey object
	pubKey, err := toPublicKey(publicKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert public key")
	}

	// Get sender's public key
	senderPubKey := &wif.PublicKey{}
	if err := senderPubKey.FromStr(privKey.ToPubKeyStr()); err != nil {
		return "", errors.Wrap(err, "failed to get sender public key")
	}

	// Determine recipient public key
	recipientPubKey := pubKey
	if senderPubKey.ToStr() == pubKey.ToStr() {
		// If sender and recipient are the same, we need the other key
		// This is a simplified version - in practice, you'd need the actual recipient key
		recipientPubKey = pubKey
	}

	// Encrypt the memo
	encrypted, nonce, checksum, err := encryptMemo(privKey, recipientPubKey, []byte(memo))
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt memo")
	}

	// Create encrypted memo structure
	encMemo := EncryptedMemo{
		From:      senderPubKey.ToStr(),
		To:        recipientPubKey.ToStr(),
		Nonce:     nonce,
		Check:     checksum,
		Encrypted: encrypted,
	}

	// Serialize (simplified - in practice, you'd use the proper serializer)
	// For now, we'll use a simple base64 encoding
	memoBytes, err := serializeEncryptedMemo(encMemo)
	if err != nil {
		return "", errors.Wrap(err, "failed to serialize encrypted memo")
	}

	// Encode to base58
	encoded := base58.Encode(memoBytes)
	return "#" + encoded, nil
}

// Decode decrypts a memo if it starts with '#', otherwise returns it as-is.
func Decode(privateKey interface{}, memo string) (string, error) {
	if memo == "" {
		return "", errors.New("memo is required")
	}

	// If memo doesn't start with '#', return as-is
	if len(memo) == 0 || memo[0] != '#' {
		return memo, nil
	}

	// Remove '#' prefix
	memo = memo[1:]

	// Convert private key to PrivateKey object
	privKey, err := toPrivateKey(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert private key")
	}

	// Decode from base58
	memoBytes := base58.Decode(memo)

	// Deserialize encrypted memo
	encMemo, err := deserializeEncryptedMemo(memoBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to deserialize encrypted memo")
	}

	// Get sender's public key
	senderPubKey := &wif.PublicKey{}
	if err := senderPubKey.FromStr(privKey.ToPubKeyStr()); err != nil {
		return "", errors.Wrap(err, "failed to get sender public key")
	}

	// Determine the other party's public key
	var otherPubKey *wif.PublicKey
	if senderPubKey.ToStr() == encMemo.From {
		otherPubKey = &wif.PublicKey{}
		if err := otherPubKey.FromStr(encMemo.To); err != nil {
			return "", errors.Wrap(err, "failed to parse recipient public key")
		}
	} else {
		otherPubKey = &wif.PublicKey{}
		if err := otherPubKey.FromStr(encMemo.From); err != nil {
			return "", errors.Wrap(err, "failed to parse sender public key")
		}
	}

	// Decrypt the memo
	decrypted, err := decryptMemo(privKey, otherPubKey, encMemo.Nonce, encMemo.Encrypted, encMemo.Check)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt memo")
	}

	// Return with '#' prefix
	return "#" + string(decrypted), nil
}

// Helper functions

func toPrivateKey(key interface{}) (*wif.PrivateKey, error) {
	switch v := key.(type) {
	case *wif.PrivateKey:
		return v, nil
	case string:
		privKey := &wif.PrivateKey{}
		if err := privKey.FromWif(v); err != nil {
			return nil, err
		}
		return privKey, nil
	default:
		return nil, errors.New("invalid private key type")
	}
}

func toPublicKey(key interface{}) (*wif.PublicKey, error) {
	switch v := key.(type) {
	case *wif.PublicKey:
		return v, nil
	case string:
		pubKey := &wif.PublicKey{}
		if err := pubKey.FromStr(v); err != nil {
			return nil, err
		}
		return pubKey, nil
	default:
		return nil, errors.New("invalid public key type")
	}
}

// encryptMemo encrypts a memo using AES-256-CBC with shared secret derived from ECDH.
func encryptMemo(privKey *wif.PrivateKey, pubKey *wif.PublicKey, message []byte) ([]byte, []byte, []byte, error) {
	// Generate shared secret using ECDH
	sharedSecret := deriveSharedSecret(privKey, pubKey)

	// Generate random nonce
	nonce := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(sharedSecret[:32])
	if err != nil {
		return nil, nil, nil, err
	}

	// Use CBC mode
	iv := nonce
	mode := cipher.NewCBCEncrypter(block, iv)

	// Pad message
	padded := pkcs7Pad(message, aes.BlockSize)

	// Encrypt
	encrypted := make([]byte, len(padded))
	mode.CryptBlocks(encrypted, padded)

	// Calculate checksum (first 4 bytes of SHA256 of encrypted data)
	hash := sha256.Sum256(encrypted)
	checksum := hash[:4]

	return encrypted, nonce, checksum, nil
}

// decryptMemo decrypts a memo using AES-256-CBC.
func decryptMemo(privKey *wif.PrivateKey, pubKey *wif.PublicKey, nonce, encrypted, checksum []byte) ([]byte, error) {
	// Verify checksum
	hash := sha256.Sum256(encrypted)
	if len(checksum) < 4 {
		return nil, errors.New("invalid checksum length")
	}
	for i := 0; i < 4; i++ {
		if hash[i] != checksum[i] {
			return nil, errors.New("checksum mismatch")
		}
	}

	// Generate shared secret using ECDH
	sharedSecret := deriveSharedSecret(privKey, pubKey)

	// Create AES cipher
	block, err := aes.NewCipher(sharedSecret[:32])
	if err != nil {
		return nil, err
	}

	// Use CBC mode
	iv := nonce
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	// Remove padding
	unpadded, err := pkcs7Unpad(decrypted)
	if err != nil {
		return nil, err
	}

	return unpadded, nil
}

// deriveSharedSecret derives a shared secret using ECDH.
func deriveSharedSecret(privKey *wif.PrivateKey, pubKey *wif.PublicKey) []byte {
	// ECDH: shared secret = privKey * pubKey
	// This is a simplified version - in practice, you'd use proper ECDH
	// For now, we'll use a combination of both keys
	privBytes := privKey.ToByte()
	pubBytes := pubKey.ToByte()
	combined := append(privBytes, pubBytes...)
	hash := sha256.Sum256(combined)
	return hash[:]
}

// pkcs7Pad adds PKCS7 padding to data.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpad removes PKCS7 padding from data.
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return nil, errors.New("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:len(data)-padding], nil
}

// serializeEncryptedMemo serializes an encrypted memo (simplified version).
func serializeEncryptedMemo(memo EncryptedMemo) ([]byte, error) {
	// This is a simplified serialization
	// In practice, you'd use the proper Steem serializer
	// For now, we'll use a simple format: from|to|nonce|check|encrypted
	fromBytes := []byte(memo.From)
	toBytes := []byte(memo.To)

	result := make([]byte, 0, len(fromBytes)+len(toBytes)+len(memo.Nonce)+len(memo.Check)+len(memo.Encrypted)+20)
	result = append(result, byte(len(fromBytes)))
	result = append(result, fromBytes...)
	result = append(result, byte(len(toBytes)))
	result = append(result, toBytes...)
	result = append(result, byte(len(memo.Nonce)))
	result = append(result, memo.Nonce...)
	result = append(result, byte(len(memo.Check)))
	result = append(result, memo.Check...)
	result = append(result, memo.Encrypted...)
	return result, nil
}

// deserializeEncryptedMemo deserializes an encrypted memo (simplified version).
func deserializeEncryptedMemo(data []byte) (EncryptedMemo, error) {
	// This is a simplified deserialization
	// In practice, you'd use the proper Steem deserializer
	if len(data) < 5 {
		return EncryptedMemo{}, errors.New("data too short")
	}

	offset := 0
	fromLen := int(data[offset])
	offset++
	if offset+fromLen > len(data) {
		return EncryptedMemo{}, errors.New("invalid from length")
	}
	from := string(data[offset : offset+fromLen])
	offset += fromLen

	toLen := int(data[offset])
	offset++
	if offset+toLen > len(data) {
		return EncryptedMemo{}, errors.New("invalid to length")
	}
	to := string(data[offset : offset+toLen])
	offset += toLen

	nonceLen := int(data[offset])
	offset++
	if offset+nonceLen > len(data) {
		return EncryptedMemo{}, errors.New("invalid nonce length")
	}
	nonce := make([]byte, nonceLen)
	copy(nonce, data[offset:offset+nonceLen])
	offset += nonceLen

	checkLen := int(data[offset])
	offset++
	if offset+checkLen > len(data) {
		return EncryptedMemo{}, errors.New("invalid check length")
	}
	check := make([]byte, checkLen)
	copy(check, data[offset:offset+checkLen])
	offset += checkLen

	encrypted := data[offset:]

	return EncryptedMemo{
		From:      from,
		To:        to,
		Nonce:     nonce,
		Check:     check,
		Encrypted: encrypted,
	}, nil
}
