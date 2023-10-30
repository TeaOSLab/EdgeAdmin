Tea.context(function () {
	this.$delay(function () {
		let that = this

		let systemAxis = teaweb.countAxis(this.systemStats, function (v) {
			return v.count
		})
		this.reloadChart("system-chart", "操作系统", this.systemStats, function (v) {
			return v.system.name
		}, function (args) {
			return that.systemStats[args.dataIndex].system.name + ": " + teaweb.formatNumber(that.systemStats[args.dataIndex].count)
		}, systemAxis)

		let browserAxis = teaweb.countAxis(this.browserStats, function (v) {
			return v.count
		})
		this.reloadChart("browser-chart", "浏览器", this.browserStats, function (v) {
			return v.browser.name
		}, function (args) {
			return that.browserStats[args.dataIndex].browser.name + ": " + teaweb.formatNumber(that.browserStats[args.dataIndex].count)
		}, browserAxis)
	})

	this.reloadChart = function (chartId, name, stats, xFunc, tooltipFunc, axis) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = teaweb.initChart(chartBox)
		let option = {
			xAxis: {
				data: stats.map(xFunc),
				axisLabel: {
					interval: 0
				}
			},
			yAxis: {
				axisLabel: {
					formatter: function (value) {
						return value + axis.unit
					}
				}
			},
			tooltip: {
				show: true,
				trigger: "item",
				formatter: tooltipFunc
			},
			grid: {
				left: 40,
				top: 10,
				right: 20,
				bottom: 20
			},
			series: [
				{
					name: name,
					type: "bar",
					data: stats.map(function (v) {
						return v.count / axis.divider;
					}),
					itemStyle: {
						color: teaweb.DefaultChartColor
					},
					barWidth: "10em"
				}
			],
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}
})
