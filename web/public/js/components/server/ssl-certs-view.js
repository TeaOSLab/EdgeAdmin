Vue.component("ssl-certs-view", {
	props: ["v-certs"],
	data: function () {
		let certs = this.vCerts
		if (certs == null) {
			certs = []
		}
		return {
			certs: certs
		}
	},
	methods: {
		// 格式化时间
		formatTime: function (timestamp) {
			return new Date(timestamp * 1000).format("Y-m-d")
		},

		// 查看详情
		viewCert: function (certId) {
			teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
				height: "28em",
				width: "48em"
			})
		}
	},
	template: `<div>
	<div v-if="certs != null && certs.length > 0">
		<div class="ui label small basic" v-for="(cert, index) in certs">
			{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp;<a href="" title="查看" @click.prevent="viewCert(cert.id)"><i class="icon expand blue"></i></a>
		</div>
	</div>
</div>`
})