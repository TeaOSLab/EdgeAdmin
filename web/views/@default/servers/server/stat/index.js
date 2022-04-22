Tea.context(function () {
	this.$delay(function () {
		let that = this

		// 每N分钟取一次
		let newStats = []
		for (let i = 0; i < this.minutelyStats.length; i++) {
			if (i % 5 == 0) {
				newStats.push(this.minutelyStats[i])
			}
		}
		this.minutelyStats = newStats

		this.reloadRequestsChart("minutely-requests-chart", "请求数统计", this.minutelyStats, function (args) {
			if (args.seriesIndex == 0) {
				return that.minutelyStats[args.dataIndex].day + " " + that.minutelyStats[args.dataIndex].minute + "<br/>请求数: " + teaweb.formatNumber(that.minutelyStats[args.dataIndex].countRequests)
			}
			if (args.seriesIndex == 1) {
				let ratio = 0
				if (that.minutelyStats[args.dataIndex].countRequests > 0) {
					ratio = Math.round(that.minutelyStats[args.dataIndex].countCachedRequests * 10000 / that.minutelyStats[args.dataIndex].countRequests) / 100
				}
				return that.minutelyStats[args.dataIndex].day + " " + that.minutelyStats[args.dataIndex].minute + "<br/>缓存请求数: " + teaweb.formatNumber(that.minutelyStats[args.dataIndex].countCachedRequests) + "<br/>命中率：" + ratio + "%"
			}
			return ""
		})
		this.reloadTrafficChart("minutely-traffic-chart", "流量统计", this.minutelyStats, function (args) {
			if (args.seriesIndex == 0) {
				return that.minutelyStats[args.dataIndex].day + " " + that.minutelyStats[args.dataIndex].minute + "<br/>流量: " + teaweb.formatBytes(that.minutelyStats[args.dataIndex].bytes)
			}
			if (args.seriesIndex == 1) {
				let ratio = 0
				if (that.minutelyStats[args.dataIndex].bytes > 0) {
					ratio = Math.round(that.minutelyStats[args.dataIndex].cachedBytes * 10000 / that.minutelyStats[args.dataIndex].bytes) / 100
				}
				return that.minutelyStats[args.dataIndex].day + " " + that.minutelyStats[args.dataIndex].minute + "<br/>缓存流量: " + teaweb.formatBytes(that.minutelyStats[args.dataIndex].cachedBytes) + "<br/>命中率：" + ratio + "%"
			}
			return ""
		})
	})

	this.reloadRequestsChart = function (chartId, name, stats, tooltipFunc) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}

		let axis = teaweb.countAxis(stats, function (v) {
			return Math.max(v.countRequests, v.countCachedRequests)
		})

		let chart = teaweb.initChart(chartBox)
		let option = {
			xAxis: {
				data: stats.map(function (v) {
					return v.minute
				})
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
				left: 50,
				top: 40,
				right: 20,
				bottom: 20
			},
			series: [
				{
					name: "请求数",
					type: "line",
					data: stats.map(function (v) {
						return v.countRequests / axis.divider
					}),
					itemStyle: {
						color: teaweb.DefaultChartColor
					},
					areaStyle: {
						color: teaweb.DefaultChartColor
					},
					smooth: true
				},
				{
					name: "缓存请求数",
					type: "line",
					data: stats.map(function (v) {
						return v.countCachedRequests / axis.divider
					}),
					itemStyle: {
						color: "#61A0A8"
					},
					areaStyle: {
						color: "#61A0A8"
					},
					smooth: true
				}
			],
			legend: {
				data: ['请求数', '缓存请求数']
			},
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}

	this.reloadTrafficChart = function (chartId, name, stats, tooltipFunc) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}

		let axis = teaweb.bytesAxis(stats, function (v) {
			return Math.max(v.bytes, v.cachedBytes)
		})

		let chart = teaweb.initChart(chartBox)
		let option = {
			xAxis: {
				data: stats.map(function (v) {
					return v.minute
				})
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
				left: 50,
				top: 40,
				right: 20,
				bottom: 20
			},
			series: [
				{
					name: "流量",
					type: "line",
					data: stats.map(function (v) {
						return v.bytes / axis.divider
					}),
					itemStyle: {
						color: teaweb.DefaultChartColor
					},
					areaStyle: {
						color: teaweb.DefaultChartColor
					},
					smooth: true
				},
				{
					name: "缓存流量",
					type: "line",
					data: stats.map(function (v) {
						return v.cachedBytes / axis.divider
					}),
					itemStyle: {
						color: "#61A0A8"
					},
					areaStyle: {
						color: "#61A0A8"
					},
					smooth: true
				}
			],
			legend: {
				data: ['流量', '缓存流量']
			},
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}
})
