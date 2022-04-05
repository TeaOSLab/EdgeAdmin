Tea.context(function () {
	this.tab = "summary"

	this.switchTab = function (tab) {
		this.tab = tab
	}

	this.requestHeaders = []
	if (this.accessLog.header != null) {
		for (let k in this.accessLog.header) {
			let v = this.accessLog.header[k]
			if (typeof (v) != "object") {
				continue
			}
			this.requestHeaders.push({
				name: k,
				values: v.values
			})
		}
	}
	this.requestHeaders.sort(function (v1, v2) {
		return (v1.name < v2.name) ? -1 : 1
	})

	this.responseHeaders = []

	if (this.accessLog.sentHeader != null) {
		for (let k in this.accessLog.sentHeader) {
			let v = this.accessLog.sentHeader[k]
			if (typeof (v) != "object") {
				continue
			}
			this.responseHeaders.push({
				name: k,
				values: v.values
			})
		}
	}
	this.responseHeaders.sort(function (v1, v2) {
		return (v1.name < v2.name) ? -1 : 1
	})
})