Tea.context(function () {
	// store information to local
	localStorage.setItem("sid", this.localSid)
	localStorage.setItem("ip", this.ip)

	if (this.errorMsg.length == 0) {
		if (this.redirect.length > 0) {
			window.location = this.redirect
		} else {
			window.location = "/dashboard"
		}
	}
})