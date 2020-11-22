Tea.context(function () {
	this.success = NotifyPopup

	this.addresses = [];
	if (this.serverConfig != null && this.serverConfig.http != null && this.serverConfig.http.listen != null) {
		this.addresses = this.serverConfig.http.listen
	}
})