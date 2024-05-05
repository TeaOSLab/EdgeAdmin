Tea.context(function () {
	this.isLoading = true
	this.trafficTab = "hourly"
	this.metricCharts = []
	this.dashboard = {}
	this.localLowerVersionAPINode = null
	this.countWeakAdmins = 0
	this.todayCountIPsFormat = "0"

	this.$delay(function () {
		this.$post("$")
			.success(function (resp) {
				for (let k in resp.data) {
					this[k] = resp.data[k]
				}

				this.todayCountIPsFormat = teaweb.formatNumber(this.todayCountIPs)

				this.isLoading = false

				this.$delay(function () {
					this.reloadHourlyTrafficChart()
					this.reloadTopDomainsChart()
				})
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
		let stats = this.hourlyTrafficStats
		this.reloadTrafficChart("hourly-traffic-chart-box", stats, function (args) {
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
		let stats = this.dailyTrafficStats
		this.reloadTrafficChart("daily-traffic-chart-box", stats, function (args) {
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

	this.reloadTrafficChart = function (chartId, stats, tooltipFunc) {
		let axis = teaweb.bytesAxis(stats, function (v) {
			return v.bytes
		})
		let chartBox = document.getElementById(chartId)
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
					formatter: function (v) {
						return v + axis.unit
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
					name: "总流量",
					type: "line",
					data: stats.map(function (v) {
						return v.bytes / axis.divider;
					}),
					itemStyle: {
						color: teaweb.DefaultChartColor
					},
					lineStyle: {
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
						return v.cachedBytes / axis.divider;
					}),
					itemStyle: {
						color: "#61A0A8"
					},
					lineStyle: {
						color: "#61A0A8"
					},
					areaStyle: {},
					smooth: true
				},
				{
					name: "攻击流量",
					type: "line",
					data: stats.map(function (v) {
						return v.attackBytes / axis.divider;
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
				data: ["总流量", "缓存流量", "攻击流量"]
			},
			animation: false
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
			axis: axis,
			click: function (args, stats) {
				let index = args.dataIndex
				window.location = "/servers/server?serverId=" + stats[index].serverId
			}
		})
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

	// 重启本地API节点
	this.isRestartingLocalAPINode = false
	this.restartAPINode = function () {
		if (this.isRestartingLocalAPINode) {
			return
		}
		if (this.localLowerVersionAPINode == null) {
			return
		}
		this.isRestartingLocalAPINode = true
		this.localLowerVersionAPINode.isRestarting = true
		this.$post("/dashboard/restartLocalAPINode")
			.params({
				"exePath": this.localLowerVersionAPINode.exePath
			})
			.timeout(300)
			.success(function () {
				this.$delay(function () {
					teaweb.reload()
				}, 5000)
			})
			.done(function () {
				this.isRestartingLocalAPINode = false
				this.localLowerVersionAPINode.isRestarting = false
			})
	}

	// 关闭XFF提示
	this.dismissXFFPrompt = function () {
		this.$post("/settings/security/dismissXFFPrompt")
			.success(function () {
				teaweb.reload()
			})
	}
})
