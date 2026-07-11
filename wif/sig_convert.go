package wif

import (
	"github.com/pkg/errors"
)

// compactSigSize is the length of a secp256k1 compact signature: a 1-byte
// recovery code followed by 32-byte R and 32-byte S values.
const compactSigSize = 65

// dsteemRecoveryBase is the base value dsteem / @steemit/libcrypto write into
// the first byte of a compact signature: recovery (0 or 1) + 31, so the byte
// is always 31 or 32 (compressed keys only).
//
// This is byte-for-byte identical to the btcec / decred secp256k1 layout,
// which encodes the first byte as compactSigMagicOffset(27) + recoveryCode
// + 4(compression) for compressed keys, i.e. 27 + 0..3 + 4 = 31..34.
//
// References:
//   - libcrypto-js/src/codecSteemit.js signRecoverably (recoveryParam = 31, +1 if odd)
//   - dsteem/src/crypto.ts Signature.toBuffer (recovery + 31)
//   - decred secp256k1/v4 ecdsa/signature.go (compactSigMagicOffset = 27)
//   - steem fc/src/crypto/elliptic_impl_pub.cpp (accepts nV in [27, 35))
//
// Because both layouts agree, the conversion functions below are identity
// transforms whose only job is to validate the byte is in the compressed
// range [31, 35). No offset arithmetic is applied.
const dsteemRecoveryBase = 31

// validateCompactSig checks that sig is a 65-byte compact signature with a
// first byte in the compressed-key range [31, 35). This range matches all four
// implementations (libcrypto, dsteem, btcec, steemd) for compressed keys.
func validateCompactSig(sig []byte) error {
	if len(sig) != compactSigSize {
		return errors.Errorf("invalid compact signature length: expected %d, got %d", compactSigSize, len(sig))
	}
	if sig[0] < dsteemRecoveryBase || sig[0] >= dsteemRecoveryBase+4 {
		return errors.Errorf("invalid signature recovery byte: expected %d..%d, got %d",
			dsteemRecoveryBase, dsteemRecoveryBase+3, sig[0])
	}
	return nil
}

// DsteemSigToBtcec validates a 65-byte dsteem/libcrypto-format compact
// signature and returns it unmodified.
//
// The dsteem/libcrypto layout (byte0 = recovery+31) and the btcec/decred
// layout (byte0 = 27+recoveryCode+4 for compressed keys) are byte-for-byte
// identical for compressed keys, so no conversion is needed. This function
// exists to guard against malformed input before handing the signature to
// RecoverPublicKeyFromSignature, and to preserve a stable API at the boundary
// between conveyor's verification path and the wif package.
//
// The input must be 65 bytes with byte0 in [31, 35). It is returned as-is on
// success.
func DsteemSigToBtcec(dsteemSig []byte) ([]byte, error) {
	if err := validateCompactSig(dsteemSig); err != nil {
		return nil, errors.Wrap(err, "failed to validate dsteem signature")
	}
	out := make([]byte, compactSigSize)
	copy(out, dsteemSig)
	return out, nil
}

// BtcecSigToDsteem is the inverse of DsteemSigToBtcec. Since the two layouts
// are identical for compressed keys, this is likewise a validating identity
// transform. The input must be 65 bytes with byte0 in [31, 35).
func BtcecSigToDsteem(btcecSig []byte) ([]byte, error) {
	if err := validateCompactSig(btcecSig); err != nil {
		return nil, errors.Wrap(err, "failed to validate btcec signature")
	}
	out := make([]byte, compactSigSize)
	copy(out, btcecSig)
	return out, nil
}
