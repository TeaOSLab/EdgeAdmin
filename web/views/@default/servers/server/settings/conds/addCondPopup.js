Tea.context(function () {
	this.isUpdating = false
	this.cond = null

	this.success = NotifyPopup
	this.condType = (this.components.length > 0) ? this.components[0].type : ""

	// 是否正在修改
	if (window.parent.UPDATING_COND != null) {
		this.isUpdating = true
		this.condType = window.parent.UPDATING_COND.type
		this.cond = window.parent.UPDATING_COND
	}
})