// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

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
var Licence map[string]string
var Remind map[string]string
var TradeOption map[string]string
var SecretOption map[string]string

func init() {
	LoadConfig()
	LoadOption()
	LoadLicence()
	LoadRemind()
	LoadSecretOption()
	//fmt.Println(Config)
	//fmt.Println(Option)
	//fmt.Println(Licence)
	//fmt.Println(Remind)
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

func LoadLicence() error {
	_Licence, err := load_config("/conf/licence.dat")
	if err != nil {
		return err
	}
	Licence = make(map[string]string)
	Licence = _Licence

	return nil
}

func SaveLicence() error {
	return save_config("/conf/licence.dat", Licence)
}

func LoadRemind() (err error) {
	_Remind, err := load_config("/conf/remind.json")
	if err != nil {
		return (err)
	}
	Remind = make(map[string]string)
	Remind = _Remind
	return nil
}

func SaveRemind() error {
	return save_config("/conf/remind.json", Remind)
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
