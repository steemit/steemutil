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
