Tea.context(function () {
	this.success = NotifyPopup

	this.type = this.policy.type
	this.authDescription = ""

	this.$delay(function () {
		this.changeType()
	})

	this.changeType = function () {
		let that = this
		let authType = this.authTypes.$find(function (k, v) {
			return v.code == that.type
		})
		if (authType != null) {
			this.policy.typeName = authType.name
			this.authDescription = authType.description
		} else {
			this.authDescription = ""
		}
	}

	/**
	 * 基本认证
	 */
	this.moreBasicAuthOptionsVisible = false

	this.showMoreBasicAuthOptions = function () {
		this.moreBasicAuthOptionsVisible = !this.moreBasicAuthOptionsVisible
	}

	/**
	 * 子请求
	 */
	this.subRequestFollowRequest = (this.policy.params.method != null && this.policy.params.method.length > 0) ? 0 : 1
})