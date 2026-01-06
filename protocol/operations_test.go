package protocol

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/steemit/steemutil/encoder"
)

func TestVoteOperation_MarshalTransaction(t *testing.T) {
	op := &VoteOperation{
		Voter:    "xeroc",
		Author:   "xeroc",
		Permlink: "piston",
		Weight:   10000,
	}

	expectedHex := "00057865726f63057865726f6306706973746f6e1027"

	var b bytes.Buffer
	encoder := encoder.NewEncoder(&b)

	if err := encoder.Encode(op); err != nil {
		t.Error(err)
	}

	serializedHex := hex.EncodeToString(b.Bytes())

	if serializedHex != expectedHex {
		t.Errorf("expected %v, got %v", expectedHex, serializedHex)
	}
}

func TestCommentOperation_MarshalTransaction(t *testing.T) {
	op := &CommentOperation{
		Author:         "ety001",
		Title:          "test post",
		Permlink:       "ety001-test-post",
		ParentAuthor:   "",
		ParentPermlink: "test",
		Body:           "test post body",
		JsonMetadata:   "{}",
	}

	expectedHex := "0100047465737406657479303031106574793030312d746573742d706f7374097465737420706f73740e7465737420706f737420626f6479027b7d"

	var b bytes.Buffer
	encoder := encoder.NewEncoder(&b)

	if err := encoder.Encode(op); err != nil {
		t.Error(err)
	}

	serializedHex := hex.EncodeToString(b.Bytes())

	if serializedHex != expectedHex {
		t.Errorf("expected %v, got %v", expectedHex, serializedHex)
	}
}

func TestTransferOperation_Type(t *testing.T) {
	op := &TransferOperation{
		From:   "alice",
		To:     "bob",
		Amount: "1.000 STEEM",
		Memo:   "test memo",
	}

	if op.Type() != TypeTransfer {
		t.Errorf("expected TypeTransfer, got %v", op.Type())
	}
}

func TestTransferOperation_Data(t *testing.T) {
	op := &TransferOperation{
		From:   "alice",
		To:     "bob",
		Amount: "1.000 STEEM",
		Memo:   "test memo",
	}

	data := op.Data()
	if data != op {
		t.Error("Data() should return the operation itself")
	}
}

func TestAccountCreateOperation_Type(t *testing.T) {
	op := &AccountCreateOperation{
		Fee:            "0.000 STEEM",
		Creator:        "creator",
		NewAccountName: "newaccount",
		MemoKey:        "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA",
		JsonMetadata:   "{}",
	}

	if op.Type() != TypeAccountCreate {
		t.Errorf("expected TypeAccountCreate, got %v", op.Type())
	}
}

func TestCommentOptionsOperation_Type(t *testing.T) {
	op := &CommentOptionsOperation{
		Author:               "author",
		Permlink:             "permlink",
		MaxAcceptedPayout:    "100.000 SBD",
		PercentSteemDollars:  50,
		AllowVotes:           true,
		AllowCurationRewards: true,
		Extensions:           []interface{}{},
	}

	if op.Type() != TypeCommentOptions {
		t.Errorf("expected TypeCommentOptions, got %v", op.Type())
	}
}

func TestCustomJSONOperation_Type(t *testing.T) {
	op := &CustomJSONOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{"alice"},
		ID:                   "follow",
		JSON:                 `["follow",{"follower":"alice","following":"bob","what":["blog"]}]`,
	}

	if op.Type() != TypeCustomJSON {
		t.Errorf("expected TypeCustomJSON, got %v", op.Type())
	}
}

func TestDeleteCommentOperation_Type(t *testing.T) {
	op := &DeleteCommentOperation{
		Author:   "author",
		Permlink: "permlink",
	}

	if op.Type() != TypeDeleteComment {
		t.Errorf("expected TypeDeleteComment, got %v", op.Type())
	}
}

