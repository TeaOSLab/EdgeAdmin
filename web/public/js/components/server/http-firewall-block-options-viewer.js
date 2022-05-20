Vue.component("http-firewall-block-options-viewer", {
	props: ["v-block-options"],
	data: function () {
		return {
			blockOptions: this.vBlockOptions,
			statusCode: this.vBlockOptions.statusCode,
			timeout: this.vBlockOptions.timeout
		}
	},
	watch: {
		statusCode: function (v) {
			let statusCode = parseInt(v)
			if (isNaN(statusCode)) {
				this.blockOptions.statusCode = 403
			} else {
				this.blockOptions.statusCode = statusCode
			}
		},
		timeout: function (v) {
			let timeout = parseInt(v)
			if (isNaN(timeout)) {
				this.blockOptions.timeout = 0
			} else {
				this.blockOptions.timeout = timeout
			}
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	状态码：{{statusCode}} / 提示内容：<span v-if="blockOptions.body != null && blockOptions.body.length > 0">[{{blockOptions.body.length}}字符]</span><span v-else class="disabled">[无]</span>  / 超时时间：{{timeout}}秒
</div>	
`
})