Vue.component("http-redirect-to-https-box", {
	props: ["v-redirect-to-https-config", "v-is-location"],
	data: function () {
		let redirectToHttpsConfig = this.vRedirectToHttpsConfig
		if (redirectToHttpsConfig == null) {
			redirectToHttpsConfig = {
				isPrior: false,
				isOn: false,
				host: "",
				port: 0,
				status: 0,
				onlyDomains: [],
				exceptDomains: []
			}
		} else {
			if (redirectToHttpsConfig.onlyDomains == null) {
				redirectToHttpsConfig.onlyDomains = []
			}
			if (redirectToHttpsConfig.exceptDomains == null) {
				redirectToHttpsConfig.exceptDomains = []
			}
		}
		return {
			redirectToHttpsConfig: redirectToHttpsConfig,
			portString: (redirectToHttpsConfig.port > 0) ? redirectToHttpsConfig.port.toString() : "",
			moreOptionsVisible: false,
			statusOptions: [
				{"code": 301, "text": "Moved Permanently"},
				{"code": 308, "text": "Permanent Redirect"},
				{"code": 302, "text": "Found"},
				{"code": 303, "text": "See Other"},
				{"code": 307, "text": "Temporary Redirect"}
			]
		}
	},
	watch: {
		"redirectToHttpsConfig.status": function () {
			this.redirectToHttpsConfig.status = parseInt(this.redirectToHttpsConfig.status)
		},
		portString: function (v) {
			let port = parseInt(v)
			if (!isNaN(port)) {
				this.redirectToHttpsConfig.port = port
			} else {
				this.redirectToHttpsConfig.port = 0
			}
		}
	},
	methods: {
		changeMoreOptions: function (isVisible) {
			this.moreOptionsVisible = isVisible
		},
		changeOnlyDomains: function (values) {
			this.redirectToHttpsConfig.onlyDomains = values
			this.$forceUpdate()
		},
		changeExceptDomains: function (values) {
			this.redirectToHttpsConfig.exceptDomains = values
			this.$forceUpdate()
		}
	},
	template: `<div>
	<input type="hidden" name="redirectToHTTPSJSON" :value="JSON.stringify(redirectToHttpsConfig)"/>
	
	<!-- Location -->
	<table class="ui table selectable definition" v-if="vIsLocation">
		<prior-checkbox :v-config="redirectToHttpsConfig"></prior-checkbox>
		<tbody v-show="redirectToHttpsConfig.isPrior">
			<tr>
				<td class="title">自动跳转到HTTPS</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="redirectToHttpsConfig.isOn"/>
						<label></label>
					</div>
					<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上，<more-options-angle @change="changeMoreOptions"></more-options-angle></p>
					
					<!--  TODO 如果已经设置了特殊设置，需要在界面上显示 -->
					<table class="ui table" v-show="moreOptionsVisible">
						<tr>
							<td class="title">状态码</td>
							<td>
								<select class="ui dropdown auto-width" v-model="redirectToHttpsConfig.status">
									<option value="0">[使用默认]</option>
									<option v-for="option in statusOptions" :value="option.code">{{option.code}} {{option.text}}</option>
								</select>
							</td>
						</tr>
						<tr>
							<td>域名或IP地址</td>
							<td>
								<input type="text" name="host" v-model="redirectToHttpsConfig.host"/>
								<p class="comment">默认和用户正在访问的域名或IP地址一致。</p>
							</td>
						</tr>
						<tr>
							<td>端口</td>
							<td>
								<input type="text" name="port" v-model="portString" maxlength="5" style="width:6em"/>
								<p class="comment">默认端口为443。</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>	
		</tbody>
	</table>
	
	<!-- 非Location -->
	<div v-if="!vIsLocation">
		<div class="ui checkbox">
			<input type="checkbox" v-model="redirectToHttpsConfig.isOn"/>
			<label></label>
		</div>
		<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上，<more-options-angle @change="changeMoreOptions"></more-options-angle></p>
		
		<!--  TODO 如果已经设置了特殊设置，需要在界面上显示 -->
		<table class="ui table" v-show="moreOptionsVisible">
			<tr>
				<td class="title">状态码</td>
				<td>
					<select class="ui dropdown auto-width" v-model="redirectToHttpsConfig.status">
						<option value="0">[使用默认]</option>
						<option v-for="option in statusOptions" :value="option.code">{{option.code}} {{option.text}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>跳转后域名或IP地址</td>
				<td>
					<input type="text" name="host" v-model="redirectToHttpsConfig.host"/>
					<p class="comment">默认和用户正在访问的域名或IP地址一致，不填写就表示使用当前的域名。</p>
				</td>
			</tr>
			<tr>
				<td>端口</td>
				<td>
					<input type="text" name="port" v-model="portString" maxlength="5" style="width:6em"/>
					<p class="comment">默认端口为443。</p>
				</td>
			</tr>
			<tr>
				<td>允许的域名</td>
				<td>
					<domains-box :v-domains="redirectToHttpsConfig.onlyDomains" @change="changeOnlyDomains"></domains-box>
					<p class="comment">如果填写了允许的域名，那么只有这些域名可以自动跳转。</p>
				</td>
			</tr>
			<tr>
				<td>排除的域名</td>
				<td>
					<domains-box :v-domains="redirectToHttpsConfig.exceptDomains" @change="changeExceptDomains"></domains-box>
					<p class="comment">如果填写了排除的域名，那么这些域名将不跳转。</p>
				</td>
			</tr>
		</table>
	</div>
	<div class="margin"></div>
</div>`
})