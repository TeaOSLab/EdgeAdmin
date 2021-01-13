Tea.context(function () {
	this.$delay(function () {
		this.load()
	})

	this.hasMore = false
	this.accessLogs = []
	this.isLoaded = false

	this.load = function () {
		this.$post("$")
			.params({
				serverId: this.serverId,
				requestId: this.requestId
			})
			.success(function (resp) {
				this.accessLogs = resp.data.accessLogs.concat(this.accessLogs)

                // 添加区域信息
                this.accessLogs.forEach(function (accessLog) {
                    if (typeof (resp.data.regions[accessLog.remoteAddr]) == "string") {
                        accessLog.region = resp.data.regions[accessLog.remoteAddr]
                    } else {
                        accessLog.region = ""
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
})