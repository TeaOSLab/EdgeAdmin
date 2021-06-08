Tea.context(function () {
	this.success = NotifyPopup

	this.group = {
		connector: "and", // 默认为and，更符合用户的直觉
		description: "",
		isReverse: false,
		conds: [],
		isOn: true
	}

	this.isUpdating = false

	// 是否在修改
	this.$delay(function () {
		if (window.parent.UPDATING_COND_GROUP != null) {
			this.group = window.parent.UPDATING_COND_GROUP
			this.isUpdating = true
		} else if (this.group.conds.length == 0) {
			// 如果尚未有条件，则自动弹出添加界面
			this.addCond()
		}
	})

	// 条件类型名称
	this.typeName = function (cond) {
		let c = this.components.$find(function (k, v) {
			return v.type == cond.type
		})
		if (c != null) {
			return c.name;
		}
		return cond.param + " " + cond.operator
	}

	// 添加条件
	this.addCond = function () {
		window.UPDATING_COND = null

		let that = this

		teaweb.popup("/servers/server/settings/conds/addCondPopup", {
			width: "32em",
			height: "22em",
			callback: function (resp) {
				that.group.conds.push(resp.data.cond)
			}
		})
	}

	// 删除条件
	this.removeCond = function (condIndex) {
		let that = this
		teaweb.confirm("确定要删除此条件？", function () {
			that.group.conds.$remove(condIndex)
		})
	}

	// 修改条件
	this.updateCond = function (condIndex, cond) {
		window.UPDATING_COND = cond
		let that = this

		teaweb.popup("/servers/server/settings/conds/addCondPopup", {
			width: "32em",
			height: "22em",
			callback: function (resp) {
				Vue.set(that.group.conds, condIndex, resp.data.cond)
			}
		})
	}
})