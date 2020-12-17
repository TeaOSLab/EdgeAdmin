Tea.context(function () {
	this.selectPolicy = function (firewallPolicy) {
		NotifyPopup({
			code: 200,
			data: {
				firewallPolicy: firewallPolicy
			},
			message: ""
		})
	}
})