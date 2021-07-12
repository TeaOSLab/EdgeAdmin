Vue.component("keyword", {
	props: ["v-word"],
	data: function () {
		let word = this.vWord
		if (word == null) {
			word = ""
		} else {
			word = word.replace(/\)/, "\\)")
			word = word.replace(/\(/, "\\(")
			word = word.replace(/\+/, "\\+")
			word = word.replace(/\^/, "\\^")
			word = word.replace(/\$/, "\\$")
		}

		let slot = this.$slots["default"][0]
		let text = this.encodeHTML(slot.text)
		if (word.length > 0) {
			text = text.replace(new RegExp("(" + word + ")", "ig"), "<span style=\"border: 1px #ccc dashed; color: #ef4d58\">$1</span>")
		}

		return {
			word: word,
			text: text
		}
	},
	methods: {
		encodeHTML: function (s) {
			s = s.replace("&", "&amp;")
			s = s.replace("<", "&lt;")
			s = s.replace(">", "&gt;")
			return s
		}
	},
	template: `<span><span style="display: none"><slot></slot></span><span v-html="text"></span></span>`
})