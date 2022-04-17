Tea.context(function () {
	this.tab = "summary"
	this.teaweb= teaweb

	this.switchTab = function (tab) {
		this.tab = tab
	}

	// 请求Header
	this.requestHeaders = []
	if (this.accessLog.header != null) {
		for (let k in this.accessLog.header) {
			let v = this.accessLog.header[k]
			if (typeof (v) != "object") {
				continue
			}
			this.requestHeaders.push({
				name: k,
				values: v.values,
				isGeneral: !k.startsWith("X-")
			})
		}
	}
	this.requestHeaders.sort(function (v1, v2) {
		return (v1.name < v2.name) ? -1 : 1
	})

	// 响应Header
	this.responseHeaders = []

	if (this.accessLog.sentHeader != null) {
		for (let k in this.accessLog.sentHeader) {
			let v = this.accessLog.sentHeader[k]
			if (typeof (v) != "object") {
				continue
			}
			this.responseHeaders.push({
				name: k,
				values: v.values,
				isGeneral: !k.startsWith("X-")
			})
		}
	}
	this.responseHeaders.sort(function (v1, v2) {
		return (v1.name < v2.name) ? -1 : 1
	})

	// Cookie
	this.cookies = []
	if (this.accessLog.cookie != null) {
		for (let k in this.accessLog.cookie) {
			let v = this.accessLog.cookie[k]
			if (typeof (v) != "string") {
				continue
			}
			this.cookies.push({
				name: k,
				value: v
			})
		}
	}
	this.cookies.sort(function (v1, v2) {
		if (v1.name.startsWith("_")) {
			if (v2.name.startsWith("_")) {
				return (v1.name < v2.name) ? -1 : 1
			}
			return -1
		}
		if (v2.name.startsWith("_")) {
			return 1
		}
		return (v1.name.toUpperCase() < v2.name.toUpperCase()) ? -1 : 1
	})
})