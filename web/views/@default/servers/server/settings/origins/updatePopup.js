Tea.context(function () {
	this.changeAddr = function () {
		if (this.serverType == "httpProxy") {
			if (this.origin.addr.startsWith("http://")) {
				this.origin.protocol = "http"
			} else if (this.origin.addr.startsWith("https://")) {
				this.origin.protocol = "https"
			}
		}
	}
})