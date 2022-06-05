Tea.context(function () {
	this.deleteTask = function (taskId) {
		teaweb.confirm("确定要删除此任务吗？", function () {
			this.$post(".deleteTask")
				.params({
					taskId: taskId
				})
				.success(function () {
					window.location = Tea.url(".tasks")
				})
		})
	}

	this.resetTask = function (taskId) {
		teaweb.confirm("确定要重置任务状态吗？", function () {
			this.$post(".resetTask")
				.params({
					taskId: taskId
				})
				.refresh()
		})
	}
})