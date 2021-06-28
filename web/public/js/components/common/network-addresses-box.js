Vue.component("network-addresses-box", {
	props: ["v-server-type", "v-addresses", "v-protocol", "v-name", "v-from"],
	data: function () {
		let addresses = this.vAddresses
		if (addresses == null) {
			addresses = []
		}
		let protocol = this.vProtocol
		if (protocol == null) {
			protocol = ""
		}

		let name = this.vName
		if (name == null) {
			name = "addresses"
		}

		let from = this.vFrom
		if (from == null) {
			from = ""
		}

		return {
			addresses: addresses,
			protocol: protocol,
			name: name,
			from: from
		}
	},
	watch: {
		"vServerType": function () {
			this.addresses = []
		},
		"vAddresses": function () {
			if (this.vAddresses != null) {
				this.addresses = this.vAddresses
			}
		}
	},
	methods: {
		addAddr: function () {
			let that = this
			window.UPDATING_ADDR = null
			teaweb.popup("/servers/addPortPopup?serverType=" + this.vServerType + "&protocol=" + this.protocol + "&from=" + this.from, {
				height: "16em",
				callback: function (resp) {
					var addr = resp.data.address
					that.addresses.push(addr)
					if (["https", "https4", "https6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "HTTPS"
					} else if (["tls", "tls4", "tls6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "TLS"
					}

					// 发送事件
					that.$emit("change", that.addresses)
				}
			})
		},
		removeAddr: function (index) {
			this.addresses.$remove(index);

			// 发送事件
			this.$emit("change", this.addresses)
		},
		updateAddr: function (index, addr) {
			let that = this
			window.UPDATING_ADDR = addr
			teaweb.popup("/servers/addPortPopup?serverType=" + this.vServerType + "&protocol=" + this.protocol, {
				height: "16em",
				callback: function (resp) {
					var addr = resp.data.address
					Vue.set(that.addresses, index, addr)

					if (["https", "https4", "https6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "HTTPS"
					} else if (["tls", "tls4", "tls6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "TLS"
					}

					// 发送事件
					that.$emit("change", that.addresses)
				}
			})

			// 发送事件
			this.$emit("change", this.addresses)
		}
	},
	template: `<div>
	<input type="hidden" :name="name" :value="JSON.stringify(addresses)"/>
	<div v-if="addresses.length > 0">
		<div class="ui label small basic" v-for="(addr, index) in addresses">
			{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host}}</span><span v-if="addr.host.length == 0">*</span>:{{addr.portRange}}
			<a href="" @click.prevent="updateAddr(index, addr)" title="修改"><i class="icon pencil small"></i></a>
			<a href="" @click.prevent="removeAddr(index)" title="删除"><i class="icon remove"></i></a> </div>
		<div class="ui divider"></div>
	</div>
	<a href="" @click.prevent="addAddr()">[添加端口绑定]</a>
</div>`
})