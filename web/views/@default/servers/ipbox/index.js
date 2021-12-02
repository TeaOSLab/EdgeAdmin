Tea.context(function () {
	this.blackListsVisible = false
	this.allPublicBlackIPLists = this.publicBlackIPLists.$copy()


	this.defaultItemExpiredAt = Math.floor(new Date().getTime() / 1000) + 86400
	this.showBlackLists = function () {
		this.defaultItemExpiredAt = Math.floor(new Date().getTime() / 1000) + 86400

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

	this.selectedListId = 0
	this.addBlackIP = function (list) {
		this.selectedListId = list.id
		let expiredAt = this.$refs.itemExpiredTimestamp.resultTimestamp()
		let that = this
		teaweb.confirm("确定要将此IP添加到黑名单'" + list.name + "'吗？", function () {
			that.$post(".addIP")
				.params({
					listId: list.id,
					ip: that.ip,
					expiredAt: expiredAt
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}

	this.deleteFromList = function (listId, itemId) {
		teaweb.confirm("确定要从此名单中删除此IP吗？", function () {
			this.$post(".deleteFromList")
				.params({
					listId: listId,
					itemId: itemId
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}
})