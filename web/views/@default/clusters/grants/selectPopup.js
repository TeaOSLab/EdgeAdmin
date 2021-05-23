Tea.context(function () {
	this.grantId = 0
	this.keyword = ""
	let allGrants = this.grants.$copy()

	this.selectGrant = function (grant) {
		NotifyPopup({
			code: 200,
			data: {
				grant: grant
			}
		})
	}

	this.$delay(function () {
		let that = this
		this.$watch("keyword", function (keyword) {
			if (keyword.length > 0) {
				that.grants = allGrants.$findAll(function (k, grant) {
					return teaweb.match(grant.name, keyword)
						|| teaweb.match(grant.description, keyword)
						|| teaweb.match(grant.username, keyword)
				})
			} else {
				that.grants = allGrants
			}
		})
	})
})