package api

// OrderPrice models the order_price field of a limit order, expressed as a
// ratio of two asset strings (e.g. base "1.000 SBD", quote "1.000 STEEM").
// Used by conveyor/src/price.ts.
type OrderPrice struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

// Order models a single entry in a database_api.get_order_book response (one
// of the asks or bids).
type Order struct {
	OrderPrice OrderPrice `json:"order_price"`
}

// OrderBook models a database_api.get_order_book response. conveyor averages
// the order prices across both sides to derive the STEEM<>SBD market price.
type OrderBook struct {
	Asks []Order `json:"asks"`
	Bids []Order `json:"bids"`
}

// CurrentMedianHistoryPrice models a witness-reported price entry, expressed
// as a base/quote asset ratio. Used by conveyor/src/price.ts to derive the
// STEEM<>USD price.
type CurrentMedianHistoryPrice struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

// FeedHistory models a database_api.get_feed_history response. conveyor reads
// the last entry of price_history as the current witness price feed.
type FeedHistory struct {
	PriceHistory []CurrentMedianHistoryPrice `json:"price_history"`
}

// ToPrice 解析 OrderPrice 的 base/quote 字符串为精确的 Price 原语。
// 是 wire string-asset 字段（"1.000 SBD"）与 Price 之间的桥梁。
// 调用方拿到 Price 后即可做零浮点的 Convert/Compare 运算。
func (o OrderPrice) ToPrice() (Price, error) {
	return ParsePrice(o.Base, o.Quote)
}

// ToPrice 解析 CurrentMedianHistoryPrice 的 base/quote 字符串为精确的 Price。
// conveyor 用 feed_history 最后一项推导 STEEM<>USD 价格。
func (c CurrentMedianHistoryPrice) ToPrice() (Price, error) {
	return ParsePrice(c.Base, c.Quote)
}
