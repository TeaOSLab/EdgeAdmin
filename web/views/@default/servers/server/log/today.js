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

	this.query = function (args) {
		// 初始化时页面尚未设置Vue变量，所以使用全局的变量获取
		let that = TEA.ACTION.data

		if (that.serverId == null) {
			that.serverId = 0
		}
		if (that.keyword == null) {
			that.keyword = ""
		}
		if (that.ip == null) {
			that.ip = ""
		}
		if (that.domain == null) {
			that.domain = ""
		}
		if (that.pageSize == null) {
			that.pageSize = ""
		}
		let query = 'serverId=' + that.serverId + '&keyword=' + encodeURIComponent(that.keyword) + '&ip=' + that.ip + '&domain=' + that.domain + '&pageSize=' + that.pageSize
		if (args != null && args.length > 0) {
			query += "&" + args
		}
		return query
	}

	this.allQuery = function () {
		if (this.query == null) {
			// 尚未初始化完成
			return
		}
		let query = this.query()
		if (this.hasError == 1) {
			query += "&hasError=1"
		}
		if (this.hasWAF == 1) {
			query += "&hasWAF=1"
		}
		return query
	}

	this.currentQuery = this.allQuery()
})