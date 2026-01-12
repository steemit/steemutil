package protocol

import (
	"math"
	"testing"
)

func TestRepLog10(t *testing.T) {
	tests := []struct {
		name       string
		reputation int64
		expected   int64
	}{
		{
			name:       "zero reputation",
			reputation: 0,
			expected:   0,
		},
		{
			name:       "positive reputation 1e9",
			reputation: 1000000000,
			expected:   0, // log10(1e9) - 9 = 9 - 9 = 0
		},
		{
			name:       "positive reputation 1e10",
			reputation: 10000000000,
			expected:   1, // log10(1e10) - 9 = 10 - 9 = 1
		},
		{
			name:       "positive reputation 1e11",
			reputation: 100000000000,
			expected:   2, // log10(1e11) - 9 = 11 - 9 = 2
		},
		{
			name:       "positive reputation 1e8",
			reputation: 100000000,
			expected:   -1, // log10(1e8) - 9 = 8 - 9 = -1
		},
		{
			name:       "negative reputation",
			reputation: -1000000000,
			expected:   0, // Should return -RepLog10(1000000000) = 0
		},
		{
			name:       "negative reputation large",
			reputation: -10000000000,
			expected:   -1, // Should return -RepLog10(10000000000) = -1
		},
		{
			name:       "small positive reputation",
			reputation: 1,
			expected:   int64(math.Log10(1)) - 9, // log10(1) - 9 = 0 - 9 = -9
		},
		{
			name:       "medium positive reputation",
			reputation: 5000000000,
			expected:   int64(math.Log10(5000000000)) - 9, // log10(5e9) - 9 ≈ 9.699 - 9 = 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RepLog10(tt.reputation)
			if result != tt.expected {
				t.Errorf("RepLog10(%d) = %d, expected %d", tt.reputation, result, tt.expected)
			}
		})
	}
}

// TestRepLog10EdgeCases tests edge cases and boundary conditions
func TestRepLog10EdgeCases(t *testing.T) {
	// Test maximum int64 value
	maxInt64 := int64(math.MaxInt64)
	result := RepLog10(maxInt64)
	if result < 0 {
		t.Errorf("RepLog10(MaxInt64) should be positive, got %d", result)
	}

	// Test minimum int64 value (negative)
	// MinInt64's absolute value is MaxInt64 + 1, which requires special handling
	minInt64 := int64(math.MinInt64)
	result = RepLog10(minInt64)
	// Should handle negative values correctly (result should be negative or zero)
	if result > 0 {
		t.Errorf("RepLog10(MinInt64) should be negative or zero, got %d", result)
	}
	// Verify it doesn't cause overflow/panic
	expectedMin := int64(math.Log10(float64(uint64(math.MaxInt64)+1))) - 9
	if result != -expectedMin {
		t.Errorf("RepLog10(MinInt64) = %d, expected %d", result, -expectedMin)
	}

	// Test very small positive values
	smallValues := []int64{1, 10, 100, 1000, 10000}
	for _, val := range smallValues {
		result := RepLog10(val)
		expected := int64(math.Log10(float64(val))) - 9
		if result != expected {
			t.Errorf("RepLog10(%d) = %d, expected %d", val, result, expected)
		}
	}
}
