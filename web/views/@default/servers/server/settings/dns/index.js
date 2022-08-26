Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.regenerateCNAME = function () {
		let serverId = this.serverId
		teaweb.confirm("确定要重新生成此服务的CNAME吗？", function () {
			this.$post(".regenerateCNAME")
				.params({
					serverId: serverId
				})
				.refresh()
		})
	}

	this.updateCNAME = function () {
		teaweb.popup("/servers/server/settings/dns/updateCNAMEPopup?serverId=" + this.serverId, {
			callback: function () {
				teaweb.successRefresh("保存成功")
			}
		})
	}
})