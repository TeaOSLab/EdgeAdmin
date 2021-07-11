Tea.context(function () {
	this.$delay(function () {
		this.reloadDailyStats()
		this.reloadCPUChart()
		this.reloadTopTrafficChart()
	})

	this.reloadDailyStats = function () {
		let axis = teaweb.countAxis(this.dailyStats, function (v) {
			return v.count
		})
		let max = axis.max
		if (max < 10) {
			max = 10
		} else if (max < 100) {
			max = 100
		}
		teaweb.renderLineChart({
			id: "daily-stat-chart",
			name: "用户",
			values: this.dailyStats,
			x: function (v) {
				return v.day.substring(4, 6) + "-" + v.day.substring(6)
			},
			tooltip: function (args, stats) {
				let index = args.dataIndex
				return stats[index].day.substring(4, 6) + "-" + stats[index].day.substring(6) + "：" + stats[index].count
			},
			value: function (v) {
				return v.count;
			},
			axis: axis,
			max: max
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

	// 流量排行
	this.reloadTopTrafficChart = function () {
		let that = this
		let axis = teaweb.bytesAxis(this.topTrafficStats, function (v) {
			return v.bytes
		})
		teaweb.renderBarChart({
			id: "top-traffic-chart",
			name: "流量",
			values: this.topTrafficStats,
			x: function (v) {
				return v.userName
			},
			tooltip: function (args, stats) {
				let index = args.dataIndex
				return stats[index].userName + "<br/>请求数：" + " " + teaweb.formatNumber(stats[index].countRequests) + "<br/>流量：" + teaweb.formatBytes(stats[index].bytes)
			},
			value: function (v) {
				return v.bytes / axis.divider;
			},
			axis: axis
		})
	}
})