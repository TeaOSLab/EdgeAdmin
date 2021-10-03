Tea.context(function () {
	this.createUser = function () {
		teaweb.popup(Tea.url(".createPopup"), {
			height: "27em",
			width: "44em",
			callback: function () {
				teaweb.success("创建成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateUser = function (userId) {
		teaweb.popup("/servers/certs/acme/users/updatePopup?userId=" + userId, {
			height: "27em",
			width: "44em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteUser = function (userId) {
		let that = this
		teaweb.confirm("确定要删除此用户吗？", function () {
			that.$post(".delete")
				.params({
					userId: userId
				})
				.refresh()
		})
	}
})