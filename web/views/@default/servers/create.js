Tea.context(function () {
	this.serverType = "httpProxy";
	this.tlsProtocolName = ""
	this.origins = [];

	this.success = NotifySuccess("保存成功", "/servers");

	this.changeServerType = function () {
		this.origins = [];
		this.tlsProtocolName = "";
	};

	this.addOrigin = function () {
		teaweb.popup("/servers/addOriginPopup?serverType=" + this.serverType, {
			callback: function (resp) {
				this.origins.push(resp.data.origin);
			}
		});
	};

	this.removeOrigin = function (index) {
		this.origins.$remove(index);
	};
});