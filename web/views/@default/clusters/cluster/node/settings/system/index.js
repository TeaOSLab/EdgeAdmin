Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.changeBypassMobile = function(bypassMobile) {
		if (bypassMobile.target.checked) {
			this.node.bypassMobile = 17;
		} else {
			this.node.bypassMobile = 0;
		}
	}
})