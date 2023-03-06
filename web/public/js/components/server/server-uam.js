// UAM模式配置
Vue.component("uam-config-box", {
	props: ["v-uam-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vUamConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				onlyURLPatterns: [],
				exceptURLPatterns: []
			}
		}
		if (config.onlyURLPatterns == null) {
			config.onlyURLPatterns = []
		}
		if (config.exceptURLPatterns == null) {
			config.exceptURLPatterns = []
		}
		return {
			config: config,
			moreOptionsVisible: false
		}
	},
	methods: {
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		}
	},
	template: `<div>
<input type="hidden" name="uamJSON" :value="JSON.stringify(config)"/>
<table class="ui table definition selectable">
	<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
	<tbody v-show="((!vIsLocation && !vIsGroup) || config.isPrior)">
		<tr>
			<td class="title">启用5秒盾</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment"><plus-label></plus-label>启用后，访问网站时，自动检查浏览器环境，阻止非正常访问。</p>
			</td>
		</tr>
	</tbody>
	<tbody>
		<tr>
			<td colspan="2"><more-options-indicator @change="showMoreOptions"></more-options-indicator></td>
		</tr>
	</tbody>
	<tbody v-show="moreOptionsVisible">
		<tr>
			<td>例外URL</td>
			<td>
				<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
				<p class="comment">如果填写了例外URL，表示这些URL跳过5秒盾不做处理。</p>
			</td>
		</tr>
		<tr>
			<td>限制URL</td>
			<td>
				<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
				<p class="comment">如果填写了支持URL，表示只对这些URL进行5秒盾处理；如果不填则表示支持所有的URL。</p>
			</td>
		</tr>	
	</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`
})