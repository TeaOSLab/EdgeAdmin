Vue.component("download-link", {
	props: ["v-element", "v-file"],
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
			let e = document.getElementById(this.vElement)
			if (e == null) {
				teaweb.warn("<download-link>找不到要下载的内容")
				return
			}
			let text = e.innerText
			if (text == null) {
				text = e.textContent
			}
			return Tea.url("/ui/download", {
				file: this.file,
				text: text
			})
		}
	},
	template: `<a :href="url" target="_blank"><slot></slot></a>`,
})