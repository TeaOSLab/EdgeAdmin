Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")
	this.checkResult = null

	this.$post(".status")
		.params({
			nodeId: this.nodeId
		})
		.success(function (resp) {
			let results = resp.data.results
			if (results.length > 0) {
				this.checkResult = results[0]
			}
		})
})