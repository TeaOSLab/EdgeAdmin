Tea.context(function () {
	this.success = NotifyPopup
	this.regionId = 0

	this.selectRegion = function (region) {
		this.regionId = region.id
	}
})