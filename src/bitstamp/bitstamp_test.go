package bitstamp

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func Init() (api *Api) {
	filename := os.ExpandEnv("$BITSTAMP_CONFIG")
	if filename == "" {
		fmt.Println("Please set $BITSTAMP_CONFIG to a proper configuration file")
		syscall.Exit(0)
	}
	api, err := NewFromConfig(filename)
	if err != nil {
		panic(err)
	}
	return
}

func TestTicker(t *testing.T) {
	api := Init()
	ticker, err := api.GetTicker()
	if err != nil {
		t.Errorf("Could not fetch ticker :", err)
	}
	if ticker.Last == 0 {
		t.Errorf("Ticker probably wrongly filled")
	}
}

func TestOrderBook(t *testing.T) {
	api := Init()
	orderbook, err := api.GetOrderBook()
	if err != nil {
		t.Errorf("Could not fetch orderbook :", err)
	}
	if orderbook.Orders[0].Price == 0. {
		t.Errorf("Orderbook probably wrongly filled")
	}
}

func TestTrades(t *testing.T) {
	api := Init()
	trades, err := api.GetTrades()
	if err != nil {
		t.Errorf("Could not fetch trades :", err)
	}
	if len(trades) == 0 || trades[0].Price == 0. {
		t.Errorf("trades probably wrongly filled")
	}
}

func TestTradesParams(t *testing.T) {
	api := Init()
	trades, err := api.GetTradesParams(1, 10, "desc")
	if err != nil {
		t.Errorf("Could not fetch trades with params:", err)
	}
	if len(trades) == 0 || trades[0].Price == 0. {
		t.Errorf("trades with params probably wrongly filled")
	}
}
