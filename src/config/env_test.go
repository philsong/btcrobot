/*
  btcbot is a Bitcoin trading bot for HUOBI.com written
  in golang, it features multiple trading methods using
  technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package config

import (
	"testing"
)

func Test_SetENV(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", "development"},
		{"not_development", "not_development"},
	}

	for _, test := range tests {
		setENV(test.in)
		if Env != test.out {
			expect(t, Env, test.out)
		}
	}
}

func Test_Root(t *testing.T) {
	if len(Root) == 0 {
		t.Errorf("Expected root path will be set")
	}
}
