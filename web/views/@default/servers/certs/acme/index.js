Tea.context(function () {
	this.viewCert = function (certId) {
		teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
			height: "28em",
			width: "48em"
		})
	}

	this.updateTask = function (taskId) {
		teaweb.popup("/servers/certs/acme/updateTaskPopup?taskId=" + taskId, {
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
})