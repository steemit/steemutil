package api

import (
	"testing"
)

// TestParsePrice 校验 Price 构造与合法性约束。
func TestParsePrice(t *testing.T) {
	// 正常构造。
	p, err := ParsePrice("1.000 SBD", "1.000 STEEM")
	if err != nil {
		t.Fatalf("ParsePrice: %v", err)
	}
	if p.Base.Symbol != "SBD" || p.Quote.Symbol != "STEEM" {
		t.Errorf("unexpected price: %+v", p)
	}

	// 同 symbol 必须报错。
	if _, err := ParsePrice("1.000 STEEM", "2.000 STEEM"); err == nil {
		t.Error("ParsePrice same symbol: expected error")
	}
	// 零值必须报错。
	if _, err := ParsePrice("0.000 SBD", "1.000 STEEM"); err == nil {
		t.Error("ParsePrice zero base: expected error")
	}
	if _, err := ParsePrice("1.000 SBD", "0.000 STEEM"); err == nil {
		t.Error("ParsePrice zero quote: expected error")
	}
}

// TestPrice_Convert_Golden1 教学向量 1（计划文档）：
//
//	amount   = 3000000 (3000.000 STEEM，即 3×10^6 atoms)
//	quote    = 5       (0.000005 VESTS)
//	base     = 7       (0.007 STEEM)
//	result   = 3000000 * 5 / 7 = 2142857 (截断) → 2.142857 VESTS
//
// 校验 steemd asset*price 语义：先乘后除、向零截断。
// 注：计划文档原注释 "3.000 STEEM" 与 "3000000 atoms" 矛盾——3.000 STEEM 实为
// 3000 atoms；此处取 3000.000 STEEM (3000000 atoms) 以匹配文档的 amount 与 result。
func TestPrice_Convert_Golden1(t *testing.T) {
	// base=0.007 STEEM (7 atoms), quote=0.000005 VESTS (5 atoms)
	p, err := ParsePrice("0.007 STEEM", "0.000005 VESTS")
	if err != nil {
		t.Fatalf("ParsePrice: %v", err)
	}
	// 输入 3000.000 STEEM (3000000 atoms)，符号匹配 base，结果符号应为 quote(VESTS)。
	a, _ := ParseAsset("3000.000 STEEM")
	res, err := p.Convert(a)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if res.Amount != 2142857 {
		t.Errorf("Convert golden1 atoms = %d want 2142857", res.Amount)
	}
	if res.Symbol != "VESTS" {
		t.Errorf("Convert golden1 symbol = %q want VESTS", res.Symbol)
	}
	if got := res.String(); got != "2.142857 VESTS" {
		t.Errorf("Convert golden1 string = %q want %q", got, "2.142857 VESTS")
	}
}

// TestPrice_Convert_Golden2_128bit 教学向量 2（128 位中间值）：
//
//	asset.amount            = 1000           (1.000 STEEM)
//	total_vesting_shares    = 8×10^16        (quote, VESTS atoms)
//	total_vesting_fund_steem= 1.5×10^14      (base, STEEM atoms)
//	中间乘积 = 1000 × 8×10^16 = 8×10^19，超 int64 上限（9.2×10^18）
//	result = 8×10^19 / 1.5×10^14 = 533333 (截断) → 0.533333 VESTS
//
// 用 math/big.Int 不溢出，结果精确。校验本实现能正确处理 128 位中间值。
func TestPrice_Convert_Golden2_128bit(t *testing.T) {
	// fund  = 150000000000.000 STEEM = 150000000000 atoms (1.5×10^14)
	// shares= 80000000000.000000 VESTS = 80000000000000000 atoms (8×10^16)
	p, err := ParsePrice("150000000000.000 STEEM", "80000000000.000000 VESTS")
	if err != nil {
		t.Fatalf("ParsePrice: %v", err)
	}
	a, _ := ParseAsset("1.000 STEEM")
	res, err := p.Convert(a)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	// 1000 × 8×10^16 / 1.5×10^14 = 533333
	if res.Amount != 533333 {
		t.Errorf("Convert golden2 atoms = %d want 533333", res.Amount)
	}
	if res.Symbol != "VESTS" {
		t.Errorf("Convert golden2 symbol = %q want VESTS", res.Symbol)
	}
	if got := res.String(); got != "0.533333 VESTS" {
		t.Errorf("Convert golden2 string = %q want %q", got, "0.533333 VESTS")
	}
}

// TestPrice_Convert_Reverse 用 quote 分支而非 base 分支，
// 确认结果符号变为 base（Invert 后的方向）。
func TestPrice_Convert_Reverse(t *testing.T) {
	p, err := ParsePrice("1.000 SBD", "2.000 STEEM")
	if err != nil {
		t.Fatalf("ParsePrice: %v", err)
	}
	// 输入 STEEM（匹配 quote）→ 结果应为 base 符号 SBD。
	// result = 1.000 STEEM(1000) * 1.000 SBD(1000) / 2.000 STEEM(2000) = 500
	a, _ := ParseAsset("1.000 STEEM")
	res, err := p.Convert(a)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if res.Amount != 500 || res.Symbol != "SBD" {
		t.Errorf("Convert reverse = %+v want {500 SBD}", res)
	}
}

