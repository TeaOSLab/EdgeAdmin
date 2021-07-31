Tea.context(function () {
	this.updateCluster = function (clusterId) {
		teaweb.popup("/dns/updateClusterPopup?clusterId=" + clusterId, {
			height: "25em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateNode = function (clusterId, nodeId) {
		teaweb.popup("/dns/issues/updateNodePopup?clusterId=" + clusterId + "&nodeId=" + nodeId, {
			height: "26em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.isSyncing = false
	this.syncCluster = function (clusterId) {
		let that = this
		teaweb.confirm("确定要执行数据同步吗？", function () {
			that.isSyncing = true
			that.$post(".sync")
				.params({clusterId: clusterId})
				.done(function () {
					that.isSyncing = false
					that.dnsHasChanges = false
				})
				.success(function () {
					teaweb.success("同步成功", function () {
						teaweb.reload()
					})
				})
		})
	}
})