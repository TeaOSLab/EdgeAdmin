// 指标图表
Vue.component("metric-chart", {
	props: ["v-chart", "v-stats", "v-period-unit", "v-value-type"],
	mounted: function () {
		this.load()
	},
	data: function () {
		let stats = this.vStats
		if (stats == null) {
			stats = []
		}
		if (stats.length > 0) {
			let sum = stats.$sum(function (k, v) {
				return v.value
			})
			if (sum < stats[0].total) {
				stats.push({
					keys: ["其他"],
					value: stats[0].total - sum,
					total: stats[0].total,
					time: stats[0].time
				})
			}
		}
		if (this.vChart.maxItems > 0) {
			stats = stats.slice(0, this.vChart.maxItems)
		}

		stats.$rsort(function (v1, v2) {
			return v1.value - v2.value
		})

		let widthPercent = 100
		if (this.vChart.widthDiv > 0) {
			widthPercent = 100 / this.vChart.widthDiv
		}

		return {
			chart: this.vChart,
			stats: stats,
			width: widthPercent + "%",
			chartId: "metric-chart-" + this.vChart.id
		}
	},
	methods: {
		load: function () {
			var el = document.getElementById(this.chartId)
			if (el == null || el.offsetWidth == 0 || el.offsetHeight == 0) {
				setTimeout(this.load, 100)
			} else {
				this.render(el)
			}
		},
		render: function (el) {
			let chart = echarts.init(el)
			switch (this.chart.type) {
				case "pie":
					this.renderPie(chart)
					break
				case "bar":
					this.renderBar(chart)
					break
				case "timeBar":
					this.renderTimeBar(chart)
					break
				case "timeLine":
					this.renderTimeLine(chart)
					break
				case "table":
					this.renderTable(chart)
					break
			}
		},
		renderPie: function (chart) {
			let values = this.stats.map(function (v) {
				return {
					name: v.keys[0],
					value: v.value
				}
			})
			let that = this
			chart.setOption({
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let percent = 0
						if (stat.total > 0) {
							percent = Math.round((stat.value * 100 / stat.total) * 100) / 100
						}
						return stat.keys[0] + ":" + stat.value + "，占比：" + percent + "%"
					}
				},
				series: [
					{
						name: name,
						type: "pie",
						data: values,
						areaStyle: {}
					}
				]
			})
		},
		renderTimeBar: function (chart) {
			this.stats.$sort(function (v1, v2) {
				return (v1.time < v2.time) ? -1 : 1
			})
			let values = this.stats.map(function (v) {
				return v.value
			})

			let axis = {unit: "", divider: 1}
			switch (this.vValueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "bytes":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}

			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return that.formatTime(v.time)
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
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						return that.formatTime(stat.time) + ": " + stat.value
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
						type: "bar",
						data: values.map(function (v){
							return v/axis.divider
						}),
						itemStyle: {
							color: "#9DD3E8"
						},
						areaStyle: {},
						barWidth: "20em"
					}
				]
			})
		},
		renderTimeLine: function (chart) {
			this.stats.$sort(function (v1, v2) {
				return (v1.time < v2.time) ? -1 : 1
			})
			let values = this.stats.map(function (v) {
				return v.value
			})

			let axis = {unit: "", divider: 1}
			switch (this.vValueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "bytes":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}

			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return that.formatTime(v.time)
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
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						return that.formatTime(stat.time) + ": " + stat.value
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
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: "#9DD3E8"
						},
						areaStyle: {}
					}
				]
			})
		},
		renderBar: function (chart) {
			let values = this.stats.map(function (v) {
				return v.value
			})
			let axis = {unit: "", divider: 1}
			switch (this.vValueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "bytes":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}

			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return v.keys[0]
					}),
					axisLabel: {
						interval: 0
					}
				},
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let percent = 0
						if (stat.total > 0) {
							percent = Math.round((stat.value * 100 / stat.total) * 100) / 100
						}
						return stat.keys[0] + ": " + stat.value + "，占比：" + percent + "%"
					}
				},
				yAxis: {
					axisLabel: {
						formatter: function (value) {
							return value + axis.unit
						}
					}
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
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: "#9DD3E8"
						},
						areaStyle: {},
						barWidth: "20em"
					}
				]
			})
		},
		renderTable: function (chart) {
			let table = `<table class="ui table celled">
	<thead>
		<tr>
			<th>对象</th>
			<th>数值</th>
			<th>占比</th>
		</tr>
	</thead>`
			this.stats.forEach(function (v) {
				table += "<tr><td>" + v.keys[0] + "</td><td>" + v.value + "</td>"
				let percent = 0
				if (v.total > 0) {
					percent = Math.round((v.value * 100 / v.total) * 100) / 100
				}
				table += "<td><div class=\"ui progress blue\"><div class=\"bar\" style=\"min-width: 0; height: 4px; width: " + percent + "%\"></div></div>" + percent + "%</td>"
				table += "</tr>"
			})

			table += `</table>`
			document.getElementById(this.chartId).innerHTML = table
		},
		formatTime: function (time) {
			if (time == null) {
				return ""
			}
			switch (this.vPeriodUnit) {
				case "month":
					return time.substring(0, 4) + "-" + time.substring(4, 6)
				case "week":
					return time.substring(0, 4) + "-" + time.substring(4, 6)
				case "day":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8)
				case "hour":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8) + " " + time.substring(8, 10)
				case "minute":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8) + " " + time.substring(8, 10) + ":" + time.substring(10, 12)
			}
			return time
		}
	},
	template: `<div style="float: left" :style="{'width': width}">
	<h4>{{chart.name}} <span>（指标）</span></h4>
	<div class="ui divider"></div>
	<div style="height: 20em; padding-bottom: 1em; " :id="chartId" :class="{'scroll-box': chart.type == 'table'}"></div>
</div>`
})

Vue.component("metric-board", {
	template: `<div><slot></slot></div>`
})