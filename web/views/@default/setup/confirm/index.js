Tea.context(function () {
	this.isRequesting = false

	this.success = function () {
		teaweb.success("完成确认，现在跳转到首页", function () {
			window.location = "/"
		})
	}

	this.before = function () {
		this.isRequesting = true
	}

	this.done = function () {
		this.isRequesting = false
	}
})