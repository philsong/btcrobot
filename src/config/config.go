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
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"path"
	"process"
)

// 项目根目录
var ROOT string

var Config map[string]string
var Option map[string]string
var TradeOption map[string]string
var SecretOption map[string]string

func init() {
	LoadAll()
	//fmt.Println(Config)
	//fmt.Println(Option)
	//fmt.Println(Licence)
	//fmt.Println(Remind)
}

func LoadAll() {
	LoadConfig()
	LoadOption()
	LoadSecretOption()
}

func load_config(file string) (config map[string]string, err error) {
	binDir, err := process.ExecutableDir()
	if err != nil {
		return nil, (err)
	}
	ROOT = path.Dir(binDir)

	// Load 全局配置文件
	configFile := ROOT + file
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, (err)
	}
	config = make(map[string]string)
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, (err)
	}

	return config, nil
}

func save_config(file string, config map[string]string) (err error) {
	binDir, err := process.ExecutableDir()
	if err != nil {
		return (err)
	}
	ROOT = path.Dir(binDir)

	var content []byte

	//
	content, err = json.Marshal(&config)
	if err != nil {
		return (err)
	}

	configFile := ROOT + file
	err = ioutil.WriteFile(configFile, content, 666)
	if err != nil {
		return (err)
	}

	return nil
}

func LoadConfig() error {
	_Config, err := load_config("/conf/config.json")
	if err != nil {
		return (err)
	}
	Config = make(map[string]string)
	Config = _Config
	return nil
}

func LoadOption() error {
	_Option, err := load_config("/conf/option.json")
	if err != nil {
		return (err)
	}
	Option = make(map[string]string)
	Option = _Option
	return nil
}

func SaveOption() error {
	return save_config("/conf/option.json", Option)
}

func LoadTrade() (err error) {
	_TradeOption, err := load_config("/conf/trade.json")
	if err != nil {
		return (err)
	}
	TradeOption = make(map[string]string)
	TradeOption = _TradeOption
	return nil
}

func SaveTrade() error {
	return save_config("/conf/trade.json", TradeOption)
}

func LoadSecretOption() (err error) {
	_SecretOption, err := load_config("/conf/secret.json")
	if err != nil {
		return (err)
	}
	SecretOption = make(map[string]string)
	SecretOption = _SecretOption
	return nil
}

func SaveSecretOption() error {
	return save_config("/conf/secret.json", SecretOption)
}
