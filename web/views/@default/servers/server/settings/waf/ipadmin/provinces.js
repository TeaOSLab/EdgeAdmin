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

	this.success = function () {
		teaweb.success("保存成功", function () {
			teaweb.reload()
		})
	}


	this.change = function () {
		this.countSelectedProvinces = this.provinces.$count(function (k, province) {
			return province.isChecked
		})
	}
})