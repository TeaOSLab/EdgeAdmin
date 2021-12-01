Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")
	this.auditSuccess = NotifyReloadSuccess("保存成功")

	this.allServerNames = []
	let that = this
	this.serverNames.forEach(function (v) {
		if (v.subNames == null || v.subNames.length == 0) {
			that.allServerNames.push({
				name: v.name,
				isPassed: that.passedDomains.$contains(v.name)
			})
		} else {
			v.subNames.forEach(function (subName) {
				that.allServerNames.push({
					name: subName,
					isPassed: that.passedDomains.$contains(subName)
				})
			})
		}
	})

	this.hasPassedDomains = false
	this.allServerNames.forEach(function (serverName) {
		if (serverName.isPassed) {
			that.hasPassedDomains = true
		}
	})

	this.auditing = 1
})