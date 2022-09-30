Vue.component("download-link", {
	props: ["v-element", "v-file", "v-value"],
	created: function () {
		let that = this
		setTimeout(function () {
			that.url = that.composeURL()
		}, 1000)
	},
	data: function () {
		let filename = this.vFile
		if (filename == null || filename.length == 0) {
			filename = "unknown-file"
		}
		return {
			file: filename,
			url: this.composeURL()
		}
	},
	methods: {
		composeURL: function () {
			let text = ""
			if (this.vValue != null) {
				text = this.vValue
			} else {
				let e = document.getElementById(this.vElement)
				if (e == null) {
					// 不提示错误，因为此时可能页面未加载完整
					return
				}
				text = e.innerText
				if (text == null) {
					text = e.textContent
				}
			}
			return Tea.url("/ui/download", {
				file: this.file,
				text: text
			})
		}
	},
	template: `<a :href="url" target="_blank" style="font-weight: normal"><slot></slot></a>`,
})