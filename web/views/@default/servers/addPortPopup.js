Tea.context(function () {
	this.success = NotifyPopup;

	this.isUpdating = false

	this.address = ""
	this.protocol = this.protocols[0].code

	if (window.parent.UPDATING_ADDR != null) {
		this.isUpdating = true
		let addr = window.parent.UPDATING_ADDR
		this.protocol = addr.protocol
		if (addr.host.length == 0) {
			this.address = addr.portRange
		} else {
			this.address = addr.host + ":" + addr.portRange
		}
	}
});