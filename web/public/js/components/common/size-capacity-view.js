Vue.component("size-capacity-view", {
	props:["v-default-text", "v-value"],
	template: `<div>
	<span v-if="vValue != null && vValue.count > 0">{{vValue.count}}{{vValue.unit.toUpperCase().replace(/(.)B/, "$1iB")}}</span>
	<span v-else>{{vDefaultText}}</span>
</div>`
})