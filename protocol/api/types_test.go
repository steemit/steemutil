package api

import (
	"encoding/json"
	"testing"
)

func TestUnmarshal_FollowCount(t *testing.T) {
	src := `{"account":"alice","follower_count":123,"following_count":45}`
	var fc FollowCountReturn
	if err := json.Unmarshal([]byte(src), &fc); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if fc.Account != "alice" || fc.FollowerCount != 123 || fc.FollowingCount != 45 {
		t.Errorf("unexpected: %+v", fc)
	}
}

func TestUnmarshal_FollowReturn(t *testing.T) {
	src := `{"follower":"alice","following":"bob","what":["blog"]}`
	var fr FollowReturn
	if err := json.Unmarshal([]byte(src), &fr); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if fr.Follower != "alice" || fr.Following != "bob" {
		t.Errorf("unexpected follower/following: %+v", fr)
	}
	if len(fr.What) != 1 || fr.What[0] != "blog" {
		t.Errorf("unexpected what: %+v", fr.What)
	}
}

func TestUnmarshal_OrderBook(t *testing.T) {
	src := `{
		"asks": [{"order_price":{"base":"1.000 SBD","quote":"1.000 STEEM"}}],
		"bids": [{"order_price":{"base":"0.990 SBD","quote":"1.000 STEEM"}}]
	}`
	var ob OrderBook
	if err := json.Unmarshal([]byte(src), &ob); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(ob.Asks) != 1 || len(ob.Bids) != 1 {
		t.Fatalf("unexpected counts: asks=%d bids=%d", len(ob.Asks), len(ob.Bids))
	}
	if ob.Asks[0].OrderPrice.Base != "1.000 SBD" {
		t.Errorf("ask base: %q", ob.Asks[0].OrderPrice.Base)
	}
	if ob.Bids[0].OrderPrice.Base != "0.990 SBD" {
		t.Errorf("bid base: %q", ob.Bids[0].OrderPrice.Base)
	}
}

func TestUnmarshal_FeedHistory(t *testing.T) {
	src := `{
		"price_history": [
			{"base":"0.500 SBD","quote":"1.000 STEEM"},
			{"base":"0.510 SBD","quote":"1.000 STEEM"}
		]
	}`
	var fh FeedHistory
	if err := json.Unmarshal([]byte(src), &fh); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(fh.PriceHistory) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(fh.PriceHistory))
	}
	last := fh.PriceHistory[len(fh.PriceHistory)-1]
	if last.Base != "0.510 SBD" {
		t.Errorf("last base: %q", last.Base)
	}
}
