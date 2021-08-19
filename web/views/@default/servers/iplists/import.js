Tea.context(function () {
	this.success = function (resp) {
		let message = "成功导入" + resp.data.count + "个IP"
		if (resp.data.countIgnore > 0) {
			message += "，并忽略" + resp.data.countIgnore + "个格式错误的IP"
		}
		teaweb.success(message, function () {
			teaweb.reload()
		})
	}
})