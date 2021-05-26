Tea.context(function () {
	this.clusterId = 0;
	if (this.node.cluster != null && this.node.cluster.id > 0) {
		this.clusterId = this.node.cluster.id;
	}

	this.success = NotifySuccess("保存成功", "/ns/clusters/cluster/node?clusterId=" + this.clusterId + "&nodeId=" + this.node.id);
});