Tea.context(function () {
	this.bitsFrom = 0
	this.bitsFromMB = ""

	this.bitsTo = 0
	this.bitsToMB = ""

	this.$delay(function () {
		let that = this
		this.$watch("bitsFrom", function (v) {
			this.bitsFromMB = that.formatBits(v)
		})
		this.$watch("bitsTo", function (v) {
			this.bitsToMB = that.formatBits(v)
		})
	})

	this.formatBits = function (bits) {
		bits = parseInt(bits)
		if (isNaN(bits)) {
			bits = 0
		}

		if (bits < 1000) {
			return bits + "MB"
		}

		if (bits < 1000 * 1000) {
			return (bits / 1000) + "GB"
		}

		if (bits < 1000 * 1000 * 1000) {
			return (bits / 1000 / 1000) + "TB"
		}

		if (bits < 1000 * 1000 * 1000 * 1000) {
			return (bits / 1000 / 1000 / 1000) + "PB"
		}

		if (bits < 1000 * 1000 * 1000 * 1000 * 1000) {
			return (bits / 1000 / 1000 / 1000 / 1000) + "EB"
		}

		return ""
	}
})