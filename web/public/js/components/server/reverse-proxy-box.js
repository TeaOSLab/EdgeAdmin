Vue.component("reverse-proxy-box", {
	props: ["v-reverse-proxy-ref", "v-reverse-proxy-config", "v-is-location"],
	data: function () {
		let reverseProxyRef = this.vReverseProxyRef
		if (reverseProxyRef == null) {
			reverseProxyRef = {
				isPrior: false,
				isOn: false,
				reverseProxyId: 0
			}
		}

		let reverseProxyConfig = this.vReverseProxyConfig
		if (reverseProxyConfig == null) {
			reverseProxyConfig = {
				requestPath: "",
				stripPrefix: "",
				requestURI: ""
			}
		}
		return {
			reverseProxyRef: reverseProxyRef,
			reverseProxyConfig: reverseProxyConfig,
			advancedVisible: false
		}
	},
	methods:  {
		isOn: function () {
			return (!this.vIsLocation || this.reverseProxyRef.isPrior) && this.reverseProxyRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		}
	},
	template: `<div>
	<input type="hidden" name="reverseProxyRefJSON" :value="JSON.stringify(reverseProxyRef)"/>
	<input type="hidden" name="reverseProxyJSON" :value="JSON.stringify(reverseProxyConfig)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="reverseProxyRef" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="!vIsLocation || reverseProxyRef.isPrior">
			<tr>
				<td class="title">是否启用反向代理</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td>请求主机名<em>（Host）</em></td>
				<td>
					<input type="text" placeholder="比如example.com" v-model="reverseProxyConfig.requestHost"/>
					<p class="comment">请求后端服务器时的Host，用于修改后端服务器接收到的域名，默认和客户端请求的主机名一致，通常不必填写，支持请求变量。</p>
				</td>
			</tr>
			<tr>
				<td>请求URI</td>
				<td>
					<input type="text" placeholder="\${requestURI}" v-model="reverseProxyConfig.requestURI"/>
					<p class="comment">\${requestURI}为完整的请求URI，可以使用类似于"\${requestURI}?arg1=value1&arg2=value2"的形式添加你的参数。</p>
				</td>
			</tr>
			<tr>
				<td>去除URL前缀</td>
				<td>
					<input type="text" v-model="reverseProxyConfig.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})