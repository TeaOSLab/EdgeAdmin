Vue.component("health-check-config-box", {
	props: ["v-health-check-config"],
	data: function () {
		let healthCheckConfig = this.vHealthCheckConfig
		let urlProtocol = "http"
		let urlPort = ""
		let urlRequestURI = "/"
		let urlHost = ""

		if (healthCheckConfig == null) {
			healthCheckConfig = {
				isOn: false,
				url: "",
				interval: {count: 60, unit: "second"},
				statusCodes: [200],
				timeout: {count: 10, unit: "second"},
				countTries: 3,
				tryDelay: {count: 100, unit: "ms"},
				autoDown: true,
				countUp: 1,
				countDown: 3
			}
			let that = this
			setTimeout(function () {
				that.changeURL()
			}, 500)
		} else {
			try {
				let url = new URL(healthCheckConfig.url)
				urlProtocol = url.protocol.substring(0, url.protocol.length - 1)

				// 域名
				urlHost = url.host
				if (urlHost == "%24%7Bhost%7D") {
					urlHost = "${host}"
				}
				let colonIndex = urlHost.indexOf(":")
				if (colonIndex > 0) {
					urlHost = urlHost.substring(0, colonIndex)
				}

				urlPort = url.port
				urlRequestURI = url.pathname
				if (url.search.length > 0) {
					urlRequestURI += url.search
				}
			} catch (e) {
			}

			if (healthCheckConfig.statusCodes == null) {
				healthCheckConfig.statusCodes = [200]
			}
			if (healthCheckConfig.interval == null) {
				healthCheckConfig.interval = {count: 60, unit: "second"}
			}
			if (healthCheckConfig.timeout == null) {
				healthCheckConfig.timeout = {count: 10, unit: "second"}
			}
			if (healthCheckConfig.tryDelay == null) {
				healthCheckConfig.tryDelay = {count: 100, unit: "ms"}
			}
			if (healthCheckConfig.countUp == null || healthCheckConfig.countUp < 1) {
				healthCheckConfig.countUp = 1
			}
			if (healthCheckConfig.countDown == null || healthCheckConfig.countDown < 1) {
				healthCheckConfig.countDown = 3
			}
		}
		return {
			healthCheck: healthCheckConfig,
			advancedVisible: false,
			urlProtocol: urlProtocol,
			urlHost: urlHost,
			urlPort: urlPort,
			urlRequestURI: urlRequestURI
		}
	},
	watch: {
		urlRequestURI: function () {
			if (this.urlRequestURI.length > 0 && this.urlRequestURI[0] != "/") {
				this.urlRequestURI = "/" + this.urlRequestURI
			}
			this.changeURL()
		},
		urlPort: function (v) {
			let port = parseInt(v)
			if (!isNaN(port)) {
				this.urlPort = port.toString()
			} else {
				this.urlPort = ""
			}
			this.changeURL()
		},
		urlProtocol: function () {
			this.changeURL()
		},
		urlHost: function () {
			this.changeURL()
		},
		"healthCheck.countTries": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countTries = count
			} else {
				this.healthCheck.countTries = 0
			}
		},
		"healthCheck.countUp": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countUp = count
			} else {
				this.healthCheck.countUp = 0
			}
		},
		"healthCheck.countDown": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countDown = count
			} else {
				this.healthCheck.countDown = 0
			}
		}
	},
	methods: {
		showAdvanced: function () {
			this.advancedVisible = !this.advancedVisible
		},
		changeURL: function () {
			let urlHost = this.urlHost
			if (urlHost.length == 0) {
				urlHost = "${host}"
			}
			this.healthCheck.url = this.urlProtocol + "://" + urlHost + ((this.urlPort.length > 0) ? ":" + this.urlPort : "") + this.urlRequestURI
		},
		changeStatus: function (values) {
			this.healthCheck.statusCodes = values.$map(function (k, v) {
				let status = parseInt(v)
				if (isNaN(status)) {
					return 0
				} else {
					return status
				}
			})
		}
	},
	template: `<div>
<input type="hidden" name="healthCheckJSON" :value="JSON.stringify(healthCheck)"/>
<table class="ui table definition selectable">
	<tbody>
		<tr>
			<td class="title">是否启用</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td>URL *</td>
			<td>
			    <table class="ui table">
			         <tr>
			            <td class="title">协议</td> 
			            <td>
			            	<select class="ui dropdown auto-width" v-model="urlProtocol">
							<option value="http">http://</option>
							<option value="https">https://</option>
						    </select>
                        </td>
                    </tr>
                    <tr>
                        <td>域名</td>
                        <td>
                            <input type="text" v-model="urlHost"/>
							<p class="comment">在此集群上可以访问到的一个域名。</p>
                        </td>
                    </tr>
                    <tr>
                        <td>端口</td>
                        <td>
                            <input type="text" maxlength="5" style="width:5.4em" placeholder="端口" v-model="urlPort"/>
                        </td>
                    </tr>
                    <tr>
                        <td>RequestURI</td>
                        <td><input type="text" v-model="urlRequestURI" placeholder="/" style="width:20em"/></td>
                    </tr>
                </table>
				<div class="ui divider"></div>
				<p class="comment" v-if="healthCheck.url.length > 0">拼接后的URL：<code-label>{{healthCheck.url}}</code-label>，其中\${host}指的是域名。</p>
			</td>
		</tr>
		<tr>
		    <td></td>
        </tr>
		<tr>
			<td>检测时间间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.interval"></time-duration-box>
			</td>
		</tr>
		<tr>
			<td>是否自动下线</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.autoDown"/>
					<label></label>
				</div>
				<p class="comment">选中后系统会根据健康检查的结果自动标记节点的上线/下线状态，并可能自动同步DNS设置。</p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续上线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countUp" style="width:5em" maxlength="6"/>
				<p class="comment">连续N次检查成功后自动恢复上线。</p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续下线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countDown" style="width:5em" maxlength="6"/>
				<p class="comment">连续N次检查失败后自动下线。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td colspan="2"><more-options-angle @change="showAdvanced"></more-options-angle></td>
		</tr>
	</tbody>
	<tbody v-show="advancedVisible && healthCheck.isOn">
		<tr>
			<td>允许的状态码</td>
			<td>
				<values-box :values="healthCheck.statusCodes" maxlength="3" @change="changeStatus"></values-box>
			</td>
		</tr>
		<tr>
			<td>超时时间</td>
			<td>
				<time-duration-box :v-value="healthCheck.timeout"></time-duration-box>
			</td>	
		</tr>
		<tr>
			<td>连续尝试次数</td>
			<td>
				<input type="text" v-model="healthCheck.countTries" style="width: 5em" maxlength="2"/>
			</td>
		</tr>
		<tr>
			<td>每次尝试间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.tryDelay"></time-duration-box>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`
})