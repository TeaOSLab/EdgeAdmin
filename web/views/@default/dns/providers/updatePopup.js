Tea.context(function () {
	this.typeDescription = ""

	let that = this
	this.types.forEach(function (v) {
		if (v.code == that.provider.type) {
			that.typeDescription = v.description
		}
	})
})