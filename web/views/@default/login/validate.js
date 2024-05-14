Tea.context(function () {
	this.$delay(function () {
		let sid = localStorage.getItem("sid")
		let ip = localStorage.getItem("ip")

		if (sid == null || sid.length == 0 || ip == null || ip.length == 0) {
			window.location = "/logout"
			return
		}

		this.$post("$")
			.params({localSid: sid, "ip": ip})
			.success(function (resp) {
				if (!resp.data.isOk) {
					window.location = "/logout"
					return
				}

				// renew local data
				localStorage.setItem("sid", resp.data.localSid)
				localStorage.setItem("ip", resp.data.ip)

				// redirect back (MUST delay)
				this.$delay(function () {
					if (this.from.length > 0) {
						window.location = this.from
					} else {
						window.location = "/dashboard"
					}
				}, 100)
			})
	})
})