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

	this.round = function (f) {
		return Math.round(f * 100) / 100
	}

	this.moreOpsVisible = false
	this.showMoreOps = function () {
		this.moreOpsVisible = !this.moreOpsVisible
	}

	this.uninstallNode = function () {
		let that = this
		teaweb.confirm("html:确定要卸载当前节点吗？<br/>此操作将会删除节点上除缓存以外的相关文件。", function () {
			that.$post("/clusters/cluster/node/uninstall")
				.params({
					nodeId: that.node.id
				})
				.success(function () {
					teaweb.success("执行成功", function () {
						teaweb.reload()
					})
				})
		})
	}
})