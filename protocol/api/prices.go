package api

import (
	"github.com/pkg/errors"
)

// PricesResult 对应 conveyor src/price.ts get_prices 的三项输出，但类型是
// 精确的 Price/Asset 而非 float64。conveyor 侧按需自行转 float64 用于展示。
//
// 字段语义对齐 conveyor：
//   - SteemSbd: 单位 STEEM 折算的 SBD 原子值（整数平均）。
//   - SteemUsd: feed_history 最后一项的价格比值（STEEM<>SBD 视作 USD 近似）。
//   - SteemVest: total_vesting_fund_steem / total_vesting_shares 的 vesting 价格。
type PricesResult struct {
	// SteemSbd 是 order book 各 order Convert(1.000 STEEM) 得到的 SBD 原子值的
	// 整数算术平均（sum/count）。原 conveyor TS 用 float64 算术平均，此处换成
	// int64 整数平均以保持精确；展示层转 float64 下沉到 conveyor。
	SteemSbd Asset
	// SteemUsd 直接取 feed_history.price_history 最后一项的比值，不做转换。
	SteemUsd Price
	// SteemVest = ParsePrice(total_vesting_fund_steem, total_vesting_shares)。
	SteemVest Price
}

// ComputePrices 从 wire 类型计算三项精确价格，对应 conveyor src/price.ts get_prices。
//
// 入参均为指针：nil 表示对应数据源缺失，相应跳过该字段并返回 error
// （三项均依赖独立数据源，缺一即视为不可用）。这与 conveyor 行为一致：
// 任一关键来源缺失都使价格不可信。
//
// 全程零浮点；任何错误（空订单簿、symbol 不匹配等）立即向上返回。
func ComputePrices(ob *OrderBook, fh *FeedHistory, dgp *DynamicGlobalProperties) (PricesResult, error) {
	var result PricesResult

	// --- SteemSbd: order book 平均 ---
	if ob == nil {
		return PricesResult{}, errors.New("compute prices: order book is nil")
	}
	totalOrders := len(ob.Asks) + len(ob.Bids)
	if totalOrders == 0 {
		return PricesResult{}, errors.New("compute prices: order book is empty")
	}
	oneSteem, err := ParseAsset("1.000 STEEM")
	if err != nil {
		// 不可达：常量字符串。
		return PricesResult{}, errors.Wrap(err, "compute prices: parse 1.000 STEEM")
	}

	// 累加 Convert(1.000 STEEM) 的 SBD 原子值。
	// 用 int64 累加；单笔上限 MaxSatoshis ≈ 4.6e18，订单数远小于 2^63/MaxSatoshis
	// 时安全。若出现溢出说明订单数异常巨大或单价异常，按链端语义视为错误。
	var sumSbd int64
	convert := func(o Order) error {
		p, err := o.OrderPrice.ToPrice()
		if err != nil {
			return errors.Wrapf(err, "order price base=%q quote=%q", o.OrderPrice.Base, o.OrderPrice.Quote)
		}
		conv, err := p.Convert(oneSteem)
		if err != nil {
			return errors.Wrapf(err, "convert 1.000 STEEM via base=%q quote=%q", o.OrderPrice.Base, o.OrderPrice.Quote)
		}
		// 溢出保护：累加前检查。
		if conv.Amount > 0 && sumSbd > MaxSatoshis-conv.Amount {
			return errors.Errorf("steem_sbd sum overflow at %d + %d", sumSbd, conv.Amount)
		}
		sumSbd += conv.Amount
		return nil
	}
	for _, o := range ob.Asks {
		if err := convert(o); err != nil {
			return PricesResult{}, err
		}
	}
	for _, o := range ob.Bids {
		if err := convert(o); err != nil {
			return PricesResult{}, err
		}
	}

	// 找一个结果 symbol（从首个有效 order 推断；已在 convert 中保证一致）。
	resultSymbol := "SBD"
	// 整数算术平均（向零截断），贴近 conveyor 原 TS 算术平均语义。
	avgSbd := sumSbd / int64(totalOrders)
	result.SteemSbd = Asset{Amount: avgSbd, Symbol: resultSymbol}

	// --- SteemUsd: feed_history 最后一项 ---
	if fh == nil {
		return PricesResult{}, errors.New("compute prices: feed history is nil")
	}
	if len(fh.PriceHistory) == 0 {
		return PricesResult{}, errors.New("compute prices: feed history price_history is empty")
	}
	last := fh.PriceHistory[len(fh.PriceHistory)-1]
	usdPrice, err := last.ToPrice()
	if err != nil {
		return PricesResult{}, errors.Wrapf(err, "feed history last entry base=%q quote=%q", last.Base, last.Quote)
	}
	result.SteemUsd = usdPrice

	// --- SteemVest: total_vesting_fund_steem / total_vesting_shares ---
	if dgp == nil {
		return PricesResult{}, errors.New("compute prices: dynamic global properties is nil")
	}
	vestPrice, err := ParsePrice(dgp.TotalVestingFundSteem, dgp.TotalVestingShares)
	if err != nil {
		return PricesResult{}, errors.Wrapf(err,
			"vesting price fund=%q shares=%q",
			dgp.TotalVestingFundSteem, dgp.TotalVestingShares,
		)
	}
	result.SteemVest = vestPrice

	return result, nil
}
