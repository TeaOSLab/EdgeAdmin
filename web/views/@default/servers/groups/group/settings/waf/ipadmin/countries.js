Tea.context(function () {
	this.letterGroups = [
		"ABC", "DEF", "GHI", "JKL", "MNO", "PQR", "STU", "VWX", "YZ"
	];
	this.selectedGroup = "ABC"
	this.letterCountries = {}
	let that = this
	this.countSelectedCountries = this.countries.$count(function (k, country) {
		return country.isChecked
	})
	this.countries.forEach(function (country) {
		if (typeof (that.letterCountries[country.letter]) == "undefined") {
			that.letterCountries[country.letter] = []
		}
		that.letterCountries[country.letter].push(country)
	})
	this.isCheckingAll = false

	this.selectGroup = function (group) {
		this.selectedGroup = group
	}

	this.selectCountry = function (country) {
		country.isChecked = !country.isChecked
		this.change()
	}

	this.deselectCountry = function (country) {
		country.isChecked = false
		this.change()
	}

	this.checkAll = function () {
		this.isCheckingAll = !this.isCheckingAll

		this.countries.forEach(function (country) {
			country.isChecked = that.isCheckingAll
		})

		this.change()
	}

	this.success = function () {
		teaweb.success("保存成功", function () {
			teaweb.reload()
		})
	}

	this.change = function () {
		this.countSelectedCountries = this.countries.$count(function (k, country) {
			return country.isChecked
		})
	}
})