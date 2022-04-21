Tea.context(function () {
    let that = this
    this.accessLogs.forEach(function (accessLog) {
        if (typeof (that.regions[accessLog.remoteAddr]) == "string") {
            accessLog.region = that.regions[accessLog.remoteAddr]
        } else {
            accessLog.region = ""
        }
		if (accessLog.firewallRuleSetId > 0 && typeof (that.wafInfos[accessLog.firewallRuleSetId]) == "object") {
			accessLog.wafInfo = that.wafInfos[accessLog.firewallRuleSetId]
		} else {
			accessLog.wafInfo = null
		}
    })
})