Tea.context(function () {
	this.windowWidth = window.innerWidth
	this.miniWidth = 760
	this.columnWidth1 = 800
	this.columnWidth2 = 900
	this.columnWidth3 = 1000
	this.columnWidth4 = 1100
	this.columnWidth5 = 1200

	this.servers.forEach(function (v) {
		v["status"] = {
			isOk: false,
			message: "",
			todo: ""
		}
	})

	this.$delay(function () {
		if (this.checkDNS) {
			this.loadStatus()
		}

		let that = this
		this.$watch("checkDNS", function (v) {
			if (v) {
				that.loadStatus()
			}
		})
	})

	this.loadStatus = function () {
		let serverIds = this.servers.map(function (v) {
			return v.id
		})
		this.$post(".status")
			.params({
				serverIds: serverIds
			})
			.timeout(300)
			.success(function (resp) {
				let status = resp.data.status
				this.servers.forEach(function (server) {
					if (typeof status[server.id] === "object") {
						server.status = status[server.id]
					}
				})
			})
	}

	/**
	 * 最近使用
	 */
	this.latestVisible = false

	this.showLatest = function () {
		this.latestVisible = !this.latestVisible
	}

	/**
	 * 全选
	 */
	this.checkedServerIds = []
	this.changeAllChecked = function (checked) {
		for (let checkbox of this.$refs.serverCheckboxes) {
			if (checked) {
				checkbox.check()
			} else {
				checkbox.uncheck()
			}
		}
		this.updateCheckedServers()
	}

	this.changeServerChecked = function () {
		this.updateCheckedServers()
	}

	this.updateCheckedServers = function () {
		let serverIds = []
		for (let checkbox of this.$refs.serverCheckboxes) {
			if (checkbox.isChecked()) {
				serverIds.push(checkbox.vValue)
			}
		}
		this.checkedServerIds = serverIds
	}

	this.resetCheckedServers = function () {
		this.$refs.allCheckedCheckboxes.uncheck()
		for (let checkbox of this.$refs.serverCheckboxes) {
			checkbox.uncheck()
		}
		this.updateCheckedServers()
	}

	this.deleteServers = function () {
		let that = this
		teaweb.confirm("确定要删除所选的" + (this.checkedServerIds.length) + "个网站吗？", function () {
			that.$post(".deleteServers")
				.params({
					serverIds: this.checkedServerIds
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}
})