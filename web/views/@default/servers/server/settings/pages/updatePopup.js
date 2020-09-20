Tea.context(function () {
	this.success = NotifyPopup

	this.newStatus = ""
	if (this.pageConfig.newStatus > 0) {
		this.newStatus = this.pageConfig.newStatus
	}
})