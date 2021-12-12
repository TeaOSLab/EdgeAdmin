Tea.context(function () {
	this.success = function (resp) {
		window.location = "/servers/components/waf/exportDownload?key=" + resp.data.key + "&policyId=" + resp.data.id
	}
})