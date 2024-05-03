Tea.context(function () {
	this.isStarted = false
	this.isChecking = false
	this.result = {isOk: false, message: "", hasNew: false, dlURL: ""}
	this.isUpgraded = false

	this.$delay(function () {
		if (this.doCheck) {
			this.start()
		}
	})

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

	this.changeAutoCheck = function () {
		this.$post(".update")
			.params({
				autoCheck: this.config.autoCheck ? 1 : 0
			})
			.success(function () {
				teaweb.successToast("已保存")
			})
	}

	this.ignoreVersion = function (version) {
		teaweb.confirm("确定要忽略版本 v" + version + " 版本更新吗？", function () {
			this.$post(".ignoreVersion")
				.params({version: version})
				.success(function () {
					teaweb.reload()
				})
		})
	}

	this.resetIgnoredVersion = function (version) {
		teaweb.confirm("确定要重置已忽略版本吗？", function () {
			this.$post(".resetIgnoredVersion")
				.success(function () {
					teaweb.reload()
				})
		})
	}

	this.install = function (dlURL) {
		this.$post(".upgrade")
			.params({
				url: dlURL
			})
			.timeout(3600)
			.success(function () {
				this.$delay(function () {
					let msg = "下载覆盖成功"
					if (this.isUpgraded) {
						msg = "升级成功"
					}
					teaweb.success(msg + "，当前管理系统将会尝试自动重启，请刷新页面查看重启状态。如果长时间没能重启成功，请使用命令手动重启。", function () {
						teaweb.reload()
					})
				}, 3000)
			})

		this.isUpgrading = true
		this.updateUpgradeProgress()
	}

	if (this.isUpgrading) {
		this.$delay(function () {
			this.updateUpgradeProgress()
		})
	}

	this.updateUpgradeProgress = function () {
		if (!this.isUpgrading) {
			return
		}
		this.$get(".upgrade")
			.success(function (resp) {
				this.upgradeProgress = resp.data.upgradeProgress
				this.isUpgrading = resp.data.isUpgrading
				this.isUpgradingDB = resp.data.isUpgradingDB
				if (resp.data.isUpgradingDB) {
					this.isUpgraded = true
				}
			})
			.done(function () {
				this.$delay(function () {
					this.updateUpgradeProgress()
				}, 2000)
			})
	}
})