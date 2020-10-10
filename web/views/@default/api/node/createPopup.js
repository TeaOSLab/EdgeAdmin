Tea.context(function () {
	this.success = NotifyPopup

	this.hasHTTPS = false
	this.changeListens = function (addrs) {
		this.hasHTTPS = addrs.$any(function (k, v) {
			return v.protocol == "https"
		})
	}
})