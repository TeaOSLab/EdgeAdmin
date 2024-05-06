Tea.context(function () {
	this.success = NotifyPopup

	// rules
	this.rules = []

	// connector
	this.selectedConnector = this.connectors[1].value
	this.selectedConnectorDescription = ""
	this.changeConnector = function () {
		let that = this
		this.selectedConnectorDescription = this.connectors.$find(function (k, v) {
			return v.value == that.selectedConnector
		}).description
	}
	this.changeConnector()

	// action
	this.action = "block"

	// action:go_group
	this.actionGroupId = 0

	// action:go_set
	this.actionSetId = 0
	this.groupSets = function (groupId) {
		let group = null
		this.firewallPolicy.inbound.groups.forEach(function (v) {
			if (v.id == groupId) {
				group = v
			}
		})
		this.firewallPolicy.outbound.groups.forEach(function (v) {
			if (v.id == groupId) {
				group = v
			}
		})
		if (group == null) {
			return []
		}
		return group.sets
	}

	// 使用代码
	this.useCode = false
	this.switchToCode = function () {
		this.useCode = !this.useCode

		if (this.useCode) {
			this.$delay(function () {
				this.$refs.codeInput.focus()
			})
		}
	}
})