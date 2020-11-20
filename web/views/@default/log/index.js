Tea.context(function () {
	this.$delay(function () {
		teaweb.datepicker("day-from-picker")
		teaweb.datepicker("day-to-picker")
	})

	this.logs.forEach(function (v) {
		v.moreVisible = false
	})

	this.showMore = function (log) {
		log.moreVisible = !log.moreVisible
	}
})