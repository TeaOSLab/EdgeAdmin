Tea.context(function () {
	this.$delay(function () {
		this.reloadHourlyTrafficChart()
		this.reloadTopDomainsChart()
		this.reloadTopNodesChart()
		this.reloadCPUChart()
	})

	/**
	 * 流量统计
	 */
	this.trafficTab = "hourly"

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
			return stats[args.dataIndex].day + " " + stats[args.dataIndex].hour + "时 流量: " + teaweb.formatBytes(stats[args.dataIndex].bytes)
		})
	}

	this.reloadDailyTrafficChart = function () {
		let stats = this.dailyStats
		this.reloadTrafficChart("daily-traffic-chart", "流量统计", stats, function (args) {
			return stats[args.dataIndex].day + " 流量: " + teaweb.formatBytes(stats[args.dataIndex].bytes)
		})
	}

	this.reloadTrafficChart = function (chartId, name, stats, tooltipFunc) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}

		let axis = teaweb.bytesAxis(stats, function (v) {
			return v.bytes
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
				}
			],
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}

	// 域名排行
	this.reloadTopDomainsChart = function () {
		let that = this
		let axis = teaweb.countAxis(this.topDomainStats, function (v) {
			return v.countRequests
		})
		teaweb.renderBarChart({
			id: "top-domains-chart",
			name: "域名",
			values: this.topDomainStats,
			x: function (v) {
				return v.domainName
			},
			tooltip: function (args, stats) {
				return stats[args.dataIndex].domainName + "<br/>请求数：" + " " + teaweb.formatNumber(stats[args.dataIndex].countRequests) + "<br/>流量：" + teaweb.formatBytes(stats[args.dataIndex].bytes)
			},
			value: function (v) {
				return v.countRequests / axis.divider;
			},
			axis: axis,
			click: function (args, stats) {
				window.location = "/ns/domains/domain?domainId=" + stats[args.dataIndex].domainId
			}
		})
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
				window.location = "/ns/clusters/cluster/node?nodeId=" + stats[args.dataIndex].nodeId + "&clusterId=" + stats[args.dataIndex].clusterId
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