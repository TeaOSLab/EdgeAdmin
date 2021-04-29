Tea.context(function () {
	this.$delay(function () {
		this.loadTrafficInChart()
		this.loadTrafficOutChart()
		this.loadConnectionsChart()
		this.loadCPUChart()
		this.loadMemoryChart()
		this.loadLoadChart()

		let that = this
		window.addEventListener("resize", function () {
			that.resizeChart("traffic-in-chart")
			that.resizeChart("traffic-out-chart")
			that.resizeChart("connections-chart")
			that.resizeChart("cpu-chart")
			that.resizeChart("memory-chart")
			that.resizeChart("load-chart")
		})
	})

	this.loadTrafficInChart = function () {
		this.$post(".trafficIn")
			.params({
				nodeId: this.nodeId
			})
			.success(function (resp) {
				let values = resp.data.values
				let maxFunc = function () {
					let max = values.map(function (v) {
						return v.value
					}).$max() / 1024 / 1024
					if (max < 1) {
						return 1
					}
					if (max < 10) {
						return 10
					}
					if (max < 100) {
						return 100
					}
					return null
				}
				let valueFunc = function (v) {
					return v.value / 1024 / 1024
				}
				this.reloadChart("traffic-in-chart", "", values, "M", maxFunc, valueFunc)
			})
			.done(function () {
				this.$delay(function () {
					this.loadTrafficInChart()
				}, 30000)
			})
	}

	this.loadTrafficOutChart = function () {
		this.$post(".trafficOut")
			.params({
				nodeId: this.nodeId
			})
			.success(function (resp) {
				let values = resp.data.values
				let maxFunc = function () {
					let max = values.map(function (v) {
						return v.value
					}).$max() / 1024 / 1024
					if (max < 1) {
						return 1
					}
					if (max < 10) {
						return 10
					}
					if (max < 100) {
						return 100
					}
					return null
				}
				let valueFunc = function (v) {
					return v.value / 1024 / 1024
				}
				this.reloadChart("traffic-out-chart", "", values, "M", maxFunc, valueFunc)
			})
			.done(function () {
				this.$delay(function () {
					this.loadTrafficOutChart()
				}, 30000)
			})
	}

	this.loadConnectionsChart = function () {
		this.$post(".connections")
			.params({
				nodeId: this.nodeId
			})
			.success(function (resp) {
				let values = resp.data.values
				let maxFunc = function () {
					let max = values.map(function (v) {
						return v.value
					}).$max()
					if (max < 10) {
						return 10
					}
					if (max < 100) {
						return 100
					}
					if (max < 1000) {
						return 1000
					}
					return null
				}
				let valueFunc = function (v) {
					return v.value
				}
				this.reloadChart("connections-chart", "", values, "", maxFunc, valueFunc)
			})
			.done(function () {
				this.$delay(function () {
					this.loadConnectionsChart()
				}, 30000)
			})
	}

	this.loadCPUChart = function () {
		this.$post(".cpu")
			.params({
				nodeId: this.nodeId
			})
			.success(function (resp) {
				let values = resp.data.values
				let maxFunc = function () {
					return 100
				}
				let valueFunc = function (v) {
					return v.value
				}
				this.reloadChart("cpu-chart", "", values, "%", maxFunc, valueFunc)
			})
			.done(function () {
				this.$delay(function () {
					this.loadCPUChart()
				}, 30000)
			})
	}

	this.loadMemoryChart = function () {
		this.$post(".memory")
			.params({
				nodeId: this.nodeId
			})
			.success(function (resp) {
				let values = resp.data.values
				let maxFunc = function () {
					return 100
				}
				let valueFunc = function (v) {
					return v.value
				}
				this.reloadChart("memory-chart", "", values, "%", maxFunc, valueFunc)
			})
			.done(function () {
				this.$delay(function () {
					this.loadMemoryChart()
				}, 30000)
			})
	}

	this.loadLoadChart = function () {
		this.$post(".load")
			.params({
				nodeId: this.nodeId
			})
			.success(function (resp) {
				let values = resp.data.values
				let maxFunc = function () {
					let max = values.map(function (v) {
						return v.value
					}).$max()
					if (max < 10) {
						return 10
					}
					return null
				}
				let valueFunc = function (v) {
					return v.value
				}
				this.reloadChart("load-chart", "5分钟", values, "", maxFunc, valueFunc)
			})
			.done(function () {
				this.$delay(function () {
					this.loadLoadChart()
				}, 30000)
			})
	}

	this.reloadChart = function (chartId, name, stats, unit, maxFunc, valueFunc) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = echarts.init(chartBox)
		let option = {
			xAxis: {
				data: stats.map(function (stat) {
					return stat.label
				})
			},
			yAxis: {
				max: maxFunc(),
				axisLabel: {
					formatter: function (value) {
						return value + unit
					}
				}
			},
			tooltip: {
				show: true,
				trigger: "item",
				formatter: function (args) {
					return stats[args.dataIndex].label + ": " + stats[args.dataIndex].text
				}
			},
			grid: {
				left: 50,
				top: 10,
				right: 20,
				bottom: 20
			},
			series: [
				{
					name: name,
					type: "line",
					data: stats.map(valueFunc),
					itemStyle: {
						color: "#9DD3E8"
					},
					areaStyle: {}
				}
			],
			animation: true
		}
		chart.setOption(option)
		chart.resize()
	}

	this.resizeChart = function (chartId) {
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = echarts.init(chartBox)
		chart.resize()
	}
})