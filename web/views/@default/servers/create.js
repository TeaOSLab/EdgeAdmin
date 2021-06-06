Tea.context(function () {
	this.serverType = "httpProxy"
	this.tlsProtocolName = ""
	this.origins = []
	this.defaultAddresses = []

	this.success = NotifySuccess("保存成功", "/servers");

	this.changeServerType = function () {
		this.origins = []
		this.tlsProtocolName = ""
		this.addDefaultAddresses()
	}

	// 初始化调用
	this.$delay(function () {
		this.changeServerType()
	})

	this.addOrigin = function () {
		teaweb.popup("/servers/addOriginPopup?serverType=" + this.serverType, {
			callback: function (resp) {
				this.origins.push(resp.data.origin);
			}
		})
	}

	this.removeOrigin = function (index) {
		this.origins.$remove(index)
	}

	this.addDefaultAddresses = function () {
		// 默认绑定的端口地址
		this.defaultAddresses = []
		if (this.serverType == "httpProxy" || this.serverType == "httpWeb") {
			this.defaultAddresses.push({
				"host": "",
				"portRange": "80",
				"protocol": "http"
			})
			this.defaultAddresses.push({
				"host": "",
				"portRange": "443",
				"protocol": "https"
			})
			this.tlsProtocolName = "https"
		}
	}
})