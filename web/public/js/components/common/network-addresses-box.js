Vue.component("network-addresses-box", {
	props: ["v-server-type", "v-addresses", "v-protocol"],
	data: function () {
		let addresses = this.vAddresses
		if (addresses == null) {
			addresses = []
		}
		let protocol = this.vProtocol
		if (protocol == null) {
			protocol = ""
		}
		return {
			addresses: addresses,
			protocol: protocol
		}
	},
	watch: {
		"vServerType": function () {
			this.addresses = []
		}
	},
	methods: {
		addAddr: function () {
			let that = this
			teaweb.popup("/servers/addPortPopup?serverType=" + this.vServerType + "&protocol=" + this.protocol, {
				callback: function (resp) {
					var addr = resp.data.address;
					that.addresses.push(addr);
					if (["https", "https4", "https6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "HTTPS";
					} else if (["tls", "tls4", "tls6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "TLS";
					}
				}
			})
		},
		removeAddr: function (index) {
			this.addresses.$remove(index);
		}
	},
	template: `<div>
	<input type="hidden" name="addresses" :value="JSON.stringify(addresses)"/>
	<div v-if="addresses.length > 0">
		<div class="ui label small" v-for="(addr, index) in addresses">
			{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host}}</span><span v-if="addr.host.length == 0">*</span>:{{addr.portRange}}
			<a href="" @click.prevent="removeAddr(index)" title="删除"><i class="icon remove"></i></a> </div>
		<div class="ui divider"></div>
	</div>
	<a href="" @click.prevent="addAddr()">[添加端口绑定]</a>
</div>`
})