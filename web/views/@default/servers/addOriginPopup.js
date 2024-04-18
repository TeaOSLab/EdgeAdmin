Tea.context(function () {
	this.addr = ""
	this.protocol = ""
	this.isOSS = false

	this.addrError = ""

	// 当前网站协议
	this.isHTTP = (this.serverType == "httpProxy" || this.serverType == "httpWeb")
	if (this.serverType == "httpProxy") {
		this.protocol = "http"
	} else if (this.serverType == "tcpProxy") {
		this.protocol = "tcp"
	} else if (this.serverType == "udpProxy") {
		this.protocol = "udp"
	}

	this.changeProtocol = function () {
		this.isOSS = this.protocol.startsWith("oss:")

		if (this.protocol == "http") {
			this.detectHTTPS()
		} else {
			this.adviceHTTPS = false
		}

		this.checkPort()
	}

	this.changeAddr = function () {
		this.adviceHTTPS = false

		if (this.serverType == "httpProxy") {
			if (this.addr.startsWith("http://")) {
				this.protocol = "http"
			} else if (this.addr.startsWith("https://")) {
				this.protocol = "https"
			}
		}

		this.checkPort()
	}

	this.checkPort = function () {
		this.addrError = ""

		// HTTP
		if (this.protocol == "http") {
			if (this.addr.endsWith(":443")) {
				this.addrError = "443通常是HTTPS协议端口，请确认源站协议选择是否正确。"
			} else if (this.addr.endsWith(":8443")) {
				this.addrError = "8443通常是HTTPS协议端口，请确认源站协议选择是否正确。"
			}
		}

		// HTTPS
		if (this.protocol == "https") {
			if (this.addr.endsWith(":80")) {
				this.addrError = "80通常是HTTP协议端口，请确认源站协议选择是否正确。"
			} else if (this.addr.endsWith(":8080")) {
				this.addrError = "8080通常是HTTP协议端口，请确认源站协议选择是否正确。"
			}
		}
	}

	this.adviceHTTPS = false

	var isDetectingHTTPS = false
	this.detectHTTPS = function () {
		if (isDetectingHTTPS) {
			return
		}
		isDetectingHTTPS = true

		this.adviceHTTPS = false
		if (this.protocol == "http") {
			this.$post("/servers/server/settings/origins/detectHTTPS")
				.params({
					addr: this.addr
				})
				.success(function (resp) {
					this.adviceHTTPS = resp.data.isOk
					if (resp.data.isOk) {
						this.addr = resp.data.addr
					}
				})
				.done(function () {
					isDetectingHTTPS = false
				})
		} else {
			isDetectingHTTPS = false
		}
	}

	this.switchToHTTPS = function () {
		this.adviceHTTPS = false
		this.protocol = "https"

		if (this.addr.endsWith(":80")) {
			this.addr = this.addr.substring(0, this.addr.length - (":80").length)
		}
	}
})