Vue.component("time-duration-box", {
	props: ["v-name", "v-value", "v-count", "v-unit", "placeholder", "v-min-unit", "maxlength"],
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

		let minUnit = this.vMinUnit
		let units = [
			{
				code: "ms",
				name: "毫秒"
			},
			{
				code: "second",
				name: "秒"
			},
			{
				code: "minute",
				name: "分钟"
			},
			{
				code: "hour",
				name: "小时"
			},
			{
				code: "day",
				name: "天"
			}
		]
		let minUnitIndex = -1
		if (minUnit != null && typeof minUnit == "string" && minUnit.length > 0) {
			for (let i = 0; i < units.length; i++) {
				if (units[i].code == minUnit) {
					minUnitIndex = i
					break
				}
			}
		}
		if (minUnitIndex > -1) {
			units = units.slice(minUnitIndex)
		}

		let maxLength = parseInt(this.maxlength)
		if (typeof maxLength != "number") {
			maxLength = 10
		}

		return {
			duration: v,
			countString: (v.count >= 0) ? v.count.toString() : "",
			units: units,
			realMaxLength: maxLength
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
		<input type="text" v-model="countString" :maxlength="realMaxLength" :size="realMaxLength" :placeholder="placeholder" @keypress.enter.prevent="1"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="duration.unit" @change="change">
			<option v-for="unit in units" :value="unit.code">{{unit.name}}</option>
		</select>
	</div>
</div>`
})

Vue.component("time-duration-text", {
	props: ["v-value"],
	methods: {
		unitName: function (unit) {
			switch (unit) {
				case "ms":
					return "毫秒"
				case "second":
					return "秒"
				case "minute":
					return "分钟"
				case "hour":
					return "小时"
				case "day":
					return "天"
			}
		}
	},
	template: `<span>
	{{vValue.count}} {{unitName(vValue.unit)}}
</span>`
})