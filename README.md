# steemutil

> This is underconstruction. The methods will be change in the future !!!!

Package steemutil provides Steem blockchain-specific convenience functions and types for Go applications.

## Overview

steemutil is a comprehensive Go library that provides low-level utilities for interacting with the Steem blockchain. It includes cryptographic functions, transaction handling, protocol definitions, and RPC authentication capabilities.

## Features

- **🔐 RPC Authentication**: SignedCall support for authenticated API requests
- **🔑 Cryptographic Operations**: Key generation, signing, and verification
- **📦 Transaction Handling**: Transaction creation, signing, and serialization
- **🌐 Protocol Support**: Complete Steem protocol operation definitions
- **🛡️ Security**: Built-in replay protection and signature validation
- **⚡ Performance**: Optimized for high-performance applications

## Installation

```bash
go get github.com/steemit/steemutil
```

## Quick Start

### Basic Key Operations

```go
package main

import (
    "fmt"
    "github.com/steemit/steemutil/auth"
)

func main() {
    // Generate private key from account and password
    privateKey, err := auth.ToWif("username", "password", "active")
    if err != nil {
        panic(err)
    }
    
    // Convert to public key
    publicKey, err := auth.WifToPublic(privateKey)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Private Key: %s\n", privateKey)
    fmt.Printf("Public Key: %s\n", publicKey)
}
```

### SignedCall RPC Authentication

```go
package main

import (
    "fmt"
    "github.com/steemit/steemutil/rpc"
)

func main() {
    // Create RPC request
    request := &rpc.RpcRequest{
        Method: "condenser_api.get_accounts",
        Params: []interface{}{[]string{"username"}},
        ID:     1,
    }
    
    // Sign the request
    signedRequest, err := rpc.Sign(request, "username", []string{"private-key-wif"})
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Signed Request: %+v\n", signedRequest)
}
```

### Transaction Operations

```go
package main

import (
    "github.com/steemit/steemutil/protocol"
    "github.com/steemit/steemutil/transaction"
)

func main() {
    // Create a vote operation
    voteOp := &protocol.VoteOperation{
        Voter:    "voter",
        Author:   "author", 
        Permlink: "permlink",
        Weight:   10000,
    }
    
    // Create transaction
    tx := &transaction.Transaction{}
    tx.PushOperation(voteOp)
    
    // Sign transaction
    signedTx := transaction.NewSignedTransaction(tx)
    // ... add signing logic
}
```

## Package Structure

### Core Packages

- **`auth/`** - Authentication and key management utilities
- **`rpc/`** - RPC authentication and signed call support  
- **`protocol/`** - Steem protocol definitions and operations
- **`transaction/`** - Transaction creation and signing
- **`wif/`** - Wallet Import Format key handling
- **`encoder/`** - Binary serialization utilities

### Protocol Support

- **`protocol/api/`** - API method definitions
- **`protocol/broadcast/`** - Broadcast operation definitions
- **`consts/`** - Protocol constants and chain parameters
- **`jsonrpc2/`** - JSON-RPC client implementation

## RPC Authentication (SignedCall)

The `rpc` package provides full support for Steem's signed RPC calls, compatible with steem-js:

### Features

- **🔒 Cryptographic Signing**: ECDSA signature generation and verification
- **🛡️ Security**: Nonce generation, timestamp expiration (60s), replay protection
- **🔄 Compatibility**: 100% compatible with steem-js signedCall format
- **🔑 Multi-Key Support**: Sign with multiple private keys simultaneously
- **⏰ Time Validation**: Automatic signature expiration handling

### Security Features

- **Unique Nonces**: 8-byte random nonce for each request
- **Timestamp Expiration**: Signatures expire after 60 seconds
- **Cross-Protocol Protection**: Protocol-specific signing constant K
- **No Key Transmission**: Private keys never leave your application

### Example Usage

```go
// Sign a request
request := &rpc.RpcRequest{
    Method: "condenser_api.get_account_history",
    Params: []interface{}{"username", -1, 100},
    ID:     1,
}

signedRequest, err := rpc.Sign(request, "username", []string{"active-key-wif"})
if err != nil {
    return err
}

// Validate a signed request
params, err := rpc.Validate(signedRequest, rpc.DefaultVerifyFunc)
if err != nil {
    return err
}
```

## Cryptographic Operations

### Key Generation

