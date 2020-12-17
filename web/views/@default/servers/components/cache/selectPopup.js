Tea.context(function () {
	this.selectPolicy = function (cachePolicy) {
		NotifyPopup({
			code: 200,
			data: {
				cachePolicy: cachePolicy
			},
			message: ""
		})
	}
})