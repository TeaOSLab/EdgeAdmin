Tea.context(function () {
	this.success = NotifyPopup

	this.addresses = [];
	if (this.serverConfig != null && this.serverConfig.https != null && this.serverConfig.https.listen != null) {
		this.addresses = this.serverConfig.https.listen
	}
})