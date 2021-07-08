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
		let text = slot.text
		if (word.length > 0) {
			text = text.replace(new RegExp("(" + word + ")", "ig"), "<span style=\"border: 1px #ccc dashed; color: #ef4d58\">$1</span>")
		}

		return {
			word: word,
			text: text
		}
	},
	template: `<span><span style="display: none"><slot></slot></span><span v-html="text"></span></span>`
})