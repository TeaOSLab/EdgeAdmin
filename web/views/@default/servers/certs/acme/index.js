Tea.context(function () {
	this.viewCert = function (certId) {
		teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
			height: "28em",
			width: "48em"
		})
	}

	this.updateTask = function (taskId) {
		teaweb.popup("/servers/certs/acme/updateTaskPopup?taskId=" + taskId, {
			width: "45em",
			height: "26em",
			callback: function () {
				teaweb.success("保存成功，如果证书域名发生了改变，请重新执行生成新证书", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteTask = function (taskId) {
		let that = this
		teaweb.confirm("确定要删除此任务吗？", function () {
			that.$post("/servers/certs/acme/deleteTask")
				.params({
					taskId: taskId
				})
				.refresh()
		})
	}


	this.isRunning = false
	this.runningIndex = -1

	this.runTask = function (index, task) {
		let that = this

		teaweb.confirm("html:确定要立即执行此任务吗？<br/>将会重新发起证书申请。", function () {
			that.isRunning = true
			that.runningIndex = index

			that.$post(".run")
				.timeout(300)
				.params({
					taskId: task.id
				})
				.success(function (resp) {
					teaweb.success("任务执行成功", function () {
						teaweb.reload()
					})
				})
				.done(function () {
					that.isRunning = false
					that.runningIndex = -1
				})
		})
	}

	this.showError = function (err) {
		teaweb.popupTip("任务执行失败：" + err)
	}
})