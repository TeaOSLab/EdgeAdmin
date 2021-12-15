Tea.context(function () {
	this.isStarted = false
	this.isChecking = true
	this.result = {isOk: false, message: "", hasNew: false, dlURL: ""}

	this.start = function () {
		this.isStarted = true
		this.isChecking = true

		this.$delay(function () {
			this.check()
		}, 1000)
	}

	this.check = function () {
		this.$post("$")
			.success(function (resp) {
				this.result = resp.data.result
			})
			.done(function () {
				this.isChecking = false
			})
	}
})