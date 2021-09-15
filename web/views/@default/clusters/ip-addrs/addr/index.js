Tea.context(function () {
	this.$delay(function () {
		this.loadChart()
	})

	this.updateUp = function (addrId, isUp) {
		let status = isUp ? "在线" : "离线"
		teaweb.confirm("确定要手动将节点设置为" + status + "吗？", function () {
			this.$post(".up")
				.params({
					addrId: addrId,
					isUp: isUp ? 1 : 0
				})
				.refresh()
		})
	}

	this.restoreBackup = function (addrId) {
		teaweb.confirm("确定要恢复IP地址吗？", function () {
			this.$post(".restoreBackup")
				.params({
					addrId: addrId
				})
				.refresh()
		})
	}

	this.loadChart = function () {
		if (this.results.length == 0) {
			return
		}

		let sumColor = "green"
		this.results.forEach(function (v) {
			switch (v.level) {
				case "good":
					v.color = "green"
					break
				case "normal":
					v.color = "blue"
					break
				case "bad":
					v.color = "orange"
					if (sumColor != "red") {
						sumColor = "orange"
					}
					break
				case "broken":
					v.color = "red"
					sumColor = "red"
					break
			}
		})


		let chartBox = document.getElementById("reports-chart-box")
		if (chartBox == null || chartBox.offsetHeight == 0) {
			let that = this
			setTimeout(function () {
				that.loadChart()
			})
			return
		}
		let chart = teaweb.initChart(chartBox)
		chart.setOption({
			radar: [
				{
					splitNumber: 4,
					indicator: this.results.map(function (result) {
						return {
							name: result.node.name,
							color: result.color,
							max: 5000
						}
					})
				}
			],
			series: [{
				name: '',
				type: 'radar',
				data: [
					{
						value: this.results.map(function (result) {
							return result.costMs
						})
					}
				],
				lineStyle: {
					width: "1",
					color: sumColor,
					opacity: 0.2
				},
				itemStyle: {
					opacity: 0
				},
				areaStyle: {
					color: sumColor,
					opacity: 0.2
				}
			}]
		})
	}
})