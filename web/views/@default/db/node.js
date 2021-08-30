Tea.context(function () {
	this.$delay(function () {
		this.loadStatus(this.node.id)
	})

	this.status = null
	this.loadStatus = function (nodeId) {
		this.$post(".status")
			.params({
				nodeId: nodeId
			})
			.success(function (resp) {
				this.status = resp.data.status
			})
	}
})