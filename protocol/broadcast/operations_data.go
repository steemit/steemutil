package broadcast

// OperationsData contains all Steem broadcast operation metadata.
var OperationsData = []BroadcastOperation{
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "vote",
		Params: []string{"voter", "author", "permlink", "weight"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "comment",
		Params: []string{"parent_author", "parent_permlink", "author", "permlink", "title", "body", "json_metadata"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "transfer",
		Params: []string{"from", "to", "amount", "memo"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "transfer_to_vesting",
		Params: []string{"from", "to", "amount"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "withdraw_vesting",
		Params: []string{"account", "vesting_shares"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "limit_order_create",
		Params: []string{"owner", "orderid", "amount_to_sell", "min_to_receive", "fill_or_kill", "expiration"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "limit_order_cancel",
		Params: []string{"owner", "orderid"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "price",
		Params: []string{"base", "quote"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "feed_publish",
		Params: []string{"publisher", "exchange_rate"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "convert",
		Params: []string{"owner", "requestid", "amount"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "account_create",
		Params: []string{"fee", "creator", "new_account_name", "owner", "active", "posting", "memo_key", "json_metadata"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "account_update",
		Params: []string{"account", "owner", "active", "posting", "memo_key", "json_metadata"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "witness_update",
		Params: []string{"owner", "url", "block_signing_key", "props", "fee"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "account_witness_vote",
		Params: []string{"account", "witness", "approve"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "account_witness_proxy",
		Params: []string{"account", "proxy"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "pow",
		Params: []string{"worker", "input", "signature", "work"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "custom",
		Params: []string{"required_auths", "id", "data"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "delete_comment",
		Params: []string{"author", "permlink"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "custom_json",
		Params: []string{"required_auths", "required_posting_auths", "id", "json"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "comment_options",
		Params: []string{"author", "permlink", "max_accepted_payout", "percent_steem_dollars", "allow_votes", "allow_curation_rewards", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "set_withdraw_vesting_route",
		Params: []string{"from_account", "to_account", "percent", "auto_vest"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "limit_order_create2",
		Params: []string{"owner", "orderid", "amount_to_sell", "exchange_rate", "fill_or_kill", "expiration"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "claim_account",
		Params: []string{"creator", "fee", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "create_claimed_account",
		Params: []string{"creator", "new_account_name", "owner", "active", "posting", "memo_key", "json_metadata", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "request_account_recovery",
		Params: []string{"recovery_account", "account_to_recover", "new_owner_authority", "extensions"},
	},
	{
		Roles: []string{"owner"},
		Operation: "recover_account",
		Params: []string{"account_to_recover", "new_owner_authority", "recent_owner_authority", "extensions"},
	},
	{
		Roles: []string{"owner"},
		Operation: "change_recovery_account",
		Params: []string{"account_to_recover", "new_recovery_account", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "escrow_transfer",
		Params: []string{"from", "to", "agent", "escrow_id", "sbd_amount", "steem_amount", "fee", "ratification_deadline", "escrow_expiration", "json_meta"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "escrow_dispute",
		Params: []string{"from", "to", "agent", "who", "escrow_id"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "escrow_release",
		Params: []string{"from", "to", "agent", "who", "receiver", "escrow_id", "sbd_amount", "steem_amount"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "pow2",
		Params: []string{"input", "pow_summary"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "escrow_approve",
		Params: []string{"from", "to", "agent", "who", "escrow_id", "approve"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "transfer_to_savings",
		Params: []string{"from", "to", "amount", "memo"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "transfer_from_savings",
		Params: []string{"from", "request_id", "to", "amount", "memo"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "cancel_transfer_from_savings",
		Params: []string{"from", "request_id"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "custom_binary",
		Params: []string{"id", "data"},
	},
	{
		Roles: []string{"owner"},
		Operation: "decline_voting_rights",
		Params: []string{"account", "decline"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "reset_account",
		Params: []string{"reset_account", "account_to_reset", "new_owner_authority"},
	},
	{
		Roles: []string{"owner", "posting"},
		Operation: "set_reset_account",
		Params: []string{"account", "current_reset_account", "reset_account"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "claim_reward_balance",
		Params: []string{"account", "reward_steem", "reward_sbd", "reward_vests"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "delegate_vesting_shares",
		Params: []string{"delegator", "delegatee", "vesting_shares"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "account_create_with_delegation",
		Params: []string{"fee", "delegation", "creator", "new_account_name", "owner", "active", "posting", "memo_key", "json_metadata", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "witness_set_properties",
		Params: []string{"owner", "props", "extensions"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "account_update2",
		Params: []string{"account", "owner", "active", "posting", "memo_key", "json_metadata", "posting_json_metadata", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "create_proposal",
		Params: []string{"creator", "receiver", "start_date", "end_date", "daily_pay", "subject", "permlink", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "update_proposal_votes",
		Params: []string{"voter", "proposal_ids", "approve", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "remove_proposal",
		Params: []string{"proposal_owner", "proposal_ids", "extensions"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "claim_reward_balance2",
		Params: []string{"account", "reward_tokens", "extensions"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "vote2",
		Params: []string{"voter", "author", "permlink", "rshares", "extensions"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "fill_convert_request",
		Params: []string{"owner", "requestid", "amount_in", "amount_out"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "comment_reward",
		Params: []string{"author", "permlink", "payout"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "liquidity_reward",
		Params: []string{"owner", "payout"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "interest",
		Params: []string{"owner", "interest"},
	},
	{
		Roles: []string{"active", "owner"},
		Operation: "fill_vesting_withdraw",
		Params: []string{"from_account", "to_account", "withdrawn", "deposited"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "fill_order",
		Params: []string{"current_owner", "current_orderid", "current_pays", "open_owner", "open_orderid", "open_pays"},
	},
	{
		Roles: []string{"posting", "active", "owner"},
		Operation: "fill_transfer_from_savings",
		Params: []string{"from", "to", "amount", "request_id", "memo"},
	},
}
