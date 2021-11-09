Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	/**
	 * 用户相关
	 */
	this.userId = 0
	this.plans = []
	this.userPlanId = 0

	if (this.userPlan != null) {
		this.userPlanId = this.userPlan.id
	}

	this.changeUserId = function (v) {
		this.userId = v

		if (this.userId == 0) {
			this.plans = []
			return
		}

		this.$post("/servers/users/plans")
			.params({
				userId: this.userId,
				serverId: this.serverId
			})
			.success(function (resp) {
				this.plans = resp.data.plans
			})
	}

	if (this.user != null) {
		this.changeUserId(this.user.id)
	}
})