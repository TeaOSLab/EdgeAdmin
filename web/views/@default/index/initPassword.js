Tea.context(function () {
	this.success = function () {
		teaweb.success("用户名和密码保存成功，现在去登录", function () {
			window.location = "/"
		})
	}
})