// TestPrice_Convert_UnknownSymbol 输入 symbol 既非 base 也非 quote → error。
func TestPrice_Convert_UnknownSymbol(t *testing.T) {
	p, _ := ParsePrice("1.000 SBD", "1.000 STEEM")
	unknown, _ := ParseAsset("1.000000 VESTS")
	if _, err := p.Convert(unknown); err == nil {
		t.Error("Convert unknown symbol: expected error")
	}
}

// TestPrice_Convert_Truncation 验证向零截断（Quo 语义，非 Div 向下取整）。
// 例：5/2=2 (positive), -5/2=-2 在 C++ 也是向零；这里只测正值，且验证非 floor。
func TestPrice_Convert_Truncation(t *testing.T) {
	// base=3 STEEM, quote=2 SBD；输入 1.000 STEEM → 1000*2/3 = 666 (截断)，而非 667。
	p, _ := ParsePrice("3.000 STEEM", "2.000 SBD")
	a, _ := ParseAsset("1.000 STEEM")
	res, err := p.Convert(a)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if res.Amount != 666 {
		t.Errorf("truncation atoms = %d want 666 (truncated, not rounded)", res.Amount)
	}
}

// TestPrice_Convert_Overflow 构造结果超过 MaxSatoshis 的输入（验收标准 3）。
// base=1 atom STEEM, quote=MaxSatoshis atoms SBD；输入 2.000 STEEM (2000 atoms)
// → 2000 * MaxSatoshis / 1 远超 int64，应 error。
func TestPrice_Convert_Overflow(t *testing.T) {
	// quote.amount = MaxSatoshis，但要能写成合法 asset 字符串。
	// 用直接构造 Price 绕过 ParseAsset 的字符串精度限制。
	p := Price{
		Base:  Asset{Amount: 1, Symbol: "STEEM"},
		Quote: Asset{Amount: MaxSatoshis, Symbol: "SBD"},
	}
	a := Asset{Amount: 2000, Symbol: "STEEM"}
	_, err := p.Convert(a)
	if err == nil {
		t.Error("Convert overflow: expected error")
	}
}

// TestPrice_Compare_Golden 教学向量（计划文档）：
//
//	A: 1.000 STEEM : 2.000 SBD   (base=1000, quote=2000)
//	B: 3.000 STEEM : 7.000 SBD   (base=3000, quote=7000)
//	A < B 吗？交叉乘：A.Base*B.Quote vs B.Base*A.Quote
//	  = 1000*7000 vs 3000*2000 = 7000000 vs 6000000 → A > B（lhs>rhs）
//
// 即 Compare(B) → -1（A 在 B 之后），实际是 p.Compare(q) 返回 lhs vs rhs。
// 本测试直接断言交叉乘结果。
func TestPrice_Compare_Golden(t *testing.T) {
	a, _ := ParsePrice("1.000 STEEM", "2.000 SBD")
	b, _ := ParsePrice("3.000 STEEM", "7.000 SBD")
	// p.Compare(q): p.Base*q.Quote vs q.Base*p.Quote
	//   = 1000*7000 vs 3000*2000 = 7000000 vs 6000000 → +1 (p > q)
	if got := a.Compare(b); got != 1 {
		t.Errorf("A.Compare(B) = %d want +1 (A>B by cross-mult)", got)
	}
	// 反向应为 -1。
	if got := b.Compare(a); got != -1 {
		t.Errorf("B.Compare(A) = %d want -1", got)
	}
}

// TestPrice_Compare_Equal 相等的交叉乘返回 0。
func TestPrice_Compare_Equal(t *testing.T) {
	a, _ := ParsePrice("1.000 STEEM", "2.000 SBD")
	b, _ := ParsePrice("2.000 STEEM", "4.000 SBD") // 同比值
	if got := a.Compare(b); got != 0 {
		t.Errorf("equal prices Compare = %d want 0", got)
	}
}

// TestPrice_Invert_RoundTrip Invert().Invert() == 原 Price（验收标准 4）。
func TestPrice_Invert_RoundTrip(t *testing.T) {
	original, _ := ParsePrice("1.000 SBD", "2.000 STEEM")
	roundTrip := original.Invert().Invert()
	if roundTrip.Base != original.Base || roundTrip.Quote != original.Quote {
		t.Errorf("Invert round-trip mismatch:\n  orig = %+v\n  rt   = %+v", original, roundTrip)
	}
	// 单次 Invert 应交换 base/quote。
	inverted := original.Invert()
	if inverted.Base != original.Quote || inverted.Quote != original.Base {
		t.Errorf("Invert did not swap: orig=%+v inv=%+v", original, inverted)
	}
}

// TestPrice_Compare_NoOverflow 用大数值确保交叉乘不会 int64 溢出。
// 两个 base.amount 接近 2^31 的 price 相乘 → ~2^62，仍安全（这里只是回归用例）。
func TestPrice_Compare_NoOverflow(t *testing.T) {
	// 用直接构造避免 ParseAsset 精度限制；amount 在 [0, MaxSatoshis] 内合法。
	a := Price{Base: Asset{Amount: 1 << 31, Symbol: "STEEM"}, Quote: Asset{Amount: 1 << 31, Symbol: "SBD"}}
	b := Price{Base: Asset{Amount: 1 << 31, Symbol: "STEEM"}, Quote: Asset{Amount: 1 << 31, Symbol: "SBD"}}
	if got := a.Compare(b); got != 0 {
		t.Errorf("large equal prices Compare = %d want 0", got)
	}
}
