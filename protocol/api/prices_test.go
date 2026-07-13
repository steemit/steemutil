package api

import (
	_ "embed"
	"encoding/json"
	"testing"
)

//go:embed testdata/order_book.json
var orderBookJSON []byte

//go:embed testdata/feed_history.json
var feedHistoryJSON []byte

//go:embed testdata/dynamic_global_properties.json
var dgpJSON []byte

// loadFixture 读取 testdata JSON 文件到对应 wire 类型，确认 fixture 可反序列化。
func loadOrderBook(t *testing.T) *OrderBook {
	t.Helper()
	var ob OrderBook
	if err := json.Unmarshal(orderBookJSON, &ob); err != nil {
		t.Fatalf("unmarshal order_book fixture: %v", err)
	}
	return &ob
}

func loadFeedHistory(t *testing.T) *FeedHistory {
	t.Helper()
	var fh FeedHistory
	if err := json.Unmarshal(feedHistoryJSON, &fh); err != nil {
		t.Fatalf("unmarshal feed_history fixture: %v", err)
	}
	return &fh
}

func loadDGP(t *testing.T) *DynamicGlobalProperties {
	t.Helper()
	var dgp DynamicGlobalProperties
	if err := json.Unmarshal(dgpJSON, &dgp); err != nil {
		t.Fatalf("unmarshal dynamic_global_properties fixture: %v", err)
	}
	return &dgp
}

// TestToPrice_OrderPrice 桥接方法把 wire string-asset 解析为 Price。
func TestToPrice_OrderPrice(t *testing.T) {
	op := OrderPrice{Base: "1.000 SBD", Quote: "1.000 STEEM"}
	p, err := op.ToPrice()
	if err != nil {
		t.Fatalf("ToPrice: %v", err)
	}
	if p.Base.Symbol != "SBD" || p.Quote.Symbol != "STEEM" {
		t.Errorf("unexpected price: %+v", p)
	}
	if p.Base.Amount != 1000 || p.Quote.Amount != 1000 {
		t.Errorf("unexpected atoms: base=%d quote=%d", p.Base.Amount, p.Quote.Amount)
	}
}

// TestToPrice_CurrentMedianHistoryPrice 桥接方法对 feed entry 同样工作。
func TestToPrice_CurrentMedianHistoryPrice(t *testing.T) {
	c := CurrentMedianHistoryPrice{Base: "0.510 SBD", Quote: "1.000 STEEM"}
	p, err := c.ToPrice()
	if err != nil {
		t.Fatalf("ToPrice: %v", err)
	}
	if p.Base.Amount != 510 || p.Quote.Amount != 1000 {
		t.Errorf("unexpected atoms: base=%d quote=%d", p.Base.Amount, p.Quote.Amount)
	}
}

// TestToPrice_Errors 非法输入向上报错。
func TestToPrice_Errors(t *testing.T) {
	// 同 symbol。
	op := OrderPrice{Base: "1.000 STEEM", Quote: "1.000 STEEM"}
	if _, err := op.ToPrice(); err == nil {
		t.Error("OrderPrice.ToPrice same symbol: expected error")
	}
	// 非法 symbol。
	c := CurrentMedianHistoryPrice{Base: "1.000 XYZ", Quote: "1.000 STEEM"}
	if _, err := c.ToPrice(); err == nil {
		t.Error("CurrentMedianHistoryPrice.ToPrice bad symbol: expected error")
	}
}

// TestComputePrices_Fixture 用 testdata JSON 跑 ComputePrices，
// 断言每项精确原子值（验收标准 5）。
//
// 期望（手算）：
//   - asks: 1.000 SBD/1.000 STEEM → Convert(1.000 STEEM)=1000 SBD atoms
//     2.000 SBD/1.000 STEEM → 2000
//   - bids: 0.500 SBD/1.000 STEEM → 500
//     0.250 SBD/1.000 STEEM → 250
//     sum=3750, count=4 → avg=937 SBD atoms → "0.937 SBD"
//   - SteemUsd: feed 最后一项 base=0.510 SBD → {base=510 SBD, quote=1000 STEEM}
//   - SteemVest: fund=150000000000.000 STEEM(1.5×10^14 atoms), shares=80000000000.000000 VESTS(8×10^16 atoms)
func TestComputePrices_Fixture(t *testing.T) {
	ob := loadOrderBook(t)
	fh := loadFeedHistory(t)
	dgp := loadDGP(t)

	res, err := ComputePrices(ob, fh, dgp)
	if err != nil {
		t.Fatalf("ComputePrices: %v", err)
	}

	// SteemSbd
	if res.SteemSbd.Amount != 937 || res.SteemSbd.Symbol != "SBD" {
		t.Errorf("SteemSbd = %+v want {937 SBD}", res.SteemSbd)
	}
	if got := res.SteemSbd.String(); got != "0.937 SBD" {
		t.Errorf("SteemSbd.String = %q want %q", got, "0.937 SBD")
	}

	// SteemUsd
	if res.SteemUsd.Base.Amount != 510 || res.SteemUsd.Quote.Amount != 1000 {
		t.Errorf("SteemUsd = %+v want base=510 quote=1000", res.SteemUsd)
	}

	// SteemVest
	if res.SteemVest.Base.Amount != 150000000000000 || res.SteemVest.Quote.Amount != 80000000000000000 {
		t.Errorf("SteemVest = %+v want base=150000000000000 quote=80000000000000000", res.SteemVest)
	}
	// 顺便验证 SteemVest 的 vesting 换算与 golden2 一致（1 STEEM → 533333 VESTS atoms）。
	oneSteem, _ := ParseAsset("1.000 STEEM")
	vestPerSteem, err := res.SteemVest.Convert(oneSteem)
	if err != nil {
		t.Fatalf("SteemVest.Convert(1.000 STEEM): %v", err)
	}
	if vestPerSteem.Amount != 533333 {
		t.Errorf("vestPerSteem = %d want 533333", vestPerSteem.Amount)
	}
}

