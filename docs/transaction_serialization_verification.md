# Transaction Serialization and Digest Calculation Verification

## Summary

After thorough analysis and testing, we have verified that:

### ✅ Transaction Serialization: Correct

1. **Signatures are excluded**: Transaction serialization correctly excludes signatures when calculating the digest
   - Test: `TestTransactionSerializationExcludesSignatures` confirms that transactions with and without signatures serialize to the same bytes
   - This matches Steem C++ behavior where `transaction::sig_digest()` serializes the `transaction` object (not `signed_transaction`)

2. **Serialization format**: The transaction serialization format matches Steem protocol:
   - `ref_block_num` (uint16)
   - `ref_block_prefix` (uint32)
   - `expiration` (time_point_sec - uint32)
   - `operations` (varint length + operations)
   - `extensions` (varint length + extensions)
   - **No signatures** in digest calculation

### ✅ Digest Calculation: Correct

1. **Chain ID format**: Chain ID is correctly decoded from hex string (64 hex chars = 32 bytes)
   - Steem C++: `fc::sha256` (32 bytes) serialized via `fc::raw::pack`
   - Go: `hex.DecodeString(chain.ID)` produces 32 bytes ✅

2. **Digest formula**: `sha256(chain_id + serialized_transaction)`
   - Steem C++: `fc::raw::pack(enc, chain_id); fc::raw::pack(enc, *this); return enc.result();`
   - Go: `msgBuffer.Write(rawChainID); msgBuffer.Write(rawTx); sha256.Sum256(msgBuffer.Bytes())` ✅

3. **Test verification**: `TestTransaction_Digest` passes with expected digest value

### Comparison with steem-js

**steem-js** (auth/index.ts:136-141):
```typescript
const chainId = (getConfig().get('chain_id') as string | undefined) || '';
const cid = Buffer.from(chainId, 'hex');
const buf = transaction.toBuffer(trx as unknown);
const sig = Signature.signBuffer(Buffer.concat([cid, buf]), key);
```

**Go (steemutil)**:
```go
rawChainID, err := hex.DecodeString(chain.ID)
msgBuffer.Write(rawChainID)
rawTx, err := tx.Serialize()
msgBuffer.Write(rawTx)
digest := sha256.Sum256(msgBuffer.Bytes())
```

Both implementations:
1. Decode chain ID from hex string to bytes
2. Serialize transaction (without signatures)
3. Concatenate chain_id + transaction
4. Calculate SHA256 digest

### Test Results

```
✅ Transaction serialization correctly excludes signatures
✅ Digests match (signatures correctly excluded)
✅ Digest calculation matches expected format
```

### Conclusion

The transaction serialization and digest calculation are **correct** and match both Steem C++ and steem-js implementations. The issue with signature verification on the blockchain is **not** related to:
- Transaction serialization format
- Digest calculation method
- Chain ID encoding

The problem is likely related to:
1. **Signature recovery**: The signature format is correct, but the blockchain may be using a different recovery method
2. **Key matching**: The private key may not match the account's posting key (though we've verified this)
3. **Transaction preparation**: The ref_block_num and ref_block_prefix may need to be set correctly before signing

### Next Steps

To debug further, we should:
1. Compare the actual signed transaction bytes with steem-js output
2. Verify the signature recovery process on the blockchain side
3. Check if there are any differences in how the blockchain validates signatures

