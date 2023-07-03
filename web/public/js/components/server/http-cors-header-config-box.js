Vue.component("http-cors-header-config-box", {
	props: ["value"],
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isOn: false,
				allowMethods: [],
				allowOrigin: "",
				allowCredentials: true,
				exposeHeaders: [],
				maxAge: 0,
				requestHeaders: [],
				requestMethod: "",
				optionsMethodOnly: false
			}
		}
		if (config.allowMethods == null) {
			config.allowMethods = []
		}
		if (config.exposeHeaders == null) {
			config.exposeHeaders = []
		}

		let maxAgeSecondsString = config.maxAge.toString()
		if (maxAgeSecondsString == "0") {
			maxAgeSecondsString = ""
		}

		return {
			config: config,

			maxAgeSecondsString: maxAgeSecondsString,

			moreOptionsVisible: false
		}
	},
	watch: {
		maxAgeSecondsString: function (v) {
			let seconds = parseInt(v)
			if (isNaN(seconds)) {
				seconds = 0
			}
			this.config.maxAge = seconds
		}
	},
	methods: {
		changeMoreOptions: function (visible) {
			this.moreOptionsVisible = visible
		},
		addDefaultAllowMethods: function () {
			let that = this
			let defaultMethods = ["PUT", "GET", "POST", "DELETE", "HEAD", "OPTIONS", "PATCH"]
			defaultMethods.forEach(function (method) {
				if (!that.config.allowMethods.$contains(method)) {
					that.config.allowMethods.push(method)
				}
			})
		}
	},
	template: `<div>
	<input type="hidden" name="corsJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">启用CORS自适应跨域</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
					<p class="comment">启用后，自动在响应报头中增加对应的<code-label>Access-Control-*</code-label>相关内容。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td colspan="2"><more-options-indicator @change="changeMoreOptions"></more-options-indicator></td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn && moreOptionsVisible">
			<tr>
				<td>允许的请求方法列表</td>
				<td>
					<http-methods-box :v-methods="config.allowMethods"></http-methods-box>
					<p class="comment"><a href="" @click.prevent="addDefaultAllowMethods">[添加默认]</a>。<code-label>Access-Control-Allow-Methods</code-label>值设置。所访问资源允许使用的方法列表，不设置则表示默认为<code-label>PUT</code-label>、<code-label>GET</code-label>、<code-label>POST</code-label>、<code-label>DELETE</code-label>、<code-label>HEAD</code-label>、<code-label>OPTIONS</code-label>、<code-label>PATCH</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>预检结果缓存时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 6em" maxlength="6" v-model="maxAgeSecondsString"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment"><code-label>Access-Control-Max-Age</code-label>值设置。预检结果缓存时间，0或者不填表示使用浏览器默认设置。注意每个浏览器有不同的缓存时间上限。</p>
				</td>
			</tr>
			<tr>
				<td>允许服务器暴露的报头</td>
				<td>
					<values-box :v-values="config.exposeHeaders"></values-box>
					<p class="comment"><code-label>Access-Control-Expose-Headers</code-label>值设置。允许服务器暴露的报头，请注意报头的大小写。</p>
				</td>
			</tr>
			<tr>
				<td>实际请求方法</td>
				<td>
					<input type="text" v-model="config.requestMethod"/>
					<p class="comment"><code-label>Access-Control-Request-Method</code-label>值设置。实际请求服务器时使用的方法，比如<code-label>POST</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>仅OPTIONS有效</td>
				<td>
					<checkbox v-model="config.optionsMethodOnly"></checkbox>
					<p class="comment">选中后，表示当前CORS设置仅在OPTIONS方法请求时有效。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})