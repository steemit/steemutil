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
