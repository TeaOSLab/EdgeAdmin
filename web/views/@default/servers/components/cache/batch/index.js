Tea.context(function () {
	this.isRequesting = false
	this.isOk = false
	this.message = ""
	this.failKeys = []

	this.$delay(function () {
		this.$refs.keysBox.focus()
		this.$watch("keyType", function () {
			this.$refs.keysBox.focus()
		})
	})

	this.before = function () {
		this.isRequesting = true
		this.isOk = false
		this.message = ""
		this.failKeys = []
	}

	this.success = function () {
		this.isOk = true
		let that = this
		teaweb.success("任务提交成功", function () {
			window.location = window.location.pathname + "?keyType=" + that.keyType
		})
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

	/**
	 * 操作类型
	 */
	if (this.keyType == null || this.keyType.length == 0) {
		this.keyType = "key" // key | prefix
	}
})