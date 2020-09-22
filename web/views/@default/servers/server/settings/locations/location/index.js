Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.$delay(function () {
		this.changePatternType(this.type)
	})

	this.selectedType = null

	this.changePatternType = function (type) {
		this.selectedType = this.patternTypes.$find(function (k, v) {
			return v.type == type;
		})
	}
})