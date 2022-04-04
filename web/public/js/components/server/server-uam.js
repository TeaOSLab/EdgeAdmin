// UAM模式配置
Vue.component("uam-config-box", {
	props: ["v-uam-config"],
	data: function () {
		let config = this.vUamConfig
		if (config == null) {
			config = {
				isOn: false
			}
		}
		return {
			config: config
		}
	},
	template: `<div>
<input type="hidden" name="uamJSON" :value="JSON.stringify(config)"/>
<table class="ui table definition selectable">
	<tr>
		<td class="title">启用5秒盾</td>
		<td>
			<checkbox v-model="config.isOn"></checkbox>
			<p class="comment"><plus-label></plus-label>启用后，访问网站时，自动检查浏览器环境，阻止非正常访问。</p>
		</td>
	</tr>
</table>
<div class="margin"></div>
</div>`
})