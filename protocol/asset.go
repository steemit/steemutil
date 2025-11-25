package protocol

import (
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/steemit/steemutil/encoder"

	"github.com/pkg/errors"
)

// Asset represents a Steem asset (amount + symbol).
// Format: int64 amount (without decimal point) + uint8 precision + 7 bytes symbol (null-padded)
type Asset struct {
	Amount    int64  `json:"amount"`
	Precision uint8  `json:"precision"`
	Symbol    string `json:"symbol"`
}

// ParseAsset parses an asset string like "0.001 STEEM" into an Asset struct.
func ParseAsset(assetStr string) (*Asset, error) {
	parts := strings.Split(strings.TrimSpace(assetStr), " ")
	if len(parts) != 2 {
		return nil, errors.Errorf("invalid asset format: %s (expected 'amount symbol')", assetStr)
	}

	amountStr := parts[0]
	symbol := strings.ToUpper(parts[1])

	// Parse amount and calculate precision
	dotIndex := strings.Index(amountStr, ".")
	var amountInt int64
	var precision uint8

	if dotIndex == -1 {
		// No decimal point
		var err error
		amountInt, err = strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse amount: %s", amountStr)
		}
		precision = 0
	} else {
		// Has decimal point
		// Remove decimal point and parse as integer
		amountWithoutDot := strings.Replace(amountStr, ".", "", 1)
		var err error
		amountInt, err = strconv.ParseInt(amountWithoutDot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse amount: %s", amountStr)
		}
		precision = uint8(len(amountStr) - dotIndex - 1)
	}

	return &Asset{
		Amount:    amountInt,
		Precision: precision,
		Symbol:    symbol,
	}, nil
}

// String returns the asset as a string like "0.001 STEEM".
func (a *Asset) String() string {
	if a == nil {
		return "0.000 STEEM"
	}

	amountStr := strconv.FormatInt(a.Amount, 10)
	if a.Precision > 0 {
		// Pad with zeros if needed
		for len(amountStr) < int(a.Precision) {
			amountStr = "0" + amountStr
		}
		// Insert decimal point
		dotPos := len(amountStr) - int(a.Precision)
		if dotPos == 0 {
			// If dotPos is 0, we need to add "0." prefix
			amountStr = "0." + amountStr
		} else {
			amountStr = amountStr[:dotPos] + "." + amountStr[dotPos:]
		}
	}

	return amountStr + " " + a.Symbol
}

// MarshalTransaction implements the asset binary serialization format.
// Format: int64 amount (little-endian) + uint8 precision + 7 bytes symbol (null-padded)
func (a *Asset) MarshalTransaction(encoderObj *encoder.Encoder) error {
	if a == nil {
		return errors.New("cannot marshal nil asset")
	}

	// Encode amount as int64 (little-endian)
	if err := encoderObj.EncodeNumber(a.Amount); err != nil {
		return errors.Wrap(err, "failed to encode asset amount")
	}

	// Encode precision as uint8
	if err := encoderObj.EncodeNumber(a.Precision); err != nil {
		return errors.Wrap(err, "failed to encode asset precision")
	}

	// Encode symbol as 7 bytes (null-padded)
	symbolBytes := make([]byte, 7)
	symbolUpper := strings.ToUpper(a.Symbol)
	copy(symbolBytes, symbolUpper)
	// Remaining bytes are already zero (null-padded)

	// Write symbol bytes directly (7 bytes, null-padded)
	if err := encoderObj.WriteBytes(symbolBytes); err != nil {
		return errors.Wrap(err, "failed to encode asset symbol")
	}

	return nil
}

// UnmarshalAsset unmarshals an asset from binary data.
func UnmarshalAsset(data []byte) (*Asset, error) {
	if len(data) < 16 {
		return nil, errors.New("insufficient data for asset (need at least 16 bytes)")
	}

	// Read int64 amount (little-endian)
	amount := int64(binary.LittleEndian.Uint64(data[0:8]))

	// Read uint8 precision
	precision := data[8]

	// Read 7 bytes symbol (null-terminated)
	symbolBytes := data[9:16]
	symbol := strings.TrimRight(string(symbolBytes), "\x00")

	return &Asset{
		Amount:    amount,
		Precision: precision,
		Symbol:    symbol,
	}, nil
}

