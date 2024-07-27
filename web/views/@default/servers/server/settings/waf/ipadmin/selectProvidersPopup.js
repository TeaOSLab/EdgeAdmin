Tea.context(function () {
	this.isCheckingAll = false

	this.countSelectedProviders = this.providers.$count(function (k, provider) {
		return provider.isChecked
	})

	this.selectProvider = function (provider) {
		provider.isChecked = !provider.isChecked
		this.change()
	}

	this.deselectProvider = function (provider) {
		provider.isChecked = false
		this.change()
	}

	this.checkAll = function () {
		this.isCheckingAll = !this.isCheckingAll
		let that = this
		this.providers.forEach(function (provider) {
			provider.isChecked = that.isCheckingAll
		})

		this.change()
	}

	this.change = function () {

	}

	this.submit = function () {
		let selectedProviders = []
		this.providers.forEach(function (provider) {
			if (provider.isChecked) {
				selectedProviders.push(provider)
			}
		})
		NotifyPopup({
			code: 200,
			data: {
				selectedProviders: selectedProviders
			}
		})
	}
})