Tea.context(function () {
	this.result = {
		isInstalling: false,
		isInstalled: false,
		isOk: false,
		err: "",

		user: "",
		password: "",
		dir: "",

		logs: []
	}

	this.$delay(function () {
		this.checkStatus()
	})

	this.install = function () {
		this.result.isInstalling = true
		this.result.isInstalled = false
		this.result.logs = []

		this.$post(".installPopup")
			.timeout(3600)
			.success(function (resp) {
				this.result.isOk = resp.data.isOk
				if (!resp.data.isOk) {
					this.result.err = resp.data.err
				} else {
					this.result.user = resp.data.user
					this.result.password = resp.data.password
					this.result.dir = resp.data.dir
				}
				this.result.isInstalled = true
				this.result.isInstalling = false
			})
	}

	this.checkStatus = function () {
		if (!this.result.isInstalling) {
			this.$delay(function () {
				this.checkStatus()
			}, 1000)
			return
		}

		this.$post(".installLogs")
			.success(function (resp) {
				let that = this
				resp.data.logs.forEach(function (log) {
					that.result.logs.unshift(log)
				})
			})
			.done(function () {
				this.$delay(function () {
					this.checkStatus()
				}, 2000)
			})
	}

	this.finish = function () {
		teaweb.closePopup()
	}
})