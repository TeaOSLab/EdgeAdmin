Tea.context(function () {
	this.success = NotifyPopup
	this.isCA = this.certConfig.isCA ? 1 : 0
	this.textMode = false

	this.switchTextMode = function () {
		this.textMode = !this.textMode
		if (this.textMode) {
			this.$delay(function () {
				this.$refs.certTextField.focus()
			})
		}
	}
})