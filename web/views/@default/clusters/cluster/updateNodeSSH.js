Tea.context(function () {
	if (this.params.port <= 0) {
		this.params.port = 22
	}
})