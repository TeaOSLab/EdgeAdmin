Tea.context(function () {
	this.clusterId = 0

	this.changeCluster = function(clusterId) {
		this.clusterId = clusterId
	}

	this.goNext = function () {
		if (this.clusterId > 0) {
			window.location = "/clusters/cluster/createNode?clusterId=" + this.clusterId
		}
	}
})