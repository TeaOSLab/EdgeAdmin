Tea.context(function () {
	this.success = NotifySuccess("保存成功", "/api")

	this.hasHTTPS = false
	this.changeListens = function (addrs) {
		this.hasHTTPS = addrs.$any(function (k, v) {
			return v.protocol == "https"
		})
	}
})