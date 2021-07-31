Tea.context(function () {
	this.select = function (cluster) {
		NotifyPopup({
			code: 200,
			message: "",
			data: {
				cluster: cluster
			}
		})
	}
})