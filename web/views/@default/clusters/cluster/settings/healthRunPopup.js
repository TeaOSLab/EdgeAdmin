Tea.context(function () {
	this.success = NotifyPopup

	this.isRequesting = true
	this.results = []
	this.countSuccess = 0
	this.countFail = 0

	this.$delay(function () {
		this.run()
	})

	this.run = function () {
		this.isRequesting = true

		this.$post("$")
			.params({
				clusterId: this.clusterId
			})
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
			.done(function () {
				this.isRequesting = false
			})
	}
})