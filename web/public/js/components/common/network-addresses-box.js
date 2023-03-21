Vue.component("network-addresses-box", {
	props: ["v-server-type", "v-addresses", "v-protocol", "v-name", "v-from", "v-support-range", "v-url"],
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
			from: from,
			isEditing: false
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
			this.isEditing = true

			let that = this
			window.UPDATING_ADDR = null

			let url = this.vUrl
			if (url == null) {
				url = "/servers/addPortPopup"
			}

			teaweb.popup(url + "?serverType=" + this.vServerType + "&protocol=" + this.protocol + "&from=" + this.from + "&supportRange=" + (this.supportRange() ? 1 : 0), {
				height: "18em",
				callback: function (resp) {
					var addr = resp.data.address
					if (that.addresses.$find(function (k, v) {
						return addr.host == v.host && addr.portRange == v.portRange && addr.protocol == v.protocol
					}) != null) {
						teaweb.warn("要添加的网络地址已经存在")
						return
					}
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

			let url = this.vUrl
			if (url == null) {
				url = "/servers/addPortPopup"
			}

			teaweb.popup(url + "?serverType=" + this.vServerType + "&protocol=" + this.protocol + "&from=" + this.from + "&supportRange=" + (this.supportRange() ? 1 : 0), {
				height: "18em",
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
		},
		supportRange: function () {
			return this.vSupportRange || (this.vServerType == "tcpProxy" || this.vServerType == "udpProxy")
		},
		edit: function () {
			this.isEditing = true
		}
	},
	template: `<div>
	<input type="hidden" :name="name" :value="JSON.stringify(addresses)"/>
	<div v-show="!isEditing">
		<div v-if="addresses.length > 0">
			<div class="ui label small basic" v-for="(addr, index) in addresses">
				{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host.quoteIP()}}</span><span v-if="addr.host.length == 0">*</span>:<span v-if="addr.portRange.indexOf('-')<0">{{addr.portRange}}</span><span v-else style="font-style: italic">{{addr.portRange}}</span>
			</div>
			&nbsp; &nbsp; <a href="" @click.prevent="edit" style="font-size: 0.9em">[修改]</a>
		</div>	
	</div>
	<div v-show="isEditing || addresses.length == 0">
		<div v-if="addresses.length > 0">
			<div class="ui label small basic" v-for="(addr, index) in addresses">
				{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host.quoteIP()}}</span><span v-if="addr.host.length == 0">*</span>:<span v-if="addr.portRange.indexOf('-')<0">{{addr.portRange}}</span><span v-else style="font-style: italic">{{addr.portRange}}</span>
				<a href="" @click.prevent="updateAddr(index, addr)" title="修改"><i class="icon pencil small"></i></a>
				<a href="" @click.prevent="removeAddr(index)" title="删除"><i class="icon remove"></i></a> 
			</div>
			<div class="ui divider"></div>
		</div>
		<a href="" @click.prevent="addAddr()">[添加端口绑定]</a>
	</div>
</div>`
})