Tea.context(function () {
	this.success = NotifyPopup

	// rules
	this.rules = this.setConfig.rules

	// connector
	this.selectedConnector = this.setConfig.connector
	this.selectedConnectorDescription = ""
	this.changeConnector = function () {
		let that = this
		this.selectedConnectorDescription = this.connectors.$find(function (k, v) {
			return v.value == that.selectedConnector
		}).description
	}
	this.changeConnector()

	// action
	this.action = this.setConfig.action

	// action:go_group
	this.actionGroupId = 0
	if (this.action == "go_group" || this.action == "go_set" && this.setConfig.actionOptions != null) {
		this.$delay(function () {
			this.actionGroupId = this.setConfig.actionOptions["groupId"]
		})
	}

	// action:go_set
	this.actionSetId = 0
	if (this.action == "go_set" && this.setConfig.actionOptions != null) {
		this.$delay(function () {
			this.actionSetId = this.setConfig.actionOptions["setId"]
		})
	}

	this.groupSets = function (groupId) {
		if (this.firewallPolicy == null) {
			return
		}
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
})