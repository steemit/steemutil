package protocol

import (
	"strings"
	"time"

	"github.com/steemit/steemutil/encoder"
)

const Layout = `"2006-01-02T15:04:05"`
const LayoutWithoutQuotes = `2006-01-02T15:04:05`

type Time struct {
	Time *time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(Layout)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	// Remove quotes if present
	timeStr := strings.Trim(string(data), `"`)
	
	// Try parsing with the layout (without quotes)
	parsed, err := time.ParseInLocation(LayoutWithoutQuotes, timeStr, time.UTC)
	if err != nil {
		// Try RFC3339 format as fallback
		parsed, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			// Try RFC3339Nano as another fallback
			parsed, err = time.Parse(time.RFC3339Nano, timeStr)
			if err != nil {
				return err
			}
		}
	}
	t.Time = &parsed
	return nil
}

func (t *Time) MarshalTransaction(encoderObj *encoder.Encoder) error {
	return encoderObj.Encode(uint32(t.Time.Unix()))
}
