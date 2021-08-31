Tea.context(function () {
	this.deleteAddr = function (addrId) {
		teaweb.confirm("确定要删除这个IP地址吗？", function () {
			this.$post(".addr.delete")
				.params({
					addrId: addrId
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}
})