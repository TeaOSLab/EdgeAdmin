// 指标图表
Vue.component("metric-chart", {
	props: ["v-chart", "v-stats", "v-item", "v-column" /** in column? **/],
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
				if (this.vChart.type == "pie") {
					stats.push({
						keys: ["其他"],
						value: stats[0].total - sum,
						total: stats[0].total,
						time: stats[0].time
					})
				}
			}
		}
		if (this.vChart.maxItems > 0) {
			stats = stats.slice(0, this.vChart.maxItems)
		} else {
			stats = stats.slice(0, 10)
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
			item: this.vItem,
			width: widthPercent + "%",
			chartId: "metric-chart-" + this.vChart.id,
			valueTypeName: (this.vItem != null && this.vItem.valueTypeName != null && this.vItem.valueTypeName.length > 0) ? this.vItem.valueTypeName : ""
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
			window.addEventListener("resize", function () {
				chart.resize()
			})
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
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
							case "count":
								value = teaweb.formatNumber(value)
								break
						}
						return stat.keys[0] + "<br/>" + that.valueTypeName + ": " + value + "<br/>占比：" + percent + "%"
					}
				},
				series: [
					{
						name: name,
						type: "pie",
						data: values,
						areaStyle: {},
						color: ["#9DD3E8", "#B2DB9E", "#F39494", "#FBD88A", "#879BD7"]
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
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
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
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
						}
						return that.formatTime(stat.time) + ": " + value
					}
				},
				grid: {
					left: 50,
					top: 10,
					right: 20,
					bottom: 25
				},
				series: [
					{
						name: name,
						type: "bar",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: teaweb.DefaultChartColor
						},
						areaStyle: {},
						barWidth: "10em"
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
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
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
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
						}
						return that.formatTime(stat.time) + ": " + value
					}
				},
				grid: {
					left: 50,
					top: 10,
					right: 20,
					bottom: 25
				},
				series: [
					{
						name: name,
						type: "line",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: teaweb.DefaultChartColor
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
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}
			let bottom = 24
			let rotate = 0
			let result = teaweb.xRotation(chart, this.stats.map(function (v) {
				return v.keys[0]
			}))
			if (result != null) {
				bottom = result[0]
				rotate = result[1]
			}
			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return v.keys[0]
					}),
					axisLabel: {
						interval: 0,
						rotate: rotate
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
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
							case "count":
								value = teaweb.formatNumber(value)
								break
						}
						return stat.keys[0] + "<br/>" + that.valueTypeName + "：" + value + "<br/>占比：" + percent + "%"
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
					bottom: bottom
				},
				series: [
					{
						name: name,
						type: "bar",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: teaweb.DefaultChartColor
						},
						areaStyle: {},
						barWidth: "10em"
					}
				]
			})

			if (this.item.keys != null) {
				// IP相关操作
				if (this.item.keys.$contains("${remoteAddr}")) {
					let that = this
					chart.on("click", function (args) {
						let index = that.item.keys.$indexesOf("${remoteAddr}")[0]
						let value = that.stats[args.dataIndex].keys[index]
						teaweb.popup("/servers/ipbox?ip=" + value, {
							width: "50em",
							height: "30em"
						})
					})
				}
			}
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
			let that = this
			this.stats.forEach(function (v) {
				let value = v.value
				switch (that.item.valueType) {
					case "byte":
						value = teaweb.formatBytes(value)
						break
				}
				table += "<tr><td>" + v.keys[0] + "</td><td>" + value + "</td>"
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
			switch (this.item.periodUnit) {
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
	template: `<div style="float: left" :style="{'width': this.vColumn ?  '' : width}" :class="{'ui column':this.vColumn}">
	<h4>{{chart.name}} <span>（{{valueTypeName}}）</span></h4>
	<div class="ui divider"></div>
	<div style="height: 14em; padding-bottom: 1em; " :id="chartId" :class="{'scroll-box': chart.type == 'table'}"></div>
</div>`
})

Vue.component("metric-board", {
	template: `<div><slot></slot></div>`
})