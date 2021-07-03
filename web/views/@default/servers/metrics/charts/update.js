Tea.context(function () {
	this.type = this.chart.type
	this.typeDefinition = null

	this.$delay(function () {
		this.changeType()
	})

	this.changeType = function () {
		let that = this
		this.typeDefinition = this.types.$find(function (k, v) {
			return v.code == that.type
		})
	}

	this.success = NotifySuccess("保存成功", Tea.url(".chart", {chartId: this.chart.id}))
})