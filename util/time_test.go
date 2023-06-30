package util_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/util"
)

var (
	timeTestData        = "2023-06-29T08:05:30"
	timeHexLittleEndian = "4a3b9d64"
)

func TestTime_MarshalJSON(t *testing.T) {
	expectObj, err := time.Parse(util.Layout, timeTestData)
	if err != nil {
		t.Error(err)
	}
	testTimeObj := &util.Time{
		Time: &expectObj,
	}
	got, err := testTimeObj.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	expected := []byte(timeTestData)

	if !bytes.Equal(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
		fmt.Println(expected, got)
	}
}

func TestTime_UnmarshalJSON(t *testing.T) {
	testTimeObj := &util.Time{}
	err := testTimeObj.UnmarshalJSON([]byte(timeTestData))
	if err != nil {
		t.Error(err)
	}
	got, err := testTimeObj.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	expected := []byte(timeTestData)
	if !bytes.Equal(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestTime_Serialize(t *testing.T) {
	testTimeObj := &util.Time{}
	err := testTimeObj.UnmarshalJSON([]byte(timeTestData))
	if err != nil {
		t.Error(err)
	}
	var b bytes.Buffer
	encoderObj := encoder.NewEncoder(&b)
	err = testTimeObj.Serialize(encoderObj)
	if err != nil {
		t.Error(err)
	}
	encodedString := hex.EncodeToString(b.Bytes())
	if encodedString != timeHexLittleEndian {
		t.Errorf("expected %v, got %v", timeHexLittleEndian, encodedString)
	}
}
