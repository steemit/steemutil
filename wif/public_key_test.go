package wif

import (
	"testing"
)

func TestFromStrToStr(t *testing.T) {
	for _, d := range data {
		p := &PublicKey{}
		err := p.FromStr(d.PublicKey)
		if err != nil {
			t.Error(err)
		}

		expected := d.PublicKey
		got := p.ToStr()

		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}

func TestPublicKeyFromWif(t *testing.T) {
	for _, d := range data {
		p := &PublicKey{}
		err := p.FromWif(d.WIF)
		if err != nil {
			t.Error(err)
		}

		expected := d.PublicKey
		got := p.ToStr()

		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}
