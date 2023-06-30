package util

import (
	"time"

	"github.com/steemit/steemutil/encoder"
)

const Layout = "2006-01-02T15:04:05"

type Time struct {
	Time *time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(Layout)), nil
}

// Time location MUST be UTC.
func (t *Time) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(Layout, string(data))
	if err != nil {
		return err
	}
	t.Time = &parsed
	return nil
}

func (t *Time) Serialize(encoderObj *encoder.Encoder) error {
	return encoderObj.Encode(uint32(t.Time.Unix()))
}
