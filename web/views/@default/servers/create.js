Tea.context(function () {
	this.serverType = "httpProxy";
	this.addresses = [];
	this.tlsProtocolName = ""

	this.serverNames = [];
	this.origins = [];


	this.success = NotifySuccess("保存成功", "/servers");

	this.changeServerType = function () {
		this.addresses = [];
		this.origins = [];
		this.tlsProtocolName = "";
	};

	this.addPort = function () {
		teaweb.popup("/servers/addPortPopup?serverType=" + this.serverType, {
			callback: function (resp) {
				var addr = resp.data.address;
				this.addresses.push(addr);
				if (["https", "https4", "https6"].$contains(addr.protocol)) {
					this.tlsProtocolName = "HTTPS";
				} else if (["tls", "tls4", "tls6"].$contains(addr.protocol)) {
					this.tlsProtocolName = "TLS";
				}
			}
		})
	};

	this.removeAddr = function (index) {
		this.addresses.$remove(index);
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