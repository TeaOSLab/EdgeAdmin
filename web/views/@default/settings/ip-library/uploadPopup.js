Tea.context(function () {
	this.isRequesting = false
	this.selectedTypeCode = this.types[0].code
	this.selectedTypeDescription = this.types[0].description
	this.selectedTypeExt = this.types[0].ext

	this.success = NotifyPopup

	this.before = function () {
		this.isRequesting = true
	}

	this.done = function () {
		this.isRequesting = false
	}

	this.changeType = function () {
		let that = this
		let selectedType = this.types.$find(function (k, v) {
			return v.code == that.selectedTypeCode
		})
		if (selectedType == null) {
			return
		}
		this.selectedTypeDescription = selectedType.description
		this.selectedTypeExt = selectedType.ext
	}
})