package api

// MethodsData contains all Steem API method definitions.
var MethodsData = []APIMethod{
	{
		API: "database_api",
		Method: "set_subscribe_callback",
		Params: []string{"callback", "clearFilter"},
	},
	{
		API: "database_api",
		Method: "set_pending_transaction_callback",
		Params: []string{"cb"},
	},
	{
		API: "database_api",
		Method: "set_block_applied_callback",
		Params: []string{"cb"},
	},
	{
		API: "database_api",
		Method: "cancel_all_subscriptions",
	},
	{
		API: "database_api",
		Method: "get_trending_tags",
		Params: []string{"afterTag", "limit"},
	},
	{
		API: "database_api",
		Method: "get_tags_used_by_author",
		Params: []string{"author"},
	},
	{
		API: "database_api",
		Method: "get_post_discussions_by_payout",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_comment_discussions_by_payout",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_trending",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_trending30",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_created",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_active",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_cashout",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_payout",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_votes",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_children",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_hot",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_feed",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_blog",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_comments",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_promoted",
		Params: []string{"query"},
	},
	{
		API: "database_api",
		Method: "get_block_header",
		Params: []string{"blockNum"},
	},
	{
		API: "database_api",
		Method: "get_block",
		Params: []string{"blockNum"},
	},
	{
		API: "database_api",
		Method: "get_ops_in_block",
		Params: []string{"blockNum", "onlyVirtual"},
	},
	{
		API: "database_api",
		Method: "get_state",
		Params: []string{"path"},
	},
	{
		API: "database_api",
		Method: "get_trending_categories",
		Params: []string{"after", "limit"},
	},
	{
		API: "database_api",
		Method: "get_best_categories",
		Params: []string{"after", "limit"},
	},
	{
		API: "database_api",
		Method: "get_active_categories",
		Params: []string{"after", "limit"},
	},
	{
		API: "database_api",
		Method: "get_recent_categories",
		Params: []string{"after", "limit"},
	},
	{
		API: "database_api",
		Method: "get_config",
	},
	{
		API: "database_api",
		Method: "get_dynamic_global_properties",
	},
	{
		API: "database_api",
		Method: "get_chain_properties",
	},
	{
		API: "database_api",
		Method: "get_feed_history",
	},
	{
		API: "database_api",
		Method: "get_current_median_history_price",
	},
	{
		API: "database_api",
		Method: "get_witness_schedule",
	},
	{
		API: "database_api",
		Method: "get_hardfork_version",
	},
	{
		API: "database_api",
		Method: "get_next_scheduled_hardfork",
	},
	{
		API: "account_by_key_api",
		Method: "get_key_references",
		Params: []string{"key"},
	},
	{
		API: "database_api",
		Method: "get_accounts",
		Params: []string{"names"},
	},
	{
		API: "database_api",
		Method: "get_account_references",
		Params: []string{"accountId"},
	},
	{
		API: "database_api",
		Method: "lookup_account_names",
		Params: []string{"accountNames"},
	},
	{
		API: "database_api",
		Method: "lookup_accounts",
		Params: []string{"lowerBoundName", "limit"},
	},
	{
		API: "database_api",
		Method: "get_account_count",
	},
	{
		API: "database_api",
		Method: "get_conversion_requests",
		Params: []string{"accountName"},
	},
	{
		API: "database_api",
		Method: "get_account_history",
		Params: []string{"account", "from", "limit"},
	},
	{
		API: "database_api",
		Method: "get_owner_history",
		Params: []string{"account"},
	},
	{
		API: "database_api",
		Method: "get_recovery_request",
		Params: []string{"account"},
	},
	{
		API: "database_api",
		Method: "get_escrow",
		Params: []string{"from", "escrowId"},
	},
	{
		API: "database_api",
		Method: "get_withdraw_routes",
		Params: []string{"account", "withdrawRouteType"},
	},
	{
		API: "database_api",
		Method: "get_account_bandwidth",
		Params: []string{"account", "bandwidthType"},
	},
	{
		API: "database_api",
		Method: "get_savings_withdraw_from",
		Params: []string{"account"},
	},
	{
		API: "database_api",
		Method: "get_savings_withdraw_to",
		Params: []string{"account"},
	},
	{
		API: "database_api",
		Method: "get_order_book",
		Params: []string{"limit"},
	},
	{
		API: "database_api",
		Method: "get_open_orders",
		Params: []string{"owner"},
	},
	{
		API: "database_api",
		Method: "get_liquidity_queue",
		Params: []string{"startAccount", "limit"},
	},
	{
		API: "database_api",
		Method: "get_transaction_hex",
		Params: []string{"trx"},
	},
	{
		API: "database_api",
		Method: "get_transaction",
		Params: []string{"trxId"},
	},
	{
		API: "database_api",
		Method: "get_required_signatures",
		Params: []string{"trx", "availableKeys"},
	},
	{
		API: "database_api",
		Method: "get_potential_signatures",
		Params: []string{"trx"},
	},
	{
		API: "database_api",
		Method: "verify_authority",
		Params: []string{"trx"},
	},
	{
		API: "database_api",
		Method: "verify_account_authority",
		Params: []string{"nameOrId", "signers"},
	},
	{
		API: "database_api",
		Method: "get_active_votes",
		Params: []string{"author", "permlink"},
	},
	{
		API: "database_api",
		Method: "get_account_votes",
		Params: []string{"voter"},
	},
	{
		API: "database_api",
		Method: "get_content",
		Params: []string{"author", "permlink"},
	},
	{
		API: "database_api",
		Method: "get_content_replies",
		Params: []string{"author", "permlink"},
	},
	{
		API: "database_api",
		Method: "get_discussions_by_author_before_date",
		Params: []string{"author", "startPermlink", "beforeDate", "limit"},
	},
	{
		API: "database_api",
		Method: "get_replies_by_last_update",
		Params: []string{"startAuthor", "startPermlink", "limit"},
	},
	{
		API: "database_api",
		Method: "get_witnesses",
		Params: []string{"witnessIds"},
	},
	{
		API: "database_api",
		Method: "get_witness_by_account",
		Params: []string{"accountName"},
	},
	{
		API: "database_api",
		Method: "get_witnesses_by_vote",
		Params: []string{"from", "limit"},
	},
	{
		API: "database_api",
		Method: "lookup_witness_accounts",
		Params: []string{"lowerBoundName", "limit"},
	},
	{
		API: "database_api",
		Method: "get_witness_count",
	},
	{
		API: "database_api",
		Method: "get_active_witnesses",
	},
	{
		API: "database_api",
		Method: "get_miner_queue",
	},
	{
		API: "database_api",
		Method: "get_reward_fund",
		Params: []string{"name"},
	},
	{
		API: "database_api",
		Method: "get_vesting_delegations",
		Params: []string{"account", "from", "limit"},
	},
	{
		API: "login_api",
		Method: "login",
		Params: []string{"username", "password"},
	},
	{
		API: "login_api",
		Method: "get_api_by_name",
		Params: []string{"database_api"},
	},
	{
		API: "login_api",
		Method: "get_version",
	},
	{
		API: "follow_api",
		Method: "get_followers",
		Params: []string{"following", "startFollower", "followType", "limit"},
	},
	{
		API: "follow_api",
		Method: "get_following",
		Params: []string{"follower", "startFollowing", "followType", "limit"},
	},
	{
		API: "follow_api",
		Method: "get_follow_count",
		Params: []string{"account"},
	},
	{
		API: "follow_api",
		Method: "get_feed_entries",
		Params: []string{"account", "entryId", "limit"},
	},
	{
		API: "follow_api",
		Method: "get_feed",
		Params: []string{"account", "entryId", "limit"},
	},
	{
		API: "follow_api",
		Method: "get_blog_entries",
		Params: []string{"account", "entryId", "limit"},
	},
	{
		API: "follow_api",
		Method: "get_blog",
		Params: []string{"account", "entryId", "limit"},
	},
	{
		API: "follow_api",
		Method: "get_account_reputations",
		Params: []string{"lowerBoundName", "limit"},
	},
	{
		API: "follow_api",
		Method: "get_reblogged_by",
		Params: []string{"author", "permlink"},
	},
	{
		API: "follow_api",
		Method: "get_blog_authors",
		Params: []string{"blogAccount"},
	},
	{
		API: "network_broadcast_api",
		Method: "broadcast_transaction",
		Params: []string{"trx"},
	},
	{
		API: "network_broadcast_api",
		Method: "broadcast_transaction_with_callback",
		Params: []string{"confirmationCallback", "trx"},
	},
	{
		API: "network_broadcast_api",
		Method: "broadcast_transaction_synchronous",
		Params: []string{"trx"},
	},
	{
		API: "network_broadcast_api",
		Method: "broadcast_block",
		Params: []string{"b"},
	},
	{
		API: "network_broadcast_api",
		Method: "set_max_block_age",
		Params: []string{"maxBlockAge"},
	},
	{
		API: "market_history_api",
		Method: "get_ticker",
	},
	{
		API: "market_history_api",
		Method: "get_volume",
	},
	{
		API: "market_history_api",
		Method: "get_order_book",
		MethodName: "getMarketOrderBook",
		Params: []string{"limit"},
	},
	{
		API: "market_history_api",
		Method: "get_trade_history",
		Params: []string{"start", "end", "limit"},
	},
	{
		API: "market_history_api",
		Method: "get_recent_trades",
		Params: []string{"limit"},
	},
	{
		API: "market_history_api",
		Method: "get_market_history",
		Params: []string{"bucket_seconds", "start", "end"},
	},
	{
		API: "market_history_api",
		Method: "get_market_history_buckets",
	},
	{
		API: "condenser_api",
		Method: "find_proposals",
		Params: []string{"id_set"},
	},
	{
		API: "condenser_api",
		Method: "list_proposals",
		Params: []string{"start", "limit", "order_by", "order_direction", "status"},
	},
	{
		API: "condenser_api",
		Method: "list_proposal_votes",
		Params: []string{"start", "limit", "order_by", "order_direction", "status"},
	},
	{
		API: "condenser_api",
		Method: "get_nai_pool",
	},
	{
		API: "rc_api",
		Method: "find_rc_accounts",
		Params: []string{"accounts"},
		IsObject: true,
	},
	{
		API: "condenser_api",
		Method: "get_expiring_vesting_delegations",
		Params: []string{"account", "start", "limit"},
	},
	{
		API: "database_api",
		Method: "find_change_recovery_account_requests",
		Params: []string{"account"},
		IsObject: true,
	},
}
