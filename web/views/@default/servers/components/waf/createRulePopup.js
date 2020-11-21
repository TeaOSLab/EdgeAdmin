Tea.context(function () {
	this.success = NotifyPopup

	this.isUpdating = (window.parent.UPDATING_RULE != null)
	this.rule = {
		id: 0,
		param: "",
		paramFilters: [],
		checkpointPrefix: "",
		checkpointParam: "",
		value: "",
		isCaseInsensitive: false,
		operator: "match",
		checkpointOptions: null,
		description: "",
		isOn: true
	}
	if (window.parent.UPDATING_RULE != null) {
		this.rule = window.parent.UPDATING_RULE

		let param = this.rule.param.substring(this.rule.param.indexOf("${") + 2, this.rule.param.indexOf("}"))
		let index = param.indexOf(".")
		if (index > 0) {
			this.rule.checkpointPrefix = param.substring(0, index)
			this.rule.checkpointParam = param.substring(index + 1)
		} else {
			this.rule.checkpointPrefix = param
		}
		this.$delay(function () {
			this.changeCheckpoint()
			if (this.rule.checkpointOptions != null && this.checkpoint != null && this.checkpoint.options != null) {
				let that = this
				this.checkpoint.options.forEach(function (option) {
					if (typeof (that.rule.checkpointOptions[option.code]) != "undefined") {
						option.value = that.rule.checkpointOptions[option.code]
					}
				})
			}
		})
	}

	/**
	 * checkpoint
	 */
	this.checkpoint = null
	this.changeCheckpoint = function () {
		if (this.rule.checkpointPrefix.length == 0) {
			this.checkpoint = null
			return
		}
		let that = this
		this.checkpoint = this.checkpoints.$find(function (k, v) {
			return v.prefix == that.rule.checkpointPrefix
		})
	}


	/**
	 * operator
	 */
	this.changeOperator = function () {
		let that = this;
		this.operator = this.operators.$find(function (k, v) {
			return v.code == that.rule.operator
		})
		if (!this.isUpdating) {
			this.rule.isCaseInsensitive = (this.operator.case == "yes")
		}
	};
	this.changeOperator()
})