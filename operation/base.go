package operation

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// OpType represents a Steem operation type, i.e. vote, comment, pow and so on.
type OpType string

// Code returns the operation code associated with the given operation type.
func (kind OpType) Code() uint16 {
	return opCodes[kind]
}

const (
	TypeVote                        OpType = "vote"
	TypeComment                     OpType = "comment"
	TypeTransfer                    OpType = "transfer"
	TypeTransferToVesting           OpType = "transfer_to_vesting"
	TypeWithdrawVesting             OpType = "withdraw_vesting"
	TypeLimitOrderCreate            OpType = "limit_order_create"
	TypeLimitOrderCancel            OpType = "limit_order_cancel"
	TypeFeedPublish                 OpType = "feed_publish"
	TypeConvert                     OpType = "convert"
	TypeAccountCreate               OpType = "account_create"
	TypeAccountUpdate               OpType = "account_update"
	TypeWitnessUpdate               OpType = "witness_update"
	TypeAccountWitnessVote          OpType = "account_witness_vote"
	TypeAccountWitnessProxy         OpType = "account_witness_proxy"
	TypePOW                         OpType = "pow"
	TypeCustom                      OpType = "custom"
	TypeReportOverProduction        OpType = "report_over_production"
	TypeDeleteComment               OpType = "delete_comment"
	TypeCustomJSON                  OpType = "custom_json"
	TypeCommentOptions              OpType = "comment_options"
	TypeSetWithdrawVestingRoute     OpType = "set_withdraw_vesting_route"
	TypeLimitOrderCreate2           OpType = "limit_order_create2"
	TypeClaimAccount                OpType = "claim_account"
	TypeCreateClaimedAccount        OpType = "create_claimed_account"
	TypeRequestAccountRecovery      OpType = "request_account_recovery"
	TypeRecoverAccount              OpType = "recover_account"
	TypeChangeRecoveryAccount       OpType = "change_recovery_account"
	TypeEscrowTransfer              OpType = "escrow_transfer"
	TypeEscrowDispute               OpType = "escrow_dispute"
	TypeEscrowRelease               OpType = "escrow_release"
	TypePOW2                        OpType = "pow2"
	TypeEscrowApprove               OpType = "escrow_approve"
	TypeTransferToSavings           OpType = "transfer_to_savings"
	TypeTransferFromSavings         OpType = "transfer_from_savings"
	TypeCancelTransferFromSavings   OpType = "cancel_transfer_from_savings"
	TypeCustomBinary                OpType = "custom_binary"
	TypeDeclineVotingRights         OpType = "decline_voting_rights"
	TypeResetAccount                OpType = "reset_account"
	TypeSetResetAccount             OpType = "set_reset_account"
	TypeClaimRewardBalance          OpType = "claim_reward_balance"
	TypeDelegateVestingShares       OpType = "delegate_vesting_shares"
	TypeAccountCreateWithDelegation OpType = "account_create_with_delegation"
	TypeWitnessSetProperties        OpType = "witness_set_properties"
	TypeAccountUpdate2              OpType = "account_update2"
	TypeCreateProposal              OpType = "create_proposal"
	TypeUpdateProposalVotes         OpType = "update_proposal_votes"
	TypeRemoveProposal              OpType = "remove_proposal"
	TypeClaimRewardBalance2         OpType = "claim_reward_balance2"
	TypeVote2                       OpType = "vote2"
	TypeSmtSetup                    OpType = "smt_setup"
	TypeSmtSetupEmissions           OpType = "smt_setup_emissions"
	TypeSmtSetupIcoTier             OpType = "smt_setup_ico_tier"
	TypeSmtSetSetupParameters       OpType = "smt_set_setup_parameters"
	TypeSmtSetRuntimeParameters     OpType = "smt_set_runtime_parameters"
	TypeSmtCreate                   OpType = "smt_create"
	TypeSmtContribute               OpType = "smt_contribute"
	TypeFillConvertRequest          OpType = "fill_convert_request"
	TypeAuthorReward                OpType = "author_reward"
	TypeCurationReward              OpType = "curation_reward"
	TypeCommentReward               OpType = "comment_reward"
	TypeLiquidityReward             OpType = "liquidity_reward"
	TypeInterest                    OpType = "interest"
	TypeFillVestingWithdraw         OpType = "fill_vesting_withdraw"
	TypeFillOrder                   OpType = "fill_order"
	TypeShutdownWitness             OpType = "shutdown_witness"
	TypeFillTransferFromSavings     OpType = "fill_transfer_from_savings"
	TypeHardfork                    OpType = "hardfork"
	TypeCommentPayoutUpdate         OpType = "comment_payout_update"
	TypeReturnVestingDelegation     OpType = "return_vesting_delegation"
	TypeCommentBenefactorReward     OpType = "comment_benefactor_reward"
	TypeProducerReward              OpType = "producer_reward"
	TypeClearNullAccountBalance     OpType = "clear_null_account_balance"
	TypeProposalPay                 OpType = "proposal_pay"
	TypeSpsFund                     OpType = "sps_fund"
)

