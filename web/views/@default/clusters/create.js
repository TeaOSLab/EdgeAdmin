Tea.context(function () {
	this.success = function (resp) {
		teaweb.success("集群创建成功", function () {
			window.location = "/clusters/cluster/nodes?clusterId=" + resp.data.clusterId
		})
	}

	this.domain = {id: 0, name: ""}
	this.changeDomain = function (domain) {
		this.domain.id = domain.id
		this.domain.name = domain.name
	}
})