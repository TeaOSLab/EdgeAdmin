Tea.context(function () {
	this.success = NotifySuccess("添加成功", "/servers/server/settings/locations?serverId=" + this.serverId)

	this.type = 1
	this.selectedType = this.patternTypes[0]


	this.changePatternType = function (type) {
		this.selectedType = this.patternTypes.$find(function (k, v) {
			return v.type == type;
		})
	}
})