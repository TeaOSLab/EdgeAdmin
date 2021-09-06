Tea.context(function () {
	this.syncDomain = function (domainId) {
		this.$post(".syncDomain")
			.params({
				domainId: domainId
			})
			.success(function () {
				teaweb.success("从服务商获取线路成功", function () {
					window.location.reload()
				})
			})
	}
})