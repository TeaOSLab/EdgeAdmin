Tea.context(function () {
	this.success = NotifyPopup
	this.type = "dnsla"
	this.typeDescription = "DNS.LA提供的DNS服务。"

	this.changeType = function () {
		let that = this
		let t = this.types.$find(function (k, v) {
			return v.code == that.type
		})
		if (t != null) {
			this.typeDescription = t.description
		} else {
			this.typeDescription = ""
		}
	}

	// DNSPod
	this.paramDNSPodAPIType = "tencentDNS"
})