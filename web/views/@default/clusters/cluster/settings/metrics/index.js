Tea.context(function () {
	this.createItem = function () {
		teaweb.popup(Tea.url(".createPopup", {
			clusterId: this.clusterId,
			category: this.category
		}), {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			},
			onClose: function () {
				teaweb.reload()
			},
			width: "50em",
			height: "25em"
		})
	}

	this.deleteItem = function (itemId) {
		let that = this
		teaweb.confirm("确定要删除这个指标吗？", function () {
			that.$post(".delete")
				.params({
					clusterId: that.clusterId,
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