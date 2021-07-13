Tea.context(function () {
	this.$delay(function () {
		this.sort()
	}, 1000)

	// 删除路由规则
	this.deleteLocation = function (locationId) {
		teaweb.confirm("确定要删除此路由规则吗？", function () {
			this.$post(".delete")
				.params({
					webId: this.webId,
					locationId: locationId
				})
				.refresh()
		})
	}

	// 排序
	this.sort = function () {
		if (this.locations.length == 0) {
			return
		}

		let box = this.$find("#sortable-table")[0]
		let that = this
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {
			},
			onUpdate: function (event) {
				let rows = box.querySelectorAll("tbody")
				let locationIds = []
				rows.forEach(function (row) {
					locationIds.push(parseInt(row.getAttribute("v-id")))
				})
				that.$post(".sort")
					.params({
						webId: that.webId,
						locationIds: locationIds
					})
					.success(function () {
						teaweb.success("保存成功")
					})
			}
		})
	}
})