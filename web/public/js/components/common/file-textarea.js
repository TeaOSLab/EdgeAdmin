Vue.component("file-textarea", {
	props: ["value"],
	data: function () {
		let value = this.value
		if (typeof value != "string") {
			value = ""
		}
		return {
			realValue: value
		}
	},
	mounted: function () {
	},
	methods: {
		dragover: function () {},
		drop: function (e) {
			let that = this
			e.dataTransfer.items[0].getAsFile().text().then(function (data) {
				that.setValue(data)
			})
		},
		setValue: function (value) {
			this.realValue = value
		},
		focus: function () {
			this.$refs.textarea.focus()
		}
	},
	template: `<textarea @drop.prevent="drop" @dragover.prevent="dragover" ref="textarea" v-model="realValue"></textarea>`
})