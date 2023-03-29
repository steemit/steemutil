package protocol

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"

	"github.com/steemit/steemutil/encoder"
)

func unmarshalUInt(data []byte) (uint64, error) {
	if len(data) == 0 {
		return 0, errors.New("types: empty data received when unmarshalling an unsigned integer")
	}

	var (
		i   uint64
		err error
	)
	if data[0] == '"' {
		d := data[1:]
		d = d[:len(d)-1]
		i, err = strconv.ParseUint(string(d), 10, 64)
	} else {
		err = json.Unmarshal(data, &i)
	}
	return i, errors.Wrapf(err, "types: failed to unmarshal unsigned integer: %v", data)
}

func unmarshalInt(data []byte) (int64, error) {
	if len(data) == 0 {
		return 0, errors.New("types: empty data received when unmarshalling an integer")
	}

	var (
		i   int64
		err error
	)
	if data[0] == '"' {
		d := data[1:]
		d = d[:len(d)-1]
		i, err = strconv.ParseInt(string(d), 10, 64)
	} else {
		err = json.Unmarshal(data, &i)
	}
	return i, errors.Wrapf(err, "types: failed to unmarshal integer: %v", data)
}

type UInt uint

func (num *UInt) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt(v)
	return nil
}

func (num UInt) Marshalencoder(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(uint(num))
}

type UInt8 uint8

func (num *UInt8) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt8(v)
	return nil
}

func (num UInt8) MarshalTransaction(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(uint8(num))
}

type UInt16 uint16

func (num *UInt16) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt16(v)
	return nil
}

func (num UInt16) MarshalTransaction(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(uint16(num))
}

type UInt32 uint32

func (num *UInt32) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt32(v)
	return nil
}

func (num UInt32) MarshalTransaction(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(uint32(num))
}

type UInt64 uint64

func (num *UInt64) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt64(v)
	return nil
}

func (num UInt64) MarshalTransaction(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(uint64(num))
}

type Int int

func (num *Int) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int(v)
	return nil
}

type Int8 int8

func (num *Int8) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int8(v)
	return nil
}

func (num Int8) MarshalTransaction(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(int(num))
}

type Int16 int16

func (num *Int16) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int16(v)
	return nil
}

func (num Int16) MarshalTransaction(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(int16(num))
}

type Int32 int32

func (num *Int32) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int32(v)
	return nil
}

func (num Int32) MarshalTransaction(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(int32(num))
}

type Int64 int64

func (num *Int64) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int64(v)
	return nil
}

func (num Int64) MarshalTransaction(encoder *encoder.Encoder) error {
	return encoder.EncodeNumber(int64(num))
}
