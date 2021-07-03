Tea.context(function () {
	this.createChart = function () {
		teaweb.popup(Tea.url(".createPopup?itemId=" + this.item.id), {
			callback: function () {
				teaweb.successRefresh("保存成功")
			},
			height: "27em"
		})
	}

	this.deleteChart = function (chartId) {
		let that = this
		teaweb.confirm("确定要删除这个图表吗？", function () {
			that.$post(".delete")
				.params({ chartId: chartId })
				.success(function () {
					teaweb.successRefresh("保存成功")
				})
		})
	}
})