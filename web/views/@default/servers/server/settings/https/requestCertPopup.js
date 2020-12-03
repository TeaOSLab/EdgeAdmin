Tea.context(function () {
	this.isRequesting = false
	this.userId = 0

	this.remove = function (index) {
		this.serverNames.$remove(index)
	}

	this.beforeSubmit = function () {
		this.isRequesting = true
	}

	this.fail = function (resp) {
		this.isRequesting = false
		teaweb.warn(resp.message)
		if (resp.data.acmeUser != null) {
			this.users.push({
				id: resp.data.acmeUser.id,
				email: resp.data.acmeUser.email,
				description: ""
			})
			this.userId = resp.data.acmeUser.id
		}
	}
})