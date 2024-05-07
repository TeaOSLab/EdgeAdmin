Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功 ")

	this.policyType = this.cachePolicy.type

	this.fileOpenFileCacheMax = 0
	if (this.cachePolicy.type == "file" && this.cachePolicy.options.openFileCache != null && this.cachePolicy.options.openFileCache.isOn && this.cachePolicy.options.openFileCache.max > 0) {
		this.fileOpenFileCacheMax = this.cachePolicy.options.openFileCache.max
	}

	this.changePolicyType = function () {
		if (this.policyType == "file") {
			let options = this.cachePolicy.options
			if (options != null && typeof options == "object" && typeof options["dir"] === "undefined") {
				options["enableMMAP"] = false
				options["enableIncompletePartialContent"] = true
				options["dir"] = "/opt/cache"
				options["memoryPolicy"] = {
					capacity: {
						unit: "gb",
						count: 2
					}
				}
			}
		}
	}
})