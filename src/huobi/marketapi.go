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

package huobi

import (
	"bufio"
	. "config"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"math/rand"
	"net/http"
	"os"
	"strategy"
	"strconv"
	"strings"
)

/*
	txcode := map[int]string{
		0:  `买单已委托，<a href="/trade/index.php?a=delegation">查看结果</a>`,
		2:  `没有足够的人民币`,
		10:	`没有足够的比特币`,
		16: `您需要登录才能继续`,
		17: `没有权限`,
		42:	`该委托已经取消，不能修改`,
		44:	`交易价钱太低`,
		56:`卖出价格不能低于限价的95%`}

	logger.Traceln(txcode[m.Code])
*/

func (w *Huobi) AnalyzeKLinePeroid(peroid int) (ret bool) {
	req, err := http.NewRequest("GET", fmt.Sprintf(Config["trade_kline_url"], peroid, rand.Float64()), nil)
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Referer", Config["trade_flash_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	if w.client == nil {
		w.client = &http.Client{nil, nil, nil}
	}
	resp, err := w.client.Do(req)
	if err != nil {
		logger.Errorln(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Errorln("read the http stream failed")
				return false
			} else {
				body = string(bodyByte)

				ioutil.WriteFile(fmt.Sprintf("cache/TradeKLine_%03d.data", peroid), bodyByte, 644)
			}
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		ret := strings.Contains(body, "您需要登录才能继续")
		if ret {
			logger.Traceln("您需要登录才能继续")
			return false
		} else {
			return w.analyzePeroidLine(fmt.Sprintf("cache/TradeKLine_%03d.data", peroid), body)
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *Huobi) AnalyzeKLineMinute() (ret bool) {
	req, err := http.NewRequest("GET", Config["trade_fenshi"], nil)
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Referer", Config["trade_flash_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	if w.client == nil {
		w.client = &http.Client{nil, nil, nil}
	}
	resp, err := w.client.Do(req)
	if err != nil {
		logger.Errorln(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Errorln("read the http stream failed")
				return false
			} else {
				body = string(bodyByte)

				ioutil.WriteFile(fmt.Sprintf("cache/TradeKLine_minute.data"), bodyByte, 644)
			}
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		ret := strings.Contains(body, "您需要登录才能继续")
		if ret {
			logger.Traceln("您需要登录才能继续")
			return false
		} else {
			return w.analyzeMinuteLine(fmt.Sprintf("cache/TradeKLine_minute.data"), body)
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

type PeroidRecord struct {
	Date   string
	Time   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volumn float64
	Amount float64
}

type MinuteRecord struct {
	Time   string
	Price  float64
	Volumn float64
	Amount float64
}

func (w *Huobi) analyzePeroidLine(filename string, content string) bool {
	//logger.Infoln(content)
	//logger.Infoln(filename)
	PeroidRecords := parsePeroidCSV(filename)

	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range PeroidRecords {
		Time = append(Time, v.Date+" "+v.Time)
		Price = append(Price, v.Close)
		Volumn = append(Volumn, v.Volumn)
		//Price = append(Price, (v.Close+v.Open+v.High+v.Low)/4.0)
		//Price = append(Price, v.Low)
	}
	w.Time = Time
	w.Price = Price
	w.Volumn = Volumn
	strategyName := Option["strategy"]
	strategy.Perform(strategyName, *w, Time, Price, Volumn)

	return true
}

func (w *Huobi) analyzeMinuteLine(filename string, content string) bool {
	//logger.Infoln(content)
	//logger.Debugln(filename)
	MinuteRecords := parseMinuteCSV(filename)
	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range MinuteRecords {
		Time = append(Time, v.Time)
		Price = append(Price, v.Price)
		Volumn = append(Volumn, v.Volumn)
	}

	w.Time = Time
	w.Price = Price
	w.Volumn = Volumn

	strategyName := Option["strategy"]
	strategy.Perform(strategyName, *w, Time, Price, Volumn)
	return true
}

// reads a whole file into memory and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writes the lines to the given file.
func writeLines(lines []string, path string, skipline int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for i, line := range lines {
		if i < skipline {
			continue
		}

		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

// convert to standard csv file
func data2csv(filename string, skipline int) {
	lines, err := readLines(filename)
	if err != nil {
		logger.Fatalf("readLines: %s", err)
	}
	/*
		for i, line := range lines {
			fmt.Println(i, line)
		}
	*/

	if err := writeLines(lines, filename+".csv", skipline); err != nil {
		logger.Fatalf("writeLines: %s", err)
	}
}

func parseMinuteCSV(filename string) (MinuteRecords []MinuteRecord) {

	// convert to standard csv file
	data2csv(filename, 3)

	file, err := os.Open(filename + ".csv")
	if err != nil {
		fmt.Println("ParseMinuteCSV Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	/*
		record, err := reader.ReadAll()
		fmt.Println(record)
		return
	*/

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(record) == 0 {
			fmt.Println("Error:", "record is zero, maybe it is not a cvs format!!!")
			return
		}

		var minRecord MinuteRecord
		minRecord.Time = record[0]
		minRecord.Price, err = strconv.ParseFloat(record[1], 64)
		if err != nil {
			logger.Fatalln(record)
			logger.Fatal("ParseMinuteCSV item price is not number")
		}
		minRecord.Volumn, err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			logger.Fatal("ParseMinuteCSV item Volumn is not number")
		}
		minRecord.Amount, err = strconv.ParseFloat(record[3], 64)
		if err != nil {
			logger.Fatal("ParseMinuteCSV item Amount is not number")
		}

		MinuteRecords = append(MinuteRecords, minRecord)
	}

	return
}

func parsePeroidCSV(filename string) (PeroidRecords []PeroidRecord) {
	// convert to standard csv file
	data2csv(filename, 2)

	file, err := os.Open(filename + ".csv")
	if err != nil {
		fmt.Println("ParsePeroidCSV Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	/*
		record, err := reader.ReadAll()
		fmt.Println(record)
		return
	*/

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(record) < 8 {
			fmt.Println("Error:", "record is zero, maybe it is not a cvs format!!!", len(record))
			return
		}

		var peroidRecord PeroidRecord
		peroidRecord.Date = record[0]
		peroidRecord.Time = record[1]
		peroidRecord.Open, err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		peroidRecord.High, err = strconv.ParseFloat(record[3], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		peroidRecord.Low, err = strconv.ParseFloat(record[4], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		peroidRecord.Close, err = strconv.ParseFloat(record[5], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}

		peroidRecord.Volumn, err = strconv.ParseFloat(record[6], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item Volumn is not number")
		}
		peroidRecord.Amount, err = strconv.ParseFloat(record[7], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item Amount is not number")
		}

		PeroidRecords = append(PeroidRecords, peroidRecord)
	}
	return
}
