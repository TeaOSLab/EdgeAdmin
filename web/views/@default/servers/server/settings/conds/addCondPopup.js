Tea.context(function () {
	this.isUpdating = false
	this.cond = null
	this.paramsTitle = ""

	this.success = NotifyPopup
	this.condType = (this.components.length > 0) ? this.components[0].type : ""

	// 是否正在修改
	if (window.parent.UPDATING_COND != null) {
		this.isUpdating = true
		this.condType = window.parent.UPDATING_COND.type
		this.cond = window.parent.UPDATING_COND
	}

	this.changeCondType = function () {
		let that = this
		let c = this.components.$find(function (k, v) {
			return v.type == that.condType
		})
		if (c == null || c.paramsTitle.length == 0) {
			that.paramsTitle = "条件参数"
		} else {
			that.paramsTitle = c.paramsTitle
		}
	}

	this.$delay(function () {
		this.changeCondType()
	})
})