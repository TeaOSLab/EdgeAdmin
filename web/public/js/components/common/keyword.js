Vue.component("keyword", {
	props: ["v-word"],
	data: function () {
		let word = this.vWord
		if (word == null) {
			word = ""
		}

		let slot = this.$slots["default"][0]
		let text = slot.text
		if (word.length > 0) {
			text = text.replace(new RegExp(word, "g"), "<span style=\"border: 1px #ccc dashed; color: #ef4d58\">" + word + "</span>")
		}

		return {
			word: word,
			text: text
		}
	},
	template: `<span><span style="display: none"><slot></slot></span><span v-html="text"></span></span>`
})