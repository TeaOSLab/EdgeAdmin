Tea.context(function () {
	this.type = this.types[0].code
	this.typeDefinition = null

	this.$delay(function () {
		this.changeType()
	})

	this.changeType = function () {
		let that = this
		this.typeDefinition = this.types.$find(function (k, v) {
			return v.code == that.type
		})
	}
})