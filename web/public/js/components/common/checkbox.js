let checkboxId = 0
Vue.component("checkbox", {
	props: ["name", "value", "v-value", "id"],
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

		return {
			elementId: elementId,
			elementValue: elementValue,
			newValue: this.value
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.newValue)
		}
	},
	template: `<div class="ui checkbox">
	<input type="checkbox" :name="name" :value="elementValue" :id="elementId" @change="change" v-model="newValue"/>
	<label :for="elementId"><slot></slot></label>
</div>`
})