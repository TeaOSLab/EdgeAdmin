Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.changeMethods = function (config) {
		Tea.action("$")
			.form(this.$refs.authForm)
			.post()
		teaweb.successRefresh("保存成功")
	}
})