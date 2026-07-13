package api

import (
	"math/bits"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// 精度常量（对齐 steemd STEEM_PRECISION_*）。
// STEEM 与 SBD 是 3 位小数，VESTS 是 6 位小数。
const (
	PrecisionSteem = 3
	PrecisionSbd   = 3
	PrecisionVests = 6
)

// MaxSatoshis 对齐 STEEM_MAX_SATOSHIS = 2^62-1（steemd config.hpp:241）。
// 链上任何 asset 的原子单位绝对值都不允许超过此值。
const MaxSatoshis int64 = 4611686018427387903 // (1 << 62) - 1

// Asset 表示 Steem 链上资产：int64 原子单位 + symbol。
// 精度（小数位数）从 symbol 派生，仅在 ParseAsset/String 时参与，不单独存储，
// 这样 Asset 始终是"链端真相"的精确表达，杜绝 float 中转。
//
// 与 protocol 根包的 Asset（用于交易二进制序列化）不同：本类型是
// protocol/api 层的"解释型"原语，面向 wire JSON 字段（如 "1.500 STEEM"）。
type Asset struct {
	Amount int64
	Symbol string // STEEM | SBD | VESTS | TESTS | TBD
}

// symbolPrecision 返回 symbol 的小数位数；非法 symbol 报错。
// 对齐 steemd asset_symbol::decimals()。
func symbolPrecision(symbol string) (int, error) {
	switch strings.ToUpper(symbol) {
	case "STEEM", "SBD", "TESTS", "TBD":
		return PrecisionSteem, nil
	case "VESTS":
		return PrecisionVests, nil
	default:
		return 0, errors.Errorf("unknown asset symbol: %q", symbol)
	}
}

// ParseAsset 解析 "1.500 STEEM" 形式的 asset 字符串为 Asset。
//
// 流程（全程零 float）：
//   - 按单空格切分 [amountStr, symbol]；symbol 转大写后校验合法性。
//   - 按 '.' 切分整数/小数部分，拼成纯数字串后 strconv.ParseInt（不经 float 中转）。
//   - 校验小数位数 == symbolPrecision(symbol)（VESTS 给 3 位小数即报错）。
//   - 校验 0 <= amount <= MaxSatoshis（对齐 steemd asset::validate）。
//
// 不接受负数：steemd wire 字段在本应用场景下均为非负（余额、供给量等）。
func ParseAsset(s string) (Asset, error) {
	parts := strings.Split(strings.TrimSpace(s), " ")
	if len(parts) != 2 {
		return Asset{}, errors.Errorf("invalid asset format: %q (expected 'amount symbol')", s)
	}

	amountStr := parts[0]
	symbol := strings.ToUpper(parts[1])

	prec, err := symbolPrecision(symbol)
	if err != nil {
		return Asset{}, err
	}

	// 切出整数与小数部分，纯字符串操作，杜绝 float。
	var neg bool
	intPart, fracPart := amountStr, ""
	if strings.HasPrefix(amountStr, "-") {
		neg = true
		intPart = amountStr[1:]
	} else if strings.HasPrefix(amountStr, "+") {
		intPart = amountStr[1:]
	}
	if dot := strings.Index(intPart, "."); dot >= 0 {
		intPart, fracPart = intPart[:dot], intPart[dot+1:]
	}

	// 校验整数/小数部分都是纯数字（空串也算非法，防止 "." / "1." / ".5"）。
	if intPart == "" && fracPart == "" {
		return Asset{}, errors.Errorf("invalid asset amount: %q", amountStr)
	}
	for _, ch := range intPart + fracPart {
		if ch < '0' || ch > '9' {
			return Asset{}, errors.Errorf("invalid asset amount: %q", amountStr)
		}
	}

	// 校验小数位数与 symbol 精度一致。
	if len(fracPart) != prec {
		return Asset{}, errors.Errorf(
			"asset %q has %d decimal places, but symbol %s requires %d",
			s, len(fracPart), symbol, prec,
		)
	}

	// 拼成纯数字串解析（去前导零交给 ParseInt 处理）。
	combined := intPart + fracPart
	if combined == "" {
		combined = "0"
	}
	parsed, err := strconv.ParseInt(combined, 10, 64)
	if err != nil {
		// 大概率是溢出 int64。
		return Asset{}, errors.Wrapf(err, "asset amount overflow: %q", amountStr)
	}
	if neg {
		parsed = -parsed
	}

	// 链端校验：|amount| <= MaxSatoshis，且本场景限定非负。
	if parsed < 0 || parsed > MaxSatoshis {
		return Asset{}, errors.Errorf(
			"asset amount %d out of range [0, %d]", parsed, MaxSatoshis,
		)
	}

	return Asset{Amount: parsed, Symbol: symbol}, nil
}

// Precision 返回该 asset symbol 的小数位数。
//
// 前置条件：Symbol 必须是经 ParseAsset 校验过的合法 symbol。直接用 Asset{}
// 字面量构造非法 symbol 属于编程错误，对未知 symbol 直接 panic 而非静默退化
// （静默返回 0 会让 String() 输出丢精度且无报错，难排查）。
func (a Asset) Precision() int {
	prec, err := symbolPrecision(a.Symbol)
	if err != nil {
		panic(err)
	}
	return prec
}

// String 从 int64 重建 "1.500 STEEM" 字符串。
// 纯整数除法 + 取模重建小数点，避免 strconv.FormatFloat 的 round-half-even。
// 不经任何 float 中转，与 ParseAsset 严格互逆（round-trip）。
func (a Asset) String() string {
	prec := a.Precision()
	amount := a.Amount
	negative := amount < 0
	if negative {
		amount = -amount
	}

	// 整数部分与小数（原子）部分。
	var intPart, fracPart string
	if prec == 0 {
		intPart = strconv.FormatInt(amount, 10)
	} else {
		divisor := int64(1)
		for i := 0; i < prec; i++ {
			divisor *= 10
		}
		intVal := amount / divisor
		fracVal := amount % divisor
		intPart = strconv.FormatInt(intVal, 10)
		// 小数部分左侧补零到 prec 位。
		fracPart = strconv.FormatInt(fracVal, 10)
		if len(fracPart) < prec {
			fracPart = strings.Repeat("0", prec-len(fracPart)) + fracPart
		}
	}

	out := intPart
	if prec > 0 {
		out = out + "." + fracPart
	}
	if negative {
		out = "-" + out
	}
	return out + " " + a.Symbol
}

// Add 同 symbol 加法（带溢出检查），对齐 steemd asset::operator+ 的 safe<int64_t> 语义。
// 不同 symbol 报错（链端 assert base.asset == addend.asset）。
//
// 用 math/bits.Add64 做符号无关的精确溢出检查：拿到进位标志即可判定，无需
// 针对正负组合分别写条件（避免"只在 b>0 时检查"一类遗漏）。
func (a Asset) Add(b Asset) (Asset, error) {
	if a.Symbol != b.Symbol {
		return Asset{}, errors.Errorf(
			"cannot add assets of different symbols: %s vs %s", a.Symbol, b.Symbol,
		)
	}
	sum, carry := bits.Add64(uint64(a.Amount), uint64(b.Amount), 0)
	if carry != 0 || sum > uint64(MaxSatoshis) {
		return Asset{}, errors.Errorf("asset addition overflow: %d + %d", a.Amount, b.Amount)
	}
	return Asset{Amount: int64(sum), Symbol: a.Symbol}, nil
}

// Sub 同 symbol 减法（带溢出检查），对齐 steemd asset::operator-。
//
// 结果为负即报错：与 ParseAsset 的"非负"不变式保持一致——所有 Asset（无论来源）
// 始终代表链上合法值 [0, MaxSatoshis]。需要"净流入"等可能为负的场景由调用方在
// Sub 前比较大小、自行交换操作数；这样 Asset 的 round-trip 契约（String ↔ ParseAsset）
// 对所有合法值都严格成立。
func (a Asset) Sub(b Asset) (Asset, error) {
	if a.Symbol != b.Symbol {
		return Asset{}, errors.Errorf(
			"cannot subtract assets of different symbols: %s vs %s", a.Symbol, b.Symbol,
		)
	}
	if b.Amount > a.Amount {
		return Asset{}, errors.Errorf(
			"asset subtraction would be negative: %d - %d", a.Amount, b.Amount,
		)
	}
	return Asset{Amount: a.Amount - b.Amount, Symbol: a.Symbol}, nil
}
