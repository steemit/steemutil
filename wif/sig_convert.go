package wif

import (
	"github.com/pkg/errors"
)

// compactSigSize is the length of a secp256k1 compact signature: a 1-byte
// recovery code followed by 32-byte R and 32-byte S values.
const compactSigSize = 65

// compactSigRecoveryBase is the smallest valid first byte for a compressed-key
// compact signature across all implementations in the Steem stack.
//
// All four implementations agree on the same byte layout for compressed keys:
//   - @steemit/libcrypto signRecoverably: byte0 = recovery + 31  -> {31, 32}
//   - dsteem Signature.toBuffer:        byte0 = recovery + 31  -> {31, 32}
//   - btcec/decred SignCompact:         byte0 = 27 + rc + 4    -> {31..34}
//   - steem fc recover:                 accepts nV in [27, 35)
//
// Because the layouts are byte-for-byte identical for compressed keys, no
// conversion is needed; this constant only defines the valid [31, 35) range
// used by the validator below.
const compactSigRecoveryBase = 31

// ValidateCompactSignature checks that sig is a 65-byte compact signature
// whose first byte is in the compressed-key range [31, 35) and returns a copy
// of it. This is a format guard, not a conversion: the dsteem/libcrypto and
// btcec/decred compact layouts are byte-for-byte identical for compressed keys,
// so the signature is returned unchanged on success.
//
// VerifySignedRpc runs signatures through this before handing them to
// RecoverPublicKeyFromSignature, so that malformed input is rejected with a
// clear error instead of producing a confusing recovery failure.
func ValidateCompactSignature(sig []byte) ([]byte, error) {
	if len(sig) != compactSigSize {
		return nil, errors.Errorf("invalid compact signature length: expected %d, got %d", compactSigSize, len(sig))
	}
	if sig[0] < compactSigRecoveryBase || sig[0] >= compactSigRecoveryBase+4 {
		return nil, errors.Errorf("invalid signature recovery byte: expected %d..%d, got %d",
			compactSigRecoveryBase, compactSigRecoveryBase+3, sig[0])
	}
	out := make([]byte, compactSigSize)
	copy(out, sig)
	return out, nil
}
