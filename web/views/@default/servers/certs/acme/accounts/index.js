Tea.context(function () {
	this.createAccount = function () {
		teaweb.popup(".createPopup", {
			height: "24em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateAccount = function (accountId) {
		teaweb.popup(".updatePopup?accountId=" + accountId, {
			height: "24em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteAccount = function (accountId) {
		teaweb.confirm("确定要删除此账号吗？", function () {
			this.$post(".delete")
				.params({
					accountId: accountId
				})
				.refresh()
		})
	}
})