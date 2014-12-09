package main

import (
	. "common"
	"compress/gzip"
	. "config"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	. "util"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

func parseEndpoint(u *url.URL) string {
	path := u.Path
	if l := len(path); l > 0 && path[len(path)-1] == '/' {
		path = path[:l-1]
	}
	lastPath := strings.LastIndex(path, "/")
	endpoint := ""
	if lastPath >= 0 {
		path := path[lastPath:]
		if len(path) > 0 {
			endpoint = path
		}
	}
	return endpoint
}

func Dial(origin string) (*websocket.Conn, error) {
	u, err := url.Parse(origin)
	if err != nil {
		return nil, err
	}
	//endpoint := parseEndpoint(u)
	u.Path = fmt.Sprintf("/socket.io/%d/", ProtocolVersion)

	url_ := u.String()
	r, err := http.Get(url_)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return nil, errors.New("invalid status: " + r.Status)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(string(body), ":", 4)
	if len(parts) != 4 {
		return nil, errors.New("invalid handshake: " + string(body))
	}
	if !strings.Contains(parts[3], "websocket") {
		return nil, errors.New("server does not support websockets")
	}
	sessionId := parts[0]
	u.Scheme = "ws" + u.Scheme[4:]
	u.Path = fmt.Sprintf("%swebsocket/%s", u.Path, sessionId)
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	//timeout, err := strconv.ParseInt(parts[1], 10, 64)
	//if err != nil {
	//	return nil, err
	//}

	return ws, err
}

type Client struct {
	ws       *websocket.Conn
	endpoint string
	timeout  int64
}

func NumberToFloat(num json.Number) (value float64) {
	var err error
	value, err = num.Float64()
	if err != nil {
		value = 0.0
	}
	return
}

func GetKLine(from, to time.Time, period int64, symbolId string) (records []Record, err error) {
	var periodHuobi string

	switch period {
	case 1:
		periodHuobi = "1min"
	case 5:
		periodHuobi = "5min"
	case 15:
		periodHuobi = "15min"
	case 30:
		periodHuobi = "30min"
	case 60:
		periodHuobi = "60min"
	case 100:
		periodHuobi = "1day"
	}

	filename := from.Format("2006-01-02") + "_" + to.Format("2006-01-02")
	filename = "\\test\\" + filename + "_" + periodHuobi + "_" + symbolId + ".kline"

	if Exist(ROOT + filename) {
		kline, _ := LoadFile(filename)
		klineJson, _ := simplejson.NewJson(kline)
		klineArray, _ := klineJson.Array()
		fmt.Println(len(klineArray))
		for i := 0; i < len(klineArray); i++ {
			var record Record

			line := klineArray[i].(map[string]interface{})
			record.Time = int64(NumberToFloat(line["Time"].(json.Number)))
			t := time.Unix(record.Time, 0)
			record.TimeStr = t.Format("2006-01-02 15:04:05")
			record.Open = NumberToFloat(line["Open"].(json.Number))
			record.High = NumberToFloat(line["High"].(json.Number))
			record.Low = NumberToFloat(line["Low"].(json.Number))
			record.Close = NumberToFloat(line["Close"].(json.Number))
			record.Volumn = NumberToFloat(line["Volumn"].(json.Number))

			records = append(records, record)
		}
		return
	}

	var shift int64
	switch periodHuobi {
	case "1min":
		shift = 60 * ShiftNumber
	case "5min":
		shift = 300 * ShiftNumber
	case "15min":
		shift = 900 * ShiftNumber
	case "30min":
		shift = 1800 * ShiftNumber
	case "60min":
		shift = 3600 * ShiftNumber
	case "1day":
		shift = 3600 * 24 * ShiftNumber
	}

	reqmap := make(map[string]interface{})
	reqmap["symbolId"] = symbolId
	reqmap["version"] = 1
	reqmap["msgType"] = "reqKLine"
	reqmap["period"] = periodHuobi
	reqmap["from"] = from.Unix() - shift
	reqmap["to"] = to.Unix()

ReDial:
	ws, err := Dial("http://hq.huobi.com:80")
	if err != nil {
		fmt.Printf("Dial: %v", err)
		return
	}
	defer ws.Close()

L:
	for {
		buf, _ := json.Marshal(reqmap)
		reqstr := "5:::{\"name\":\"request\",\"args\":[" + string(buf) + "]}\n"

		ws.WriteMessage(websocket.TextMessage, []byte(reqstr))

		for {
			_, msgBuf, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
				ws.Close()
				goto ReDial
			}
			index := strings.Index(string(msgBuf), "5:::")
			if index != -1 {
				msg := string(msgBuf)
				msg = msg[strings.Index(msg, "payload")+9 : len(msg)-3]
				jsonKLine, err := simplejson.NewJson([]byte(msg))
				if err != nil {
					fmt.Println(err)
					break L
				}

				jsonTime, _ := jsonKLine.CheckGet("time")
				jsonOpen, _ := jsonKLine.CheckGet("priceOpen")
				jsonHigh, _ := jsonKLine.CheckGet("priceHigh")
				jsonLow, _ := jsonKLine.CheckGet("priceLow")
				jsonLast, _ := jsonKLine.CheckGet("priceLast")
				jsonAmount, _ := jsonKLine.CheckGet("amount")

				arrayTime, _ := jsonTime.Array()
				arrayOpen, _ := jsonOpen.Array()
				arrayHigh, _ := jsonHigh.Array()
				arrayLow, _ := jsonLow.Array()
				arrayLast, _ := jsonLast.Array()
				arrayAmount, _ := jsonAmount.Array()

				for i := 0; i < len(arrayTime); i++ {
					var record Record

					record.Time = int64(NumberToFloat(arrayTime[i].(json.Number)))
					t := time.Unix(record.Time, 0)
					record.TimeStr = t.Format("2006-01-02 15:04:05")
					record.Open = NumberToFloat(arrayOpen[i].(json.Number))
					record.High = NumberToFloat(arrayHigh[i].(json.Number))
					record.Low = NumberToFloat(arrayLow[i].(json.Number))
					record.Close = NumberToFloat(arrayLast[i].(json.Number))
					record.Volumn = NumberToFloat(arrayAmount[i].(json.Number))
					records = append(records, record)
					if record.Time >= to.Unix() {
						break L
					}
				}
				fmt.Println("读取k线完成：", records[len(records)-1].TimeStr)
				fromUix := records[len(records)-1].Time + (records[1].Time - records[0].Time)
				if fromUix >= to.Unix() {
					break L
				}
				reqmap["from"] = fromUix
				break
			}
		}
	}

	SaveContent(filename, records)

	return
}

