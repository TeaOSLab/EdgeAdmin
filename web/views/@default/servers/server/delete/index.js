Tea.context(function () {
	this.deleteServer = function (serverId) {
		teaweb.confirm("html:确定要删除当前网站吗？<br/>请慎重操作，删除后无法恢复！", function () {
			this.$post("$")
				.params({
					"serverId": serverId
				})
				.success(function () {
					teaweb.success("删除成功", function () {
						window.location = "/servers"
					})
				})
		})
	}
})