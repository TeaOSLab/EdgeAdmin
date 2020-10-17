Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.run = function () {
		teaweb.confirm("确定要对当前集群下的所有节点进行健康检查吗？", function () {
			teaweb.popup("/clusters/cluster/settings/healthRun?clusterId=" + this.clusterId, {

			})
		})
	}
})