func TestLimitOrderCreateOperation_Type(t *testing.T) {
	op := &LimitOrderCreateOperation{
		Owner:        "owner",
		OrderID:      1,
		AmountToSell: "1.000 STEEM",
		MinToReceive: "1.000 SBD",
		FillOrKill:   false,
		Expiration:   nil,
	}

	if op.Type() != TypeLimitOrderCreate {
		t.Errorf("expected TypeLimitOrderCreate, got %v", op.Type())
	}
}

func TestLimitOrderCancelOperation_Type(t *testing.T) {
	op := &LimitOrderCancelOperation{
		Owner:   "owner",
		OrderID: 1,
	}

	if op.Type() != TypeLimitOrderCancel {
		t.Errorf("expected TypeLimitOrderCancel, got %v", op.Type())
	}
}

func TestAccountWitnessVoteOperation_Type(t *testing.T) {
	op := &AccountWitnessVoteOperation{
		Account: "account",
		Witness: "witness",
		Approve: true,
	}

	if op.Type() != TypeAccountWitnessVote {
		t.Errorf("expected TypeAccountWitnessVote, got %v", op.Type())
	}
}

func TestAccountWitnessProxyOperation_Type(t *testing.T) {
	op := &AccountWitnessProxyOperation{
		Account: "account",
		Proxy:   "proxy",
	}

	if op.Type() != TypeAccountWitnessProxy {
		t.Errorf("expected TypeAccountWitnessProxy, got %v", op.Type())
	}
}

func TestWitnessUpdateOperation_Type(t *testing.T) {
	op := &WitnessUpdateOperation{
		Owner:           "owner",
		URL:             "https://example.com",
		BlockSigningKey: "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA",
		Props: &ChainProperties{
			AccountCreationFee: "0.100 STEEM",
			MaximumBlockSize:   65536,
			SBDInterestRate:    1000,
		},
		Fee: "0.000 STEEM",
	}

	if op.Type() != TypeWitnessUpdate {
		t.Errorf("expected TypeWitnessUpdate, got %v", op.Type())
	}
}

func TestSetWithdrawVestingRouteOperation_Type(t *testing.T) {
	op := &SetWithdrawVestingRouteOperation{
		FromAccount: "from",
		ToAccount:   "to",
		Percent:     50,
		AutoVest:    true,
	}

	if op.Type() != TypeSetWithdrawVestingRoute {
		t.Errorf("expected TypeSetWithdrawVestingRoute, got %v", op.Type())
	}
}

func TestLimitOrderCreate2Operation_Type(t *testing.T) {
	op := &LimitOrderCreate2Operation{
		Owner:        "owner",
		OrderID:      1,
		AmountToSell: "1.000 STEEM",
		FillOrKill:   false,
		Expiration:   nil,
	}

	if op.Type() != TypeLimitOrderCreate2 {
		t.Errorf("expected TypeLimitOrderCreate2, got %v", op.Type())
	}
}

func TestClaimAccountOperation_Type(t *testing.T) {
	op := &ClaimAccountOperation{
		Creator:    "creator",
		Fee:        "0.000 STEEM",
		Extensions: []interface{}{},
	}

	if op.Type() != TypeClaimAccount {
		t.Errorf("expected TypeClaimAccount, got %v", op.Type())
	}
}

func TestCreateClaimedAccountOperation_Type(t *testing.T) {
	op := &CreateClaimedAccountOperation{
		Creator:        "creator",
		NewAccountName: "newaccount",
		MemoKey:        "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA",
		JsonMetadata:   "{}",
		Extensions:     []interface{}{},
	}

	if op.Type() != TypeCreateClaimedAccount {
		t.Errorf("expected TypeCreateClaimedAccount, got %v", op.Type())
	}
}

func TestRequestAccountRecoveryOperation_Type(t *testing.T) {
	op := &RequestAccountRecoveryOperation{
		RecoveryAccount:  "recovery",
		AccountToRecover: "account",
		NewOwnerAuthority: &Authority{
			WeightThreshold: 1,
		},
		Extensions: []interface{}{},
	}

	if op.Type() != TypeRequestAccountRecovery {
		t.Errorf("expected TypeRequestAccountRecovery, got %v", op.Type())
	}
}

