Tea.context(function () {
	this.success = NotifyPopup;

	this.isUpdating = false

	this.address = ""
	this.protocol = this.protocols[0].code

	// 初始化
	// from 用来标记是否为特殊的节点
	if (this.from.length == 0) {
		if (this.protocol == "http") {
			this.address = "80"
		} else if (this.protocol == "https") {
			this.address = "443"
		}
	}

	if (window.parent.UPDATING_ADDR != null) {
		this.isUpdating = true
		let addr = window.parent.UPDATING_ADDR
		this.protocol = addr.protocol
		if (addr.host.length == 0) {
			this.address = addr.portRange
		} else {
			this.address = addr.host.quoteIP() + ":" + addr.portRange
		}
	}

	this.changeProtocol = function () {
		if (this.from.length > 0) {
			return
		}
		switch (this.protocol) {
			case "http":
				this.address = "80"
				break
			case "https":
				this.address = "443"
				break
		}
	}

	this.addPort = function (port) {
		this.address = port
	}
});