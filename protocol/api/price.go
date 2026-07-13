package api

import (
	"math/big"

	"github.com/pkg/errors"
)

// Price 表示两个 asset 的比值，复刻 steemd price {base, quote}（asset.hpp）。
// 链上语义：base 个 Quote.Symbol 兑换 quote 个 Base.Symbol；例如
//
//	Price{Base: "1.000 SBD", Quote: "1.000 STEEM"}
//
// 表示"1 STEEM 值 1 SBD"。
//
// 价格运算全程使用 math/big.Int（标准库），覆盖 128 位中间值且免手写 uint128，
// 并复刻 steemd "先乘后除 + 溢出检查"的语义，与链端逐位一致。
// 零浮点、零第三方依赖。
type Price struct {
	Base, Quote Asset
}

// ParsePrice 解析两个 asset 字符串并构造 Price。
// 校验对齐 steemd price::validate：
//   - Base.Symbol != Quote.Symbol（否则比值退化为 1，无意义）。
//   - Base 与 Quote 均非零（分母为零无意义，分子为零说明价格未初始化）。
func ParsePrice(base, quote string) (Price, error) {
	b, err := ParseAsset(base)
	if err != nil {
		return Price{}, errors.Wrapf(err, "parse price base %q", base)
	}
	q, err := ParseAsset(quote)
	if err != nil {
		return Price{}, errors.Wrapf(err, "parse price quote %q", quote)
	}
	if b.Symbol == q.Symbol {
		return Price{}, errors.Errorf(
			"price base and quote must have different symbols, both %s", b.Symbol,
		)
	}
	if b.Amount == 0 || q.Amount == 0 {
		return Price{}, errors.New("price base and quote must be non-zero")
	}
	return Price{Base: b, Quote: q}, nil
}

// Convert 复刻 steemd asset*price（asset.cpp:285-302）：
//
//	r = (a.Amount * 分子.Amount) / 分母.Amount
//
// 其中分子/分母取决于 a 的符号匹配哪一边：
//   - a.Symbol == Base.Symbol：结果符号 = Quote.Symbol，分子 = Quote.Amount，分母 = Base.Amount
//   - a.Symbol == Quote.Symbol：结果符号 = Base.Symbol，分子 = Base.Amount，分母 = Quote.Amount
//
// 实现要点（与链端逐位一致）：
//   - 用 math/big.Int 计算，中间乘积可超 int64（128 位）。
//   - Quo（向零截断）而非 Div（向下取整），与 C++ 整数除法一致。
//   - 溢出检查：除法后结果必须 fit int64 且 >= 0（对齐 result.hi == 0 断言）。
//   - 返回的 Asset 是结果 symbol 的整数原子单位。
func (p Price) Convert(a Asset) (Asset, error) {
	var (
		numeratorAmount   int64
		denominatorAmount int64
		resultSymbol      string
	)
	switch a.Symbol {
	case p.Base.Symbol:
		resultSymbol = p.Quote.Symbol
		numeratorAmount = p.Quote.Amount
		denominatorAmount = p.Base.Amount
	case p.Quote.Symbol:
		resultSymbol = p.Base.Symbol
		numeratorAmount = p.Base.Amount
		denominatorAmount = p.Quote.Amount
	default:
		return Asset{}, errors.Errorf(
			"asset symbol %s matches neither price side (%s / %s)",
			a.Symbol, p.Base.Symbol, p.Quote.Symbol,
		)
	}
	if denominatorAmount == 0 {
		return Asset{}, errors.New("price convert: zero denominator")
	}

	// a.Amount 可能为 0，结果自然为 0；但分子/分母已在校验中保证非零。
	num := new(big.Int).Mul(big.NewInt(a.Amount), big.NewInt(numeratorAmount))
	// Quo = 向零截断（truncated division），与 C++ 整数除法语义一致。
	res := new(big.Int).Quo(num, big.NewInt(denominatorAmount))

	// 链端断言：result.hi == 0 即结果 fit uint64。本实现更严：必须 fit [0, MaxSatoshis]。
	if res.Sign() < 0 {
		return Asset{}, errors.Errorf("price convert result negative: %s", res.String())
	}
	if !res.IsUint64() || res.Uint64() > uint64(MaxSatoshis) {
		return Asset{}, errors.Errorf("price convert result overflow int62: %s", res.String())
	}
	return Asset{Amount: int64(res.Uint64()), Symbol: resultSymbol}, nil
}

// Compare 复刻 steemd price 比较（asset.cpp:263-283）：
//
//	a < b  ⟺  a.Base*b.Quote < b.Base*a.Quote   （交叉相乘后比分子）
//
// 返回 -1/0/1。零除法零浮点。
// 用 math/big.Int 做交叉乘，避免 int64 乘法溢出。
//
// 注意：steemd 要求两个 price 描述同一种"兑换方向"才可比较；
// 本实现不做方向归一化，由调用方保证语义一致（与 steemd operator< 前置条件一致）。
func (p Price) Compare(q Price) int {
	// a.Base * b.Quote
	lhs := new(big.Int).Mul(big.NewInt(p.Base.Amount), big.NewInt(q.Quote.Amount))
	// b.Base * a.Quote
	rhs := new(big.Int).Mul(big.NewInt(q.Base.Amount), big.NewInt(p.Quote.Amount))
	return lhs.Cmp(rhs)
}

// Invert 复刻 steemd operator~（取倒数）：交换 Base 与 Quote。
// Invert().Invert() == 原 Price（round-trip 恒等）。
func (p Price) Invert() Price {
	return Price{Base: p.Quote, Quote: p.Base}
}
