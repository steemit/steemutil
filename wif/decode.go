package wif

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"
)

// Decode can be used to turn WIF into a raw private key (32 bytes).
func DecodeWif(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode WIF")
	}

	return w.PrivKey.Serialize(), nil
}

// GetPublicKey returns the public key associated with the given WIF
// in the 33-byte compressed format.
func GetPublicKey(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode WIF")
	}

	return w.PrivKey.PubKey().SerializeCompressed(), nil
}
