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
			this.loadCheckpoint()
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

	this.loadCheckpoint = function () {
		if (this.rule.checkpointPrefix.length == 0) {
			this.checkpoint = null
			return
		}
		let that = this
		this.checkpoint = this.checkpoints.$find(function (k, v) {
			return v.prefix == that.rule.checkpointPrefix
		})
	}

	this.changeCheckpoint = function () {
		if (this.rule.checkpointPrefix.length == 0) {
			this.checkpoint = null
			return
		}
		let that = this
		this.checkpoint = this.checkpoints.$find(function (k, v) {
			return v.prefix == that.rule.checkpointPrefix
		})
		if (this.checkpoint == null) {
			return
		}
		switch (this.checkpoint.dataType) {
			case "bool":
				this.rule.operator = "eq"
				break
			case "number":
				this.rule.operator = "eq"
				break
			default:
				this.rule.operator = "match"
		}
	}


	/**
	 * operator
	 */
	this.changeOperator = function () {
		let that = this;
		this.operator = this.operators.$find(function (k, v) {
			return v.code == that.rule.operator
		})
		if (this.operator == null) {
			return
		}
		if (!this.isUpdating) {
			this.rule.isCaseInsensitive = (this.operator.case == "yes")
		}
	};
	this.changeOperator()

	/**
	 * caseInsensitive
	 */
	this.changeCaseInsensitive = function () {
		if (this.rule.operator == "match" || this.rule.operator == "not match") {
			if (this.regexpTestIsOn) {
				this.changeRegexpTestBody()
			}
		}
	}

	/**
	 * value
	 */
	this.changeRuleValue = function () {
		if (this.rule.operator == "match" || this.rule.operator == "not match") {
			if (this.regexpTestIsOn) {
				this.changeRegexpTestBody()
			}
		} else {
			this.regexpTestIsOn = false
			this.regexpTestResult = {isOk: false, message: ""}
		}
	}

	this.convertValueLine = function () {
		let value = this.rule.value
		if (value != null && value.length > 0) {
			let lines = value.split(/\n/)
			let resultLines = []
			lines.forEach(function (line) {
				line = line.trim()
				if (line.length > 0) {
					resultLines.push(line)
				}
			})
			this.rule.value = resultLines.join("|")
		}
	}

	/**
	 * 正则测试
	 */
	this.regexpTestIsOn = false
	this.regexpTestBody = ""
	this.regexpTestResult = {isOk: false, message: ""}

	this.changeRegexpTestIsOn = function () {
		this.regexpTestIsOn = !this.regexpTestIsOn
		if (this.regexpTestIsOn) {
			this.$delay(function () {
				this.$refs.regexpTestBody.focus()
			})
		}
	}

	this.changeRegexpTestBody = function () {
		this.$post(".testRegexp")
			.params({
				"regexp": this.rule.value,
				"body": this.regexpTestBody,
				"isCaseInsensitive": this.rule.isCaseInsensitive
			})
			.success(function (resp) {
				this.regexpTestResult = resp.data.result
			})
	}

	// isp
	this.selectISPName = function (isp) {
		if (isp == null) {
			return
		}

		let ispName = isp.name
		this.$refs.ispComboBox.clear()

		if (this.rule.value.length == 0) {
			this.rule.value = ispName
		} else {
			this.rule.value += "|" + ispName
		}
	}

	// country
	this.selectGeoCountryName = function (country) {
		if (country == null) {
			return
		}

		let countryName = country.name
		this.$refs.countryComboBox.clear()

		if (this.rule.value.length == 0) {
			this.rule.value = countryName
		} else {
			this.rule.value += "|" + countryName
		}
	}

	// province
	this.selectGeoProvinceName = function (province) {
		if (province == null) {
			return
		}

		let provinceName = province.name
		this.$refs.provinceComboBox.clear()

		if (this.rule.value.length == 0) {
			this.rule.value = provinceName
		} else {
			this.rule.value += "|" + provinceName
		}
	}

	// city
	this.selectGeoCityName = function (city) {
		if (city == null) {
			return
		}

		let cityName = city.name
		this.$refs.cityComboBox.clear()

		if (this.rule.value.length == 0) {
			this.rule.value = cityName
		} else {
			this.rule.value += "|" + cityName
		}
	}
})