package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/protocol"
	"github.com/steemit/steemutil/transaction"
)

type Fixture struct {
	Name        string                 `json:"name"`
	Tx          map[string]interface{} `json:"tx"`
	ExpectedHex string                 `json:"expected_hex"`
}

func mustTime(s string) *protocol.Time {
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &protocol.Time{Time: &parsed}
}

func makeTxFixture(name string, opName string, op protocol.Operation) Fixture {
	// Use fixed header so JS side can construct identical tx objects.
	refBlockNum := protocol.UInt16(19297)
	refBlockPrefix := protocol.UInt32(1608085982)
	expiration := mustTime("2016-03-23T22:41:21Z")

	tx := &transaction.Transaction{
		RefBlockNum:    refBlockNum,
		RefBlockPrefix: refBlockPrefix,
		Expiration:     expiration,
		Operations:     protocol.Operations{op},
		Extensions:     []interface{}{},
	}

	var buf bytes.Buffer
	enc := encoder.NewEncoder(&buf)
	if err := enc.Encode(tx); err != nil {
		panic(fmt.Errorf("failed to encode tx for %s: %w", name, err))
	}

	expectedHex := hex.EncodeToString(buf.Bytes())

	// Build JSON tx payload matching steem-js serializeTransaction input.
	// operations must be [[opName, opData]] where opData uses JSON field names.
	opJSON, err := json.Marshal(op)
	if err != nil {
		panic(fmt.Errorf("failed to marshal op %s to JSON: %w", name, err))
	}
	var opData map[string]interface{}
	if err := json.Unmarshal(opJSON, &opData); err != nil {
		panic(fmt.Errorf("failed to unmarshal op %s JSON: %w", name, err))
	}

	txJSON := map[string]interface{}{
		"ref_block_num":    uint16(refBlockNum),
		"ref_block_prefix": uint32(refBlockPrefix),
		"expiration":       "2016-03-23T22:41:21", // without Z, matches steem-js tests
		"operations": []interface{}{
			[]interface{}{opName, opData},
		},
		"extensions": []interface{}{},
	}

	return Fixture{
		Name:        name,
		Tx:          txJSON,
		ExpectedHex: expectedHex,
	}
}

func writeFixture(dir string, f Fixture) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(fmt.Errorf("failed to create dir %s: %w", dir, err))
	}
	path := filepath.Join(dir, fmt.Sprintf("%s.json", f.Name))

	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal fixture %s: %w", f.Name, err))
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		panic(fmt.Errorf("failed to write fixture %s: %w", path, err))
	}
}

