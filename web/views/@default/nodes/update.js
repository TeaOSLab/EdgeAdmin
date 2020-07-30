Tea.context(function () {
	this.grantId = 0;
	this.grant = null;


	this.clusterId = 0;
	if (this.node.cluster != null && this.node.cluster.id > 0) {
		this.clusterId = this.node.cluster.id;
	}

	this.sshHost = "";
	this.sshPort = "";
	this.loginId = 0;
	if (this.node.login != null) {
		this.loginId = this.node.login.id;

		if (this.node.login.params != null) {
			this.sshHost = this.node.login.params.host;
			this.sshPort = this.node.login.params.port;
		}

		if (this.node.login.grant != null) {
			this.grantId = this.node.login.grant.id;
			this.grant = {
				id: this.node.login.grant.id,
				name: this.node.login.grant.name,
				method: this.node.login.grant.method,
				methodName: this.node.login.grant.methodName
			};
		}
	}

	this.success = NotifySuccess("保存成功", "/nodes/node?nodeId=" + this.node.id);

	this.selectGrant = function () {
		var that = this;
		teaweb.popup("/nodes/grants/selectPopup", {
			callback: function (resp) {
				that.grantId = resp.data.grant.id;
				if (that.grantId > 0) {
					that.grant = resp.data.grant;
				}
			}
		});
	};

	this.createGrant = function () {
		var that = this;
		teaweb.popup("/nodes/grants/createPopup", {
			height: "31em",
			callback: function (resp) {
				that.grantId = resp.data.grant.id;
				if (that.grantId > 0) {
					that.grant = resp.data.grant;
				}
			}
		});
	};

	this.removeGrant = function () {
		this.grant = null;
		this.grantId = 0;
	};
});