Vue.component("node-ip-addresses-box", {
	props: ["vIpAddresses"],
	data: function () {
		return {
			ipAddresses: (this.vIpAddresses == null) ? [] : this.vIpAddresses
		}
	},
	methods: {
		// 添加IP地址
		addIPAddress: function () {
			window.UPDATING_NODE_IP_ADDRESS = null

			let that = this;
			teaweb.popup("/nodes/ipAddresses/createPopup", {
				callback: function (resp) {
					that.ipAddresses.push(resp.data.ipAddress);
				}
			})
		},

		// 修改地址
		updateIPAddress: function (index, address) {
			window.UPDATING_NODE_IP_ADDRESS = address

			let that = this;
			teaweb.popup("/nodes/ipAddresses/updatePopup", {
				callback: function (resp) {
					Vue.set(that.ipAddresses, index, resp.data.ipAddress);
				}
			})
		},

		// 删除IP地址
		removeIPAddress: function (index) {
			this.ipAddresses.$remove(index);
		},

		// 判断是否为IPv6
		isIPv6: function (ip) {
			return ip.indexOf(":") > -1
		}
	},
	template: `<div>
	<input type="hidden" name="ipAddressesJSON" :value="JSON.stringify(ipAddresses)"/>
	<div v-if="ipAddresses.length > 0">
		<div>
			<div v-for="(address, index) in ipAddresses" class="ui label tiny basic">
				<span v-if="isIPv6(address.ip)" class="grey">[IPv6]</span> {{address.ip}}
				<span class="small grey" v-if="address.name.length > 0">（{{address.name}}<span v-if="!address.canAccess">，不可访问</span>）</span>
				<span class="small grey" v-if="address.name.length == 0 && !address.canAccess">（不可访问）</span>
				<a href="" title="修改" @click.prevent="updateIPAddress(index, address)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeIPAddress(index)"><i class="icon remove"></i></a>
			</div>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button small" type="button" @click.prevent="addIPAddress()">+</button>
	</div>
</div>`
})