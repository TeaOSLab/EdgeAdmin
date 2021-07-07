Tea.context(function () {
	this.formatCount = function (count) {
		if (count < 1000) {
			return count.toString()
		}
		if (count < 1000 * 1000) {
			return (Math.round(count / 1000 * 100) / 100) + "K"
		}
		return (Math.round(count / 1000 / 1000 * 100) / 100) + "M"
	}

	/**
	 * 流量统计
	 */
	this.trafficTab = "hourly"

	this.$delay(function () {
		this.reloadHourlyTrafficChart()
		this.reloadHourlyRequestsChart()
		this.reloadTopDomainsChart()
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
		let stats = this.hourlyStats
		this.reloadTrafficChart("hourly-traffic-chart", "流量统计", stats, function (args) {
			if (args.seriesIndex == 0) {
				return stats[args.dataIndex].day + " " + stats[args.dataIndex].hour + "时 流量: " + teaweb.formatBytes(stats[args.dataIndex].bytes)
			}
			if (args.seriesIndex == 1) {
				let ratio = 0
				if (stats[args.dataIndex].bytes > 0) {
					ratio = Math.round(stats[args.dataIndex].cachedBytes * 10000 / stats[args.dataIndex].bytes) / 100
				}
				return stats[args.dataIndex].day + " " + stats[args.dataIndex].hour + "时 缓存流量: " + teaweb.formatBytes(stats[args.dataIndex].cachedBytes) + ", 命中率：" + ratio + "%"
			}
			return ""
		})
	}

	this.reloadDailyTrafficChart = function () {
		let stats = this.dailyStats
		this.reloadTrafficChart("daily-traffic-chart", "流量统计", stats, function (args) {
			if (args.seriesIndex == 0) {
				return stats[args.dataIndex].day + " 流量: " + teaweb.formatBytes(stats[args.dataIndex].bytes)
			}
			if (args.seriesIndex == 1) {
				let ratio = 0
				if (stats[args.dataIndex].bytes > 0) {
					ratio = Math.round(stats[args.dataIndex].cachedBytes * 10000 / stats[args.dataIndex].bytes) / 100
				}
				return stats[args.dataIndex].day + " 缓存流量: " + teaweb.formatBytes(stats[args.dataIndex].cachedBytes) + ", 命中率：" + ratio + "%"
			}
			return ""
		})
	}

	this.reloadTrafficChart = function (chartId, name, stats, tooltipFunc) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}

		let axis = teaweb.bytesAxis(stats, function (v) {
			return Math.max(v.bytes, v.cachedBytes)
		})

		let chart = echarts.init(chartBox)
		let option = {
			xAxis: {
				data: stats.map(function (v) {
					if (v.hour != null) {
						return v.hour
					}
					return v.day
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
						color: "#9DD3E8"
					},
					areaStyle: {
						color: "#9DD3E8"
					}
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
					}
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

	/**
	 * 请求数统计
	 */
	this.requestsTab = "hourly"

	this.selectRequestsTab = function (tab) {
		this.requestsTab = tab
		if (tab == "hourly") {
			this.$delay(function () {
				this.reloadHourlyRequestsChart()
			})
		} else if (tab == "daily") {
			this.$delay(function () {
				this.reloadDailyRequestsChart()
			})
		}
	}

	this.reloadHourlyRequestsChart = function () {
		let stats = this.hourlyStats
		this.reloadRequestsChart("hourly-requests-chart", "请求数统计", stats, function (args) {
			if (args.seriesIndex == 0) {
				return stats[args.dataIndex].day + " " + stats[args.dataIndex].hour + "时 请求数: " + teaweb.formatNumber(stats[args.dataIndex].countRequests)
			}
			if (args.seriesIndex == 1) {
				let ratio = 0
				if (stats[args.dataIndex].countRequests > 0) {
					ratio = Math.round(stats[args.dataIndex].countCachedRequests * 10000 / stats[args.dataIndex].countRequests) / 100
				}
				return stats[args.dataIndex].day + " " + stats[args.dataIndex].hour + "时 缓存请求数: " + teaweb.formatNumber(stats[args.dataIndex].countCachedRequests) + ", 命中率：" + ratio + "%"
			}
			return ""
		})
	}

	this.reloadDailyRequestsChart = function () {
		let stats = this.dailyStats
		this.reloadRequestsChart("daily-requests-chart", "请求数统计", stats, function (args) {
			if (args.seriesIndex == 0) {
				return stats[args.dataIndex].day + " 请求数: " + teaweb.formatNumber(stats[args.dataIndex].countRequests)
			}
			if (args.seriesIndex == 1) {
				let ratio = 0
				if (stats[args.dataIndex].countRequests > 0) {
					ratio = Math.round(stats[args.dataIndex].countCachedRequests * 10000 / stats[args.dataIndex].countRequests) / 100
				}
				return stats[args.dataIndex].day + " 缓存请求数: " + teaweb.formatNumber(stats[args.dataIndex].countCachedRequests) + ", 命中率：" + ratio + "%"
			}
			return ""
		})
	}

	this.reloadRequestsChart = function (chartId, name, stats, tooltipFunc) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}

		let axis = teaweb.countAxis(stats, function (v) {
			return Math.max(v.countRequests, v.countCachedRequests)
		})

		let chart = echarts.init(chartBox)
		let option = {
			xAxis: {
				data: stats.map(function (v) {
					if (v.hour != null) {
						return v.hour
					}
					if (v.day != null) {
						return v.day
					}
					return ""
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
						color: "#9DD3E8"
					},
					areaStyle: {
						color: "#9DD3E8"
					}
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
					}
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

	// 域名排行
	this.reloadTopDomainsChart = function () {
		let axis = teaweb.countAxis(this.topDomainStats, function (v) {
			return v.countRequests
		})
		teaweb.renderBarChart({
			id: "top-domains-chart",
			name: "域名",
			values: this.topDomainStats,
			x: function (v) {
				return v.domain
			},
			tooltip: function (args, stats) {
				return stats[args.dataIndex].domain + "<br/>请求数：" + " " + teaweb.formatNumber(stats[args.dataIndex].countRequests) + "<br/>流量：" + teaweb.formatBytes(stats[args.dataIndex].bytes)
			},
			value: function (v) {
				return v.countRequests / axis.divider;
			},
			axis: axis
		})
	}
})