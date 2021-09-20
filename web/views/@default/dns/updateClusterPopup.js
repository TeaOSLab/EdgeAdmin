Tea.context(function () {
	this.$delay(function () {
		this.changeProviderType()
		this.changeProvider()

		this.$watch("providerId", function () {
			this.changeProvider()
		})
		this.$watch("domainId", function () {
			this.changeDomain()
		})
	})

	this.success = NotifyPopup

	// 初始化的内容
	// this.domainId = 0
	// this.domain = ""
	// this.providerId = 0

	if (this.providerType == "") {
		this.providerType = this.providerTypes[0].code
	}
	this.providers = []
	this.domains = []

	this.changeProviderType = function () {
		this.$post(".providerOptions")
			.params({
				type: this.providerType
			})
			.success(function (resp) {
				this.providers = resp.data.providers

				// 检查providerId
				if (this.providers.length == 0) {
					this.providerId = 0
					return
				}
				let that = this
				if (this.providers.$find(function (k, v) {
					return v.id == that.providerId
				}) == null) {
					this.providerId = this.providers[0].id
				}
				this.changeProvider()
			})
	}

	this.changeProvider = function () {
		this.$post(".domainOptions")
			.params({
				providerId: this.providerId
			})
			.success(function (resp) {
				this.domains = resp.data.domains
				this.changeDomain()
			})
	}

	this.changeDomain = function () {
		if (this.domains.length == 0) {
			this.domainId = 0
			this.domain = ""
			return
		}

		let domainId = this.domainId
		let domainInfo = this.domains.$find(function (k, v) {
			return v.id == domainId
		})
		if (domainInfo == null) {
			// 默认选取第一个
			this.domainId = this.domains[0].id
			this.domain = this.domains[0].name
		} else {
			this.domain = domainInfo.name
		}
	}

	/**
	 * 自动设置CNAME
	 */
	this.addCnameRecord = function (name) {
		this.$refs.cnameRecords.addValue(name)
	}
})