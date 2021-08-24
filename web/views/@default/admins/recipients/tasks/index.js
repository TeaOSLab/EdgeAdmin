Tea.context(function () {
	this.deleteTask = function (taskId) {
		let that = this
		teaweb.confirm("确定要删除这个发送任务吗？", function () {
			that.$post(".delete")
				.params({
					taskId: taskId
				})
				.success(function () {
					teaweb.successToast("删除成功", null, function () {
						teaweb.reload()
					})
				})
		})
	}
})