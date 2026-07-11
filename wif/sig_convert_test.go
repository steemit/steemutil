package wif

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

// TestValidateCompactSignature_RoundTrip verifies that a real compact
// signature passes validation unchanged and still recovers the original public
// key. Signatures are produced with btcec's SignCompact (compressed), whose
// byte0 layout (27 + recovery + 4 = 31..34) is identical to the dsteem/
// libcrypto layout (recovery + 31).
func TestValidateCompactSignature_RoundTrip(t *testing.T) {
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

		out, err := ValidateCompactSignature(sig)
		if err != nil {
			t.Fatalf("ValidateCompactSignature failed: %v", err)
		}
		if !bytes.Equal(out, sig) {
			t.Fatalf("ValidateCompactSignature altered the signature")
		}

		// The validated signature must still recover the original public key.
		recovered, err := RecoverPublicKeyFromSignature(msg, out)
		if err != nil {
			t.Fatalf("RecoverPublicKeyFromSignature failed: %v", err)
		}
		if recovered.ToStr() != d.PublicKey {
			t.Fatalf("recovered key %q does not match original %q", recovered.ToStr(), d.PublicKey)
		}
	}
}

// TestValidateCompactSignature_RejectsBadLength verifies that non-65-byte
// inputs are rejected.
func TestValidateCompactSignature_RejectsBadLength(t *testing.T) {
	cases := [][]byte{
		nil,
		{},
		make([]byte, 64),
		make([]byte, 66),
	}
	for i, in := range cases {
		if _, err := ValidateCompactSignature(in); err == nil {
			t.Errorf("case %d: expected error for length %d", i, len(in))
		}
	}
}

// TestValidateCompactSignature_RejectsBadRecoveryByte verifies that a first
// byte outside the compressed range [31, 35) is rejected, while 31..34 are
// accepted.
func TestValidateCompactSignature_RejectsBadRecoveryByte(t *testing.T) {
	good := make([]byte, 65)
	copy(good[1:], make([]byte, 64))

	for _, b := range []byte{31, 32, 33, 34} {
		good[0] = b
		if _, err := ValidateCompactSignature(good); err != nil {
			t.Errorf("byte0=%d: expected accept, got error: %v", b, err)
		}
	}

	for _, b := range []byte{0, 27, 30, 35, 36, 100} {
		good[0] = b
		if _, err := ValidateCompactSignature(good); err == nil {
			t.Errorf("byte0=%d: expected reject", b)
		}
	}
}

// TestValidateCompactSignature_JSRealVector validates a signature produced by
// @steemit/libcrypto's signRecoverably over a fixed 32-byte digest. It proves
// validation works on JS-produced input, not just btcec-produced input.
// byte0 == 32 (odd recovery parity), within [31, 35).
func TestValidateCompactSignature_JSRealVector(t *testing.T) {
	sigHex := "201cba8696cdecef383702be9b4349e6418adf0973c018534e11b30012f5d893d32e083c38fc719c7470c64bba949b2d4c94bd74bafb8df80a6c4495a51a4ed054"
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		t.Fatalf("hex decode: %v", err)
	}
	if sig[0] != 32 {
		t.Fatalf("sanity: expected byte0=32, got %d", sig[0])
	}

	out, err := ValidateCompactSignature(sig)
	if err != nil {
		t.Fatalf("ValidateCompactSignature on JS vector failed: %v", err)
	}
	if !bytes.Equal(out, sig) {
		t.Fatalf("JS vector altered by validation")
	}
}
