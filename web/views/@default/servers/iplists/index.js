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

	this.readAllItems = function () {
		this.$post(".readAll")
			.success(function () {
				teaweb.reload()
			})
	}
})