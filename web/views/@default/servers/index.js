Tea.context(function () {
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
			.success(function (resp) {
				let status = resp.data.status
				this.servers.forEach(function (server) {
					if (typeof status[server.id] === "object") {
						server.status = status[server.id]
					}
				})
			})
	}
});