// opCodes keeps mapping operation type -> operation code.
var opCodes = map[OpType]uint16{
	TypeVote:                        0,
	TypeComment:                     1,
	TypeTransfer:                    2,
	TypeTransferToVesting:           3,
	TypeWithdrawVesting:             4,
	TypeLimitOrderCreate:            5,
	TypeLimitOrderCancel:            6,
	TypeFeedPublish:                 7,
	TypeConvert:                     8,
	TypeAccountCreate:               9,
	TypeAccountUpdate:               10,
	TypeWitnessUpdate:               11,
	TypeAccountWitnessVote:          12,
	TypeAccountWitnessProxy:         13,
	TypePOW:                         14,
	TypeCustom:                      15,
	TypeReportOverProduction:        16,
	TypeDeleteComment:               17,
	TypeCustomJSON:                  18,
	TypeCommentOptions:              19,
	TypeSetWithdrawVestingRoute:     20,
	TypeLimitOrderCreate2:           21,
	TypeClaimAccount:                22,
	TypeCreateClaimedAccount:        23,
	TypeRequestAccountRecovery:      24,
	TypeRecoverAccount:              25,
	TypeChangeRecoveryAccount:       26,
	TypeEscrowTransfer:              27,
	TypeEscrowDispute:               28,
	TypeEscrowRelease:               29,
	TypePOW2:                        30,
	TypeEscrowApprove:               31,
	TypeTransferToSavings:           32,
	TypeTransferFromSavings:         33,
	TypeCancelTransferFromSavings:   34,
	TypeCustomBinary:                35,
	TypeDeclineVotingRights:         36,
	TypeResetAccount:                37,
	TypeSetResetAccount:             38,
	TypeClaimRewardBalance:          39,
	TypeDelegateVestingShares:       40,
	TypeAccountCreateWithDelegation: 41,
	TypeWitnessSetProperties:        42,
	TypeAccountUpdate2:              43,
	TypeCreateProposal:              44,
	TypeUpdateProposalVotes:         45,
	TypeRemoveProposal:              46,
	TypeClaimRewardBalance2:         47,
	TypeVote2:                       48,
	TypeSmtSetup:                    49,
	TypeSmtSetupEmissions:           50,
	TypeSmtSetupIcoTier:             51,
	TypeSmtSetSetupParameters:       52,
	TypeSmtSetRuntimeParameters:     53,
	TypeSmtCreate:                   54,
	TypeSmtContribute:               55,
	TypeFillConvertRequest:          56,
	TypeAuthorReward:                57,
	TypeCurationReward:              58,
	TypeCommentReward:               59,
	TypeLiquidityReward:             60,
	TypeInterest:                    61,
	TypeFillVestingWithdraw:         62,
	TypeFillOrder:                   63,
	TypeShutdownWitness:             64,
	TypeFillTransferFromSavings:     65,
	TypeHardfork:                    66,
	TypeCommentPayoutUpdate:         67,
	TypeReturnVestingDelegation:     68,
	TypeCommentBenefactorReward:     69,
	TypeProducerReward:              70,
	TypeClearNullAccountBalance:     71,
	TypeProposalPay:                 72,
	TypeSpsFund:                     73,
}

