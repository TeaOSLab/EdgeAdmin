Tea.context(function () {
	this.typeDescription = ""

	let that = this
	this.types.forEach(function (v) {
		if (v.code == that.provider.type) {
			that.typeDescription = v.description
		}
	})

	// DNSPod
	if (this.provider.type == "dnspod" && this.provider.params != null && (this.provider.params.apiType == null || this.provider.params.apiType.length == 0)) {
		this.provider.params.apiType = "dnsPodToken"
	}
})