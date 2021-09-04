Tea.context(function () {
	this.createReporter = function () {
		teaweb.popup(".createPopup", function () {
			teaweb.success("保存成功", function () {
				teaweb.reload()
			})
		})
	}
})