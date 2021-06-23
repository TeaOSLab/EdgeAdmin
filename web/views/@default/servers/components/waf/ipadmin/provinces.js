Tea.context(function () {
	this.isCheckingAll = false

	this.countSelectedProvinces = this.provinces.$count(function (k, province) {
		return province.isChecked
	})

	this.selectProvince = function (province) {
		province.isChecked = !province.isChecked
		this.change()
	}

	this.deselectProvince = function (province) {
		province.isChecked = false
		this.change()
	}

	this.checkAll = function () {
		this.isCheckingAll = !this.isCheckingAll
		let that = this
		this.provinces.forEach(function (province) {
			province.isChecked = that.isCheckingAll
		})

		this.change()
	}

	this.success = function () {
		teaweb.success("保存成功", function () {
			teaweb.reload()
		})
	}


	this.change = function () {
		this.countSelectedProvinces = this.provinces.$count(function (k, province) {
			return province.isChecked
		})
	}

	/**
	 * 添加IP名单菜单
	 */
	this.createIP = function (type) {
		let that = this
		teaweb.popup("/servers/components/waf/ipadmin/createIPPopup?firewallPolicyId=" + this.firewallPolicyId + '&type=' + type, {
			height: "23em",
			callback: function () {
				teaweb.success("保存成功", function () {
					window.location = "/servers/components/waf/ipadmin/lists?firewallPolicyId=" + that.firewallPolicyId + "&type=" + type
				})
			}
		})
	}
})