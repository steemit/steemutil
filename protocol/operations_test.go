package protocol

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/steemit/steemutil/encoder"
)

func TestVoteOperation_MarshalTransaction(t *testing.T) {
	op := &VoteOperation{
		Voter:    "xeroc",
		Author:   "xeroc",
		Permlink: "piston",
		Weight:   10000,
	}

	expectedHex := "00057865726f63057865726f6306706973746f6e1027"

	var b bytes.Buffer
	encoder := encoder.NewEncoder(&b)

	if err := encoder.Encode(op); err != nil {
		t.Error(err)
	}

	serializedHex := hex.EncodeToString(b.Bytes())

	if serializedHex != expectedHex {
		t.Errorf("expected %v, got %v", expectedHex, serializedHex)
	}
}

func TestCommentOperation_MarshalTransaction(t *testing.T) {
	op := &CommentOperation{
		Author:         "ety001",
		Title:          "test post",
		Permlink:       "ety001-test-post",
		ParentAuthor:   "",
		ParentPermlink: "test",
		Body:           "test post body",
		JsonMetadata:   "{}",
	}

	expectedHex := "1300047465737406657479303031106574793030312d746573742d706f7374097465737420706f73740e7465737420706f737420626f6479027b7d"

	var b bytes.Buffer
	encoder := encoder.NewEncoder(&b)

	if err := encoder.Encode(op); err != nil {
		t.Error(err)
	}

	serializedHex := hex.EncodeToString(b.Bytes())

	if serializedHex != expectedHex {
		t.Errorf("expected %v, got %v", expectedHex, serializedHex)
	}
}
