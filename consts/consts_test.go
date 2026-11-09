package consts

import "testing"

func TestSetAddressPrefixValidation(t *testing.T) {
	orig := ADDRESS_PREFIX
	defer func() { ADDRESS_PREFIX = orig }()

	cases := []struct {
		in      string
		wantErr bool
	}{
		{"TST", false},
		{"STM", false},
		{"", true},
		{"tst", true},        // lowercase rejected
		{"TS1", true},        // digit rejected
		{"TSTTSTTSTT", true}, // > 8 chars rejected
	}
	for _, c := range cases {
		err := SetAddressPrefix(c.in)
		if c.wantErr && err == nil {
			t.Errorf("SetAddressPrefix(%q): expected error, got nil", c.in)
		}
		if !c.wantErr {
			if err != nil {
				t.Errorf("SetAddressPrefix(%q): unexpected error: %v", c.in, err)
			}
			if ADDRESS_PREFIX != c.in {
				t.Errorf("SetAddressPrefix(%q): ADDRESS_PREFIX = %q", c.in, ADDRESS_PREFIX)
			}
		}
	}
}