// 读取k线数据，非正式API
func OkcoinKLine(period int64, symbolId string) (records []Record, err error) {
	periodOkcoin := IntegerToString(period * 60)
	reqest, err := http.NewRequest("GET", "https://www.okcoin.cn/kline/period.do?step="+periodOkcoin+"&symbol=okcoin"+symbolId, nil)
	reqest.Header.Set("Accept", "*/*")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.103 Safari/537.36")
	reqest.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := http.DefaultClient.Do(reqest)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("Http StatusCode:" + strconv.Itoa(resp.StatusCode))
		resp.Body.Close()
		return
	}
	var body []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(resp.Body)
		body, _ = ioutil.ReadAll(reader)
	default:
		body, _ = ioutil.ReadAll(resp.Body)
	}
	resp.Body.Close()
	var KLine [][8]float64 // 时间戳，0，0，开盘，收盘，最高，最低，成交量
	dec := json.NewDecoder(strings.NewReader(string(body)))
	if err = dec.Decode(&KLine); err != nil {
		return
	}

	for i := 0; i < len(KLine); i++ {
		var record Record

		record.Time = int64(KLine[i][0])
		t := time.Unix(record.Time, 0)
		record.TimeStr = t.Format("2006-01-02 15:04:05")
		record.Open = KLine[i][3]
		record.High = KLine[i][5]
		record.Low = KLine[i][6]
		record.Close = KLine[i][4]
		record.Volumn = KLine[i][7]
		records = append(records, record)
	}

	filename := time.Now().Format("2006-01-02") + "_okcoin_"
	filename = "\\test\\" + filename + "_" + IntegerToString(period) + "_" + symbolId + ".kline"
	SaveContent(filename, records)

	return
}
