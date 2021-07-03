Tea.context(function () {
	this.createItem = function () {
		teaweb.popup(Tea.url(".createPopup?category=" + this.category), {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			},
			height: "26em",
			width: "44em"
		})
	}

	this.deleteItem = function (itemId) {
		let that = this
		teaweb.confirm("确定要删除此指标吗？", function () {
			that.$post(".delete")
				.params({
					itemId: itemId
				})
				.success(function () {
					teaweb.success("删除成功", function () {
						teaweb.reload()
					})
				})
		})
	}
})