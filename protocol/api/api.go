package api

import (
	"github.com/steemit/steemutil/operation"
	"github.com/steemit/steemutil/util"
)

type DynamicGlobalProperties struct {
	HeadBlockNumber                 util.UInt32 `json:"head_block_number"`
	HeadBlockId                     string      `json:"head_block_id"`
	Time                            string      `json:"time"`
	CurrentWitness                  string      `json:"current_witness"`
	TotalPow                        util.UInt   `json:"total_pow"`
	NumPowWitnesses                 util.UInt   `json:"num_pow_witnesses"`
	VirtualSupply                   string      `json:"virtual_supply"`
	CurrentSupply                   string      `json:"current_supply"`
	ConfidentialSupply              string      `json:"confidential_supply"`
	InitSbdSupply                   string      `json:"init_sbd_supply"`
	CurrentSbdSupply                string      `json:"current_sbd_supply"`
	ConfidentialSbdSupply           string      `json:"confidential_sbd_supply"`
	TotalVestingFundSteem           string      `json:"total_vesting_fund_steem"`
	TotalVestingShares              string      `json:"total_vesting_shares"`
	TotalRewardFundSteem            string      `json:"total_reward_fund_steem"`
	TotalRewardShares2              string      `json:"total_reward_shares2"`
	PendingRewardedVestingShares    string      `json:"pending_rewarded_vesting_shares"`
	PendingRewardedVestingSteem     string      `json:"pending_rewarded_vesting_steem"`
	SbdInterestRate                 util.UInt   `json:"sbd_interest_rate"`
	SbdPrintRate                    util.UInt   `json:"sbd_print_rate"`
	MaximumBlockSize                util.UInt   `json:"maximum_block_size"`
	RequiredActionsPartitionPercent util.UInt   `json:"required_actions_partition_percent"`
	CurrentAslot                    util.UInt   `json:"current_aslot"`
	RecentSlotsFilled               string      `json:"recent_slots_filled"`
	ParticipationCount              util.UInt   `json:"participation_count"`
	LastIrreversibleBlockNum        util.UInt   `json:"last_irreversible_block_num"`
	VotePowerReserveRate            util.UInt   `json:"vote_power_reserve_rate"`
	DelegationReturnPeriod          util.UInt   `json:"delegation_return_period"`
	ReverseAuctionSeconds           util.UInt   `json:"reverse_auction_seconds"`
	AvailableAccountSubsidies       util.UInt   `json:"available_account_subsidies"`
	SbdStopPercent                  util.UInt   `json:"sbd_stop_percent"`
	SbdStartPercent                 util.UInt   `json:"sbd_start_percent"`
	NextMaintenanceTime             string      `json:"next_maintenance_time"`
	LastBudgetTime                  string      `json:"last_budget_time"`
	ContentRewardPercent            util.UInt   `json:"content_reward_percent"`
	VestingRewardPercent            util.UInt   `json:"vesting_reward_percent"`
	SpsFundPercent                  util.UInt   `json:"sps_fund_percent"`
	SpsIntervalLedger               string      `json:"sps_interval_ledger"`
	DownvotePoolPercent             util.UInt   `json:"downvote_pool_percent"`
}

type Block struct {
	BlockId               string        `json:"block_id"`
	Previous              string        `json:"previous"`
	Witness               string        `json:"witness"`
	WitnessSignature      string        `json:"witness_signature"`
	TransactionMerkleRoot string        `json:"transaction_merkle_root"`
	Transactions          []Transaction `json:"transactions"`
	Timestamp             *util.Time    `json:"timestamp"`
	Extensions            []any         `json:"extensions"`
	SigningKey            string        `json:"signing_key"`
	TransactionIds        []string      `json:"transaction_ids"`
}

type Transaction struct {
	RefBlockNum    util.UInt16          `json:"ref_block_num"`
	RefBlockPrefix util.UInt32          `json:"ref_block_prefix"`
	Expiration     *util.Time           `json:"expiration"`
	Operations     operation.Operations `json:"operations"`
	Signatures     []string             `json:"signatures"`
	Extensions     []any                `json:"extensions"`
	TransactionId  string               `json:"transaction_id"`
	BlockNum       util.UInt32          `json:"block_num"`
	TransactionNum util.UInt            `json:"transaction_num"`
}
