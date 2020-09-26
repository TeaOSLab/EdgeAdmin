Tea.context(function () {
	this.deleteServer = function (serverId) {
		teaweb.confirm("确定要删除当前服务吗？", function () {
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