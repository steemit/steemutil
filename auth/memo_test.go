package auth

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/steemit/steemutil/wif"
)

// ----------------------------------------------------------------------------
// steem-js interop vectors
//
// Generated from steem-js `next` @ d3945f1 (src/memo + src/auth/ecc/src/aes),
// the post-#524/#525/#530/#531 security-audited state. These pin byte-for-byte
// compatibility with the canonical Steem memo format.
//
// Recipient key is derived from a fixed seed — TEST ONLY, do not use on mainnet.
// ----------------------------------------------------------------------------

const (
	// senderWif is the historical steemutil test WIF; its public key is senderPub.
	senderWif = "5JRaypasxMx1L97ZUX7YuC5Psb5EAbF821kkAGtBj7xCJFQcbLg"
	senderPub = "STM6aGPtxMUGnTPfKLSxdwCHbximSJxzrRjeQmwRW9BRCdrFotKLs"

	// recipientWif / recipientPub come from PrivateKey.fromBuffer(seed) where
	// seed = 0123456789abcdef... (32 bytes).
	recipientWif = "5HpneLQNKrcznVCQpzodYwAmZ4AoHeyjuRf9iAHAa498rP5kuWb"
	recipientPub = "STM7NBdzdQjZi7M6xo7XqcwcW4YCBLZeqjBvpj6obDpiXhbTSnhjq"

	// interopVectorNonce is the fixed 64-bit nonce used to generate the vector.
	interopVectorNonce uint64 = 14653678900875387904

	// interopVectorMemo is the plaintext WITHOUT the leading '#'.
	interopVectorMemo = "hello steem"

	// interopVectorEnc is the exact output of steem-js memo.encode(senderWif,
	// recipientPub, "#hello steem", "14653678900875387904").
	interopVectorEnc = "#Dyn1R3F3BGrJgACeq3pw3foykzH18UqUTdx1hZ7PxyXGNogQ37vr2WhzXpAY2VgM7iRqShkhAcSjnpVW2Gbpc5pXnt8GjnL3mU6jyXngtz8ekbpmdQGuBtXmKAgoCEdj1"

	// interopVectorWireHex is the base58-decoded wire bytes of interopVectorEnc.
	interopVectorWireHex = "02de0381af8fa527ceee901e2463aa15132a69b96368782a1bced60065fc271b1a034646ae5047316b4230d0086c8acec687f00b1cd9d1dc634f6cb358ac0a9a8fff0064506352535ccb8e061c0610a08212a29f599d0ce16ace3086da554c"
)

// TestEncodeDecodePlainText: plain memos pass through untouched both ways.
func TestEncodeDecodePlainText(t *testing.T) {
	memo := "plain text memo"

	encoded, err := Encode(nil, nil, memo)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if encoded != memo {
		t.Errorf("Encode should return plain text as-is, got: %s", encoded)
	}

	decoded, err := Decode(nil, memo)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != memo {
		t.Errorf("Decode should return plain text as-is, got: %s", decoded)
	}
}

// TestEncodeDecodeEmptyMemo: empty input is rejected.
func TestEncodeDecodeEmptyMemo(t *testing.T) {
	if _, err := Encode(nil, nil, ""); err == nil {
		t.Error("Encode should fail for empty memo")
	}
	if _, err := Decode(nil, ""); err == nil {
		t.Error("Decode should fail for empty memo")
	}
}

// TestEncodeByteEqualityWithSteemJS pins the core interop guarantee: Go
// EncodeWithNonce must produce the exact same ciphertext as steem-js
// memo.encode for identical inputs. A single byte drift here means we have
// broken wire-format compatibility.
func TestEncodeByteEqualityWithSteemJS(t *testing.T) {
	got, err := EncodeWithNonce(senderWif, recipientPub, "#"+interopVectorMemo, interopVectorNonce)
	if err != nil {
		t.Fatalf("EncodeWithNonce failed: %v", err)
	}
	if got != interopVectorEnc {
		t.Fatalf("byte-equality with steem-js broken:\n got:  %s\n want: %s", got, interopVectorEnc)
	}
}

// TestDecodeSteemJSMemo verifies our decoder can decrypt a memo produced by
// steem-js (the wire is steem-js's, the key is the recipient's).
func TestDecodeSteemJSMemo(t *testing.T) {
	decoded, err := Decode(recipientWif, interopVectorEnc)
	if err != nil {
		t.Fatalf("Decode of steem-js memo failed: %v", err)
	}
	if decoded != "#"+interopVectorMemo {
		t.Fatalf("unexpected plaintext: got %q, want %q", decoded, "#"+interopVectorMemo)
	}
}