// dataObjects keeps mapping operation type -> operation data object.
// This is used later on to unmarshal operation data based on the operation type.
var dataObjects = map[OpType]IOperation{
	TypeVote:    &VoteOperation{},
	TypeComment: &CommentOperation{},

	TypeTransfer:          &TransferOperation{},
	TypeTransferToVesting: &TransferToVestingOperation{},
	TypeWithdrawVesting:   &WithdrawVestingOperation{},

	TypeLimitOrderCreate: &LimitOrderCreateOperation{},
	TypeLimitOrderCancel: &LimitOrderCancelOperation{},

	TypeFeedPublish: &FeedPublishOperation{},
	TypeConvert:     &ConvertOperation{},

	TypeAccountCreate: &AccountCreateOperation{},
	TypeAccountUpdate: &AccountUpdateOperation{},

	TypeWitnessUpdate:       &WitnessUpdateOperation{},
	TypeAccountWitnessVote:  &AccountWitnessVoteOperation{},
	TypeAccountWitnessProxy: &AccountWitnessProxyOperation{},

	TypePOW: &POWOperation{},

	TypeCustom: &CustomOperation{},

	TypeReportOverProduction: &ReportOverProductionOperation{},

	TypeDeleteComment:               &DeleteCommentOperation{},
	TypeCustomJSON:                  &CustomJSONOperation{},
	TypeCommentOptions:              &CommentOptionsOperation{},
	TypeSetWithdrawVestingRoute:     &SetWithdrawVestingRouteOperation{},
	TypeLimitOrderCreate2:           &LimitOrderCreate2Operation{},
	TypeClaimAccount:                &ClaimAccountOperation{},
	TypeCreateClaimedAccount:        &CreateClaimedAccountOperation{},
	TypeRequestAccountRecovery:      &RequestAccountRecoveryOperation{},
	TypeRecoverAccount:              &RecoverAccountOperation{},
	TypeChangeRecoveryAccount:       &ChangeRecoverAccountOperation{},
	TypeEscrowTransfer:              &EscrowTransferOperation{},
	TypeEscrowDispute:               &EscrowDisputeOperation{},
	TypeEscrowRelease:               &EescrowReleaseOperation{},
	TypePOW2:                        &POW2Operation{},
	TypeEscrowApprove:               &EscrowApproveOperation{},
	TypeTransferToSavings:           &TransferToSavingsOperation{},
	TypeTransferFromSavings:         &TransferFromSavingsOperation{},
	TypeCancelTransferFromSavings:   &CancelTransferFromSavingsOperation{},
	TypeCustomBinary:                &CustomBinaryOperation{},
	TypeDeclineVotingRights:         &DeclineVotingRightsOperation{},
	TypeResetAccount:                &ResetAccountOperation{},
	TypeSetResetAccount:             &SetResetAccountOperation{},
	TypeClaimRewardBalance:          &ClaimRewardBalanceOperation{},
	TypeDelegateVestingShares:       &DelegateVestingSharesOperation{},
	TypeAccountCreateWithDelegation: &AccountCreateWithDelegationOperation{},
	TypeWitnessSetProperties:        &WitnessSetPropertiesOperation{},
	TypeAccountUpdate2:              &AccountUpdate2Operation{},
	TypeCreateProposal:              &CreateProposalOperation{},
	TypeUpdateProposalVotes:         &UpdateProposalVotesOperation{},
	TypeRemoveProposal:              &RemoveProposalOperation{},
	TypeClaimRewardBalance2:         &ClaimRewardBalance2Operation{},
	TypeVote2:                       &Vote2Operation{},

	TypeSmtSetup:                &SmtSetupOperation{},
	TypeSmtSetupEmissions:       &SmtSetupEmissionsOperation{},
	TypeSmtSetupIcoTier:         &SmtSetupIcoTierOperation{},
	TypeSmtSetSetupParameters:   &SmtSetSetupParametersOperation{},
	TypeSmtSetRuntimeParameters: &SmtSetRuntimeParaetersOperation{},
	TypeSmtCreate:               &SmtCreateOperation{},
	TypeSmtContribute:           &SmtContributeOperation{},

	TypeFillConvertRequest:      &FillConvertRequestOperation{},
	TypeAuthorReward:            &AuthorRewardOperation{},
	TypeCurationReward:          &CurationRewardOperation{},
	TypeCommentReward:           &CommentRewardOperation{},
	TypeLiquidityReward:         &LiquidityRewardOperation{},
	TypeInterest:                &InterestOperation{},
	TypeFillVestingWithdraw:     &FillVestingWithdrawOperation{},
	TypeFillOrder:               &FillOrderOperation{},
	TypeShutdownWitness:         &ShutdownWitnessOperation{},
	TypeFillTransferFromSavings: &FillTransferFromSavingsOperation{},
	TypeHardfork:                &HardforkOperation{},
	TypeCommentPayoutUpdate:     &CommentPayoutUpdateOperation{},
	TypeReturnVestingDelegation: &ReturnVestingDelegationOperation{},
	TypeCommentBenefactorReward: &CommentBenefactorRewardOperation{},
	TypeProducerReward:          &ProducerRewardOperation{},
	TypeClearNullAccountBalance: &ClearNullAccountBalanceOperation{},
	TypeProposalPay:             &ProposalPayOperation{},
	TypeSpsFund:                 &SpsFundOperation{},
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

// Operation represents an operation stored in a transaction.
type IOperation interface {
	// Type returns the operation type as present in the operation object, element [0].
	Type() OpType

	// Data returns the operation data as present in the operation object, element [1].
	//
	// When the operation type is known to this package, this field contains
	// the operation data object associated with the given operation type,
	// e.g. Type is TypeVote -> Data contains *VoteOperation.
	// Otherwise this field contains raw JSON (type *json.RawMessage).
	Data() interface{}
}

type Operations []IOperation

func (ops *Operations) UnmarshalJSON(data []byte) error {
	var tuples []*operationTuple
	if err := json.Unmarshal(data, &tuples); err != nil {
		return err
	}

	items := make([]IOperation, 0, len(tuples))
	for _, tuple := range tuples {
		items = append(items, tuple.Data)
	}

	*ops = items
	return nil
}

func (ops Operations) MarshalJSON() ([]byte, error) {
	tuples := make([]*operationTuple, 0, len(ops))
	for _, op := range ops {
		tuples = append(tuples, &operationTuple{
			Type: op.Type(),
			Data: op.Data().(IOperation),
		})
	}
	return json.Marshal(tuples)
}

type operationTuple struct {
	Type OpType
	Data IOperation
}

func (op *operationTuple) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		op.Type,
		op.Data,
	})
}

func (op *operationTuple) UnmarshalJSON(data []byte) error {
	// The operation object is [opType, opBody].
	raw := make([]*json.RawMessage, 2)
	if err := json.Unmarshal(data, &raw); err != nil {
		return errors.Wrapf(err, "failed to unmarshal operation object: %v", string(data))
	}
	if len(raw) != 2 {
		return errors.Errorf("invalid operation object: %v", string(data))
	}

	// Unmarshal the type.
	var opType OpType
	if err := json.Unmarshal(*raw[0], &opType); err != nil {
		return errors.Wrapf(err, "failed to unmarshal Operation.Type: %v", string(*raw[0]))
	}

	// Unmarshal the data.
	var opData IOperation
	template, ok := dataObjects[opType]
	if ok {
		opData = reflect.New(
			reflect.Indirect(reflect.ValueOf(template)).Type(),
		).Interface().(IOperation)

		if err := json.Unmarshal(*raw[1], opData); err != nil {
			return errors.Wrapf(err, "failed to unmarshal Operation.Data: %v", string(*raw[1]))
		}
	} else {
		opData = &UnknownOperation{opType, raw[1]}
	}

	// Update fields.
	op.Type = opType
	op.Data = opData
	return nil
}
