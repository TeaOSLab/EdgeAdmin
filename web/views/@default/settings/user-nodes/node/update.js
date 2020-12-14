Tea.context(function () {
	this.success = NotifySuccess("保存成功", "/settings/userNodes/node?nodeId=" + this.node.id)

	this.hasHTTPS = this.node.listens.$any(function (k, v) {
		return v.protocol == "https"
	})
	this.changeListens = function (addrs) {
		this.hasHTTPS = addrs.$any(function (k, v) {
			return v.protocol == "https"
		})
	}
})