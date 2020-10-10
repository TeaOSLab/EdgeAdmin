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