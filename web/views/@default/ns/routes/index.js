Tea.context(function () {
	this.$delay(function () {
		let that = this
		sortTable(function (ids) {
			that.$post(".sort")
				.params({
					routeIds: ids
				})
				.success(function () {
					teaweb.successToast("排序保存成功")
				})
		})
	})

	this.createRoute = function () {
		teaweb.popup("/ns/routes/createPopup", {
			width: "42em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateRoute = function (routeId) {
		teaweb.popup("/ns/routes/updatePopup?routeId=" + routeId, {
			width: "42em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteRoute = function (routeId) {
		let that = this
		teaweb.confirm("确定要删除此线路吗？", function () {
			that.$post(".delete")
				.params({
					routeId: routeId
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}
})