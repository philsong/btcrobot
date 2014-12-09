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
	. "common"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// 项目根目录
var ROOT string
var DebugEnv bool
var Config map[string]string
var Option map[string]string
var TradeOption map[string]string
var SecretOption map[string]string

var SimAccount map[string]string

func init() {
	binDir, err := ExecutableDir()
	if err != nil {
		return
	}
	ROOT = path.Dir(binDir)
	return
}

func LoadSimulate() error {
	var filename string

	if !GetBacktest() {
		filename = "/conf/simulate.json"
	} else {
		filename = "/test/simulate.json"
	}

	_Simulate, err := load_config(filename)
	if err != nil {
		return err
	}
	SimAccount = _Simulate
	return nil
}

func SaveSimulate() error {
	var filename string

	if !GetBacktest() {
		filename = "/conf/simulate.json"
	} else {
		filename = "/test/simulate.json"
	}

	return save_config(filename, SimAccount)
}

func init() {
	LoadAll()
	var pDebugEnv *bool
	pDebugEnv = flag.Bool("d", false, "enable debug for dev")
	flag.Parse()
	DebugEnv = *pDebugEnv
}

func LoadAll() {
	LoadConfig()
	LoadOption()
	LoadSecretOption()
}

func get_config_path(file string) (filepath string) {
	return ROOT + file
}

func load_config(file string) (config map[string]string, err error) {
	// Load 全局配置文件
	configFile := get_config_path(file)

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	config = make(map[string]string)
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func SaveContent(file string, object interface{}) error {
	content, err := json.Marshal(&object)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ROOT+file, content, 666)
}

func LoadFile(file string) (content []byte, err error) {
	file = get_config_path(file)
	content, err = ioutil.ReadFile(file)
	return
}

func LoadContent(file string, object interface{}) error {
	content, err := LoadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, &object)
}

func save_config(file string, config map[string]string) (err error) {
	return SaveContent(file, config)
}

func LoadConfig() error {
	_Config, err := load_config("/conf/config.json")
	if err != nil {
		return err
	}
	Config = make(map[string]string)
	Config = _Config
	return nil
}

func LoadOption() error {
	_Option, err := load_config("/conf/option.json")
	if err != nil {
		return err
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
		return err
	}
	TradeOption = make(map[string]string)
	TradeOption = _TradeOption
	return nil
}

func SaveTrade() error {
	return save_config("/conf/trade.json", TradeOption)
}

// filesExists returns whether or not the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func LoadSecretOption() (err error) {
	secretFile := get_config_path("/conf/secret.json")
	secretSampleFile := get_config_path("/conf/secret.sample")

	if !fileExists(secretFile) {
		if err := os.Rename(secretSampleFile, secretFile); err != nil {
			fmt.Println(err)
			err := fmt.Errorf("unable to create secret.json "+
				"root: %v", err)
			fmt.Println(err)
			return err
		}
	}

	_SecretOption, err := load_config("/conf/secret.json")
	if err != nil {
		return err
	}
	SecretOption = make(map[string]string)
	SecretOption = _SecretOption
	return nil
}

func SaveSecretOption() error {
	return save_config("/conf/secret.json", SecretOption)
}

// 获得可执行程序所在目录
func ExecutableDir() (string, error) {
	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(pathAbs), nil
}
