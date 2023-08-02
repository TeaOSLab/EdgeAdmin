Tea.context(function () {
	this.success = NotifySuccess("保存成功", "/servers/components/waf/policy?firewallPolicyId=" + this.firewallPolicyId)

	this.maxRequestBodySize = this.firewallPolicy.maxRequestBodySize
	this.maxRequestBodySizeFormat = teaweb.formatBytes(this.maxRequestBodySize)
	if (this.maxRequestBodySize == 0) {
		this.maxRequestBodySizeFormat = ""
	}

	this.changeMaxRequestBodySize = function (v) {
		if (v.toString().length == 0) {
			this.maxRequestBodySize = 0
			this.maxRequestBodySizeFormat = teaweb.formatBytes(this.maxRequestBodySize)
			return
		}
		let size = parseInt(v)
		if (!isNaN(size) && size >= 0) {
			this.maxRequestBodySize = size
			this.maxRequestBodySizeFormat = teaweb.formatBytes(this.maxRequestBodySize)
		}
	}
})