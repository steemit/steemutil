package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"io"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/pkg/errors"
	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/wif"
)

// EncryptedMemo represents an encrypted memo, serialized to the canonical
// Steem wire format (compatible with steem-js / steemd).
//
// Layout (all little-endian):
//
//	from     33 bytes  compressed secp256k1 public key of the sender
//	to       33 bytes  compressed secp256k1 public key of the recipient
//	nonce     8 bytes  uint64 nonce (little-endian)
//	check     4 bytes  uint32 checksum (little-endian)
//	n         varint   byte length of the ciphertext
//	cipher    n bytes  AES-256-CBC ciphertext
type EncryptedMemo struct {
	From      *wif.PublicKey
	To        *wif.PublicKey
	Nonce     uint64
	Check     uint32
	Encrypted []byte
}

// nonceLen is the fixed length of the on-wire nonce (uint64, little-endian).
const nonceLen = 8

// Encode encrypts a memo if it starts with '#', otherwise returns it as-is.
// privateKey can be a WIF string or a *wif.PrivateKey.
// publicKey can be a public key string (STM...) or a *wif.PublicKey.
//
// The encrypted output is byte-for-byte compatible with steem-js
// memo.encode (next branch): ECDH shared secret, SHA-512 key derivation,
// AES-256-CBC, and the canonical EncryptedMemo serialization.
func Encode(privateKey interface{}, publicKey interface{}, memo string) (string, error) {
	return EncodeWithNonce(privateKey, publicKey, memo, UniqueNonce())
}

// EncodeWithNonce is like Encode but with a caller-supplied nonce. It exists
// primarily for deterministic tests; production callers should use Encode.
func EncodeWithNonce(privateKey interface{}, publicKey interface{}, memo string, nonce uint64) (string, error) {
	if memo == "" {
		return "", errors.New("memo is required")
	}

	// Plain memos (no '#' prefix) pass through unmodified.
	if memo[0] != '#' {
		return memo, nil
	}
	plain := memo[1:]

	privKey, err := toPrivateKey(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert private key")
	}
	pubKey, err := toPublicKey(publicKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert public key")
	}

	// Sender's own public key.
	senderPubKey := &wif.PublicKey{}
	if err := senderPubKey.FromStr(privKey.ToPubKeyStr()); err != nil {
		return "", errors.Wrap(err, "failed to derive sender public key")
	}

	encrypted, checksum, err := encryptMemo(privKey, pubKey, []byte(plain), nonce)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt memo")
	}

	encMemo := EncryptedMemo{
		From:      senderPubKey,
		To:        pubKey,
		Nonce:     nonce,
		Check:     checksum,
		Encrypted: encrypted,
	}

	memoBytes, err := serializeEncryptedMemo(encMemo)
	if err != nil {
		return "", errors.Wrap(err, "failed to serialize encrypted memo")
	}

	return "#" + base58.Encode(memoBytes), nil
}

// Decode decrypts a memo if it starts with '#', otherwise returns it as-is.
// privateKey can be a WIF string or a *wif.PrivateKey.
//
// The decryption is symmetric with Encode and interoperates with memos
// produced by steem-js memo.encode: given a memo encrypted by A for B, B (and
// only B) can recover the plaintext.
func Decode(privateKey interface{}, memo string) (string, error) {
	if memo == "" {
		return "", errors.New("memo is required")
	}
	if memo[0] != '#' {
		return memo, nil
	}

	privKey, err := toPrivateKey(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert private key")
	}

	memoBytes := base58.Decode(memo[1:])

	encMemo, err := deserializeEncryptedMemo(memoBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to deserialize encrypted memo")
	}

	// Identify the other party's public key. We are either the sender (our
	// pubkey == From) or the recipient (our pubkey == To).
	senderPubKey := &wif.PublicKey{}
	if err := senderPubKey.FromStr(privKey.ToPubKeyStr()); err != nil {
		return "", errors.Wrap(err, "failed to derive our public key")
	}

	var otherPubKey *wif.PublicKey
	switch {
	case senderPubKey.ToStr() == encMemo.From.ToStr():
		otherPubKey = encMemo.To
	case senderPubKey.ToStr() == encMemo.To.ToStr():
		otherPubKey = encMemo.From
	default:
		return "", errors.New("memo was not encrypted for this key")
	}

	plain, err := decryptMemo(privKey, otherPubKey, encMemo.Nonce, encMemo.Encrypted, encMemo.Check)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt memo")
	}

	return "#" + string(plain), nil
}

