Tea.context(function () {
	this.secretType = "clear"
	this.secret = ""

	this.$delay(function () {
		this.$watch("secretType", function () {
			this.secret = ""
		})
	})

	this.generateSecret = function () {
		this.$post(".generateSecret")
			.params({
				secretType: this.secretType
			})
			.success(function (resp) {
				this.secret = resp.data.secret
			})
	}
})