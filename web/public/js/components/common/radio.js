let radioId = 0
Vue.component("radio", {
	props: ["name", "value", "v-value", "id"],
	data: function () {
		radioId++
		let elementId = this.id
		if (elementId == null) {
			elementId = "radio" + radioId
		}
		return {
			"elementId": elementId
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.vValue)
		}
	},
	template: `<div class="ui checkbox radio">
	<input type="radio" :name="name" :value="vValue" :id="elementId" @change="change" :checked="(vValue == value)"/>
	<label :for="elementId"><slot></slot></label>
</div>`
})