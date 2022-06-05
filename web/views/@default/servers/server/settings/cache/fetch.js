Tea.context(function () {
	this.isRequesting = false
	this.isOk = false
	this.message = ""
	this.failKeys = []

	this.before = function () {
		this.isRequesting = true
		this.isOk = false
		this.message = ""
		this.failKeys = []
	}

	this.success = function (resp) {
		this.isOk = true

		let f = NotifyReloadSuccess("任务提交成功")
		f()
	}

	this.fail = function (resp) {
		this.message = resp.message

		if (resp.data.failKeys != null) {
			this.failKeys = resp.data.failKeys
		}
	}

	this.done = function () {
		this.isRequesting = false
	}
});