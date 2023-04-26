Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.run = function () {
		teaweb.confirm("确定要对当前集群下的所有节点进行健康检查吗？", function () {
			teaweb.popup("/clusters/cluster/settings/health/runPopup?clusterId=" + this.clusterId, {
				height: "30em"
			})
		})
	}
})