package protocol

import (
	"bytes"
	"strings"
	"testing"

	"github.com/steemit/steemutil/encoder"
)

func encodeOpErr(op interface{}) error {
	var buf bytes.Buffer
	enc := encoder.NewEncoder(&buf)
	return enc.Encode(op)
}

// TestPOWOperationSerialization locks the pow_operation wire format to the
// C++ chain reference (FC_REFLECT(steem::protocol::pow_operation,
// (worker_account)(block_id)(nonce)(work)(props)) and FC_REFLECT(pow,
// (worker)(input)(signature)(work))):
//
//	worker_account  string          (varint length + bytes)
//	block_id        block_id_type   20 raw bytes (fc::ripemd160)
//	nonce           uint64          8 bytes little-endian, no presence byte
//	work.worker     public_key_type 33 raw bytes
//	work.input      digest_type     32 raw bytes
//	work.signature  compact_sig     65 raw bytes
//	work.work       digest_type     32 raw bytes
//	props           legacy_chain_properties (fee asset, uint32, uint16)
//
// The worker's 33-byte binary was derived independently of this package via
// base58 + ripemd160-checksum decoding of the STM string.
func TestPOWOperationSerialization(t *testing.T) {
	nonce := UInt64(427)
	op := &POWOperation{
		WorkerAccount: "nxt4",
		BlockID:       "0000044666219088eff80258e4d2c73523a5203c",
		Nonce:         &nonce,
		Work: &POW{
			Worker:    "STM5gzvDurFRmVUUs38TDtTtGVAEz8TcWMt4xLVbxwP2PP8b9q7P4",
			Input:     "8afebe79fb50fab989ca5a5bd8ebdbbbab838e8e8dc8bb6386889bf6c2344bc7",
			Signature: "1f6dad80034431283996e4a4f95d5130423cffc6d18e7d7ecb89345fe1be7931320c634aa945124c622ca99ab52a3358def3118ca5d12bf3f54d82a539795b707a",
			Work:      "002495e36694a3733737138bffaacc4cf425ca868b02214323f70f996934d2c5",
		},
		Props: &ChainProperties{
			AccountCreationFee: "0.200 STEEM",
			MaximumBlockSize:   131072,
			SBDInterestRate:    1000,
		},
	}
	assertHex(t, "POWOperation", encodeOp(t, op),
		"0e"+ // op code 14 (pow_operation in the C++ static_variant)
			"046e787434"+ // worker_account "nxt4"
			"0000044666219088eff80258e4d2c73523a5203c"+ // block_id, 20 raw bytes
			"ab01000000000000"+ // nonce 427, uint64 LE
			"02699b26170706662df8bcdc97a504b5cca8b2f3a3ee14f3eeececd264b8491d4f"+ // worker, 33 raw bytes
			"8afebe79fb50fab989ca5a5bd8ebdbbbab838e8e8dc8bb6386889bf6c2344bc7"+ // input, 32 raw bytes
			"1f6dad80034431283996e4a4f95d5130423cffc6d18e7d7ecb89345fe1be7931320c634aa945124c622ca99ab52a3358def3118ca5d12bf3f54d82a539795b707a"+ // signature, 65 raw bytes
			"002495e36694a3733737138bffaacc4cf425ca868b02214323f70f996934d2c5"+ // work, 32 raw bytes
			"c80000000000000003535445454d0000"+ // props.account_creation_fee "0.200 STEEM"
			"00000200"+ // props.maximum_block_size 131072
			"e803") // props.sbd_interest_rate 1000
}

