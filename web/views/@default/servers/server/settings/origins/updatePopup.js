Tea.context(function () {
	this.addrError = ""
	this.isOSS = this.origin != null && this.origin.protocol != null && this.origin.protocol.startsWith("oss:")

	// 预先设置oss选项
	if (!this.isOSS) {
		this.origin.oss = {
			type: "",
			options: {}
		}
	}

	this.$delay(function () {
		this.checkPort()
	})

	this.changeProtocol = function () {
		this.isOSS = this.origin.protocol.startsWith("oss:")

		this.checkPort()
	}

	this.changeAddr = function () {
		if (this.serverType == "httpProxy") {
			if (this.origin.addr.startsWith("http://")) {
				this.origin.protocol = "http"
			} else if (this.origin.addr.startsWith("https://")) {
				this.origin.protocol = "https"
			}
		}

		this.checkPort()
	}

	this.checkPort = function () {
		this.addrError = ""

		// HTTP
		if (this.origin.protocol == "http") {
			if (this.origin.addr.endsWith(":443")) {
				this.addrError = "443通常是HTTPS协议端口，请确认源站协议选择是否正确。"
			}
		}

		// HTTPS
		if (this.origin.protocol == "https") {
			if (this.origin.addr.endsWith(":80")) {
				this.addrError = "80通常是HTTP协议端口，请确认源站协议选择是否正确。"
			}
		}
	}
})