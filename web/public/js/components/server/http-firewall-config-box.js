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
			moreOptionsVisible: false,
			execGlobalRules: !firewall.ignoreGlobalRules
		}
	},
	watch: {
		execGlobalRules: function (v) {
			this.firewall.ignoreGlobalRules = !v
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
					<p class="comment">选中后，表示启用当前网站服务的WAF功能。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeOptionsVisible"></more-options-tbody>
		<tbody v-show="moreOptionsVisible">
			<tr>
				<td>启用系统全局规则</td>
				<td>
					<checkbox v-model="execGlobalRules"></checkbox>
					<p class="comment">选中后，表示使用系统全局WAF策略中定义的规则。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})