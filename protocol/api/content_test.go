package api

import (
	"encoding/json"
	"testing"
)

// TestUnmarshal_Content_FullFixture decodes a trimmed-but-type-faithful
// get_content response and asserts representative fields from each group.
//
// The fixture is captured from a live condenser_api.get_content call. It is
// trimmed (body shortened, active_votes cut to 3 entries) but deliberately
// keeps both JSON forms of the share_type fields that steemd emits
// inconsistently:
//   - top-level net_rshares / abs_rshares / vote_rshares /
//     children_abs_rshares / author_reputation come back as JSON STRINGS
//     (large values, quoted);
//   - top-level total_vote_weight comes back as a JSON NUMBER (small value);
//   - within active_votes, vote[0].rshares is a STRING while vote[1].rshares
//     is a NUMBER — both forms coexisting in one response is the headline
//     reason these fields are json.RawMessage rather than a typed scalar.
func TestUnmarshal_Content_FullFixture(t *testing.T) {
	const src = `{"id":113900934,"author":"rme","permlink":"fun-meme-puss-logo","category":"hive-129948","parent_author":"","parent_permlink":"hive-129948","title":" FUN MEME $PUSS logos","body":"<trimmed>","json_metadata":"{\"tags\":[\"meme\"]}","last_update":"2026-08-12T13:00:06","created":"2026-08-12T13:00:06","active":"2026-08-12T13:15:48","last_payout":"1970-01-01T00:00:00","depth":0,"children":1,"net_rshares":"1605718787099101","abs_rshares":"1605718787099101","vote_rshares":"1605718787099101","children_abs_rshares":"1605718787099101","cashout_time":"2026-08-19T13:00:06","max_cashout_time":"1969-12-31T23:59:59","total_vote_weight":39390711,"reward_weight":10000,"total_payout_value":"0.000 SBD","curator_payout_value":"0.000 SBD","author_rewards":0,"net_votes":271,"root_author":"rme","root_permlink":"fun-meme-puss-logo","root_title":" FUN MEME $PUSS logos","max_accepted_payout":"1000000.000 SBD","percent_steem_dollars":10000,"allow_replies":true,"allow_votes":true,"allow_curation_rewards":true,"beneficiaries":[],"url":"/hive-129948/@rme/fun-meme-puss-logo","pending_payout_value":"229.377 SBD","total_pending_payout_value":"0.000 STEEM","promoted":"0.000 STEEM","body_length":0,"replies":[],"author_reputation":"30678045089626224","active_votes":[{"voter":"xpilar","weight":8597,"rshares":"221646137396","reputation":0,"time":"2026-08-12T13:01:30","percent":1000},{"voter":"royalmacro","weight":201520,"rshares":"8514469954525","reputation":0,"time":"2026-08-12T13:05:09","percent":10000},{"voter":"xiaohui","weight":28,"rshares":2125781185,"reputation":0,"time":"2026-08-12T13:06:09","percent":5672}],"reblogged_by":[]}`

	var c Content
	if err := json.Unmarshal([]byte(src), &c); err != nil {
		t.Fatalf("unmarshal full fixture failed: %v", err)
	}

	// --- scalar string fields: identities / assets / timestamps ---
	if c.ID != 113900934 {
		t.Errorf("id: got %d", c.ID)
	}
	if c.Author != "rme" || c.Permlink != "fun-meme-puss-logo" {
		t.Errorf("author/permlink: got %q/%q", c.Author, c.Permlink)
	}
	if c.Category != "hive-129948" {
		t.Errorf("category: got %q", c.Category)
	}
	if c.ParentAuthor != "" {
		t.Errorf("parent_author: want empty, got %q", c.ParentAuthor)
	}
	if c.RootAuthor != "rme" || c.RootPermlink != "fun-meme-puss-logo" {
		t.Errorf("root_author/root_permlink: got %q/%q", c.RootAuthor, c.RootPermlink)
	}
	if c.Created != "2026-08-12T13:00:06" {
		t.Errorf("created: got %q", c.Created)
	}
	if c.CashoutTime != "2026-08-19T13:00:06" {
		t.Errorf("cashout_time: got %q", c.CashoutTime)
	}
	if c.LastPayout != "1970-01-01T00:00:00" {
		t.Errorf("last_payout: got %q", c.LastPayout)
	}
	if c.PendingPayoutValue != "229.377 SBD" {
		t.Errorf("pending_payout_value: got %q", c.PendingPayoutValue)
	}
	if c.MaxAcceptedPayout != "1000000.000 SBD" {
		t.Errorf("max_accepted_payout: got %q", c.MaxAcceptedPayout)
	}
	if c.Promoted != "0.000 STEEM" {
		t.Errorf("promoted: got %q", c.Promoted)
	}
	if c.URL != "/hive-129948/@rme/fun-meme-puss-logo" {
		t.Errorf("url: got %q", c.URL)
	}

	// --- scalar int/bool fields ---
	if c.Depth != 0 || c.Children != 1 || c.NetVotes != 271 {
		t.Errorf("depth/children/net_votes: got %d/%d/%d", c.Depth, c.Children, c.NetVotes)
	}
	if c.RewardWeight != 10000 || c.PercentSteemDollars != 10000 {
		t.Errorf("reward_weight/percent_steem_dollars: got %d/%d", c.RewardWeight, c.PercentSteemDollars)
	}
	if !c.AllowReplies || !c.AllowVotes || !c.AllowCurationRewards {
		t.Errorf("allow_*: want all true, got replies=%v votes=%v curation=%v", c.AllowReplies, c.AllowVotes, c.AllowCurationRewards)
	}

	// --- RawMessage share_type group: captured verbatim regardless of form.
	// net_rshares/abs_rshares/etc. here are JSON STRINGS (quoted, large values);
	// total_vote_weight is a JSON NUMBER (small value). Both must decode and
	// round-trip bytes. ---
	if string(c.NetRshares) != `"1605718787099101"` {
		t.Errorf("net_rshares raw: got %q", string(c.NetRshares))
	}
	if string(c.AbsRshares) != `"1605718787099101"` {
		t.Errorf("abs_rshares raw: got %q", string(c.AbsRshares))
	}
	if string(c.VoteRshares) != `"1605718787099101"` {
		t.Errorf("vote_rshares raw: got %q", string(c.VoteRshares))
	}
	if string(c.ChildrenAbsRshares) != `"1605718787099101"` {
		t.Errorf("children_abs_rshares raw: got %q", string(c.ChildrenAbsRshares))
	}
	if string(c.AuthorReputation) != `"30678045089626224"` {
		t.Errorf("author_reputation raw: got %q", string(c.AuthorReputation))
	}
	if string(c.TotalVoteWeight) != "39390711" {
		t.Errorf("total_vote_weight raw: got %q", string(c.TotalVoteWeight))
	}
	// author_rewards=0 arrives as a JSON number here; must also decode.
	if string(c.AuthorRewards) != "0" {
		t.Errorf("author_rewards raw: got %q", string(c.AuthorRewards))
	}

	// --- active_votes: MIXED rshares forms within one response.
	// vote[0].rshares is a STRING ("221646137396"), vote[2].rshares is a
	// NUMBER (2125781185) — the headline reason VoteState.Rshares is
	// json.RawMessage rather than a typed field. ---
	if len(c.ActiveVotes) != 3 {
		t.Fatalf("active_votes len: got %d", len(c.ActiveVotes))
	}
	if c.ActiveVotes[0].Voter != "xpilar" {
		t.Errorf("active_votes[0].voter: got %q", c.ActiveVotes[0].Voter)
	}
	if c.ActiveVotes[0].Percent != 1000 {
		t.Errorf("active_votes[0].percent: got %d", c.ActiveVotes[0].Percent)
	}
	if c.ActiveVotes[0].Time != "2026-08-12T13:01:30" {
		t.Errorf("active_votes[0].time: got %q", c.ActiveVotes[0].Time)
	}
	if string(c.ActiveVotes[0].Rshares) != `"221646137396"` {
		t.Errorf("active_votes[0].rshares (string form): got %q", string(c.ActiveVotes[0].Rshares))
	}
	if string(c.ActiveVotes[0].Weight) != "8597" {
		t.Errorf("active_votes[0].weight (number form): got %q", string(c.ActiveVotes[0].Weight))
	}
	if string(c.ActiveVotes[2].Rshares) != "2125781185" {
		t.Errorf("active_votes[2].rshares (number form): got %q", string(c.ActiveVotes[2].Rshares))
	}
	// reputation=0 arrives as a JSON number; must decode too.
	if string(c.ActiveVotes[0].Reputation) != "0" {
		t.Errorf("active_votes[0].reputation (number form): got %q", string(c.ActiveVotes[0].Reputation))
	}

	// --- nested collections: empty arrays stay non-nil / length 0 ---
	if len(c.Beneficiaries) != 0 {
		t.Errorf("beneficiaries: want len 0, got %d", len(c.Beneficiaries))
	}
	if len(c.Replies) != 0 {
		t.Errorf("replies: want len 0, got %d", len(c.Replies))
	}
	if len(c.RebloggedBy) != 0 {
		t.Errorf("reblogged_by: want len 0, got %d", len(c.RebloggedBy))
	}

	// --- optional fields: omitted on chain -> nil pointers ---
	if c.FirstRebloggedBy != nil {
		t.Errorf("first_reblogged_by: want nil, got %v", c.FirstRebloggedBy)
	}
	if c.FirstRebloggedOn != nil {
		t.Errorf("first_reblogged_on: want nil, got %v", c.FirstRebloggedOn)
	}
}

