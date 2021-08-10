Tea.context(function () {
	this.deleteDomain = function (domainId) {
		let that = this
		teaweb.confirm("确定要删除此域名吗？", function () {
			that.$post("/ns/domains/delete")
				.params({
					domainId: domainId
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}
})