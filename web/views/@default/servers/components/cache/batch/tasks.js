Tea.context(function () {
	this.deleteTask = function (taskId) {
		teaweb.confirm("确定要删除此任务吗？", function () {
			this.$post(".deleteTask")
				.params({
					taskId: taskId
				})
				.refresh()
		})
	}
})