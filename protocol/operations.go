package protocol

import (
	"encoding/json"

	"github.com/steemit/steemutil/encoder"
)

// FC_REFLECT( steemit::chain::report_over_production_operation,
//             (reporter)
//             (first_block)
//             (second_block) )

type ReportOverProductionOperation struct {
	Reporter    string `json:"reporter"`
	FirstBlock  string `json:"first_block"`
	SecondBlock string `json:"second_block"`
}

func (op *ReportOverProductionOperation) Type() OpType {
	return TypeReportOverProduction
}

func (op *ReportOverProductionOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::convert_operation,
//             (owner)
//             (requestid)
//             (amount) )

type ConvertOperation struct {
	Owner     string `json:"owner"`
	RequestID uint32 `json:"requestid"`
	Amount    string `json:"amount"`
}

func (op *ConvertOperation) Type() OpType {
	return TypeConvert
}

func (op *ConvertOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::feed_publish_operation,
//             (publisher)
//             (exchange_rate) )

type FeedPublishOperation struct {
	Publisher    string `json:"publisher"`
	ExchangeRate struct {
		Base  string `json:"base"`
		Quote string `json:"quote"`
	} `json:"exchange_rate"`
}

func (op *FeedPublishOperation) Type() OpType {
	return TypeFeedPublish
}

func (op *FeedPublishOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::pow,
//             (worker)
//             (input)
//             (signature)
//             (work) )

type POW struct {
	Worker    string `json:"worker"`
	Input     string `json:"input"`
	Signature string `json:"signature"`
	Work      string `json:"work"`
}

// FC_REFLECT( steemit::chain::chain_properties,
//             (account_creation_fee)
//             (maximum_block_size)
//             (sbd_interest_rate) );

// AssetObject represents an asset in the new format (with nai field)
// Used for account_creation_fee in newer Steem versions
type AssetObject struct {
	Amount    string `json:"amount"`
	NAI       string `json:"nai"` // Native Asset Identifier
	Precision uint8  `json:"precision"`
}

type ChainProperties struct {
	AccountCreationFee string `json:"account_creation_fee"`
	MaximumBlockSize   uint32 `json:"maximum_block_size"`
	SBDInterestRate    uint16 `json:"sbd_interest_rate"`
}

// UnmarshalJSON custom unmarshaling for ChainProperties to handle
// account_creation_fee as either a string (old format) or an object (new format with nai)
func (cp *ChainProperties) UnmarshalJSON(data []byte) error {
	// First, try to unmarshal as a map to check the account_creation_fee type
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Handle account_creation_fee - can be string or object
	if feeRaw, ok := raw["account_creation_fee"]; ok {
		switch v := feeRaw.(type) {
		case string:
			// Old format: string like "0.100 STEEM"
			cp.AccountCreationFee = v
		case map[string]interface{}:
			// New format: object with amount, nai, precision
			// Convert to string format for backward compatibility
			if amount, ok := v["amount"].(string); ok {
				if nai, ok := v["nai"].(string); ok {
					if precision, ok := v["precision"].(float64); ok {
						// Format as "amount symbol" where symbol is derived from nai
						// For now, we'll store the amount with precision
						// The nai can be used to determine the symbol later if needed
						cp.AccountCreationFee = formatAssetFromObject(amount, nai, uint8(precision))
					} else {
						cp.AccountCreationFee = amount
					}
				} else {
					cp.AccountCreationFee = amount
				}
			} else {
				// Fallback: try to marshal back to JSON string
				if feeJSON, err := json.Marshal(v); err == nil {
					cp.AccountCreationFee = string(feeJSON)
				}
			}
		default:
			// Unknown type, try to convert to string
			if feeJSON, err := json.Marshal(v); err == nil {
				cp.AccountCreationFee = string(feeJSON)
			}
		}
	}

	// Handle maximum_block_size
	if mbs, ok := raw["maximum_block_size"]; ok {
		switch v := mbs.(type) {
		case float64:
			cp.MaximumBlockSize = uint32(v)
		case uint32:
			cp.MaximumBlockSize = v
		}
	}

	// Handle sbd_interest_rate
	if sir, ok := raw["sbd_interest_rate"]; ok {
		switch v := sir.(type) {
		case float64:
			cp.SBDInterestRate = uint16(v)
		case uint16:
			cp.SBDInterestRate = v
		}
	}

	return nil
}

// formatAssetFromObject formats an asset object into a string representation
// This is a helper function to convert the new format to the old string format
func formatAssetFromObject(amount, nai string, precision uint8) string {
	// Map common NAIs to symbols
	symbolMap := map[string]string{
		"@@000000021": "STEEM",
		"@@000000013": "SBD",
		"@@000000037": "VESTS",
	}

	symbol := symbolMap[nai]
	if symbol == "" {
		// Unknown NAI, use the NAI itself
		symbol = nai
	}

	// Convert amount string to decimal format
	// amount is already a string representation of the integer amount
	// We need to format it with the precision
	if precision == 0 {
		return amount + " " + symbol
	}

	// Format with decimal point
	amountInt := amount
	if len(amountInt) <= int(precision) {
		// Pad with zeros
		for len(amountInt) < int(precision) {
			amountInt = "0" + amountInt
		}
		return "0." + amountInt + " " + symbol
	}

	// Insert decimal point
	dotPos := len(amountInt) - int(precision)
	return amountInt[:dotPos] + "." + amountInt[dotPos:] + " " + symbol
}

// FC_REFLECT( steemit::chain::pow_operation,
//             (worker_account)
//             (block_id)
//             (nonce)
//             (work)
//             (props) )

type POWOperation struct {
	WorkerAccount string           `json:"worker_account"`
	BlockID       string           `json:"block_id"`
	Nonce         *Int             `json:"nonce"`
	Work          *POW             `json:"work"`
	Props         *ChainProperties `json:"props"`
}

func (op *POWOperation) Type() OpType {
	return TypePOW
}

func (op *POWOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::account_create_operation,
//             (fee)
//             (creator)
//             (new_account_name)
//             (owner)
//             (active)
//             (posting)
//             (memo_key)
//             (json_metadata) )

type AccountCreateOperation struct {
	Fee            string     `json:"fee"`
	Creator        string     `json:"creator"`
	NewAccountName string     `json:"new_account_name"`
	Owner          *Authority `json:"owner"`
	Active         *Authority `json:"active"`
	Posting        *Authority `json:"posting"`
	MemoKey        string     `json:"memo_key"`
	JsonMetadata   string     `json:"json_metadata"`
}

func (op *AccountCreateOperation) Type() OpType {
	return TypeAccountCreate
}

func (op *AccountCreateOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::account_update_operation,
//             (account)
//             (owner)
//             (active)
//             (posting)
//             (memo_key)
//             (json_metadata) )

type AccountUpdateOperation struct {
	Account      string     `json:"account"`
	Owner        *Authority `json:"owner"`
	Active       *Authority `json:"active"`
	Posting      *Authority `json:"posting"`
	MemoKey      string     `json:"memo_key"`
	JsonMetadata string     `json:"json_metadata"`
}

func (op *AccountUpdateOperation) Type() OpType {
	return TypeAccountUpdate
}

func (op *AccountUpdateOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::transfer_operation,
//             (from)
//             (to)
//             (amount)
//             (memo) )

type TransferOperation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Memo   string `json:"memo"`
}

func (op *TransferOperation) Type() OpType {
	return TypeTransfer
}

func (op *TransferOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::transfer_to_vesting_operation,
//             (from)
//             (to)
//             (amount) )

type TransferToVestingOperation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

func (op *TransferToVestingOperation) Type() OpType {
	return TypeTransferToVesting
}

func (op *TransferToVestingOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::withdraw_vesting_operation,
//             (account)
//             (vesting_shares) )

type WithdrawVestingOperation struct {
	Account       string `json:"account"`
	VestingShares string `json:"vesting_shares"`
}

func (op *WithdrawVestingOperation) Type() OpType {
	return TypeWithdrawVesting
}

func (op *WithdrawVestingOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::set_withdraw_vesting_route_operation,
//             (from_account)
//             (to_account)
//             (percent)
//             (auto_vest) )

// FC_REFLECT( steemit::chain::witness_update_operation,
//             (owner)
//             (url)
//             (block_signing_key)
//             (props)
//             (fee) )

// FC_REFLECT( steemit::chain::account_witness_vote_operation,
//             (account)
//             (witness)(approve) )

type AccountWitnessVoteOperation struct {
	Account string `json:"account"`
	Witness string `json:"witness"`
	Approve bool   `json:"approve"`
}

func (op *AccountWitnessVoteOperation) Type() OpType {
	return TypeAccountWitnessVote
}

func (op *AccountWitnessVoteOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::account_witness_proxy_operation,
//             (account)
//             (proxy) )

type AccountWitnessProxyOperation struct {
	Account string `json:"account"`
	Proxy   string `json:"proxy"`
}

func (op *AccountWitnessProxyOperation) Type() OpType {
	return TypeAccountWitnessProxy
}

func (op *AccountWitnessProxyOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::comment_operation,
//             (parent_author)
//             (parent_permlink)
//             (author)
//             (permlink)
//             (title)
//             (body)
//             (json_metadata) )

// CommentOperation represents either a new post or a comment.
//
// In case Title is filled in and ParentAuthor is empty, it is a new post.
// The post category can be read from ParentPermlink.
type CommentOperation struct {
	ParentAuthor   string `json:"parent_author"`
	ParentPermlink string `json:"parent_permlink"`
	Author         string `json:"author"`
	Permlink       string `json:"permlink"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	JsonMetadata   string `json:"json_metadata"`
}

func (op *CommentOperation) Type() OpType {
	return TypeComment
}

func (op *CommentOperation) Data() any {
	return op
}

func (op *CommentOperation) IsStoryOperation() bool {
	return op.ParentAuthor == ""
}

func (op *CommentOperation) MarshalTransaction(encoderObj *encoder.Encoder) error {
	enc := encoder.NewRollingEncoder(encoderObj)
	enc.EncodeUVarint(uint64(op.Type().Code()))
	enc.Encode(op.ParentAuthor)
	enc.Encode(op.ParentPermlink)
	enc.Encode(op.Author)
	enc.Encode(op.Permlink)
	enc.Encode(op.Title)
	enc.Encode(op.Body)
	enc.Encode(op.JsonMetadata)
	return enc.Err()
}

// FC_REFLECT( steemit::chain::vote_operation,
//             (voter)
//             (author)
//             (permlink)
//             (weight) )

type VoteOperation struct {
	Voter    string `json:"voter"`
	Author   string `json:"author"`
	Permlink string `json:"permlink"`
	Weight   Int16  `json:"weight"`
}

func (op *VoteOperation) Type() OpType {
	return TypeVote
}

func (op *VoteOperation) Data() any {
	return op
}

func (op *VoteOperation) MarshalTransaction(encoderObj *encoder.Encoder) error {
	enc := encoder.NewRollingEncoder(encoderObj)
	enc.EncodeUVarint(uint64(TypeVote.Code()))
	enc.Encode(op.Voter)
	enc.Encode(op.Author)
	enc.Encode(op.Permlink)
	enc.Encode(op.Weight)
	return enc.Err()
}

// FC_REFLECT( steemit::chain::custom_operation,
//             (required_auths)
//             (id)
//             (data) )

// FC_REFLECT( steemit::chain::limit_order_create_operation,
//             (owner)
//             (orderid)
//             (amount_to_sell)
//             (min_to_receive)
//             (fill_or_kill)
//             (expiration) )

type LimitOrderCreateOperation struct {
	Owner        string `json:"owner"`
	OrderID      uint32 `json:"orderid"`
	AmountToSell string `json:"amount_to_sell"`
	MinToReceive string `json:"min_to_receive"`
	FillOrKill   bool   `json:"fill_or_kill"`
	Expiration   *Time  `json:"expiration"`
}

func (op *LimitOrderCreateOperation) Type() OpType {
	return TypeLimitOrderCreate
}

func (op *LimitOrderCreateOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::limit_order_cancel_operation,
//             (owner)
//             (orderid) )

type LimitOrderCancelOperation struct {
	Owner   string `json:"owner"`
	OrderID uint32 `json:"orderid"`
}

func (op *LimitOrderCancelOperation) Type() OpType {
	return TypeLimitOrderCancel
}

func (op *LimitOrderCancelOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::delete_comment_operation,
//             (author)
//             (permlink) )

type DeleteCommentOperation struct {
	Author   string `json:"author"`
	Permlink string `json:"permlink"`
}

func (op *DeleteCommentOperation) Type() OpType {
	return TypeDeleteComment
}

func (op *DeleteCommentOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::comment_options_operation,
//             (author)
//             (permlink)
//             (max_accepted_payout)
//             (percent_steem_dollars)
//             (allow_votes)
//             (allow_curation_rewards)
//             (extensions) )

type CommentOptionsOperation struct {
	Author               string `json:"author"`
	Permlink             string `json:"permlink"`
	MaxAcceptedPayout    string `json:"max_accepted_payout"`
	PercentSteemDollars  uint16 `json:"percent_steem_dollars"`
	AllowVotes           bool   `json:"allow_votes"`
	AllowCurationRewards bool   `json:"allow_curation_rewards"`
	Extensions           []any  `json:"extensions"`
}

func (op *CommentOptionsOperation) Type() OpType {
	return TypeCommentOptions
}

func (op *CommentOptionsOperation) Data() any {
	return op
}

type Authority struct {
	AccountAuths    StringInt64Map `json:"account_auths"`
	KeyAuths        StringInt64Map `json:"key_auths"`
	WeightThreshold uint32         `json:"weight_threshold"`
}

// FC_REFLECT( steemit::chain::witness_update_operation,
//             (owner)
//             (url)
//             (block_signing_key)
//             (props)
//             (fee) )

type WitnessUpdateOperation struct {
	Owner           string           `json:"owner"`
	URL             string           `json:"url"`
	BlockSigningKey string           `json:"block_signing_key"`
	Props           *ChainProperties `json:"props"`
	Fee             string           `json:"fee"`
}

func (op *WitnessUpdateOperation) Type() OpType {
	return TypeWitnessUpdate
}

func (op *WitnessUpdateOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::set_withdraw_vesting_route_operation,
//             (from_account)
//             (to_account)
//             (percent)
//             (auto_vest) )

type SetWithdrawVestingRouteOperation struct {
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Percent     uint16 `json:"percent"`
	AutoVest    bool   `json:"auto_vest"`
}

func (op *SetWithdrawVestingRouteOperation) Type() OpType {
	return TypeSetWithdrawVestingRoute
}

func (op *SetWithdrawVestingRouteOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::limit_order_create2_operation,
//             (owner)
//             (orderid)
//             (amount_to_sell)
//             (exchange_rate)
//             (fill_or_kill)
//             (expiration) )

type LimitOrderCreate2Operation struct {
	Owner        string `json:"owner"`
	OrderID      uint32 `json:"orderid"`
	AmountToSell string `json:"amount_to_sell"`
	ExchangeRate struct {
		Base  string `json:"base"`
		Quote string `json:"quote"`
	} `json:"exchange_rate"`
	FillOrKill bool  `json:"fill_or_kill"`
	Expiration *Time `json:"expiration"`
}

func (op *LimitOrderCreate2Operation) Type() OpType {
	return TypeLimitOrderCreate2
}

func (op *LimitOrderCreate2Operation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::claim_account_operation,
//             (creator)
//             (fee)
//             (extensions) )

type ClaimAccountOperation struct {
	Creator    string `json:"creator"`
	Fee        string `json:"fee"`
	Extensions []any  `json:"extensions"`
}

func (op *ClaimAccountOperation) Type() OpType {
	return TypeClaimAccount
}

func (op *ClaimAccountOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::create_claimed_account_operation,
//             (creator)
//             (new_account_name)
//             (owner)
//             (active)
//             (posting)
//             (memo_key)
//             (json_metadata)
//             (extensions) )

type CreateClaimedAccountOperation struct {
	Creator        string     `json:"creator"`
	NewAccountName string     `json:"new_account_name"`
	Owner          *Authority `json:"owner"`
	Active         *Authority `json:"active"`
	Posting        *Authority `json:"posting"`
	MemoKey        string     `json:"memo_key"`
	JsonMetadata   string     `json:"json_metadata"`
	Extensions     []any      `json:"extensions"`
}

func (op *CreateClaimedAccountOperation) Type() OpType {
	return TypeCreateClaimedAccount
}

func (op *CreateClaimedAccountOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::request_account_recovery_operation,
//             (recovery_account)
//             (account_to_recover)
//             (new_owner_authority)
//             (extensions) )

type RequestAccountRecoveryOperation struct {
	RecoveryAccount   string     `json:"recovery_account"`
	AccountToRecover  string     `json:"account_to_recover"`
	NewOwnerAuthority *Authority `json:"new_owner_authority"`
	Extensions        []any      `json:"extensions"`
}

func (op *RequestAccountRecoveryOperation) Type() OpType {
	return TypeRequestAccountRecovery
}

func (op *RequestAccountRecoveryOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::recover_account_operation,
//             (account_to_recover)
//             (new_owner_authority)
//             (recent_owner_authority)
//             (extensions) )

type RecoverAccountOperation struct {
	AccountToRecover     string     `json:"account_to_recover"`
	NewOwnerAuthority    *Authority `json:"new_owner_authority"`
	RecentOwnerAuthority *Authority `json:"recent_owner_authority"`
	Extensions           []any      `json:"extensions"`
}

func (op *RecoverAccountOperation) Type() OpType {
	return TypeRecoverAccount
}

func (op *RecoverAccountOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::change_recovery_account_operation,
//             (account_to_recover)
//             (new_recovery_account)
//             (extensions) )

type ChangeRecoveryAccountOperation struct {
	AccountToRecover   string `json:"account_to_recover"`
	NewRecoveryAccount string `json:"new_recovery_account"`
	Extensions         []any  `json:"extensions"`
}

func (op *ChangeRecoveryAccountOperation) Type() OpType {
	return TypeChangeRecoveryAccount
}

func (op *ChangeRecoveryAccountOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::escrow_transfer_operation,
//             (from)
//             (to)
//             (agent)
//             (escrow_id)
//             (sbd_amount)
//             (steem_amount)
//             (fee)
//             (ratification_deadline)
//             (escrow_expiration)
//             (json_meta) )

type EscrowTransferOperation struct {
	From                 string `json:"from"`
	To                   string `json:"to"`
	SBDAmount            string `json:"sbd_amount"`
	SteemAmount          string `json:"steem_amount"`
	EscrowID             uint32 `json:"escrow_id"`
	Agent                string `json:"agent"`
	Fee                  string `json:"fee"`
	JsonMeta             string `json:"json_meta"`
	RatificationDeadline *Time  `json:"ratification_deadline"`
	EscrowExpiration     *Time  `json:"escrow_expiration"`
}

func (op *EscrowTransferOperation) Type() OpType {
	return TypeEscrowTransfer
}

func (op *EscrowTransferOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::escrow_dispute_operation,
//             (from)
//             (to)
//             (agent)
//             (who)
//             (escrow_id) )

type EscrowDisputeOperation struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Agent    string `json:"agent"`
	Who      string `json:"who"`
	EscrowID uint32 `json:"escrow_id"`
}

func (op *EscrowDisputeOperation) Type() OpType {
	return TypeEscrowDispute
}

func (op *EscrowDisputeOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::escrow_release_operation,
//             (from)
//             (to)
//             (agent)
//             (who)
//             (receiver)
//             (escrow_id)
//             (sbd_amount)
//             (steem_amount) )

type EscrowReleaseOperation struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Agent       string `json:"agent"`
	Who         string `json:"who"`
	Receiver    string `json:"receiver"`
	EscrowID    uint32 `json:"escrow_id"`
	SBDAmount   string `json:"sbd_amount"`
	SteemAmount string `json:"steem_amount"`
}

func (op *EscrowReleaseOperation) Type() OpType {
	return TypeEscrowRelease
}

func (op *EscrowReleaseOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::escrow_approve_operation,
//             (from)
//             (to)
//             (agent)
//             (who)
//             (escrow_id)
//             (approve) )

type EscrowApproveOperation struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Agent    string `json:"agent"`
	Who      string `json:"who"`
	EscrowID uint32 `json:"escrow_id"`
	Approve  bool   `json:"approve"`
}

func (op *EscrowApproveOperation) Type() OpType {
	return TypeEscrowApprove
}

func (op *EscrowApproveOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::pow2_operation,
//             (input)
//             (pow_summary) )

type POW2Operation struct {
	Input      string `json:"input"`
	POWSummary any    `json:"pow_summary"`
}

func (op *POW2Operation) Type() OpType {
	return TypePOW2
}

func (op *POW2Operation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::transfer_to_savings_operation,
//             (from)
//             (to)
//             (amount)
//             (memo) )

type TransferToSavingsOperation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Memo   string `json:"memo"`
}

func (op *TransferToSavingsOperation) Type() OpType {
	return TypeTransferToSavings
}

func (op *TransferToSavingsOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::transfer_from_savings_operation,
//             (from)
//             (request_id)
//             (to)
//             (amount)
//             (memo) )

type TransferFromSavingsOperation struct {
	From      string `json:"from"`
	RequestID uint32 `json:"request_id"`
	To        string `json:"to"`
	Amount    string `json:"amount"`
	Memo      string `json:"memo"`
}

func (op *TransferFromSavingsOperation) Type() OpType {
	return TypeTransferFromSavings
}

func (op *TransferFromSavingsOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::cancel_transfer_from_savings_operation,
//             (from)
//             (request_id) )

type CancelTransferFromSavingsOperation struct {
	From      string `json:"from"`
	RequestID uint32 `json:"request_id"`
}

func (op *CancelTransferFromSavingsOperation) Type() OpType {
	return TypeCancelTransferFromSavings
}

func (op *CancelTransferFromSavingsOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::custom_binary_operation,
//             (id)
//             (data) )

type CustomBinaryOperation struct {
	ID        string `json:"id"`
	DataBytes string `json:"data"`
}

func (op *CustomBinaryOperation) Type() OpType {
	return TypeCustomBinary
}

func (op *CustomBinaryOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::decline_voting_rights_operation,
//             (account)
//             (decline) )

type DeclineVotingRightsOperation struct {
	Account string `json:"account"`
	Decline bool   `json:"decline"`
}

func (op *DeclineVotingRightsOperation) Type() OpType {
	return TypeDeclineVotingRights
}

func (op *DeclineVotingRightsOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::reset_account_operation,
//             (reset_account)
//             (account_to_reset)
//             (new_owner_authority) )

type ResetAccountOperation struct {
	ResetAccount      string     `json:"reset_account"`
	AccountToReset    string     `json:"account_to_reset"`
	NewOwnerAuthority *Authority `json:"new_owner_authority"`
}

func (op *ResetAccountOperation) Type() OpType {
	return TypeResetAccount
}

func (op *ResetAccountOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::set_reset_account_operation,
//             (account)
//             (current_reset_account)
//             (reset_account) )

type SetResetAccountOperation struct {
	Account             string `json:"account"`
	CurrentResetAccount string `json:"current_reset_account"`
	ResetAccount        string `json:"reset_account"`
}

func (op *SetResetAccountOperation) Type() OpType {
	return TypeSetResetAccount
}

func (op *SetResetAccountOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::claim_reward_balance_operation,
//             (account)
//             (reward_steem)
//             (reward_sbd)
//             (reward_vests) )

type ClaimRewardBalanceOperation struct {
	Account     string `json:"account"`
	RewardSteem string `json:"reward_steem"`
	RewardSBD   string `json:"reward_sbd"`
	RewardVests string `json:"reward_vests"`
}

func (op *ClaimRewardBalanceOperation) Type() OpType {
	return TypeClaimRewardBalance
}

func (op *ClaimRewardBalanceOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::delegate_vesting_shares_operation,
//             (delegator)
//             (delegatee)
//             (vesting_shares) )

type DelegateVestingSharesOperation struct {
	Delegator     string `json:"delegator"`
	Delegatee     string `json:"delegatee"`
	VestingShares string `json:"vesting_shares"`
}

func (op *DelegateVestingSharesOperation) Type() OpType {
	return TypeDelegateVestingShares
}

func (op *DelegateVestingSharesOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::account_create_with_delegation_operation,
//             (fee)
//             (delegation)
//             (creator)
//             (new_account_name)
//             (owner)
//             (active)
//             (posting)
//             (memo_key)
//             (json_metadata)
//             (extensions) )

type AccountCreateWithDelegationOperation struct {
	Fee            string     `json:"fee"`
	Delegation     string     `json:"delegation"`
	Creator        string     `json:"creator"`
	NewAccountName string     `json:"new_account_name"`
	Owner          *Authority `json:"owner"`
	Active         *Authority `json:"active"`
	Posting        *Authority `json:"posting"`
	MemoKey        string     `json:"memo_key"`
	JsonMetadata   string     `json:"json_metadata"`
	Extensions     []any      `json:"extensions"`
}

func (op *AccountCreateWithDelegationOperation) Type() OpType {
	return TypeAccountCreateWithDelegation
}

func (op *AccountCreateWithDelegationOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::witness_set_properties_operation,
//             (owner)
//             (props)
//             (extensions) )

type WitnessSetPropertiesOperation struct {
	Owner      string         `json:"owner"`
	Props      StringBytesMap `json:"props"`
	Extensions []any          `json:"extensions"`
}

func (op *WitnessSetPropertiesOperation) Type() OpType {
	return TypeWitnessSetProperties
}

func (op *WitnessSetPropertiesOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::account_update2_operation,
//             (account)
//             (owner)
//             (active)
//             (posting)
//             (memo_key)
//             (json_metadata)
//             (posting_json_metadata)
//             (extensions) )

type AccountUpdate2Operation struct {
	Account             string     `json:"account"`
	Owner               *Authority `json:"owner"`
	Active              *Authority `json:"active"`
	Posting             *Authority `json:"posting"`
	MemoKey             string     `json:"memo_key"`
	JsonMetadata        string     `json:"json_metadata"`
	PostingJsonMetadata string     `json:"posting_json_metadata"`
	Extensions          []any      `json:"extensions"`
}

func (op *AccountUpdate2Operation) Type() OpType {
	return TypeAccountUpdate2
}

func (op *AccountUpdate2Operation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::create_proposal_operation,
//             (creator)
//             (receiver)
//             (start_date)
//             (end_date)
//             (daily_pay)
//             (subject)
//             (permlink)
//             (extensions) )

type CreateProposalOperation struct {
	Creator    string `json:"creator"`
	Receiver   string `json:"receiver"`
	StartDate  *Time  `json:"start_date"`
	EndDate    *Time  `json:"end_date"`
	DailyPay   string `json:"daily_pay"`
	Subject    string `json:"subject"`
	Permlink   string `json:"permlink"`
	Extensions []any  `json:"extensions"`
}

func (op *CreateProposalOperation) Type() OpType {
	return TypeCreateProposal
}

func (op *CreateProposalOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::update_proposal_votes_operation,
//             (voter)
//             (proposal_ids)
//             (approve)
//             (extensions) )

type UpdateProposalVotesOperation struct {
	Voter       string   `json:"voter"`
	ProposalIDs []uint64 `json:"proposal_ids"`
	Approve     bool     `json:"approve"`
	Extensions  []any    `json:"extensions"`
}

func (op *UpdateProposalVotesOperation) Type() OpType {
	return TypeUpdateProposalVotes
}

func (op *UpdateProposalVotesOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::remove_proposal_operation,
//             (proposal_owner)
//             (proposal_ids)
//             (extensions) )

type RemoveProposalOperation struct {
	ProposalOwner string   `json:"proposal_owner"`
	ProposalIDs   []uint64 `json:"proposal_ids"`
	Extensions    []any    `json:"extensions"`
}

func (op *RemoveProposalOperation) Type() OpType {
	return TypeRemoveProposal
}

func (op *RemoveProposalOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::claim_reward_balance2_operation,
//             (account)
//             (reward_tokens)
//             (extensions) )

type ClaimRewardBalance2Operation struct {
	Account      string `json:"account"`
	Extensions   []any  `json:"extensions"`
	RewardTokens []any  `json:"reward_tokens"`
}

func (op *ClaimRewardBalance2Operation) Type() OpType {
	return TypeClaimRewardBalance2
}

func (op *ClaimRewardBalance2Operation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::vote2_operation,
//             (voter)
//             (author)
//             (permlink)
//             (rshares)
//             (extensions) )

type Vote2Operation struct {
	Voter      string `json:"voter"`
	Author     string `json:"author"`
	Permlink   string `json:"permlink"`
	RShares    string `json:"rshares"`
	Extensions []any  `json:"extensions"`
}

func (op *Vote2Operation) Type() OpType {
	return TypeVote2
}

func (op *Vote2Operation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::fill_convert_request_operation,
//             (owner)
//             (requestid)
//             (amount_in)
//             (amount_out) )

type FillConvertRequestOperation struct {
	Owner     string `json:"owner"`
	RequestID uint32 `json:"requestid"`
	AmountIn  string `json:"amount_in"`
	AmountOut string `json:"amount_out"`
}

func (op *FillConvertRequestOperation) Type() OpType {
	return TypeFillConvertRequest
}

func (op *FillConvertRequestOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::comment_reward_operation,
//             (author)
//             (permlink)
//             (payout) )

type CommentRewardOperation struct {
	Author   string `json:"author"`
	Permlink string `json:"permlink"`
	Payout   string `json:"payout"`
}

func (op *CommentRewardOperation) Type() OpType {
	return TypeCommentReward
}

func (op *CommentRewardOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::liquidity_reward_operation,
//             (owner)
//             (payout) )

type LiquidityRewardOperation struct {
	Owner  string `json:"owner"`
	Payout string `json:"payout"`
}

func (op *LiquidityRewardOperation) Type() OpType {
	return TypeLiquidityReward
}

func (op *LiquidityRewardOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::interest_operation,
//             (owner)
//             (interest) )

type InterestOperation struct {
	Owner    string `json:"owner"`
	Interest string `json:"interest"`
}

func (op *InterestOperation) Type() OpType {
	return TypeInterest
}

func (op *InterestOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::fill_vesting_withdraw_operation,
//             (from_account)
//             (to_account)
//             (withdrawn)
//             (deposited) )

type FillVestingWithdrawOperation struct {
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Withdrawn   string `json:"withdrawn"`
	Deposited   string `json:"deposited"`
}

func (op *FillVestingWithdrawOperation) Type() OpType {
	return TypeFillVestingWithdraw
}

func (op *FillVestingWithdrawOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::fill_order_operation,
//             (current_owner)
//             (current_orderid)
//             (current_pays)
//             (open_owner)
//             (open_orderid)
//             (open_pays) )

type FillOrderOperation struct {
	CurrentOwner   string `json:"current_owner"`
	CurrentOrderID uint32 `json:"current_orderid"`
	CurrentPays    string `json:"current_pays"`
	OpenOwner      string `json:"open_owner"`
	OpenOrderID    uint32 `json:"open_orderid"`
	OpenPays       string `json:"open_pays"`
}

func (op *FillOrderOperation) Type() OpType {
	return TypeFillOrder
}

func (op *FillOrderOperation) Data() any {
	return op
}

// FC_REFLECT( steemit::chain::fill_transfer_from_savings_operation,
//             (from)
//             (to)
//             (amount)
//             (request_id)
//             (memo) )

type FillTransferFromSavingsOperation struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    string `json:"amount"`
	RequestID uint32 `json:"request_id"`
	Memo      string `json:"memo"`
}

func (op *FillTransferFromSavingsOperation) Type() OpType {
	return TypeFillTransferFromSavings
}

func (op *FillTransferFromSavingsOperation) Data() any {
	return op
}

type UnknownOperation struct {
	kind OpType
	data *json.RawMessage
}

func (op *UnknownOperation) Type() OpType {
	return op.kind
}

func (op *UnknownOperation) Data() any {
	return op.data
}
