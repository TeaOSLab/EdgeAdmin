Tea.context(function () {
	this.success = function (resp) {
		teaweb.success("集群创建成功", function () {
			window.location = "/clusters/cluster/nodes?clusterId=" + resp.data.clusterId
		} )
	}
})