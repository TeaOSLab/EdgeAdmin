Tea.context(function () {
	this.isCheckingAll = false

	this.selectProvince = function (province) {
		province.isChecked = !province.isChecked
	}

	this.deselectProvince = function (province) {
		province.isChecked = false
	}

	this.checkAll = function () {
		this.isCheckingAll = !this.isCheckingAll

		this.provinces.forEach(function (province) {
			province.isChecked = that.isCheckingAll
		})
	}

	this.success = function () {
		teaweb.successToast("保存成功")
	}
})