// TestComputePrices_Errors 各 nil / 空输入必须报错。
func TestComputePrices_Errors(t *testing.T) {
	fh := loadFeedHistory(t)
	dgp := loadDGP(t)

	// nil order book。
	if _, err := ComputePrices(nil, fh, dgp); err == nil {
		t.Error("nil order book: expected error")
	}
	// empty order book。
	emptyOB := &OrderBook{}
	if _, err := ComputePrices(emptyOB, fh, dgp); err == nil {
		t.Error("empty order book: expected error")
	}
	// nil feed history。
	ob := loadOrderBook(t)
	if _, err := ComputePrices(ob, nil, dgp); err == nil {
		t.Error("nil feed history: expected error")
	}
	// empty price_history。
	if _, err := ComputePrices(ob, &FeedHistory{}, dgp); err == nil {
		t.Error("empty price_history: expected error")
	}
	// nil dgp。
	if _, err := ComputePrices(ob, fh, nil); err == nil {
		t.Error("nil dgp: expected error")
	}
}

// TestComputePrices_SymbolMismatch order book 中混入不匹配 STEEM 的 order 必须报错。
func TestComputePrices_SymbolMismatch(t *testing.T) {
	ob := &OrderBook{
		Asks: []Order{{OrderPrice: OrderPrice{Base: "1.000 STEEM", Quote: "1.000 STEEM"}}}, // 非法
	}
	fh := loadFeedHistory(t)
	dgp := loadDGP(t)
	if _, err := ComputePrices(ob, fh, dgp); err == nil {
		t.Error("symbol mismatch in order: expected error")
	}
}

// === 审计修复 #1 回归 ===

// TestComputePrices_RejectsNonSbdOrder 修复 #1：order book 含非 STEEM/SBD 对
// （如 STEEM/VESTS）时，Convert 结果符号是 VESTS，必须报错而非贴硬编码 "SBD"
// 撒谎。此前的实现会返回 SteemSbd.Symbol="SBD" 但 Amount 实为 VESTS 原子，
// 导致 String() 输出连数值都错（VESTS 6 位 vs SBD 3 位精度差 1000 倍）。
func TestComputePrices_RejectsNonSbdOrder(t *testing.T) {
	ob := &OrderBook{
		Asks: []Order{
			{OrderPrice: OrderPrice{Base: "0.000001 VESTS", Quote: "1.000 STEEM"}},
		},
	}
	fh := loadFeedHistory(t)
	dgp := loadDGP(t)
	_, err := ComputePrices(ob, fh, dgp)
	if err == nil {
		t.Fatal("ComputePrices with VESTS order: expected error, got nil (would silently mislabel symbol)")
	}
}

// TestComputePrices_RejectsInconsistentSymbols 修复 #1：所有 order 的结果 symbol
// 必须一致。混入一个结果非 SBD 的 order 必须报错。
func TestComputePrices_RejectsInconsistentSymbols(t *testing.T) {
	ob := &OrderBook{
		Asks: []Order{
			{OrderPrice: OrderPrice{Base: "1.000 SBD", Quote: "1.000 STEEM"}}},
		Bids: []Order{
			// Convert(1.000 STEEM) 结果是 VESTS，与 asks 的 SBD 不一致。
			{OrderPrice: OrderPrice{Base: "0.000001 VESTS", Quote: "1.000 STEEM"}}},
	}
	fh := loadFeedHistory(t)
	dgp := loadDGP(t)
	if _, err := ComputePrices(ob, fh, dgp); err == nil {
		t.Error("ComputePrices with mixed SBD/VESTS orders: expected error")
	}
}

// TestComputePrices_SbdSymbolNotHardcoded 修复 #1 的正向验证：纯 STEEM/SBD 对的
// order book 正常工作，且结果 symbol 是从 Convert 推断的 "SBD"（而非硬编码）。
// 关键断言：SteemSbd.Symbol 必须等于 "SBD"，且数值与手算一致。
func TestComputePrices_SbdSymbolNotHardcoded(t *testing.T) {
	// 反向 order（base=STEEM quote=SBD）也能正确推导出 SBD 结果。
	ob := &OrderBook{
		Asks: []Order{
			// Convert(1.000 STEEM): STEEM 匹配 base → 结果 quote=SBD
			// amount = 1000(steem) × 1000(sbd) / 1000(steem) = 1000 SBD atoms
			{OrderPrice: OrderPrice{Base: "1.000 STEEM", Quote: "1.000 SBD"}}},
	}
	fh := loadFeedHistory(t)
	dgp := loadDGP(t)
	res, err := ComputePrices(ob, fh, dgp)
	if err != nil {
		t.Fatalf("ComputePrices: %v", err)
	}
	if res.SteemSbd.Symbol != "SBD" {
		t.Errorf("SteemSbd.Symbol = %q want SBD (must be inferred, not hardcoded)", res.SteemSbd.Symbol)
	}
	if res.SteemSbd.Amount != 1000 {
		t.Errorf("SteemSbd.Amount = %d want 1000", res.SteemSbd.Amount)
	}
}
