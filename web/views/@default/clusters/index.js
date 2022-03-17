Tea.context(function () {
	this.latestVisible = false

	this.showLatest = function () {
		this.latestVisible = !this.latestVisible
	}

	this.pin = function (clusterId, isPinned) {
		this.$post(".pin")
			.params({
				clusterId: clusterId,
				isPinned: isPinned
			})
			.success(function () {
				teaweb.reload()
			})
	}
})