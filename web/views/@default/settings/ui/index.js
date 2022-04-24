Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.$delay(function () {
		this.$watch("config.showOpenSourceInfo", function (v) {
			this.teaShowOpenSourceInfo = v
		})
	})

	// 时区
	this.timeZoneGroupCode = "asia"
	if (this.timeZoneLocation != null) {
		this.timeZoneGroupCode = this.timeZoneLocation.group
	}

	let oldTimeZoneGroupCode = this.timeZoneGroupCode
	let oldTimeZoneName = ""
	if (this.timeZoneLocation != null) {
		oldTimeZoneName = this.timeZoneLocation.name
	}

	this.$delay(function () {
		this.$watch("timeZoneGroupCode", function (groupCode) {
			if (groupCode == oldTimeZoneGroupCode && oldTimeZoneName.length > 0) {
				this.config.timeZone = oldTimeZoneName
				return
			}
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
				this.config.timeZone = firstLocation.name
			}
		})
	})
})