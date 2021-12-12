Tea.context(function () {
	this.isUpdating = false
	this.cond = null
	this.paramsTitle = ""
	this.paramsCaseInsensitive = false
	this.isCaseInsensitive = false

	this.success = NotifyPopup
	this.condType = (this.components.length > 0) ? this.components[0].type : ""

	// 是否正在修改
	if (window.parent.UPDATING_COND != null) {
		this.isUpdating = true
		this.condType = window.parent.UPDATING_COND.type
		this.cond = window.parent.UPDATING_COND
		if (typeof (this.cond.isCaseInsensitive) == "boolean") {
			this.isCaseInsensitive = this.cond.isCaseInsensitive
		}
	}

	this.changeCondType = function () {
		let that = this
		let c = this.components.$find(function (k, v) {
			return v.type == that.condType
		})
		if (c == null || c.paramsTitle.length == 0) {
			that.paramsTitle = "条件参数"
			that.paramsCaseInsensitive = false
		} else {
			that.paramsTitle = c.paramsTitle
			if (typeof (c.caseInsensitive) != "undefined") {
				that.paramsCaseInsensitive = c.caseInsensitive
				that.$delay(function () {
					that.changeCaseInsensitive()
				})
			} else {
				that.paramsCaseInsensitive = false
			}
		}
	}

	this.$delay(function () {
		this.changeCondType()
	})

	this.changeCaseInsensitive = function () {
		let componentRef = this.$refs.component
		if (componentRef == null) {
			return
		}
		if (typeof (componentRef.changeCaseInsensitive) == "function") {
			componentRef.changeCaseInsensitive(this.isCaseInsensitive)
		}
	}
})