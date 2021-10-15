Tea.context(function () {
	this.$delay(function () {
		teaweb.datepicker("day-from-picker")
		teaweb.datepicker("day-to-picker")
	})

	this.updateRead = function (logId) {
		this.$post(".readLogs")
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
		this.$post(".readLogs")
			.params({
				logIds: logIds
			})
			.success(function () {
				teaweb.reload()
			})
	}
})