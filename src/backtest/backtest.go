package main

import (
	. "common"
	. "config"
	"fmt"
	"logger"
	"os"
	"simulate"
	"strategy"
	"time"
	. "util"
)

const (
	ProtocolVersion = 1
	ShiftNumber     = 60
	MaxKLineLength  = 120
	BackTestCNY     = 20000
)

func init() {
	os.Mkdir(ROOT+"/test/", 0777)
}

func main() {
	backtesting()
	return
}

func backtesting() {
	fmt.Println("back testing begin...")

	SetBacktest(true)
	SimAccount = make(map[string]string)
	SimAccount["CNY"] = IntegerToString(BackTestCNY)
	SimAccount["BTC"] = "0"
	SimAccount["LTC"] = "0"
	SaveSimulate()

	btEnd := time.Now()
	btEnd = time.Date(btEnd.Year(), btEnd.Month(), btEnd.Day(), btEnd.Hour(), btEnd.Minute(), 0, 0, time.Local)
	//btEnd = btEnd.AddDate(0, 0, -1)
	//btStart := btEnd.AddDate(0, -3, 0)

	btEnd = btEnd.Add(0 - time.Second*3600*2)
	btStart := btEnd.Add(0 - time.Second*3600*24*30)

	day := time.Second * 3600 * 24
	dayshift := 0
	btEnd = btEnd.Add(0 - day*time.Duration(dayshift))
	btStart = btStart.Add(0 - day*time.Duration(dayshift))

	var symbolId string
	if Option["symbol"] == "btc_cny" {
		symbolId = "btccny"
	} else {
		symbolId = "ltccny"
	}

	fmt.Println("回测开始时间", btStart.Format("2006-01-02 15:04:05"))
	fmt.Println("回测结束时间", btEnd.Format("2006-01-02 15:04:05"))

	period := StringToInteger(Option["tick_interval"])

	var records []Record
	var err error
	datacenter := Option["datacenter"]
	if datacenter == "huobi" {
		records, err = GetKLine(btStart, btEnd, period, symbolId)
	} else {
		records, err = OkcoinKLine(period, symbolId)
	}
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(len(records))
	}

	if Option["tick_interval"] != "1" {
		// check exception data
		for i := 0; i < len(records)-1; i++ {
			if strategy.CheckException(records[i], records[i+1]) == false {
				logger.Errorln("detect exception data ",
					records[i].Close, records[i+1].Close, records[i+1].Volumn)
				return
			}
		}
	}

	filepath := btStart.Format("2006-01-02") + "_" + btEnd.Format("2006-01-02")
	filepath = "\\test\\" + filepath + "_" + IntegerToString(period) + "_" + symbolId
	filepath = filepath + "_" + Option["strategy"] + ".txt"
	logger.SetBacktestFile(filepath)
	DeleteFile(ROOT + filepath)

	for i := ShiftNumber; i < len(records); i++ {
		var rec []Record
		if i < MaxKLineLength {
			rec = records[0:i]
		} else {
			rec = records[i-MaxKLineLength : i]
		}
		SetBtTime(rec[len(rec)-1].Time)
		SetBtPrice(records[i].Close)
		strategy.Tick(simulate.NewSimulate(), rec)
	}

	var coin float64
	var total float64
	if symbolId == "btccny" {
		coin = StringToFloat(SimAccount["BTC"])
		total = coin*records[len(records)-1].Close + StringToFloat(SimAccount["CNY"])
	} else {
		coin = StringToFloat(SimAccount["LTC"])
		total = coin*records[len(records)-1].Close + StringToFloat(SimAccount["CNY"])
	}
	logger.Backtestf("资金总计：%f    盈亏比：%f \n", total, (total-BackTestCNY)/BackTestCNY*100)

	//peroids := []int{1, 5, 15, 30, 60, 100}
	//for _, peroid := range peroids {
	//	if huobi.AnalyzeKLinePeroid("btc_cny", peroid) == true {
	//	} else {
	//		logger.Errorln("TradeKLine failed.")
	//	}
	//}
	//fmt.Println("生成 1/5/15/30/60分钟及1天 周期的后向测试报告于log/reportxxx.log文件中,请查看")

	fmt.Println("back testing end ...")
}
