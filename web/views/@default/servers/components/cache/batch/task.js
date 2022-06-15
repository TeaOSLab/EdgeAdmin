Tea.context(function () {
	this.$delay(function () {
		this.reload()
	})

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

	this.reload = function () {
		this.$post("$")
			.params({
				taskId: this.task.id
			})
			.success(function (resp) {
				this.task = resp.data.task
			})
			.done(function () {
				this.$delay(function () {
					this.reload()
				}, 10000)
			})
	}
})