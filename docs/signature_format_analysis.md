# Signature Format Analysis

## Summary

After analyzing the signature format used by `btcec.SignCompact` and comparing it with `steem-js` and Steem C++ implementations, we found:

### Key Findings

1. **Recovery Parameter Format**: ✅ Correct
   - `btcec.SignCompact` returns recovery ID in format: `27 + 4 + recovery_param` for compressed keys (31-34)
   - This matches Steem's expected format
   - Steem C++ extracts recovery param as: `(recoveryID - 27) & 3`

2. **Canonical Signature**: ✅ Correct
   - `btcec.SignCompact` automatically generates canonical signatures
   - Test confirms signatures pass `is_fc_canonical` check:
     - `!(c.data[1] & 0x80)` - r[0] doesn't have high bit set
     - `!(c.data[1] == 0 && !(c.data[2] & 0x80))` - if r[0] == 0, r[1] must have high bit set
     - `!(c.data[33] & 0x80)` - s[0] doesn't have high bit set
     - `!(c.data[33] == 0 && !(c.data[34] & 0x80))` - if s[0] == 0, s[1] must have high bit set

3. **Public Key Recovery**: ✅ Correct
   - Recovered public key matches expected public key
   - Compression flag is correctly set

### Comparison with steem-js

**old-steem-js** (signature.js:88-90):
```javascript
i = ecdsa.calcPubKeyRecoveryParam(curve, e, ecsignature, private_key.toPublicKey().Q);
i += 4;  // compressed
i += 27; // compact
// Final: 27 + 4 + (0-3) = 31-34
```

**steem-js** (signature.ts:88-89):
```typescript
const i = calcPubKeyRecoveryParam(secp256k1, new BN(buf_sha256), ecsignature, privKey.toPublic().Q!);
return new Signature(ecsignature.r, ecsignature.s, i + 27);
// Note: toCompact() method adds 4 for compressed keys
```

**Go (steemutil)**:
```go
sig, err := ecdsa.SignCompact(privKey.Raw.PrivKey, digest, true)
// Returns recovery ID in format: 27 + 4 + recovery_param (31-34) for compressed keys
```

### Steem C++ Implementation

**elliptic_secp256k1.cpp:164**:
```cpp
secp256k1_ecdsa_recover_compact(
    detail::_get_context(),
    (unsigned char*) digest.data(),
    (unsigned char*) c.begin() + 1,  // Skip recovery ID byte
    (unsigned char*) my->_key.begin(),
    (int*) &pk_len,
    1,  // compressed flag
    (*c.begin() - 27) & 3  // Extract recovery param: (recoveryID - 27) & 3
);
```

**elliptic_impl_pub.cpp:343-348**:
```cpp
if (nV >= 31) {
    EC_KEY_set_conv_form(my->_key, POINT_CONVERSION_COMPRESSED);
    nV -= 4;  // Remove compressed flag
}
```

### Conclusion

The `btcec.SignCompact` implementation is **correct** and matches Steem's expected format. The signature format, recovery parameter encoding, and canonical signature requirements are all properly handled.

If signature verification fails on the blockchain, the issue is likely:
1. **Not related to signature format** - the format is correct
2. **May be related to transaction serialization** - ensure transaction bytes match exactly
3. **May be related to digest calculation** - ensure chain ID and transaction bytes are concatenated correctly
4. **May be related to key matching** - ensure the private key matches the account's posting key

### Test Results

```
✅ Recovery ID 32 is in expected range for compressed keys (31-34)
✅ Signature is canonical (fc_canonical format)
✅ Recovered public key matches!
```

