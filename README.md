BTC Robot
===========
BTC操盘手自动化交易引擎

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

原理来自Goomboo的分析,EMA crossover：https://bitcointalk.org/index.php?topic=60501.0

BTC捐赠地址：1NDnnWCUu926z4wxA3sNBGYWNQD3mKyes8


# 本地搭建 #

1、下载 btcrobot 代码
	
	git clone https://github.com/philsong/btcrobot

2、下载安装依赖库（如果依赖库下载不下来可以联系我）

	cd btcrobot/websites/code/thirdparty
	// windows 下执行
	getpkg.bat
	// linux/mac 下执行
	sh getpkg

3、编译 btcrobot

先编译

	// 接着上一步
	cd ../btcbot/
	// windows 下执行
	install.bat
	// linux/mac 下执行
	sh install
	
这样便编译好了 btcrobot

4、配置秘钥API文件

修改btcrobot/websites/code/btcbot/conf目录下的secret.sample文件名为secret.json

5、运行 btcrobot。

	// windows 下执行
	start.bat
	// linux/mac 下执行
	sh start

一切顺利的话，btcrobot应该就启动了。

6、浏览器中查看

在浏览器中输入：http://127.0.0.1:9090

应该就能看到了。

此时，你可以通过WEB界面配置各种参数，参数实时生效。


注：在第5步运行前可以根据需要修改 conf目录里的 配置，亦可在第6步配置。

/*

 *BTC操盘手自动化交易引擎

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


 *@feedback http://weibo.com/bocaicfa

 *@Version 0.10
 *@time 2014-01-06 support auto process: login/buy/sell/update/cancel/delegation query/auto trade

 *@Version 0.11
 *@time 2014-01-09 support query details transaction data:buy/sell/trade/topsell/topbuy/new/high/low/last/...

  @Version 0.12
 *@time 2014-01-10 support query my trade info, 5min/1day/1week/... kline data and real trx data
 				   support okcoin.com K-line via highchart

 *@Version 0.13
 *@time 2014-01-13 support EMA line to indict the time to buy and sell simulately, support diff factor
 				   support web interface to display MA/EMA/MACD/Trender line too

 *@Version 0.14
 *@time 2014-01-14 support EMA line to indict the time to buy and sell in huobi.com

 *@Version 0.15
 *@time 2014-01-16 support send alert email when triger buy/sell point

 *@Version 0.16
 *@time 2014-01-27 support huibi official API,optimize MACD enter point

 *@Version 0.17
 *@time 2014-02-08 support the 5mintes momentum theory in fx

 *@Version 0.18
 *@time 2014-02-10 simplify the 5mintes momentum theory, only keep three key points:"enter"/"stop"/"exit" 

 *@Version 0.19
 *@time 2014-03-01 add the web UI to config option
 

 *@go语言(模拟登陆huobi.com平台)+(官方API)实现自动化套利

 */
