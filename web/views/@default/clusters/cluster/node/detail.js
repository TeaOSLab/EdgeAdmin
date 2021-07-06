Tea.context(function () {
	this.isStarting = false
	this.startNode = function () {
		this.isStarting = true
		this.$post("/clusters/cluster/node/start")
			.params({
				nodeId: this.node.id
			})
			.success(function () {
				teaweb.success("启动成功", function () {
					teaweb.reload()
				})
			})
			.done(function () {
				this.isStarting = false
			})
	}

	this.isStopping = false
	this.stopNode = function () {
		this.isStopping = true
		this.$post("/clusters/cluster/node/stop")
			.params({
				nodeId: this.node.id
			})
			.success(function () {
				teaweb.success("执行成功", function () {
					teaweb.reload()
				})
			})
			.done(function () {
				this.isStopping = false
			})
	}
})