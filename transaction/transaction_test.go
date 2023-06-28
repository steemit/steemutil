package transaction_test

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/steemit/steemutil/chain"
	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/operation"
	"github.com/steemit/steemutil/transaction"
	"github.com/steemit/steemutil/util"
	"github.com/steemit/steemutil/wif"
)

var (
	tx   *transaction.Transaction
	wifs = []string{
		"5JLw5dgQAx6rhZEgNN5C2ds1V47RweGshynFSWFbaMohsYsBvE8",
	}
	privateKeys = make([]*wif.PrivateKey, 0, len(wifs))
	publicKeys  = make([]*wif.PublicKey, 0, len(wifs))
)

func init() {
	// Prepare the transaction.
	expiration := time.Date(2016, 8, 8, 12, 24, 17, 0, time.UTC)
	tx = &transaction.Transaction{
		RefBlockNum:    36029,
		RefBlockPrefix: 1164960351,
		Expiration:     &util.Time{Time: &expiration},
	}
	tx.PushOperation(&operation.VoteOperation{
		Voter:    "xeroc",
		Author:   "xeroc",
		Permlink: "piston",
		Weight:   10000,
	})
	// Prepare the private keys
	for _, v := range wifs {
		privKey := &wif.PrivateKey{}
		err := privKey.FromWif(v)
		if err != nil {
			panic(err)
		}
		privateKeys = append(privateKeys, privKey)
	}
	// Prepare the public keys
	for _, v := range wifs {
		pubKey := &wif.PublicKey{}
		err := pubKey.FromWif(v)
		if err != nil {
			panic(err)
		}
		publicKeys = append(publicKeys, pubKey)
	}
}

func TestTransaction_MarshalTransaction(t *testing.T) {
	// The result we expect.
	expected := "bd8c5fe26f45f179a8570100057865726f63057865726f6306706973746f6e102700"

	// Marshal the transaction.
	var b bytes.Buffer
	encoder := encoder.NewEncoder(&b)

	if err := tx.MarshalTransaction(encoder); err != nil {
		t.Error(err)
	}
	got := hex.EncodeToString(b.Bytes())

	// Compare that we got with what we expect to get.
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestTransaction_Digest(t *testing.T) {
	expected := "582176b1daf89984bc8b4fdcb24ff1433d1eb114a8c4bf20fb22ad580d035889"
	digest, err := tx.Digest(chain.SteemChain)
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
	if err := tx.Sign(privateKeys, chain.SteemChain); err != nil {
		t.Error(err)
	}
	if len(tx.Signatures) != 1 {
		t.Error("expected signatures not appended to the transaction")
	}
	ok, err := tx.Verify(publicKeys, chain.SteemChain)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("verification failed")
	}
}
