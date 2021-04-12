Tea.context(function () {
	this.selectedRecipient = null
	this.selectedGroup = null


	this.selectRecipient = function (recipient) {
		this.selectedRecipient = recipient
		this.selectedGroup = null
	}

	this.selectGroup = function (group) {
		this.selectedRecipient = null
		this.selectedGroup = group
	}

	this.confirm = function () {
		if (this.selectedRecipient != null) {
			NotifyPopup({
				code: 200,
				data: {
					id: this.selectedRecipient.id,
					name: this.selectedRecipient.name,
					subName: this.selectedRecipient.instanceName,
					type: "recipient"
				}
			})
		} else if (this.selectedGroup != null) {
			NotifyPopup({
				code: 200,
				data: {
					id: this.selectedGroup.id,
					name: this.selectedGroup.name,
					type: "group"
				}
			})
		} else {
			teaweb.closePopup()
		}
	}
})