func main() {
	outDir := "/tmp/steem-serializer-fixtures"

	// Representative operations from each functional group.
	var fixtures []Fixture

	// transfer
	fixtures = append(fixtures, makeTxFixture(
		"transfer_basic",
		"transfer",
		&protocol.TransferOperation{
			From:   "alice",
			To:     "bob",
			Amount: "1.000 STEEM",
			Memo:   "hello",
		},
	))

	// account_create_with_delegation
	fixtures = append(fixtures, makeTxFixture(
		"account_create_with_delegation_basic",
		"account_create_with_delegation",
		&protocol.AccountCreateWithDelegationOperation{
			Fee:        "0.000 STEEM",
			Delegation: "0.000000 VESTS",
			Creator:    "initminer",
			NewAccountName: "alice",
			Owner: &protocol.Authority{
				WeightThreshold: 1,
				AccountAuths:    protocol.StringInt64Map{},
				KeyAuths:        protocol.StringInt64Map{"STM8k1f8fvHxLrCTqMdRUJcK2rCE3y7SQBb8PremyadWvVWMeedZy": 1},
			},
			Active: &protocol.Authority{
				WeightThreshold: 1,
				AccountAuths:    protocol.StringInt64Map{},
				KeyAuths:        protocol.StringInt64Map{"STM8k1f8fvHxLrCTqMdRUJcK2rCE3y7SQBb8PremyadWvVWMeedZy": 1},
			},
			Posting: &protocol.Authority{
				WeightThreshold: 1,
				AccountAuths:    protocol.StringInt64Map{},
				KeyAuths:        protocol.StringInt64Map{"STM8k1f8fvHxLrCTqMdRUJcK2rCE3y7SQBb8PremyadWvVWMeedZy": 1},
			},
			MemoKey:      "STM6ppNVEFmvBW4jEkzxXnGKuKuwYjMUrhz2WX1kHeGSchGdWJEDQ",
			JsonMetadata: "",
			Extensions:   []interface{}{},
		},
	))

	// withdraw_vesting
	fixtures = append(fixtures, makeTxFixture(
		"withdraw_vesting_basic",
		"withdraw_vesting",
		&protocol.WithdrawVestingOperation{
			Account:       "alice",
			VestingShares: "10.000000 VESTS",
		},
	))

	// limit_order_create2
	fixtures = append(fixtures, makeTxFixture(
		"limit_order_create2_basic",
		"limit_order_create2",
		&protocol.LimitOrderCreate2Operation{
			Owner:        "alice",
			OrderID:      123,
			AmountToSell: "1.000 STEEM",
			ExchangeRate: struct {
				Base  string `json:"base"`
				Quote string `json:"quote"`
			}{
				Base:  "1.000 STEEM",
				Quote: "1.000 SBD",
			},
			FillOrKill: false,
			Expiration: mustTime("2016-03-24T16:05:00Z"),
		},
	))

	// escrow_transfer
	fixtures = append(fixtures, makeTxFixture(
		"escrow_transfer_basic",
		"escrow_transfer",
		&protocol.EscrowTransferOperation{
			From:        "alice",
			To:          "bob",
			SBDAmount:   "1.000 SBD",
			SteemAmount: "0.000 STEEM",
			EscrowID:    1,
			Agent:       "carol",
			Fee:         "0.001 STEEM",
			JsonMeta:    "",
			RatificationDeadline: mustTime("2016-03-24T16:10:00Z"),
			EscrowExpiration:     mustTime("2016-03-25T16:10:00Z"),
		},
	))

	// claim_reward_balance
	fixtures = append(fixtures, makeTxFixture(
		"claim_reward_balance_basic",
		"claim_reward_balance",
		&protocol.ClaimRewardBalanceOperation{
			Account:     "alice",
			RewardSteem: "1.000 STEEM",
			RewardSBD:   "0.500 SBD",
			RewardVests: "10.000000 VESTS",
		},
	))

	// pow
	nonce := protocol.UInt64(427)
	fixtures = append(fixtures, makeTxFixture(
		"pow_basic",
		"pow",
		&protocol.POWOperation{
			WorkerAccount: "nxt4",
			BlockID:       "0000044666219088eff80258e4d2c73523a5203c",
			Nonce:         &nonce,
			Work: &protocol.POW{
				Worker:    "STM5gzvDurFRmVUUs38TDtTtGVAEz8TcWMt4xLVbxwP2PP8b9q7P4",
				Input:     "8afebe79fb50fab989ca5a5bd8ebdbbbab838e8e8dc8bb6386889bf6c2344bc7",
				Signature: "1f6dad80034431283996e4a4f95d5130423cffc6d18e7d7ecb89345fe1be7931320c634aa945124c622ca99ab52a3358def3118ca5d12bf3f54d82a539795b707a",
				Work:      "002495e36694a3733737138bffaacc4cf425ca868b02214323f70f996934d2c5",
			},
			Props: &protocol.ChainProperties{
				AccountCreationFee: "0.200 STEEM",
				MaximumBlockSize:   131072,
				SBDInterestRate:    1000,
			},
		},
	))

	// witness_update
	fixtures = append(fixtures, makeTxFixture(
		"witness_update_basic",
		"witness_update",
		&protocol.WitnessUpdateOperation{
			Owner:           "initminer",
			URL:             "https://example.com",
			BlockSigningKey: "STM5gzvDurFRmVUUs38TDtTtGVAEz8TcWMt4xLVbxwP2PP8b9q7P4",
			Props: &protocol.ChainProperties{
				AccountCreationFee: "0.200 STEEM",
				MaximumBlockSize:   65536,
				SBDInterestRate:    1000,
			},
			Fee: "0.000 STEEM",
		},
	))

	// custom_json
	fixtures = append(fixtures, makeTxFixture(
		"custom_json_basic",
		"custom_json",
		&protocol.CustomJSONOperation{
			RequiredAuths:        []string{"alice"},
			RequiredPostingAuths: []string{},
			ID:                   "follow",
			JSON:                 `{"foo":"bar"}`,
		},
	))

	// custom_binary
	fixtures = append(fixtures, makeTxFixture(
		"custom_binary_basic",
		"custom_binary",
		&protocol.CustomBinaryOperation{
			ID:        "test",
			DataBytes: "01020304",
		},
	))

	for _, f := range fixtures {
		writeFixture(outDir, f)
		fmt.Fprintf(os.Stderr, "wrote fixture %s\n", f.Name)
	}
}

