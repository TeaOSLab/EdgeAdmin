Tea.context(function () {
	this.updateItem = function (itemId) {
		teaweb.popup(Tea.url(".updateIPPopup", {itemId: itemId}), {
			height: "30em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteItem = function (itemId) {
		let that = this
		teaweb.confirm("确定要删除这个IP吗？", function () {
			that.$post(".deleteIP")
				.params({
					"itemId": itemId
				})
				.refresh()
		})
	}

	/**
	 * 添加IP名单菜单
	 */
	this.createIP = function () {
		teaweb.popup(Tea.url(".createIPPopup", {listId: this.list.id}), {
			height: "30em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}
})