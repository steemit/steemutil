package api

import (
	"testing"
)

// TestParseAsset_RoundTrip 覆盖验收标准 1：ParseAsset/String 严格互逆。
// 全程不经 float，VESTS 6 位小数也能精确重建。
func TestParseAsset_RoundTrip(t *testing.T) {
	cases := []string{
		"0.001 STEEM",
		"1.000 SBD",
		"100.500 SBD",
		"0.000001 VESTS",
		"12345.678 STEEM",
		"0.123456 VESTS",
		"9999999.999 SBD",
	}
	for _, s := range cases {
		a, err := ParseAsset(s)
		if err != nil {
			t.Fatalf("ParseAsset(%q): %v", s, err)
		}
		if got := a.String(); got != s {
			t.Errorf("round-trip mismatch: in=%q out=%q (atoms=%d)", s, got, a.Amount)
		}
	}
}

// TestParseAsset_AtomValues 手算原子单位，确认小数→int64 拼接无误。
func TestParseAsset_AtomValues(t *testing.T) {
	cases := []struct {
		in     string
		amount int64
		symbol string
	}{
		{"0.001 STEEM", 1, "STEEM"},
		{"1.000 SBD", 1000, "SBD"},
		{"0.000001 VESTS", 1, "VESTS"},
		{"1.000000 VESTS", 1000000, "VESTS"},
		{"1.500 STEEM", 1500, "STEEM"},
		{"1000.000 STEEM", 1000000, "STEEM"},
	}
	for _, c := range cases {
		a, err := ParseAsset(c.in)
		if err != nil {
			t.Fatalf("ParseAsset(%q): %v", c.in, err)
		}
		if a.Amount != c.amount {
			t.Errorf("%q: amount=%d want %d", c.in, a.Amount, c.amount)
		}
		if a.Symbol != c.symbol {
			t.Errorf("%q: symbol=%q want %q", c.in, a.Symbol, c.symbol)
		}
	}
}

// TestParseAsset_SymbolLowercase 校验大小写不敏感（symbol 转大写）。
func TestParseAsset_SymbolLowercase(t *testing.T) {
	a, err := ParseAsset("1.000 steem")
	if err != nil {
		t.Fatalf("ParseAsset lowercase: %v", err)
	}
	if a.Symbol != "STEEM" {
		t.Errorf("expected uppercased STEEM, got %q", a.Symbol)
	}
}

// TestParseAsset_Errors 覆盖验收标准 1 的非法输入分支。
func TestParseAsset_Errors(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"missing space", "1.000STEEM"},
		{"too many parts", "1.000 STEEM EXTRA"},
		{"unknown symbol", "1.000 XYZ"},
		{"vests wrong precision", "0.001 VESTS"},  // VESTS 需 6 位
		{"steem wrong precision", "0.0001 STEEM"}, // STEEM 需 3 位
		{"empty amount", " STEEM"},
		{"non-numeric", "a.000 STEEM"},
		{"negative", "-1.000 STEEM"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			_, err := ParseAsset(c.in)
			if err == nil {
				t.Errorf("ParseAsset(%q): expected error, got nil", c.in)
			}
		})
	}
}

// TestParseAsset_ExceedsMaxSatoshis 构造超 MaxSatoshis 的输入。
// MaxSatoshis = 2^62-1 = 4611686018427387903。
// 4611686018427387904 = 2^62 即越界。
func TestParseAsset_ExceedsMaxSatoshis(t *testing.T) {
	over := "4611686018427387.904 STEEM" // 原子 = 2^62，超过 MaxSatoshis
	_, err := ParseAsset(over)
	if err == nil {
		t.Errorf("ParseAsset(%q): expected range error, got nil", over)
	}
	// 边界值：恰好等于 MaxSatoshis 应通过。
	maxAtoms := "4611686018427387.903 STEEM"
	a, err := ParseAsset(maxAtoms)
	if err != nil {
		t.Fatalf("ParseAsset(MaxSatoshis) unexpected err: %v", err)
	}
	if a.Amount != MaxSatoshis {
		t.Errorf("MaxSatoshis atoms = %d want %d", a.Amount, MaxSatoshis)
	}
}

