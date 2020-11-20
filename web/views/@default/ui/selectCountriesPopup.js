Tea.context(function () {
	this.letterGroups = [
		"ABC", "DEF", "GHI", "JKL", "MNO", "PQR", "STU", "VWX", "YZ"
	];
	this.selectedGroup = "ABC"
	this.letterCountries = {}
	let that = this
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

	this.success = NotifyPopup
})