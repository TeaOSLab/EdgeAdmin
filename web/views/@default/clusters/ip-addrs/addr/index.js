Tea.context(function () {
	this.updateUp = function (addrId, isUp) {
		let status = isUp ? "在线" : "离线"
		teaweb.confirm("确定要手动将节点设置为" + status + "吗？", function () {
			this.$post(".up")
				.params({
					addrId: addrId,
					isUp: isUp ? 1 : 0
				})
				.refresh()
		})
	}

	this.restoreBackup = function (addrId) {
		teaweb.confirm("确定要恢复IP地址吗？", function () {
			this.$post(".restoreBackup")
				.params({
					addrId: addrId
				})
				.refresh()
		})
	}
})