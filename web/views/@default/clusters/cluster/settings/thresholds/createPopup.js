Tea.context(function () {
	this.success = NotifyPopup
	this.threshold = {
		item: this.items[0].code,
		param: "",
		operator: this.operators[0].code
	}
	this.$delay(function () {
		this.changeItem()
		this.changeParam()
	})

	this.itemDescription = ""
	this.itemParams = []

	this.changeItem = function () {
		let that = this
		this.threshold.param = ""
		this.items.forEach(function (v) {
			if (v.code == that.threshold.item) {
				that.itemDescription = v.description
				that.itemParams = v.params
				that.threshold.param = v.params[0].code
			}
		})
	}

	this.paramDescription = ""

	this.changeParam = function () {
		let that = this
		this.items.forEach(function (v) {
			if (v.code == that.threshold.item) {
				v.params.forEach(function (param) {
					if (param.code == that.threshold.param) {
						that.paramDescription = param.description
					}
				})
			}
		})
	}
})