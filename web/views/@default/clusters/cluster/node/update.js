Tea.context(function () {
	this.clusterId = 0
	if (this.node.cluster != null && this.node.cluster.id > 0) {
		this.clusterId = this.node.cluster.id
	}

	this.success = function () {
		let that = this
		teaweb.success("保存成功", function () {
			window.location = "/clusters/cluster/node/detail?clusterId=" + that.clusterId + "&nodeId=" + that.node.id
		})
	}

	// IP地址相关
	this.ipAddresses = this.node.ipAddresses

	this.changeClusters = function (info) {
		this.clusterId = info.clusterId
	}

	/**
	 * 集群相关
	 */
	this.showClustersBox = false

	this.updateClusters = function () {
		this.showClustersBox = !this.showClustersBox
	}

	/**
	 * 级别相关
	 */
	this.nodeLevel = this.node.level
	this.changeLevel = function (level) {
		this.nodeLevel = level
	}
})