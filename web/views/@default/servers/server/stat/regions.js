Tea.context(function () {
	this.$delay(function () {
		let that = this

		// 地区
		let countryAxis = teaweb.countAxis(this.countryStats, function (v) {
			return v.count
		})
		this.reloadChart("country-chart", "地区", this.countryStats, function (v) {
			return v.country.name
		}, function (args) {
			return that.countryStats[args.dataIndex].country.name + ": " + teaweb.formatNumber(that.countryStats[args.dataIndex].count)
		}, countryAxis)

		// 省份
		let provinceAxis = teaweb.countAxis(this.provinceStats, function (v) {
			return v.count
		})
		this.reloadChart("province-chart", "省市", this.provinceStats, function (v) {
			return v.province.name
		}, function (args) {
			return that.provinceStats[args.dataIndex].country.name + ": " + that.provinceStats[args.dataIndex].province.name + " " + teaweb.formatNumber(that.provinceStats[args.dataIndex].count)
		}, provinceAxis)

		// 城市
		let cityAxis = teaweb.countAxis(this.cityStats, function (v) {
			return v.count
		})
		this.reloadChart("city-chart", "城市", this.cityStats, function (v) {
			return v.city.name
		}, function (args) {
			return that.cityStats[args.dataIndex].country.name + ": " + that.cityStats[args.dataIndex].province.name + " " + that.cityStats[args.dataIndex].city.name + " " + teaweb.formatNumber(that.cityStats[args.dataIndex].count)
		}, cityAxis)
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