// TestCrossKeyEncryption is the regression test for the original [高] finding.
//
// The bug: deriveSharedSecret did SHA-256(privBytes || pubBytes), which is NOT
// symmetric — sender's key was SHA-256(privS||pubR) while the recipient's was
// SHA-256(privR||pubS), so the recipient could never decrypt. The old tests
// passed only because they used the SAME WIF for encrypt and decrypt
// (self-to-self), hiding the asymmetry.
//
// This test uses two DIFFERENT keys: sender encrypts, recipient decrypts. It
// must succeed now (real ECDH is symmetric) and would fail against the old code.
func TestCrossKeyEncryption(t *testing.T) {
	const plaintext = "#secret message between two parties"

	// Sender encrypts to the recipient's public key.
	enc, err := EncodeWithNonce(senderWif, recipientPub, plaintext, interopVectorNonce)
	if err != nil {
		t.Fatalf("Encode (sender) failed: %v", err)
	}
	if enc == plaintext {
		t.Fatal("memo was not encrypted")
	}

	// Recipient decrypts with its own private key.
	dec, err := Decode(recipientWif, enc)
	if err != nil {
		t.Fatalf("Decode (recipient) failed: %v — this is exactly the original bug", err)
	}
	if dec != plaintext {
		t.Fatalf("cross-key decryption mismatch: got %q, want %q", dec, plaintext)
	}

	// Symmetry the other way: recipient encrypts to sender, sender decrypts.
	enc2, err := EncodeWithNonce(recipientWif, senderPub, plaintext, interopVectorNonce)
	if err != nil {
		t.Fatalf("Encode (recipient) failed: %v", err)
	}
	dec2, err := Decode(senderWif, enc2)
	if err != nil {
		t.Fatalf("Decode (sender) failed: %v", err)
	}
	if dec2 != plaintext {
		t.Fatalf("reverse cross-key decryption mismatch: got %q, want %q", dec2, plaintext)
	}
}

// TestEncodeAcceptsKeyObjects verifies the interface{} entry points accept
// *wif.PrivateKey / *wif.PublicKey as well as strings.
func TestEncodeAcceptsKeyObjects(t *testing.T) {
	priv, err := toPrivateKey(senderWif)
	if err != nil {
		t.Fatalf("toPrivateKey: %v", err)
	}
	pub, err := toPublicKey(recipientPub)
	if err != nil {
		t.Fatalf("toPublicKey: %v", err)
	}

	enc, err := EncodeWithNonce(priv, pub, "#obj test", interopVectorNonce)
	if err != nil {
		t.Fatalf("EncodeWithNonce objects: %v", err)
	}
	// Same inputs as strings must yield the same ciphertext.
	encStr, err := EncodeWithNonce(senderWif, recipientPub, "#obj test", interopVectorNonce)
	if err != nil {
		t.Fatalf("EncodeWithNonce strings: %v", err)
	}
	if enc != encStr {
		t.Errorf("object vs string inputs diverged:\n obj: %s\n str: %s", enc, encStr)
	}
}

// TestDecodeRejectsWrongKey: a third party (neither sender nor recipient)
// cannot decrypt the memo — Decode must report it was not encrypted for them.
func TestDecodeRejectsWrongKey(t *testing.T) {
	// A completely unrelated key.
	const strangerWif = "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	if _, err := Decode(strangerWif, interopVectorEnc); err == nil {
		t.Fatal("Decode with an unrelated key should fail, but succeeded")
	}
}

// TestDecodeRejectsTamperedMemo: flipping any byte of the ciphertext must
// produce a checksum/padding failure, not a silently wrong plaintext.
func TestDecodeRejectsTamperedMemo(t *testing.T) {
	enc, err := EncodeWithNonce(senderWif, recipientPub, "#tamper me", interopVectorNonce)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	// Flip the last data character of the base58 payload (not the '#').
	tampered := enc[:len(enc)-1] + string(flipBase58Char(enc[len(enc)-1]))
	if tampered == enc {
		t.Skip("could not produce a differing tampered string")
	}
	if _, err := Decode(recipientWif, tampered); err == nil {
		t.Fatal("Decode of a tampered memo should fail")
	}
}

// TestRoundTripUniqueNonce verifies the default Encode path (random nonce)
// still round-trips through Decode.
func TestRoundTripUniqueNonce(t *testing.T) {
	const plaintext = "#random nonce round trip"

	enc, err := Encode(senderWif, recipientPub, plaintext)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	// Two calls with random nonces must (overwhelmingly likely) differ.
	enc2, err := Encode(senderWif, recipientPub, plaintext)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if enc == enc2 {
		t.Error("expected distinct ciphertexts from random nonces, got identical")
	}

	dec, err := Decode(recipientWif, enc)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if dec != plaintext {
		t.Fatalf("round-trip mismatch: got %q want %q", dec, plaintext)
	}
}

