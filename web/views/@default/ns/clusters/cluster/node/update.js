Tea.context(function () {
	this.clusterId = 0;
	if (this.node.cluster != null && this.node.cluster.id > 0) {
		this.clusterId = this.node.cluster.id;
	}

	this.success = NotifySuccess("保存成功", "/ns/clusters/cluster/node?clusterId=" + this.clusterId + "&nodeId=" + this.node.id);

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
				methodName: this.node.login.grant.methodName,
				username: this.node.login.grant.username
			}
		}
	}
})