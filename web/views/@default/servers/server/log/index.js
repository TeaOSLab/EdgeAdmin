Tea.context(function () {
	this.$delay(function () {
		this.load()
	})

	this.hasMore = false
	this.accessLogs = []
	this.isLoaded = false

	this.load = function () {
		// 如果有弹窗时，暂时不更新
		if (teaweb.hasPopup()) {
			this.$delay(function () {
				this.load()
			}, 5000)
			return
		}

		this.$post("$")
			.params({
				serverId: this.serverId,
				requestId: this.requestId,
				keyword: this.keyword,
				ip: this.ip,
				domain: this.domain,
				clusterId: this.clusterId,
				nodeId: this.nodeId
			})
			.success(function (resp) {
				this.accessLogs = resp.data.accessLogs.concat(this.accessLogs)

				// 添加区域信息
				let that = this
				this.accessLogs.forEach(function (accessLog) {
					that.formatTime(accessLog)
					if (typeof (resp.data.regions[accessLog.remoteAddr]) == "string") {
						accessLog.region = resp.data.regions[accessLog.remoteAddr]
					} else {
						accessLog.region = ""
					}
					if (accessLog.firewallRuleSetId > 0 && typeof (resp.data.wafInfos[accessLog.firewallRuleSetId]) == "object") {
						accessLog.wafInfo = resp.data.wafInfos[accessLog.firewallRuleSetId]
					} else {
						accessLog.wafInfo = null
					}
				})

				let max = 100
				if (this.accessLogs.length > max) {
					this.accessLogs = this.accessLogs.slice(0, max)
				}
				this.hasMore = resp.data.hasMore
				this.requestId = resp.data.requestId
			})
			.done(function () {
				if (!this.isLoaded) {
					this.$delay(function () {
						this.isLoaded = true
					})
				}

				// 自动刷新
				this.$delay(function () {
					this.load()
				}, 5000)
			})
	}


	this.formatTime = function (accessLog) {
		let elapsedSeconds = Math.ceil(new Date().getTime() / 1000) - accessLog.timestamp
		if (elapsedSeconds >= 0) {
			if (elapsedSeconds < 60) {
				accessLog.humanTime = elapsedSeconds + "秒前"
			} else if (elapsedSeconds < 3600) {
				accessLog.humanTime = Math.ceil(elapsedSeconds / 60) + "分钟前"
			} else if (elapsedSeconds < 3600 * 24) {
				accessLog.humanTime = Math.ceil(elapsedSeconds / 3600) + "小时前"
			}
		}
	}
})