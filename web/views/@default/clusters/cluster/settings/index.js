Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.timeZoneGroupCode = "asia"
	if (this.timeZoneLocation != null) {
		this.timeZoneGroupCode = this.timeZoneLocation.group
	}

	this.$delay(function () {
		this.$watch("timeZoneGroupCode", function (groupCode) {
			let firstLocation = null
			this.timeZoneLocations.forEach(function (v) {
				if (firstLocation != null) {
					return
				}
				if (v.group == groupCode) {
					firstLocation = v
				}
			})
			if (firstLocation != null) {
				this.cluster.timeZone = firstLocation.name
			}
		})
	})
})