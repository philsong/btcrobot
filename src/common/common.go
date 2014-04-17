package common

//trade interface type and method

type Account struct {
	Available_cny string
	Available_btc string
	Available_ltc string
	Frozen_cny    string
	Frozen_btc    string
	Frozen_ltc    string
}

type Record struct {
	TimeStr string
	Time    int64
	Open    float64
	High    float64
	Low     float64
	Close   float64
	Volumn  float64
}

type _MarketOrder struct {
	Price  float64 //价格
	Amount float64 //委单量
}

//price from high to low: asks[0] > .....>asks[9] > bids[0] > ......> bids[9]
type OrderBook struct {
	Asks [10]_MarketOrder //sell
	Bids [10]_MarketOrder //buy
}

type Order struct {
	Id          int
	Amount      float64
	Deal_amount float64
}

type MarketAPI interface {
	GetKLine(peroid int) (ret bool, records []Record)
}

type TradeAPI interface {
	Buy(price, amount string) string
	Sell(price, amount string) string
	GetOrder(order_id string) (ret bool, order Order)
	CancelOrder(order_id string) bool
	GetAccount() (Account, bool)
	GetOrderBook() (ret bool, orderBook OrderBook)
}
