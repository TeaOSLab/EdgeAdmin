Tea.context(function () {
	this.createThreshold = function () {
		teaweb.popup(Tea.url("/clusters/cluster/settings/thresholds/createPopup", {
			clusterId: this.clusterId,
			nodeId: this.nodeId
		}), {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateThreshold = function (thresholdId) {
		teaweb.popup(Tea.url("/clusters/cluster/settings/thresholds/updatePopup", {
			thresholdId: thresholdId
		}), {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteThreshold = function (thresholdId) {
		let that = this
		teaweb.confirm("确定要删除这个阈值吗？", function () {
			that.$post("/clusters/cluster/settings/thresholds/delete")
				.params({
					thresholdId: thresholdId
				})
				.success(function () {
					teaweb.success("删除成功", function () {
						teaweb.reload()
					})
				})
		})
	}
})