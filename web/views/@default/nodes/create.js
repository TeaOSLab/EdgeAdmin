Tea.context(function () {
	this.grantId = 0;
	this.grant = null;

	this.success = NotifySuccess("保存成功", "/nodes");

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