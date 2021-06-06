Tea.context(function () {
	this.success = NotifyPopup;

	this.isUpdating = false

	this.address = ""
	this.protocol = this.protocols[0].code

	// 初始化
	if (this.protocol == "http") {
		this.address = "80"
	} else if (this.protocol == "https") {
		this.address = "443"
	}

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

	this.changeProtocol = function () {
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