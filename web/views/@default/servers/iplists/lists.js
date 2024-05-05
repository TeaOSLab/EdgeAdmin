Tea.context(function () {
	this.createList = function () {
		teaweb.popup(Tea.url(".createPopup", {type: this.type}), {
			height: "32em",
			callback: function (resp) {
				teaweb.success("保存成功", function () {
					window.location = "/servers/iplists/lists?type=" + resp.data.list.type
				})
			}
		})
	}

	this.deleteList = function (listId) {
		let that = this
		teaweb.confirm("确定要删除此IP名单吗？", function () {
			that.$post(".delete")
				.params({
					listId: listId
				})
				.success(function () {
					teaweb.success("删除成功", function () {
						teaweb.reload()
					})
				})
		})
	}
})