// TestUniqueNonceMatchesSteemJSLayout sanity-checks that UniqueNonce packs the
// time into the high 32 bits and randomness into the low 32 bits, matching
// steem-js Aes.uniqueNonce(). (We assert structure, not an exact value.)
func TestUniqueNonceMatchesSteemJSLayout(t *testing.T) {
	hi := interopVectorNonce >> 32
	_ = hi // vector's high bits are time-ish ms (non-zero), low bits random.

	n := UniqueNonce()
	if n>>32 == 0 {
		t.Error("UniqueNonce high 32 bits are zero")
	}
	if n&0xffffffff == 0 {
		t.Error("UniqueNonce low 32 bits are zero (no randomness)")
	}
}

// TestWireFormatLayout decodes the steem-js vector by hand and asserts the
// exact byte layout the Go serializer must reproduce:
//
//	from(33) || to(33) || nonce(uint64 LE) || check(uint32 LE) || varint(len) || ct
func TestWireFormatLayout(t *testing.T) {
	wire, err := hex.DecodeString(interopVectorWireHex)
	if err != nil {
		t.Fatalf("bad wire hex: %v", err)
	}

	const pubLen = 33
	fromBytes := wire[0:pubLen]
	toBytes := wire[pubLen : pubLen*2]
	off := pubLen * 2

	// nonce (uint64 LE)
	var nonceBytes [8]byte
	copy(nonceBytes[:], wire[off:off+8])
	off += 8

	// check (uint32 LE)
	var checkBytes [4]byte
	copy(checkBytes[:], wire[off:off+4])
	off += 4

	// varint length of ciphertext
	length, n := varintDecode(wire[off:])
	off += n

	ciphertext := wire[off:]

	// Cross-check each field against independently-known values.
	fromPub := &wif.PublicKey{}
	if err := fromPub.FromByte(fromBytes); err != nil {
		t.Fatalf("parse from: %v", err)
	}
	if fromPub.ToStr() != senderPub {
		t.Errorf("from pubkey: got %s, want %s", fromPub.ToStr(), senderPub)
	}
	toPub := &wif.PublicKey{}
	if err := toPub.FromByte(toBytes); err != nil {
		t.Fatalf("parse to: %v", err)
	}
	if toPub.ToStr() != recipientPub {
		t.Errorf("to pubkey: got %s, want %s", toPub.ToStr(), recipientPub)
	}
	if fromBytes[0] != 0x02 && fromBytes[0] != 0x03 {
		t.Errorf("from key prefix byte unexpected: 0x%02x", fromBytes[0])
	}

	// nonce little-endian must equal interopVectorNonce.
	var nonceLE [8]byte
	for i := 0; i < 8; i++ {
		nonceLE[i] = byte(interopVectorNonce >> (8 * i))
	}
	if nonceBytes != nonceLE {
		t.Errorf("nonce LE mismatch: got %x, want %x", nonceBytes, nonceLE)
	}

	// Re-serialize the parsed memo and require byte-for-byte equality, which
	// simultaneously validates the varint length and the check field.
	memo := EncryptedMemo{From: fromPub, To: toPub}
	memo.Nonce = interopVectorNonce
	memo.Check = leUint32(checkBytes)
	memo.Encrypted = ciphertext
	reserialized, err := serializeEncryptedMemo(memo)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if hex.EncodeToString(reserialized) != interopVectorWireHex {
		t.Errorf("reserialized wire differs from steem-js:\n got  %x\n want %s", reserialized, interopVectorWireHex)
	}
	if int(length) != len(ciphertext) {
		t.Errorf("varint length %d != ciphertext len %d", length, len(ciphertext))
	}
}

// flipBase58Char returns a different base58 character (for tamper tests).
func flipBase58Char(c byte) byte {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	idx := strings.IndexByte(alphabet, c)
	if idx < 0 {
		return 'z'
	}
	return alphabet[(idx+1)%len(alphabet)]
}

// varintDecode reads an unsigned LEB128 varint, returning value and byte count.
func varintDecode(b []byte) (uint64, int) {
	var x uint64
	var s uint
	for i, c := range b {
		if i >= 10 {
			return 0, 0
		}
		if c < 0x80 {
			if i == 9 && c > 1 {
				return 0, 0
			}
			return x | uint64(c)<<s, i + 1
		}
		x |= uint64(c&0x7f) << s
		s += 7
	}
	return 0, 0
}

// leUint32 decodes a 4-byte little-endian uint32.
func leUint32(b [4]byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
