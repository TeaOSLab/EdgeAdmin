Tea.context(function () {
	this.updateRead = function (logId) {
		this.$post(".readEvents")
			.params({
				logIds: [logId]
			})
			.success(function () {
				teaweb.reload()
			})
	}

	this.updatePageRead = function () {
		let logIds = this.logs.map(function (v) {
			return v.id
		})
		this.$post(".readEvents")
			.params({
				logIds: logIds
			})
			.success(function () {
				teaweb.reload()
			})
	}
})