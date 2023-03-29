package api

type DynamicGlobalProperties struct {
	HeadBlockNumber                 uint   `json:"head_block_number"`
	HeadBlockId                     string `json:"head_block_id"`
	Time                            string `json:"time"`
	CurrentWitness                  string `json:"current_witness"`
	TotalPow                        uint   `json:"total_pow"`
	NumPowWitnesses                 uint   `json:"num_pow_witnesses"`
	VirtualSupply                   string `json:"virtual_supply"`
	CurrentSupply                   string `json:"current_supply"`
	ConfidentialSupply              string `json:"confidential_supply"`
	InitSbdSupply                   string `json:"init_sbd_supply"`
	CurrentSbdSupply                string `json:"current_sbd_supply"`
	ConfidentialSbdSupply           string `json:"confidential_sbd_supply"`
	TotalVestingFundSteem           string `json:"total_vesting_fund_steem"`
	TotalVestingShares              string `json:"total_vesting_shares"`
	TotalRewardFundSteem            string `json:"total_reward_fund_steem"`
	TotalRewardShares2              string `json:"total_reward_shares2"`
	PendingRewardedVestingShares    string `json:"pending_rewarded_vesting_shares"`
	PendingRewardedVestingSteem     string `json:"pending_rewarded_vesting_steem"`
	SbdInterestRate                 uint   `json:"sbd_interest_rate"`
	SbdPrintRate                    uint   `json:"sbd_print_rate"`
	MaximumBlockSize                uint   `json:"maximum_block_size"`
	RequiredActionsPartitionPercent uint   `json:"required_actions_partition_percent"`
	CurrentAslot                    uint   `json:"current_aslot"`
	RecentSlotsFilled               string `json:"recent_slots_filled"`
	ParticipationCount              uint   `json:"participation_count"`
	LastIrreversibleBlockNum        uint   `json:"last_irreversible_block_num"`
	VotePowerReserveRate            uint   `json:"vote_power_reserve_rate"`
	DelegationReturnPeriod          uint   `json:"delegation_return_period"`
	ReverseAuctionSeconds           uint   `json:"reverse_auction_seconds"`
	AvailableAccountSubsidies       uint   `json:"available_account_subsidies"`
	SbdStopPercent                  uint   `json:"sbd_stop_percent"`
	SbdStartPercent                 uint   `json:"sbd_start_percent"`
	NextMaintenanceTime             string `json:"next_maintenance_time"`
	LastBudgetTime                  string `json:"last_budget_time"`
	ContentRewardPercent            uint   `json:"content_reward_percent"`
	VestingRewardPercent            uint   `json:"vesting_reward_percent"`
	SpsFundPercent                  uint   `json:"sps_fund_percent"`
	SpsIntervalLedger               string `json:"sps_interval_ledger"`
	DownvotePoolPercent             uint   `json:"downvote_pool_percent"`
}
