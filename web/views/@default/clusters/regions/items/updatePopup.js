Tea.context(function () {
	this.bitsFrom = this.item.bitsFrom / 1000 / 1000
	this.bitsFromMB = ""

	this.bitsTo = this.item.bitsTo / 1000 / 1000
	this.bitsToMB = ""

	this.$delay(function () {
		let that = this
		this.$watch("bitsFrom", function (v) {
			this.bitsFromMB = that.formatBits(v)
		})
		this.$watch("bitsTo", function (v) {
			this.bitsToMB = that.formatBits(v)
		})

		this.bitsFromMB = that.formatBits(this.bitsFrom)
		this.bitsToMB = that.formatBits(this.bitsTo)
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