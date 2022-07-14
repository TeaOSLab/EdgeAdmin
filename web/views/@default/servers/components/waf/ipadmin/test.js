Tea.context(function () {
	this.ip = ""
	this.result = {
		isDone: false,
		isOk: false,
		isFound: false,
		isAllowed: false,
		error: "",
		province: null,
		country: null,
		ipItem: null,
		ipList: null
	}

	this.$delay(function () {
		this.$watch("ip", function () {
			this.result.isDone = false
		})
	})

	this.success = function (resp) {
		this.result = resp.data.result
	}

	this.updateItem = function (itemId) {
		teaweb.popup(Tea.url(".updateIPPopup?firewallPolicyId=" + this.firewallPolicyId, {itemId: itemId}), {
			height: "30em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	/**
	 * 添加IP名单菜单
	 */
	this.createIP = function (type) {
		let that = this
		teaweb.popup("/servers/iplists/createIPPopup?listId=" + this.listId + '&type=' + type, {
			height: "30em",
			callback: function () {
				teaweb.success("保存成功", function () {
					window.location = "/servers/components/waf/ipadmin/lists?firewallPolicyId=" + that.firewallPolicyId + "&type=" + type
				})
			}
		})
	}
})