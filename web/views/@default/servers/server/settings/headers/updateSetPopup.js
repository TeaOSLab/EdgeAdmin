Tea.context(function () {
	this.shouldReplace = this.headerConfig.shouldReplace
	this.statusList = []
	if (this.headerConfig.status != null && this.headerConfig.status.codes != null) {
		this.statusList = this.headerConfig.status.codes
	}

	this.selectHeader = function (headerName) {
		this.headerConfig.name = headerName
	}
})