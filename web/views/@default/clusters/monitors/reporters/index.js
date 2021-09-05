Tea.context(function () {
	this.createReporter = function () {
		teaweb.popup(".createPopup", {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteReporter = function (reporterId) {
		teaweb.confirm("确定要删除此终端吗？", function () {
			this.$post(".reporter.delete")
				.params({
					reporterId: reporterId
				})
				.success(function () {
					teaweb.success("删除成功", function () {
						teaweb.reload()
					})
				})
		})
	}
})