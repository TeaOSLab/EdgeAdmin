Tea.context(function () {
	this.$delay(function () {
		let that = this
		sortTable(function () {
			let groupIds = []
			document.querySelectorAll("*[data-group-id]").forEach(function (element) {
				groupIds.push(element.getAttribute("data-group-id"))
			})
			that.$post("/servers/components/groups/sort")
				.params({
					groupIds: groupIds
				})
				.success(function () {
					teaweb.successToast("保存成功")
				})
		})
	})

	this.createGroup = function () {
		teaweb.popup("/servers/components/groups/createPopup", {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateGroup = function (groupId) {
		teaweb.popup("/servers/components/groups/updatePopup?groupId=" + groupId, {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteGroup = function (groupId) {
		let that = this
		teaweb.confirm("确定要删除这个分组吗？", function () {
			that.$post("/servers/components/groups/delete")
				.params({
					groupId: groupId
				})
				.success(function () {
					teaweb.success("删除成功", function () {
						teaweb.reload()
					})
				})
		})
	}
})