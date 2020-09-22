Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.updateOn = function (b) {
		teaweb.confirm(b ? "确定要启用反向代理服务吗？" : "确定要停用反向代理服务吗？", function () {
			this.$post(".updateOn")
				.params({
					"serverId": this.serverId,
					"isOn": b ? 1 : 0,
					"reverseProxyId": this.reverseProxyId
				})
				.success(function () {
					window.location.reload()
				})
		})
	}
})