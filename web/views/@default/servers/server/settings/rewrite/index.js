Tea.context(function () {
	// 创建重写规则
	this.createRewriteRule = function () {
		teaweb.popup("/servers/server/settings/rewrite/createPopup?webId=" + this.webId, {
			height: "26em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}
})