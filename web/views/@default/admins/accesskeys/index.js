Tea.context(function () {
	this.createAccessKey = function () {
		teaweb.popup("/admins/accesskeys/createPopup?adminId=" + this.admin.id, {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateAccessKeyIsOn = function (accessKeyId, isOn) {
		let that = this

		let message = ""
		if (isOn) {
			message = "确定要启用此AccessKey吗？"
		} else {
			message = "确定要禁用此AccessKey吗？"
		}
		teaweb.confirm(message, function () {
			that.$post(".updateIsOn")
				.params({
					accessKeyId: accessKeyId,
					isOn: isOn ? 1 : 0
				})
				.refresh()
		})
	}

	this.deleteAccessKey = function (accessKeyId) {
		let that = this
		teaweb.confirm("确定要删除此AccessKey吗？", function () {
			that.$post(".delete")
				.params({
					accessKeyId: accessKeyId
				})
				.refresh()
		})
	}
})