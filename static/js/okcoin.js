// Copyright 2014 The btcrobot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/btcrobot

var wsUri = "wss://real.okcoin.cn:10440/websocket/okcoinapi";
var output;

function init() {
	output = document.getElementById("output");
	testWebSocket();
}

function testWebSocket() {
	websocket = new WebSocket(wsUri);
	websocket.onopen = function(evt) {
		onOpen(evt)
	};
	websocket.onclose = function(evt) {
		onClose(evt)
	};
	websocket.onmessage = function(evt) {
		onMessage(evt);
		websocket.close();
	};
	websocket.onerror = function(evt) {
		onError(evt)
		
	};
}

function onOpen(evt) {
	//writeToScreen("CONNECTED");
	doSend("{'event':'addChannel','channel':'ok_btccny_kline_1hour'}");
}

function onClose(evt) {
	//writeToScreen("DISCONNECTED");
}

function onMessage(evt) {
	//writeToScreen('<span style="color: blue;">RESPONSE: ' + evt.data + '</span>');

	// split the data set into ohlc and volume
	var datas = [];

	// Split the lines
	//console.log(evt.data);
	var lines = JSON.parse(evt.data);
	//var lines = evt.data.data;
	console.log(lines[0].data);

	//writeToScreen('<span style="color: blue;">RESPONSE: ' + lines[0].data+ '</span>');
	//lines = lines.slice(-290) //okcoin total:1440, but huobi only 4hour50minutes=290minutes

	// Iterate over the lines and add categories or series
	$.each(lines[0].data, function(lineNo, line) {
		//console.log(line[0])
		var d = new Date();
		d.setTime(line[0])

		datas.push([
			d.getTime(),
			parseFloat(line[4])
		]);
	});

	// Create the chart
	$('#indictorcontainer').highcharts('StockChart', {


		title: {
			text: 'OKcoin趋势指标'
		},
		subtitle: {
			text: '时间序列价格分析系统'
		},
		xAxis: {
			type: 'datetime'
		},

		yAxis: [{
			title: {
				text: '价格'
			},
			height: 200,
			plotLines: [{
				value: 0,
				width: 1,
				color: '#808080'
			}]
		}, {
			title: {
				text: 'MACD'
			},
			top: 300,
			height: 100,
			offset: 0,
			lineWidth: 0.1
		}],

		tooltip: {
			crosshairs: true,
			shared: true
		},

		rangeSelector: {
			selected: 1
		},

		legend: {
			enabled: true,
			layout: 'vertical',
			align: 'right',
			verticalAlign: 'middle',
			borderWidth: 0
		},

		plotOptions: {
			series: {
				marker: {
					enabled: false,
				}
			}
		},


		series: [{
			name: '价格',
			type: 'line',
			id: 'primary',
			data: datas
		}, {
			name: 'MACD',
			linkedTo: 'primary',
			yAxis: 1,
			showInLegend: true,
			type: 'trendline',
			algorithm: 'MACD'

		}, {
			name: 'Signal line信号线',
			linkedTo: 'primary',
			yAxis: 1,
			showInLegend: true,
			type: 'trendline',
			algorithm: 'signalLine'

		}, {
			name: 'Histogram柱状图',
			linkedTo: 'primary',
			yAxis: 1,
			showInLegend: true,
			type: 'histogram'

		}, {
			name: '10-EMA',
			linkedTo: 'primary',
			showInLegend: true,
			type: 'trendline',
			algorithm: 'EMA',
			periods: 10
		}, {
			name: '21-EMA',
			linkedTo: 'primary',
			showInLegend: true,
			type: 'trendline',
			algorithm: 'EMA',
			periods: 21
		}, {
			name: 'Linear Trendline趋势线',
			linkedTo: 'primary',
			showInLegend: true,
			enableMouseTracking: false,
			type: 'trendline',
			algorithm: 'linear'
		}]
	});

	//websocket.close();
}

function onError(evt) {
	writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
}

function doSend(message) {
	//writeToScreen("SENT: " + message);
	websocket.send(message);
}

function writeToScreen(message) {
	var pre = document.createElement("p");
	pre.style.wordWrap = "break-word";
	pre.innerHTML = message;
	output.appendChild(pre);
}
window.addEventListener("load", init, false);

