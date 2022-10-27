Tea.context(function () {
	this.selectedProvider = null
	this.accounts = []
	this.accountId = 0

	this.changeProvider = function () {
		this.accountId = 0

		if (this.providerCode.length == 0) {
			return
		}

		let that = this
		let provider = this.providers.$find(function (k, v) {
			return v.code == that.providerCode
		})
		if (provider == null) {
			return
		}

		this.selectedProvider = provider

		this.$post(".accountsWithCode")
			.params({
				code: provider.code
			})
			.success(function (resp) {
				this.accounts = resp.data.accounts
			})
	}

	if (this.providerCode.length > 0) {
		this.changeProvider()
	}

	this.addAccount = function () {
		let that = this
		teaweb.popup("/servers/certs/acme/accounts/createPopup?providerCode=" + this.providerCode, {
			height: "24em",
			callback: function () {
				teaweb.successToast("创建成功，已自动选中", 1500, function () {
					that.$post(".accountsWithCode")
						.params({
							code: that.providerCode
						})
						.success(function (resp) {
							that.accounts = resp.data.accounts

							if (that.accounts.length > 0) {
								that.accountId = that.accounts[0].id
							}
						})
				})
			}
		})
	}
})