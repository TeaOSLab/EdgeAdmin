Tea.context(function () {
	this.success = NotifyPopup;

	this.address = window.parent.UPDATING_NODE_IP_ADDRESS
	if (this.address != null) {
		this.address.isUp = (this.address.isUp ? 1 : 0)

		// 专属集群
		if (this.address.clusters != null) {
			let selectedClusterIds = this.address.clusters.map(function (cluster) {
				return cluster.id
			})
			this.clusters.forEach(function (cluster) {
				cluster.isChecked = selectedClusterIds.$contains(cluster.id)
			})
		}
	}
})