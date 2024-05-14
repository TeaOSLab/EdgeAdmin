Tea.context(function () {
	this.windowWidth = window.innerWidth
	this.columnWidth1 = 1000

	this.createUser = function () {
		teaweb.popup(Tea.url(".createPopup"), {
			height: "30em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteUser = function (userId) {
		let that = this
		teaweb.confirm("确定要删除这个用户吗？", function () {
			that.$post(".delete")
				.params({
					userId: userId
				})
				.refresh()
		})
	}
})