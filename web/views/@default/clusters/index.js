Tea.context(function () {
	this.deleteCluster = function (clusterId) {
		let that = this
		teaweb.confirm("确定要删除此集群吗？", function () {
			that.$post("/clusters/delete")
				.params({
					clusterId: clusterId
				})
				.refresh()
		})
	}
})