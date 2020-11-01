Tea.context(function () {
	this.success = NotifyPopup;
	this.isUpdating = false
	this.mode = "single" // single|multiple
	this.serverName = {
		name: "",
		subNames: []
	}
	this.multipleServerNames = ""
	if (window.parent.UPDATING_SERVER_NAME != null) {
		this.isUpdating = true
		this.serverName = window.parent.UPDATING_SERVER_NAME
		if (this.serverName.subNames != null && this.serverName.subNames.length > 0) {
			this.mode = "multiple"
			this.multipleServerNames = this.serverName.subNames.join("\n")
		}
	}

	this.switchMode = function (mode) {
		this.mode = mode
		this.$delay(function () {
			if (mode == "single") {
				this.$refs.focus.focus()
			} else {
				this.$refs.serverNames.focus()
			}
		})
	}
});