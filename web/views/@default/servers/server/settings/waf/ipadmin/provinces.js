Tea.context(function () {
	this.success = function () {
		teaweb.success("保存成功", function () {
			teaweb.reload()
		})
	}

	this.changeAllowedProvinces = function (event) {
		this.allowedProvinces = event.provinces
	}
})