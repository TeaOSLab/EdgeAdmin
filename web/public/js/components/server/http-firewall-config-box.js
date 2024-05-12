Vue.component("http-firewall-config-box", {
	props: ["v-firewall-config", "v-is-location", "v-is-group", "v-firewall-policy"],
	data: function () {
		let firewall = this.vFirewallConfig
		if (firewall == null) {
			firewall = {
				isPrior: false,
				isOn: false,
				firewallPolicyId: 0,
				ignoreGlobalRules: false,
				defaultCaptchaType: "none"
			}
		}

		if (firewall.defaultCaptchaType == null || firewall.defaultCaptchaType.length == 0) {
			firewall.defaultCaptchaType = "none"
		}

		let allCaptchaTypes = window.WAF_CAPTCHA_TYPES.$copy()

		// geetest
		let geeTestIsOn = false
		if (this.vFirewallPolicy != null && this.vFirewallPolicy.captchaAction != null && this.vFirewallPolicy.captchaAction.geeTestConfig != null) {
			geeTestIsOn = this.vFirewallPolicy.captchaAction.geeTestConfig.isOn
		}

		// 如果没有启用geetest，则还原
		if (!geeTestIsOn && firewall.defaultCaptchaType == "geetest") {
			firewall.defaultCaptchaType = "none"
		}

		return {
			firewall: firewall,
			moreOptionsVisible: false,
			execGlobalRules: !firewall.ignoreGlobalRules,
			captchaTypes: allCaptchaTypes,
			geeTestIsOn: geeTestIsOn
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
					<p class="comment">当前网站所在集群的设置。</p>
				</div>
				<span v-else class="red">当前集群没有设置WAF策略，当前配置无法生效。</span>
			</td>
		</tr>
	</table>
	
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="firewall" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || firewall.isPrior">
			<tr>
				<td class="title">启用Web防火墙</td>
				<td>
					<checkbox v-model="firewall.isOn"></checkbox>
					<p class="comment">选中后，表示启用当前网站的WAF功能。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeOptionsVisible" v-show="firewall.isOn"></more-options-tbody>
		<tbody v-show="moreOptionsVisible">
			<tr>
				<td>人机识别验证方式</td>
				<td>
					<select class="ui dropdown auto-width" v-model="firewall.defaultCaptchaType">
						<option value="none">默认</option>
						<option v-for="captchaType in captchaTypes" v-if="captchaType.code != 'geetest' || geeTestIsOn" :value="captchaType.code">{{captchaType.name}}</option>
					</select>
					<p class="comment" v-if="firewall.defaultCaptchaType == 'none'">使用系统默认的设置。你需要在入站规则中添加规则集来决定哪些请求需要人机识别验证。</p>
					<p class="comment" v-for="captchaType in captchaTypes" v-if="captchaType.code == firewall.defaultCaptchaType">{{captchaType.description}}你需要在入站规则中添加规则集来决定哪些请求需要人机识别验证。</p>
				</td>
			</tr>
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