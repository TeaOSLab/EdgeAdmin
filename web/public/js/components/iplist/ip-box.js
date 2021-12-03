Vue.component("ip-box", {
	props: ["v-ip"],
	methods: {
		popup: function () {
			let ip = this.vIp
			if (ip == null || ip.length == 0) {
				let e = this.$refs.container
				ip = e.innerText
				if (ip == null) {
					ip = e.textContent
				}
			}

			teaweb.popup("/servers/ipbox?ip=" + ip, {
				width: "50em",
				height: "30em"
			})
		}
	},
	template: `<span @click.prevent="popup()" ref="container"><slot></slot></span>`
})