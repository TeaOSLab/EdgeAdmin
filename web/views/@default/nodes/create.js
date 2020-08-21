Tea.context(function () {
	this.success = NotifySuccess("保存成功", "/nodes");

	// IP地址相关
	this.ipAddresses = [];

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

	// 授权相关
	this.grantId = 0;
	this.grant = null;

	// 选择授权
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

	// 创建授权
	this.createGrant = function () {
		teaweb.popup("/nodes/grants/createPopup", {
			height: "31em",
			callback: function (resp) {
				this.grantId = resp.data.grant.id;
				if (this.grantId > 0) {
					this.grant = resp.data.grant;
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

	// 删除已选择授权
	this.removeGrant = function () {
		this.grant = null;
		this.grantId = 0;
	};
});