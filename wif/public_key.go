package wif

import (
	"bytes"
	"hash"

	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/steemit/steemutil/consts"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/base58"

	//nolint:staticcheck // RIPEMD-160 is required by Steem protocol for public key checksum
	"golang.org/x/crypto/ripemd160"

	"github.com/pkg/errors"
)

type PublicKey struct {
	Raw *secp256k1.PublicKey
}

// Get a raw public key from []byte.
func (p *PublicKey) FromByte(raw []byte) (err error) {
	pubKey, err := btcec.ParsePubKey(raw)
	if err != nil {
		return errors.Wrap(err, "failed to decode Public Key")
	}
	p.Raw = pubKey
	return nil
}

// Get a raw public key from Public String (STM...).
func (p *PublicKey) FromStr(pubKey string) (err error) {
	// check prefix
	prefixLen := len(consts.ADDRESS_PREFIX)
	prefix := pubKey[0:prefixLen]
	if prefix != consts.ADDRESS_PREFIX {
		return errors.New("public key has an error prefix")
	}
	// get pub key without prefix
	pubKeyWithoutPrefix := pubKey[prefixLen:]
	pubKeyByte := base58.Decode(pubKeyWithoutPrefix)

	// check checksum
	pubKeyOri := pubKeyByte[0 : len(pubKeyByte)-4]
	checkSum := pubKeyByte[len(pubKeyByte)-4:]
	calcCheckSum := calcHash(pubKeyOri, ripemd160.New())
	if !bytes.Equal(checkSum, calcCheckSum[0:4]) {
		return errors.New("public key checksum failed")
	}

	// save
	err = p.FromByte(pubKeyOri)
	return err
}

func (p *PublicKey) FromWif(wif string) (err error) {
	tmp := &PrivateKey{}
	err = tmp.FromWif(wif)
	if err != nil {
		return
	}
	err = p.FromStr(tmp.ToPubKeyStr())
	return
}

func (p *PublicKey) ToStr() string {
	checkSum := calcHash(p.Raw.SerializeCompressed(), ripemd160.New())
	pubByte := append(p.Raw.SerializeCompressed(), checkSum[0:4]...)
	pubStr := base58.Encode(pubByte)
	return consts.ADDRESS_PREFIX + pubStr
}

func (p *PublicKey) ToByte() []byte {
	return p.Raw.SerializeCompressed()
}

func calcHash(buf []byte, hasher hash.Hash) []byte {
	_, _ = hasher.Write(buf)
	return hasher.Sum(nil)
}
