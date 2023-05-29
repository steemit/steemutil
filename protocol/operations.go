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
	Reporter string `json:"reporter"`
}

func (op *ReportOverProductionOperation) Type() OpType {
	return TypeReportOverProduction
}

func (op *ReportOverProductionOperation) Data() interface{} {
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

func (op *ConvertOperation) Data() interface{} {
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

func (op *FeedPublishOperation) Data() interface{} {
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

type ChainProperties struct {
	AccountCreationFee string `json:"account_creation_fee"`
	MaximumBlockSize   uint32 `json:"maximum_block_size"`
	SBDInterestRate    uint16 `json:"sbd_interest_rate"`
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

func (op *POWOperation) Data() interface{} {
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

func (op *AccountCreateOperation) Data() interface{} {
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

func (op *AccountUpdateOperation) Data() interface{} {
	return op
}

func (op *AccountUpdateOperation) MarshalTransaction(encoderObj *encoder.Encoder) error {
	enc := encoder.NewRollingEncoder(encoderObj)
	enc.EncodeUVarint(uint64(TypeAccountUpdate.Code()))
	enc.Encode(op.Account)
	enc.Encode(op.Owner)
	enc.Encode(op.Active)
	enc.Encode(op.Posting)
	enc.Encode(op.MemoKey)
	enc.Encode(op.JsonMetadata)
	return enc.Err()
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

func (op *TransferOperation) Data() interface{} {
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

func (op *TransferToVestingOperation) Data() interface{} {
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

func (op *WithdrawVestingOperation) Data() interface{} {
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

func (op *AccountWitnessVoteOperation) Data() interface{} {
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

func (op *AccountWitnessProxyOperation) Data() interface{} {
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

func (op *CommentOperation) Data() interface{} {
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

func (op *VoteOperation) Data() interface{} {
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

func (op *LimitOrderCreateOperation) Data() interface{} {
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

func (op *LimitOrderCancelOperation) Data() interface{} {
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

func (op *DeleteCommentOperation) Data() interface{} {
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
	Author               string        `json:"author"`
	Permlink             string        `json:"permlink"`
	MaxAcceptedPayout    string        `json:"max_accepted_payout"`
	PercentSteemDollars  uint16        `json:"percent_steem_dollars"`
	AllowVotes           bool          `json:"allow_votes"`
	AllowCurationRewards bool          `json:"allow_curation_rewards"`
	Extensions           []interface{} `json:"extensions"`
}

func (op *CommentOptionsOperation) Type() OpType {
	return TypeCommentOptions
}

func (op *CommentOptionsOperation) Data() interface{} {
	return op
}

type WitnessUpdateOperation struct {
}

func (op *WitnessUpdateOperation) Type() OpType {
	return TypeWitnessUpdate
}

func (op *WitnessUpdateOperation) Data() interface{} {
	return op
}

type CustomOperation struct {
}

func (op *CustomOperation) Type() OpType {
	return TypeCustom
}

func (op *CustomOperation) Data() interface{} {
	return op
}

type Authority struct {
	AccountAuths    StringInt64Map `json:"account_auths"`
	KeyAuths        StringInt64Map `json:"key_auths"`
	WeightThreshold uint32         `json:"weight_threshold"`
}

func (props *Authority) MarshalTransaction(encoderObj *encoder.Encoder) error {
	enc := encoder.NewRollingEncoder(encoderObj)
	enc.Encode(props.AccountAuths)
	enc.Encode(props.KeyAuths)
	enc.Encode(props.WeightThreshold)
	return enc.Err()
}

type SetWithdrawVestingRouteOperation struct {
}

func (op *SetWithdrawVestingRouteOperation) Type() OpType {
	return TypeSetWithdrawVestingRoute
}
func (op *SetWithdrawVestingRouteOperation) Data() interface{} {
	return op
}

type LimitOrderCreate2Operation struct {
}

func (op *LimitOrderCreate2Operation) Type() OpType {
	return TypeLimitOrderCreate2
}
func (op *LimitOrderCreate2Operation) Data() interface{} {
	return op
}

type ClaimAccountOperation struct {
}

func (op *ClaimAccountOperation) Type() OpType {
	return TypeClaimAccount
}
func (op *ClaimAccountOperation) Data() interface{} {
	return op
}

type CreateClaimedAccountOperation struct {
}

func (op *CreateClaimedAccountOperation) Type() OpType {
	return TypeCreateClaimedAccount
}
func (op *CreateClaimedAccountOperation) Data() interface{} {
	return op
}

type RequestAccountRecoveryOperation struct {
}

func (op *RequestAccountRecoveryOperation) Type() OpType {
	return TypeRequestAccountRecovery
}
func (op *RequestAccountRecoveryOperation) Data() interface{} {
	return op
}

type RecoverAccountOperation struct {
}

func (op *RecoverAccountOperation) Type() OpType {
	return TypeRecoverAccount
}
func (op *RecoverAccountOperation) Data() interface{} {
	return op
}

type ChangeRecoverAccountOperation struct {
}

func (op *ChangeRecoverAccountOperation) Type() OpType {
	return TypeChangeRecoveryAccount
}
func (op *ChangeRecoverAccountOperation) Data() interface{} {
	return op
}

type EscrowTransferOperation struct {
}

func (op *EscrowTransferOperation) Type() OpType {
	return TypeEscrowTransfer
}
func (op *EscrowTransferOperation) Data() interface{} {
	return op
}

type EscrowDisputeOperation struct {
}

func (op *EscrowDisputeOperation) Type() OpType {
	return TypeEscrowDispute
}
func (op *EscrowDisputeOperation) Data() interface{} {
	return op
}

type EescrowReleaseOperation struct {
}

func (op *EescrowReleaseOperation) Type() OpType {
	return TypeEscrowRelease
}
func (op *EescrowReleaseOperation) Data() interface{} {
	return op
}

type POW2Operation struct {
}

func (op *POW2Operation) Type() OpType {
	return TypePOW2
}
func (op *POW2Operation) Data() interface{} {
	return op
}

type EscrowApproveOperation struct {
}

func (op *EscrowApproveOperation) Type() OpType {
	return TypeEscrowApprove
}
func (op *EscrowApproveOperation) Data() interface{} {
	return op
}

type TransferToSavingsOperation struct {
}

func (op *TransferToSavingsOperation) Type() OpType {
	return TypeTransferToSavings
}
func (op *TransferToSavingsOperation) Data() interface{} {
	return op
}

type TransferFromSavingsOperation struct {
}

func (op *TransferFromSavingsOperation) Type() OpType {
	return TypeTransferFromSavings
}
func (op *TransferFromSavingsOperation) Data() interface{} {
	return op
}

type CancelTransferFromSavingsOperation struct {
}

func (op *CancelTransferFromSavingsOperation) Type() OpType {
	return TypeCancelTransferFromSavings
}
func (op *CancelTransferFromSavingsOperation) Data() interface{} {
	return op
}

type CustomBinaryOperation struct {
}

func (op *CustomBinaryOperation) Type() OpType {
	return TypeCustomBinary
}
func (op *CustomBinaryOperation) Data() interface{} {
	return op
}

type DeclineVotingRightsOperation struct {
}

func (op *DeclineVotingRightsOperation) Type() OpType {
	return TypeDeclineVotingRights
}
func (op *DeclineVotingRightsOperation) Data() interface{} {
	return op
}

type ResetAccountOperation struct {
}

func (op *ResetAccountOperation) Type() OpType {
	return TypeResetAccount
}
func (op *ResetAccountOperation) Data() interface{} {
	return op
}

type SetResetAccountOperation struct {
}

func (op *SetResetAccountOperation) Type() OpType {
	return TypeSetResetAccount
}
func (op *SetResetAccountOperation) Data() interface{} {
	return op
}

type ClaimRewardBalanceOperation struct {
}

func (op *ClaimRewardBalanceOperation) Type() OpType {
	return TypeClaimRewardBalance
}
func (op *ClaimRewardBalanceOperation) Data() interface{} {
	return op
}

type DelegateVestingSharesOperation struct {
}

func (op *DelegateVestingSharesOperation) Type() OpType {
	return TypeDelegateVestingShares
}
func (op *DelegateVestingSharesOperation) Data() interface{} {
	return op
}

type AccountCreateWithDelegationOperation struct {
}

func (op *AccountCreateWithDelegationOperation) Type() OpType {
	return TypeAccountCreateWithDelegation
}
func (op *AccountCreateWithDelegationOperation) Data() interface{} {
	return op
}

type WitnessSetPropertiesOperation struct {
}

func (op *WitnessSetPropertiesOperation) Type() OpType {
	return TypeWitnessSetProperties
}
func (op *WitnessSetPropertiesOperation) Data() interface{} {
	return op
}

type AccountUpdate2Operation struct {
}

func (op *AccountUpdate2Operation) Type() OpType {
	return TypeAccountUpdate2
}
func (op *AccountUpdate2Operation) Data() interface{} {
	return op
}

type CreateProposalOperation struct {
}

func (op *CreateProposalOperation) Type() OpType {
	return TypeCreateProposal
}
func (op *CreateProposalOperation) Data() interface{} {
	return op
}

type UpdateProposalVotesOperation struct {
}

func (op *UpdateProposalVotesOperation) Type() OpType {
	return TypeUpdateProposalVotes
}
func (op *UpdateProposalVotesOperation) Data() interface{} {
	return op
}

type RemoveProposalOperation struct {
}

func (op *RemoveProposalOperation) Type() OpType {
	return TypeRemoveProposal
}
func (op *RemoveProposalOperation) Data() interface{} {
	return op
}

type ClaimRewardBalance2Operation struct {
}

func (op *ClaimRewardBalance2Operation) Type() OpType {
	return TypeClaimRewardBalance2
}
func (op *ClaimRewardBalance2Operation) Data() interface{} {
	return op
}

type Vote2Operation struct {
}

func (op *Vote2Operation) Type() OpType {
	return TypeVote2
}
func (op *Vote2Operation) Data() interface{} {
	return op
}

type SmtSetupOperation struct {
}

func (op *SmtSetupOperation) Type() OpType {
	return TypeSmtSetup
}
func (op *SmtSetupOperation) Data() interface{} {
	return op
}

type SmtSetupEmissionsOperation struct {
}

func (op *SmtSetupEmissionsOperation) Type() OpType {
	return TypeSmtSetupEmissions
}
func (op *SmtSetupEmissionsOperation) Data() interface{} {
	return op
}

type SmtSetupIcoTierOperation struct {
}

func (op *SmtSetupIcoTierOperation) Type() OpType {
	return TypeSmtSetupIcoTier
}
func (op *SmtSetupIcoTierOperation) Data() interface{} {
	return op
}

type SmtSetSetupParametersOperation struct {
}

func (op *SmtSetSetupParametersOperation) Type() OpType {
	return TypeSmtSetSetupParameters
}
func (op *SmtSetSetupParametersOperation) Data() interface{} {
	return op
}

type SmtSetRuntimeParaetersOperation struct {
}

func (op *SmtSetRuntimeParaetersOperation) Type() OpType {
	return TypeSmtSetRuntimeParameters
}
func (op *SmtSetRuntimeParaetersOperation) Data() interface{} {
	return op
}

type SmtCreateOperation struct {
}

func (op *SmtCreateOperation) Type() OpType {
	return TypeSmtCreate
}
func (op *SmtCreateOperation) Data() interface{} {
	return op
}

type SmtContributeOperation struct {
}

func (op *SmtContributeOperation) Type() OpType {
	return TypeSmtContribute
}
func (op *SmtContributeOperation) Data() interface{} {
	return op
}

type FillConvertRequestOperation struct {
}

func (op *FillConvertRequestOperation) Type() OpType {
	return TypeFillConvertRequest
}
func (op *FillConvertRequestOperation) Data() interface{} {
	return op
}

type AuthorRewardOperation struct {
}

func (op *AuthorRewardOperation) Type() OpType {
	return TypeAuthorReward
}
func (op *AuthorRewardOperation) Data() interface{} {
	return op
}

type CurationRewardOperation struct {
}

func (op *CurationRewardOperation) Type() OpType {
	return TypeCurationReward
}
func (op *CurationRewardOperation) Data() interface{} {
	return op
}

type CommentRewardOperation struct {
}

func (op *CommentRewardOperation) Type() OpType {
	return TypeCommentReward
}
func (op *CommentRewardOperation) Data() interface{} {
	return op
}

type LiquidityRewardOperation struct {
}

func (op *LiquidityRewardOperation) Type() OpType {
	return TypeLiquidityReward
}
func (op *LiquidityRewardOperation) Data() interface{} {
	return op
}

type InterestOperation struct {
}

func (op *InterestOperation) Type() OpType {
	return TypeInterest
}
func (op *InterestOperation) Data() interface{} {
	return op
}

type FillVestingWithdrawOperation struct {
}

func (op *FillVestingWithdrawOperation) Type() OpType {
	return TypeFillVestingWithdraw
}
func (op *FillVestingWithdrawOperation) Data() interface{} {
	return op
}

type FillOrderOperation struct {
}

func (op *FillOrderOperation) Type() OpType {
	return TypeFillOrder
}
func (op *FillOrderOperation) Data() interface{} {
	return op
}

type ShutdownWitnessOperation struct {
}

func (op *ShutdownWitnessOperation) Type() OpType {
	return TypeShutdownWitness
}
func (op *ShutdownWitnessOperation) Data() interface{} {
	return op
}

type FillTransferFromSavingsOperation struct {
}

func (op *FillTransferFromSavingsOperation) Type() OpType {
	return TypeFillTransferFromSavings
}
func (op *FillTransferFromSavingsOperation) Data() interface{} {
	return op
}

type HardforkOperation struct {
}

func (op *HardforkOperation) Type() OpType {
	return TypeHardfork
}
func (op *HardforkOperation) Data() interface{} {
	return op
}

type CommentPayoutUpdateOperation struct {
}

func (op *CommentPayoutUpdateOperation) Type() OpType {
	return TypeCommentPayoutUpdate
}
func (op *CommentPayoutUpdateOperation) Data() interface{} {
	return op
}

type ReturnVestingDelegationOperation struct {
}

func (op *ReturnVestingDelegationOperation) Type() OpType {
	return TypeReturnVestingDelegation
}
func (op *ReturnVestingDelegationOperation) Data() interface{} {
	return op
}

type CommentBenefactorRewardOperation struct {
}

func (op *CommentBenefactorRewardOperation) Type() OpType {
	return TypeCommentBenefactorReward
}
func (op *CommentBenefactorRewardOperation) Data() interface{} {
	return op
}

type ProducerRewardOperation struct {
}

func (op *ProducerRewardOperation) Type() OpType {
	return TypeProducerReward
}
func (op *ProducerRewardOperation) Data() interface{} {
	return op
}

type ClearNullAccountBalanceOperation struct {
}

func (op *ClearNullAccountBalanceOperation) Type() OpType {
	return TypeClearNullAccountBalance
}
func (op *ClearNullAccountBalanceOperation) Data() interface{} {
	return op
}

type ProposalPayOperation struct {
}

func (op *ProposalPayOperation) Type() OpType {
	return TypeProposalPay
}
func (op *ProposalPayOperation) Data() interface{} {
	return op
}

type SpsFundOperation struct {
}

func (op *SpsFundOperation) Type() OpType {
	return TypeSpsFund
}
func (op *SpsFundOperation) Data() interface{} {
	return op
}

type UnknownOperation struct {
	kind OpType
	data *json.RawMessage
}

func (op *UnknownOperation) Type() OpType {
	return op.kind
}

func (op *UnknownOperation) Data() interface{} {
	return op.data
}