// TestUnmarshal_Content_OldPost verifies a post whose share_type values are all
// small (e.g. an already-paid-out legacy post) decodes with everything arriving
// as JSON NUMBERS — the other half of the ambiguity coverage.
func TestUnmarshal_Content_OldPost(t *testing.T) {
	const src = `{"id":0,"author":"steemit","permlink":"firstpost","category":"meta","parent_author":"","parent_permlink":"meta","title":"Welcome to Steem!","body":"Steemit is a social media platform.","json_metadata":"","last_update":"2016-03-30T18:30:18","created":"2016-03-30T18:30:18","active":"2024-02-23T00:39:36","last_payout":"2016-08-24T19:59:42","depth":0,"children":441,"net_rshares":0,"abs_rshares":0,"vote_rshares":0,"children_abs_rshares":"28804959899290","cashout_time":"1969-12-31T23:59:59","max_cashout_time":"1969-12-31T23:59:59","total_vote_weight":0,"reward_weight":10000,"total_payout_value":"0.942 SBD","curator_payout_value":"0.756 SBD","author_rewards":3548,"net_votes":90,"root_author":"steemit","root_permlink":"firstpost","root_title":"Welcome to Steem!","max_accepted_payout":"1000000.000 SBD","percent_steem_dollars":10000,"allow_replies":true,"allow_votes":true,"allow_curation_rewards":true,"beneficiaries":[],"url":"/meta/@steemit/firstpost","pending_payout_value":"0.000 SBD","total_pending_payout_value":"0.000 STEEM","promoted":"0.000 STEEM","body_length":0,"replies":[],"author_reputation":"12944616889","active_votes":[{"voter":"dantheman","weight":"32866333630","rshares":375241,"reputation":0,"time":"2016-04-07T19:15:36","percent":100}],"reblogged_by":[]}`

	var c Content
	if err := json.Unmarshal([]byte(src), &c); err != nil {
		t.Fatalf("unmarshal old post failed: %v", err)
	}

	// net_rshares=0 arrives as a JSON number here (small value).
	if string(c.NetRshares) != "0" {
		t.Errorf("net_rshares (number form): got %q", string(c.NetRshares))
	}
	if string(c.AuthorRewards) != "3548" {
		t.Errorf("author_rewards (number form): got %q", string(c.AuthorRewards))
	}
	// children_abs_rshares still arrives as a string even on this old post.
	if string(c.ChildrenAbsRshares) != `"28804959899290"` {
		t.Errorf("children_abs_rshares (string form): got %q", string(c.ChildrenAbsRshares))
	}
	// author_reputation is a string on both old and new posts.
	if string(c.AuthorReputation) != `"12944616889"` {
		t.Errorf("author_reputation (string form): got %q", string(c.AuthorReputation))
	}

	// active_votes[0].weight is a STRING here ("32866333630"), rshares a NUMBER.
	if len(c.ActiveVotes) != 1 {
		t.Fatalf("active_votes len: got %d", len(c.ActiveVotes))
	}
	if string(c.ActiveVotes[0].Weight) != `"32866333630"` {
		t.Errorf("active_votes[0].weight (string form): got %q", string(c.ActiveVotes[0].Weight))
	}
	if string(c.ActiveVotes[0].Rshares) != "375241" {
		t.Errorf("active_votes[0].rshares (number form): got %q", string(c.ActiveVotes[0].Rshares))
	}

	if c.ID != 0 || c.Author != "steemit" {
		t.Errorf("id/author: got %d/%q", c.ID, c.Author)
	}
}
