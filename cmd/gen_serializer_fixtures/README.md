# Serializer Fixtures Generator

This tool generates JSON fixtures for cross-language serialization verification between Go SDK (steemutil) and JavaScript SDK (steem-js).

## Purpose

When implementing or modifying transaction serialization logic, it's critical to ensure that both Go and JavaScript SDKs produce identical binary output. This tool generates test fixtures that can be used by steem-js tests to verify serialization compatibility.

## Usage

```bash
# Run the generator
go run ./cmd/gen_serializer_fixtures/

# Output files will be written to /tmp/steem-serializer-fixtures/
ls /tmp/steem-serializer-fixtures/
```

## Output Format

Each fixture is a JSON file with the following structure:

```json
{
  "name": "transfer_basic",
  "tx": {
    "ref_block_num": 19297,
    "ref_block_prefix": 1608085982,
    "expiration": "2016-03-23T22:41:21",
    "operations": [
      ["transfer", {
        "from": "alice",
        "to": "bob",
        "amount": "1.000 STEEM",
        "memo": "hello"
      }]
    ],
    "extensions": []
  },
  "expected_hex": "..."
}
```

- `name`: Fixture identifier
- `tx`: Transaction in JSON format (compatible with steem-js `serializeTransaction` input)
- `expected_hex`: Expected hexadecimal output from Go SDK serialization

## Verification Workflow

1. **Generate fixtures** with this tool:
   ```bash
   go run ./cmd/gen_serializer_fixtures/
   ```

2. **Copy fixtures** to steem-js test directory

3. **Run steem-js tests** that:
   - Parse the `tx` JSON
   - Serialize using steem-js serializer
   - Compare output with `expected_hex`

4. **If hex doesn't match**, there's a serialization discrepancy that needs investigation

## Available Fixtures

| Fixture Name | Operation Type | Description |
|--------------|----------------|-------------|
| `transfer_basic` | transfer | Basic token transfer |
| `account_create_with_delegation_basic` | account_create_with_delegation | Account creation with delegation |
| `withdraw_vesting_basic` | withdraw_vesting | Vesting withdrawal |
| `limit_order_create2_basic` | limit_order_create2 | Limit order creation |
| `escrow_transfer_basic` | escrow_transfer | Escrow transfer |
| `claim_reward_balance_basic` | claim_reward_balance | Claim rewards |
| `pow_basic` | pow | Proof of work |
| `witness_update_basic` | witness_update | Witness update |
| `custom_json_basic` | custom_json | Custom JSON operation |
| `custom_binary_basic` | custom_binary | Custom binary operation |
| `comment_options_beneficiaries` | comment_options | Comment options with single beneficiary |
| `comment_options_beneficiaries_multiple` | comment_options | Comment options with multiple beneficiaries |

## Adding New Fixtures

To add a new fixture, edit `main.go` and append to the `fixtures` slice:

```go
fixtures = append(fixtures, makeTxFixture(
    "your_fixture_name",
    "operation_type",
    &protocol.YourOperation{
        // operation fields
    },
))
```

## Notes

- All fixtures use fixed `ref_block_num`, `ref_block_prefix`, and `expiration` values for reproducibility
- Beneficiaries in `comment_options` must be sorted alphabetically by account name (Steem protocol requirement)
- Asset fields (Amount, Fee, etc.) are automatically encoded in binary format
