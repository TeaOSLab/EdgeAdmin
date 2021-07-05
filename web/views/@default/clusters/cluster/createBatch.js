Tea.context(function () {
	this.success = NotifySuccess("保存成功", "/clusters/cluster/nodes?clusterId=" + this.clusterId)

	this.$delay(function () {
		this.$refs.ipList.focus()
	})
})