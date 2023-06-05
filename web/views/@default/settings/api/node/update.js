Tea.context(function () {
	this.success = NotifySuccess("保存成功", "/settings/api/node?nodeId=" + this.node.id)

	this.hasHTTPS = this.node.listens.$any(function (k, v) {
		return v.protocol == "https"
	}) || (this.node.restIsOn && this.node.restListens.$any(function (k, v) {
		return v.protocol == "https"
	}))
	this.grpcAddrs = []
	this.restAddrs = []

	this.changeListens = function (addrs) {
		this.grpcAddrs = addrs

		this.hasHTTPS = this.grpcAddrs.$any(function (k, v) {
			return v.protocol == "https"
		}) || (this.node.restIsOn && this.restAddrs.$any(function (k, v) {
			return v.protocol == "https"
		}))
	}

	this.changeRestListens = function (addrs) {
		this.restAddrs = addrs

		this.hasHTTPS = this.grpcAddrs.$any(function (k, v) {
			return v.protocol == "https"
		}) || (this.node.restIsOn && this.restAddrs.$any(function (k, v) {
			return v.protocol == "https"
		}))
	}
})