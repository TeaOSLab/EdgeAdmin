Tea.context(function () {
	this.success = NotifySuccess("保存成功", "/clusters/cluster?clusterId=" + this.clusterId)

	this.$delay(function () {
		this.$refs.ipList.focus()
	})
})