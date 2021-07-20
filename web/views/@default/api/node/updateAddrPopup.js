Tea.context(function () {
	let addr = window.parent.UPDATING_ADDR
	this.protocol = addr.protocol
	this.addr = addr.host.quoteIP() + ":" + addr.portRange
})