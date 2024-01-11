Vue.component("bandwidth-size-capacity-view", {
	props: ["v-value"],
	data: function () {
		let capacity = this.vValue
		if (capacity != null && capacity.count > 0 && typeof capacity.unit === "string") {
			capacity.unit = capacity.unit[0].toUpperCase() + capacity.unit.substring(1) + "ps"
		}
		return {
			capacity: capacity
		}
	},
	template: `<span>
	<span v-if="capacity != null && capacity.count > 0">{{capacity.count}}{{capacity.unit}}</span>
</span>`
})