Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.change = function (values) {
		this.$post("$")
			.params({
				webId: this.webId,
				hostRedirectsJSON: JSON.stringify(values)
			})
			.success(function () {
				NotifyReloadSuccess("保存成功")()
			})
	}
})