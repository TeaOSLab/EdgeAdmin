Tea.context(function () {
	this.clusterId = 0
	if (this.node.cluster != null && this.node.cluster.id > 0) {
		this.clusterId = this.node.cluster.id
	}

	this.success = function () {
		let that = this
		teaweb.success("保存成功", function () {
			window.location = "/clusters/cluster/node/detail?clusterId=" + that.clusterId + "&nodeId=" + that.node.id
		})
	}

	// IP地址相关
	this.ipAddresses = this.node.ipAddresses

	// 认证相关
	this.grant = null

	this.sshHost = ""
	this.sshPort = ""
	this.loginId = 0
	if (this.node.login != null) {
		this.loginId = this.node.login.id

		if (this.node.login.params != null) {
			this.sshHost = this.node.login.params.host
			if (this.node.login.params.port > 0) {
				this.sshPort = this.node.login.params.port
			}
		}

		if (this.node.login.grant != null && typeof this.node.login.grant.id != "undefined") {
			this.grant = {
				id: this.node.login.grant.id,
				name: this.node.login.grant.name,
				method: this.node.login.grant.method,
				methodName: this.node.login.grant.methodName
			}
		}
	}

	this.changeClusters = function (info) {
		this.clusterId = info.clusterId
	}
})