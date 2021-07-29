Tea.context(function () {
	this.createPolicy = function () {
		teaweb.popup(Tea.url(".createPopup", {}), {
			height: "24em",
			callback: NotifyReloadSuccess("保存成功")
		})
	}

	this.deletePolicy = function (policyId) {
		teaweb.confirm("确定要删除这个日志策略吗？", function () {
			this.$post(".delete")
				.params({
					policyId: policyId
				})
				.success(function () {
					teaweb.successRefresh("保存成功")
				})
		})
	}
})