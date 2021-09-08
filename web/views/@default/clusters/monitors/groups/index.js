Tea.context(function () {
	this.createGroup = function () {
		teaweb.popup(".createPopup", {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateGroup = function (groupId) {
		teaweb.popup(".group.updatePopup?groupId=" + groupId, {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteGroup = function (groupId) {
		teaweb.confirm("确定要删除此分组吗？", function () {
			this.$post(".group.delete")
				.params({
					groupId: groupId
				})
				.refresh()
		})
	}
})