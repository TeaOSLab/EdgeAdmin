Tea.context(function () {
	let that = this
	this.accessLogs.forEach(function (accessLog) {
		// 区域
		if (typeof (that.regions[accessLog.remoteAddr]) == "string") {
			accessLog.region = that.regions[accessLog.remoteAddr]
		} else {
			accessLog.region = ""
		}

		// 节点
		if (typeof (that.nodes[accessLog.nsNodeId]) != "undefined") {
			accessLog["node"] = that.nodes[accessLog.nsNodeId]
		} else {
			accessLog["node"] = null
		}

		// 域名
		if (typeof (that.domains[accessLog.nsDomainId]) != "undefined") {
			accessLog["domain"] = that.domains[accessLog.nsDomainId]
		} else {
			accessLog["domain"] = null
		}
	})

	this.$delay(function () {
		let that = this
		teaweb.datepicker("day-input", function (v) {
			that.day = v
		})
	})
})