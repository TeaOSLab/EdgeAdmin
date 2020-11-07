Tea.context(function () {
	this.letterGroups = [
		"ABC", "DEF", "GHI", "JKL", "MNO", "PQR", "STU", "VWX", "YZ"
	];
	this.selectedGroup = "ABC"
	this.letterCountries = {}
	let that = this
	this.countries.forEach(function (country) {
		if (typeof (that.letterCountries[country.letter]) == "undefined") {
			that.letterCountries[country.letter] = []
		}
		that.letterCountries[country.letter].push(country)
	})
	this.isCheckingAll = false

	this.selectGroup = function (group) {
		this.selectedGroup = group
	}

	this.selectCountry = function (country) {
		country.isChecked = !country.isChecked
	}

	this.deselectCountry = function (country) {
		country.isChecked = false
	}

	this.checkAll = function () {
		this.isCheckingAll = !this.isCheckingAll

		this.countries.forEach(function (country) {
			country.isChecked = that.isCheckingAll
		})
	}

	this.success = function () {
		teaweb.successToast("保存成功")
	}


	/**
	 * 添加IP名单菜单
	 */
	this.createIP = function (type) {
		teaweb.popup("/servers/components/waf/ipadmin/createIPPopup?firewallPolicyId=" + this.firewallPolicyId + '&type=' + type, {
			height: "23em",
			callback: function () {
				window.location = "/servers/components/waf/ipadmin/lists?firewallPolicyId=" + this.firewallPolicyId + "&type=" + type
			}
		})
	}
})