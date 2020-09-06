Tea.context(function () {
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
})