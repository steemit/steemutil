package operation

// FC_REFLECT( steemit::chain::feed_publish_operation,
//             (publisher)
//             (exchange_rate) )

type FeedPublishOperation struct {
	Publisher    string `json:"publisher"`
	ExchangeRate struct {
		Base  string `json:"base"`
		Quote string `json:"quote"`
	} `json:"exchange_rate"`
}

func (op *FeedPublishOperation) Type() OpType {
	return TypeFeedPublish
}

func (op *FeedPublishOperation) Data() interface{} {
	return op
}
