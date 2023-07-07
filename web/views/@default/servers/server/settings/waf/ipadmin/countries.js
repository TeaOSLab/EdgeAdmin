Tea.context(function () {
	this.success = function () {
		teaweb.success("保存成功", function () {
			teaweb.reload()
		})
	}

	this.changeAllowedCountries = function (event) {
		this.allowedCountries = event.countries
	}
})