Vue.component("http-remote-addr-config-box", {
	props: ["v-remote-addr-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vRemoteAddrConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				value: "${remoteAddr}"
			}
		}
		return {
			config: config
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		}
	},
	template: `<div>
	<input type="hidden" name="remoteAddrJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示使用自定义的请求变量获取客户端IP。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>请求变量</td>
				<td>
					<input type="text" v-model="config.value" maxlength="100"/>
					<p class="comment">通过此变量获取用户的IP地址。具体可用的请求变量列表可参考官方网站文档。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>		
</div>`
})