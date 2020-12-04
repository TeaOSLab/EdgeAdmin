Tea.context(function () {
	this.success = NotifySuccess("保存成功", "/users/user?userId=" + this.user.id)

	this.passwordEditing = false

	this.changePasswordEditing = function () {
		this.passwordEditing = !this.passwordEditing
	}
})