// ----------------------------------------------------------------------------
// key conversion helpers
// ----------------------------------------------------------------------------

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

// ----------------------------------------------------------------------------
// cryptography (Steem memo spec, matches steem-js)
// ----------------------------------------------------------------------------

// sharedSecret computes the ECDH shared secret S = SHA-512(x), where x is the
// 32-byte X coordinate of privKey * pubKey on secp256k1. This is symmetric:
// SHA-512(privA * pubB) == SHA-512(privB * pubA).
//
// wif.PrivateKey.Raw.PrivKey is *btcec.PrivateKey, which is a type alias for
// *secp256k1.PrivateKey (btcec/v2 re-exports decred's secp256k1/v4), so it can
// be passed to GenerateSharedSecret with no conversion.
func sharedSecret(privKey *wif.PrivateKey, pubKey *wif.PublicKey) []byte {
	x := secp256k1.GenerateSharedSecret(privKey.Raw.PrivKey, pubKey.Raw)
	s := sha512.Sum512(x)
	return s[:] // 64 bytes
}

// deriveEncryptionKey derives the 64-byte encryption key material:
//
//	ek = SHA-512( uint64_le(nonce) || S )
//
// The first 32 bytes are the AES-256 key; bytes [32:48] are the AES-CBC IV;
// the checksum is derived separately (see encryptMemo).
func deriveEncryptionKey(privKey *wif.PrivateKey, pubKey *wif.PublicKey, nonce uint64) []byte {
	s := sharedSecret(privKey, pubKey)
	var nonceBuf [nonceLen]byte
	binary.LittleEndian.PutUint64(nonceBuf[:], nonce)
	h := sha512.New()
	h.Write(nonceBuf[:])
	h.Write(s)
	return h.Sum(nil) // 64 bytes
}

// checksumOf returns the 4-byte little-endian checksum = SHA-256(ek)[0:4]
// reinterpreted as a uint32 (matching steem-js Aes.encrypt).
func checksumOf(ek []byte) uint32 {
	sum := sha256.Sum256(ek)
	return binary.LittleEndian.Uint32(sum[0:4])
}

// encryptMemo encrypts message for pubKey under a fresh key derived from the
// ECDH shared secret and nonce. Returns (ciphertext, checksum).
func encryptMemo(privKey *wif.PrivateKey, pubKey *wif.PublicKey, message []byte, nonce uint64) ([]byte, uint32, error) {
	ek := deriveEncryptionKey(privKey, pubKey, nonce)
	key := ek[0:32]
	iv := ek[32:48]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, 0, err
	}

	// Steem prefixes the plaintext with a varint length (bytebuffer
	// writeVString) before encryption.
	prefixed := prefixVarString(message)
	padded := pkcs7Pad(prefixed, aes.BlockSize)

	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	return ciphertext, checksumOf(ek), nil
}

// decryptMemo is the inverse of encryptMemo.
func decryptMemo(privKey *wif.PrivateKey, pubKey *wif.PublicKey, nonce uint64, ciphertext []byte, checksum uint32) ([]byte, error) {
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("invalid ciphertext length")
	}

	ek := deriveEncryptionKey(privKey, pubKey, nonce)
	if checksumOf(ek) != checksum {
		return nil, errors.New("checksum mismatch")
	}

	block, err := aes.NewCipher(ek[0:32])
	if err != nil {
		return nil, err
	}

	padded := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, ek[32:48]).CryptBlocks(padded, ciphertext)

	plain, err := pkcs7Unpad(padded)
	if err != nil {
		return nil, err
	}

	// Strip the varint length prefix that writeVString added at encrypt time.
	return stripVarString(plain)
}

// prefixVarString prepends an unsigned LEB128 varint encoding of len(b),
// matching bytebuffer's writeVString (which uses unsigned writeVarint32).
func prefixVarString(b []byte) []byte {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], uint64(len(b)))
	out := make([]byte, n+len(b))
	copy(out, buf[:n])
	copy(out[n:], b)
	return out
}

