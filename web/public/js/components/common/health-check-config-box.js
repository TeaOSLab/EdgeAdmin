Vue.component("health-check-config-box", {
	props: ["v-health-check-config", "v-check-domain-url", "v-is-plus"],
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
				countDown: 3,
				userAgent: "",
				onlyBasicRequest: true,
				accessLogIsOn: true
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
			urlRequestURI: urlRequestURI,
			urlIsEditing: healthCheckConfig.url.length == 0,

			hostErr: ""
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
			this.hostErr = ""
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
		},
		onChangeURLHost: function () {
			let checkDomainURL = this.vCheckDomainUrl
			if (checkDomainURL == null || checkDomainURL.length == 0) {
				return
			}

			let that = this
			Tea.action(checkDomainURL)
				.params({host: this.urlHost})
				.success(function (resp) {
					if (!resp.data.isOk) {
						that.hostErr = "在当前集群中找不到此域名，可能会影响健康检查结果。"
					} else {
						that.hostErr = ""
					}
				})
				.post()
		},
		editURL: function () {
			this.urlIsEditing = !this.urlIsEditing
		}
	},
	template: `<div>
<input type="hidden" name="healthCheckJSON" :value="JSON.stringify(healthCheck)"/>
<table class="ui table definition selectable">
	<tbody>
		<tr>
			<td class="title">启用健康检查</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.isOn"/>
					<label></label>
				</div>
				<p class="comment">通过访问节点上的网站URL来确定节点是否健康。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td>检测URL *</td>
			<td>
				<div v-if="healthCheck.url.length > 0" style="margin-bottom: 1em"><code-label>{{healthCheck.url}}</code-label> &nbsp; <a href="" @click.prevent="editURL"><span class="small">修改 <i class="icon angle" :class="{down: !urlIsEditing, up: urlIsEditing}"></i></span></a> </div>
				<div v-show="urlIsEditing">
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
								<input type="text" v-model="urlHost" @change="onChangeURLHost"/>
								<p class="comment"><span v-if="hostErr.length > 0" class="red">{{hostErr}}</span>已经部署到当前集群的一个域名；如果为空则使用节点IP作为域名。<span class="red" v-if="urlProtocol == 'https' && urlHost.length == 0">如果协议是https，这里必须填写一个已经设置了SSL证书的域名。</span></p>
							</td>
						</tr>
						<tr>
							<td>端口</td>
							<td>
								<input type="text" maxlength="5" style="width:5.4em" placeholder="端口" v-model="urlPort"/>
								<p class="comment">域名或者IP的端口，可选项，默认为80/443。</p>
							</td>
						</tr>
						<tr>
							<td>RequestURI</td>
							<td><input type="text" v-model="urlRequestURI" placeholder="/" style="width:20em"/>
								<p class="comment">请求的路径，可以带参数，可选项。</p>
							</td>
						</tr>
					</table>
					<div class="ui divider"></div>
					<p class="comment" v-if="healthCheck.url.length > 0">拼接后的检测URL：<code-label>{{healthCheck.url}}</code-label>，其中\${host}指的是域名。</p>
				</div>
			</td>
		</tr>
		<tr>
			<td>检测时间间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.interval"></time-duration-box>
				<p class="comment">两次检查之间的间隔。</p>
			</td>
		</tr>
		<tr>
			<td>自动上/下线<span v-if="vIsPlus">IP</span></td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.autoDown"/>
					<label></label>
				</div>
				<p class="comment">选中后系统会根据健康检查的结果自动标记<span v-if="vIsPlus">节点IP</span><span v-else>节点</span>的上线/下线状态，并可能自动同步DNS设置。<span v-if="!vIsPlus">注意：免费版的只能整体上下线整个节点，商业版的可以下线单个IP。</span></p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续上线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countUp" style="width:5em" maxlength="6"/>
				<p class="comment">连续{{healthCheck.countUp}}次检查成功后自动恢复上线。</p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续下线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countDown" style="width:5em" maxlength="6"/>
				<p class="comment">连续{{healthCheck.countDown}}次检查失败后自动下线。</p>
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
				<p class="comment">允许检测URL返回的状态码列表。</p>
			</td>
		</tr>
		<tr>
			<td>超时时间</td>
			<td>
				<time-duration-box :v-value="healthCheck.timeout"></time-duration-box>
				<p class="comment">读取检测URL超时时间。</p>
			</td>	
		</tr>
		<tr>
			<td>连续尝试次数</td>
			<td>
				<input type="text" v-model="healthCheck.countTries" style="width: 5em" maxlength="2"/>
				<p class="comment">如果读取检测URL失败后需要再次尝试的次数。</p>
			</td>
		</tr>
		<tr>
			<td>每次尝试间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.tryDelay"></time-duration-box>
				<p class="comment">如果读取检测URL失败后再次尝试时的间隔时间。</p>
			</td>
		</tr>
		<tr>
			<td>终端信息<em>（User-Agent）</em></td>
			<td>
				<input type="text" v-model="healthCheck.userAgent" maxlength="200"/>
				<p class="comment">发送到服务器的User-Agent值，不填写表示使用默认值。</p>
			</td>
		</tr>
		<tr>
			<td>只基础请求</td>
			<td>
				<checkbox v-model="healthCheck.onlyBasicRequest"></checkbox>
				<p class="comment">只做基础的请求，不处理反向代理（不检查源站）、WAF等。</p>
			</td>
		</tr>
		<tr>
			<td>记录访问日志</td>
			<td>
				<checkbox v-model="healthCheck.accessLogIsOn"></checkbox>
				<p class="comment">记录健康检查的访问日志。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`
})