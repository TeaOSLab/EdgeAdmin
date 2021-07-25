Vue.component("ip-box", {
	props: [],
	methods: {
		popup: function () {
			let e = this.$refs.container
			let text = e.innerText
			if (text == null) {
				text = e.textContent
			}

			teaweb.popup("/servers/ipbox?ip=" + text, {
				width: "50em",
				height: "30em"
			})
		}
	},
	template: `<span @click.prevent="popup()" ref="container"><slot></slot></span>`
})