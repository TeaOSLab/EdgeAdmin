Vue.component("http-firewall-js-cookie-options-viewer", {
	props: ["v-js-cookie-options"],
	mounted: function () {
		this.updateSummary()
	},
	data: function () {
		let options = this.vJsCookieOptions
		if (options == null) {
			options = {
				life: 0,
				maxFails: 0,
				failBlockTimeout: 0,
				failBlockScopeAll: false,
				scope: ""
			}
		}
		return {
			options: options,
			summary: ""
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
				summaryList.push("尝试全局封禁")
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