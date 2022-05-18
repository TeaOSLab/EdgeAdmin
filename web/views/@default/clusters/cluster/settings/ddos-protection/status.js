Tea.context(function () {
	this.isLoading = true
	this.results = []

	this.$delay(function () {
		this.reload()
	})

	this.reload = function () {
		this.isLoading = true
		this.$post("$")
			.params({ clusterId: this.clusterId })
			.success(function (resp) {
				this.results = resp.data.results
			})
			.done(function () {
				this.isLoading = false
			})
	}
})