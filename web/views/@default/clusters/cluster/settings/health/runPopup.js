Tea.context(function () {
	this.success = NotifyPopup

	this.isRequesting = true
	this.results = []
	this.countSuccess = 0
	this.countFail = 0
	this.errorString = ""

	this.$delay(function () {
		if (this.hasServers) {
			this.run()
		}
	})

	this.run = function () {
		this.isRequesting = true
		this.errorString = ""

		this.$post("$")
			.params({
				clusterId: this.clusterId
			})
			.timeout(60)
			.success(function (resp) {
				this.results = resp.data.results
				let that = this
				this.results.forEach(function (v) {
					v.costMs = Math.ceil(v.costMs)
					if (isNaN(v.costMs)) {
						v.costMs = 0
					}
					if (v.isOk) {
						that.countSuccess++
					} else {
						that.countFail++
					}
				})
			})
			.error(function () {
				this.errorString = "执行健康检查超时，请重试"
			})
			.done(function () {
				this.isRequesting = false
			})
	}
})