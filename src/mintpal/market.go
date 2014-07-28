package mintpal

import (
	"errors"
	"fmt"
)

var (
	PrefixAddr          = "https://api.mintpal.com/v2/"
	MarketSummaryAddr   = PrefixAddr + "market/summary/"
	MarketStatsAddr     = PrefixAddr + "market/stats/"
	MarketTradesAddr    = PrefixAddr + "market/trades/"
	MarketOrdersAddr    = PrefixAddr + "market/orders/"
	MarketChartDataAddr = PrefixAddr + "market/chartdata/"
)

func (p *manager) GetMarketSummary(exchange string) (summary *MarketSummary, err error) {
	if exchange != "BTC" && exchange != "LTC" && exchange != "" {
		err = errors.New("invalid exchange")
		return
	}

	summary = new(MarketSummary)
	err = get(MarketSummaryAddr+exchange, summary)
	if err != nil {
		fmt.Println(err)
	}
	debug(summary)
	return
}

func (p *manager) GetMarketStats(code, exchange string) (stats *MarketStats, err error) {
	if code == "" || (exchange != "BTC" && exchange != "LTC") {
		err = errors.New("invalid exchange")
		return
	}
	stats = new(MarketStats)
	err = get(MarketStatsAddr+code+"/"+exchange, stats)
	if err != nil {
		fmt.Println(err)
	}
	debug(stats)
	return
}

func (p *manager) GetMarketTrades(code, exchange string) (trades *MarketTrades, err error) {
	if code == "" || (exchange != "BTC" && exchange != "LTC") {
		err = errors.New("invalid exchange")
		return
	}
	trades = new(MarketTrades)
	err = get(MarketTradesAddr+code+"/"+exchange, trades)
	if err != nil {
		fmt.Println(err)
	}
	debug(trades)
	return
}

func (p *manager) GetMarketOrders(code, exchange, tp string) (orders *MarketOrders, err error) {
	if code == "" || (exchange != "BTC" && exchange != "LTC") || (tp != "BUY" && tp != "SELL") {
		err = errors.New("invalid exchange")
		return
	}

	orders = new(MarketOrders)
	err = get(MarketOrdersAddr+code+"/"+exchange+"/"+tp, orders)
	if err != nil {
		fmt.Println(err)
	}
	debug(orders)
	return
}

func (p *manager) GetMarketChartData(code, exchange, period string) (chart *MarketChartData, err error) {
	if code == "" || (exchange != "BTC" && exchange != "LTC") ||
		(period != "6hh" && period != "1DD" && period != "3DD" &&
			period != "7DD" && period != "MAX") {
		err = errors.New("invalid exchange")
		return
	}

	chart = new(MarketChartData)
	err = get(MarketChartDataAddr+code+"/"+exchange+"/"+period, chart)
	if err != nil {
		fmt.Println(err)
	}
	debug(chart)
	return
}