func TestRecoverAccountOperation_Type(t *testing.T) {
	op := &RecoverAccountOperation{
		AccountToRecover:     "account",
		NewOwnerAuthority:    &Authority{WeightThreshold: 1},
		RecentOwnerAuthority: &Authority{WeightThreshold: 1},
		Extensions:           []interface{}{},
	}

	if op.Type() != TypeRecoverAccount {
		t.Errorf("expected TypeRecoverAccount, got %v", op.Type())
	}
}

func TestChangeRecoveryAccountOperation_Type(t *testing.T) {
	op := &ChangeRecoveryAccountOperation{
		AccountToRecover:   "account",
		NewRecoveryAccount: "newrecovery",
		Extensions:         []interface{}{},
	}

	if op.Type() != TypeChangeRecoveryAccount {
		t.Errorf("expected TypeChangeRecoveryAccount, got %v", op.Type())
	}
}

func TestEscrowTransferOperation_Type(t *testing.T) {
	op := &EscrowTransferOperation{
		From:        "from",
		To:          "to",
		Agent:       "agent",
		EscrowID:    1,
		SBDAmount:   "1.000 SBD",
		SteemAmount: "1.000 STEEM",
		Fee:         "0.000 STEEM",
		JsonMeta:    "{}",
	}

	if op.Type() != TypeEscrowTransfer {
		t.Errorf("expected TypeEscrowTransfer, got %v", op.Type())
	}
}

func TestEscrowDisputeOperation_Type(t *testing.T) {
	op := &EscrowDisputeOperation{
		From:     "from",
		To:       "to",
		Agent:    "agent",
		Who:      "who",
		EscrowID: 1,
	}

	if op.Type() != TypeEscrowDispute {
		t.Errorf("expected TypeEscrowDispute, got %v", op.Type())
	}
}

func TestEscrowReleaseOperation_Type(t *testing.T) {
	op := &EscrowReleaseOperation{
		From:        "from",
		To:          "to",
		Agent:       "agent",
		Who:         "who",
		Receiver:    "receiver",
		EscrowID:    1,
		SBDAmount:   "1.000 SBD",
		SteemAmount: "1.000 STEEM",
	}

	if op.Type() != TypeEscrowRelease {
		t.Errorf("expected TypeEscrowRelease, got %v", op.Type())
	}
}

func TestEscrowApproveOperation_Type(t *testing.T) {
	op := &EscrowApproveOperation{
		From:     "from",
		To:       "to",
		Agent:    "agent",
		Who:      "who",
		EscrowID: 1,
		Approve:  true,
	}

	if op.Type() != TypeEscrowApprove {
		t.Errorf("expected TypeEscrowApprove, got %v", op.Type())
	}
}

func TestTransferToSavingsOperation_Type(t *testing.T) {
	op := &TransferToSavingsOperation{
		From:   "from",
		To:     "to",
		Amount: "1.000 STEEM",
		Memo:   "memo",
	}

	if op.Type() != TypeTransferToSavings {
		t.Errorf("expected TypeTransferToSavings, got %v", op.Type())
	}
}

func TestTransferFromSavingsOperation_Type(t *testing.T) {
	op := &TransferFromSavingsOperation{
		From:      "from",
		RequestID: 1,
		To:        "to",
		Amount:    "1.000 STEEM",
		Memo:      "memo",
	}

	if op.Type() != TypeTransferFromSavings {
		t.Errorf("expected TypeTransferFromSavings, got %v", op.Type())
	}
}

func TestCancelTransferFromSavingsOperation_Type(t *testing.T) {
	op := &CancelTransferFromSavingsOperation{
		From:      "from",
		RequestID: 1,
	}

	if op.Type() != TypeCancelTransferFromSavings {
		t.Errorf("expected TypeCancelTransferFromSavings, got %v", op.Type())
	}
}

func TestCustomBinaryOperation_Type(t *testing.T) {
	op := &CustomBinaryOperation{
		ID:        "test",
		DataBytes: "data",
	}

	if op.Type() != TypeCustomBinary {
		t.Errorf("expected TypeCustomBinary, got %v", op.Type())
	}
}

