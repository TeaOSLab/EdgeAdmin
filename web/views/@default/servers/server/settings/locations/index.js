Tea.context(function () {
	this.deleteLocation = function (locationId) {
		teaweb.confirm("确定要删除此路径规则吗？", function () {
			this.$post(".delete")
				.params({
					webId: this.webId,
					locationId: locationId
				})
				.refresh()
		})
	}
})