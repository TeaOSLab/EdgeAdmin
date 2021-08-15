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
		let that = this
		teaweb.confirm("确定要将此IP添加到此黑名单吗？", function () {
			that.$post(".addIP")
				.params({
					listId: list.id,
					ip: that.ip
				})
				.success(function () {
					that.ipLists.push(list)
					that.blackListsVisible = false
				})
		})
	}
})