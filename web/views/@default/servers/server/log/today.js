Tea.context(function () {
    let that = this
    this.accessLogs.forEach(function (accessLog) {
        if (typeof (that.regions[accessLog.remoteAddr]) == "string") {
            accessLog.region = that.regions[accessLog.remoteAddr]
        } else {
            accessLog.region = ""
        }
    })
})