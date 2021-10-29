Vue.component("datepicker", {
	props: ["v-name", "v-value", "v-bottom-left"],
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
			name = "day"
		}

		let day = this.vValue
		if (day == null) {
			day = ""
		}

		return {
			name: name,
			day: day
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.day)
		}
	},
	template: `<div style="display: inline-block">
	<input type="text" :name="name" v-model="day" placeholder="YYYY-MM-DD" style="width:8.6em" maxlength="10" @input="change" ref="dayInput" autocomplete="off"/>
</div>`
})