Tea.context(function () {
	this.$delay(function () {
		let that = this

		let axis = teaweb.countAxis(this.providerStats, function (v) {
			return v.count
		})
		this.reloadChart("provider-chart", "运营商", this.providerStats, function (v) {
			return v.provider.name
		}, function (args) {
			return that.providerStats[args.dataIndex].provider.name + ": " + teaweb.formatNumber(that.providerStats[args.dataIndex].count)
		}, axis)
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
