Tea.context(function () {
	this.createNode = function () {
		teaweb.popup("/api/node/createPopup", {
			width: "50em",
			height: "30em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}
})