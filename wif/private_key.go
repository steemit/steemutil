package wif

import (
	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
)

type PrivateKey struct {
	Raw *btcutil.WIF
}

// Get a raw private key (32 bytes) from WIF.
func (p *PrivateKey) FromWif(wif string) (err error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return errors.Wrap(err, "failed to decode WIF")
	}
	p.Raw = w
	return
}

func (p *PrivateKey) FromByte(raw []byte) (err error) {
	privKey, _ := btcec.PrivKeyFromBytes(raw)
	// Steem uses 0x80 as WIF version byte (same as Bitcoin mainnet) and uncompressed format
	steemParams := &chaincfg.Params{
		PrivateKeyID: 0x80,
	}
	tmp, err := btcutil.NewWIF(privKey, steemParams, false)
	if err != nil {
		return errors.Wrap(err, "failed to create new WIF struct")
	}
	p.Raw = tmp
	return
}

func (p *PrivateKey) ToByte() []byte {
	return p.Raw.PrivKey.Serialize()
}

func (p *PrivateKey) ToWif() string {
	return p.Raw.String()
}

func (p *PrivateKey) ToPubKeyStr() string {
	_, pubKey := btcec.PrivKeyFromBytes(p.ToByte())
	publicKey := &PublicKey{
		Raw: pubKey,
	}
	return publicKey.ToStr()
}
