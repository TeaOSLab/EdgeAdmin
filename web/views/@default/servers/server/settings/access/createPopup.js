Tea.context(function () {
	this.success = NotifyPopup

	this.type = ""
	this.authDescription = ""
	this.rawDescription = ""

	this.changeType = function () {
		let that = this
		let authType = this.authTypes.$find(function (k, v) {
			return v.code == that.type
		})
		if (authType != null) {
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
	 * 基本认证
	 */
	this.moreBasicAuthOptionsVisible = false

	this.showMoreBasicAuthOptions = function () {
		this.moreBasicAuthOptionsVisible = !this.moreBasicAuthOptionsVisible
	}

	/**
	 * 子请求
	 */
	this.subRequestFollowRequest = 1
})