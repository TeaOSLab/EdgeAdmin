Tea.context(function () {
	this.success = NotifyPopup

	this.type = this.policy.type
	this.authDescription = ""
	this.rawDescription = ""

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
			this.rawDescription = authType.description
		} else {
			this.authDescription = ""
			this.rawDescription = ""
		}
	}

	/**
	 * TypeA
	 */
	this.typeASecret = ""
	this.typeASignParamName = "sign"

	if (this.policy.type == "typeA") {
		this.typeASecret = this.policy.params.secret
		this.typeASignParamName = this.policy.params.signParamName
		this.$delay(function () {
			this.changeTypeASignParamName()
		})
	}

	this.generateTypeASecret = function () {
		this.$post(".random")
			.success(function (resp) {
				this.typeASecret = resp.data.random
			})
	}

	this.changeTypeASignParamName = function () {
		this.authDescription = this.rawDescription.replace("sign=", this.typeASignParamName + "=")
	}

	/**
	 * TypeB
	 */
	this.typeBSecret = ""

	if (this.policy.type == "typeB") {
		this.typeBSecret = this.policy.params.secret
	}

	this.generateTypeBSecret = function () {
		this.$post(".random")
			.success(function (resp) {
				this.typeBSecret = resp.data.random
			})
	}

	/**
	 * TypeC
	 */
	this.typeCSecret = ""

	if (this.policy.type == "typeC") {
		this.typeCSecret = this.policy.params.secret
	}

	this.generateTypeCSecret = function () {
		this.$post(".random")
			.success(function (resp) {
				this.typeCSecret = resp.data.random
			})
	}

	/**
	 * TypeD
	 */
	this.typeDSecret = ""
	this.typeDSignParamName = "sign"
	this.typeDTimestampParamName = "t"

	if (this.policy.type == "typeD") {
		this.typeDSecret = this.policy.params.secret
		this.typeDSignParamName = this.policy.params.signParamName
		this.typeDTimestampParamName = this.policy.params.timestampParamName
		this.$delay(function () {
			this.changeTypeDSignParamName()
			this.changeTypeDTimestampParamName()
		})
	}

	this.generateTypeDSecret = function () {
		this.$post(".random")
			.success(function (resp) {
				this.typeDSecret = resp.data.random
			})
	}

	this.changeTypeDSignParamName = function () {
		this.authDescription = this.rawDescription.replace("sign=", this.typeDSignParamName + "=")
		this.authDescription = this.authDescription.replace("t=", this.typeDTimestampParamName + "=")
	}

	this.changeTypeDTimestampParamName = function () {
		this.authDescription = this.rawDescription.replace("sign=", this.typeDSignParamName + "=")
		this.authDescription = this.authDescription.replace("t=", this.typeDTimestampParamName + "=")
	}

	/**
	 * 基本鉴权
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