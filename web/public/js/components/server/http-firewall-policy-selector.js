Vue.component("http-firewall-policy-selector", {
	props: ["v-http-firewall-policy"],
	mounted: function () {
		let that = this
		Tea.action("/servers/components/waf/count")
			.post()
			.success(function (resp) {
				that.count = resp.data.count
			})
	},
	data: function () {
		let firewallPolicy = this.vHttpFirewallPolicy
		return {
			count: 0,
			firewallPolicy: firewallPolicy
		}
	},
	methods: {
		remove: function () {
			this.firewallPolicy = null
		},
		select: function () {
			let that = this
			teaweb.popup("/servers/components/waf/selectPopup", {
				height: "26em",
				callback: function (resp) {
					that.firewallPolicy = resp.data.firewallPolicy
				}
			})
		},
		create: function () {
			let that = this
			teaweb.popup("/servers/components/waf/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.firewallPolicy = resp.data.firewallPolicy
				}
			})
		}
	},
	template: `<div>
	<div v-if="firewallPolicy != null" class="ui label basic">
		<input type="hidden" name="httpFirewallPolicyId" :value="firewallPolicy.id"/>
		{{firewallPolicy.name}} &nbsp; <a :href="'/servers/components/waf/policy?firewallPolicyId=' + firewallPolicy.id" target="_blank" title="修改"><i class="icon pencil small"></i></a>&nbsp; <a href="" @click.prevent="remove()" title="删除"><i class="icon remove small"></i></a>
	</div>
	<div v-if="firewallPolicy == null">
		<span v-if="count > 0"><a href="" @click.prevent="select">[选择已有策略]</a> &nbsp; &nbsp; </span><a href="" @click.prevent="create">[创建新策略]</a>
	</div>
</div>`
})