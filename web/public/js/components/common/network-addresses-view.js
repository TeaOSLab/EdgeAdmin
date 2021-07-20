Vue.component("network-addresses-view", {
	props: ["v-addresses"],
	template: `<div>
	<div class="ui label tiny basic" v-if="vAddresses != null" v-for="addr in vAddresses">
		{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host.quoteIP()}}</span><span v-else>*</span>:{{addr.portRange}}
	</div>
</div>`
})