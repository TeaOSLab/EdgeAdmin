Vue.component("health-check-config-box", {
	props: ["v-health-check-config"],
	data: function () {
		let healthCheckConfig = this.vHealthCheckConfig
		let urlProtocol = "http"
		let urlPort = ""
		let urlRequestURI = "/"

		if (healthCheckConfig == null) {
			healthCheckConfig = {
				isOn: false,
				url: "",
				interval: {count: 60, unit: "second"},
				statusCodes: [200],
				timeout: {count: 10, unit: "second"},
				countTries: 3,
				tryDelay: {count: 100, unit: "ms"}
			}
			let that = this
			setTimeout(function () {
				that.changeURL()
			}, 500)
		} else {
			let url = new URL(healthCheckConfig.url)
			urlProtocol = url.protocol.substring(0, url.protocol.length - 1)
			urlPort = url.port
			urlRequestURI = url.pathname
			if (url.search.length > 0) {
				urlRequestURI += url.search
			}
		}
		return {
			healthCheck: healthCheckConfig,
			advancedVisible: false,
			urlProtocol: urlProtocol,
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
		"healthCheck.countTries": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countTries = count
			} else {
				this.healthCheck.countTries = 0
			}
		}
	},
	methods: {
		showAdvanced: function () {
			this.advancedVisible = !this.advancedVisible
		},
		changeURL: function () {
			this.healthCheck.url = this.urlProtocol + "://${host}" + ((this.urlPort.length > 0) ? ":" + this.urlPort : "") + this.urlRequestURI
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
			console.log(this.healthCheck.statusCodes)
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
				<div class="ui fields inline">
					<div class="ui field">
						<select class="ui dropdown" v-model="urlProtocol">
							<option value="http">http://</option>
							<option value="https">https://</option>
						</select>
					</div>
					<div class="ui field">
						<var style="color:grey">\${host}</var>
					</div>
					<div class="ui field">:</div>
					<div class="ui field">
						<input type="text" maxlength="5" style="width:5.4em" placeholder="端口" v-model="urlPort"/>
					</div>
					<div class="ui field">
						<input type="text" v-model="urlRequestURI" placeholder="/" style="width:23em"/>
					</div>
				</div>
				<div class="ui divider"></div>
				<p class="comment" v-if="healthCheck.url.length > 0">拼接后的URL：<code-label>{{healthCheck.url}}</code-label>，其中\${host}指的是节点地址。</p>
			</td>
		</tr>
		<tr>
			<td>检测时间间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.interval"></time-duration-box>
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