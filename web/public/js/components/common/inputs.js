Vue.component("digit-input", {
	props: ["value", "maxlength", "size", "min", "max", "required", "placeholder"],
	mounted: function () {
		let that = this
		setTimeout(function () {
			that.check()
		})
	},
	data: function () {
		let realMaxLength = this.maxlength
		if (realMaxLength == null) {
			realMaxLength = 20
		}

		let realSize = this.size
		if (realSize == null) {
			realSize = 6
		}

		return {
			realValue: this.value,
			realMaxLength: realMaxLength,
			realSize: realSize,
			isValid: true
		}
	},
	watch: {
		realValue: function (v) {
			this.notifyChange()
		}
	},
	methods: {
		notifyChange: function () {
			let v = parseInt(this.realValue.toString(), 10)
			if (isNaN(v)) {
				v = 0
			}
			this.check()
			this.$emit("input", v)
		},
		check: function () {
			if (this.realValue == null) {
				return
			}
			let s = this.realValue.toString()
			if (!/^\d+$/.test(s)) {
				this.isValid = false
				return
			}
			let v = parseInt(s, 10)
			if (isNaN(v)) {
				this.isValid = false
			} else {
				if (this.required) {
					this.isValid = (this.min == null || this.min <= v) && (this.max == null || this.max >= v)
				} else {
					this.isValid = (v == 0 || (this.min == null || this.min <= v) && (this.max == null || this.max >= v))
				}
			}
		}
	},
	template: `<input type="text" v-model="realValue" :maxlength="realMaxLength" :size="realSize" :class="{error: !this.isValid}" :placeholder="placeholder" autocomplete="off"/>`
})