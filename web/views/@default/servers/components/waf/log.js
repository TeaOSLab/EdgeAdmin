Tea.context(function () {
	this.$delay(function () {
		let that = this
		teaweb.datepicker("day-input", function (day) {
			that.day = day
		})
	})

    let that = this
    this.accessLogs.forEach(function (accessLog) {
        if (typeof (that.regions[accessLog.remoteAddr]) == "string") {
            accessLog.region = that.regions[accessLog.remoteAddr]
        } else {
            accessLog.region = ""
        }
    })
})