// TestPOWWireFormatRejectsBadHex verifies the hexbytes tags enforce the
// fixed sizes from the C++ types (block_id 20, input/work 32, signature 65).
func TestPOWWireFormatRejectsBadHex(t *testing.T) {
	goodWork := &POW{
		Worker:    "STM5gzvDurFRmVUUs38TDtTtGVAEz8TcWMt4xLVbxwP2PP8b9q7P4",
		Input:     "8afebe79fb50fab989ca5a5bd8ebdbbbab838e8e8dc8bb6386889bf6c2344bc7",
		Signature: "1f6dad80034431283996e4a4f95d5130423cffc6d18e7d7ecb89345fe1be7931320c634aa945124c622ca99ab52a3358def3118ca5d12bf3f54d82a539795b707a",
		Work:      "002495e36694a3733737138bffaacc4cf425ca868b02214323f70f996934d2c5",
	}
	nonce := UInt64(427)
	newOp := func() *POWOperation {
		return &POWOperation{
			WorkerAccount: "nxt4",
			BlockID:       "0000044666219088eff80258e4d2c73523a5203c",
			Nonce:         &nonce,
			Work: &POW{
				Worker: goodWork.Worker, Input: goodWork.Input,
				Signature: goodWork.Signature, Work: goodWork.Work,
			},
			Props: &ChainProperties{AccountCreationFee: "0.200 STEEM"},
		}
	}

	cases := []struct {
		name string
		mut  func(*POWOperation)
		want string
	}{
		{"block_id too short", func(op *POWOperation) { op.BlockID = "00" }, "exactly 20 bytes"},
		{"input wrong size", func(op *POWOperation) { op.Work.Input = "0102" }, "exactly 32 bytes"},
		{"signature not hex", func(op *POWOperation) { op.Work.Signature = "zz" }, "failed to decode hex"},
		{"work wrong size", func(op *POWOperation) { op.Work.Work = "00" }, "exactly 32 bytes"},
		{"bad worker key", func(op *POWOperation) { op.Work.Worker = "STM5invalid" }, "too short"},
	}
	for _, tc := range cases {
		op := newOp()
		tc.mut(op)
		err := encodeOpErr(op)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: expected error containing %q, got %v", tc.name, tc.want, err)
		}
	}
}

// TestCustomBinarySerialization locks the custom_binary_operation wire format
// to FC_REFLECT(steem::protocol::custom_binary_operation,
// (required_owner_auths)(required_active_auths)(required_posting_auths)
// (required_auths)(id)(data)): the four auth containers must appear even
// when empty. steem-js omits them, so its fixture for this op is not golden.
func TestCustomBinarySerialization(t *testing.T) {
	op := &CustomBinaryOperation{
		ID:        "test",
		DataBytes: "01020304",
	}
	assertHex(t, "CustomBinary empty auths", encodeOp(t, op),
		"23"+ // op code 35
			"00"+ // required_owner_auths: empty
			"00"+ // required_active_auths: empty
			"00"+ // required_posting_auths: empty
			"00"+ // required_auths: empty
			"0474657374"+ // id "test"
			"0401020304") // data

	opAuths := &CustomBinaryOperation{
		// flat_set: sorted on the wire regardless of input order.
		RequiredOwnerAuths:   []string{"zoe", "alice"},
		RequiredActiveAuths:  []string{"bob"},
		RequiredPostingAuths: []string{},
		RequiredAuths: []*Authority{
			{WeightThreshold: 2, KeyAuths: StringInt64Map{
				"STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA": 2,
			}},
		},
		ID:        "x",
		DataBytes: "",
	}
	assertHex(t, "CustomBinary with auths", encodeOp(t, opAuths),
		"23"+
			"02"+ // owner auths: 2 entries
			"05616c696365"+ // "alice" (sorted first)
			"037a6f65"+ // "zoe"
			"01"+ // active auths: 1 entry
			"03626f62"+ // "bob"
			"00"+ // posting auths empty
			"01"+ // one authority
			"02000000"+ // weight_threshold 2 (uint32 LE)
			"00"+ // account_auths empty
			"01"+ // one key_auth
			"03fdf4907810a9f5d9462a1ae09feee5ab205d32798b0ffcc379442021f84c5bbf"+ // 33-byte key
			"0200"+ // weight 2 (uint16 LE)
			"0178"+ // id "x"
			"00") // data empty
}

// TestAuthorityWeightRange verifies out-of-range weights (weight_type is
// uint16 on the chain) and unparseable keys are rejected instead of being
// silently truncated or reordered.
func TestAuthorityWeightRange(t *testing.T) {
	for _, tc := range []struct {
		name string
		auth *Authority
	}{
		{
			name: "key_auths weight too large",
			auth: &Authority{WeightThreshold: 1, KeyAuths: StringInt64Map{
				"STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA": 1 << 20,
			}},
		},
		{
			name: "key_auths weight negative",
			auth: &Authority{WeightThreshold: 1, KeyAuths: StringInt64Map{
				"STM8m5UgaFAAYQRuaNejYdS8FVLVp9Ss3K1qAVk5de6F8s3HnVbvA": -1,
			}},
		},
		{
			name: "account_auths weight too large",
			auth: &Authority{WeightThreshold: 1, AccountAuths: StringInt64Map{"alice": 1 << 20}},
		},
		{
			name: "invalid key_auths key",
			auth: &Authority{WeightThreshold: 1, KeyAuths: StringInt64Map{"STMnotakey": 1}},
		},
	} {
		if err := encodeOpErr(tc.auth); err == nil {
			t.Errorf("%s: expected error, got nil", tc.name)
		}
	}
}
