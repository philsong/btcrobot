BTC/LTC Robot
===========
BTC/LTC自动化算法交易引擎

![btcrobot logo](https://raw.githubusercontent.com/philsong/btcrobot/master/static/img/hacking-bitcoin-with-go.png)

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


支持多个交易平台-火币、OKCoin，自动买卖，机器人EMA算法，MACD算法，MIX算法

程序达成如下目标：上涨时获取中间段的利润；下跌时逃顶，跑赢大盘，比大盘跌的少；横盘时可能高买低卖，需要自行配置参数调整。。
还有一种策略叫买入并持有，buy and hold，亦是一种选择。


举个例子，假设一个IT男，上班时根本没时间一直看盘。
往往去厕所的时间大盘爬上去了，下班的路上，大盘暴跌了，所以如何在利润锁定的范围内，降低这些突发性的风险？
而且，人是很贪婪的，涨的很高，都不想卖，看着下跌，就等死套牢,这个机器人就是辅助你的眼睛，控制你的贪心，至于高买低卖，那个波动，因为这是为了获取拉升的一个机会风险,是否在你可忍受范围内，需要自我把握，

配置可能是伴随市场动态调整的，所以也没什么赚钱的标准配置，需要自己体会。

这个机器人有共振效应，假设很多人用同一个参数配置，很容易产生共振，被人利用，
本机器人完全是开源的，如果你会点golang程序，可以自己写点策略，实现自己的操盘目标
每个人的目标都不太一样
有的人为了赚法币，有的为了赚比特币；）
那算法可能不太一样


若你有自己的方法，大把时间， 不一定非要用机器人
有时候你自己操盘利润可能更多，不然操盘手都失业了。。

使用指南如下:
https://github.com/philsong/btcrobot/wiki/%E6%9C%BA%E5%99%A8%E4%BA%BA%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97 

BTC捐助地址：1NDnnWCUu926z4wxA3sNBGYWNQD3mKyes8

LTC捐助地址：LXwX5XeZeVfXM2b4GRs6HM1mNn4K9En3F4

捐助清单如下：
https://github.com/philsong/btcrobot/wiki/%E6%84%9F%E8%B0%A2%E6%8D%90%E5%8A%A9%E5%BC%80%E5%8F%91%E7%9A%84%E4%BA%BA


安装说明（限win系统）如下：
https://github.com/philsong/btcrobot/wiki/%E5%AE%89%E8%A3%85%E8%AF%B4%E6%98%8E%EF%BC%88%E9%99%90win%E7%B3%BB%E7%BB%9F%EF%BC%89

# 本地搭建 #

0、

0.0安装golang开发运行环境，选择适合自己电脑操作系统的安装包
  
  http://code.google.com/p/go/downloads/list

0.1安装git环境

  http://code.google.com/p/msysgit/downloads/list

1、下载 btcrobot 代码
	
	git clone https://github.com/philsong/btcrobot

2、下载安装依赖库

	cd btcrobot/thirdparty
	// windows/DOS下执行
	getpkg.bat
	// linux/mac或者Git Bash 下执行
	sh getpkg

3、编译 btcrobot

先编译

	// 接着上一步
	cd ../
	// windows/DOS 下执行
	install.bat
	// linux/mac或者Git Bash 下执行
	sh install
	
这样便编译好了 btcrobot

4、配置秘钥API文件

修改btcrobot/conf目录下的secret.sample文件名为secret.json

向火币或者OKcoin申请交易API，并填入secret.json中

5、运行 btcrobot。

	// windows 下执行
	start.bat
	// linux/mac 下执行
	sh start

一切顺利的话，btcrobot应该就启动了。

6、浏览器中查看

在浏览器中输入：http://127.0.0.1:9090

用默认用户名admin, 密码是123456 登录

应该就能开始自动化交易之旅了。

此时，你可以通过WEB界面配置各种参数，参数重启生效。


注：在第5步运行前可以根据需要修改 conf目录里的 配置，亦可在第6步配置。

