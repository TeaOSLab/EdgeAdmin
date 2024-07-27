Tea.context(function () {
	this.success = function () {
		teaweb.success("保存成功", function () {
			teaweb.reload()
		})
	}

	this.changeAllowedProviders = function (event) {
		this.allowedProviders = event.providers
	}
})