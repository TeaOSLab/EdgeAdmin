Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")
	this.auditSuccess = NotifyReloadSuccess("保存成功")

	this.allServerNames = []
	let that = this
	this.serverNames.forEach(function (v) {
		if (v.subNames == null || v.subNames.length == 0) {
			that.allServerNames.push(v.name)
		} else {
			that.allServerNames.$pushAll(v.subNames)
		}
	})
	this.auditing = 1
})