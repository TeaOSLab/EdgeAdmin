Tea.context(function () {
	var commonGroupName = "常用"
	this.letterGroups = [
		{
			"code": commonGroupName,
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "ABC",
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "DEF",
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "GHI",
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "JKL",
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "MNO",
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "PQR",
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "STU",
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "VWX",
			"count": 0,
			"countMatched": 0
		},
		{
			"code": "YZ",
			"count": 0,
			"countMatched": 0
		}
	]
	this.commonGroupName = commonGroupName
	this.selectedGroup = commonGroupName
	this.letterCountries = {}
	let that = this
	this.countSelectedCountries = this.countries.$count(function (k, country) {
		return country.isChecked
	})
	this.countries.forEach(function (country) {
		// letter
		if (typeof (that.letterCountries[country.letter]) == "undefined") {
			that.letterCountries[country.letter] = []
		}
		that.letterCountries[country.letter].push(country)

		// common
		if (country.isCommon) {
			if (typeof that.letterCountries[commonGroupName] == "undefined") {
				that.letterCountries[commonGroupName] = []
			}
			that.letterCountries[commonGroupName].push(country)
		}
	})
	this.isCheckingAll = false

	this.$delay(function () {
		this.change()
	})

	this.checkAll = function () {
		this.isCheckingAll = !this.isCheckingAll

		this.countries.forEach(function (country) {
			country.isChecked = that.isCheckingAll
		})

		this.change()
	}

	this.selectGroup = function (group) {
		this.selectedGroup = group.code
	}

	this.selectCountry = function (country) {
		country.isChecked = !country.isChecked
		this.change()
	}

	this.deselectCountry = function (country) {
		country.isChecked = false
		this.change()
	}

	this.change = function () {
		let that = this
		this.letterGroups.forEach(function (group) {
			group.count = 0
			group.countMatched = 0
		})
		this.countries.forEach(function (country) {
			that.letterGroups.forEach(function (group) {
				if (group.code.indexOf(country.letter) >= 0 || (group.code == commonGroupName && country.isCommon)) {
					if (country.isChecked) {
						group.count++
					}
					if (that.matchCountry(country)) {
						country.isMatched = (that.keyword.length > 0)
						group.countMatched++
					} else {
						country.isMatched = false
					}
				}
			})
		})
	}

	this.submit = function () {
		let selectedCountries = []
		this.countries.forEach(function (country) {
			if (country.isChecked) {
				selectedCountries.push(country)
			}
		})
		NotifyPopup({
			code: 200,
			data: {
				selectedCountries: selectedCountries
			}
		})
	}

	/**
	 * searching
	 */
	this.searchBoxVisible = false
	this.keyword = ""

	this.showSearchBox = function () {
		this.searchBoxVisible = true
		this.$delay(function () {
			this.$refs.searchBox.focus()
		})
	}

	this.changeKeyword = function (event) {
		this.keyword = event.value.trim()
		this.change()
	}

	this.matchCountry = function (country) {
		if (this.keyword.length == 0) {
			return true
		}

		if (teaweb.match(country.name, this.keyword)) {
			return true
		}
		if (country.pinyin != null && country.pinyin.length > 0) {
			let matched = false
			let that = this
			country.pinyin.forEach(function (code) {
				if (teaweb.match(code, that.keyword)) {
					matched = true
				}
			})
			if (matched) {
				return true
			}
		}
		if (country.codes != null && country.codes.length > 0) {
			let matched = false
			let that = this
			country.codes.forEach(function (code) {
				if (teaweb.match(code, that.keyword)) {
					matched = true
				}
			})
			if (matched) {
				return true
			}
		}
		return false
	}
})