package wif

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

// TestSigConvert_RoundTrip verifies that DsteemSigToBtcec / BtcecSigToDsteem
// are validating identity transforms: a real compact signature survives a
// round-trip unchanged. The signatures are produced with btcec's SignCompact
// (compressed), whose byte0 layout (27 + recovery + 4 = 31..34) is identical
// to the dsteem/libcrypto layout (recovery + 31).
func TestSigConvert_RoundTrip(t *testing.T) {
	for _, d := range data {
		priv := &PrivateKey{}
		if err := priv.FromWif(d.WIF); err != nil {
			t.Fatalf("FromWif failed: %v", err)
		}

		msg := make([]byte, 32)
		for i := range msg {
			msg[i] = byte(i)
		}

		sig, err := ecdsa.SignCompact(priv.Raw.PrivKey, msg, true)
		if err != nil {
			t.Fatalf("SignCompact failed: %v", err)
		}

		// Dsteem -> btcec -> dsteem must be a no-op on valid input.
		btcecSig, err := DsteemSigToBtcec(sig)
		if err != nil {
			t.Fatalf("DsteemSigToBtcec failed: %v", err)
		}
		if !bytes.Equal(btcecSig, sig) {
			t.Fatalf("DsteemSigToBtcec altered the signature")
		}

		back, err := BtcecSigToDsteem(btcecSig)
		if err != nil {
			t.Fatalf("BtcecSigToDsteem failed: %v", err)
		}
		if !bytes.Equal(back, sig) {
			t.Fatalf("round-trip altered the signature")
		}

		// The converted signature must still recover the original public key.
		recovered, err := RecoverPublicKeyFromSignature(msg, btcecSig)
		if err != nil {
			t.Fatalf("RecoverPublicKeyFromSignature failed: %v", err)
		}
		if recovered.ToStr() != d.PublicKey {
			t.Fatalf("recovered key %q does not match original %q", recovered.ToStr(), d.PublicKey)
		}
	}
}

// TestSigConvert_RejectsBadLength verifies that non-65-byte inputs are rejected.
func TestSigConvert_RejectsBadLength(t *testing.T) {
	cases := [][]byte{
		nil,
		{},
		make([]byte, 64),
		make([]byte, 66),
	}
	for i, in := range cases {
		if _, err := DsteemSigToBtcec(in); err == nil {
			t.Errorf("case %d: expected error for length %d", i, len(in))
		}
		if _, err := BtcecSigToDsteem(in); err == nil {
			t.Errorf("case %d: expected error for length %d", i, len(in))
		}
	}
}

// TestSigConvert_RejectsBadRecoveryByte verifies that a first byte outside the
// compressed range [31, 35) is rejected, while 31..34 are accepted.
func TestSigConvert_RejectsBadRecoveryByte(t *testing.T) {
	good := make([]byte, 65)
	copy(good[1:], make([]byte, 64))

	for _, b := range []byte{31, 32, 33, 34} {
		good[0] = b
		if _, err := DsteemSigToBtcec(good); err != nil {
			t.Errorf("byte0=%d: expected accept, got error: %v", b, err)
		}
		if _, err := BtcecSigToDsteem(good); err != nil {
			t.Errorf("byte0=%d: expected accept, got error: %v", b, err)
		}
	}

	for _, b := range []byte{0, 27, 30, 35, 36, 100} {
		good[0] = b
		if _, err := DsteemSigToBtcec(good); err == nil {
			t.Errorf("byte0=%d: expected reject", b)
		}
		if _, err := BtcecSigToDsteem(good); err == nil {
			t.Errorf("byte0=%d: expected reject", b)
		}
	}
}

// helper: a PublicKey view of the private key for comparisons.
func (p *PrivateKey) ToPubKey() *PublicKey {
	return &PublicKey{Raw: p.Raw.PrivKey.PubKey()}
}

// JSLibcryptoVector is a real signature produced by @steemit/libcrypto's
// signRecoverably over a fixed 32-byte digest. It proves the identity transform
// works on JS-produced input, not just btcec-produced input. byte0 == 32 (odd
// recovery parity), within [31, 35).
func TestSigConvert_JSRealVector(t *testing.T) {
	sigHex := "201cba8696cdecef383702be9b4349e6418adf0973c018534e11b30012f5d893d32e083c38fc719c7470c64bba949b2d4c94bd74bafb8df80a6c4495a51a4ed054"
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		t.Fatalf("hex decode: %v", err)
	}
	if sig[0] != 32 {
		t.Fatalf("sanity: expected byte0=32, got %d", sig[0])
	}

	out, err := DsteemSigToBtcec(sig)
	if err != nil {
		t.Fatalf("DsteemSigToBtcec on JS vector failed: %v", err)
	}
	if !bytes.Equal(out, sig) {
		t.Fatalf("JS vector altered by identity transform")
	}
}
