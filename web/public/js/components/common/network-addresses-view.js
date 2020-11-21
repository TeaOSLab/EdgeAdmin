Vue.component("network-addresses-view", {
	props: ["v-addresses"],
	template: `<div>
	<div class="ui label tiny basic" v-if="vAddresses != null" v-for="addr in vAddresses">
		{{addr.protocol}}://{{addr.host}}:{{addr.portRange}}
	</div>
</div>`
})