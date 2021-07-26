Tea.context(function () {
	this.isCreating = true
	if (window.parent.UPDATING_REDIRECT != null) {
		this.isCreating = false
		this.redirect = window.parent.UPDATING_REDIRECT
	} else {
		this.redirect = {
			status: 0,
			beforeURL: "",
			afterURL: "",
			matchPrefix: false,
			matchRegexp: false,
			keepRequestURI: false,
			conds: null,
			isOn: true
		}
	}

	this.mode = ""
	if (this.redirect.matchPrefix) {
		this.mode = "matchPrefix"
	} else if (this.redirect.matchRegexp) {
		this.mode = "matchRegexp"
	} else {
		this.mode = "matchPrefix"
		this.redirect.matchPrefix = true
	}

	this.$delay(function () {
		let that = this
		this.$watch("mode", function (v) {
			if (v == "matchPrefix") {
				that.redirect.matchPrefix = true
				that.redirect.matchRegexp = false
			} else if (v == "matchRegexp") {
				that.redirect.matchPrefix = false
				that.redirect.matchRegexp = true
			} else {
				that.redirect.matchPrefix = false
				that.redirect.matchRegexp = false
			}
		})
	})

	this.changeConds = function (conds) {
		this.redirect.conds = conds
	}
})