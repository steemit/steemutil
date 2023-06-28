package util_test

import (
	"bytes"
	"testing"

	"github.com/steemit/steemutil/util"
)

var (
	expected = "[[\"abc\",1],[\"def\",2],[\"ghi\",3]]"
	testObj  = &util.StringInt64Map{
		"abc": 1,
		"def": 2,
		"ghi": 3,
	}
)

func TestStringInt64Map_MarshalJSON(t *testing.T) {
	result, err := testObj.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	got := string(result)
	if expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestStringInt64Map_UnmarshalJSON(t *testing.T) {
	obj := &util.StringInt64Map{}
	obj.UnmarshalJSON([]byte(expected))
	got, err := obj.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	expected, err := testObj.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
