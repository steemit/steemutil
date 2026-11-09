package wif

import (
	"bytes"
	"strings"
	"testing"

	"github.com/steemit/steemutil/consts"
)

// Known key pair: the deterministic Steem testnet initminer key.
// Derived from sha256("init_key"); the public key payload is identical on
// mainnet (STM...) and testnet (TST...) — only the prefix differs.
const testWIF = "5JNHfZYKGaomSFvd4NUdQ9qMcEAC43kujbfjueTHpVapX1Kzq2n"
const testPubPayload = "6LLegbAgLAy28EHrffBVuANFWcFgmqRMW13wBmTExqFE9SCkg4"

func TestAddressPrefixDefaultSTM(t *testing.T) {
	if consts.ADDRESS_PREFIX != "STM" {
		t.Fatalf("default ADDRESS_PREFIX = %q, want STM", consts.ADDRESS_PREFIX)
	}
	pub := &PublicKey{}
	if err := pub.FromWif(testWIF); err != nil {
		t.Fatalf("FromWif: %v", err)
	}
	got := pub.ToStr()
	if want := "STM" + testPubPayload; got != want {
		t.Fatalf("ToStr() = %q, want %q", got, want)
	}
}

func TestAddressPrefixTestnetTSTRoundTrip(t *testing.T) {
	if err := consts.SetAddressPrefix("TST"); err != nil {
		t.Fatalf("SetAddressPrefix: %v", err)
	}
	defer func() {
		if err := consts.SetAddressPrefix("STM"); err != nil {
			t.Fatalf("restore prefix: %v", err)
		}
	}()

	// Encoding under TST must only change the prefix, never the payload.
	pub := &PublicKey{}
	if err := pub.FromWif(testWIF); err != nil {
		t.Fatalf("FromWif: %v", err)
	}
	tstStr := pub.ToStr()
	if want := "TST" + testPubPayload; tstStr != want {
		t.Fatalf("ToStr() under TST = %q, want %q", tstStr, want)
	}

	// A TST-prefixed key must parse back to the same raw key bytes.
	parsed := &PublicKey{}
	if err := parsed.FromStr(tstStr); err != nil {
		t.Fatalf("FromStr(TST key) failed: %v", err)
	}
	if !bytes.Equal(parsed.ToByte(), pub.ToByte()) {
		t.Fatalf("round-trip key bytes differ: %x vs %x", parsed.ToByte(), pub.ToByte())
	}

	// An STM-prefixed key must be rejected while TST is active.
	stmStr := "STM" + testPubPayload
	if err := (&PublicKey{}).FromStr(stmStr); err == nil {
		t.Fatal("FromStr(STM key) under TST prefix: expected prefix error, got nil")
	} else if !strings.Contains(err.Error(), "prefix") {
		t.Fatalf("FromStr(STM key) under TST prefix: unexpected error: %v", err)
	}
}
