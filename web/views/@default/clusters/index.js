Tea.context(function () {
	this.latestVisible = false

	this.showLatest = function () {
		this.latestVisible = !this.latestVisible
	}
})