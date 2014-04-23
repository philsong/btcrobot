package webui

import (
	. "config"
	"fmt"
	"github.com/codegangsta/martini-contrib/auth"
	"github.com/go-martini/martini"
	"logger"
	"net/http"

	"config"
	"encoding/json"
)

func webui() {
	m := martini.Classic()
	m.Get("/secret", func() string {
		// show something
		Option_json, err := json.Marshal(config.SecretOption)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取引擎配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/secret", func() {
		// create something
	})

	m.Get("/engine", func() string {
		// show something
		Option_json, err := json.Marshal(config.Option)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取引擎配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/engine", func() {
		// create something
	})

	m.Get("/trade", func() string {
		// show something
		Option_json, err := json.Marshal(config.TradeOption)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取引擎配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/trade", func() {
		// create something
	})

	m.Use(auth.Basic(config.SecretOption["username"], config.SecretOption["password"]))
	m.Use(martini.Static("./static"))

	logger.Infoln(http.ListenAndServe(Config["host"], m))

	m.Run()

	fmt.Println("[ ---------------------------------------------------------->>> ")
	fmt.Printf("start web server failed, please check if %s is already used.", Config["host"])
	fmt.Println(" <<<----------------------------------------------------------] ")
}

func RunServer() {
	webui()
}
