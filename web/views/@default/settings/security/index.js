Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.addDefaultClientIPHeaderNames = function (headerNames) {
		if (this.config.clientIPHeaderNames == null || this.config.clientIPHeaderNames.length == 0) {
			this.config.clientIPHeaderNames = headerNames
		} else {
			this.config.clientIPHeaderNames += " " + headerNames
		}
	}
})