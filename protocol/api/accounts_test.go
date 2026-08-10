package api

import (
	"encoding/json"
	"testing"
)

// TestUnmarshal_ExtendedAccount verifies a realistic condenser_api.get_accounts
// response entry deserializes into ExtendedAccount, including:
//   - reputation as both a string and a number (kept as RawMessage)
//   - the posting authority with key_auths and weight_threshold
func TestUnmarshal_ExtendedAccount(t *testing.T) {
	// reputation as a JSON string (legacy large-number form)
	srcStr := `{
		"name": "alice",
		"created": "2016-03-24T17:00:21",
		"reputation": "6217887123456",
		"voting_power": 9800,
		"balance": "1234.567 STEEM",
		"posting": {
			"weight_threshold": 1,
			"account_auths": [],
			"key_auths": [["STM7jNh5ejQoqHqWcGWFJ1v4F5CzsG3EiBuz1VooCng1cH5QpJD27", 1]]
		},
		"active": {"weight_threshold": 1, "account_auths": [], "key_auths": []},
		"owner": {"weight_threshold": 1, "account_auths": [], "key_auths": []}
	}`

	var acct ExtendedAccount
	if err := json.Unmarshal([]byte(srcStr), &acct); err != nil {
		t.Fatalf("unmarshal (reputation as string) failed: %v", err)
	}
	if acct.Name != "alice" {
		t.Errorf("name: want alice, got %q", acct.Name)
	}
	if acct.VotingPower != 9800 {
		t.Errorf("voting_power: want 9800, got %d", acct.VotingPower)
	}
	if acct.Balance != "1234.567 STEEM" {
		t.Errorf("balance mismatch: %q", acct.Balance)
	}
	// reputation raw message retains the quoted string
	if string(acct.Reputation) != `"6217887123456"` {
		t.Errorf("reputation raw: %q", string(acct.Reputation))
	}
	if acct.Posting.WeightThreshold != 1 || len(acct.Posting.KeyAuths) != 1 {
		t.Errorf("posting authority wrong: %+v", acct.Posting)
	}
	if acct.Posting.KeyAuths[0].PubKey != "STM7jNh5ejQoqHqWcGWFJ1v4F5CzsG3EiBuz1VooCng1cH5QpJD27" {
		t.Errorf("posting key wrong: %q", acct.Posting.KeyAuths[0].PubKey)
	}
	if acct.Posting.KeyAuths[0].Weight != 1 {
		t.Errorf("posting key weight wrong: %d", acct.Posting.KeyAuths[0].Weight)
	}

	// reputation as a JSON number must also decode without error
	srcNum := `{
		"name": "bob",
		"created": "2020-01-01T00:00:00",
		"reputation": 6217887123456,
		"voting_power": 0,
		"balance": "0.000 STEEM",
		"posting": {"weight_threshold": 1, "account_auths": [], "key_auths": []}
	}`
	var acct2 ExtendedAccount
	if err := json.Unmarshal([]byte(srcNum), &acct2); err != nil {
		t.Fatalf("unmarshal (reputation as number) failed: %v", err)
	}
	if string(acct2.Reputation) != "6217887123456" {
		t.Errorf("reputation numeric raw: %q", string(acct2.Reputation))
	}
}

// TestKeyAuth_MarshalRoundTrip verifies MarshalJSON is the inverse of
// UnmarshalJSON: the nested-array wire shape survives a round trip.
func TestKeyAuth_MarshalRoundTrip(t *testing.T) {
	in := KeyAuth{PubKey: "STM7jNh5ejQoqHqWcGWFJ1v4F5CzsG3EiBuz1VooCng1cH5QpJD27", Weight: 3}

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	// wire shape must be the nested array, not a flat object
	const want = `["STM7jNh5ejQoqHqWcGWFJ1v4F5CzsG3EiBuz1VooCng1cH5QpJD27",3]`
	if string(data) != want {
		t.Errorf("marshal shape\nwant: %s\ngot:  %s", want, string(data))
	}

	var out KeyAuth
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out != in {
		t.Errorf("round-trip mismatch: want %+v, got %+v", in, out)
	}
}

