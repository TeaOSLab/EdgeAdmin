Tea.context(function () {
	this.serverType = "httpProxy"
	this.tlsProtocolName = ""
	this.origins = []
	this.defaultAddresses = []

	this.success = NotifySuccess("保存成功", "/servers");

	this.fail = function (resp) {
		if (resp.errors != null && resp.errors.length > 0) {
			let isFiltered = false

			let that = this
			resp.errors.forEach(function (err) {
				if (err.param == "emptyDomain") {
					isFiltered = true
					teaweb.warn(err.messages[0], function () {
						that.$refs.serverNameBox.addServerName()
					})
				} else if (err.param == "emptyOrigin") {
					isFiltered = true
					teaweb.warn(err.messages[0], function () {
						that.addOrigin()
					})
				}
			})

			if (isFiltered) {
				return
			}
		}
		Tea.failResponse(resp)
	}

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
			width: "45em",
			height: "27em",
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

	/**
	 * 用户相关
	 */
	this.userId = 0
	this.plans = []

	this.changeUserId = function (v) {
		this.userId = v

		if (this.userId == 0) {
			this.plans = []
			return
		}

		this.$post("/servers/users/plans")
			.params({
				userId: this.userId
			})
			.success(function (resp) {
				this.plans = resp.data.plans
			})
	}

	/**
	 * 证书相关
	 */
	this.findServerNames = function () {
		return this.$refs.serverNameBox.allServerNames()
	}
})