Tea.context(function () {
	this.isCheckingAll = false

	this.selectProvince = function (province) {
		province.isChecked = !province.isChecked
	}

	this.deselectProvince = function (province) {
		province.isChecked = false
	}

	this.checkAll = function () {
		this.isCheckingAll = !this.isCheckingAll

		this.provinces.forEach(function (province) {
			province.isChecked = that.isCheckingAll
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