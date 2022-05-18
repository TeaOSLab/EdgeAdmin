Vue.component("node-ddos-protection-config-box", {
	props: ["v-ddos-protection-config", "v-default-configs", "v-is-node", "v-cluster-is-on"],
	data: function () {
		let config = this.vDdosProtectionConfig
		if (config == null) {
			config = {
				tcp: {
					isPrior: false,
					isOn: false,
					maxConnections: 0,
					maxConnectionsPerIP: 0,
					newConnectionsRate: 0,
					allowIPList: [],
					ports: []
				}
			}
		}

		// initialize
		if (config.tcp == null) {
			config.tcp = {
				isPrior: false,
				isOn: false,
				maxConnections: 0,
				maxConnectionsPerIP: 0,
				newConnectionsRate: 0,
				allowIPList: [],
				ports: []
			}
		}


		return {
			config: config,
			defaultConfigs: this.vDefaultConfigs,
			isNode: this.vIsNode,

			isAddingPort: false
		}
	},
	methods: {
		changeTCPPorts: function (ports) {
			this.config.tcp.ports = ports
		},
		changeTCPAllowIPList: function (ipList) {
			this.config.tcp.allowIPList = ipList
		}
	},
	template: `<div>
 <input type="hidden" name="ddosProtectionJSON" :value="JSON.stringify(config)"/>

 <p class="comment">功能说明：此功能为<strong>试验性质</strong>，目前仅能防御简单的DDoS攻击，试验期间建议仅在被攻击时启用，仅支持已安装<code-label>nftables v0.9</code-label>以上的Linux系统。<pro-warning-label></pro-warning-label></p>

 <div class="ui message" v-if="vClusterIsOn">当前节点所在集群已设置DDoS防护。</div>

 <h4>TCP设置</h4>
 <table class="ui table definition selectable">
 	<prior-checkbox :v-config="config.tcp" v-if="isNode"></prior-checkbox>
 	<tbody v-show="config.tcp.isPrior || !isNode">
		<tr>
			<td class="title">启用</td>
			<td>
				<checkbox v-model="config.tcp.isOn"></checkbox>
			</td>
		</tr>
	</tbody>
	<tbody v-show="config.tcp.isOn && (config.tcp.isPrior || !isNode)">
		<tr>
			<td class="title">单节点TCP最大连接数</td>
			<td>
				<digit-input name="tcpMaxConnections" v-model="config.tcp.maxConnections" maxlength="6" size="6" style="width: 6em"></digit-input>
				<p class="comment">单个节点可以接受的TCP最大连接数。如果为0，则默认为{{defaultConfigs.tcpMaxConnections}}。</p>
			</td>
		</tr>
		<tr>
			<td>单IP TCP最大连接数</td>
			<td>
				<digit-input name="tcpMaxConnectionsPerIP" v-model="config.tcp.maxConnectionsPerIP" maxlength="6" size="6" style="width: 6em"></digit-input>
				<p class="comment">单个IP可以连接到节点的TCP最大连接数。如果为0，则默认为{{defaultConfigs.tcpMaxConnectionsPerIP}}；最小值为{{defaultConfigs.tcpMinConnectionsPerIP}}。</p>
			</td>
		</tr>
		<tr>
			<td>单IP TCP新连接速率</td>
			<td>
				<div class="ui input right labeled">
					<digit-input name="tcpNewConnectionsRate" v-model="config.tcp.newConnectionsRate" maxlength="6" size="6" style="width: 6em" :min="defaultConfigs.tcpNewConnectionsMinRate"></digit-input>
					<span class="ui label">个新连接/每分钟</span>
				</div>
				<p class="comment">单个IP可以创建TCP新连接的速率。如果为0，则默认为{{defaultConfigs.tcpNewConnectionsRate}}；最小值为{{defaultConfigs.tcpNewConnectionsMinRate}}。</p>
			</td>
		</tr>
		<tr>
			<td>TCP端口列表</td>
			<td>
				<ddos-protection-ports-config-box :v-ports="config.tcp.ports" @change="changeTCPPorts"></ddos-protection-ports-config-box>
				<p class="comment">默认为80和443两个端口。</p>
			</td>
		</tr>
		<tr>
			<td>IP白名单</td>
			<td>
				<ddos-protection-ip-list-config-box :v-ip-list="config.tcp.allowIPList" @change="changeTCPAllowIPList"></ddos-protection-ip-list-config-box>
				<p class="comment">在白名单中的IP不受当前设置的限制。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`
})