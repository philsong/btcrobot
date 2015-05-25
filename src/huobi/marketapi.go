/*
  btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang,
  it features multiple trading methods using technical analysis.

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
	. "common"
	. "config"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"os"
	"strconv"
	"util"
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

func (w *Huobi) AnalyzeKLinePeroid(symbol string, peroid int) (ret bool, records []Record) {
	ret = false
	var huobisymbol string
	if symbol == "btc_cny" {
		huobisymbol = "huobibtccny"
	} else {
		huobisymbol = "huobiltccny"
		logger.Fatal("huobi does not support LTC by now, wait for huobi provide it.", huobisymbol)
		return
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["hb_kline_url"], peroid), nil)
	if err != nil {
		logger.Fatal(err)
		return
	}

	req.Header.Set("Referer", Config["hb_base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	c := util.NewTimeoutClient()

	logger.Tracef("HTTP req begin AnalyzeKLinePeroid")
	resp, err := c.Do(req)
	logger.Tracef("HTTP req end AnalyzeKLinePeroid")
	if err != nil {
		logger.Traceln(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Tracef("HTTP returned status %v", resp)
		return
	}
	var body string

	contentEncoding := resp.Header.Get("Content-Encoding")
	logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
	logger.Traceln(resp.Header.Get("Content-Type"))
	switch contentEncoding {
	case "gzip":
		body = util.DumpGZIP(resp.Body)
	default:
		bodyByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Errorln("read the http stream failed")
			return
		} else {
			body = string(bodyByte)
		}
	}

	ioutil.WriteFile(fmt.Sprintf("cache/hbKLine_%03d.data", peroid), []byte(body), 0644)

	return analyzePeroidLine(fmt.Sprintf("cache/hbKLine_%03d.data", peroid))
}

func analyzePeroidLine(filename string) (ret bool, records []Record) {
	// convert to standard csv file
	data2csv(filename, 2)

	ret = false
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
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Fatal("Error:", err)
			return
		}

		if len(line) < 8 {
			logger.Fatal("Error:", "record is zero, maybe it is not a cvs format!!!", len(line))
			return
		}

		var record Record
		record.TimeStr = line[0] + " " + line[1]
		record.Open, err = strconv.ParseFloat(line[2], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		record.High, err = strconv.ParseFloat(line[3], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		record.Low, err = strconv.ParseFloat(line[4], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		record.Close, err = strconv.ParseFloat(line[5], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}

		record.Volumn, err = strconv.ParseFloat(line[6], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item Volumn is not number")
		}
		_, err = strconv.ParseFloat(line[7], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item Amount is not number")
		}

		records = append(records, record)
	}

	ret = true
	return
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
