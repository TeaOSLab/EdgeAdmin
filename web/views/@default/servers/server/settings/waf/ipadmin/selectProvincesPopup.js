Tea.context(function () {
	this.isCheckingAll = false

	this.countSelectedProvinces = this.provinces.$count(function (k, province) {
		return province.isChecked
	})

	this.selectProvince = function (province) {
		province.isChecked = !province.isChecked
		this.change()
	}

	this.deselectProvince = function (province) {
		province.isChecked = false
		this.change()
	}

	this.checkAll = function () {
		this.isCheckingAll = !this.isCheckingAll
		let that = this
		this.provinces.forEach(function (province) {
			province.isChecked = that.isCheckingAll
		})

		this.change()
	}

	this.change = function () {

	}

	this.submit = function () {
		let selectedProvinces = []
		this.provinces.forEach(function (province) {
			if (province.isChecked) {
				selectedProvinces.push(province)
			}
		})
		NotifyPopup({
			code: 200,
			data: {
				selectedProvinces: selectedProvinces
			}
		})
	}
})