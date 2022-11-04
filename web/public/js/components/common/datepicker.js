Vue.component("datepicker", {
	props: ["value", "v-name", "name", "v-value", "v-bottom-left", "placeholder"],
	mounted: function () {
		let that = this
		teaweb.datepicker(this.$refs.dayInput, function (v) {
			that.day = v
			that.change()
		}, !!this.vBottomLeft)
	},
	data: function () {
		let name = this.vName
		if (name == null) {
			name = this.name
		}
		if (name == null) {
			name = "day"
		}

		let day = this.vValue
		if (day == null) {
			day = this.value
			if (day == null) {
				day = ""
			}
		}

		let placeholder = "YYYY-MM-DD"
		if (this.placeholder != null) {
			placeholder = this.placeholder
		}

		return {
			realName: name,
			realPlaceholder: placeholder,
			day: day
		}
	},
	watch: {
		value: function (v) {
			this.day = v

			let picker = this.$refs.dayInput.picker
			if (picker != null) {
				if (v != null && /^\d+-\d+-\d+$/.test(v)) {
					picker.setDate(v)
				}
			}
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.day) // support v-model，事件触发需要在 change 之前
			this.$emit("change", this.day)
		}
	},
	template: `<div style="display: inline-block">
	<input type="text" :name="realName" v-model="day" :placeholder="realPlaceholder" style="width:8.6em" maxlength="10" @input="change" ref="dayInput" autocomplete="off"/>
</div>`
})