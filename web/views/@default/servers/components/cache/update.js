Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功 ")

	this.policyType = this.cachePolicy.type

	this.fileOpenFileCacheMax = 0
	if (this.cachePolicy.type == "file" && this.cachePolicy.options.openFileCache != null && this.cachePolicy.options.openFileCache.isOn && this.cachePolicy.options.openFileCache.max > 0) {
		this.fileOpenFileCacheMax = this.cachePolicy.options.openFileCache.max
	}
})