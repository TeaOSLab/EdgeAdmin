Tea.context(function () {
	this.isLoading = true
	this.board = {}
	this.metricCharts = []


	/**
	 * 流量统计
	 */
	this.trafficTab = "hourly"

	this.$delay(function () {
		this.$post("$")
			.params({
				clusterId: this.clusterId
			})
			.timeout(30)
			.success(function (resp) {
				for (let k in resp.data) {
					this[k] = resp.data[k]
				}

				this.reloadHourlyTrafficChart()
				this.reloadHourlyRequestsChart()
				this.reloadTopNodesChart()
				this.reloadTopDomainsChart()
				this.reloadCPUChart()

				this.isLoading = false
			})
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
			let index = args.dataIndex
			let cachedRatio = 0
			let attackRatio = 0
			if (stats[index].bytes > 0) {
				cachedRatio = Math.round(stats[index].cachedBytes * 10000 / stats[index].bytes) / 100
				attackRatio = Math.round(stats[index].attackBytes * 10000 / stats[index].bytes) / 100
			}

			return stats[index].day + " " + stats[index].hour + "时<br/>总流量：" + teaweb.formatBytes(stats[index].bytes) + "<br/>缓存流量：" + teaweb.formatBytes(stats[index].cachedBytes) + "<br/>缓存命中率：" + cachedRatio + "%<br/>拦截攻击流量：" + teaweb.formatBytes(stats[index].attackBytes) + "<br/>拦截比例：" + attackRatio + "%"
		})
	}

	this.reloadDailyTrafficChart = function () {
		let stats = this.dailyStats
		this.reloadTrafficChart("daily-traffic-chart", "流量统计", stats, function (args) {
			let index = args.dataIndex
			let cachedRatio = 0
			let attackRatio = 0
			if (stats[index].bytes > 0) {
				cachedRatio = Math.round(stats[index].cachedBytes * 10000 / stats[index].bytes) / 100
				attackRatio = Math.round(stats[index].attackBytes * 10000 / stats[index].bytes) / 100
			}

			return stats[index].day + "<br/>总流量：" + teaweb.formatBytes(stats[index].bytes) + "<br/>缓存流量：" + teaweb.formatBytes(stats[index].cachedBytes) + "<br/>缓存命中率：" + cachedRatio + "%<br/>拦截攻击流量：" + teaweb.formatBytes(stats[index].attackBytes) + "<br/>拦截比例：" + attackRatio + "%"
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

		let chart = teaweb.initChart(chartBox)
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
				},
				{
					name: "攻击流量",
					type: "line",
					data: stats.map(function (v) {
						return v.attackBytes / axis.divider
					}),
					itemStyle: {
						color: "#F39494"
					},
					areaStyle: {
						color: "#F39494"
					},
					smooth: true
				}
			],
			legend: {
				data: ["流量", "缓存流量", "攻击流量"]
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
			let index = args.dataIndex
			let cachedRatio = 0
			let attackRatio = 0
			if (stats[index].countRequests > 0) {
				cachedRatio = Math.round(stats[index].countCachedRequests * 10000 / stats[index].countRequests) / 100
				attackRatio = Math.round(stats[index].countAttackRequests * 10000 / stats[index].countRequests) / 100
			}

			return stats[index].day + " " + stats[index].hour + "时<br/>总请求数：" + stats[index].countRequests + "<br/>缓存请求数：" + stats[index].countCachedRequests + "<br/>缓存命中率：" + cachedRatio + "%<br/>拦截攻击数：" + stats[index].countAttackRequests + "<br/>拦截比例：" + attackRatio + "%"
		})
	}

	this.reloadDailyRequestsChart = function () {
		let stats = this.dailyStats
		this.reloadRequestsChart("daily-requests-chart", "请求数统计", stats, function (args) {
			let index = args.dataIndex
			let cachedRatio = 0
			let attackRatio = 0
			if (stats[index].countRequests > 0) {
				cachedRatio = Math.round(stats[index].countCachedRequests * 10000 / stats[index].countRequests) / 100
				attackRatio = Math.round(stats[index].countAttackRequests * 10000 / stats[index].countRequests) / 100
			}

			return stats[index].day + "<br/>总请求数：" + stats[index].countRequests + "<br/>缓存请求数：" + stats[index].countCachedRequests + "<br/>缓存命中率：" + cachedRatio + "%<br/>拦截攻击数：" + stats[index].countAttackRequests + "<br/>拦截比例：" + attackRatio + "%"
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

		let chart = teaweb.initChart(chartBox)
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
				},
				{
					name: "攻击请求数",
					type: "line",
					data: stats.map(function (v) {
						return v.countAttackRequests / axis.divider;
					}),
					itemStyle: {
						color: "#F39494"
					},
					areaStyle: {
						color: "#F39494"
					},
					smooth: true
				}
			],
			legend: {
				data: ["请求数", "缓存请求数", "攻击请求数"]
			},
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}

	// 节点排行
	this.reloadTopNodesChart = function () {
		let that = this
		let axis = teaweb.countAxis(this.topNodeStats, function (v) {
			return v.countRequests
		})
		teaweb.renderBarChart({
			id: "top-nodes-chart",
			name: "节点",
			values: this.topNodeStats,
			x: function (v) {
				return v.nodeName
			},
			tooltip: function (args, stats) {
				return stats[args.dataIndex].nodeName + "<br/>请求数：" + " " + teaweb.formatNumber(stats[args.dataIndex].countRequests) + "<br/>流量：" + teaweb.formatBytes(stats[args.dataIndex].bytes)
			},
			value: function (v) {
				return v.countRequests / axis.divider;
			},
			axis: axis,
			click: function (args, stats) {
				window.location = "/clusters/cluster/node?nodeId=" + stats[args.dataIndex].nodeId + "&clusterId=" + that.clusterId
			}
		})
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
			axis: axis,
			click: function (args, stats) {
				window.location = "/servers/server?serverId=" + stats[args.dataIndex].serverId
			}
		})
	}

	/**
	 * 系统信息
	 */
	this.nodeStatusTab = "cpu"

	this.selectNodeStatusTab = function (tab) {
		this.nodeStatusTab = tab
		this.$delay(function () {
			switch (tab) {
				case "cpu":
					this.reloadCPUChart()
					break
				case "memory":
					this.reloadMemoryChart()
					break
				case "load":
					this.reloadLoadChart()
					break
			}
		})
	}

	this.reloadCPUChart = function () {
		let axis = {unit: "%", divider: 1}
		teaweb.renderLineChart({
			id: "cpu-chart",
			name: "CPU",
			values: this.cpuValues,
			x: function (v) {
				return v.time
			},
			tooltip: function (args, stats) {
				return stats[args.dataIndex].time + "：" + (Math.ceil(stats[args.dataIndex].value * 100 * 100) / 100) + "%"
			},
			value: function (v) {
				return v.value * 100;
			},
			axis: axis,
			max: 100
		})
	}

	this.reloadMemoryChart = function () {
		let axis = {unit: "%", divider: 1}
		teaweb.renderLineChart({
			id: "memory-chart",
			name: "内存",
			values: this.memoryValues,
			x: function (v) {
				return v.time
			},
			tooltip: function (args, stats) {
				return stats[args.dataIndex].time + "：" + (Math.ceil(stats[args.dataIndex].value * 100 * 100) / 100) + "%"
			},
			value: function (v) {
				return v.value * 100;
			},
			axis: axis,
			max: 100
		})
	}

	this.reloadLoadChart = function () {
		let axis = {unit: "", divider: 1}
		let max = this.loadValues.$map(function (k, v) {
			return v.value
		}).$max()
		if (max < 10) {
			max = 10
		} else if (max < 20) {
			max = 20
		} else if (max < 100) {
			max = 100
		} else {
			max = null
		}
		teaweb.renderLineChart({
			id: "load-chart",
			name: "负载",
			values: this.loadValues,
			x: function (v) {
				return v.time
			},
			tooltip: function (args, stats) {
				return stats[args.dataIndex].time + "：" + (Math.ceil(stats[args.dataIndex].value * 100) / 100)
			},
			value: function (v) {
				return v.value;
			},
			axis: axis,
			max: max
		})
	}
})