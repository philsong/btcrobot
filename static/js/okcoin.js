// Copyright 2014 The btcrobot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/btcrobot

$(function() {
	myTimer();
});

var myVar = setInterval(function() {
	myTimer()
}, 10000);

function myStopFunction() {
	clearInterval(myVar);
}

function myStartFunction() {
	clearInterval(myVar);
	myVar = setInterval(function() {
		myTimer()
	}, 10000);
}

function myTimer() {
	$.ajaxSetup({
		// Disable caching of AJAX responses
		cache: false
	});

	$.get('https://www.okcoin.com/kline/period.do?step=60&symbol=okcoinbtccny', function(data) {
		// split the data set into ohlc and volume
		var datas = [];

		// Split the lines
		var lines = JSON.parse(data);
		lines = lines.slice(-290) //okcoin total:1440, but huobi only 4hour50minutes=290minutes

		// Iterate over the lines and add categories or series
		$.each(lines, function(lineNo, line) {
			//console.log(line[0])
			var d = new Date();
			d.setTime(1000*line[0])
			
			datas.push([
				d.getTime(),
				parseFloat(line[6])
			]);	
		});

		// Create the chart
		$('#container').highcharts('StockChart', {


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
	}, "text");
}