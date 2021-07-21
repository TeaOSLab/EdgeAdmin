Tea.context(function () {
	this.trafficTab = "hourly"

	this.$delay(function () {
		this.reloadHourlyTrafficChart()
	})

	this.selectTrafficTab = function (tab) {
		this.trafficTab = tab
		if (tab == "hourly") {
			this.$delay(function () {
				this.reloadHourlyTrafficChart()
			})
		} else if (tab == "daily") {
			this.$delay(function () {
				this.reloadDailyTrafficChart()
			})
		}
	}

	this.reloadHourlyTrafficChart = function () {
		let axis = teaweb.bytesAxis(this.hourlyTrafficStats, function (v) {
			return v.bytes
		})
		let chartBox = document.getElementById("hourly-traffic-chart-box")
		let chart = teaweb.initChart(chartBox)
		let that = this
		let option = {
			xAxis: {
				data: this.hourlyTrafficStats.map(function (v) {
					return v.hour;
				})
			},
			yAxis: {
				axisLabel: {
					formatter: function (v) {
						return v + axis.unit
					}
				}
			},
			tooltip: {
				show: true,
				trigger: "item",
				formatter: function (args) {
					let index = args.dataIndex
					return that.hourlyTrafficStats[index].hour + "时：" + teaweb.formatBytes(that.hourlyTrafficStats[index].bytes)
				}
			},
			grid: {
				left: 40,
				top: 10,
				right: 20
			},
			series: [
				{
					name: "流量",
					type: "line",
					data: this.hourlyTrafficStats.map(function (v) {
						return v.bytes / axis.divider;
					}),
					itemStyle: {
						color: "#9DD3E8"
					},
					lineStyle: {
						color: "#9DD3E8"
					},
					areaStyle: {
						color: "#9DD3E8"
					}
				}
			],
			animation: false
		}
		chart.setOption(option)
		chart.resize()
	}

	this.reloadDailyTrafficChart = function () {
		let axis = teaweb.bytesAxis(this.dailyTrafficStats, function (v) {
			return v.bytes
		})
		let chartBox = document.getElementById("daily-traffic-chart-box")
		let chart = teaweb.initChart(chartBox)

		let that = this
		let option = {
			xAxis: {
				data: this.dailyTrafficStats.map(function (v) {
					return v.day;
				})
			},
			yAxis: {
				axisLabel: {
					formatter: function (v) {
						return v + axis.unit
					}
				}
			},
			tooltip: {
				show: true,
				trigger: "item",
				formatter: function (args) {
					let index = args.dataIndex
					return that.dailyTrafficStats[index].day + "：" + teaweb.formatBytes(that.dailyTrafficStats[index].bytes)
				}
			},
			grid: {
				left: 40,
				top: 10,
				right: 20
			},
			series: [
				{
					name: "流量",
					type: "line",
					data: this.dailyTrafficStats.map(function (v) {
						return v.bytes / axis.divider;
					}),
					itemStyle: {
						color: "#9DD3E8"
					},
					lineStyle: {
						color: "#9DD3E8"
					},
					areaStyle: {
						color: "#9DD3E8"
					}
				}
			],
			animation: false
		}
		chart.setOption(option)
		chart.resize()
	}

	/**
	 * 升级提醒
	 */
	this.closeMessage = function (e) {
		let target = e.target
		while (true) {
			target = target.parentNode
			if (target.tagName.toUpperCase() == "DIV") {
				target.style.cssText = "display: none"
				break
			}
		}
	}
})
