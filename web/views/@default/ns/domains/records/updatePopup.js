Tea.context(function () {
	this.type = this.record.type
	this.typeDescription = ""

	this.changeType = function () {
		let that = this
		this.types.forEach(function (v) {
			if (v.type == that.type) {
				that.typeDescription = v.description
			}
		})
	}

	this.changeType()
})