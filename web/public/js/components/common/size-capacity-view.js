Vue.component("size-capacity-view", {
	props:["v-default-text", "v-value"],
	methods: {
		composeCapacity: function (capacity) {
			return teaweb.convertSizeCapacityToString(capacity)
		}
	},
	template: `<div>
	<span v-if="vValue != null && vValue.count > 0">{{composeCapacity(vValue)}}</span>
	<span v-else>{{vDefaultText}}</span>
</div>`
})