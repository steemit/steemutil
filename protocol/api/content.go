package api

import "encoding/json"

// Beneficiary models one entry of a Content's beneficiaries array
// (steemd chain::beneficiary_route_type: { account, weight }). weight is a
// uint16_t permille (0..SMARTS_AUTHOR_REWARD_DENOM) and always a small JSON
// number, so a plain uint16 decodes cleanly.
type Beneficiary struct {
	Account string `json:"account"`
	Weight  uint16 `json:"weight"`
}

// VoteState models one entry of a Content's active_votes array
// (steemd tags_api::vote_state). voter and time are strings; percent is int16
// and always a small JSON number.
//
// weight (uint64_t), rshares (int64_t), and reputation (share_type) are raw
// integer C++ types that steemd serializes inconsistently: values within the
// int32 range emit as a JSON number, larger values emit as a JSON string.
// Within a single response different voters can have rshares in either form,
// so a typed field would fail to decode on one or the other. All three use
// json.RawMessage so both forms decode; same strategy as ExtendedAccount's
// reputation/withdrawn group and Manabar.CurrentMana.
type VoteState struct {
	Voter      string          `json:"voter"`
	Weight     json.RawMessage `json:"weight"`
	Rshares    json.RawMessage `json:"rshares"`
	Percent    int             `json:"percent"`
	Reputation json.RawMessage `json:"reputation"`
	Time       string          `json:"time"`
}

// Content models the return value of condenser_api.get_content (the "discussion"
// object). It mirrors steemd's condenser_api::discussion (which extends
// condenser_api::api_comment_object), carrying every key the chain emits so
// steemdb-sync's comment rescanner captures a complete snapshot in a single
// decode.
//
// Field conventions, matching the on-chain JSON exactly (verified against live
// get_content responses and steemd source):
//   - Assets ("100.000 SBD") and ISO-8601 timestamps are kept as raw strings;
//     parsing is the consumer's responsibility (see ParseAsset).
//   - The share_type / large-int group (net_rshares, abs_rshares, vote_rshares,
//     children_abs_rshares, author_rewards, author_reputation, total_vote_weight)
//     are decoded as json.RawMessage because steemd serializes them
//     inconsistently: within the int32 range as a JSON number, larger as a JSON
//     string — a typed field would fail to decode on one or the other. Same
//     strategy as ExtendedAccount's reputation/withdrawn group and Manabar.
//   - active_votes[].weight / rshares / reputation use json.RawMessage for the
//     same reason (both forms coexist within a single response).
//   - allow_votes / allow_replies / allow_curation_rewards are plain bool.
//   - first_reblogged_by / first_reblogged_on are optional on chain (omitted
//     entirely when absent); they are pointers here so absence decodes to nil.
type Content struct {
	// --- scalar string fields: identities, assets, timestamps ---
	ID                      int    `json:"id"`
	Author                  string `json:"author"`
	Permlink                string `json:"permlink"`
	Category                string `json:"category"`
	ParentAuthor            string `json:"parent_author"`
	ParentPermlink          string `json:"parent_permlink"`
	RootAuthor              string `json:"root_author"`
	RootPermlink            string `json:"root_permlink"`
	Title                   string `json:"title"`
	Body                    string `json:"body"`
	JSONMetadata            string `json:"json_metadata"`
	RootTitle               string `json:"root_title"`
	URL                     string `json:"url"`
	LastUpdate              string `json:"last_update"`
	Created                 string `json:"created"`
	Active                  string `json:"active"`
	LastPayout              string `json:"last_payout"`
	CashoutTime             string `json:"cashout_time"`
	MaxCashoutTime          string `json:"max_cashout_time"`
	Mode                    string `json:"mode"`
	PendingPayoutValue      string `json:"pending_payout_value"`
	TotalPendingPayoutValue string `json:"total_pending_payout_value"`
	TotalPayoutValue        string `json:"total_payout_value"`
	CuratorPayoutValue      string `json:"curator_payout_value"`
	MaxAcceptedPayout       string `json:"max_accepted_payout"`
	Promoted                string `json:"promoted"`

	// --- scalar int/bool fields (always small JSON numbers / true|false) ---
	BodyLength           int  `json:"body_length"`
	Depth                int  `json:"depth"`
	Children             int  `json:"children"`
	NetVotes             int  `json:"net_votes"`
	RewardWeight         int  `json:"reward_weight"`
	PercentSteemDollars  int  `json:"percent_steem_dollars"`
	AllowReplies         bool `json:"allow_replies"`
	AllowVotes           bool `json:"allow_votes"`
	AllowCurationRewards bool `json:"allow_curation_rewards"`

	// --- RawMessage: share_type / large-int, string-or-number on chain ---
	AbsRshares         json.RawMessage `json:"abs_rshares"`
	NetRshares         json.RawMessage `json:"net_rshares"`
	VoteRshares        json.RawMessage `json:"vote_rshares"`
	ChildrenAbsRshares json.RawMessage `json:"children_abs_rshares"`
	AuthorRewards      json.RawMessage `json:"author_rewards"`
	AuthorReputation   json.RawMessage `json:"author_reputation"`
	TotalVoteWeight    json.RawMessage `json:"total_vote_weight"`

	// --- nested collections ---
	ActiveVotes      []VoteState   `json:"active_votes"`
	Replies          []string      `json:"replies"`
	Beneficiaries    []Beneficiary `json:"beneficiaries"`
	RebloggedBy      []string      `json:"reblogged_by"`
	FirstRebloggedBy *string       `json:"first_reblogged_by"`
	FirstRebloggedOn *string       `json:"first_reblogged_on"`
}
