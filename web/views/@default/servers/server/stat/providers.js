Tea.context(function () {
	this.$delay(function () {
		let that = this

		let axis = teaweb.countAxis(this.providerStats, function (v) {
			return v.count
		})
		this.providerStats.forEach(function (v) {
			v.count /= axis.divider
		})
		this.reloadChart("provider-chart", "运营商", this.providerStats, function (v) {
			return v.provider.name
		}, function (args) {
			return that.providerStats[args.dataIndex].provider.name + ": " + teaweb.formatNumber(that.providerStats[args.dataIndex].rawCount)
		}, axis.unit)
		window.addEventListener("resize", function () {
			that.resizeChart("provider-chart")
		})
	})

	this.reloadChart = function (chartId, name, stats, xFunc, tooltipFunc, unit) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = echarts.init(chartBox)
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
						return value + unit
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
						return v.count;
					}),
					itemStyle: {
						color: "#9DD3E8"
					},
					barWidth: "20em"
				}
			],
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}

	this.resizeChart = function (chartId) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = echarts.init(chartBox)
		chart.resize()
	}
})
