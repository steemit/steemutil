package util_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/util"
)

type testIntStruct struct {
	ItemUInt   *util.UInt   `json:"item_uint,omitempty"`
	ItemUInt8  *util.UInt8  `json:"item_uint8,omitempty"`
	ItemUInt16 *util.UInt16 `json:"item_uint16,omitempty"`
	ItemUInt32 *util.UInt32 `json:"item_uint32,omitempty"`
	ItemUInt64 *util.UInt64 `json:"item_uint64,omitempty"`
	ItemInt    *util.Int    `json:"item_int,omitempty"`
	ItemInt8   *util.Int8   `json:"item_int8,omitempty"`
	ItemInt16  *util.Int16  `json:"item_int16,omitempty"`
	ItemInt32  *util.Int32  `json:"item_int32,omitempty"`
	ItemInt64  *util.Int64  `json:"item_int64,omitempty"`
}

func TestUInt_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_uint": "345"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemUInt != 345 {
		t.Errorf("got %v", *testdata.ItemUInt)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_uint": 345}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemUInt != 345 {
		t.Errorf("got %v", *testdataInt.ItemUInt)
	}
}

func TestUInt_MarshalJSON(t *testing.T) {
	var data util.UInt = 345
	var expected string = `{"item_uint":345}`
	testdata := &testIntStruct{
		ItemUInt: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestUInt_Serialize(t *testing.T) {
	var data util.UInt = 345
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("59010000")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestUInt8_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_uint8": "35"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemUInt8 != 35 {
		t.Errorf("got %v", *testdata.ItemUInt8)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_uint8": 35}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemUInt8 != 35 {
		t.Errorf("got %v", *testdataInt.ItemUInt8)
	}
}

func TestUInt8_MarshalJSON(t *testing.T) {
	var data util.UInt8 = 35
	var expected string = `{"item_uint8":35}`
	testdata := &testIntStruct{
		ItemUInt8: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestUInt8_Serialize(t *testing.T) {
	var data util.UInt8 = 35
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("23")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestUInt16_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_uint16": "35"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemUInt16 != 35 {
		t.Errorf("got %v", *testdata.ItemUInt16)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_uint16": 35}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemUInt16 != 35 {
		t.Errorf("got %v", *testdataInt.ItemUInt16)
	}
}

func TestUInt16_MarshalJSON(t *testing.T) {
	var data util.UInt16 = 35
	var expected string = `{"item_uint16":35}`
	testdata := &testIntStruct{
		ItemUInt16: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestUInt16_Serialize(t *testing.T) {
	var data util.UInt16 = 35
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("2300")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestUInt32_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_uint32": "35"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemUInt32 != 35 {
		t.Errorf("got %v", *testdata.ItemUInt32)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_uint32": 35}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemUInt32 != 35 {
		t.Errorf("got %v", *testdataInt.ItemUInt32)
	}
}

func TestUInt32_MarshalJSON(t *testing.T) {
	var data util.UInt32 = 35
	var expected string = `{"item_uint32":35}`
	testdata := &testIntStruct{
		ItemUInt32: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestUInt32_Serialize(t *testing.T) {
	var data util.UInt32 = 35
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("23000000")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestUInt64_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_uint64": "35"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemUInt64 != 35 {
		t.Errorf("got %v", *testdata.ItemUInt64)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_uint64": 35}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemUInt64 != 35 {
		t.Errorf("got %v", *testdataInt.ItemUInt64)
	}
}

func TestUInt64_MarshalJSON(t *testing.T) {
	var data util.UInt64 = 35
	var expected string = `{"item_uint64":35}`
	testdata := &testIntStruct{
		ItemUInt64: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestUInt64_Serialize(t *testing.T) {
	var data util.UInt64 = 35
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("2300000000000000")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestInt_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_int": "345"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemInt != 345 {
		t.Errorf("got %v", *testdata.ItemInt)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_int": 345}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemInt != 345 {
		t.Errorf("got %v", *testdataInt.ItemInt)
	}
}

func TestInt_MarshalJSON(t *testing.T) {
	var data util.Int = 345
	var expected string = `{"item_int":345}`
	testdata := &testIntStruct{
		ItemInt: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestInt_Serialize(t *testing.T) {
	var data util.Int = 345
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("59010000")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestInt8_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_int8": "35"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemInt8 != 35 {
		t.Errorf("got %v", *testdata.ItemInt8)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_int8": 35}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemInt8 != 35 {
		t.Errorf("got %v", *testdataInt.ItemInt8)
	}
}

func TestInt8_MarshalJSON(t *testing.T) {
	var data util.Int8 = 35
	var expected string = `{"item_int8":35}`
	testdata := &testIntStruct{
		ItemInt8: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestInt8_Serialize(t *testing.T) {
	var data util.Int8 = 35
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("23")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestInt16_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_int16": "35"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemInt16 != 35 {
		t.Errorf("got %v", *testdata.ItemInt16)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_int16": 35}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemInt16 != 35 {
		t.Errorf("got %v", *testdataInt.ItemInt16)
	}
}

func TestInt16_MarshalJSON(t *testing.T) {
	var data util.Int16 = 35
	var expected string = `{"item_int16":35}`
	testdata := &testIntStruct{
		ItemInt16: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestInt16_Serialize(t *testing.T) {
	var data util.Int16 = 35
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("2300")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestInt32_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_int32": "35"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemInt32 != 35 {
		t.Errorf("got %v", *testdata.ItemInt32)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_int32": 35}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemInt32 != 35 {
		t.Errorf("got %v", *testdataInt.ItemInt32)
	}
}

func TestInt32_MarshalJSON(t *testing.T) {
	var data util.Int32 = 35
	var expected string = `{"item_int32":35}`
	testdata := &testIntStruct{
		ItemInt32: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestInt32_Serialize(t *testing.T) {
	var data util.Int32 = 35
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("23000000")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}

func TestInt64_UnmarshalJSON(t *testing.T) {
	var testdata testIntStruct
	err := json.Unmarshal([]byte(`{"item_int64": "35"}`), &testdata)
	if err != nil {
		t.Error(err)
	}
	if *testdata.ItemInt64 != 35 {
		t.Errorf("got %v", *testdata.ItemInt64)
	}

	var testdataInt testIntStruct
	err = json.Unmarshal([]byte(`{"item_int64": 35}`), &testdataInt)
	if err != nil {
		t.Error(err)
	}
	if *testdataInt.ItemInt64 != 35 {
		t.Errorf("got %v", *testdataInt.ItemInt64)
	}
}

func TestInt64_MarshalJSON(t *testing.T) {
	var data util.Int64 = 35
	var expected string = `{"item_int64":35}`
	testdata := &testIntStruct{
		ItemInt64: &data,
	}
	got, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
	}
	if string(got) != expected {
		t.Errorf("got %v, expected %v", string(got), expected)
	}
}

func TestInt64_Serialize(t *testing.T) {
	var data util.Int64 = 35
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err := data.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	expected, err := hex.DecodeString("2300000000000000")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, b.Bytes()) {
		t.Errorf("got %v, expected %v", b.Bytes(), expected)
	}
}
