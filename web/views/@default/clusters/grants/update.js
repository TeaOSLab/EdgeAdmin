Tea.context(function () {
	this.method = this.grant.method;

	this.success = NotifySuccess("保存成功", "/clusters/grants/grant?grantId=" + this.grant.id);
});