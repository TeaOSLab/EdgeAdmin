Tea.context(function () {
	this.addr = ""
	this.protocol = ""

	if (this.isHTTP) {
		this.protocol = "http"
	}

	this.changeAddr = function () {
		if (this.serverType == "httpProxy") {
			if (this.addr.startsWith("http://")) {
				this.protocol = "http"
			} else if (this.addr.startsWith("https://")) {
				this.protocol = "https"
			}
		}
	}
})