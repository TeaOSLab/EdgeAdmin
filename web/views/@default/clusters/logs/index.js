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

	this.updateNodeRead = function (nodeId) {
		this.$post(".readLogs")
			.params({
				nodeId: nodeId
			})
			.success(function () {
				teaweb.reload()
			})
	}

	this.updateAllRead = function () {
		this.$post(".readAllLogs")
			.params({})
			.success(function () {
				teaweb.reload()
			})
	}

	this.changeCluster = function (clusterId) {
		this.clusterId = clusterId
	}

	this.fixLog = function (logId) {
		this.$post(".fix")
			.params({
				logIds: [logId]
			})
			.success(function () {
				teaweb.reload()
			})
	}

	this.fixPageLogs = function () {
		let logIds = this.logs.map(function (v) {
			return v.id
		})
		teaweb.confirm("确定已修复并消除当前页的问题？", function () {
			this.$post(".fix")
				.params({
					logIds: logIds
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}

	this.fixAllLogs = function () {
		teaweb.confirm("确定已修复并消除所有的问题？", function () {
			this.$post(".fixAll")
				.success(function () {
					teaweb.reload()
				})
		})
	}
})