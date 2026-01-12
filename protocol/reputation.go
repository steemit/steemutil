package protocol

import (
	"math"
)

// RepLog10 calculates the reputation score in log10 format.
// This function converts a raw reputation value (typically a large integer)
// into a more readable format used by Steem.
//
// The formula is: log10(abs(reputation)) - 9
// For negative reputation values, the result is negated.
//
// Examples:
//   - RepLog10(0) returns 0
//   - RepLog10(1000000000) returns 0 (log10(1e9) - 9 = 0)
//   - RepLog10(10000000000) returns 1 (log10(1e10) - 9 = 1)
//   - RepLog10(-1000000000) returns 0 (negative values are handled)
//
// This implementation is based on the standard Steem reputation calculation.
func RepLog10(reputation int64) int64 {
	if reputation == 0 {
		return 0
	}

	// Handle negative reputation values
	// Special case for MinInt64 to avoid overflow when negating
	isNegative := reputation < 0
	var absRep uint64
	if reputation == math.MinInt64 {
		// MinInt64's absolute value is MaxInt64 + 1, which overflows int64
		// Use uint64 to handle this case safely
		absRep = uint64(math.MaxInt64) + 1
	} else if isNegative {
		absRep = uint64(-reputation)
	} else {
		absRep = uint64(reputation)
	}

	// Calculate log10 and subtract 9 to normalize the result
	// Use math.Log10 with float64 conversion for precision
	logValue := math.Log10(float64(absRep))
	result := int64(logValue) - 9

	// Negate result if original reputation was negative
	if isNegative {
		return -result
	}

	return result
}
