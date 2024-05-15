Tea.context(function () {
	this.$delay(function () {
		let that = this
		sortTable(function () {
			let setIds = []
			document
				.querySelectorAll("tbody[data-set-id]")
				.forEach(function (v) {
					setIds.push(v.getAttribute("data-set-id"))
				})
			that.$post("/servers/components/waf/sortSets")
				.params({
					groupId: that.group.id,
					setIds: setIds
				})
				.success(function () {
					teaweb.successToast("排序保存成功")
				})
		})
	})

	// 更改分组
	this.updateGroup = function (groupId) {
		teaweb.popup("/servers/components/waf/updateGroupPopup?groupId=" + groupId, {
			height: "20em",
			callback: function () {
				teaweb.success("保存成功", function () {
					window.location.reload()
				})
			}
		})
	}

	// 创建规则集
	this.createSet = function (groupId) {
		teaweb.popup("/servers/components/waf/createSetPopup?firewallPolicyId=" + this.firewallPolicyId + "&groupId=" + groupId + "&type=" + this.type, {
			width: "50em",
			height: "40em",
			callback: function () {
				teaweb.success("保存成功", function () {
					window.location.reload()
				})
			}
		})
	}

	// 修改规则集
	this.updateSet = function (setId) {
		teaweb.popup("/servers/components/waf/updateSetPopup?firewallPolicyId=" + this.firewallPolicyId + "&groupId=" + this.group.id + "&type=" + this.type + "&setId=" + setId, {
			width: "50em",
			height: "40em",
			callback: function () {
				teaweb.success("保存成功", function () {
					window.location.reload()
				})
			}
		})
	}

	// 停用|启用规则集
	this.updateSetOn = function (setId, isOn) {
		this.$post("/servers/components/waf/updateSetOn")
			.params({
				setId: setId,
				isOn: isOn ? 1 : 0
			})
			.refresh()
	}

	// 删除规则集
	this.deleteSet = function (setId) {
		let that = this
		teaweb.confirm("确定要删除此规则集吗？", function () {
			that.$post("/servers/components/waf/deleteSet")
				.params({
					groupId: this.group.id,
					setId: setId
				})
				.refresh()
		})
	}
})