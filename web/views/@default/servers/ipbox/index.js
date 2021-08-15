Tea.context(function () {
	this.blackListsVisible = false
	this.allPublicBlackIPLists = this.publicBlackIPLists.$copy()

	this.showBlackLists = function () {
		let that = this
		this.publicBlackIPLists = this.allPublicBlackIPLists.filter(function (allList) {
			let found = true
			that.ipLists.forEach(function (list) {
				if (allList.id == list.id) {
					found = false
				}
			})
			return found
		})
		this.blackListsVisible = !this.blackListsVisible
	}

	this.addBlackIP = function (list) {
		this.$post(".addIP")
			.params({
				listId: list.id,
				ip: this.ip
			})
			.success(function () {
				this.ipLists.push(list)
				this.blackListsVisible = false
			})
	}
})