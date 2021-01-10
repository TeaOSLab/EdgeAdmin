Tea.context(function () {
	this.isCreating = true
	if (window.parent.UPDATING_REDIRECT != null) {
		this.isCreating = false
		this.redirect = window.parent.UPDATING_REDIRECT
	} else {
		this.redirect = {
			status: 0,
			beforeURL: "",
			afterURL: ""
		}
	}
})