Tea.context(function () {
	this.change = function (values) {
		this.$post("$")
			.params({
				webId: this.webId,
				hostRedirectsJSON: JSON.stringify(values)
			})
			.success(function () {
				teaweb.successToast("保存成功", null, function () {
					teaweb.reload()
				})
			})
	}
})