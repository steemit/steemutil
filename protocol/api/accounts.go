package api

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// AuthorityWeight is the weight assigned to a key or account in a Steem
// authority. Steem encodes it as a 16-bit integer.
type AuthorityWeight uint16

// KeyAuth is a single (public key, weight) entry in an account authority.
//
// On the wire (condenser_api.get_accounts), key_auths is serialized as a JSON
// array of two-element arrays: [["STMxxx", 1], ...]. KeyAuth implements a
// custom unmarshaler that flattens each [key, weight] pair into its fields,
// so consumers can work with a typed slice rather than nested raw arrays.
type KeyAuth struct {
	PubKey string
	Weight AuthorityWeight
}

// UnmarshalJSON parses the wire form ["STMxxx", 1] into KeyAuth.
func (k *KeyAuth) UnmarshalJSON(data []byte) error {
	// A JSON null leaves the value at its zero value (standard json behavior).
	if string(data) == "null" {
		return nil
	}
	var pair []json.RawMessage
	if err := json.Unmarshal(data, &pair); err != nil {
		return errors.Wrap(err, "key_auths entry must be a [key, weight] array")
	}
	if len(pair) != 2 {
		return errors.Errorf("key_auths entry must have exactly 2 elements, got %d", len(pair))
	}
	var key string
	if err := json.Unmarshal(pair[0], &key); err != nil {
		return errors.Wrap(err, "invalid public key in key_auths")
	}
	var weight AuthorityWeight
	if err := json.Unmarshal(pair[1], &weight); err != nil {
		return errors.Wrap(err, "invalid weight in key_auths")
	}
	k.PubKey = key
	k.Weight = weight
	return nil
}

// MarshalJSON emits the wire form ["STMxxx", 1], the inverse of UnmarshalJSON,
// so a round-trip through JSON preserves the condenser_api array-of-pairs shape.
func (k KeyAuth) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{k.PubKey, k.Weight})
}

// AccountAuthEntry is a weighted account name in an account authority. On the
// wire, account_auths is [["name", weight], ...], parsed the same way as
// key_auths.
type AccountAuthEntry struct {
	Name   string
	Weight AuthorityWeight
}

// UnmarshalJSON parses the wire form ["name", 1] into AccountAuthEntry.
func (a *AccountAuthEntry) UnmarshalJSON(data []byte) error {
	// A JSON null leaves the value at its zero value (standard json behavior).
	if string(data) == "null" {
		return nil
	}
	var pair []json.RawMessage
	if err := json.Unmarshal(data, &pair); err != nil {
		return errors.Wrap(err, "account_auths entry must be a [name, weight] array")
	}
	if len(pair) != 2 {
		return errors.Errorf("account_auths entry must have exactly 2 elements, got %d", len(pair))
	}
	var name string
	if err := json.Unmarshal(pair[0], &name); err != nil {
		return errors.Wrap(err, "invalid account name in account_auths")
	}
	var weight AuthorityWeight
	if err := json.Unmarshal(pair[1], &weight); err != nil {
		return errors.Wrap(err, "invalid weight in account_auths")
	}
	a.Name = name
	a.Weight = weight
	return nil
}

// MarshalJSON emits the wire form ["name", 1], the inverse of UnmarshalJSON.
func (a AccountAuthEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{a.Name, a.Weight})
}

// Authority models a Steem account authority (owner / active / posting).
// weight_threshold is the total weight required to authorize an action under
// this authority.
type Authority struct {
	WeightThreshold uint32             `json:"weight_threshold"`
	AccountAuths    []AccountAuthEntry `json:"account_auths"`
	KeyAuths        []KeyAuth          `json:"key_auths"`
}

// Manabar models a Steem resource manabar (steemd chain/util/manabar.hpp).
// current_mana is a share_type (int64) and serializes as a JSON string because
// it can be large; last_update_time is a unix-epoch uint32 emitted as a number.
type Manabar struct {
	CurrentMana    json.RawMessage `json:"current_mana"`
	LastUpdateTime uint32           `json:"last_update_time"`
}

