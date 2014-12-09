package mintpal

// Market
type MarketSummary struct {
	Status string
	Count  int
	Data   []marketSummary
}

type MarketStats struct {
	Status string
	Data   marketSummary
}

type marketSummary struct {
	MarketID       string `json:"market_id"`
	Coin           string
	Code           string
	Exchange       string
	LastPrice      string `json:"last_price"`
	YesterdayPrice string `json:"yesterday_price"`
	Change         string
	High           string `json:"24hhigh"`
	Low            string `json:"24hlow"`
	Vol            string `json:"24hvol"`
	TopBid         string `json:"top_bid"`
	TopAsk         string `json:"top_ask"`
}

type MarketTrades struct {
	Status string
	Count  int
	Data   []marketTrades
}

type marketTrades struct {
	Type   string
	Price  string
	Amount string
	Total  string
	Time   string
}

type MarketOrders struct {
	Status string
	Count  int
	Data   []marketOrders
}

type marketOrders struct {
	Price  string
	Amount string
	Total  string
}

type MarketChartData struct {
	Status string
	Count  int
	Data   []marketChartData
}

type marketChartData struct {
	Date           string
	Open           string
	Close          string
	High           string
	Low            string
	ExchangeVolume string `json:"exchange_volume"`
	CoinVolume     string `json:"coin_volume"`
}

// trading

type TradingOrders struct {
	Status string
	count  int
	Date   []tradingOrder
}

type TradingSingleOrder struct {
	Status string
	Data   tradingOrder
}

type tradingOrder struct {
	OrderID       string `json:"order_id"`
	Market        string
	Type          string
	Price         string
	Amount        string
	Total         string
	Fee           string
	NetTotal      string
	Time          string
	TimeFormatted string
}

type TradingAddOrderReq struct {
	Coin     string
	Exchange string
	Price    int
	Amount   int
	Type     bool
}

type TradingAddOrderResp struct {
	Status  string
	Message string
	Data    []tradingOrder
}

type TradingCancelOrder struct {
	Status  string
	Message string
}

type TradingTrades struct {
	Status string
	Count  int
	Data   []tradingOrder
}

type TradingSingleTrade struct {
	Status string
	Data   tradingOrder
}
