Vue.component("http-remote-addr-config-box", {
	props: ["v-remote-addr-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vRemoteAddrConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				value: "${rawRemoteAddr}",
				type: "default",

				requestHeaderName: ""
			}
		}

		// type
		if (config.type == null || config.type.length == 0) {
			config.type = "default"
			switch (config.value) {
				case "${rawRemoteAddr}":
					config.type = "default"
					break
				case "${remoteAddrValue}":
					config.type = "default"
					break
				case "${remoteAddr}":
					config.type = "proxy"
					break
				default:
					if (config.value != null && config.value.length > 0) {
						config.type = "variable"
					}
			}
		}

		// value
		if (config.value == null || config.value.length == 0) {
			config.value = "${rawRemoteAddr}"
		}

		return {
			config: config,
			options: [
				{
					name: "直接获取",
					description: "用户直接访问边缘节点，即 \"用户 --> 边缘节点\" 模式，这时候系统会试图从直接的连接中读取到客户端IP地址。",
					value: "${rawRemoteAddr}",
					type: "default"
				},
				{
					name: "从上级代理中获取",
					description: "用户和边缘节点之间有别的代理服务转发，即 \"用户 --> [第三方代理服务] --> 边缘节点\"，这时候只能从上级代理中获取传递的IP地址；上级代理传递的请求报头中必须包含 X-Forwarded-For 或 X-Real-IP 信息。",
					value: "${remoteAddr}",
					type: "proxy"
				},
				{
					name: "从请求报头中读取",
					description: "从自定义请求报头读取客户端IP。",
					value: "",
					type: "requestHeader"
				},
				{
					name: "[自定义变量]",
					description: "通过自定义变量来获取客户端真实的IP地址。",
					value: "",
					type: "variable"
				}
			]
		}
	},
	watch: {
		"config.requestHeaderName": function (value) {
			if (this.config.type == "requestHeader"){
				this.config.value = "${header." + value.trim() + "}"
			}
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeOptionType: function () {
			let that = this

			switch(this.config.type) {
				case "default":
					this.config.value = "${rawRemoteAddr}"
					break
				case "proxy":
					this.config.value = "${remoteAddr}"
					break
				case "requestHeader":
					this.config.value = ""
					if (this.requestHeaderName != null && this.requestHeaderName.length > 0) {
						this.config.value = "${header." + this.requestHeaderName + "}"
					}

					setTimeout(function () {
						that.$refs.requestHeaderInput.focus()
					})
					break
				case "variable":
					this.config.value = "${rawRemoteAddr}"

					setTimeout(function () {
						that.$refs.variableInput.focus()
					})

					break
			}
		}
	},
	template: `<div>
	<input type="hidden" name="remoteAddrJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用访客IP设置</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后，表示使用自定义的请求变量获取客户端IP。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>获取IP方式 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="config.type" @change="changeOptionType">
						<option v-for="option in options" :value="option.type">{{option.name}}</option>
					</select>
					<p class="comment" v-for="option in options" v-if="option.type == config.type && option.description.length > 0">{{option.description}}</p>
				</td>
			</tr>
			
			<!-- read from request header -->
			<tr v-show="config.type == 'requestHeader'">
				<td>请求报头 *</td>
				<td>
					<input type="text" name="requestHeaderName" v-model="config.requestHeaderName" maxlength="100" ref="requestHeaderInput"/>
					<p class="comment">请输入包含有客户端IP的请求报头，需要注意大小写，常见的有<code-label>X-Forwarded-For</code-label>、<code-label>X-Real-IP</code-label>、<code-label>X-Client-IP</code-label>等。</p>
				</td>
			</tr>
			
			<!-- read from variable -->
			<tr v-show="config.type == 'variable'">
				<td>读取IP变量值 *</td>
				<td>
					<input type="text" name="value" v-model="config.value" maxlength="100" ref="variableInput"/>
					<p class="comment">通过此变量获取用户的IP地址。具体可用的请求变量列表可参考官方网站文档；比如通过报头传递IP的情形，可以使用<code-label>\${header.你的自定义报头}</code-label>（类似于<code-label>\${header.X-Forwarded-For}</code-label>，需要注意大小写规范）。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>		
</div>`
})