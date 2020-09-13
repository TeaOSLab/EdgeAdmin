Tea.context(function () {
	this.$delay(function () {
		this.reloadStatus(this.nodeId)
	})

	// 开始安装
	this.install = function () {
		this.$post("$")
			.params({
				nodeId: this.nodeId
			})
			.success(function () {

			})
	}

	// 设置节点安装状态
	this.updateNodeIsInstalled = function (isInstalled) {
		teaweb.confirm("确定要将当前节点修改为未安装状态？", function () {
			this.$post("/clusters/cluster/node/updateInstallStatus")
				.params({
					nodeId: this.nodeId,
					isInstalled: isInstalled ? 1 : 0
				})
				.refresh()
		})
	}

	// 刷新状态
	this.reloadStatus = function (nodeId) {
		this.$post("/clusters/cluster/node/status")
			.params({
				nodeId: nodeId
			})
			.success(function (resp) {
				this.installStatus = resp.data.installStatus
				this.node.isInstalled = resp.data.isInstalled
			})
			.done(function () {
				this.$delay(function () {
					this.reloadStatus(nodeId)
				}, 1000)
			});
	}
})