func TestDeclineVotingRightsOperation_Type(t *testing.T) {
	op := &DeclineVotingRightsOperation{
		Account: "account",
		Decline: true,
	}

	if op.Type() != TypeDeclineVotingRights {
		t.Errorf("expected TypeDeclineVotingRights, got %v", op.Type())
	}
}

func TestClaimRewardBalanceOperation_Type(t *testing.T) {
	op := &ClaimRewardBalanceOperation{
		Account:     "account",
		RewardSteem: "1.000 STEEM",
		RewardSBD:   "1.000 SBD",
		RewardVests: "1.000000 VESTS",
	}

	if op.Type() != TypeClaimRewardBalance {
		t.Errorf("expected TypeClaimRewardBalance, got %v", op.Type())
	}
}

func TestDelegateVestingSharesOperation_Type(t *testing.T) {
	op := &DelegateVestingSharesOperation{
		Delegator:     "delegator",
		Delegatee:     "delegatee",
		VestingShares: "1.000000 VESTS",
	}

	if op.Type() != TypeDelegateVestingShares {
		t.Errorf("expected TypeDelegateVestingShares, got %v", op.Type())
	}
}

func TestAccountCreateWithDelegationOperation_Type(t *testing.T) {
	op := &AccountCreateWithDelegationOperation{
		Fee:            "0.000 STEEM",
		Delegation:     "1.000000 VESTS",
		Creator:        "creator",
		NewAccountName: "newaccount",
		MemoKey:        "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA",
		JsonMetadata:   "{}",
		Extensions:     []interface{}{},
	}

	if op.Type() != TypeAccountCreateWithDelegation {
		t.Errorf("expected TypeAccountCreateWithDelegation, got %v", op.Type())
	}
}

func TestWitnessSetPropertiesOperation_Type(t *testing.T) {
	op := &WitnessSetPropertiesOperation{
		Owner:      "owner",
		Props:      StringBytesMap{"key": "value"},
		Extensions: []interface{}{},
	}

	if op.Type() != TypeWitnessSetProperties {
		t.Errorf("expected TypeWitnessSetProperties, got %v", op.Type())
	}
}

func TestAccountUpdate2Operation_Type(t *testing.T) {
	op := &AccountUpdate2Operation{
		Account:             "account",
		MemoKey:             "STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA",
		JsonMetadata:        "{}",
		PostingJsonMetadata: "{}",
		Extensions:          []interface{}{},
	}

	if op.Type() != TypeAccountUpdate2 {
		t.Errorf("expected TypeAccountUpdate2, got %v", op.Type())
	}
}

func TestVote2Operation_Type(t *testing.T) {
	op := &Vote2Operation{
		Voter:      "voter",
		Author:     "author",
		Permlink:   "permlink",
		RShares:    "1000000",
		Extensions: []interface{}{},
	}

	if op.Type() != TypeVote2 {
		t.Errorf("expected TypeVote2, got %v", op.Type())
	}
}

func TestCreateProposalOperation_Type(t *testing.T) {
	op := &CreateProposalOperation{
		Creator:    "creator",
		Receiver:   "receiver",
		DailyPay:   "1.000 STEEM",
		Subject:    "subject",
		Permlink:   "permlink",
		Extensions: []interface{}{},
	}

	if op.Type() != TypeCreateProposal {
		t.Errorf("expected TypeCreateProposal, got %v", op.Type())
	}
}

func TestUpdateProposalVotesOperation_Type(t *testing.T) {
	op := &UpdateProposalVotesOperation{
		Voter:       "voter",
		ProposalIDs: []uint64{1, 2, 3},
		Approve:     true,
		Extensions:  []interface{}{},
	}

	if op.Type() != TypeUpdateProposalVotes {
		t.Errorf("expected TypeUpdateProposalVotes, got %v", op.Type())
	}
}

func TestRemoveProposalOperation_Type(t *testing.T) {
	op := &RemoveProposalOperation{
		ProposalOwner: "owner",
		ProposalIDs:   []uint64{1, 2},
		Extensions:    []interface{}{},
	}

	if op.Type() != TypeRemoveProposal {
		t.Errorf("expected TypeRemoveProposal, got %v", op.Type())
	}
}

