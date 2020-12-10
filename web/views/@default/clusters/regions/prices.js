Tea.context(function () {
	this.createItem = function () {
		teaweb.popup(Tea.url(".items.createPopup"), {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateItem = function (itemId) {
		teaweb.popup(Tea.url(".items.updatePopup", {itemId: itemId}), {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteItem = function (itemId) {
		let that = this
		teaweb.confirm("确定要删除此价格项吗？", function () {
			that.$post(".items.delete")
				.params({
					itemId: itemId
				})
				.refresh()
		})
	}

	this.updatePrice = function (regionId, itemId) {
		teaweb.popup(Tea.url(".updatePricePopup", {regionId: regionId, itemId: itemId}), {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}
})