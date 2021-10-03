Tea.context(function () {
	this.selectedProvider = null
	this.changeProvider = function () {
		if (this.providerCode.length == 0) {
			this.selectedProvider = null
			return
		}

		let that = this
		this.selectedProvider = this.providers.$find(function (k, v) {
			return v.code == that.providerCode
		})
	}

	this.changeProvider()
})