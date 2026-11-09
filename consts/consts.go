package consts

import "fmt"

// ADDRESS_PREFIX is the public key address prefix used when encoding and
// parsing public key strings. It defaults to the Steem mainnet prefix
// ("STM").
//
// It is a var (not a const) so alternative chains can switch it to their
// own prefix — the Steem testnet uses "TST" — via SetAddressPrefix.
// Call SetAddressPrefix once at process startup, before any key is
// encoded or parsed; it is NOT safe for concurrent modification.
var ADDRESS_PREFIX = defaultAddressPrefix

const defaultAddressPrefix = "STM"

// SetAddressPrefix overrides the public key address prefix (e.g. "TST" for
// the Steem testnet). It validates the prefix and mutates the package-level
// ADDRESS_PREFIX. Intended to be called once at startup; not concurrency
// safe.
func SetAddressPrefix(prefix string) error {
	if prefix == "" {
		return fmt.Errorf("address prefix must not be empty")
	}
	if len(prefix) > 8 {
		return fmt.Errorf("address prefix too long: %q (max 8 chars)", prefix)
	}
	for i := 0; i < len(prefix); i++ {
		c := prefix[i]
		if c < 'A' || c > 'Z' {
			return fmt.Errorf("address prefix must be uppercase A-Z, got %q", prefix)
		}
	}
	ADDRESS_PREFIX = prefix
	return nil
}