// ExtendedAccount models a full condenser_api.get_accounts response entry,
// mirroring steemd's condenser_api::extended_account (which extends
// condenser_api::api_account_object). It carries every key the chain emits, so
// steemdb-sync's account refresher can capture a complete account snapshot in a
// single decode.
//
// The first eight fields (name, created, reputation, voting_power, balance,
// posting, active, owner) are the original conveyor-facing subset
// (conveyor/src/user-search/user.ts UserAccount) and MUST NOT be reordered;
// older callers depend on their order and types.
//
// Field conventions, matching the on-chain JSON exactly (verified against a live
// get_accounts(["steemit"]) response and steemd source):
//   - Assets ("100.000 STEEM") and ISO-8601 timestamps are kept as raw strings;
//     parsing is the consumer's responsibility.
//   - reputation and the share_type group below (withdrawn, to_withdraw,
//     curation_rewards, posting_rewards) are decoded as json.RawMessage because
//     the chain's share_type serializes inconsistently across nodes/versions as
//     either a JSON string (large number) or a JSON number — a typed field would
//     fail to decode on one or the other.
//   - proxied_vsf_votes is []json.RawMessage because the chain emits a MIXED
//     array, e.g. ["452574069424", 0, 0, 0] — some elements are strings, others
//     numbers — so a typed slice would fail.
//   - voting_manabar / downvote_manabar use the typed Manabar struct.
//   - The *_history / tags_usage / guest_bloggers fields are extended_account
//     collections that get_accounts currently always returns empty ([]). They
//     are captured as json.RawMessage so a future steemd that populates them
//     (FC serializes map<uint64_t,T> as [[key,value],...]) will not break the
//     decode; consumers that need their contents should decode defensively.
type ExtendedAccount struct {
	// --- conveyor-facing fields (original subset; do not reorder) ---
	Name        string          `json:"name"`
	Created     string          `json:"created"`
	Reputation  json.RawMessage `json:"reputation"`
	VotingPower int16           `json:"voting_power"`
	Balance     string          `json:"balance"`
	Posting     Authority       `json:"posting"`
	Active      Authority       `json:"active"`
	Owner       Authority       `json:"owner"`

	// --- scalar string fields: assets, timestamps, account names ---
	Proxy                         string `json:"proxy"`
	LastOwnerUpdate               string `json:"last_owner_update"`
	LastAccountUpdate             string `json:"last_account_update"`
	LastAccountRecovery           string `json:"last_account_recovery"`
	RecoveryAccount               string `json:"recovery_account"`
	ResetAccount                  string `json:"reset_account"`
	MemoKey                       string `json:"memo_key"`
	JSONMetadata                  string `json:"json_metadata"`
	PostingJSONMetadata           string `json:"posting_json_metadata"`
	SBDBalance                    string `json:"sbd_balance"`
	SavingsBalance                string `json:"savings_balance"`
	SavingsSBDBalance             string `json:"savings_sbd_balance"`
	SBDSeconds                    string `json:"sbd_seconds"`
	SBDSecondsLastUpdate          string `json:"sbd_seconds_last_update"`
	SBDLastInterestPayment        string `json:"sbd_last_interest_payment"`
	SavingsSBDSeconds             string `json:"savings_sbd_seconds"`
	SavingsSBDSecondsLastUpdate   string `json:"savings_sbd_seconds_last_update"`
	SavingsSBDLastInterestPayment string `json:"savings_sbd_last_interest_payment"`
	RewardSBDBalance              string `json:"reward_sbd_balance"`
	RewardSteemBalance            string `json:"reward_steem_balance"`
	RewardVestingBalance          string `json:"reward_vesting_balance"`
	RewardVestingSteem            string `json:"reward_vesting_steem"`
	VestingShares                 string `json:"vesting_shares"`
	DelegatedVestingShares        string `json:"delegated_vesting_shares"`
	ReceivedVestingShares         string `json:"received_vesting_shares"`
	VestingWithdrawRate           string `json:"vesting_withdraw_rate"`
	VestingBalance                string `json:"vesting_balance"`
	NextVestingWithdrawal         string `json:"next_vesting_withdrawal"`
	LastPost                      string `json:"last_post"`
	LastRootPost                  string `json:"last_root_post"`
	LastVoteTime                  string `json:"last_vote_time"`

	// --- scalar bool/int fields ---
	ID                      int      `json:"id"`
	Mined                   bool     `json:"mined"`
	CommentCount            int      `json:"comment_count"`
	LifetimeVoteCount       int      `json:"lifetime_vote_count"`
	PostCount               int      `json:"post_count"`
	CanVote                 bool     `json:"can_vote"`
	SavingsWithdrawRequests int      `json:"savings_withdraw_requests"`
	WithdrawRoutes          int      `json:"withdraw_routes"`
	WitnessesVotedFor       int      `json:"witnesses_voted_for"`
	PendingClaimedAccounts  int      `json:"pending_claimed_accounts"`
	WitnessVotes            []string `json:"witness_votes"`

	// --- manabar (typed) ---
	VotingManabar   Manabar `json:"voting_manabar"`
	DownvoteManabar Manabar `json:"downvote_manabar"`

	// --- RawMessage: share_type, may be string or number across nodes ---
	Withdrawn       json.RawMessage   `json:"withdrawn"`
	ToWithdraw      json.RawMessage   `json:"to_withdraw"`
	CurationRewards json.RawMessage   `json:"curation_rewards"`
	PostingRewards  json.RawMessage   `json:"posting_rewards"`
	PostBandwidth   json.RawMessage   `json:"post_bandwidth"`
	ProxiedVSFVotes []json.RawMessage `json:"proxied_vsf_votes"`

	// --- RawMessage: extended_account collections (currently empty [] on chain) ---
	TransferHistory json.RawMessage `json:"transfer_history"`
	MarketHistory   json.RawMessage `json:"market_history"`
	PostHistory     json.RawMessage `json:"post_history"`
	VoteHistory     json.RawMessage `json:"vote_history"`
	OtherHistory    json.RawMessage `json:"other_history"`
	TagsUsage       json.RawMessage `json:"tags_usage"`
	GuestBloggers   json.RawMessage `json:"guest_bloggers"`
}