// TestAsset_AddSubZeroFloatError 覆盖验收标准 1 的"0.1+0.2 零误差"等价。
// 用原子单位：0.100 STEEM + 0.200 STEEM = 0.300 STEEM（100+200=300），
// 零浮点 round-off。
func TestAsset_AddSubZeroFloatError(t *testing.T) {
	a, _ := ParseAsset("0.100 STEEM")
	b, _ := ParseAsset("0.200 STEEM")
	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if sum.Amount != 300 {
		t.Errorf("0.100+0.200 atoms = %d want 300", sum.Amount)
	}
	if got := sum.String(); got != "0.300 STEEM" {
		t.Errorf("sum string = %q want %q", got, "0.300 STEEM")
	}

	diff, err := sum.Sub(b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	if diff.Amount != 100 {
		t.Errorf("0.300-0.200 atoms = %d want 100", diff.Amount)
	}
	if got := diff.String(); got != "0.100 STEEM" {
		t.Errorf("diff string = %q want %q", got, "0.100 STEEM")
	}
}

// TestAsset_Add_DifferentSymbol 不同 symbol 必须报错。
func TestAsset_Add_DifferentSymbol(t *testing.T) {
	a, _ := ParseAsset("1.000 STEEM")
	b, _ := ParseAsset("1.000 SBD")
	if _, err := a.Add(b); err == nil {
		t.Error("Add across symbols: expected error")
	}
	if _, err := a.Sub(b); err == nil {
		t.Error("Sub across symbols: expected error")
	}
}

// TestAsset_Add_Overflow 构造两个接近 MaxSatoshis 的值相加溢出。
func TestAsset_Add_Overflow(t *testing.T) {
	// MaxSatoshis - 1 与 2 相加 → MaxSatoshis+1，溢出。
	a := Asset{Amount: MaxSatoshis - 1, Symbol: "STEEM"}
	b := Asset{Amount: 2, Symbol: "STEEM"}
	if _, err := a.Add(b); err == nil {
		t.Error("Add overflow: expected error")
	}
	// 边界：MaxSatoshis-1 + 1 = MaxSatoshis，恰好不溢出。
	c := Asset{Amount: 1, Symbol: "STEEM"}
	sum, err := a.Add(c)
	if err != nil {
		t.Fatalf("Add boundary: %v", err)
	}
	if sum.Amount != MaxSatoshis {
		t.Errorf("boundary sum = %d want %d", sum.Amount, MaxSatoshis)
	}
}

// TestAsset_Precision 确认精度派生正确。
func TestAsset_Precision(t *testing.T) {
	cases := []struct {
		symbol string
		prec   int
	}{
		{"STEEM", PrecisionSteem},
		{"SBD", PrecisionSbd},
		{"VESTS", PrecisionVests},
		{"TESTS", PrecisionSteem},
		{"TBD", PrecisionSteem},
	}
	for _, c := range cases {
		a := Asset{Amount: 1, Symbol: c.symbol}
		if got := a.Precision(); got != c.prec {
			t.Errorf("%s precision = %d want %d", c.symbol, got, c.prec)
		}
	}
}

// TestAsset_String_Format 精确重建小数点格式，无 round-half-even。
func TestAsset_String_Format(t *testing.T) {
	cases := []struct {
		asset Asset
		out   string
	}{
		{Asset{Amount: 0, Symbol: "STEEM"}, "0.000 STEEM"},
		{Asset{Amount: 1, Symbol: "STEEM"}, "0.001 STEEM"},
		{Asset{Amount: 1500, Symbol: "STEEM"}, "1.500 STEEM"},
		{Asset{Amount: 1000000, Symbol: "STEEM"}, "1000.000 STEEM"},
		{Asset{Amount: 1, Symbol: "VESTS"}, "0.000001 VESTS"},
		{Asset{Amount: 999999, Symbol: "VESTS"}, "0.999999 VESTS"},
	}
	for _, c := range cases {
		if got := c.asset.String(); got != c.out {
			t.Errorf("String(%d %s) = %q want %q", c.asset.Amount, c.asset.Symbol, got, c.out)
		}
	}
}
