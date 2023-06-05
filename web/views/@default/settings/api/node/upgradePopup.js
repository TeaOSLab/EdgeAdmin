Tea.context(function () {
	this.$delay(function () {
		this.checkLoop()
	})

	this.success = function () {
	}

	this.isRequesting = false

	this.before = function () {
		this.isRequesting = true
	}

	this.done = function () {
		this.isRequesting = false
	}

	this.checkLoop = function () {
		if (this.currentVersion == this.latestVersion) {
			return
		}

		this.$post(".upgradeCheck")
			.params({
				nodeId: this.nodeId
			})
			.success(function (resp) {
				if (resp.data.isOk) {
					teaweb.reload()
				}
			})
			.done(function () {
				this.$delay(function () {
					this.checkLoop()
				}, 3000)
			})
	}
})