// stripVarString removes the leading unsigned varint length and returns the
// bytes it described. It mirrors prefixVarString and tolerates only a
// well-formed length that matches the remaining bytes.
func stripVarString(b []byte) ([]byte, error) {
	length, n := binary.Uvarint(b)
	if n <= 0 {
		return nil, errors.New("invalid varint length prefix")
	}
	if uint64(len(b)-n) != length {
		return nil, errors.Errorf("varint length %d does not match payload %d", length, len(b)-n)
	}
	return b[n:], nil
}

// pkcs7Pad applies PKCS#7 padding so len(data) is a multiple of blockSize.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpad validates and removes PKCS#7 padding.
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > len(data) {
		return nil, errors.New("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:len(data)-padding], nil
}

// ----------------------------------------------------------------------------
// nonce
// ----------------------------------------------------------------------------

// UniqueNonce returns a 64-bit nonce composed of the high 32 bits of the
// current time in milliseconds and 32 random bits, matching steem-js
// Aes.uniqueNonce() (next branch, PR #531).
func UniqueNonce() uint64 {
	now := uint64(time.Now().UnixMilli()) // low 32 bits kept via shift below
	var randBuf [4]byte
	if _, err := io.ReadFull(rand.Reader, randBuf[:]); err != nil {
		// rand.Reader should never fail; fall back to time entropy if it does.
		now ^= uint64(time.Now().UnixNano())
		return now
	}
	random := uint64(randBuf[0])<<24 | uint64(randBuf[1])<<16 | uint64(randBuf[2])<<8 | uint64(randBuf[3])
	return (now << 32) | random
}

// ----------------------------------------------------------------------------
// wire format (de)serialization
// ----------------------------------------------------------------------------

// pubKeyLen is the length of a compressed secp256k1 public key.
const pubKeyLen = 33

func serializeEncryptedMemo(m EncryptedMemo) ([]byte, error) {
	if m.From == nil || m.To == nil {
		return nil, errors.New("encrypted memo missing from/to public keys")
	}
	var buf bytes.Buffer
	enc := encoder.NewEncoder(&buf)

	// from / to: raw 33-byte compressed public keys
	if err := enc.WriteBytes(m.From.ToByte()); err != nil {
		return nil, err
	}
	if err := enc.WriteBytes(m.To.ToByte()); err != nil {
		return nil, err
	}
	// nonce: uint64 little-endian
	if err := enc.EncodeNumber(m.Nonce); err != nil {
		return nil, err
	}
	// check: uint32 little-endian
	if err := enc.EncodeNumber(m.Check); err != nil {
		return nil, err
	}
	// ciphertext length: unsigned varint (varint32 in steem-js)
	if err := enc.EncodeUVarint(uint64(len(m.Encrypted))); err != nil {
		return nil, err
	}
	if err := enc.WriteBytes(m.Encrypted); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deserializeEncryptedMemo(data []byte) (EncryptedMemo, error) {
	var m EncryptedMemo

	// Minimum: two public keys (33+33) + nonce(8) + check(4) + at least one
	// varint byte. Reject short / malformed input early.
	const minLen = pubKeyLen*2 + nonceLen + 4 + 1
	if len(data) < minLen {
		return m, errors.Errorf("encrypted memo too short: %d bytes", len(data))
	}

	off := 0
	read := func(n int) []byte {
		b := data[off : off+n]
		off += n
		return b
	}

	fromBytes := read(pubKeyLen)
	toBytes := read(pubKeyLen)

	m.From = &wif.PublicKey{}
	if err := m.From.FromByte(fromBytes); err != nil {
		return m, errors.Wrap(err, "failed to parse 'from' public key")
	}
	m.To = &wif.PublicKey{}
	if err := m.To.FromByte(toBytes); err != nil {
		return m, errors.Wrap(err, "failed to parse 'to' public key")
	}

	m.Nonce = binary.LittleEndian.Uint64(read(nonceLen))
	m.Check = binary.LittleEndian.Uint32(read(4))

	length, n := binary.Uvarint(data[off:])
	if n <= 0 {
		return m, errors.New("invalid ciphertext length varint")
	}
	off += n

	if uint64(len(data)-off) != length {
		return m, errors.Errorf("ciphertext length %d does not match remaining %d bytes", length, len(data)-off)
	}
	m.Encrypted = make([]byte, len(data)-off)
	copy(m.Encrypted, data[off:])
	return m, nil
}