// TestAccountAuthEntry_MarshalRoundTrip verifies AccountAuthEntry round-trips
// through its nested-array wire form.
func TestAccountAuthEntry_MarshalRoundTrip(t *testing.T) {
	in := AccountAuthEntry{Name: "alice", Weight: 2}

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	const want = `["alice",2]`
	if string(data) != want {
		t.Errorf("marshal shape\nwant: %s\ngot:  %s", want, string(data))
	}

	var out AccountAuthEntry
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out != in {
		t.Errorf("round-trip mismatch: want %+v, got %+v", in, out)
	}
}

// TestKeyAuth_UnmarshalMalformed verifies that malformed key_auths wire input
// is rejected with an error rather than silently producing a zero value.
func TestKeyAuth_UnmarshalMalformed(t *testing.T) {
	cases := map[string]string{
		"not an array":       `"STMxxx"`,
		"single element":     `["STMxxx"]`,
		"three elements":     `["STMxxx", 1, 2]`,
		"non-string key":     `[123, 1]`,
		"non-numeric weight": `["STMxxx", "heavy"]`,
		"null":               `null`,
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			var k KeyAuth
			// `null` unmarshals to a zero value without error (standard json
			// behavior); every other malformed case must error.
			if src == `null` {
				if err := json.Unmarshal([]byte(src), &k); err != nil {
					t.Errorf("expected null to yield zero value, got error: %v", err)
				}
				return
			}
			if err := json.Unmarshal([]byte(src), &k); err == nil {
				t.Errorf("expected error for malformed input %s", src)
			}
		})
	}
}

// TestAccountAuthEntry_UnmarshalMalformed verifies malformed account_auths
// input is rejected.
func TestAccountAuthEntry_UnmarshalMalformed(t *testing.T) {
	cases := map[string]string{
		"not an array":       `"alice"`,
		"single element":     `["alice"]`,
		"three elements":     `["alice", 1, 2]`,
		"non-numeric weight": `["alice", "heavy"]`,
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			var a AccountAuthEntry
			if err := json.Unmarshal([]byte(src), &a); err == nil {
				t.Errorf("expected error for malformed input %s", src)
			}
		})
	}
}

