Tea.context(function () {
	this.updateItem = function (itemId) {
		teaweb.popup(Tea.url(".updateIPPopup?listId=" + this.listId, {itemId: itemId}), {
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
					"listId": this.listId,
					"itemId": itemId
				})
				.refresh()
		})
	}

	/**
	 * 添加IP名单菜单
	 */
	this.createIP = function (type) {
		teaweb.popup("/servers/iplists/createIPPopup?listId=" + this.listId + '&type=' + type, {
			height: "30em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}
})