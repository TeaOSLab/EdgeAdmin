Tea.context(function () {
	this.$delay(function () {
		let that = this

		this.countDailyBlock = teaweb.formatCount(this.countDailyBlock)
		this.countDailyCaptcha = teaweb.formatCount(this.countDailyCaptcha)
		this.countDailyLog = teaweb.formatCount(this.countDailyLog)
		this.countWeeklyBlock = teaweb.formatCount(this.countWeeklyBlock)
		this.countMonthlyBlock = teaweb.formatCount(this.countMonthlyBlock)

		this.totalDailyStats = this.logDailyStats.map(function (v, k) {
			return {
				day: v.day,
				count: that.logDailyStats[k].count + that.blockDailyStats[k].count + that.captchaDailyStats[k].count
			}
		})
		let dailyAxis = teaweb.countAxis(this.totalDailyStats, function (v) {
			return v.count
		})
		this.reloadLineChart("daily-chart", "规则分组", this.totalDailyStats, function (v) {
			return v.day.substring(4, 6) + "-" + v.day.substring(6)
		}, function (args) {
			return that.logDailyStats[args.dataIndex].day.substring(4, 6) + "-" + that.logDailyStats[args.dataIndex].day.substring(6) + ": 拦截: "
				+ teaweb.formatNumber(that.blockDailyStats[args.dataIndex].count) + ", 验证码: " + teaweb.formatNumber(that.captchaDailyStats[args.dataIndex].count) + ", 记录: " + teaweb.formatNumber(that.logDailyStats[args.dataIndex].count)
		}, dailyAxis)

		let groupAxis = teaweb.countAxis(this.groupStats, function (v) {
			return v.count
		})
		let total = this.groupStats.$sum(function (k, v) {
			return v.count
		})
		this.reloadBarChart("group-chart", "规则分组", this.groupStats, function (v) {
			return v.group.name
		}, function (args) {
			let percent = ""
			if (total > 0) {
				percent = ", 占比: " + (Math.round(that.groupStats[args.dataIndex].count * 100 * 100 / total) / 100) + "%"
			}
			return that.groupStats[args.dataIndex].group.name + ": " + teaweb.formatNumber(that.groupStats[args.dataIndex].count) + percent
		}, groupAxis)
	})

	this.reloadLineChart = function (chartId, name, stats, xFunc, tooltipFunc, axis) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let that = this
		let chart = teaweb.initChart(chartBox)
		let option = {
			xAxis: {
				data: stats.map(xFunc)
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
					type: "line",
					data: this.totalDailyStats.map(function (v, index) {
						return that.totalDailyStats[index].count / axis.divider;
					}),
					areaStyle: {},
					itemStyle: {
						color: teaweb.DefaultChartColor
					},
					smooth: true
				},
				{
					name: name,
					type: "line",
					data: this.logDailyStats.map(function (v) {
						return v.count / axis.divider;
					}),
					itemStyle: {
						color: "#879BD7"
					},
					smooth: true
				},
				{
					name: name,
					type: "line",
					data: this.blockDailyStats.map(function (v) {
						return v.count / axis.divider;
					}),
					itemStyle: {
						color: "#F39494"
					},
					smooth: true
				},
				{
					name: name,
					type: "line",
					data: this.captchaDailyStats.map(function (v) {
						return v.count / axis.divider;
					}),
					itemStyle: {
						color: "#FBD88A"
					},
					smooth: true
				}
			],
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}

	this.reloadBarChart = function (chartId, name, stats, xFunc, tooltipFunc, axis) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = teaweb.initChart(chartBox)
		let option = {
			xAxis: {
				data: stats.map(xFunc)
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