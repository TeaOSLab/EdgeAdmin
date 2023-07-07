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
})