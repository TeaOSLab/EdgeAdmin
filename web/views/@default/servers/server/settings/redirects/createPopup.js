Tea.context(function () {
	this.isCreating = true
	if (window.parent.UPDATING_REDIRECT != null) {
		this.isCreating = false
		this.redirect = window.parent.UPDATING_REDIRECT
		if (this.redirect.type == null || this.redirect.type.length == 0) {
			this.redirect.type = "url"
		}
	} else {
		this.redirect = {
			type: "url",
			status: 0,
			beforeURL: "",
			afterURL: "",
			matchPrefix: false,
			matchRegexp: false,
			keepRequestURI: false,
			keepArgs: true,
			conds: null,
			isOn: true,

			domainsAll: false,
			domainBefore: [],
			domainBeforeIgnorePorts: true,
			domainAfter: "",
			domainAfterScheme: "",

			portsAll: false,
			portsBefore: [],
			portAfter: 0,
			portAfterScheme: ""
		}
	}

	this.mode = ""
	if (this.redirect.matchPrefix) {
		this.mode = "matchPrefix"
	} else if (this.redirect.matchRegexp) {
		this.mode = "matchRegexp"
	} else {
		this.mode = "equal"
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