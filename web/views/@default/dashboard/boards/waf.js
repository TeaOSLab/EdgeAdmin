Tea.context(function () {
	this.$delay(function () {
		let that = this

		this.board.countDailyBlocks = teaweb.formatCount(this.board.countDailyBlocks)
		this.board.countDailyCaptcha = teaweb.formatCount(this.board.countDailyCaptcha)
		this.board.countDailyLogs = teaweb.formatCount(this.board.countDailyLogs)
		this.board.countWeeklyBlocks = teaweb.formatCount(this.board.countWeeklyBlocks)

		this.reloadHourlyChart()
		this.reloadGroupChart()
		this.reloadAccessLogs()
	})

	this.requestsTab = "hourly"

	this.selectRequestsTab = function (tab) {
		this.requestsTab = tab
		this.$delay(function () {
			switch (tab) {
				case "hourly":
					this.reloadHourlyChart()
					break
				case "daily":
					this.reloadDailyChart()
					break
			}
		})
	}

	this.reloadHourlyChart = function () {
		let axis = teaweb.countAxis(this.hourlyStats, function (v) {
			return [v.countLogs, v.countCaptcha, v.countBlocks].$max()
		})
		let that = this
		this.reloadLineChart("hourly-chart", "按小时统计", this.hourlyStats, function (v) {
			return v.hour.substring(8, 10)
		}, function (args) {
			let index = args.dataIndex
			let hour = that.hourlyStats[index].hour.substring(0, 4) + "-" + that.hourlyStats[index].hour.substring(4, 6) + "-" + that.hourlyStats[index].hour.substring(6, 8) + " " + that.hourlyStats[index].hour.substring(8)
			return hour + "时<br/>拦截: "
				+ teaweb.formatNumber(that.hourlyStats[index].countBlocks) + "<br/>验证码: " + teaweb.formatNumber(that.hourlyStats[index].countCaptcha) + "<br/>记录: " + teaweb.formatNumber(that.hourlyStats[index].countLogs)
		}, axis)
	}

	this.reloadDailyChart = function () {
		let axis = teaweb.countAxis(this.dailyStats, function (v) {
			return [v.countLogs, v.countCaptcha, v.countBlocks].$max()
		})
		let that = this
		this.reloadLineChart("daily-chart", "按天统计", this.dailyStats, function (v) {
			return v.day.substring(4, 6) + "月" + v.day.substring(6, 8) + "日"
		}, function (args) {
			let index = args.dataIndex
			let day = that.dailyStats[index].day.substring(0, 4) + "-" + that.dailyStats[index].day.substring(4, 6) + "-" + that.dailyStats[index].day.substring(6, 8)
			return day + "<br/>拦截: "
				+ teaweb.formatNumber(that.dailyStats[index].countBlocks) + "<br/>验证码: " + teaweb.formatNumber(that.dailyStats[index].countCaptcha) + "<br/>记录: " + teaweb.formatNumber(that.dailyStats[index].countLogs)
		}, axis)
	}

	this.reloadLineChart = function (chartId, name, stats, xFunc, tooltipFunc, axis) {
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
				left: 42,
				top: 10,
				right: 20,
				bottom: 20
			},
			series: [
				{
					name: name,
					type: "line",
					data: stats.map(function (v) {
						return v.countLogs / axis.divider;
					}),
					itemStyle: {
						color: "#879BD7"
					},
					areaStyle: {},
					stack: "总量"
				},
				{
					name: name,
					type: "line",
					data: stats.map(function (v) {
						return v.countCaptcha / axis.divider;
					}),
					itemStyle: {
						color: "#FBD88A"
					},
					areaStyle: {},
					stack: "总量"
				},
				{
					name: name,
					type: "line",
					data: stats.map(function (v) {
						return v.countBlocks / axis.divider;
					}),
					itemStyle: {
						color: "#F39494"
					},
					areaStyle: {},
					stack: "总量"
				}
			],
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}

	this.reloadGroupChart = function () {
		let axis = teaweb.countAxis(this.groupStats, function (v) {
			return v.count
		})
		teaweb.renderBarChart({
			id: "group-chart",
			values: this.groupStats,
			x: function (v) {
				return v.name
			},
			value: function (v) {
				return v.count / axis.divider
			},
			tooltip: function (args, stats) {
				let index = args.dataIndex
				return stats[index].name + ": " + stats[index].count
			},
			axis: axis
		})
	}

	this.accessLogs = []
	this.reloadAccessLogs = function () {
		this.$post(".wafLogs")
			.success(function (resp) {
				if (resp.data.accessLogs != null) {
					let that = this
					resp.data.accessLogs.forEach(function (v) {
						that.formatTime(v)
					})
					this.accessLogs = resp.data.accessLogs
				}
			})
			.done(function () {
				this.$delay(this.reloadAccessLogs, 10000)
			})
	}

	this.formatTime = function (accessLog) {
		let elapsedSeconds = Math.ceil(new Date().getTime() / 1000) - accessLog.timestamp
		if (elapsedSeconds >= 0) {
			if (elapsedSeconds < 60) {
				accessLog.humanTime = elapsedSeconds + "秒前"
			} else if (elapsedSeconds < 3600) {
				accessLog.humanTime = Math.ceil(elapsedSeconds / 60) + "分钟前"
			} else if (elapsedSeconds < 3600 * 24) {
				accessLog.humanTime = Math.ceil(elapsedSeconds / 3600) + "小时前"
			}
		}
	}
})