package common

type UserMoney struct {
	Available_cny string
	Available_btc string
	Available_ltc string
	Frozen_cny    string
	Frozen_btc    string
	Frozen_ltc    string
}

type TradeAPI interface {
	AnalyzeKLine(peroid int) (ret bool)
	Buy(price, amount string) bool
	Sell(price, amount string) bool
	GetTradePrice(tradeDirection string) string
	Get_account_info() (userMoney UserMoney, ret bool)
}
