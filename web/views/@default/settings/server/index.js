Tea.context(function () {
	this.updateHTTP = function () {
		teaweb.popup("/settings/server/updateHTTPPopup", {
			callback: function () {
				teaweb.success("保存成功", teaweb.reload)
			}
		})
	}

	this.updateHTTPS = function () {
		teaweb.popup("/settings/server/updateHTTPSPopup", {
			height: "26em",
			width:"50em",
			callback: function () {
				teaweb.success("保存成功", teaweb.reload)
			}
		})
	}
})