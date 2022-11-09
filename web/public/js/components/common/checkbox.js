let checkboxId = 0
Vue.component("checkbox", {
	props: ["name", "value", "v-value", "id", "checked"],
	data: function () {
		checkboxId++
		let elementId = this.id
		if (elementId == null) {
			elementId = "checkbox" + checkboxId
		}

		let elementValue = this.vValue
		if (elementValue == null) {
			elementValue = "1"
		}

		let checkedValue = this.value
		if (checkedValue == null && this.checked == "checked") {
			checkedValue = elementValue
		}

		return {
			elementId: elementId,
			elementValue: elementValue,
			newValue: checkedValue
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.newValue)
		},
		check: function () {
			this.newValue = this.elementValue
		},
		uncheck: function () {
			this.newValue = ""
		},
		isChecked: function () {
			return (typeof (this.newValue) == "boolean" && this.newValue) || this.newValue == this.elementValue
		}
	},
	watch: {
		value: function (v) {
			if (typeof v == "boolean") {
				this.newValue = v
			}
		}
	},
	template: `<div class="ui checkbox">
	<input type="checkbox" :name="name" :value="elementValue" :id="elementId" @change="change" v-model="newValue"/>
	<label :for="elementId"><slot></slot></label>
</div>`
})