// TestUnmarshal_ExtendedAccount_FullFixture decodes a real
// condenser_api.get_accounts(["steemit"]) response entry (all 65 on-chain
// keys) into ExtendedAccount and asserts representative fields from every
// group: the original conveyor subset, scalar string/bool/int fields, the typed
// Manabar, the RawMessage share_type group, the MIXED proxied_vsf_votes array,
// and the empty extended_account collections.
//
// The fixture is a verbatim public on-chain snapshot; it pins the full wire
// shape so any future field/type drift is caught here rather than downstream.
func TestUnmarshal_ExtendedAccount_FullFixture(t *testing.T) {
	// 65-key response captured from a live condenser_api.get_accounts call.
	const src = `{"id":28,"name":"steemit","owner":{"weight_threshold":1,"account_auths":[],"key_auths":[["STM5nqRQyxpAy1NzMk6Ho5cjZbZ2kUV3WTg4ckMSNR7VacHib6nAm",1]]},"active":{"weight_threshold":1,"account_auths":[],"key_auths":[["STM8mXPsLVsWBeYNun45r3pHGq3iyASN6FS27mXXjXdjrVEW5f4Ad",1]]},"posting":{"weight_threshold":1,"account_auths":[],"key_auths":[["STM6P961NP2LyqMhUyi1XraJ4gsEGwkwWy3x4oH74H4y9aC1KLEbC",1]]},"memo_key":"STM8bYJg6FVsNrxLPRpzJf2KPh32ZSjTNJyAgEbgZHsMiWczNcoNi","json_metadata":"","posting_json_metadata":"","proxy":"dev365","last_owner_update":"2020-03-21T11:53:12","last_account_update":"2020-03-21T11:53:12","created":"2016-03-24T17:00:21","mined":true,"recovery_account":"steem","reset_account":"null","last_account_recovery":"1970-01-01T00:00:00","comment_count":0,"lifetime_vote_count":0,"post_count":1,"can_vote":true,"voting_manabar":{"current_mana":"5640467421072","last_update_time":1699768485},"downvote_manabar":{"current_mana":"1410116855268","last_update_time":1699768485},"voting_power":0,"balance":"398420.141 STEEM","savings_balance":"0.000 STEEM","sbd_balance":"153099.071 SBD","sbd_seconds":"15326354649","sbd_seconds_last_update":"2023-11-13T03:39:12","sbd_last_interest_payment":"2023-11-12T05:54:45","savings_sbd_balance":"0.000 SBD","savings_sbd_seconds":"0","savings_sbd_seconds_last_update":"1970-01-01T00:00:00","savings_sbd_last_interest_payment":"1970-01-01T00:00:00","savings_withdraw_requests":0,"reward_sbd_balance":"0.000 SBD","reward_steem_balance":"0.000 STEEM","reward_vesting_balance":"0.000000 VESTS","reward_vesting_steem":"0.000 STEEM","curation_rewards":2812964,"posting_rewards":3548,"vesting_shares":"5640467.421072 VESTS","delegated_vesting_shares":"0.000000 VESTS","received_vesting_shares":"0.000000 VESTS","vesting_withdraw_rate":"0.000000 VESTS","next_vesting_withdrawal":"1969-12-31T23:59:59","withdrawn":"57269412490896164","to_withdraw":"57269412490896164","withdraw_routes":0,"proxied_vsf_votes":["452574069424",0,0,0],"witnesses_voted_for":0,"last_post":"2016-03-30T18:30:18","last_root_post":"2016-03-30T18:30:18","last_vote_time":"2020-03-27T08:46:24","post_bandwidth":0,"pending_claimed_accounts":0,"vesting_balance":"0.000 STEEM","reputation":"12944616889","transfer_history":[],"market_history":[],"post_history":[],"vote_history":[],"other_history":[],"witness_votes":[],"tags_usage":[],"guest_bloggers":[]}`

	var acct ExtendedAccount
	if err := json.Unmarshal([]byte(src), &acct); err != nil {
		t.Fatalf("unmarshal full fixture failed: %v", err)
	}

	// --- original conveyor-facing subset (must keep working) ---
	if acct.Name != "steemit" {
		t.Errorf("name: got %q", acct.Name)
	}
	if acct.Created != "2016-03-24T17:00:21" {
		t.Errorf("created: got %q", acct.Created)
	}
	if acct.Balance != "398420.141 STEEM" {
		t.Errorf("balance: got %q", acct.Balance)
	}
	if string(acct.Reputation) != `"12944616889"` {
		t.Errorf("reputation raw: got %q", string(acct.Reputation))
	}
	if acct.VotingPower != 0 {
		t.Errorf("voting_power: got %d", acct.VotingPower)
	}
	if acct.Owner.WeightThreshold != 1 || len(acct.Owner.KeyAuths) != 1 || acct.Owner.KeyAuths[0].Weight != 1 {
		t.Errorf("owner authority wrong: %+v", acct.Owner)
	}

	// --- newly added scalar string fields (assets / timestamps / names) ---
	if acct.Proxy != "dev365" {
		t.Errorf("proxy: got %q", acct.Proxy)
	}
	if acct.RewardSBDBalance != "0.000 SBD" {
		t.Errorf("reward_sbd_balance: got %q", acct.RewardSBDBalance)
	}
	if acct.RewardVestingSteem != "0.000 STEEM" {
		t.Errorf("reward_vesting_steem: got %q", acct.RewardVestingSteem)
	}
	if acct.MemoKey != "STM8bYJg6FVsNrxLPRpzJf2KPh32ZSjTNJyAgEbgZHsMiWczNcoNi" {
		t.Errorf("memo_key: got %q", acct.MemoKey)
	}
	if acct.VestingShares != "5640467.421072 VESTS" {
		t.Errorf("vesting_shares: got %q", acct.VestingShares)
	}

	// --- newly added scalar bool/int fields ---
	if acct.ID != 28 {
		t.Errorf("id: got %d", acct.ID)
	}
	if !acct.Mined {
		t.Errorf("mined: want true")
	}
	if !acct.CanVote {
		t.Errorf("can_vote: want true")
	}
	if acct.PostCount != 1 {
		t.Errorf("post_count: got %d", acct.PostCount)
	}
	if acct.SavingsWithdrawRequests != 0 {
		t.Errorf("savings_withdraw_requests: got %d", acct.SavingsWithdrawRequests)
	}

	// --- typed Manabar: current_mana is a string, last_update_time a number ---
	if acct.VotingManabar.CurrentMana != "5640467421072" {
		t.Errorf("voting_manabar.current_mana: got %q", acct.VotingManabar.CurrentMana)
	}
	if acct.VotingManabar.LastUpdateTime != 1699768485 {
		t.Errorf("voting_manabar.last_update_time: got %d", acct.VotingManabar.LastUpdateTime)
	}
	if acct.DownvoteManabar.LastUpdateTime != 1699768485 {
		t.Errorf("downvote_manabar.last_update_time: got %d", acct.DownvoteManabar.LastUpdateTime)
	}

	// --- RawMessage share_type group: captured verbatim regardless of form.
	// curation_rewards/posting_rewards here are JSON numbers; withdrawn/to_withdraw
	// are JSON strings — both must decode without error and round-trip bytes. ---
	if string(acct.Withdrawn) != `"57269412490896164"` {
		t.Errorf("withdrawn raw: got %q", string(acct.Withdrawn))
	}
	if string(acct.ToWithdraw) != `"57269412490896164"` {
		t.Errorf("to_withdraw raw: got %q", string(acct.ToWithdraw))
	}
	if string(acct.CurationRewards) != "2812964" {
		t.Errorf("curation_rewards raw: got %q", string(acct.CurationRewards))
	}
	if string(acct.PostingRewards) != "3548" {
		t.Errorf("posting_rewards raw: got %q", string(acct.PostingRewards))
	}
	if string(acct.PostBandwidth) != "0" {
		t.Errorf("post_bandwidth raw: got %q", string(acct.PostBandwidth))
	}

	// --- proxied_vsf_votes: MIXED ["452574069424",0,0,0] — the headline reason
	// this field is []json.RawMessage rather than a typed slice. ---
	if len(acct.ProxiedVSFVotes) != 4 {
		t.Fatalf("proxied_vsf_votes len: got %d", len(acct.ProxiedVSFVotes))
	}
	if want := `"452574069424"`; string(acct.ProxiedVSFVotes[0]) != want {
		t.Errorf("proxied_vsf_votes[0]: got %q want %q", string(acct.ProxiedVSFVotes[0]), want)
	}
	if want := "0"; string(acct.ProxiedVSFVotes[1]) != want {
		t.Errorf("proxied_vsf_votes[1]: got %q want %q", string(acct.ProxiedVSFVotes[1]), want)
	}

	// --- extended_account collections: currently empty [], captured raw. ---
	for _, tc := range []struct {
		name string
		raw  json.RawMessage
	}{
		{"transfer_history", acct.TransferHistory},
		{"market_history", acct.MarketHistory},
		{"post_history", acct.PostHistory},
		{"vote_history", acct.VoteHistory},
		{"other_history", acct.OtherHistory},
		{"tags_usage", acct.TagsUsage},
		{"guest_bloggers", acct.GuestBloggers},
	} {
		if string(tc.raw) != "[]" {
			t.Errorf("%s: got %q want []", tc.name, string(tc.raw))
		}
	}

	// round-trip: re-marshaling the typed Manabar must reproduce the wire shape.
	vm, err := json.Marshal(acct.VotingManabar)
	if err != nil {
		t.Fatalf("marshal voting_manabar: %v", err)
	}
	const wantVM = `{"current_mana":"5640467421072","last_update_time":1699768485}`
	if string(vm) != wantVM {
		t.Errorf("voting_manabar marshal: got %s want %s", string(vm), wantVM)
	}
}
