Vue.component("http-firewall-captcha-options-viewer", {
	props: ["v-captcha-options"],
	mounted: function () {
		this.updateSummary()
	},
	data: function () {
		let options = this.vCaptchaOptions
		if (options == null) {
			options = {
				life: 0,
				maxFails: 0,
				failBlockTimeout: 0,
				failBlockScopeAll: false,
				uiIsOn: false,
				uiTitle: "",
				uiPrompt: "",
				uiButtonTitle: "",
				uiShowRequestId: false,
				uiCss: "",
				uiFooter: "",
				uiBody: "",
				cookieId: "",
				lang: ""
			}
		}
		return {
			options: options,
			summary: "",
			captchaTypes: window.WAF_CAPTCHA_TYPES
		}
	},
	methods: {
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
				summaryList.push("全局封禁")
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
		}
	},
	template: `<div>{{summary}}</div>
`
})