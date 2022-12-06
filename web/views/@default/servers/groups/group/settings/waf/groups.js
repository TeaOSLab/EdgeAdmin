Tea.context(function () {
	// 排序
	this.$delay(function () {
		let that = this
		sortTable(function () {
			let groupIds = []
			document.querySelectorAll("tbody[data-group-id]")
				.forEach(function (v) {
					groupIds.push(v.getAttribute("data-group-id"))
				})

			that.$post("/servers/components/waf/sortGroups")
				.params({
					firewallPolicyId: that.firewallPolicyId,
					type: that.type,
					groupIds: groupIds
				})
				.success(function () {
					teaweb.successToast("排序保存成功")
				})
		})
	})

	// 启用
	this.enableGroup = function (groupId) {
		this.$post("/servers/components/waf/updateGroupOn")
			.params({
				groupId: groupId,
				isOn: 1
			})
			.refresh()

	}

	// 停用
	this.disableGroup = function (groupId) {
		this.$post("/servers/components/waf/updateGroupOn")
			.params({
				groupId: groupId,
				isOn: 0
			})
			.refresh()
	}

	// 删除
	this.deleteGroup = function (groupId) {
		teaweb.confirm("确定要删除此规则分组吗？", function () {
			this.$post("/servers/components/waf/deleteGroup")
				.params({
					firewallPolicyId: this.firewallPolicyId,
					groupId: groupId
				})
				.refresh()
		})
	}

	// 添加分组
	this.createGroup = function (type) {
		teaweb.popup("/servers/components/waf/createGroupPopup?firewallPolicyId=" + this.firewallPolicyId + "&type=" + type, {
			callback: function () {
				teaweb.success("保存成功", function () {
					window.location.reload()
				})
			}
		})
	}
})