func TestClaimRewardBalance2Operation_Type(t *testing.T) {
	op := &ClaimRewardBalance2Operation{
		Account:      "account",
		RewardTokens: []interface{}{},
		Extensions:   []interface{}{},
	}

	if op.Type() != TypeClaimRewardBalance2 {
		t.Errorf("expected TypeClaimRewardBalance2, got %v", op.Type())
	}
}

func TestChainProperties_UnmarshalJSON_StringFormat(t *testing.T) {
	// Test old format: string
	jsonData := `{"account_creation_fee":"0.100 STEEM","maximum_block_size":131072,"sbd_interest_rate":1000}`

	var cp ChainProperties
	if err := json.Unmarshal([]byte(jsonData), &cp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if cp.AccountCreationFee != "0.100 STEEM" {
		t.Errorf("expected AccountCreationFee to be '0.100 STEEM', got '%s'", cp.AccountCreationFee)
	}
	if cp.MaximumBlockSize != 131072 {
		t.Errorf("expected MaximumBlockSize to be 131072, got %d", cp.MaximumBlockSize)
	}
	if cp.SBDInterestRate != 1000 {
		t.Errorf("expected SBDInterestRate to be 1000, got %d", cp.SBDInterestRate)
	}
}

func TestChainProperties_UnmarshalJSON_ObjectFormat(t *testing.T) {
	// Test new format: object with nai (from the actual error case)
	jsonData := `{"account_creation_fee":{"amount":"100000","nai":"@@000000021","precision":3},"maximum_block_size":131072,"sbd_interest_rate":1000}`

	var cp ChainProperties
	if err := json.Unmarshal([]byte(jsonData), &cp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Should be converted to string format: "100.000 STEEM"
	expected := "100.000 STEEM"
	if cp.AccountCreationFee != expected {
		t.Errorf("expected AccountCreationFee to be '%s', got '%s'", expected, cp.AccountCreationFee)
	}
	if cp.MaximumBlockSize != 131072 {
		t.Errorf("expected MaximumBlockSize to be 131072, got %d", cp.MaximumBlockSize)
	}
	if cp.SBDInterestRate != 1000 {
		t.Errorf("expected SBDInterestRate to be 1000, got %d", cp.SBDInterestRate)
	}
}

func TestChainProperties_UnmarshalJSON_POWOperation(t *testing.T) {
	// Test the actual pow_operation format from the error
	jsonData := `{"worker_account":"nxt4","block_id":"0000044666219088eff80258e4d2c73523a5203c","nonce":427,"work":{"input":"8afebe79fb50fab989ca5a5bd8ebdbbbab838e8e8dc8bb6386889bf6c2344bc7","signature":"1f6dad80034431283996e4a4f95d5130423cffc6d18e7d7ecb89345fe1be7931320c634aa945124c622ca99ab52a3358def3118ca5d12bf3f54d82a539795b707a","work":"002495e36694a3733737138bffaacc4cf425ca868b02214323f70f996934d2c5","worker":"STM5gzvDurFRmVUUs38TDtTtGVAEz8TcWMt4xLVbxwP2PP8b9q7P4"},"props":{"account_creation_fee":{"amount":"100000","nai":"@@000000021","precision":3},"maximum_block_size":131072,"sbd_interest_rate":1000}}`

	var op POWOperation
	if err := json.Unmarshal([]byte(jsonData), &op); err != nil {
		t.Fatalf("failed to unmarshal pow_operation: %v", err)
	}

	if op.WorkerAccount != "nxt4" {
		t.Errorf("expected WorkerAccount to be 'nxt4', got '%s'", op.WorkerAccount)
	}
	if op.Props == nil {
		t.Fatal("expected Props to be non-nil")
	}
	if op.Props.AccountCreationFee != "100.000 STEEM" {
		t.Errorf("expected AccountCreationFee to be '100.000 STEEM', got '%s'", op.Props.AccountCreationFee)
	}
}
