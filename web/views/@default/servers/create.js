Tea.context(function () {
	this.serverType = "httpProxy";
	this.tlsProtocolName = ""

	this.serverNames = [];
	this.origins = [];


	this.success = NotifySuccess("保存成功", "/servers");

	this.changeServerType = function () {
		this.origins = [];
		this.tlsProtocolName = "";
	};

	this.addServerName = function () {
		teaweb.popup("/servers/addServerNamePopup", {
			callback: function (resp) {
				var serverName = resp.data.serverName;
				this.serverNames.push(serverName);
			}
		});
	};

	this.removeServerName = function (index) {
		this.serverNames.$remove(index);
	};

	this.addOrigin = function () {
		teaweb.popup("/servers/addOriginPopup?serverType=" + this.serverType, {
			callback: function (resp){
				this.origins.push(resp.data.origin);
			}
		});
	};

	this.removeOrigin = function (index) {
		this.origins.$remove(index);
	};
});