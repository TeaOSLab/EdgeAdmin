Vue.component("copy-to-clipboard", {
	props: ["v-target"],
	created: function () {
		if (typeof ClipboardJS == "undefined") {
			let jsFile = document.createElement("script")
			jsFile.setAttribute("src", "/js/clipboard.min.js")
			document.head.appendChild(jsFile)
		}
	},
	methods: {
		copy: function () {
			new ClipboardJS('[data-clipboard-target]');
			teaweb.successToast("已复制到剪切板")
		}
	},
	template: `<a href="" title="拷贝到剪切板" :data-clipboard-target="'#' + vTarget" @click.prevent="copy"><i class="ui icon copy small"></i></em></a>`
})