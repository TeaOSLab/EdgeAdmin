Tea.context(function () {
	this.clusterId = 0;
	if (this.node.cluster != null && this.node.cluster.id > 0) {
		this.clusterId = this.node.cluster.id;
	}

	this.success = NotifySuccess("保存成功", "/nodes/node?nodeId=" + this.node.id);

	// IP地址相关
	this.ipAddresses = this.node.ipAddresses;

	// 添加IP地址
	this.addIPAddress = function () {
		teaweb.popup("/nodes/ipAddresses/createPopup", {
			callback: function (resp) {
				this.ipAddresses.push(resp.data.ipAddress);
			}
		})
	};

	// 修改地址
	this.updateIPAddress = function (index, address) {
		teaweb.popup("/nodes/ipAddresses/updatePopup?addressId=" + address.id, {
			callback: function (resp) {
				Vue.set(this.ipAddresses, index, resp.data.ipAddress);
			}
		})
	}

	// 删除IP地址
	this.removeIPAddress = function (index) {
		this.ipAddresses.$remove(index);
	};

	// 认证相关
	this.grantId = 0;
	this.grant = null;

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

	// 修改授权
	this.updateGrant = function () {
		if (this.grant == null) {
			window.location.reload();
			return;
		}
		teaweb.popup("/nodes/grants/updatePopup?grantId=" + this.grant.id, {
			height: "31em",
			callback: function (resp) {
				this.grant = resp.data.grant;
			}
		})
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