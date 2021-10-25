Vue.component("http-remote-addr-config-box", {
	props: ["v-remote-addr-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vRemoteAddrConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				value: "${rawRemoteAddr}",
				isCustomized: false
			}
		}

		let optionValue = ""
		if (!config.isCustomized && (config.value == "${remoteAddr}" || config.value == "${rawRemoteAddr}")) {
			optionValue = config.value
		}

		return {
			config: config,
			options: [
				{
					name: "直接获取",
					description: "用户直接访问边缘节点，即 \"用户 --> 边缘节点\" 模式，这时候可以直接从连接中读取到真实的IP地址。",
					value: "${rawRemoteAddr}"
				},
				{
					name: "从上级代理中获取",
					description: "用户和边缘节点之间有别的代理服务转发，即 \"用户 --> [第三方代理服务] --> 边缘节点\"，这时候只能从上级代理中获取传递的IP地址。",
					value: "${remoteAddr}"
				},
				{
					name: "[自定义]",
					description: "通过自定义变量来获取客户端真实的IP地址。",
					value: ""
				}
			],
			optionValue: optionValue
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeOptionValue: function () {
			if (this.optionValue.length > 0) {
				this.config.value = this.optionValue
				this.config.isCustomized = false
			} else {
				this.config.isCustomized = true
			}
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
				<td>获取IP方式 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="optionValue" @change="changeOptionValue">
						<option v-for="option in options" :value="option.value">{{option.name}}</option>
					</select>
					<p class="comment" v-for="option in options" v-if="option.value == optionValue && option.description.length > 0">{{option.description}}</p>
				</td>
			</tr>
			<tr v-show="optionValue.length == 0">
				<td>读取IP变量值 *</td>
				<td>
					<input type="hidden" v-model="config.value" maxlength="100"/>
					<div v-if="optionValue == ''" style="margin-top: 1em">
						<input type="text" v-model="config.value" maxlength="100"/>
						<p class="comment">通过此变量获取用户的IP地址。具体可用的请求变量列表可参考官方网站文档。</p>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>		
</div>`
})