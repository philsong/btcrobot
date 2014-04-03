package common

//trade interface type and method

type UserMoney struct {
	Available_cny string
	Available_btc string
	Available_ltc string
	Frozen_cny    string
	Frozen_btc    string
	Frozen_ltc    string
}

type TradeAPI interface {
	AnalyzeKLine(peroid int) bool
	Buy(price, amount string) bool
	Sell(price, amount string) bool
	GetTradePrice(tradeDirection string) string
	Get_account_info() (UserMoney, bool)
	GetOrderBook(string) bool
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
