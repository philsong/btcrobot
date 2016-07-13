BTC/LTC Robot

============
## Golang BTC/LTC trading bot engine


![btcrobot logo](
https://raw.githubusercontent.com/philsong/btcrobot/master/static/images/demo/hacking-bitcoin-with-go.png)

  btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang,
  it features multiple trading methods using technical analysis.

  ## Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

## How You Can Help

### Donations

Donations help me to keep working on btcrobot and keep it free and open source, without having to worry about income. Any amount is really helpful! Thank you so much.

The btcrobot donation Bitcoin address is **1NDnnWCUu926z4wxA3sNBGYWNQD3mKyes8**


Again, thank you. :heart:

### Contributing

Help us build! We're in beta right now [and seeking help to find bugs]. If you are interested in contributing, jump in! Anyone is welcome to send pull requests. Issue reports are good too, but pull requests are much better. Here's how you do it:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Write the code, **and tests to confirm it works**
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request


声明：

本软件风险巨大！！！使用前务必小量测试N天！！！目前还在快速完善发展中，欢迎各路程序员贡献代码！

Weibo:http://weibo.com/bocaicfa


支持多个交易平台-火币、OKCoin，自动买卖，机器人EMA，MACD，KDJ, 对冲，高低突破，及各种MIX算法

程序达成如下目标：上涨时获取中间段的利润；下跌时逃顶，跑赢大盘，比大盘跌的少；横盘时可能高买低卖，需要提高周期或优化参数。

Update: 

2014/4/23: EMA、MACD等各种指标和其他策略，在市场深度不够时，已然失效。。。业界其他程序员即使基于robot写了策略，也极少开源，原因就是一旦公开策略，使用人数变多，该策略就失效。即日起本人也不再公开发布新的策略，改变开发方向为专注于该开源框架的完善，然后提供回测和策略接口甚至提供策略交易市场，让用户方便自己实现自己的策略并可回测，及交流。有任何建议及合作都可互相探讨。

2015/5/25: 项目已内部转化为公司买卖币及搬砖模块。

2016/7/13: 该内部项目已经 开源：

搬砖做市套利： https://github.com/philsong/bitcoin-arbitrage

代理商模型： https://github.com/philsong/bitcoin-broker

## 使用指南:
https://github.com/philsong/btcrobot/wiki/%E6%9C%BA%E5%99%A8%E4%BA%BA%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97 


## 捐助开发清单：
https://github.com/philsong/btcrobot/wiki/%E6%84%9F%E8%B0%A2%E6%8D%90%E5%8A%A9%E5%BC%80%E5%8F%91%E7%9A%84%E4%BA%BA


BTC捐助地址：1NDnnWCUu926z4wxA3sNBGYWNQD3mKyes8


# 本地搭建 #

0、

0.0安装golang开发运行环境，选择适合自己电脑操作系统的安装包
  
  http://code.google.com/p/go/downloads/list

0.1安装git环境

  http://msysgit.github.io/

  注意:windows下安装时一定要选择，把git路径加入到系统PATH中,否则后续安装无法找到git命令。

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

4、运行 btcrobot。

	// windows 下执行
	start.bat
	// linux/mac或者Git Bash 下执行
	sh start

一切顺利的话，btcrobot应该就启动了。

5、浏览器中登陆

在浏览器中打开：http://127.0.0.1:9090

用默认用户名admin, 密码123456 登录

6、配置秘钥API文件

向火币或者OKcoin申请交易API，并填入”安全配置“菜单中。


应该就能开始自动化交易之旅了。

此时，你可以通过WEB界面配置各种参数，参数重启生效。


注：在第5步运行前可以根据需要修改 conf目录里的 配置，亦可在第6步配置。

## 安装说明（限win系统）：
https://github.com/philsong/btcrobot/wiki/%E5%AE%89%E8%A3%85%E8%AF%B4%E6%98%8E%EF%BC%88%E9%99%90win%E7%B3%BB%E7%BB%9F%EF%BC%89
