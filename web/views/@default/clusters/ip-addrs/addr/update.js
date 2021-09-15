Tea.context(function () {
	this.success = NotifySuccess("保存成功", Tea.url(".", {addrId: this.addr.id}))

	this.addr.isUp = this.addr.isUp ? 1 : 0
})