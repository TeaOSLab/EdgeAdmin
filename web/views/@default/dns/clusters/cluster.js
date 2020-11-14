Tea.context(function () {
	this.updateClusterDNS = function (clusterId) {
		teaweb.popup("/dns/updateClusterPopup?clusterId=" + clusterId, {
			height: "22em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateNode = function (nodeId) {
		teaweb.popup("/dns/issues/updateNodePopup?nodeId=" + nodeId, {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}
})