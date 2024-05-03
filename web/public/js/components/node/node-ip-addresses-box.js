// 节点IP地址管理（标签形式）
Vue.component("node-ip-addresses-box", {
	props: ["v-ip-addresses", "role", "v-node-id"],
	data: function () {
		let nodeId = this.vNodeId
		if (nodeId == null) {
			nodeId = 0
		}

		return {
			ipAddresses: (this.vIpAddresses == null) ? [] : this.vIpAddresses,
			supportThresholds: this.role != "ns",
			nodeId: nodeId
		}
	},
	methods: {
		// 添加IP地址
		addIPAddress: function () {
			window.UPDATING_NODE_IP_ADDRESS = null

			let that = this;
			teaweb.popup("/nodes/ipAddresses/createPopup?nodeId=" + this.nodeId + "&supportThresholds=" + (this.supportThresholds ? 1 : 0), {
				callback: function (resp) {
					that.ipAddresses.push(resp.data.ipAddress);
				},
				height: "24em",
				width: "44em"
			})
		},

		// 修改地址
		updateIPAddress: function (index, address) {
			window.UPDATING_NODE_IP_ADDRESS = teaweb.clone(address)

			let that = this;
			teaweb.popup("/nodes/ipAddresses/updatePopup?nodeId=" + this.nodeId + "&supportThresholds=" + (this.supportThresholds ? 1 : 0), {
				callback: function (resp) {
					Vue.set(that.ipAddresses, index, resp.data.ipAddress);
				},
				height: "24em",
				width: "44em"
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
				<span class="small grey" v-if="address.name.length > 0">（备注：{{address.name}}<span v-if="!address.canAccess">，不公开访问</span>）</span>
				<span class="small grey" v-if="address.name.length == 0 && !address.canAccess">（不公开访问）</span>
				<span class="small red" v-if="!address.isOn" title="未启用">[off]</span>
				<span class="small red" v-if="!address.isUp" title="已下线">[down]</span>
				<span class="small" v-if="address.thresholds != null && address.thresholds.length > 0">[{{address.thresholds.length}}个阈值]</span>
				&nbsp;
				 <span v-if="address.clusters != null && address.clusters.length > 0">
					&nbsp; <span class="small grey">专属集群：[</span><span v-for="(cluster, index) in address.clusters" class="small grey">{{cluster.name}}<span v-if="index < address.clusters.length - 1">，</span></span><span class="small grey">]</span>
					&nbsp;
				</span>
				
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