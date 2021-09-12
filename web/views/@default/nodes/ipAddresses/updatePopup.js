Tea.context(function () {
	this.success = NotifyPopup;

	this.address = window.parent.UPDATING_NODE_IP_ADDRESS
	if (this.address != null) {
		this.address.isUp = (this.address.isUp ? 1 : 0)
	}
})