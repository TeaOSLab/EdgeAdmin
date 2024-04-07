Vue.component("http-firewall-captcha-options", {
	props: ["v-captcha-options"],
	mounted: function () {
		this.updateSummary()
	},
	data: function () {
		let options = this.vCaptchaOptions
		if (options == null) {
			options = {
				captchaType: "default",
				countLetters: 0,
				life: 0,
				maxFails: 0,
				failBlockTimeout: 0,
				failBlockScopeAll: false,
				uiIsOn: false,
				uiTitle: "",
				uiPrompt: "",
				uiButtonTitle: "",
				uiShowRequestId: true,
				uiCss: "",
				uiFooter: "",
				uiBody: "",
				cookieId: "",
				lang: "",
				geeTestConfig: {
					isOn: false,
					captchaId: "",
					captchaKey: ""
				}
			}
		}
		if (options.countLetters <= 0) {
			options.countLetters = 6
		}

		if (options.captchaType == null || options.captchaType.length == 0) {
			options.captchaType = "default"
		}


		return {
			options: options,
			isEditing: false,
			summary: "",
			uiBodyWarning: "",
			captchaTypes: window.WAF_CAPTCHA_TYPES
		}
	},
	watch: {
		"options.countLetters": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			} else if (i < 0) {
				i = 0
			} else if (i > 10) {
				i = 10
			}
			this.options.countLetters = i
		},
		"options.life": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.life = i
			this.updateSummary()
		},
		"options.maxFails": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.maxFails = i
			this.updateSummary()
		},
		"options.failBlockTimeout": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.failBlockTimeout = i
			this.updateSummary()
		},
		"options.failBlockScopeAll": function (v) {
			this.updateSummary()
		},
		"options.captchaType": function (v) {
			this.updateSummary()
		},
		"options.uiIsOn": function (v) {
			this.updateSummary()
		},
		"options.uiBody": function (v) {
			if (/<form(>|\s).+\$\{body}.*<\/form>/s.test(v)) {
				this.uiBodyWarning = "页面模板中不能使用<form></form>标签包裹\${body}变量，否则将导致验证码表单无法提交。"
			} else {
				this.uiBodyWarning = ""
			}
		},
		"options.geeTestConfig.isOn": function (v) {
			this.updateSummary()
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		},
		updateSummary: function () {
			let summaryList = []
			if (this.options.life > 0) {
				summaryList.push("有效时间" + this.options.life + "秒")
			}
			if (this.options.maxFails > 0) {
				summaryList.push("最多失败" + this.options.maxFails + "次")
			}
			if (this.options.failBlockTimeout > 0) {
				summaryList.push("失败拦截" + this.options.failBlockTimeout + "秒")
			}
			if (this.options.failBlockScopeAll) {
				summaryList.push("尝试全局封禁")
			}

			let that = this
			let typeDef = this.captchaTypes.$find(function (k, v) {
				return v.code == that.options.captchaType
			})
			if (typeDef != null) {
				summaryList.push("默认验证方式：" + typeDef.name)
			}

			if (this.options.captchaType == "default") {
				if (this.options.uiIsOn) {
					summaryList.push("定制UI")
				}
			}

			if (this.options.geeTestConfig != null && this.options.geeTestConfig.isOn) {
				summaryList.push("已配置极验")
			}

			if (summaryList.length == 0) {
				this.summary = "默认配置"
			} else {
				this.summary = summaryList.join(" / ")
			}
		},
		confirm: function () {
			this.isEditing = false
		}
	},
	template: `<div>
	<input type="hidden" name="captchaOptionsJSON" :value="JSON.stringify(options)"/>
	<a href="" @click.prevent="edit">{{summary}} <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<div v-show="isEditing" style="margin-top: 0.5em">
		<table class="ui table definition selectable">
			<tbody>
				<tr>
					<td>默认验证方式</td>
					<td>
						<select class="ui dropdown auto-width" v-model="options.captchaType">
							<option v-for="captchaDef in captchaTypes" :value="captchaDef.code">{{captchaDef.name}}</option>
						</select>
						<p class="comment" v-for="captchaDef in captchaTypes" v-if="captchaDef.code == options.captchaType">{{captchaDef.description}}</p>
					</td>
				</tr>
				<tr>
					<td class="title">有效时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">验证通过后在这个时间内不再验证，默认600秒。</p>
					</td>
				</tr>
				<tr>
					<td>最多失败次数</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.maxFails" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">次</span>
						</div>
						<p class="comment"><span v-if="options.maxFails > 0 && options.maxFails < 5" class="red">建议填入一个不小于5的数字，以减少误判几率。</span>允许用户失败尝试的最多次数，超过这个次数将被自动加入黑名单。如果为空或者为0，表示不限制。</p>
					</td>
				</tr>
				<tr>
					<td>失败拦截时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.failBlockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">在达到最多失败次数（大于0）时，自动拦截的时长；如果为0表示不自动拦截。</p>
					</td>
				</tr>
				<tr>
					<td>失败全局封禁</td>
					<td>
						<checkbox v-model="options.failBlockScopeAll"></checkbox>
						<p class="comment">选中后，表示允许系统尝试全局封禁某个IP，以提升封禁性能。</p>
					</td>
				</tr>
				
				<tr v-show="options.captchaType == 'default'">
					<td>验证码中数字个数</td>
					<td>
						<select class="ui dropdown auto-width" v-model="options.countLetters">
							<option v-for="i in 10" :value="i">{{i}}</option>
						</select>
					</td>
				</tr>
				<tr v-show="options.captchaType == 'default'">
					<td class="color-border">定制UI</td>
					<td><checkbox v-model="options.uiIsOn"></checkbox></td>
				</tr>
			</tbody>
			<tbody v-show="options.uiIsOn && options.captchaType == 'default'">
				<tr>
					<td class="color-border">页面标题</td>
					<td>
						<input type="text" v-model="options.uiTitle" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
					</td>
				</tr>
				<tr>
					<td class="color-border">按钮标题</td>
					<td>
						<input type="text" v-model="options.uiButtonTitle" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<p class="comment">类似于<code-label>提交验证</code-label>。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">显示请求ID</td>
					<td>
						<checkbox v-model="options.uiShowRequestId"></checkbox>
						<p class="comment">在界面上显示请求ID，方便用户报告问题。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">CSS样式</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiCss" rows="2"></textarea>
					</td>
				</tr>
				<tr>
					<td class="color-border">页头提示</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiPrompt" rows="2"></textarea>
						<p class="comment">类似于<code-label>请输入上面的验证码</code-label>，支持HTML。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">页尾提示</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiFooter" rows="2"></textarea>
						<p class="comment">支持HTML。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">页面模板</td>
					<td>
						<textarea spellcheck="false" rows="2" v-model="options.uiBody"></textarea>
						<p class="comment"><span v-if="uiBodyWarning.length > 0" class="red">警告：{{uiBodyWarning}}</span><span v-if="options.uiBody.length > 0 && options.uiBody.indexOf('\${body}') < 0 " class="red">模板中必须包含\${body}表示验证码表单！</span>整个页面的模板，支持HTML，其中必须使用<code-label>\${body}</code-label>变量代表验证码表单，否则将无法正常显示验证码。</p>
					</td>
				</tr>
			</tbody>
		</table>
		
		<table class="ui table definition selectable">
			<tr>
				<td class="title">允许用户使用极验</td>
				<td><checkbox v-model="options.geeTestConfig.isOn"></checkbox>
					<p class="comment">选中后，表示允许用户在WAF设置中选择极验。</p>
				</td>
			</tr>
			<tbody v-show="options.geeTestConfig.isOn">
				<tr>
					<td class="color-border">极验-验证ID *</td>
					<td>
						<input type="text" maxlength="100" name="geetestCaptchaId" v-model="options.geeTestConfig.captchaId" spellcheck="false"/>
						<p class="comment">在极验控制台--业务管理中获取。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">极验-验证Key *</td>
					<td>
						<input type="text" maxlength="100" name="geetestCaptchaKey" v-model="options.geeTestConfig.captchaKey" spellcheck="false"/>
						<p class="comment">在极验控制台--业务管理中获取。</p>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
`
})