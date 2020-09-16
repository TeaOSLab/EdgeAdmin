Vue.component("size-capacity-box", {
	props: ["v-name", "v-value", "v-count", "v-unit"],
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
			"size": v,
			countString: (v.count >= 0) ? v.count.toString() : ""
		}
	},
	watch: {
		"countString": function (newValue) {
			let value = newValue.trim()
			if (value.length == 0) {
				this.size.count = -1
				return
			}
			let count = parseInt(value)
			if (!isNaN(count)) {
				this.size.count = count
			}
		}
	},
	template: `<div class="ui fields inline">
	<input type="hidden" :name="vName" :value="JSON.stringify(size)"/>
	<div class="ui field">
		<input type="text" v-model="countString" maxlength="11" size="11"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="size.unit">
			<option value="byte">字节</option>
			<option value="kb">KB</option>
			<option value="mb">MB</option>
			<option value="gb">GB</option>
		</select>
	</div>
</div>`
})