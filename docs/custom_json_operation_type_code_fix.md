# CustomJSONOperation Operation Type Code Fix

## Problem Summary

The `CustomJSONOperation.MarshalTransaction` method was missing the operation type code encoding, causing transactions to be serialized incorrectly. This document explains why this bug wasn't caught earlier.

## Why Other Operations Worked

### Encoder Logic Flow

The `encoder.Encode()` method follows this priority:

1. **Check for `MarshalTransaction` method** (line 50):
   ```go
   if marshaller, ok := v.(TransactionMarshaller); ok {
       return marshaller.MarshalTransaction(encoder)
   }
   ```

2. **Check if it's an Operation interface** (line 55):
   ```go
   if op := encoder.checkOperation(v); op != nil {
       return encoder.encodeOperation(op)  // Automatically encodes type code
   }
   ```

### Two Paths for Operations

**Path 1: Operations WITH `MarshalTransaction` method**
- `VoteOperation`, `CommentOperation`, `CustomJSONOperation` (before fix)
- These methods are called directly
- **They must manually encode the operation type code**
- Example from `VoteOperation`:
  ```go
  func (op *VoteOperation) MarshalTransaction(encoderObj *encoder.Encoder) error {
      enc.EncodeUVarint(uint64(TypeVote.Code()))  // ✅ Encodes type code
      enc.Encode(op.Voter)
      // ...
  }
  ```

**Path 2: Operations WITHOUT `MarshalTransaction` method**
- Operations that don't implement `MarshalTransaction`
- Go through `encodeOperation()` which **automatically encodes the type code**:
  ```go
  func (encoder *Encoder) encodeOperation(op *operationInterface) error {
      // Encode the operation type code
      code := op.getTypeCode()
      encoder.EncodeUVarint(code)  // ✅ Automatically encoded
      // Then encode operation data
      return encoder.encodeByReflection(opData)
  }
  ```

### Why CustomJSONOperation Failed

`CustomJSONOperation` had a `MarshalTransaction` method but **didn't encode the type code**:

```go
// BEFORE FIX (WRONG)
func (op *CustomJSONOperation) MarshalTransaction(encoderObj *encoder.Encoder) error {
    // ❌ Missing: encoderObj.EncodeUVarint(uint64(op.Type().Code()))
    encoderObj.EncodeUVarint(uint64(len(requiredAuths)))
    // ... rest of encoding
}
```

This caused the operation type code (18 for `custom_json`) to be missing, making the blockchain interpret it as operation type 0 (`vote`), which would fail validation.

## Why Unit Tests Didn't Catch This

### Test Implementation

The unit test in `operation_custom_json_test.go` directly calls `MarshalTransaction`:

```go
err := tt.op.MarshalTransaction(enc)  // Direct call, not through encoder.Encode()
```

### What the Test Validated

The test compared the output of `MarshalTransaction` with expected hex values from `old-steem-js`:

```go
expectedHex2 := "000105616c69636506666f6c6c6f77315b22666f6c6c6f77222c7b22666f6c6c6f776572223a22616c696365222c22666f6c6c6f77696e67223a22626f62227d5d"
```

### Why It Matched

The expected hex from `old-steem-js` was generated using:

```javascript
const buf = custom_json.toBuffer(testCase.op);  // Only operation DATA, not full operation
```

This `custom_json.toBuffer()` method in steem-js **only serializes the operation data**, not the full operation tuple `[type_code, data]`. It's equivalent to Go's `MarshalTransaction` method.

So the test was correctly validating that the **operation data** serialization matched, but it wasn't testing the **full operation** serialization (which includes the type code).

### The Missing Test

The test should have also tested the full operation encoding path:

```go
// Test full operation encoding (what actually happens in transactions)
err := enc.Encode(tt.op)  // Goes through encoder.Encode() → encodeOperation()
```

This would have caught the missing type code.

## Fix

Added operation type code encoding to `CustomJSONOperation.MarshalTransaction`:

```go
// AFTER FIX (CORRECT)
func (op *CustomJSONOperation) MarshalTransaction(encoderObj *encoder.Encoder) error {
    // ✅ Encode operation type code first
    if err := encoderObj.EncodeUVarint(uint64(op.Type().Code())); err != nil {
        return errors.Wrap(err, "failed to encode operation type code")
    }
    // ... rest of encoding
}
```

## Lessons Learned

1. **Operations with `MarshalTransaction` must encode type code manually**
2. **Unit tests should test both paths**:
   - Direct `MarshalTransaction` call (operation data only)
   - Through `encoder.Encode()` (full operation with type code)
3. **When comparing with steem-js**, ensure you're comparing the same level:
   - `custom_json.toBuffer()` = operation data only
   - `operation.toBuffer([type, data])` = full operation with type code

