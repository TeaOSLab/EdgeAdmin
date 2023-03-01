Vue.component("http-firewall-config-box", {
	props: ["v-firewall-config", "v-is-location", "v-is-group", "v-firewall-policy"],
	data: function () {
		let firewall = this.vFirewallConfig
		if (firewall == null) {
			firewall = {
				isPrior: false,
				isOn: false,
				firewallPolicyId: 0,
				ignoreGlobalRules: false
			}
		}

		return {
			firewall: firewall,
			moreOptionsVisible: false
		}
	},
	methods: {
		changeOptionsVisible: function (v) {
			this.moreOptionsVisible = v
		}
	},
	template: `<div>
	<input type="hidden" name="firewallJSON" :value="JSON.stringify(firewall)"/>
	
	<table class="ui table selectable definition" v-show="!vIsGroup">
		<tr>
			<td class="title">全局WAF策略</td>
			<td>
				<div v-if="vFirewallPolicy != null">{{vFirewallPolicy.name}} <span v-if="vFirewallPolicy.modeInfo != null">&nbsp; <span :class="{green: vFirewallPolicy.modeInfo.code == 'defend', blue: vFirewallPolicy.modeInfo.code == 'observe', grey: vFirewallPolicy.modeInfo.code == 'bypass'}">[{{vFirewallPolicy.modeInfo.name}}]</span>&nbsp;</span> <link-icon :href="'/servers/components/waf/policy?firewallPolicyId=' + vFirewallPolicy.id"></link-icon>
					<p class="comment">当前服务所在集群的设置。</p>
				</div>
				<span v-else class="red">当前集群没有设置WAF策略，当前配置无法生效。</span>
			</td>
		</tr>
	</table>
	
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="firewall" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || firewall.isPrior">
			<tr>
				<td class="title">启用WAF</td>
				<td>
					<checkbox v-model="firewall.isOn"></checkbox>
					<p class="comment">启用WAF之后，各项WAF设置才会生效。</p>
				</td>
			</tr>
		</tbody>
		<tr>
			<td colspan="2"><more-options-indicator @change="changeOptionsVisible"></more-options-indicator></td>
		</tr>
		<tbody v-show="moreOptionsVisible">
			<tr>
				<td>不使用全局规则</td>
				<td>
					<checkbox v-model="firewall.ignoreGlobalRules"></checkbox>
					<p class="comment">选中后，表示<strong>不使用</strong>系统全局WAF策略中定义的规则。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})