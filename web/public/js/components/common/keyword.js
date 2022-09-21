Vue.component("keyword", {
	props: ["v-word"],
	data: function () {
		let word = this.vWord
		if (word == null) {
			word = ""
		} else {
			word = word.replace(/\)/g, "\\)")
			word = word.replace(/\(/g, "\\(")
			word = word.replace(/\+/g, "\\+")
			word = word.replace(/\^/g, "\\^")
			word = word.replace(/\$/g, "\\$")
			word = word.replace(/\?/g, "\\?")
			word = word.replace(/\*/g, "\\*")
			word = word.replace(/\[/g, "\\[")
			word = word.replace(/{/g, "\\{")
			word = word.replace(/\./g, "\\.")
		}

		let slot = this.$slots["default"][0]
		let text = slot.text
		if (word.length > 0) {
			let that = this
			let m = []  // replacement => tmp
			let tmpIndex = 0
			text = text.replaceAll(new RegExp("(" + word + ")", "ig"), function (replacement) {
				tmpIndex++
				let s = "<span style=\"border: 1px #ccc dashed; color: #ef4d58\">" + that.encodeHTML(replacement) + "</span>"
				let tmpKey = "$TMP__KEY__" + tmpIndex.toString() + "$"
				m.push([tmpKey, s])
				return tmpKey
			})
			text = this.encodeHTML(text)

			m.forEach(function (r) {
				text = text.replace(r[0], r[1])
			})

		} else {
			text = this.encodeHTML(text)
		}

		return {
			word: word,
			text: text
		}
	},
	methods: {
		encodeHTML: function (s) {
			s = s.replace(/&/g, "&amp;")
			s = s.replace(/</g, "&lt;")
			s = s.replace(/>/g, "&gt;")
			s = s.replace(/"/g, "&quot;")
			return s
		}
	},
	template: `<span><span style="display: none"><slot></slot></span><span v-html="text"></span></span>`
})