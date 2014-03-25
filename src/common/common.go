package common

type Account_info struct {
	Total                 float64
	net_asset             float64
	available_cny_display float64
	available_btc_display float64
	frozen_cny_display    float64
	frozen_btc_display    float64
	loan_cny_display      float64
	loan_btc_display      float64
}

type TradeAPI interface {
	AnalyzeKLine(peroid int) (ret bool)
	Buy(price, amount string) bool
	Sell(price, amount string) bool
	GetTradePrice(tradeDirection string) string
	Get_account_info() (account_info Account_info, ret bool)
}
