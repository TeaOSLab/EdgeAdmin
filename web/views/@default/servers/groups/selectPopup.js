Tea.context(function () {
	this.success = NotifyPopup
	this.groupId = 0

	this.selectGroup = function (group) {
		this.groupId = group.id
	}
})