package transaction

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/steemit/steemutil/protocol"
	"github.com/steemit/steemutil/wif"
)

var tx *Transaction

func init() {
	// Prepare the transaction.
	expiration := time.Date(2016, 8, 8, 12, 24, 17, 0, time.UTC)
	tx = &Transaction{
		RefBlockNum:    36029,
		RefBlockPrefix: 1164960351,
		Expiration:     &protocol.Time{Time: &expiration},
	}
	tx.PushOperation(&protocol.VoteOperation{
		Voter:    "xeroc",
		Author:   "xeroc",
		Permlink: "piston",
		Weight:   10000,
	})
}

var wifs = []string{
	"5JLw5dgQAx6rhZEgNN5C2ds1V47RweGshynFSWFbaMohsYsBvE8",
}

var privateKeys = make([]*wif.PrivateKey, 0, len(wifs))

func init() {
	for _, v := range wifs {
		privKey := &wif.PrivateKey{}
		err := privKey.FromWif(v)
		if err != nil {
			panic(err)
		}
		privateKeys = append(privateKeys, privKey)
	}
}

var publicKeys = make([]*wif.PublicKey, 0, len(wifs))

func init() {
	for _, v := range wifs {
		pubKey := &wif.PublicKey{}
		err := pubKey.FromWif(v)
		if err != nil {
			panic(err)
		}
		publicKeys = append(publicKeys, pubKey)
	}
}

func TestTransaction_Digest(t *testing.T) {
	expected := "582176b1daf89984bc8b4fdcb24ff1433d1eb114a8c4bf20fb22ad580d035889"

	stx := NewSignedTransaction(tx)

	digest, err := stx.Digest(SteemChain)
	if err != nil {
		t.Error(err)
	}

	got := hex.EncodeToString(digest)
	if got != expected {
		t.Errorf("got %v, expected %v", got, expected)
	}
}

func TestTransaction_SignAndVerify(t *testing.T) {
	tx.Signatures = nil
	defer func() {
		tx.Signatures = nil
	}()

	stx := NewSignedTransaction(tx)
	if err := stx.Sign(privateKeys, SteemChain); err != nil {
		t.Error(err)
	}

	if len(tx.Signatures) != 1 {
		t.Error("expected signatures not appended to the transaction")
	}

	ok, err := stx.Verify(publicKeys, SteemChain)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("verification failed")
	}
}

// TestTransaction_VerifyNegativeCases pins the security-critical behaviour
// that Verify must be fail-closed. Before the fix, Verify iterated over an
// empty slice and returned (true, nil) for any input, including tampered
// transactions, empty signature lists, and wrong keys.
func TestTransaction_VerifyNegativeCases(t *testing.T) {
	tx.Signatures = nil
	defer func() {
		tx.Signatures = nil
	}()

	stx := NewSignedTransaction(tx)
	if err := stx.Sign(privateKeys, SteemChain); err != nil {
		t.Fatalf("Sign failed: %v", err)
	}
	goodSigs := append([]string{}, tx.Signatures...)

	const origWeight = 10000

	// reSign returns the transaction to a clean, freshly-signed baseline.
	reSign := func() {
		tx.Signatures = nil
		if err := stx.Sign(privateKeys, SteemChain); err != nil {
			t.Fatalf("reSign failed: %v", err)
		}
	}

	// Case 1: tampered transaction — mutating the signed payload after signing
	// changes the digest, so recovery must produce a different key.
	t.Run("tampered_digest", func(t *testing.T) {
		tx.Signatures = goodSigs
		defer reSign()

		tx.PushOperation(&protocol.VoteOperation{
			Voter:    "xeroc",
			Author:   "xeroc",
			Permlink: "piston",
			Weight:   origWeight + 1, // mutate signed payload
		})

		ok, err := stx.Verify(publicKeys, SteemChain)
		if ok {
			t.Errorf("Verify returned true for a tampered transaction (err=%v)", err)
		}
		if err == nil {
			t.Errorf("Verify returned nil error for a tampered transaction")
		}
	})

	// Case 2: empty signature list must be fail-closed.
	t.Run("empty_signatures", func(t *testing.T) {
		tx.Signatures = nil
		defer reSign()

		ok, err := stx.Verify(publicKeys, SteemChain)
		if ok {
			t.Errorf("Verify returned true with no signatures (err=%v)", err)
		}
		if err == nil {
			t.Errorf("Verify returned nil error with no signatures")
		}
	})

	// Case 3: a signature that does not correspond to the supplied pubKey.
	t.Run("wrong_pubkey", func(t *testing.T) {
		tx.Signatures = goodSigs
		defer reSign()

		// An unrelated valid Steem public key (standard test vector).
		wrongKey := &wif.PublicKey{}
		if err := wrongKey.FromStr("STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA"); err != nil {
			t.Fatalf("failed to parse wrong pubkey: %v", err)
		}

		ok, err := stx.Verify([]*wif.PublicKey{wrongKey}, SteemChain)
		if ok {
			t.Errorf("Verify returned true against an unrelated pubkey (err=%v)", err)
		}
		if err == nil {
			t.Errorf("Verify returned nil error against an unrelated pubkey")
		}
	})

	// Case 4: garbage signature bytes must not pass.
	t.Run("malformed_signature", func(t *testing.T) {
		tx.Signatures = []string{"deadbeef"}
		defer reSign()

		ok, err := stx.Verify(publicKeys, SteemChain)
		if ok {
			t.Errorf("Verify returned true for a malformed signature (err=%v)", err)
		}
		if err == nil {
			t.Errorf("Verify returned nil error for a malformed signature")
		}
	})

	// Sanity: the baseline verifies, proving the negatives above are about the
	// inputs, not a broken fixture.
	t.Run("baseline_still_verifies", func(t *testing.T) {
		reSign()
		ok, err := stx.Verify(publicKeys, SteemChain)
		if err != nil || !ok {
			t.Fatalf("baseline should verify, got ok=%v err=%v", ok, err)
		}
	})
}
