Tea.context(function () {
	this.verify = function () {
		teaweb.popup(".verifyPopup?userId=" + this.user.id, {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}
})