Vue.component("time-duration-box", {
	props: ["v-name", "v-value", "v-count", "v-unit"],
	mounted: function () {
		this.change()
	},
	data: function () {
		let v = this.vValue
		if (v == null) {
			v = {
				count: this.vCount,
				unit: this.vUnit
			}
		}
		if (typeof (v["count"]) != "number") {
			v["count"] = -1
		}
		return {
			duration: v,
			countString: (v.count >= 0) ? v.count.toString() : ""
		}
	},
	watch: {
		"countString": function (newValue) {
			let value = newValue.trim()
			if (value.length == 0) {
				this.duration.count = -1
				return
			}
			let count = parseInt(value)
			if (!isNaN(count)) {
				this.duration.count = count
			}
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.duration)
		}
	},
	template: `<div class="ui fields inline" style="padding-bottom: 0; margin-bottom: 0">
	<input type="hidden" :name="vName" :value="JSON.stringify(duration)"/>
	<div class="ui field">
		<input type="text" v-model="countString" maxlength="11" size="11" @keypress.enter.prevent="1"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="duration.unit" @change="change">
			<option value="ms">毫秒</option>
			<option value="second">秒</option>
			<option value="minute">分钟</option>
			<option value="hour">小时</option>
			<option value="day">天</option>
		</select>
	</div>
</div>`
})