```go
// Generate keys for multiple roles
roles := []string{"active", "posting", "owner", "memo"}
keys, err := auth.GetPrivateKeys("username", "password", roles)

// Generate single key
activeKey, err := auth.ToWif("username", "password", "active")
```

### Message Signing

```go
// Sign arbitrary messages
privateKey := &wif.PrivateKey{}
privateKey.FromWif("private-key-wif")

message := []byte("Hello, Steem!")
signature, err := privateKey.SignSha256(message)

// Verify signatures
publicKey := &wif.PublicKey{}
publicKey.FromWif("private-key-wif") // Derives public key
isValid := publicKey.VerifySha256(message, signature)
```

## Transaction Handling

### Creating Transactions

```go
// Create operations
voteOp := &protocol.VoteOperation{
    Voter:    "voter",
    Author:   "author",
    Permlink: "post-permlink", 
    Weight:   10000,
}

transferOp := &protocol.TransferOperation{
    From:   "sender",
    To:     "receiver",
    Amount: "1.000 STEEM",
    Memo:   "Transfer memo",
}

// Build transaction
tx := &transaction.Transaction{}
tx.PushOperation(voteOp)
tx.PushOperation(transferOp)

// Sign transaction
signedTx := transaction.NewSignedTransaction(tx)
err := signedTx.Sign(privateKeys, transaction.SteemChain)
```

### Supported Operations

The library supports all Steem protocol operations:

- **Content**: `comment`, `vote`, `delete_comment`
- **Financial**: `transfer`, `transfer_to_savings`, `claim_reward_balance`
- **Account**: `account_create`, `account_update`, `recover_account`
- **Witness**: `witness_update`, `account_witness_vote`
- **Market**: `limit_order_create`, `limit_order_cancel`
- **Custom**: `custom_json`, `custom_binary`
- **And many more...**

### Working with Operation Data

The `Operation` interface provides a `Data() any` method that returns the operation data. Understanding how to use this method is important for type-safe operation handling.

#### Return Types

The `Data()` method returns different types depending on the operation:

- **Known Operations**: Returns the operation struct itself (e.g., `*VoteOperation`, `*TransferOperation`)
- **Unknown Operations**: Returns `*json.RawMessage` for operations not recognized by the package

#### Best Practices

**1. Type Assertion Based on Operation Type**

Always check the operation type before performing type assertions:

```go
import (
    "encoding/json"
    "github.com/steemit/steemutil/protocol"
)

func processOperation(op protocol.Operation) {
    switch op.Type() {
    case protocol.TypeVote:
        // Type assertion is safe after checking Type()
        voteOp := op.Data().(*protocol.VoteOperation)
        fmt.Printf("Voter: %s, Author: %s\n", voteOp.Voter, voteOp.Author)
        
    case protocol.TypeTransfer:
        transferOp := op.Data().(*protocol.TransferOperation)
        fmt.Printf("From: %s, To: %s, Amount: %s\n", 
            transferOp.From, transferOp.To, transferOp.Amount)
        
    case protocol.TypeComment:
        commentOp := op.Data().(*protocol.CommentOperation)
        fmt.Printf("Author: %s, Title: %s\n", commentOp.Author, commentOp.Title)
        
    default:
        // Handle unknown operations
        if rawJSON, ok := op.Data().(*json.RawMessage); ok {
            fmt.Printf("Unknown operation type: %s\n", op.Type())
            fmt.Printf("Raw JSON: %s\n", string(*rawJSON))
        }
    }
}
```

**2. Safe Type Assertion with Error Handling**

Use type assertions with the two-value form for safer code:

```go
func safeProcessVote(op protocol.Operation) error {
    if op.Type() != protocol.TypeVote {
        return fmt.Errorf("expected vote operation, got %s", op.Type())
    }
    
    voteOp, ok := op.Data().(*protocol.VoteOperation)
    if !ok {
        return fmt.Errorf("failed to assert vote operation data")
    }
    
    // Use voteOp safely
    fmt.Printf("Processing vote: %s -> %s/%s\n", 
        voteOp.Voter, voteOp.Author, voteOp.Permlink)
    return nil
}
```

**3. Handling Unknown Operations**

For operations not recognized by the package, `Data()` returns `*json.RawMessage`:

```go
func handleUnknownOperation(op protocol.Operation) {
    if rawJSON, ok := op.Data().(*json.RawMessage); ok {
        // This is an unknown operation type
        var data map[string]any
        if err := json.Unmarshal(*rawJSON, &data); err == nil {
            fmt.Printf("Unknown operation data: %+v\n", data)
        }
    } else {
        // This is a known operation type
        fmt.Printf("Known operation: %s\n", op.Type())
    }
}
```

**4. Direct Operation Access**

For known operations, you can directly use the operation struct without calling `Data()`:

```go
// Instead of:
data := op.Data().(*protocol.VoteOperation)

// You can directly cast the operation:
if voteOp, ok := op.(*protocol.VoteOperation); ok {
    // Use voteOp directly
    fmt.Printf("Voter: %s\n", voteOp.Voter)
}
```

**5. Iterating Over Operations**

When processing multiple operations:

```go
func processOperations(ops protocol.Operations) {
    for _, op := range ops {
        switch op.Type() {
        case protocol.TypeVote:
            voteOp := op.Data().(*protocol.VoteOperation)
            processVote(voteOp)
            
        case protocol.TypeTransfer:
            transferOp := op.Data().(*protocol.TransferOperation)
            processTransfer(transferOp)
            
        default:
            // Log or handle unknown operations
            fmt.Printf("Unhandled operation type: %s\n", op.Type())
        }
    }
}
```

#### Important Notes

- **Type Safety**: Always check `op.Type()` before performing type assertions
- **Unknown Operations**: Use `*json.RawMessage` type assertion to handle unrecognized operations
- **Performance**: Direct operation casting (e.g., `op.(*VoteOperation)`) is more efficient than using `Data()` for known types
- **Compatibility**: The `any` return type allows the library to handle both known and unknown operation types flexibly

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./rpc -v
go test ./auth -v
go test ./protocol -v

# Run with coverage
go test ./... -cover
```

## Examples

See the [examples directory](../test-gosdk/examples/) for complete working examples:

- **SignedCall Authentication**: `examples/signed_call/`
- **Key Generation**: `examples/generate_keys/`
- **Transaction Broadcasting**: `examples/transfer/`
- **Vote Operations**: `examples/vote_post/`

## API Reference

### Authentication (`auth/`)

- `ToWif(name, password, role string) (string, error)` - Generate WIF from credentials
- `GetPrivateKeys(name, password string, roles []string) (map[string]string, error)` - Generate multiple keys
- `WifToPublic(wif string) (string, error)` - Convert WIF to public key
- `IsWif(wif string) bool` - Validate WIF format
- `Verify(name, password string, auths map[string]interface{}) (bool, error)` - Verify credentials

### RPC Authentication (`rpc/`)

- `Sign(request *RpcRequest, account string, keys []string) (*SignedRequest, error)` - Sign RPC request
- `Validate(request *SignedRequest, verifyFunc func(...) error) ([]interface{}, error)` - Validate signed request
- `SignRequest(method string, params []interface{}, id int, account, key string) (*SignedRequest, error)` - Convenience function

### Transaction (`transaction/`)

- `NewSignedTransaction(tx *Transaction) *SignedTransaction` - Create signed transaction
- `(tx *SignedTransaction) Sign(keys []*wif.PrivateKey, chain *Chain) error` - Sign transaction
- `(tx *SignedTransaction) Digest(chain *Chain) ([]byte, error)` - Calculate transaction digest
- `(tx *SignedTransaction) Serialize() ([]byte, error)` - Serialize transaction

### WIF Operations (`wif/`)

- `(pk *PrivateKey) FromWif(wif string) error` - Import from WIF
- `(pk *PrivateKey) ToWif() string` - Export to WIF
- `(pk *PrivateKey) SignSha256(message []byte) ([]byte, error)` - Sign message
- `(pk *PublicKey) VerifySha256(message, signature []byte) bool` - Verify signature

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- **[steemgosdk](https://github.com/steemit/steemgosdk/)** - High-level Steem Go SDK built on steemutil
- **[steem-js](https://github.com/steemit/steem-js)** - JavaScript Steem library (compatible with our SignedCall)
- **[steem](https://github.com/steemit/steem)** - Official Steem blockchain implementation

## Support

- **Issues**: [GitHub Issues](https://github.com/steemit/steemutil/issues)
- **Documentation**: [API Documentation](https://pkg.go.dev/github.com/steemit/steemutil)

---

**Note**: This library provides low-level utilities. For high-level application development, consider using [steemgosdk](https://github.com/steemit/steemgosdk/) which provides a more convenient API built on top of steemutil.
