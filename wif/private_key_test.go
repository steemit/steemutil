package wif

import (
	"encoding/hex"
	"testing"
)

func TestFromWif(t *testing.T) {
	for _, d := range data {
		p := &PrivateKey{}
		err := p.FromWif(d.WIF)
		if err != nil {
			t.Error(err)
		}

		expected := d.PrivateKeyHex
		got := hex.EncodeToString(p.Raw.PrivKey.Serialize())

		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}

func TestToByte(t *testing.T) {
	for _, d := range data {
		p := &PrivateKey{}
		err := p.FromWif(d.WIF)
		if err != nil {
			t.Error(err)
		}

		expected := d.PrivateKeyHex
		got := hex.EncodeToString(p.ToByte())

		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}

func TestFromByte(t *testing.T) {
	for _, d := range data {
		p := &PrivateKey{}
		raw, err := hex.DecodeString(d.PrivateKeyHex)
		if err != nil {
			t.Error(err)
		}
		err = p.FromByte(raw)
		if err != nil {
			t.Error(err)
		}

		expected := d.WIF
		got := p.ToWif()
		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}

func TestToWif(t *testing.T) {
	for _, d := range data {
		p := &PrivateKey{}
		err := p.FromWif(d.WIF)
		if err != nil {
			t.Error(err)
		}

		expected := d.WIF
		got := p.ToWif()

		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}

func TestToPubKeyStr(t *testing.T) {
	for _, d := range data {
		p := &PrivateKey{}
		err := p.FromWif(d.WIF)
		if err != nil {
			t.Error(err)
		}

		expected := d.PublicKey
		got := p.ToPubKeyStr()

		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}
