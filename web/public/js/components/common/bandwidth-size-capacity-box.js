Vue.component("bandwidth-size-capacity-box", {
	props: ["v-name", "v-value", "v-count", "v-unit", "size", "maxlength", "v-supported-units"],
	data: function () {
		let v = this.vValue
		if (v == null) {
			v = {
				count: this.vCount,
				unit: this.vUnit
			}
		}
		if (v.unit == null || v.unit.length == 0) {
			v.unit = "mb"
		}

		if (typeof (v["count"]) != "number") {
			v["count"] = -1
		}

		let vSize = this.size
		if (vSize == null) {
			vSize = 6
		}

		let vMaxlength = this.maxlength
		if (vMaxlength == null) {
			vMaxlength = 10
		}

		let supportedUnits = this.vSupportedUnits
		if (supportedUnits == null) {
			supportedUnits = []
		}

		return {
			capacity: v,
			countString: (v.count >= 0) ? v.count.toString() : "",
			vSize: vSize,
			vMaxlength: vMaxlength,
			supportedUnits: supportedUnits
		}
	},
	watch: {
		"countString": function (newValue) {
			let value = newValue.trim()
			if (value.length == 0) {
				this.capacity.count = -1
				this.change()
				return
			}
			let count = parseInt(value)
			if (!isNaN(count)) {
				this.capacity.count = count
			}
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.capacity)
		}
	},
	template: `<div class="ui fields inline">
	<input type="hidden" :name="vName" :value="JSON.stringify(capacity)"/>
	<div class="ui field">
		<input type="text" v-model="countString" :maxlength="vMaxlength" :size="vSize"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="capacity.unit" @change="change">
			<option value="b" v-if="supportedUnits.length == 0 || supportedUnits.$contains('b')">Bps</option>
			<option value="kb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('kb')">Kbps</option>
			<option value="mb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('mb')">Mbps</option>
			<option value="gb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('gb')">Gbps</option>
			<option value="tb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('tb')">Tbps</option>
			<option value="pb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('pb')">Pbps</option>
			<option value="eb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('eb')">Ebps</option>
		</select>
	</div>
</div>`
})