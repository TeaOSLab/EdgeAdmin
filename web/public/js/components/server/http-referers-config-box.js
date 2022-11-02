Vue.component("http-referers-config-box", {
	props: ["v-referers-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vReferersConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				allowEmpty: true,
				allowSameDomain: true,
				allowDomains: [],
				denyDomains: []
			}
		}
		if (config.allowDomains == null) {
			config.allowDomains = []
		}
		if (config.denyDomains == null) {
			config.denyDomains = []
		}
		return {
			config: config
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeAllowDomains: function (domains) {
		},
		changeDenyDomains: function (domains) {
		}
	},
	template: `<div>
<input type="hidden" name="referersJSON" :value="JSON.stringify(config)"/>
<table class="ui table selectable definition">
	<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
	<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
		<tr>
			<td class="title">启用防盗链</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="config.isOn"/>
					<label></label>
				</div>
				<p class="comment">选中后表示开启防盗链。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="isOn()">
		<tr>
			<td class="title">允许直接访问网站</td>
			<td>
				<checkbox v-model="config.allowEmpty"></checkbox>
				<p class="comment">允许用户直接访问网站，用户第一次访问网站时来源域名通常为空。</p>
			</td>
		</tr>
		<tr>
			<td>来源域名允许一致</td>
			<td>
				<checkbox v-model="config.allowSameDomain"></checkbox>
				<p class="comment">允许来源域名和当前访问的域名一致，相当于在站内访问。</p>
			</td>
		</tr>
		<tr>
			<td>允许的来源域名</td>
			<td>
				<domains-box :v-domains="config.allowDomains" @change="changeAllowDomains">></domains-box>
				<p class="comment">允许的其他来源域名列表，比如<code-label>example.com</code-label>、<code-label>*.example.com</code-label>。单个星号<code-label>*</code-label>表示允许所有域名。</p>
			</td>
		</tr>
		<tr>
			<td>禁止的来源域名</td>
			<td>
				<domains-box :v-domains="config.denyDomains" @change="changeDenyDomains"></domains-box>
				<p class="comment">禁止的来源域名列表，比如<code-label>example.org</code-label>、<code-label>*.example.org</code-label>；除了这些禁止的来源域名外，其他域名都会被允许，除非限定了允许的来源域名。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="ui margin"></div>
</div>`
})