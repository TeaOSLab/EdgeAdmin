Tea.context(function () {
	this.$delay(function () {
		let that = this
		teaweb.datepicker("day-input", function (day) {
			that.day = day
		})
	})
})