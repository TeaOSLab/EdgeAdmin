Vue.component("traffic-map-box", {
	props: ["v-stats", "v-is-attack"],
	mounted: function () {
		this.render()
	},
	data: function () {
		let maxPercent = 0
		let isAttack = this.vIsAttack
		this.vStats.forEach(function (v) {
			let percent = parseFloat(v.percent)
			if (percent > maxPercent) {
				maxPercent = percent
			}

			v.formattedCountRequests = teaweb.formatCount(v.countRequests) + "次"
			v.formattedCountAttackRequests = teaweb.formatCount(v.countAttackRequests) + "次"
		})

		if (maxPercent < 100) {
			maxPercent *= 1.2 // 不要让某一项100%
		}

		let screenIsNarrow = window.innerWidth < 512

		return {
			isAttack: isAttack,
			stats: this.vStats,
			chart: null,
			minOpacity: 0.2,
			maxPercent: maxPercent,
			selectedCountryName: "",
			screenIsNarrow: screenIsNarrow
		}
	},
	methods: {
		render: function () {
			if (this.$el.offsetWidth < 300) {
				let that = this
				setTimeout(function () {
					that.render()
				}, 100)
				return
			}

			this.chart = teaweb.initChart(document.getElementById("traffic-map-box"));
			let that = this
			this.chart.setOption({
				backgroundColor: "white",
				grid: {
					top: 0,
					bottom: 0,
					left: 0,
					right: 0
				},
				roam: false,
				tooltip: {
					trigger: "item"
				},
				series: [{
					type: "map",
					map: "world",
					zoom: 1.3,
					selectedMode: false,
					itemStyle: {
						areaColor: "#E9F0F9",
						borderColor: "#DDD"
					},
					label: {
						show: false,
						fontSize: "10px",
						color: "#fff",
						backgroundColor: "#8B9BD3",
						padding: [2, 2, 2, 2]
					},
					emphasis: {
						itemStyle: {
							areaColor: "#8B9BD3",
							opacity: 1.0
						},
						label: {
							show: true,
							fontSize: "10px",
							color: "#fff",
							backgroundColor: "#8B9BD3",
							padding: [2, 2, 2, 2]
						}
					},
					//select: {itemStyle:{ areaColor: "#8B9BD3", opacity: 0.8 }},
					tooltip: {
						formatter: function (args) {
							let name = args.name
							let stat = null
							that.stats.forEach(function (v) {
								if (v.name == name) {
									stat = v
								}
							})

							if (stat != null) {
								return name + "<br/>流量：" + stat.formattedBytes + "<br/>流量占比：" + stat.percent + "%<br/>请求数：" + stat.formattedCountRequests + "<br/>攻击数：" + stat.formattedCountAttackRequests
							}
							return name
						}
					},
					data: this.stats.map(function (v) {
						let opacity = parseFloat(v.percent) / that.maxPercent
						if (opacity < that.minOpacity) {
							opacity = that.minOpacity
						}
						let fullOpacity = opacity * 3
						if (fullOpacity > 1) {
							fullOpacity = 1
						}
						let isAttack = that.vIsAttack
						let bgColor = "#276AC6"
						if (isAttack) {
							bgColor = "#B03A5B"
						}

						return {
							name: v.name,
							value: v.bytes,
							percent: parseFloat(v.percent),
							itemStyle: {
								areaColor: bgColor,
								opacity: opacity
							},
							emphasis: {
								itemStyle: {
									areaColor: bgColor,
									opacity: fullOpacity
								},
								label: {
									show: true,
									formatter: function (args) {
										return args.name
									}
								}
							},
							label: {
								show: false,
								formatter: function (args) {
									if (args.name == that.selectedCountryName) {
										return args.name
									}
									return ""
								},
								fontSize: "10px",
								color: "#fff",
								backgroundColor: "#8B9BD3",
								padding: [2, 2, 2, 2]
							}
						}
					}),
					nameMap: window.WorldCountriesMap
				}]
			})
			this.chart.resize()
		},
		selectCountry: function (countryName) {
			if (this.chart == null) {
				return
			}
			let option = this.chart.getOption()
			let that = this
			option.series[0].data.forEach(function (v) {
				let opacity = v.percent / that.maxPercent
				if (opacity < that.minOpacity) {
					opacity = that.minOpacity
				}

				if (v.name == countryName) {
					if (v.isSelected) {
						v.itemStyle.opacity = opacity
						v.isSelected = false
						v.label.show = false
						that.selectedCountryName = ""
						return
					}
					v.isSelected = true
					that.selectedCountryName = countryName
					opacity *= 3
					if (opacity > 1) {
						opacity = 1
					}

					// 至少是0.5，让用户能够看清
					if (opacity < 0.5) {
						opacity = 0.5
					}
					v.itemStyle.opacity = opacity
					v.label.show = true
				} else {
					v.itemStyle.opacity = opacity
					v.isSelected = false
					v.label.show = false
				}
			})
			this.chart.setOption(option)
		},
		select: function (args) {
			this.selectCountry(args.countryName)
		}
	},
	template: `<div>
	<table style="width: 100%; border: 0; padding: 0; margin: 0">
		<tbody>
       	<tr>
           <td>
               <div class="traffic-map-box" id="traffic-map-box"></div>
           </td>
           <td style="width: 14em" v-if="!screenIsNarrow">
           		<traffic-map-box-table :v-stats="stats" :v-is-attack="isAttack" @select="select"></traffic-map-box-table>
           </td>
       </tr>
       </tbody>
       <tbody v-if="screenIsNarrow">
		   <tr>
				<td colspan="2">
					<traffic-map-box-table :v-stats="stats" :v-is-attack="isAttack" :v-screen-is-narrow="true" @select="select"></traffic-map-box-table>
				</td>
			</tr>
		</tbody>
   </table>
</div>`
})

Vue.component("traffic-map-box-table", {
	props: ["v-stats", "v-is-attack", "v-screen-is-narrow"],
	data: function () {
		return {
			stats: this.vStats,
			isAttack: this.vIsAttack
		}
	},
	methods: {
		select: function (countryName) {
			this.$emit("select", {countryName: countryName})
		}
	},
	template: `<div style="overflow-y: auto" :style="{'max-height':vScreenIsNarrow ? 'auto' : '16em'}" class="narrow-scrollbar">
	   <table class="ui table selectable">
		  <thead>
			<tr>
				<th colspan="2">国家/地区排行&nbsp; <tip-icon content="只有开启了统计的网站才会有记录。"></tip-icon></th>
			</tr>
		  </thead>
		   <tbody v-if="stats.length == 0">
			   <tr>
				   <td colspan="2">暂无数据</td>
			   </tr>
		   </tbody>
		   <tbody>
			   <tr v-for="(stat, index) in stats.slice(0, 10)">
				   <td @click.prevent="select(stat.name)" style="cursor: pointer" colspan="2">
					   <div class="ui progress bar" :class="{red: vIsAttack, blue:!vIsAttack}" style="margin-bottom: 0.3em">
						   <div class="bar" style="min-width: 0; height: 4px;" :style="{width: stat.percent + '%'}"></div>
					   </div>
					  <div>{{stat.name}}</div> 
					   <div><span class="grey">{{stat.percent}}% </span>
					   <span class="small grey" v-if="isAttack">{{stat.formattedCountAttackRequests}}</span>
					   <span class="small grey" v-if="!isAttack">（{{stat.formattedBytes}}）</span></div>
				   </td>
			   </tr>
		   </tbody>
	   </table>
   </div>`
})

Vue.component("ddos-protection-ports-config-box", {
	props: ["v-ports"],
	data: function () {
		let ports = this.vPorts
		if (ports == null) {
			ports = []
		}
		return {
			ports: ports,
			isAdding: false,
			addingPort: {
				port: "",
				description: ""
			}
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingPortInput.focus()
			})
		},
		confirm: function () {
			let portString = this.addingPort.port
			if (portString.length == 0) {
				this.warn("请输入端口号")
				return
			}
			if (!/^\d+$/.test(portString)) {
				this.warn("请输入正确的端口号")
				return
			}
			let port = parseInt(portString, 10)
			if (port <= 0) {
				this.warn("请输入正确的端口号")
				return
			}
			if (port > 65535) {
				this.warn("请输入正确的端口号")
				return
			}

			let exists = false
			this.ports.forEach(function (v) {
				if (v.port == port) {
					exists = true
				}
			})
			if (exists) {
				this.warn("端口号已经存在")
				return
			}

			this.ports.push({
				port: port,
				description: this.addingPort.description
			})
			this.notifyChange()
			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.addingPort = {
				port: "",
				description: ""
			}
		},
		remove: function (index) {
			this.ports.$remove(index)
			this.notifyChange()
		},
		warn: function (message) {
			let that = this
			teaweb.warn(message, function () {
				that.$refs.addingPortInput.focus()
			})
		},
		notifyChange: function () {
			this.$emit("change", this.ports)
		}
	},
	template: `<div>
	<div v-if="ports.length > 0">
		<div class="ui label basic tiny" v-for="(portConfig, index) in ports">
			{{portConfig.port}} <span class="grey small" v-if="portConfig.description.length > 0">（{{portConfig.description}}）</span> <a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">端口</span>
					<input type="text" v-model="addingPort.port" ref="addingPortInput" maxlength="5" size="5" placeholder="端口号" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">备注</span>
					<input type="text" v-model="addingPort.description" maxlength="12" size="12" placeholder="备注（可选）" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp;<a href="" @click.prevent="cancel()">取消</a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

// 显示节点的多个集群
Vue.component("node-clusters-labels", {
	props: ["v-primary-cluster", "v-secondary-clusters", "size"],
	data: function () {
		let cluster = this.vPrimaryCluster
		let secondaryClusters = this.vSecondaryClusters
		if (secondaryClusters == null) {
			secondaryClusters = []
		}

		let labelSize = this.size
		if (labelSize == null) {
			labelSize = "small"
		}
		return {
			cluster: cluster,
			secondaryClusters: secondaryClusters,
			labelSize: labelSize
		}
	},
	template: `<div>
	<a v-if="cluster != null" :href="'/clusters/cluster?clusterId=' + cluster.id" title="主集群" style="margin-bottom: 0.3em;">
		<span class="ui label basic grey" :class="labelSize" v-if="labelSize != 'tiny'">{{cluster.name}}</span>
		<grey-label v-if="labelSize == 'tiny'">{{cluster.name}}</grey-label>
	</a>
	<a v-for="c in secondaryClusters" :href="'/clusters/cluster?clusterId=' + c.id" :class="labelSize" title="从集群">
		<span class="ui label basic grey" :class="labelSize" v-if="labelSize != 'tiny'">{{c.name}}</span>
		<grey-label v-if="labelSize == 'tiny'">{{c.name}}</grey-label>
	</a>
</div>`
})

// 单个集群选择
Vue.component("cluster-selector", {
	props: ["v-cluster-id"],
	mounted: function () {
		let that = this

		Tea.action("/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})
	},
	data: function () {
		let clusterId = this.vClusterId
		if (clusterId == null) {
			clusterId = 0
		}
		return {
			clusters: [],
			clusterId: clusterId
		}
	},
	template: `<div>
	<select class="ui dropdown" style="max-width: 10em" name="clusterId" v-model="clusterId">
		<option value="0">[选择集群]</option>
		<option v-for="cluster in clusters" :value="cluster.id">{{cluster.name}}</option>
	</select>
</div>`
})

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
					newConnectionsRateBlockTimeout: 0,
					newConnectionsSecondlyRate: 0,
					newConnectionSecondlyRateBlockTimeout: 0,
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
				newConnectionsRateBlockTimeout: 0,
				newConnectionsSecondlyRate: 0,
				newConnectionSecondlyRateBlockTimeout: 0,
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
			<td class="title">启用DDoS防护</td>
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
			<td>单IP TCP新连接速率<em>（分钟）</em></td>
			<td>
				<div class="ui fields inline">
					<div class="ui field">
						<div class="ui input right labeled">
							<digit-input name="tcpNewConnectionsRate" v-model="config.tcp.newConnectionsRate" maxlength="6" size="6" style="width: 6em" :min="defaultConfigs.tcpNewConnectionsMinRate"></digit-input>
							<span class="ui label">个新连接/每分钟</span>
						</div>
					</div>
					<div class="ui field" style="line-height: 2.4em">
						屏蔽
					</div>
					<div class="ui field">
						<div class="ui input right labeled">
							<digit-input name="tcpNewConnectionsRateBlockTimeout" v-model="config.tcp.newConnectionsRateBlockTimeout" maxlength="6" size="6" style="width: 5em"></digit-input>
							<span class="ui label">秒</span>
						</div>
					</div>
				</div>
				
				<p class="comment">单个IP每分钟可以创建TCP新连接的数量。如果为0，则默认为{{defaultConfigs.tcpNewConnectionsMinutelyRate}}；最小值为{{defaultConfigs.tcpNewConnectionsMinMinutelyRate}}。如果没有填写屏蔽时间，则只丢弃数据包。</p>
			</td>
		</tr>
		<tr>
			<td>单IP TCP新连接速率<em>（秒钟）</em></td>
			<td>
				<div class="ui fields inline">
					<div class="ui field">
						<div class="ui input right labeled">
							<digit-input name="tcpNewConnectionsSecondlyRate" v-model="config.tcp.newConnectionsSecondlyRate" maxlength="6" size="6" style="width: 6em" :min="defaultConfigs.tcpNewConnectionsMinRate"></digit-input>
							<span class="ui label">个新连接/每秒钟</span>
						</div>
					</div>
					<div class="ui field" style="line-height: 2.4em">
						屏蔽
					</div>
					<div class="ui field">
						<div class="ui input right labeled">
							<digit-input name="tcpNewConnectionsSecondlyRateBlockTimeout" v-model="config.tcp.newConnectionsSecondlyRateBlockTimeout" maxlength="6" size="6" style="width: 5em"></digit-input>
							<span class="ui label">秒</span>
						</div>
					</div>
				</div>
				
				<p class="comment">单个IP每秒钟可以创建TCP新连接的数量。如果为0，则默认为{{defaultConfigs.tcpNewConnectionsSecondlyRate}}；最小值为{{defaultConfigs.tcpNewConnectionsMinSecondlyRate}}。如果没有填写屏蔽时间，则只丢弃数据包。</p>
			</td>
		</tr>
		<tr>
			<td>TCP端口列表</td>
			<td>
				<ddos-protection-ports-config-box :v-ports="config.tcp.ports" @change="changeTCPPorts"></ddos-protection-ports-config-box>
				<p class="comment">在这些端口上使用当前配置。默认为80和443两个端口。</p>
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

Vue.component("ddos-protection-ip-list-config-box", {
	props: ["v-ip-list"],
	data: function () {
		let list = this.vIpList
		if (list == null) {
			list = []
		}
		return {
			list: list,
			isAdding: false,
			addingIP: {
				ip: "",
				description: ""
			}
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingIPInput.focus()
			})
		},
		confirm: function () {
			let ip = this.addingIP.ip
			if (ip.length == 0) {
				this.warn("请输入IP")
				return
			}

			let exists = false
			this.list.forEach(function (v) {
				if (v.ip == ip) {
					exists = true
				}
			})
			if (exists) {
				this.warn("IP '" + ip + "'已经存在")
				return
			}

			let that = this
			Tea.Vue.$post("/ui/validateIPs")
				.params({
					ips: [ip]
				})
				.success(function () {
					that.list.push({
						ip: ip,
						description: that.addingIP.description
					})
					that.notifyChange()
					that.cancel()
				})
				.fail(function () {
					that.warn("请输入正确的IP")
				})
		},
		cancel: function () {
			this.isAdding = false
			this.addingIP = {
				ip: "",
				description: ""
			}
		},
		remove: function (index) {
			this.list.$remove(index)
			this.notifyChange()
		},
		warn: function (message) {
			let that = this
			teaweb.warn(message, function () {
				that.$refs.addingIPInput.focus()
			})
		},
		notifyChange: function () {
			this.$emit("change", this.list)
		}
	},
	template: `<div>
	<div v-if="list.length > 0">
		<div class="ui label basic tiny" v-for="(ipConfig, index) in list">
			{{ipConfig.ip}} <span class="grey small" v-if="ipConfig.description.length > 0">（{{ipConfig.description}}）</span> <a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">IP</span>
					<input type="text" v-model="addingIP.ip" ref="addingIPInput" maxlength="40" size="20" placeholder="IP" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">备注</span>
					<input type="text" v-model="addingIP.description" maxlength="10" size="10" placeholder="备注（可选）" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp;<a href="" @click.prevent="cancel()">取消</a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("node-cluster-combo-box", {
	props: ["v-cluster-id"],
	data: function () {
		let that = this
		Tea.action("/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})
		return {
			clusters: []
		}
	},
	methods: {
		change: function (item) {
			if (item == null) {
				this.$emit("change", 0)
			} else {
				this.$emit("change", item.value)
			}
		}
	},
	template: `<div v-if="clusters.length > 0" style="min-width: 10.4em">
	<combo-box title="集群" placeholder="集群名称" :v-items="clusters" name="clusterId" :v-value="vClusterId" @change="change"></combo-box>
</div>`
})

// 一个节点的多个集群选择器
Vue.component("node-clusters-selector", {
	props: ["v-primary-cluster", "v-secondary-clusters"],
	data: function () {
		let primaryCluster = this.vPrimaryCluster

		let secondaryClusters = this.vSecondaryClusters
		if (secondaryClusters == null) {
			secondaryClusters = []
		}

		return {
			primaryClusterId: (primaryCluster == null) ? 0 : primaryCluster.id,
			secondaryClusterIds: secondaryClusters.map(function (v) {
				return v.id
			}),

			primaryCluster: primaryCluster,
			secondaryClusters: secondaryClusters
		}
	},
	methods: {
		addPrimary: function () {
			let that = this
			let selectedClusterIds = [this.primaryClusterId].concat(this.secondaryClusterIds)
			teaweb.popup("/clusters/selectPopup?selectedClusterIds=" + selectedClusterIds.join(",") + "&mode=single", {
				height: "30em",
				width: "50em",
				callback: function (resp) {
					if (resp.data.cluster != null) {
						that.primaryCluster = resp.data.cluster
						that.primaryClusterId = that.primaryCluster.id
						that.notifyChange()
					}
				}
			})
		},
		removePrimary: function () {
			this.primaryClusterId = 0
			this.primaryCluster = null
			this.notifyChange()
		},
		addSecondary: function () {
			let that = this
			let selectedClusterIds = [this.primaryClusterId].concat(this.secondaryClusterIds)
			teaweb.popup("/clusters/selectPopup?selectedClusterIds=" + selectedClusterIds.join(",") + "&mode=multiple", {
				height: "30em",
				width: "50em",
				callback: function (resp) {
					if (resp.data.cluster != null) {
						that.secondaryClusterIds.push(resp.data.cluster.id)
						that.secondaryClusters.push(resp.data.cluster)
						that.notifyChange()
					}
				}
			})
		},
		removeSecondary: function (index) {
			this.secondaryClusterIds.$remove(index)
			this.secondaryClusters.$remove(index)
			this.notifyChange()
		},
		notifyChange: function () {
			this.$emit("change", {
				clusterId: this.primaryClusterId
			})
		}
	},
	template: `<div>
	<input type="hidden" name="primaryClusterId" :value="primaryClusterId"/>
	<input type="hidden" name="secondaryClusterIds" :value="JSON.stringify(secondaryClusterIds)"/>
	<table class="ui table">
		<tr>
			<td class="title">主集群</td>
			<td>
				<div v-if="primaryCluster != null">
					<div class="ui label basic small">{{primaryCluster.name}} &nbsp; <a href="" title="删除" @click.prevent="removePrimary"><i class="icon remove small"></i></a> </div>
				</div>
				<div style="margin-top: 0.6em" v-if="primaryClusterId == 0">
					<button class="ui button tiny" type="button" @click.prevent="addPrimary">+</button>
				</div>
				<p class="comment">多个集群配置有冲突时，优先使用主集群配置。</p>
			</td>
		</tr>
		<tr>
			<td>从集群</td>
			<td>
				<div v-if="secondaryClusters.length > 0">
					<div class="ui label basic small" v-for="(cluster, index) in secondaryClusters"><span class="grey">{{cluster.name}}</span> &nbsp; <a href="" title="删除" @click.prevent="removeSecondary(index)"><i class="icon remove small"></i></a> </div>
				</div>
				<div style="margin-top: 0.6em">
					<button class="ui button tiny" type="button" @click.prevent="addSecondary">+</button>
				</div>
			</td>
		</tr>
	</table>
</div>`
})

Vue.component("message-media-selector", {
    props: ["v-media-type"],
    mounted: function () {
        let that = this
        Tea.action("/admins/recipients/mediaOptions")
            .post()
            .success(function (resp) {
                that.medias = resp.data.medias

                // 初始化简介
                if (that.mediaType.length > 0) {
                    let media = that.medias.$find(function (_, media) {
                        return media.type == that.mediaType
                    })
                    if (media != null) {
                        that.description = media.description
                    }
                }
            })
    },
    data: function () {
        let mediaType = this.vMediaType
        if (mediaType == null) {
            mediaType = ""
        }
        return {
            medias: [],
            description: "",
            mediaType: mediaType
        }
    },
    watch: {
        mediaType: function (v) {
            let media = this.medias.$find(function (_, media) {
                return media.type == v
            })
            if (media == null) {
                this.description = ""
            } else {
                this.description = media.description
            }
            this.$emit("change", media)
        },
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="mediaType" v-model="mediaType">
        <option value="">[选择媒介类型]</option>
        <option v-for="media in medias" :value="media.type">{{media.name}}</option>
    </select>
    <p class="comment" v-html="description"></p>
</div>`
})

// 消息接收人设置
Vue.component("message-receivers-box", {
	props: ["v-node-cluster-id"],
	mounted: function () {
		let that = this
		Tea.action("/clusters/cluster/settings/message/selectedReceivers")
			.params({
				clusterId: this.clusterId
			})
			.post()
			.success(function (resp) {
				that.receivers = resp.data.receivers
			})
	},
	data: function () {
		let clusterId = this.vNodeClusterId
		if (clusterId == null) {
			clusterId = 0
		}
		return {
			clusterId: clusterId,
			receivers: []
		}
	},
	methods: {
		addReceiver: function () {
			let that = this
			let recipientIdStrings = []
			let groupIdStrings = []
			this.receivers.forEach(function (v) {
				if (v.type == "recipient") {
					recipientIdStrings.push(v.id.toString())
				} else if (v.type == "group") {
					groupIdStrings.push(v.id.toString())
				}
			})

			teaweb.popup("/clusters/cluster/settings/message/selectReceiverPopup?recipientIds=" + recipientIdStrings.join(",") + "&groupIds=" + groupIdStrings.join(","), {
				callback: function (resp) {
					that.receivers.push(resp.data)
				}
			})
		},
		removeReceiver: function (index) {
			this.receivers.$remove(index)
		}
	},
	template: `<div>
        <input type="hidden" name="receiversJSON" :value="JSON.stringify(receivers)"/>           
        <div v-if="receivers.length > 0">
            <div v-for="(receiver, index) in receivers" class="ui label basic small">
               <span v-if="receiver.type == 'group'">分组：</span>{{receiver.name}} <span class="grey small" v-if="receiver.subName != null && receiver.subName.length > 0">({{receiver.subName}})</span> &nbsp; <a href="" title="删除" @click.prevent="removeReceiver(index)"><i class="icon remove"></i></a>
            </div>
             <div class="ui divider"></div>
        </div>
      <button type="button" class="ui button tiny" @click.prevent="addReceiver">+</button>
</div>`
})

Vue.component("message-recipient-group-selector", {
    props: ["v-groups"],
    data: function () {
        let groups = this.vGroups
        if (groups == null) {
            groups = []
        }
        let groupIds = []
        if (groups.length > 0) {
            groupIds = groups.map(function (v) {
                return v.id.toString()
            }).join(",")
        }

        return {
            groups: groups,
            groupIds: groupIds
        }
    },
    methods: {
        addGroup: function () {
            let that = this
            teaweb.popup("/admins/recipients/groups/selectPopup?groupIds=" + this.groupIds, {
                callback: function (resp) {
                    that.groups.push(resp.data.group)
                    that.update()
                }
            })
        },
        removeGroup: function (index) {
            this.groups.$remove(index)
            this.update()
        },
        update: function () {
            let groupIds = []
            if (this.groups.length > 0) {
                this.groups.forEach(function (v) {
                    groupIds.push(v.id)
                })
            }
            this.groupIds = groupIds.join(",")
        }
    },
    template: `<div>
    <input type="hidden" name="groupIds" :value="groupIds"/>
    <div v-if="groups.length > 0">
        <div>
            <div v-for="(group, index) in groups" class="ui label small basic">
                {{group.name}} &nbsp; <a href="" title="删除" @click.prevent="removeGroup(index)"><i class="icon remove"></i></a>
            </div>
        </div>
        <div class="ui divider"></div>
    </div>   
    <button class="ui button tiny" type="button" @click.prevent="addGroup()">+</button>
</div>`
})

Vue.component("message-media-instance-selector", {
    props: ["v-instance-id"],
    mounted: function () {
        let that = this
        Tea.action("/admins/recipients/instances/options")
            .post()
            .success(function (resp) {
                that.instances = resp.data.instances

                // 初始化简介
                if (that.instanceId > 0) {
                    let instance = that.instances.$find(function (_, instance) {
                        return instance.id == that.instanceId
                    })
                    if (instance != null) {
                        that.description = instance.description
                        that.update(instance.id)
                    }
                }
            })
    },
    data: function () {
        let instanceId = this.vInstanceId
        if (instanceId == null) {
            instanceId = 0
        }
        return {
            instances: [],
            description: "",
            instanceId: instanceId
        }
    },
    watch: {
        instanceId: function (v) {
            this.update(v)
        }
    },
    methods: {
        update: function (v) {
            let instance = this.instances.$find(function (_, instance) {
                return instance.id == v
            })
            if (instance == null) {
                this.description = ""
            } else {
                this.description = instance.description
            }
            this.$emit("change", instance)
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="instanceId" v-model="instanceId">
        <option value="0">[选择媒介]</option>
        <option v-for="instance in instances" :value="instance.id">{{instance.name}} ({{instance.media.name}})</option>
    </select>
    <p class="comment" v-html="description"></p>
</div>`
})

Vue.component("message-row", {
	props: ["v-message", "v-can-close"],
	data: function () {
		let paramsJSON = this.vMessage.params
		let params = null
		if (paramsJSON != null && paramsJSON.length > 0) {
			params = JSON.parse(paramsJSON)
		}

		return {
			message: this.vMessage,
			params: params,
			isClosing: false
		}
	},
	methods: {
		viewCert: function (certId) {
			teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
				height: "28em",
				width: "48em"
			})
		},
		readMessage: function (messageId) {
			let that = this

			Tea.action("/messages/readPage")
				.params({"messageIds": [messageId]})
				.post()
				.success(function () {
					// 刷新父级页面Badge
					if (window.parent.Tea != null && window.parent.Tea.Vue != null) {
						window.parent.Tea.Vue.checkMessagesOnce()
					}

					// 刷新当前页面
					if (that.vCanClose && typeof (NotifyPopup) != "undefined") {
						that.isClosing = true
						setTimeout(function () {
							NotifyPopup({})
						}, 1000)
					} else {
						teaweb.reload()
					}
				})
		}
	},
	template: `<div>
<table class="ui table selectable" v-if="!isClosing">
	<tr :class="{error: message.level == 'error', positive: message.level == 'success', warning: message.level == 'warning'}">
		<td style="position: relative">
			<strong>{{message.datetime}}</strong>
			<span v-if="message.cluster != null && message.cluster.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster?clusterId=' + message.cluster.id" target="_top" v-if="message.role == 'node'">集群：{{message.cluster.name}}</a>
				<a :href="'/ns/clusters/cluster?clusterId=' + message.cluster.id" target="_top" v-if="message.role == 'dns'">DNS集群：{{message.cluster.name}}</a>
			</span>
			<span v-if="message.node != null && message.node.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id" target="_top" v-if="message.role == 'node'">节点：{{message.node.name}}</a>
				<a :href="'/ns/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id" target="_top" v-if="message.role == 'dns'">DNS节点：{{message.node.name}}</a>
			</span>
			<a href=""  style="position: absolute; right: 1em" @click.prevent="readMessage(message.id)" title="标为已读"><i class="icon check"></i></a>
		</td>
	</tr>
	<tr :class="{error: message.level == 'error', positive: message.level == 'success', warning: message.level == 'warning'}">
		<td>
			<pre style="padding: 0; margin:0; word-break: break-all;">{{message.body}}</pre>
			
			<!-- 健康检查 -->
			<div v-if="message.type == 'HealthCheckFailed'" style="margin-top: 0.8em">
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + param.node.id" v-for="param in params" class="ui label small basic" style="margin-bottom: 0.5em" target="_top">{{param.node.name}}: {{param.error}}</a>
			</div>
			
			<!-- 集群DNS设置 -->
			<div v-if="message.type == 'ClusterDNSSyncFailed'" style="margin-top: 0.8em">
				<a :href="'/dns/clusters/cluster?clusterId=' + message.cluster.id" target="_top">查看问题 &raquo;</a>
			</div>
			
			<!-- 证书即将过期 -->
			<div v-if="message.type == 'SSLCertExpiring'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a><span v-if="params != null && params.acmeTaskId > 0"> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" target="_top">查看任务&raquo;</a></span>
			</div>
			
			<!-- 证书续期成功 -->
			<div v-if="message.type == 'SSLCertACMETaskSuccess'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0" target="_top">查看任务&raquo;</a>
			</div>
			
			<!-- 证书续期失败 -->
			<div v-if="message.type == 'SSLCertACMETaskFailed'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0" target="_top">查看任务&raquo;</a>
			</div>
			
			<!-- 网站域名审核 -->
			<div v-if="message.type == 'serverNamesRequireAuditing'" style="margin-top: 0.8em">
				<a :href="'/servers/server/settings/serverNames?serverId=' + params.serverId" target="_top">去审核</a></a>
			</div>

			<!-- 节点调度 -->
			<div v-if="message.type == 'NodeSchedule'" style="margin-top: 0.8em">
				<a :href="'/clusters/cluster/node/settings/schedule?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id" target="_top">查看调度状态 &raquo;</a>
			</div>
			
			<!-- 节点租期结束 -->
			<div v-if="message.type == 'NodeOfflineDay'" style="margin-top: 0.8em">
				<a :href="'/clusters/cluster/node/detail?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id" target="_top">查看详情 &raquo;</a>
			</div>
		</td>
	</tr>
</table>
<div class="margin"></div>
</div>`
})

Vue.component("ns-domain-group-selector", {
	props: ["v-domain-group-id"],
	data: function () {
		let groupId = this.vDomainGroupId
		if (groupId == null) {
			groupId = 0
		}
		return {
			userId: 0,
			groupId: groupId
		}
	},
	methods: {
		change: function (group) {
			if (group != null) {
				this.$emit("change", group.id)
			} else {
				this.$emit("change", 0)
			}
		},
		reload: function (userId) {
			this.userId = userId
			this.$refs.comboBox.clear()
			this.$refs.comboBox.setDataURL("/ns/domains/groups/options?userId=" + userId)
			this.$refs.comboBox.reloadData()
		}
	},
	template: `<div>
	<combo-box 
		data-url="/ns/domains/groups/options" 
		placeholder="选择分组" 
		data-key="groups" 
		name="groupId"
		:v-value="groupId" 
		@change="change"
		ref="comboBox">	
	</combo-box>
</div>`
})

// 选择多个线路
Vue.component("ns-routes-selector", {
	props: ["v-routes", "name"],
	mounted: function () {
		let that = this
		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.routes = resp.data.routes

				// provinces
				let provinces = {}
				if (resp.data.provinces != null && resp.data.provinces.length > 0) {
					for (const province of resp.data.provinces) {
						let countryCode = province.countryCode
						if (typeof provinces[countryCode] == "undefined") {
							provinces[countryCode] = []
						}
						provinces[countryCode].push({
							name: province.name,
							code: province.code
						})
					}
				}
				that.provinces = provinces
			})
	},
	data: function () {
		let selectedRoutes = this.vRoutes
		if (selectedRoutes == null) {
			selectedRoutes = []
		}

		let inputName = this.name
		if (typeof inputName != "string" || inputName.length == 0) {
			inputName = "routeCodes"
		}

		return {
			routeCode: "default",
			inputName: inputName,
			routes: [],

			provinces: {}, // country code => [ province1, province2, ... ]
			provinceRouteCode: "",

			isAdding: false,
			routeType: "default",
			selectedRoutes: selectedRoutes,
		}
	},
	watch: {
		routeType: function (v) {
			this.routeCode = ""
			let that = this
			this.routes.forEach(function (route) {
				if (route.type == v && that.routeCode.length == 0) {
					that.routeCode = route.code
				}
			})
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.routeType = "default"
			this.routeCode = "default"
			this.provinceRouteCode = ""
			this.$emit("add")
		},
		cancel: function () {
			this.isAdding = false
			this.$emit("cancel")
		},
		confirm: function () {
			if (this.routeCode.length == 0) {
				return
			}

			let that = this

			// route
			let selectedRoute = null
			for (const route of this.routes) {
				if (route.code == this.routeCode) {
					selectedRoute = route
					break
				}
			}

			if (selectedRoute != null) {
				// province route
				if (this.provinceRouteCode.length > 0 && this.provinces[this.routeCode] != null) {
					for (const province of this.provinces[this.routeCode]) {
						if (province.code == this.provinceRouteCode) {
							selectedRoute = {
								name: selectedRoute.name + "-" + province.name,
								code: province.code
							}
							break
						}
					}
				}

				that.selectedRoutes.push(selectedRoute)
			}

			this.$emit("change", this.selectedRoutes)
			this.cancel()
		},
		remove: function (index) {
			this.selectedRoutes.$remove(index)
			this.$emit("change", this.selectedRoutes)
		}
	}
	,
	template: `<div>
	<div v-show="selectedRoutes.length > 0">
		<div class="ui label basic text small" v-for="(route, index) in selectedRoutes" style="margin-bottom: 0.3em">
			<input type="hidden" :name="inputName" :value="route.code"/>
			{{route.name}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding" style="margin-bottom: 1em">
		<table class="ui table">
			<tr>
				<td class="title">选择类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="routeType">
						<option value="default">[默认线路]</option>
						<option value="user">自定义线路</option>
						<option value="isp">运营商</option>
						<option value="china">中国省市</option>
						<option value="world">全球国家地区</option>
						<option value="agent">搜索引擎</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>选择线路 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="routeCode">
						<option v-for="route in routes" :value="route.code" v-if="route.type == routeType">{{route.name}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="routeCode.length > 0 && provinces[routeCode] != null">
				<td>选择省/州</td>
				<td>
					<select class="ui dropdown auto-width" v-model="provinceRouteCode">
						<option value="">[全域]</option>
						<option v-for="province in provinces[routeCode]" :value="province.code">{{province.name}}</option>
					</select>
				</td>
			</tr>
		</table>
		<div>
			<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
			&nbsp; <a href="" title="取消" @click.prevent="cancel">取消</a>
		</div>	
	</div>
	<button class="ui button tiny" type="button" @click.prevent="add" v-if="!isAdding">+</button>
</div>`
})

// 递归DNS设置
Vue.component("ns-recursion-config-box", {
	props: ["v-recursion-config"],
	data: function () {
		let recursion = this.vRecursionConfig
		if (recursion == null) {
			recursion = {
				isOn: false,
				hosts: [],
				allowDomains: [],
				denyDomains: [],
				useLocalHosts: false
			}
		}
		if (recursion.hosts == null) {
			recursion.hosts = []
		}
		if (recursion.allowDomains == null) {
			recursion.allowDomains = []
		}
		if (recursion.denyDomains == null) {
			recursion.denyDomains = []
		}
		return {
			config: recursion,
			hostIsAdding: false,
			host: "",
			updatingHost: null
		}
	},
	methods: {
		changeHosts: function (hosts) {
			this.config.hosts = hosts
		},
		changeAllowDomains: function (domains) {
			this.config.allowDomains = domains
		},
		changeDenyDomains: function (domains) {
			this.config.denyDomains = domains
		},
		removeHost: function (index) {
			this.config.hosts.$remove(index)
		},
		addHost: function () {
			this.updatingHost = null
			this.host = ""
			this.hostIsAdding = !this.hostIsAdding
			if (this.hostIsAdding) {
				var that = this
				setTimeout(function () {
					let hostRef = that.$refs.hostRef
					if (hostRef != null) {
						hostRef.focus()
					}
				}, 200)
			}
		},
		updateHost: function (host) {
			this.updatingHost = host
			this.host = host.host
			this.hostIsAdding = !this.hostIsAdding

			if (this.hostIsAdding) {
				var that = this
				setTimeout(function () {
					let hostRef = that.$refs.hostRef
					if (hostRef != null) {
						hostRef.focus()
					}
				}, 200)
			}
		},
		confirmHost: function () {
			if (this.host.length == 0) {
				teaweb.warn("请输入DNS地址")
				return
			}

			// TODO 校验Host
			// TODO 可以输入端口号
			// TODO 可以选择协议

			this.hostIsAdding = false
			if (this.updatingHost == null) {
				this.config.hosts.push({
					host: this.host
				})
			} else {
				this.updatingHost.host = this.host
			}
		},
		cancelHost: function () {
			this.hostIsAdding = false
		}
	},
	template: `<div>
	<input type="hidden" name="recursionJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="isOn" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">启用后，如果找不到某个域名的解析记录，则向上一级DNS查找。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>从节点本机读取<br/>上级DNS主机</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="useLocalHosts" value="1" v-model="config.useLocalHosts"/>
						<label></label>
					</div>
					<p class="comment">选中后，节点会试图从<code-label>/etc/resolv.conf</code-label>文件中读取DNS配置。 </p>
				</td>
			</tr>
			<tr v-show="!config.useLocalHosts">
				<td>上级DNS主机地址 *</td>
				<td>
					<div v-if="config.hosts.length > 0">
						<div v-for="(host, index) in config.hosts" class="ui label tiny basic">
							{{host.host}} &nbsp;
							<a href="" title="修改" @click.prevent="updateHost(host)"><i class="icon pencil tiny"></i></a>
							<a href="" title="删除" @click.prevent="removeHost(index)"><i class="icon remove small"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<div v-if="hostIsAdding">
						<div class="ui fields inline">
							<div class="ui field">
								<input type="text" placeholder="DNS主机地址" v-model="host" ref="hostRef" @keyup.enter="confirmHost" @keypress.enter.prevent="1"/>
							</div>
							<div class="ui field">
								<button class="ui button tiny" type="button" @click.prevent="confirmHost">确认</button> &nbsp; <a href="" title="取消" @click.prevent="cancelHost"><i class="icon remove small"></i></a>
							</div>
						</div>
					</div>
					<div style="margin-top: 0.5em">
						<button type="button" class="ui button tiny" @click.prevent="addHost">+</button>
					</div>
				</td>
			</tr>
			<tr>
				<td>允许的域名</td>
				<td><values-box name="allowDomains" :values="config.allowDomains" @change="changeAllowDomains"></values-box>
					<p class="comment">支持星号通配符，比如<code-label>*.example.org</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>不允许的域名</td>
				<td>
					<values-box name="denyDomains" :values="config.denyDomains" @change="changeDenyDomains"></values-box>
					<p class="comment">支持星号通配符，比如<code-label>*.example.org</code-label>。优先级比允许的域名高。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("ns-access-log-ref-box", {
	props: ["v-access-log-ref", "v-is-parent"],
	data: function () {
		let config = this.vAccessLogRef
		if (config == null) {
			config = {
				isOn: false,
				isPrior: false,
				logMissingDomains: false
			}
		}
		if (typeof (config.logMissingDomains) == "undefined") {
			config.logMissingDomains = false
		}
		return {
			config: config
		}
	},
	template: `<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="!vIsParent"></prior-checkbox>
		<tbody v-show="vIsParent || config.isPrior">
			<tr>
				<td class="title">启用</td>
				<td>
					<checkbox name="isOn" value="1" v-model="config.isOn"></checkbox>
				</td>
			</tr>
			<tr>
				<td>只记录失败查询</td>
				<td>
					<checkbox v-model="config.missingRecordsOnly"></checkbox>
					<p class="comment">选中后，表示只记录查询失败的日志。</p>
				</td>
			</tr>
			<tr>
				<td>包含未添加的域名</td>
				<td>
					<checkbox name="logMissingDomains" value="1" v-model="config.logMissingDomains"></checkbox>
					<p class="comment">选中后，表示日志中包含对没有在系统里创建的域名访问。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("ns-records-health-check-config-box", {
	props:["value"],
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isOn: false,
				port: 80,
				timeoutSeconds: 5,
				countUp: 1,
				countDown: 3
			}
		}
		return {
			config: config,
			portString: config.port.toString(),
			timeoutSecondsString: config.timeoutSeconds.toString(),
			countUpString: config.countUp.toString(),
			countDownString: config.countDown.toString()
		}
	},
	watch: {
		portString: function (value) {
			let port = parseInt(value.toString())
			if (isNaN(port) || port > 65535 || port < 1) {
				this.config.port = 80
			} else {
				this.config.port = port
			}
		},
		timeoutSecondsString: function (value) {
			let timeoutSeconds = parseInt(value.toString())
			if (isNaN(timeoutSeconds) || timeoutSeconds > 1000 || timeoutSeconds < 1) {
				this.config.timeoutSeconds = 5
			} else {
				this.config.timeoutSeconds = timeoutSeconds
			}
		},
		countUpString: function (value) {
			let countUp = parseInt(value.toString())
			if (isNaN(countUp) || countUp > 1000 || countUp < 1) {
				this.config.countUp = 1
			} else {
				this.config.countUp = countUp
			}
		},
		countDownString: function (value) {
			let countDown = parseInt(value.toString())
			if (isNaN(countDown) || countDown > 1000 || countDown < 1) {
				this.config.countDown = 3
			} else {
				this.config.countDown = countDown
			}
		}
	},
	template: `<div>
	<input type="hidden" name="recordsHealthCheckJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">启用健康检查</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
					<p class="comment">选中后，表示启用当前域名下A/AAAA记录的健康检查；启用此设置后，你仍需设置单个A/AAAA记录的健康检查。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>默认检测端口</td>
				<td>
					<input type="text" v-model="portString" maxlength="5" style="width: 5em"/>
					<p class="comment">通过尝试连接A/AAAA记录中的IP加此端口来确定当前记录是否健康。</p>
				</td>
			</tr>
			<tr>
				<td>默认超时时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 4em" v-model="timeoutSecondsString" maxlength="3"/>
						<span class="ui label">秒</span>
					</div>
				</td>
			</tr>
			<tr>
				<td>默认连续上线次数</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 4em" v-model="countUpString" maxlength="3"/>
						<span class="ui label">次</span>
					</div>
					<p class="comment">连续检测<span v-if="config.countUp > 0">{{config.countUp}}</span><span v-else>N</span>次成功后，认为当前记录是在线的。</p>
				</td>
			</tr>
			<tr>
				<td>默认连续下线次数</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 4em" v-model="countDownString" maxlength="3"/>
						<span class="ui label">次</span>
					</div>
					<p class="comment">连续检测<span v-if="config.countDown > 0">{{config.countDown}}</span><span v-else>N</span>次失败后，认为当前记录是离线的。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("ns-node-ddos-protection-config-box", {
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
					newConnectionsRateBlockTimeout: 0,
					newConnectionsSecondlyRate: 0,
					newConnectionSecondlyRateBlockTimeout: 0,
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
				newConnectionsRateBlockTimeout: 0,
				newConnectionsSecondlyRate: 0,
				newConnectionSecondlyRateBlockTimeout: 0,
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
			<td>单IP TCP新连接速率<em>（分钟）</em></td>
			<td>
				<div class="ui fields inline">
					<div class="ui field">
						<div class="ui input right labeled">
							<digit-input name="tcpNewConnectionsRate" v-model="config.tcp.newConnectionsRate" maxlength="6" size="6" style="width: 6em" :min="defaultConfigs.tcpNewConnectionsMinRate"></digit-input>
							<span class="ui label">个新连接/每分钟</span>
						</div>
					</div>
					<div class="ui field" style="line-height: 2.4em">
						屏蔽
					</div>
					<div class="ui field">
						<div class="ui input right labeled">
							<digit-input name="tcpNewConnectionsRateBlockTimeout" v-model="config.tcp.newConnectionsRateBlockTimeout" maxlength="6" size="6" style="width: 5em"></digit-input>
							<span class="ui label">秒</span>
						</div>
					</div>
				</div>
				
				<p class="comment">单个IP每分钟可以创建TCP新连接的数量。如果为0，则默认为{{defaultConfigs.tcpNewConnectionsMinutelyRate}}；最小值为{{defaultConfigs.tcpNewConnectionsMinMinutelyRate}}。如果没有填写屏蔽时间，则只丢弃数据包。</p>
			</td>
		</tr>
		<tr>
			<td>单IP TCP新连接速率<em>（秒钟）</em></td>
			<td>
				<div class="ui fields inline">
					<div class="ui field">
						<div class="ui input right labeled">
							<digit-input name="tcpNewConnectionsSecondlyRate" v-model="config.tcp.newConnectionsSecondlyRate" maxlength="6" size="6" style="width: 6em" :min="defaultConfigs.tcpNewConnectionsMinRate"></digit-input>
							<span class="ui label">个新连接/每秒钟</span>
						</div>
					</div>
					<div class="ui field" style="line-height: 2.4em">
						屏蔽
					</div>
					<div class="ui field">
						<div class="ui input right labeled">
							<digit-input name="tcpNewConnectionsSecondlyRateBlockTimeout" v-model="config.tcp.newConnectionsSecondlyRateBlockTimeout" maxlength="6" size="6" style="width: 5em"></digit-input>
							<span class="ui label">秒</span>
						</div>
					</div>
				</div>
				
				<p class="comment">单个IP每秒钟可以创建TCP新连接的数量。如果为0，则默认为{{defaultConfigs.tcpNewConnectionsSecondlyRate}}；最小值为{{defaultConfigs.tcpNewConnectionsMinSecondlyRate}}。如果没有填写屏蔽时间，则只丢弃数据包。</p>
			</td>
		</tr>
		<tr>
			<td>TCP端口列表</td>
			<td>
				<ddos-protection-ports-config-box :v-ports="config.tcp.ports" @change="changeTCPPorts"></ddos-protection-ports-config-box>
				<p class="comment">在这些端口上使用当前配置。默认为53端口。</p>
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

Vue.component("ns-route-ranges-box", {
	props: ["v-ranges"],
	data: function () {
		let ranges = this.vRanges
		if (ranges == null) {
			ranges = []
		}
		return {
			ranges: ranges,
			isAdding: false,
			isAddingBatch: false,

			// 类型
			rangeType: "ipRange",
			isReverse: false,

			// IP范围
			ipRangeFrom: "",
			ipRangeTo: "",

			batchIPRange: "",

			// CIDR
			ipCIDR: "",
			batchIPCIDR: "",

			// region
			regions: [],
			regionType: "country",
			regionConnector: "OR"
		}
	},
	methods: {
		addIPRange: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.ipRangeFrom.focus()
			}, 100)
		},
		addCIDR: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.ipCIDR.focus()
			}, 100)
		},
		addRegions: function () {
			this.isAdding = true
		},
		addRegion: function (regionType) {
			this.regionType = regionType
		},
		remove: function (index) {
			this.ranges.$remove(index)
		},
		cancelIPRange: function () {
			this.isAdding = false
			this.ipRangeFrom = ""
			this.ipRangeTo = ""
			this.isReverse = false
		},
		cancelIPCIDR: function () {
			this.isAdding = false
			this.ipCIDR = ""
			this.isReverse = false
		},
		cancelRegions: function () {
			this.isAdding = false
			this.regions = []
			this.regionType = "country"
			this.regionConnector = "OR"
			this.isReverse = false
		},
		confirmIPRange: function () {
			// 校验IP
			let that = this
			this.ipRangeFrom = this.ipRangeFrom.trim()
			if (!this.validateIP(this.ipRangeFrom)) {
				teaweb.warn("开始IP填写错误", function () {
					that.$refs.ipRangeFrom.focus()
				})
				return
			}

			this.ipRangeTo = this.ipRangeTo.trim()
			if (!this.validateIP(this.ipRangeTo)) {
				teaweb.warn("结束IP填写错误", function () {
					that.$refs.ipRangeTo.focus()
				})
				return
			}

			this.ranges.push({
				type: "ipRange",
				params: {
					ipFrom: this.ipRangeFrom,
					ipTo: this.ipRangeTo,
					isReverse: this.isReverse
				}
			})
			this.cancelIPRange()
		},
		confirmIPCIDR: function () {
			let that = this
			if (this.ipCIDR.length == 0) {
				teaweb.warn("请填写CIDR", function () {
					that.$refs.ipCIDR.focus()
				})
				return
			}
			if (!this.validateCIDR(this.ipCIDR)) {
				teaweb.warn("请输入正确的CIDR", function () {
					that.$refs.ipCIDR.focus()
				})
				return
			}


			this.ranges.push({
				type: "cidr",
				params: {
					cidr: this.ipCIDR,
					isReverse: this.isReverse
				}
			})
			this.cancelIPCIDR()
		},
		confirmRegions: function () {
			if (this.regions.length == 0) {
				this.cancelRegions()
				return
			}
			this.ranges.push({
				type: "region",
				connector: this.regionConnector,
				params: {
					regions: this.regions,
					isReverse: this.isReverse
				}
			})
			this.cancelRegions()
		},
		addBatchIPRange: function () {
			this.isAddingBatch = true
			let that = this
			setTimeout(function () {
				that.$refs.batchIPRange.focus()
			}, 100)
		},
		addBatchCIDR: function () {
			this.isAddingBatch = true
			let that = this
			setTimeout(function () {
				that.$refs.batchIPCIDR.focus()
			}, 100)
		},
		cancelBatchIPRange: function () {
			this.isAddingBatch = false
			this.batchIPRange = ""
			this.isReverse = false
		},
		cancelBatchIPCIDR: function () {
			this.isAddingBatch = false
			this.batchIPCIDR = ""
			this.isReverse = false
		},
		confirmBatchIPRange: function () {
			let that = this
			let rangesText = this.batchIPRange
			if (rangesText.length == 0) {
				teaweb.warn("请填写要加入的IP范围", function () {
					that.$refs.batchIPRange.focus()
				})
				return
			}

			let validRanges = []
			let invalidLine = ""
			rangesText.split("\n").forEach(function (line) {
				line = line.trim()
				if (line.length == 0) {
					return
				}
				line = line.replace("，", ",")
				let pieces = line.split(",")
				if (pieces.length != 2) {
					invalidLine = line
					return
				}
				let ipFrom = pieces[0].trim()
				let ipTo = pieces[1].trim()
				if (!that.validateIP(ipFrom) || !that.validateIP(ipTo)) {
					invalidLine = line
					return
				}
				validRanges.push({
					type: "ipRange",
					params: {
						ipFrom: ipFrom,
						ipTo: ipTo,
						isReverse: that.isReverse
					}
				})
			})
			if (invalidLine.length > 0) {
				teaweb.warn("'" + invalidLine + "'格式错误", function () {
					that.$refs.batchIPRange.focus()
				})
				return
			}
			validRanges.forEach(function (v) {
				that.ranges.push(v)
			})
			this.cancelBatchIPRange()
		},
		confirmBatchIPCIDR: function () {
			let that = this
			let rangesText = this.batchIPCIDR
			if (rangesText.length == 0) {
				teaweb.warn("请填写要加入的CIDR", function () {
					that.$refs.batchIPCIDR.focus()
				})
				return
			}

			let validRanges = []
			let invalidLine = ""
			rangesText.split("\n").forEach(function (line) {
				let cidr = line.trim()
				if (cidr.length == 0) {
					return
				}
				if (!that.validateCIDR(cidr)) {
					invalidLine = line
					return
				}
				validRanges.push({
					type: "cidr",
					params: {
						cidr: cidr,
						isReverse: that.isReverse
					}
				})
			})
			if (invalidLine.length > 0) {
				teaweb.warn("'" + invalidLine + "'格式错误", function () {
					that.$refs.batchIPCIDR.focus()
				})
				return
			}
			validRanges.forEach(function (v) {
				that.ranges.push(v)
			})
			this.cancelBatchIPCIDR()
		},
		selectRegionCountry: function (country) {
			if (country == null) {
				return
			}
			this.regions.push({
				type: "country",
				id: country.id,
				name: country.name
			})
			this.$refs.regionCountryComboBox.clear()
		},
		selectRegionProvince: function (province) {
			if (province == null) {
				return
			}
			this.regions.push({
				type: "province",
				id: province.id,
				name: province.name
			})
			this.$refs.regionProvinceComboBox.clear()
		},
		selectRegionCity: function (city) {
			if (city == null) {
				return
			}
			this.regions.push({
				type: "city",
				id: city.id,
				name: city.name
			})
			this.$refs.regionCityComboBox.clear()
		},
		selectRegionProvider: function (provider) {
			if (provider == null) {
				return
			}
			this.regions.push({
				type: "provider",
				id: provider.id,
				name: provider.name
			})
			this.$refs.regionProviderComboBox.clear()
		},
		removeRegion: function (index) {
			this.regions.$remove(index)
		},
		validateIP: function (ip) {
			if (ip.length == 0) {
				return
			}

			// IPv6
			if (ip.indexOf(":") >= 0) {
				let pieces = ip.split(":")
				if (pieces.length > 8) {
					return false
				}
				let isOk = true
				pieces.forEach(function (piece) {
					if (!/^[\da-fA-F]{0,4}$/.test(piece)) {
						isOk = false
					}
				})

				return isOk
			}

			if (!ip.match(/^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/)) {
				return false
			}
			let pieces = ip.split(".")
			let isOk = true
			pieces.forEach(function (v) {
				let v1 = parseInt(v)
				if (v1 > 255) {
					isOk = false
				}
			})
			return isOk
		},
		validateCIDR: function (cidr) {
			let pieces = cidr.split("/")
			if (pieces.length != 2) {
				return false
			}
			let ip = pieces[0]
			if (!this.validateIP(ip)) {
				return false
			}
			let mask = pieces[1]
			if (!/^\d{1,3}$/.test(mask)) {
				return false
			}
			mask = parseInt(mask, 10)
			if (cidr.indexOf(":") >= 0) { // IPv6
				return mask <= 128
			}
			return mask <= 32
		},
		updateRangeType: function (rangeType) {
			this.rangeType = rangeType
		}
	},
	template: `<div>
	<input type="hidden" name="rangesJSON" :value="JSON.stringify(ranges)"/>
	<div v-if="ranges.length > 0">
		<div class="ui label tiny basic" v-for="(range, index) in ranges" style="margin-bottom: 0.3em">
			<span class="red" v-if="range.params.isReverse">[排除]</span>
			<span v-if="range.type == 'ipRange'">IP范围：</span>
			<span v-if="range.type == 'cidr'">CIDR：</span>
			<span v-if="range.type == 'region'"></span>
			<span v-if="range.type == 'ipRange'">{{range.params.ipFrom}} - {{range.params.ipTo}}</span>
			<span v-if="range.type == 'cidr'">{{range.params.cidr}}</span>
			<span v-if="range.type == 'region'">
				<span v-for="(region, index) in range.params.regions">
					<span v-if="region.type == 'country'">国家/地区</span>
					<span v-if="region.type == 'province'">省份</span>
					<span v-if="region.type == 'city'">城市</span>
					<span v-if="region.type == 'provider'">ISP</span>
					：{{region.name}}
					<span v-if="index < range.params.regions.length - 1" class="grey">
						&nbsp;
						<span v-if="range.connector == 'OR' || range.connector == '' || range.connector == null">或</span>
						<span v-if="range.connector == 'AND'">且</span>
						&nbsp;
					</span>
				</span>
			</span>
			 &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- IP范围 -->
	<div v-if="rangeType == 'ipRange'">
		<!-- 添加单个IP范围 -->
		<div style="margin-bottom: 1em" v-show="isAdding">
			<table class="ui table">
				<tr>
					<td class="title">开始IP *</td>
					<td>
						<input type="text" placeholder="开始IP" maxlength="40" size="40" style="width: 15em" v-model="ipRangeFrom" ref="ipRangeFrom"  @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
					</td>
				</tr>
				<tr>
					<td>结束IP *</td>
					<td>
						<input type="text" placeholder="结束IP" maxlength="40" size="40" style="width: 15em" v-model="ipRangeTo" ref="ipRangeTo" @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmIPRange">确定</button> &nbsp;
					<a href="" @click.prevent="cancelIPRange" title="取消"><i class="icon remove small"></i></a>
		</div>
	
		<!-- 添加多个IP范围 -->
		<div style="margin-bottom: 1em" v-show="isAddingBatch">
			<table class="ui table">
				<tr>
					<td class="title">IP范围列表 *</td>
					<td>
						<textarea rows="5" ref="batchIPRange" v-model="batchIPRange"></textarea>	
						<p class="comment">每行一条，格式为<code-label>开始IP,结束IP</code-label>，比如<code-label>192.168.1.100,192.168.1.200</code-label>。</p>	
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmBatchIPRange">确定</button> &nbsp;
				<a href="" @click.prevent="cancelBatchIPRange" title="取消"><i class="icon remove small"></i></a>
		</div>
		
		<div v-if="!isAdding && !isAddingBatch">
			<button class="ui button tiny" type="button" @click.prevent="addIPRange">添加单个IP范围</button> &nbsp;
			<button class="ui button tiny" type="button" @click.prevent="addBatchIPRange">批量添加IP范围</button>
		</div>
	</div>
	
	<!-- CIDR -->
	<div v-if="rangeType == 'cidr'">
		<!-- 添加单个IP范围 -->
		<div style="margin-bottom: 1em" v-show="isAdding">
			<table class="ui table">
				<tr>
					<td class="title">CIDR *</td>
					<td>
						<input type="text" placeholder="IP/MASK" maxlength="40" size="40" style="width: 15em" v-model="ipCIDR" ref="ipCIDR"  @keyup.enter="confirmIPCIDR" @keypress.enter.prevent="1"/>
						<p class="comment">类似于<code-label>192.168.2.1/24</code-label>。</p>
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmIPCIDR">确定</button> &nbsp;
					<a href="" @click.prevent="cancelIPCIDR" title="取消"><i class="icon remove small"></i></a>
		</div>
		
		<!-- 添加多个IP范围 -->
		<div style="margin-bottom: 1em" v-show="isAddingBatch">
			<table class="ui table">
				<tr>
					<td class="title">IP范围列表 *</td>
					<td>
						<textarea rows="5" ref="batchIPCIDR" v-model="batchIPCIDR"></textarea>	
						<p class="comment">每行一条，格式为<code-label>IP/MASK</code-label>，比如<code-label>192.168.2.1/24</code-label>。</p>	
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmBatchIPCIDR">确定</button> &nbsp;
				<a href="" @click.prevent="cancelBatchIPCIDR" title="取消"><i class="icon remove small"></i></a>
		</div>
		
		<div v-if="!isAdding && !isAddingBatch">
			<button class="ui button tiny" type="button" @click.prevent="addCIDR">添加单个CIDR</button> &nbsp;
			<button class="ui button tiny" type="button" @click.prevent="addBatchCIDR">批量添加CIDR</button>
		</div>
	</div>
	
	<!-- 区域 -->
	<div v-if="rangeType == 'region'">
		<!-- 添加区域 -->
		<div v-if="isAdding">
			<table class="ui table">
				<tr>
					<td>已添加</td>
					<td>
						<span v-for="(region, index) in regions">
							<span class="ui label small basic">
								<span v-if="region.type == 'country'">国家/地区</span>
								<span v-if="region.type == 'province'">省份</span>
								<span v-if="region.type == 'city'">城市</span>
								<span v-if="region.type == 'provider'">ISP</span>
								：{{region.name}} <a href="" title="删除" @click.prevent="removeRegion(index)"><i class="icon remove small"></i></a>
							</span>
							<span v-if="index < regions.length - 1" class="grey">
								&nbsp;
								<span v-if="regionConnector == 'OR' || regionConnector == ''">或</span>
								<span v-if="regionConnector == 'AND'">且</span>
								&nbsp;
							</span>
						</span>
					</td>
				</tr>
				<tr>
					<td class="title">添加新<span v-if="regionType == 'country'">国家/地区</span><span v-if="regionType == 'province'">省份</span><span v-if="regionType == 'city'">城市</span><span v-if="regionType == 'provider'">ISP</span>
					
					 *</td>
					<td>
					 	<!-- region country name -->
						<div v-if="regionType == 'country'">
							<combo-box title="" width="14em" data-url="/ui/countryOptions" data-key="countries" placeholder="点这里选择国家/地区" @change="selectRegionCountry" ref="regionCountryComboBox" key="combo-box-country"></combo-box>
						</div>
			
						<!-- region province name -->
						<div v-if="regionType == 'province'" >
							<combo-box title="" data-url="/ui/provinceOptions" data-key="provinces" placeholder="点这里选择省份" @change="selectRegionProvince" ref="regionProvinceComboBox" key="combo-box-province"></combo-box>
						</div>
			
						<!-- region city name -->
						<div v-if="regionType == 'city'" >
							<combo-box title="" data-url="/ui/cityOptions" data-key="cities" placeholder="点这里选择城市" @change="selectRegionCity" ref="regionCityComboBox" key="combo-box-city"></combo-box>
						</div>
			
						<!-- ISP Name -->
						<div v-if="regionType == 'provider'" >
							<combo-box title="" data-url="/ui/providerOptions" data-key="providers" placeholder="点这里选择ISP" @change="selectRegionProvider" ref="regionProviderComboBox" key="combo-box-isp"></combo-box>
						</div>
						
						<div style="margin-top: 1em">
							<button class="ui button tiny basic" :class="{blue: regionType == 'country'}" type="button" @click.prevent="addRegion('country')">添加国家/地区</button> &nbsp;
							<button class="ui button tiny basic" :class="{blue: regionType == 'province'}" type="button" @click.prevent="addRegion('province')">添加省份</button> &nbsp;
							<button class="ui button tiny basic" :class="{blue: regionType == 'city'}" type="button" @click.prevent="addRegion('city')">添加城市</button> &nbsp;
							<button class="ui button tiny basic" :class="{blue: regionType == 'provider'}" type="button" @click.prevent="addRegion('provider')">ISP</button> &nbsp;
						</div>
					</td>	
				</tr>
				<tr>
					<td>区域之间关系</td>
					<td>
						<select class="ui dropdown auto-width" v-model="regionConnector">
							<option value="OR">或</option>
							<option value="AND">且</option>
						</select>
						<p class="comment" v-if="regionConnector == 'OR'">匹配所选任一区域即认为匹配成功。</p>
						<p class="comment" v-if="regionConnector == 'AND'">匹配所有所选区域才认为匹配成功。</p>
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmRegions">确定</button> &nbsp;
				<a href="" @click.prevent="cancelRegions" title="取消"><i class="icon remove small"></i></a>
		</div>
		<div v-if="!isAdding && !isAddingBatch">
			<button class="ui button tiny" type="button" @click.prevent="addRegions">添加区域</button> &nbsp;
		</div>	
	</div>
</div>`
})

Vue.component("ns-record-health-check-config-box", {
	props:["value", "v-parent-config"],
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isOn: false,
				port: 0,
				timeoutSeconds: 0,
				countUp: 0,
				countDown: 0
			}
		}

		let parentConfig = this.vParentConfig

		return {
			config: config,
			portString: config.port.toString(),
			timeoutSecondsString: config.timeoutSeconds.toString(),
			countUpString: config.countUp.toString(),
			countDownString: config.countDown.toString(),

			portIsEditing: config.port > 0,
			timeoutSecondsIsEditing: config.timeoutSeconds > 0,
			countUpIsEditing: config.countUp > 0,
			countDownIsEditing: config.countDown > 0,

			parentConfig: parentConfig
		}
	},
	watch: {
		portString: function (value) {
			let port = parseInt(value.toString())
			if (isNaN(port) || port > 65535 || port < 1) {
				this.config.port = 0
			} else {
				this.config.port = port
			}
		},
		timeoutSecondsString: function (value) {
			let timeoutSeconds = parseInt(value.toString())
			if (isNaN(timeoutSeconds) || timeoutSeconds > 1000 || timeoutSeconds < 1) {
				this.config.timeoutSeconds = 0
			} else {
				this.config.timeoutSeconds = timeoutSeconds
			}
		},
		countUpString: function (value) {
			let countUp = parseInt(value.toString())
			if (isNaN(countUp) || countUp > 1000 || countUp < 1) {
				this.config.countUp = 0
			} else {
				this.config.countUp = countUp
			}
		},
		countDownString: function (value) {
			let countDown = parseInt(value.toString())
			if (isNaN(countDown) || countDown > 1000 || countDown < 1) {
				this.config.countDown = 0
			} else {
				this.config.countDown = countDown
			}
		}
	},
	template: `<div>
	<input type="hidden" name="recordHealthCheckJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">启用当前记录健康检查</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>检测端口</td>
				<td>
					<span v-if="!portIsEditing" class="grey">
						默认{{parentConfig.port}}
						&nbsp; <a href="" @click.prevent="portIsEditing = true; portString = parentConfig.port">[修改]</a>
					</span>
					<div v-show="portIsEditing">
						<div style="margin-bottom: 0.5em">
							<a href="" @click.prevent="portIsEditing = false; portString = '0'">[使用默认]</a>
						</div>
						<input type="text" v-model="portString" maxlength="5" style="width: 5em"/>
						<p class="comment">通过尝试连接A/AAAA记录中的IP加此端口来确定当前记录是否健康。</p>
					</div>
				</td>
			</tr>
			<tr>
				<td>超时时间</td>
				<td>
					<span v-if="!timeoutSecondsIsEditing" class="grey">
						默认{{parentConfig.timeoutSeconds}}秒
						&nbsp; <a href="" @click.prevent="timeoutSecondsIsEditing = true; timeoutSecondsString = parentConfig.timeoutSeconds">[修改]</a>
					</span>
					<div v-show="timeoutSecondsIsEditing">
						<div style="margin-bottom: 0.5em">
							<a href="" @click.prevent="timeoutSecondsIsEditing = false; timeoutSecondsString = '0'">[使用默认]</a>
						</div>
						<div class="ui input right labeled">
							<input type="text" style="width: 4em" v-model="timeoutSecondsString" maxlength="3"/>
							<span class="ui label">秒</span>
						</div>
					</div>
				</td>
			</tr>
			<tr>
				<td>默认连续上线次数</td>
				<td>
					<span v-if="!countUpIsEditing" class="grey">
						默认{{parentConfig.countUp}}次
						&nbsp; <a href="" @click.prevent="countUpIsEditing = true; countUpString = parentConfig.countUp">[修改]</a>
					</span>
					<div v-show="countUpIsEditing">
						<div style="margin-bottom: 0.5em">
							<a href="" @click.prevent="countUpIsEditing = false; countUpString = '0'">[使用默认]</a>
						</div>
						<div class="ui input right labeled">
							<input type="text" style="width: 4em" v-model="countUpString" maxlength="3"/>
							<span class="ui label">次</span>
						</div>
						<p class="comment">连续检测<span v-if="config.countUp > 0">{{config.countUp}}</span><span v-else>N</span>次成功后，认为当前记录是在线的。</p>
					</div>
				</td>
			</tr>
			<tr>
				<td>默认连续下线次数</td>
				<td>
					<span v-if="!countDownIsEditing" class="grey">
						默认{{parentConfig.countDown}}次
						&nbsp; <a href="" @click.prevent="countDownIsEditing = true; countDownString = parentConfig.countDown">[修改]</a>
					</span>
					<div v-show="countDownIsEditing">
						<div style="margin-bottom: 0.5em">
							<a href="" @click.prevent="countDownIsEditing = false; countDownString = '0'">[使用默认]</a>
						</div>
						<div class="ui input right labeled">
							<input type="text" style="width: 4em" v-model="countDownString" maxlength="3"/>
							<span class="ui label">次</span>
						</div>
						<p class="comment">连续检测<span v-if="config.countDown > 0">{{config.countDown}}</span><span v-else>N</span>次失败后，认为当前记录是离线的。</p>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("ns-create-records-table", {
	props: ["v-types"],
	data: function () {
		let types = this.vTypes
		if (types == null) {
			types = []
		}
		return {
			types: types,
			records: [
				{
					name: "",
					type: "A",
					value: "",
					routeCodes: [],
					ttl: 600,
					index: 0
				}
			],
			lastIndex: 0,
			isAddingRoutes: false // 是否正在添加线路
		}
	},
	methods: {
		add: function () {
			this.records.push({
				name: "",
				type: "A",
				value: "",
				routeCodes: [],
				ttl: 600,
				index: ++this.lastIndex
			})
			let that = this
			setTimeout(function () {
				that.$refs.nameInputs.$last().focus()
			}, 100)
		},
		remove: function (index) {
			this.records.$remove(index)
		},
		addRoutes: function () {
			this.isAddingRoutes = true
		},
		cancelRoutes: function () {
			let that = this
			setTimeout(function () {
				that.isAddingRoutes = false
			}, 1000)
		},
		changeRoutes: function (record, routes) {
			if (routes == null) {
				record.routeCodes = []
			} else {
				record.routeCodes = routes.map(function (route) {
					return route.code
				})
			}
		}
	},
	template: `<div>
<input type="hidden" name="recordsJSON" :value="JSON.stringify(records)"/>
<table class="ui table selectable celled" style="max-width: 60em">
	<thead class="full-width">
		<tr>
			<th style="width:10em">记录名</th>
			<th style="width:7em">记录类型</th>
			<th>线路</th>
			<th v-if="!isAddingRoutes">记录值</th>
			<th v-if="!isAddingRoutes">TTL</th>
			<th class="one op" v-if="!isAddingRoutes">操作</th>
		</tr>
	</thead>
	<tr v-for="(record, index) in records" :key="record.index">
		<td>
			<input type="text" style="width:10em" v-model="record.name" ref="nameInputs"/>		
		</td>
		<td>
			<select class="ui dropdown auto-width" v-model="record.type">
				<option v-for="type in types" :value="type.type">{{type.type}}</option>
			</select>
		</td>
		<td>
			<ns-routes-selector @add="addRoutes" @cancel="cancelRoutes" @change="changeRoutes(record, $event)"></ns-routes-selector>
		</td>
		<td v-if="!isAddingRoutes">
		  <input type="text" style="width:10em" maxlength="512" v-model="record.value"/>
		</td>
		<td v-if="!isAddingRoutes">
			<div class="ui input right labeled">
				<input type="text" v-model="record.ttl" style="width:5em" maxlength="8"/>
				<span class="ui label">秒</span>
			</div>
		</td>
		<td v-if="!isAddingRoutes">
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a>
		</td>
	</tr>
</table>
<button class="ui button tiny" type="button" @click.prevent="add">+</button>
</div>`,
})

// 选择单一线路
Vue.component("ns-route-selector", {
	props: ["v-route-code"],
	mounted: function () {
		let that = this
		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.routes = resp.data.routes
			})
	},
	data: function () {
		let routeCode = this.vRouteCode
		if (routeCode == null) {
			routeCode = ""
		}
		return {
			routeCode: routeCode,
			routes: []
		}
	},
	template: `<div>
	<div v-if="routes.length > 0">
		<select class="ui dropdown" name="routeCode" v-model="routeCode">
			<option value="">[线路]</option>
			<option v-for="route in routes" :value="route.code">{{route.name}}</option>
		</select>
	</div>
</div>`
})

Vue.component("ns-user-selector", {
	props: ["v-user-id"],
	data: function () {
		return {}
	},
	methods: {
		change: function (userId) {
			this.$emit("change", userId)
		}
	},
	template: `<div>
	<user-selector :v-user-id="vUserId" data-url="/ns/users/options" @change="change"></user-selector>
</div>`
})

Vue.component("ns-access-log-box", {
	props: ["v-access-log", "v-keyword"],
	data: function () {
		let accessLog = this.vAccessLog
		let isFailure = false

		if (accessLog.isRecursive) {
			if (accessLog.recordValue == null || accessLog.recordValue.length == 0) {
				isFailure = true
			}
		} else {
			if (accessLog.recordType == "SOA" || accessLog.recordType == "NS") {
				if (accessLog.recordValue == null || accessLog.recordValue.length == 0) {
					isFailure = true
				}
			}

			// 没有找到记录的不需要高亮显示，防止管理员看到红色就心理恐慌
		}

		return {
			accessLog: accessLog,
			isFailure: isFailure
		}
	},
	methods: {
		showLog: function () {
			let that = this
			let requestId = this.accessLog.requestId
			this.$parent.$children.forEach(function (v) {
				if (v.deselect != null) {
					v.deselect()
				}
			})
			this.select()

			teaweb.popup("/ns/clusters/accessLogs/viewPopup?requestId=" + requestId, {
				width: "50em",
				height: "24em",
				onClose: function () {
					that.deselect()
				}
			})
		},
		select: function () {
			this.$refs.box.parentNode.style.cssText = "background: rgba(0, 0, 0, 0.1)"
		},
		deselect: function () {
			this.$refs.box.parentNode.style.cssText = ""
		}
	},
	template: `<div class="access-log-row" :style="{'color': isFailure ? '#dc143c' : ''}" ref="box">
	<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey">[{{accessLog.region}}]</span> <keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword> [{{accessLog.timeLocal}}] [{{accessLog.networking}}] <em>{{accessLog.questionType}} <keyword :v-word="vKeyword">{{accessLog.questionName}}</keyword></em> -&gt; 
	
	<span v-if="accessLog.recordType != null && accessLog.recordType.length > 0"><em>{{accessLog.recordType}} <keyword :v-word="vKeyword">{{accessLog.recordValue}}</keyword></em></span>
	<span v-else class="disabled">&nbsp;[没有记录]</span>
	
	<!-- &nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>-->
	<div v-if="(accessLog.nsRoutes != null && accessLog.nsRoutes.length > 0) || accessLog.isRecursive" style="margin-top: 0.3em">
		<span class="ui label tiny basic grey" v-for="route in accessLog.nsRoutes">线路: {{route.name}}</span>
		<span class="ui label tiny basic grey" v-if="accessLog.isRecursive">递归DNS</span>
	</div>
	<div v-if="accessLog.error != null && accessLog.error.length > 0" style="color:#dc143c">
		<i class="icon warning circle"></i>错误：[{{accessLog.error}}]
	</div>
</div>`
})

Vue.component("ns-cluster-selector", {
	props: ["v-cluster-id"],
	mounted: function () {
		let that = this

		Tea.action("/ns/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})
	},
	data: function () {
		let clusterId = this.vClusterId
		if (clusterId == null) {
			clusterId = 0
		}
		return {
			clusters: [],
			clusterId: clusterId
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="clusterId" v-model="clusterId">
		<option value="0">[选择集群]</option>
		<option v-for="cluster in clusters" :value="cluster.id">{{cluster.name}}</option>
	</select>
</div>`
})

Vue.component("ns-cluster-combo-box", {
	props: ["v-cluster-id", "name"],
	data: function () {
		let that = this
		Tea.action("/ns/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})


		let inputName = "clusterId"
		if (this.name != null && this.name.length > 0) {
			inputName = this.name
		}

		return {
			clusters: [],
			inputName: inputName
		}
	},
	methods: {
		change: function (item) {
			if (item == null) {
				this.$emit("change", 0)
			} else {
				this.$emit("change", item.value)
			}
		}
	},
	template: `<div v-if="clusters.length > 0" style="min-width: 10.4em">
	<combo-box title="集群" placeholder="集群名称" :v-items="clusters" :name="inputName" :v-value="vClusterId" @change="change"></combo-box>
</div>`
})

Vue.component("plan-user-selector", {
	props: ["v-user-id"],
	data: function () {
		return {}
	},
	methods: {
		change: function (userId) {
			this.$emit("change", userId)
		}
	},
	template: `<div>
	<user-selector :v-user-id="vUserId" data-url="/plans/users/options" @change="change"></user-selector>
</div>`
})

// 显示流量限制说明
Vue.component("plan-limit-view", {
	props: ["value", "v-single-mode"],
	data: function () {
		let config = this.value

		let hasLimit = false
		if (!this.vSingleMode) {
			if (config.trafficLimit != null && config.trafficLimit.isOn && ((config.trafficLimit.dailySize != null && config.trafficLimit.dailySize.count > 0) || (config.trafficLimit.monthlySize != null && config.trafficLimit.monthlySize.count > 0))) {
				hasLimit = true
			} else if (config.dailyRequests > 0 || config.monthlyRequests > 0) {
				hasLimit = true
			}
		}

		return {
			config: config,
			hasLimit: hasLimit
		}
	},
	methods: {
		formatNumber: function (n) {
			return teaweb.formatNumber(n)
		},
		composeCapacity: function (capacity) {
			return teaweb.convertSizeCapacityToString(capacity)
		}
	},
	template: `<div style="font-size: 0.8em; color: grey">
	<div class="ui divider" v-if="hasLimit"></div>
	<div v-if="config.trafficLimit != null && config.trafficLimit.isOn">
		<span v-if="config.trafficLimit.dailySize != null && config.trafficLimit.dailySize.count > 0">日流量限制：{{composeCapacity(config.trafficLimit.dailySize)}}<br/></span>
		<span v-if="config.trafficLimit.monthlySize != null && config.trafficLimit.monthlySize.count > 0">月流量限制：{{composeCapacity(config.trafficLimit.monthlySize)}}<br/></span>
	</div>
	<div v-if="config.dailyRequests > 0">单日请求数限制：{{formatNumber(config.dailyRequests)}}</div>
	<div v-if="config.monthlyRequests > 0">单月请求数限制：{{formatNumber(config.monthlyRequests)}}</div>
	<div v-if="config.dailyWebsocketConnections > 0">单日Websocket限制：{{formatNumber(config.dailyWebsocketConnections)}}</div>
	<div v-if="config.monthlyWebsocketConnections > 0">单月Websocket限制：{{formatNumber(config.monthlyWebsocketConnections)}}</div>
	<div v-if="config.maxUploadSize != null && config.maxUploadSize.count > 0">文件上传限制：{{composeCapacity(config.maxUploadSize)}}</div>
</div>`
})

Vue.component("plan-price-view", {
	props: ["v-plan"],
	data: function () {
		return {
			plan: this.vPlan
		}
	},
	template: `<div>
	 <span v-if="plan.priceType == 'period'">
	 	按时间周期计费
	 	<div>
	 		<span class="grey small">
				<span v-if="plan.monthlyPrice > 0">月度：￥{{plan.monthlyPrice}}元<br/></span>
				<span v-if="plan.seasonallyPrice > 0">季度：￥{{plan.seasonallyPrice}}元<br/></span>
				<span v-if="plan.yearlyPrice > 0">年度：￥{{plan.yearlyPrice}}元</span>
			</span>
		</div>
	</span>
	<span v-if="plan.priceType == 'traffic'">
		按流量计费
		<div>
			<span class="grey small">基础价格：￥{{plan.trafficPrice.base}}元/GiB</span>
		</div>
	</span>
	<div v-if="plan.priceType == 'bandwidth' && plan.bandwidthPrice != null && plan.bandwidthPrice.percentile > 0">
		按{{plan.bandwidthPrice.percentile}}th带宽计费 
		<div>
			<div v-for="range in plan.bandwidthPrice.ranges">
				<span class="small grey">{{range.minMB}} - <span v-if="range.maxMB > 0">{{range.maxMB}}MiB</span><span v-else>&infin;</span>： <span v-if="range.totalPrice > 0">{{range.totalPrice}}元</span><span v-else="">{{range.pricePerMB}}元/MiB</span></span>
			</div>
		</div>
	</div>
</div>`
})

// 套餐价格配置
Vue.component("plan-price-config-box", {
	props: ["v-price-type", "v-monthly-price", "v-seasonally-price", "v-yearly-price", "v-traffic-price", "v-bandwidth-price", "v-disable-period"],
	data: function () {
		let priceType = this.vPriceType
		if (priceType == null) {
			priceType = "bandwidth"
		}

		// 按时间周期计费
		let monthlyPriceNumber = 0
		let monthlyPrice = this.vMonthlyPrice
		if (monthlyPrice == null || monthlyPrice <= 0) {
			monthlyPrice = ""
		} else {
			monthlyPrice = monthlyPrice.toString()
			monthlyPriceNumber = parseFloat(monthlyPrice)
			if (isNaN(monthlyPriceNumber)) {
				monthlyPriceNumber = 0
			}
		}

		let seasonallyPriceNumber = 0
		let seasonallyPrice = this.vSeasonallyPrice
		if (seasonallyPrice == null || seasonallyPrice <= 0) {
			seasonallyPrice = ""
		} else {
			seasonallyPrice = seasonallyPrice.toString()
			seasonallyPriceNumber = parseFloat(seasonallyPrice)
			if (isNaN(seasonallyPriceNumber)) {
				seasonallyPriceNumber = 0
			}
		}

		let yearlyPriceNumber = 0
		let yearlyPrice = this.vYearlyPrice
		if (yearlyPrice == null || yearlyPrice <= 0) {
			yearlyPrice = ""
		} else {
			yearlyPrice = yearlyPrice.toString()
			yearlyPriceNumber = parseFloat(yearlyPrice)
			if (isNaN(yearlyPriceNumber)) {
				yearlyPriceNumber = 0
			}
		}

		// 按流量计费
		let trafficPrice = this.vTrafficPrice
		let trafficPriceBaseNumber = 0
		if (trafficPrice != null) {
			trafficPriceBaseNumber = trafficPrice.base
		} else {
			trafficPrice = {
				base: 0
			}
		}
		let trafficPriceBase = ""
		if (trafficPriceBaseNumber > 0) {
			trafficPriceBase = trafficPriceBaseNumber.toString()
		}

		// 按带宽计费
		let bandwidthPrice = this.vBandwidthPrice
		if (bandwidthPrice == null) {
			bandwidthPrice = {
				percentile: 95,
				ranges: []
			}
		} else if (bandwidthPrice.ranges == null) {
			bandwidthPrice.ranges = []
		}

		return {
			priceType: priceType,
			monthlyPrice: monthlyPrice,
			seasonallyPrice: seasonallyPrice,
			yearlyPrice: yearlyPrice,

			monthlyPriceNumber: monthlyPriceNumber,
			seasonallyPriceNumber: seasonallyPriceNumber,
			yearlyPriceNumber: yearlyPriceNumber,

			trafficPriceBase: trafficPriceBase,
			trafficPrice: trafficPrice,

			bandwidthPrice: bandwidthPrice,
			bandwidthPercentile: bandwidthPrice.percentile
		}
	},
	methods: {
		changeBandwidthPriceRanges: function (ranges) {
			this.bandwidthPrice.ranges = ranges
		}
	},
	watch: {
		monthlyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.monthlyPriceNumber = price
		},
		seasonallyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.seasonallyPriceNumber = price
		},
		yearlyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.yearlyPriceNumber = price
		},
		trafficPriceBase: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.trafficPrice.base = price
		},
		bandwidthPercentile: function (v) {
			let percentile = parseInt(v)
			if (isNaN(percentile) || percentile <= 0) {
				percentile = 95
			} else if (percentile > 100) {
				percentile = 100
			}
			this.bandwidthPrice.percentile = percentile
		}
	},
	template: `<div>
	<input type="hidden" name="priceType" :value="priceType"/>
	<input type="hidden" name="monthlyPrice" :value="monthlyPriceNumber"/>
	<input type="hidden" name="seasonallyPrice" :value="seasonallyPriceNumber"/>
	<input type="hidden" name="yearlyPrice" :value="yearlyPriceNumber"/>
	<input type="hidden" name="trafficPriceJSON" :value="JSON.stringify(trafficPrice)"/>
	<input type="hidden" name="bandwidthPriceJSON" :value="JSON.stringify(bandwidthPrice)"/>
	
	<div>
		<radio :v-value="'bandwidth'" :value="priceType" v-model="priceType">&nbsp;按带宽</radio> &nbsp;
		<radio :v-value="'traffic'" :value="priceType" v-model="priceType">&nbsp;按流量</radio> &nbsp;
		<radio :v-value="'period'" :value="priceType" v-model="priceType" v-show="typeof(vDisablePeriod) != 'boolean' || !vDisablePeriod">&nbsp;按时间周期</radio>
	</div>
	
	<!-- 按时间周期 -->
	<div v-show="priceType == 'period'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">月度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="monthlyPrice"/>
						<span class="ui label">元</span>
					</div>
					<p class="comment">如果为0表示免费。</p>
				</td>
			</tr>
			<tr>
				<td class="title">季度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="seasonallyPrice"/>
						<span class="ui label">元</span>
					</div>
					<p class="comment">如果为0表示免费。</p>
				</td>
			</tr>
			<tr>
				<td class="title">年度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="yearlyPrice"/>
						<span class="ui label">元</span>
					</div>
					<p class="comment">如果为0表示免费。</p>
				</td>
			</tr>
		</table>
	</div>
	
	<!-- 按流量 -->
	<div v-show="priceType =='traffic'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">基础流量费用 *</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" v-model="trafficPriceBase" maxlength="10" style="width: 7em"/>
						<span class="ui label">元/GB</span>
					</div>
				</td>
			</tr>
		</table>
	</div>
	
	<!-- 按带宽 -->
	<div v-show="priceType == 'bandwidth'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">带宽百分位 *</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 4em" maxlength="3" v-model="bandwidthPercentile"/>
						<span class="ui label">th</span>
					</div>
				</td>
			</tr>
			<tr>
				<td>带宽价格</td>
				<td>
					<plan-bandwidth-ranges v-model="bandwidthPrice.ranges" @change="changeBandwidthPriceRanges"></plan-bandwidth-ranges>
				</td>
			</tr>
		</table>
	</div>
</div>`
})

Vue.component("plan-price-traffic-config-box", {
	props: ["v-plan-price-traffic-config"],
	data: function () {
		let config = this.vPlanPriceTrafficConfig
		if (config == null) {
			config = {
				base: 0,
				ranges: [],
				supportRegions: false
			}
		}

		if (config.ranges == null) {
			config.ranges = []
		}

		return {
			config: config,
			priceBase: config.base,
			isEditing: false
		}
	},
	watch: {
		priceBase: function (v) {
			let f = parseFloat(v)
			if (isNaN(f) || f < 0) {
				this.config.base = 0
			} else {
				this.config.base = f
			}
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	<input type="hidden" name="trafficPriceJSON" :value="JSON.stringify(config)"/>
	<div>
		基础流量价格：<span v-if="config.base > 0">{{config.base}}元/GB</span><span v-else class="disabled">没有设置</span> &nbsp; | &nbsp; 
		阶梯价格：<span v-if="config.ranges.length > 0">{{config.ranges.length}}段</span><span v-else class="disabled">没有设置</span>  &nbsp; <span v-if="config.supportRegions">| &nbsp;支持区域流量计费</span>
		<div style="margin-top: 0.5em">
			<a href="" @click.prevent="edit">修改 <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
		</div>		
	</div>
	<div v-show="isEditing" style="margin-top: 0.5em">
		<table class="ui table definition selectable" style="margin-top: 0">
			<tr>
				<td class="title">基础流量费用</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" v-model="priceBase" maxlength="10" style="width: 7em"/>
						<span class="ui label">元/GB</span>
					</div>
					<p class="comment">没有定义流量阶梯价格时，使用此价格。</p>
				</td>
			</tr>
			<tr>
				<td>流量阶梯价格</td>
				<td>
					<plan-traffic-ranges v-model="config.ranges"></plan-traffic-ranges>
				</td>
			</tr>
			<tr>
				<td>支持按区域流量计费</td>
				<td>
					<checkbox v-model="config.supportRegions"></checkbox>
					<p class="comment">选中后，表示可以根据节点所在区域设置不同的流量价格；并且开启此项后才可以使用流量包。</p>
				</td>
			</tr>
		</table>
	</div>
</div>`
})

Vue.component("plan-bandwidth-limit-view", {
	props: ["value"],
	template: `<div style="font-size: 0.8em; color: grey" v-if="value != null && value.bandwidthLimitPerNode != null && value.bandwidthLimitPerNode.count > 0">
	带宽限制：<bandwidth-size-capacity-view :v-value="value.bandwidthLimitPerNode"></bandwidth-size-capacity-view>
</div>`
})

Vue.component("plan-bandwidth-ranges", {
	props: ["value"],
	data: function () {
		let ranges = this.value
		if (ranges == null) {
			ranges = []
		}
		return {
			ranges: ranges,
			isAdding: false,

			minMB: "",
			minMBUnit: "mb",

			maxMB: "",
			maxMBUnit: "mb",

			pricePerMB: "",
			totalPrice: "",
			addingRange: {
				minMB: 0,
				maxMB: 0,
				pricePerMB: 0,
				totalPrice: 0
			}
		}
	},
	methods: {
		add: function () {
			this.isAdding = !this.isAdding
			let that = this
			setTimeout(function () {
				that.$refs.minMB.focus()
			})
		},
		cancelAdding: function () {
			this.isAdding = false
		},
		confirm: function () {
			if (this.addingRange.minMB < 0) {
				teaweb.warn("带宽下限需要大于0")
				return
			}
			if (this.addingRange.maxMB < 0) {
				teaweb.warn("带宽上限需要大于0")
				return
			}
			if (this.addingRange.pricePerMB <= 0) {
				teaweb.warn("请设置单位价格或者总价格")
				return
			}

			this.isAdding = false
			this.minMB = ""
			this.maxMB = ""
			this.pricePerMB = ""
			this.totalPrice = ""
			this.ranges.push(this.addingRange)
			this.ranges.$sort(function (v1, v2) {
				if (v1.minMB < v2.minMB) {
					return -1
				}
				if (v1.minMB == v2.minMB) {
					if (v2.maxMB == 0 || v1.maxMB < v2.maxMB) {
						return -1
					}
					return 0
				}
				return 1
			})
			this.change()
			this.addingRange = {
				minMB: 0,
				maxMB: 0,
				pricePerMB: 0,
				totalPrice: 0
			}
		},
		remove: function (index) {
			this.ranges.$remove(index)
			this.change()
		},
		change: function () {
			this.$emit("change", this.ranges)
		},
		formatMB: function (mb) {
			return teaweb.formatBits(mb * 1024 * 1024)
		},
		changeMinMB: function (v) {
			let minMB = parseFloat(v.toString())
			if (isNaN(minMB) || minMB < 0) {
				minMB = 0
			}
			switch (this.minMBUnit) {
				case "gb":
					minMB *= 1024
					break
				case "tb":
					minMB *= 1024 * 1024
					break
			}
			this.addingRange.minMB = minMB
		},
		changeMaxMB: function (v) {
			let maxMB = parseFloat(v.toString())
			if (isNaN(maxMB) || maxMB < 0) {
				maxMB = 0
			}
			switch (this.maxMBUnit) {
				case "gb":
					maxMB *= 1024
					break
				case "tb":
					maxMB *= 1024 * 1024
					break
			}
			this.addingRange.maxMB = maxMB
		}
	},
	watch: {
		minMB: function (v) {
			this.changeMinMB(v)
		},
		minMBUnit: function () {
			this.changeMinMB(this.minMB)
		},
		maxMB: function (v) {
			this.changeMaxMB(v)
		},
		maxMBUnit: function () {
			this.changeMaxMB(this.maxMB)
		},
		pricePerMB: function (v) {
			let pricePerMB = parseFloat(v.toString())
			if (isNaN(pricePerMB) || pricePerMB < 0) {
				pricePerMB = 0
			}
			this.addingRange.pricePerMB = pricePerMB
		},
		totalPrice: function (v) {
			let totalPrice = parseFloat(v.toString())
			if (isNaN(totalPrice) || totalPrice < 0) {
				totalPrice = 0
			}
			this.addingRange.totalPrice = totalPrice
		}
	},
	template: `<div>
	<!-- 已有价格 -->
	<div v-if="ranges.length > 0">
		<div class="ui label basic small" v-for="(range, index) in ranges" style="margin-bottom: 0.5em">
			{{formatMB(range.minMB)}} - <span v-if="range.maxMB > 0">{{formatMB(range.maxMB)}}</span><span v-else>&infin;</span> &nbsp;  价格：<span v-if="range.totalPrice > 0">{{range.totalPrice}}元</span><span v-else="">{{range.pricePerMB}}元/Mbps</span>
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- 添加 -->
	<div v-if="isAdding">
		<table class="ui table">
			<tr>
				<td class="title">带宽下限 *</td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" placeholder="最小带宽" style="width: 7em" maxlength="10" ref="minMB" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="minMB"/>
						</div>
						<div class="ui field">
							<select class="ui dropdown" v-model="minMBUnit">
								<option value="mb">Mbps</option>
								<option value="gb">Gbps</option>
								<option value="tb">Tbps</option>
							</select>
						</div>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">带宽上限 *</td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" placeholder="最大带宽" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="maxMB"/>
						</div>
						<div class="ui field">
							<select class="ui dropdown" v-model="maxMBUnit">
								<option value="mb">Mbps</option>
								<option value="gb">Gbps</option>
								<option value="tb">Tbps</option>
							</select>
						</div>
					</div>
					<p class="comment">如果填0，表示上不封顶。</p>
				</td>
			</tr>
			<tr>
				<td class="title">单位价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="单位价格" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="pricePerMB"/>
						<span class="ui label">元/Mbps</span>
					</div>
					<p class="comment">和总价格二选一。如果设置了单位价格，那么"总价格 = 单位价格 x 带宽/Mbps"。</p>
				</td>
			</tr>
			<tr>
				<td>总价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="总价格" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="totalPrice"/>
						<span class="ui label">元</span>
					</div>
					<p class="comment">固定的总价格，和单位价格二选一。</p>
				</td>
			</tr>
		</table>
		<button class="ui button small" type="button" @click.prevent="confirm">确定</button> &nbsp;
		<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove small"></i></a>
	</div>
	
	<!-- 按钮 -->
	<div v-if="!isAdding">
		<button class="ui button small" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("plan-price-bandwidth-config-box", {
	props: ["v-plan-price-bandwidth-config"],
	data: function () {
		let config = this.vPlanPriceBandwidthConfig
		if (config == null) {
			config = {
				percentile: 95,
				base: 0,
				ranges: [],
				supportRegions: false
			}
		}

		if (config.ranges == null) {
			config.ranges = []
		}

		return {
			config: config,
			bandwidthPercentile: config.percentile,
			priceBase: config.base,
			isEditing: false
		}
	},
	watch: {
		priceBase: function (v) {
			let f = parseFloat(v)
			if (isNaN(f) || f < 0) {
				this.config.base = 0
			} else {
				this.config.base = f
			}
		},
		bandwidthPercentile: function (v) {
			let i = parseInt(v)
			if (isNaN(i) || i < 0) {
				this.config.percentile = 0
			} else {
				this.config.percentile = i
			}
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
<input type="hidden" name="bandwidthPriceJSON" :value="JSON.stringify(config)"/>
<div>
	带宽百分位：<span v-if="config.percentile > 0">{{config.percentile}}th</span><span v-else class="disabled">没有设置</span> &nbsp; | &nbsp;
	基础带宽价格：<span v-if="config.base > 0">{{config.base}}元/Mbps</span><span v-else class="disabled">没有设置</span> &nbsp; | &nbsp; 
	阶梯价格：<span v-if="config.ranges.length > 0">{{config.ranges.length}}段</span><span v-else class="disabled">没有设置</span>  &nbsp; <span v-if="config.supportRegions">| &nbsp;支持区域带宽计费</span>
	<span v-if="config.bandwidthAlgo == 'avg'"> &nbsp;| &nbsp;使用平均带宽算法</span>
	<div style="margin-top: 0.5em">
		<a href="" @click.prevent="edit">修改 <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	</div>		
</div>
<div v-show="isEditing" style="margin-top: 0.5em">
	<table class="ui table definition selectable" style="margin-top: 0">
		<tr>
			<td class="title">带宽百分位 *</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" style="width: 4em" maxlength="3" v-model="bandwidthPercentile"/>
					<span class="ui label">th</span>
				</div>
				<p class="comment">带宽计费位置，在1-100之间。</p>
			</td>
		</tr>
		<tr>
			<td class="title">基础带宽费用</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="priceBase" maxlength="10" style="width: 7em"/>
					<span class="ui label">元/Mbps</span>
				</div>
				<p class="comment">没有定义带宽阶梯价格时，使用此价格。</p>
			</td>
		</tr>	
		<tr>
			<td>带宽阶梯价格</td>
			<td>
				<plan-bandwidth-ranges v-model="config.ranges"></plan-bandwidth-ranges>
			</td>
		</tr>
		<tr>
			<td>支持按区域带宽计费</td>
			<td>
				<checkbox v-model="config.supportRegions"></checkbox>
				<p class="comment">选中后，表示可以根据节点所在区域设置不同的带宽价格。</p>
			</td>
		</tr>
		<tr>
			<td>带宽算法</td>
			<td>
				<select class="ui dropdown auto-width" v-model="config.bandwidthAlgo">
                    <option value="secondly">峰值带宽</option>
                    <option value="avg">平均带宽</option>
                </select>
                <p class="comment" v-if="config.bandwidthAlgo == 'secondly'">按在计时时间段内（5分钟）最高带宽峰值计算，比如5分钟内最高的某个时间点带宽为100Mbps，那么就认为此时间段内的峰值带宽为100Mbps。修改此选项会同时影响到用量统计图表。</p>
                <p class="comment" v-if="config.bandwidthAlgo == 'avg'">按在计时时间段内（5分钟）平均带宽计算，即此时间段内的总流量除以时间段的秒数，比如5分钟（300秒）内总流量600MiB，那么带宽即为<code-label>600MiB * 8bit/300s = 16Mbps</code-label>；通常平均带宽算法要比峰值带宽要少很多。修改此选项会同时影响到用量统计图表。</p>
			</td>
		</tr>
	</table>
</div>	
</div>`
})

Vue.component("plan-traffic-ranges", {
	props: ["value"],
	data: function () {
		let ranges = this.value
		if (ranges == null) {
			ranges = []
		}
		return {
			ranges: ranges,
			isAdding: false,

			minGB: "",
			minGBUnit: "gb",

			maxGB: "",
			maxGBUnit: "gb",

			pricePerGB: "",
			totalPrice: "",
			addingRange: {
				minGB: 0,
				maxGB: 0,
				pricePerGB: 0,
				totalPrice: 0
			}
		}
	},
	methods: {
		add: function () {
			this.isAdding = !this.isAdding
			let that = this
			setTimeout(function () {
				that.$refs.minGB.focus()
			})
		},
		cancelAdding: function () {
			this.isAdding = false
		},
		confirm: function () {
			if (this.addingRange.minGB < 0) {
				teaweb.warn("流量下限需要大于0")
				return
			}
			if (this.addingRange.maxGB < 0) {
				teaweb.warn("流量上限需要大于0")
				return
			}
			if (this.addingRange.pricePerGB <= 0 && this.addingRange.totalPrice <= 0) {
				teaweb.warn("请设置单位价格或者总价格")
				return;
			}

			this.isAdding = false
			this.minGB = ""
			this.maxGB = ""
			this.pricePerGB = ""
			this.totalPrice = ""
			this.ranges.push(this.addingRange)
			this.ranges.$sort(function (v1, v2) {
				if (v1.minGB < v2.minGB) {
					return -1
				}
				if (v1.minGB == v2.minGB) {
					if (v2.maxGB == 0 || v1.maxGB < v2.maxGB) {
						return -1
					}
					return 0
				}
				return 1
			})
			this.change()
			this.addingRange = {
				minGB: 0,
				maxGB: 0,
				pricePerGB: 0,
				totalPrice: 0
			}
		},
		remove: function (index) {
			this.ranges.$remove(index)
			this.change()
		},
		change: function () {
			this.$emit("change", this.ranges)
		},
		formatGB: function (gb) {
			return teaweb.formatBytes(gb * 1024 * 1024 * 1024)
		},
		changeMinGB: function (v) {
			let minGB = parseFloat(v.toString())
			if (isNaN(minGB) || minGB < 0) {
				minGB = 0
			}
			switch (this.minGBUnit) {
				case "tb":
					minGB *= 1024
					break
				case "pb":
					minGB *= 1024 * 1024
					break
				case "eb":
					minGB *= 1024 * 1024 * 1024
					break
			}
			this.addingRange.minGB = minGB
		},
		changeMaxGB: function (v) {
			let maxGB = parseFloat(v.toString())
			if (isNaN(maxGB) || maxGB < 0) {
				maxGB = 0
			}
			switch (this.maxGBUnit) {
				case "tb":
					maxGB *= 1024
					break
				case "pb":
					maxGB *= 1024 * 1024
					break
				case "eb":
					maxGB *= 1024 * 1024 * 1024
					break
			}
			this.addingRange.maxGB = maxGB
		}
	},
	watch: {
		minGB: function (v) {
			this.changeMinGB(v)
		},
		minGBUnit: function (v) {
			this.changeMinGB(this.minGB)
		},
		maxGB: function (v) {
			this.changeMaxGB(v)
		},
		maxGBUnit: function (v) {
			this.changeMaxGB(this.maxGB)
		},
		pricePerGB: function (v) {
			let pricePerGB = parseFloat(v.toString())
			if (isNaN(pricePerGB) || pricePerGB < 0) {
				pricePerGB = 0
			}
			this.addingRange.pricePerGB = pricePerGB
		},
		totalPrice: function (v) {
			let totalPrice = parseFloat(v.toString())
			if (isNaN(totalPrice) || totalPrice < 0) {
				totalPrice = 0
			}
			this.addingRange.totalPrice = totalPrice
		}
	},
	template: `<div>
	<!-- 已有价格 -->
	<div v-if="ranges.length > 0">
		<div class="ui label basic small" v-for="(range, index) in ranges" style="margin-bottom: 0.5em">
			{{formatGB(range.minGB)}} - <span v-if="range.maxGB > 0">{{formatGB(range.maxGB)}}</span><span v-else>&infin;</span> &nbsp;  价格：<span v-if="range.totalPrice > 0">{{range.totalPrice}}元</span><span v-else="">{{range.pricePerGB}}元/GB</span>
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- 添加 -->
	<div v-if="isAdding">
		<table class="ui table">
			<tr>
				<td class="title">流量下限 *</td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" placeholder="最小流量" style="width: 7em" maxlength="10" ref="minGB" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="minGB"/>
						</div>
						<div class="ui field">
							<select class="ui dropdown" v-model="minGBUnit">
								<option value="gb">GB</option>
								<option value="tb">TB</option>
								<option value="pb">PB</option>
								<option value="eb">EB</option>
							</select>
						</div>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">流量上限 *</td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" placeholder="最大流量" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="maxGB"/>
						</div>
						<div class="ui field">
							<select class="ui dropdown" v-model="maxGBUnit">
								<option value="gb">GB</option>
								<option value="tb">TB</option>
								<option value="pb">PB</option>
								<option value="eb">EB</option>
							</select>
						</div>
					</div>
					<p class="comment">如果填0，表示上不封顶。</p>
				</td>
			</tr>
			<tr>
				<td class="title">单位价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="单位价格" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="pricePerGB"/>
						<span class="ui label">元/GB</span>
					</div>
					<p class="comment">和总价格二选一。如果设置了单位价格，那么"总价格 = 单位价格 x 流量/GB"。</p>
				</td>
			</tr>
			<tr>
				<td>总价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="总价格" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="totalPrice"/>
						<span class="ui label">元</span>
					</div>
					<p class="comment">固定的总价格，和单位价格二选一。</p>
				</td>
			</tr>
		</table>
		<button class="ui button small" type="button" @click.prevent="confirm">确定</button> &nbsp;
		<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove small"></i></a>
	</div>
	
	<!-- 按钮 -->
	<div v-if="!isAdding">
		<button class="ui button small" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("http-stat-config-box", {
	props: ["v-stat-config", "v-is-location", "v-is-group"],
	data: function () {
		let stat = this.vStatConfig
		if (stat == null) {
			stat = {
				isPrior: false,
				isOn: false
			}
		}
		return {
			stat: stat
		}
	},
	template: `<div>
	<input type="hidden" name="statJSON" :value="JSON.stringify(stat)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="stat" v-if="vIsLocation || vIsGroup" ></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || stat.isPrior">
			<tr>
				<td class="title">启用统计</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="stat.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
<div class="margin"></div>
</div>`
})

Vue.component("http-firewall-page-options-viewer", {
	props: ["v-page-options"],
	data: function () {
		return {
			options: this.vPageOptions
		}
	},
	template: `<div>
	<span v-if="options == null">默认设置</span>
	<div v-else>
		状态码：{{options.status}} / 提示内容：<span v-if="options.body != null && options.body.length > 0">[{{options.body.length}}字符]</span>
	</div>
</div>	
`
})

Vue.component("http-request-conds-box", {
	props: ["v-conds"],
	data: function () {
		let conds = this.vConds
		if (conds == null) {
			conds = {
				isOn: true,
				connector: "or",
				groups: []
			}
		}
		if (conds.groups == null) {
			conds.groups = []
		}
		return {
			conds: conds,
			components: window.REQUEST_COND_COMPONENTS
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.conds)
		},
		addGroup: function () {
			window.UPDATING_COND_GROUP = null

			let that = this
			teaweb.popup("/servers/server/settings/conds/addGroupPopup", {
				height: "30em",
				callback: function (resp) {
					that.conds.groups.push(resp.data.group)
					that.change()
				}
			})
		},
		updateGroup: function (groupIndex, group) {
			window.UPDATING_COND_GROUP = group
			let that = this
			teaweb.popup("/servers/server/settings/conds/addGroupPopup", {
				height: "30em",
				callback: function (resp) {
					Vue.set(that.conds.groups, groupIndex, resp.data.group)
					that.change()
				}
			})
		},
		removeGroup: function (groupIndex) {
			let that = this
			teaweb.confirm("确定要删除这一组条件吗？", function () {
				that.conds.groups.$remove(groupIndex)
				that.change()
			})
		},
		typeName: function (cond) {
			let c = this.components.$find(function (k, v) {
				return v.type == cond.type
			})
			if (c != null) {
				return c.name;
			}
			return cond.param + " " + cond.operator
		}
	},
	template: `<div>
		<input type="hidden" name="condsJSON" :value="JSON.stringify(conds)"/>
		<div v-if="conds.groups.length > 0">
			<table class="ui table">
				<tr v-for="(group, groupIndex) in conds.groups">
					<td class="title" :class="{'color-border':conds.connector == 'and'}" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">分组{{groupIndex+1}}</td>
					<td style="background: white; word-break: break-all" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">
						<var v-for="(cond, index) in group.conds" style="font-style: normal;display: inline-block; margin-bottom:0.5em">
							<span class="ui label tiny">
								<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
								<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{typeName(cond)}}: </var>
								{{cond.value}}
								<sup v-if="cond.isCaseInsensitive" title="不区分大小写"><i class="icon info small"></i></sup>
							</span>
							
							<var v-if="index < group.conds.length - 1"> {{group.connector}} &nbsp;</var>
						</var>
					</td>
					<td style="width: 5em; background: white" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">
						<a href="" title="修改分组" @click.prevent="updateGroup(groupIndex, group)"><i class="icon pencil small"></i></a> <a href="" title="删除分组" @click.prevent="removeGroup(groupIndex)"><i class="icon remove"></i></a>
					</td>
				</tr>
			</table>
			<div class="ui divider"></div>
		</div>
		
		<!-- 分组之间关系 -->
		<table class="ui table" v-if="conds.groups.length > 1">
			<tr>
				<td class="title">分组之间关系</td>
				<td>
					<select class="ui dropdown auto-width" v-model="conds.connector">
						<option value="and">和</option>
						<option value="or">或</option>
					</select>
					<p class="comment">
						<span v-if="conds.connector == 'or'">只要满足其中一个条件分组即可。</span>
						<span v-if="conds.connector == 'and'">需要满足所有条件分组。</span>
					</p>	
				</td>
			</tr>
		</table>
		
		<div>
			<button class="ui button tiny basic" type="button" @click.prevent="addGroup()">+添加分组</button>
		</div>
	</div>	
</div>`
})

Vue.component("ssl-config-box", {
	props: [
		"v-ssl-policy",
		"v-protocol",
		"v-server-id",
		"v-support-http3"
	],
	created: function () {
		let that = this
		setTimeout(function () {
			that.sortableCipherSuites()
		}, 100)
	},
	data: function () {
		let policy = this.vSslPolicy
		if (policy == null) {
			policy = {
				id: 0,
				isOn: true,
				certRefs: [],
				certs: [],
				clientCARefs: [],
				clientCACerts: [],
				clientAuthType: 0,
				minVersion: "TLS 1.1",
				hsts: null,
				cipherSuitesIsOn: false,
				cipherSuites: [],
				http2Enabled: true,
				http3Enabled: false,
				ocspIsOn: false
			}
		} else {
			if (policy.certRefs == null) {
				policy.certRefs = []
			}
			if (policy.certs == null) {
				policy.certs = []
			}
			if (policy.clientCARefs == null) {
				policy.clientCARefs = []
			}
			if (policy.clientCACerts == null) {
				policy.clientCACerts = []
			}
			if (policy.cipherSuites == null) {
				policy.cipherSuites = []
			}
		}

		let hsts = policy.hsts
		let hstsMaxAgeString = "31536000"
		if (hsts == null) {
			hsts = {
				isOn: false,
				maxAge: 31536000,
				includeSubDomains: false,
				preload: false,
				domains: []
			}
		}
		if (hsts.maxAge != null) {
			hstsMaxAgeString = hsts.maxAge.toString()
		}

		return {
			policy: policy,

			// hsts
			hsts: hsts,
			hstsOptionsVisible: false,
			hstsDomainAdding: false,
			hstsMaxAgeString: hstsMaxAgeString,
			addingHstsDomain: "",
			hstsDomainEditingIndex: -1,

			// 相关数据
			allVersions: window.SSL_ALL_VERSIONS,
			allCipherSuites: window.SSL_ALL_CIPHER_SUITES.$copy(),
			modernCipherSuites: window.SSL_MODERN_CIPHER_SUITES,
			intermediateCipherSuites: window.SSL_INTERMEDIATE_CIPHER_SUITES,
			allClientAuthTypes: window.SSL_ALL_CLIENT_AUTH_TYPES,
			cipherSuitesVisible: false,

			// 高级选项
			moreOptionsVisible: false
		}
	},
	watch: {
		hsts: {
			deep: true,
			handler: function () {
				this.policy.hsts = this.hsts
			}
		}
	},
	methods: {
		// 删除证书
		removeCert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前网站不再使用此证书。", function () {
				that.policy.certRefs.$remove(index)
				that.policy.certs.$remove(index)
			})
		},

		// 选择证书
		selectCert: function () {
			let that = this
			let selectedCertIds = []
			if (this.policy != null && this.policy.certs.length > 0) {
				this.policy.certs.forEach(function (cert) {
					selectedCertIds.push(cert.id.toString())
				})
			}
			let serverId = this.vServerId
			if (serverId == null) {
				serverId = 0
			}
			teaweb.popup("/servers/certs/selectPopup?selectedCertIds=" + selectedCertIds + "&serverId=" + serverId, {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					if (resp.data.cert != null && resp.data.certRef != null) {
						that.policy.certRefs.push(resp.data.certRef)
						that.policy.certs.push(resp.data.cert)
					}
					if (resp.data.certs != null && resp.data.certRefs != null) {
						that.policy.certRefs.$pushAll(resp.data.certRefs)
						that.policy.certs.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 上传证书
		uploadCert: function () {
			let that = this
			let serverId = this.vServerId
			if (typeof serverId != "number" && typeof serverId != "string") {
				serverId = 0
			}
			teaweb.popup("/servers/certs/uploadPopup?serverId=" + serverId, {
				height: "30em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						that.policy.certRefs.push(resp.data.certRef)
						that.policy.certs.push(resp.data.cert)
					})
				}
			})
		},

		// 批量上传
		uploadBatch: function () {
			let that = this
			let serverId = this.vServerId
			if (typeof serverId != "number" && typeof serverId != "string") {
				serverId = 0
			}
			teaweb.popup("/servers/certs/uploadBatchPopup?serverId=" + serverId, {
				callback: function (resp) {
					if (resp.data.cert != null) {
						that.policy.certRefs.push(resp.data.certRef)
						that.policy.certs.push(resp.data.cert)
					}
					if (resp.data.certs != null) {
						that.policy.certRefs.$pushAll(resp.data.certRefs)
						that.policy.certs.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 申请证书
		requestCert: function () {
			// 已经在证书中的域名
			let excludeServerNames = []
			if (this.policy != null && this.policy.certs.length > 0) {
				this.policy.certs.forEach(function (cert) {
					excludeServerNames.$pushAll(cert.dnsNames)
				})
			}

			let that = this
			teaweb.popup("/servers/server/settings/https/requestCertPopup?serverId=" + this.vServerId + "&excludeServerNames=" + excludeServerNames.join(","), {
				callback: function () {
					that.policy.certRefs.push(resp.data.certRef)
					that.policy.certs.push(resp.data.cert)
				}
			})
		},

		// 更多选项
		changeOptionsVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},

		// 格式化时间
		formatTime: function (timestamp) {
			return new Date(timestamp * 1000).format("Y-m-d")
		},

		// 格式化加密套件
		formatCipherSuite: function (cipherSuite) {
			return cipherSuite.replace(/(AES|3DES)/, "<var style=\"font-weight: bold\">$1</var>")
		},

		// 添加单个套件
		addCipherSuite: function (cipherSuite) {
			if (!this.policy.cipherSuites.$contains(cipherSuite)) {
				this.policy.cipherSuites.push(cipherSuite)
			}
			this.allCipherSuites.$removeValue(cipherSuite)
		},

		// 删除单个套件
		removeCipherSuite: function (cipherSuite) {
			let that = this
			teaweb.confirm("确定要删除此套件吗？", function () {
				that.policy.cipherSuites.$removeValue(cipherSuite)
				that.allCipherSuites = window.SSL_ALL_CIPHER_SUITES.$findAll(function (k, v) {
					return !that.policy.cipherSuites.$contains(v)
				})
			})
		},

		// 清除所选套件
		clearCipherSuites: function () {
			let that = this
			teaweb.confirm("确定要清除所有已选套件吗？", function () {
				that.policy.cipherSuites = []
				that.allCipherSuites = window.SSL_ALL_CIPHER_SUITES.$copy()
			})
		},

		// 批量添加套件
		addBatchCipherSuites: function (suites) {
			var that = this
			teaweb.confirm("确定要批量添加套件？", function () {
				suites.$each(function (k, v) {
					if (that.policy.cipherSuites.$contains(v)) {
						return
					}
					that.policy.cipherSuites.push(v)
				})
			})
		},

		/**
		 * 套件拖动排序
		 */
		sortableCipherSuites: function () {
			var box = document.querySelector(".cipher-suites-box")
			Sortable.create(box, {
				draggable: ".label",
				handle: ".icon.handle",
				onStart: function () {

				},
				onUpdate: function (event) {

				}
			})
		},

		// 显示所有套件
		showAllCipherSuites: function () {
			this.cipherSuitesVisible = !this.cipherSuitesVisible
		},

		// 显示HSTS更多选项
		showMoreHSTS: function () {
			this.hstsOptionsVisible = !this.hstsOptionsVisible;
			if (this.hstsOptionsVisible) {
				this.changeHSTSMaxAge()
			}
		},

		// 监控HSTS有效期修改
		changeHSTSMaxAge: function () {
			var v = parseInt(this.hstsMaxAgeString)
			if (isNaN(v) || v < 0) {
				this.hsts.maxAge = 0
				this.hsts.days = "-"
				return
			}
			this.hsts.maxAge = v
			this.hsts.days = v / 86400
			if (this.hsts.days == 0) {
				this.hsts.days = "-"
			}
		},

		// 设置HSTS有效期
		setHSTSMaxAge: function (maxAge) {
			this.hstsMaxAgeString = maxAge.toString()
			this.changeHSTSMaxAge()
		},

		// 添加HSTS域名
		addHstsDomain: function () {
			this.hstsDomainAdding = true
			this.hstsDomainEditingIndex = -1
			let that = this
			setTimeout(function () {
				that.$refs.addingHstsDomain.focus()
			}, 100)
		},

		// 修改HSTS域名
		editHstsDomain: function (index) {
			this.hstsDomainEditingIndex = index
			this.addingHstsDomain = this.hsts.domains[index]
			this.hstsDomainAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingHstsDomain.focus()
			}, 100)
		},

		// 确认HSTS域名添加
		confirmAddHstsDomain: function () {
			this.addingHstsDomain = this.addingHstsDomain.trim()
			if (this.addingHstsDomain.length == 0) {
				return;
			}
			if (this.hstsDomainEditingIndex > -1) {
				this.hsts.domains[this.hstsDomainEditingIndex] = this.addingHstsDomain
			} else {
				this.hsts.domains.push(this.addingHstsDomain)
			}
			this.cancelHstsDomainAdding()
		},

		// 取消HSTS域名添加
		cancelHstsDomainAdding: function () {
			this.hstsDomainAdding = false
			this.addingHstsDomain = ""
			this.hstsDomainEditingIndex = -1
		},

		// 删除HSTS域名
		removeHstsDomain: function (index) {
			this.cancelHstsDomainAdding()
			this.hsts.domains.$remove(index)
		},

		// 选择客户端CA证书
		selectClientCACert: function () {
			let that = this
			teaweb.popup("/servers/certs/selectPopup?isCA=1", {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					if (resp.data.cert != null && resp.data.certRef != null) {
						that.policy.clientCARefs.push(resp.data.certRef)
						that.policy.clientCACerts.push(resp.data.cert)
					}
					if (resp.data.certs != null && resp.data.certRefs != null) {
						that.policy.clientCARefs.$pushAll(resp.data.certRefs)
						that.policy.clientCACerts.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 上传CA证书
		uploadClientCACert: function () {
			let that = this
			teaweb.popup("/servers/certs/uploadPopup?isCA=1", {
				height: "28em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						that.policy.clientCARefs.push(resp.data.certRef)
						that.policy.clientCACerts.push(resp.data.cert)
					})
				}
			})
		},

		// 删除客户端CA证书
		removeClientCACert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前网站不再使用此证书。", function () {
				that.policy.clientCARefs.$remove(index)
				that.policy.clientCACerts.$remove(index)
			})
		}
	},
	template: `<div>
	<h4>SSL/TLS相关配置</h4>
	<input type="hidden" name="sslPolicyJSON" :value="JSON.stringify(policy)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr v-show="vProtocol == 'https'">
				<td class="title">启用HTTP/2</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.http2Enabled"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="vProtocol == 'https' && vSupportHttp3">
				<td class="title">启用HTTP/3</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.http3Enabled"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">设置证书</td>
				<td>
					<div v-if="policy.certs != null && policy.certs.length > 0">
						<div class="ui label small basic" v-for="(cert, index) in policy.certs" style="margin-top: 0.2em">
							{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<div v-else>
						<span class="red">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
					<span class="disabled">|</span> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadBatch()">批量上传证书</button> &nbsp;
					<span class="disabled">|</span> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="requestCert()" v-if="vServerId != null && vServerId > 0">申请免费证书</button>
				</td>
			</tr>
			<tr>
				<td>TLS最低版本</td>
				<td>
					<select v-model="policy.minVersion" class="ui dropdown auto-width">
						<option v-for="version in allVersions" :value="version">{{version}}</option>
					</select>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeOptionsVisible"></more-options-tbody>
		<tbody v-show="moreOptionsVisible">
			<!-- 加密套件 -->
			<tr>
				<td>加密算法套件<em>（CipherSuites）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.cipherSuitesIsOn" />
						<label>是否要自定义</label>
					</div>
					<div v-show="policy.cipherSuitesIsOn">
						<div class="ui divider"></div>
						<div class="cipher-suites-box">
							已添加套件({{policy.cipherSuites.length}})：
							<div v-for="cipherSuite in policy.cipherSuites" class="ui label tiny basic" style="margin-bottom: 0.5em">
								<input type="hidden" name="cipherSuites" :value="cipherSuite"/>
								<span v-html="formatCipherSuite(cipherSuite)"></span> &nbsp; <a href="" title="删除套件" @click.prevent="removeCipherSuite(cipherSuite)"><i class="icon remove"></i></a>
								<a href="" title="拖动改变顺序"><i class="icon bars handle"></i></a>
							</div>
						</div>
						<div>
							<div class="ui divider"></div>
							<span v-if="policy.cipherSuites.length > 0"><a href="" @click.prevent="clearCipherSuites()">[清除所有已选套件]</a> &nbsp; </span>
							<a href="" @click.prevent="addBatchCipherSuites(modernCipherSuites)">[添加推荐套件]</a> &nbsp;
							<a href="" @click.prevent="addBatchCipherSuites(intermediateCipherSuites)">[添加兼容套件]</a>
							<div class="ui divider"></div>
						</div>
		
						<div class="cipher-all-suites-box">
							<a href="" @click.prevent="showAllCipherSuites()"><span v-if="policy.cipherSuites.length == 0">所有</span>可选套件({{allCipherSuites.length}}) <i class="icon angle" :class="{down:!cipherSuitesVisible, up:cipherSuitesVisible}"></i></a>
							<a href="" v-if="cipherSuitesVisible" v-for="cipherSuite in allCipherSuites" class="ui label tiny" title="点击添加到自定义套件中" @click.prevent="addCipherSuite(cipherSuite)" v-html="formatCipherSuite(cipherSuite)" style="margin-bottom:0.5em"></a>
						</div>
						<p class="comment" v-if="cipherSuitesVisible">点击可选套件添加。</p>
					</div>
				</td>
			</tr>
			
			<!-- HSTS -->
			<tr v-show="vProtocol == 'https'">
				<td :class="{'color-border':hsts.isOn}">开启HSTS</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsOn" v-model="hsts.isOn" value="1"/>
						<label></label>
					</div>
					<p class="comment">
						 开启后，会自动在响应Header中加入
						 <span class="ui label small">Strict-Transport-Security:
							 <var v-if="!hsts.isOn">...</var>
							 <var v-if="hsts.isOn"><span>max-age=</span>{{hsts.maxAge}}</var>
							 <var v-if="hsts.isOn && hsts.includeSubDomains">; includeSubDomains</var>
							 <var v-if="hsts.isOn && hsts.preload">; preload</var>
						 </span>
						  <span v-if="hsts.isOn">
							<a href="" @click.prevent="showMoreHSTS()">修改<i class="icon angle" :class="{down:!hstsOptionsVisible, up:hstsOptionsVisible}"></i> </a>
						 </span>
					</p>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS有效时间<em>（max-age）</em></td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" name="hstsMaxAge" v-model="hstsMaxAgeString" maxlength="10" size="10" @input="changeHSTSMaxAge()"/>
						</div>
						<div class="ui field">
							秒
						</div>
						<div class="ui field">{{hsts.days}}天</div>
					</div>
					<p class="comment">
						<a href="" @click.prevent="setHSTSMaxAge(31536000)" :class="{active:hsts.maxAge == 31536000}">[1年/365天]</a> &nbsp; &nbsp;
						<a href="" @click.prevent="setHSTSMaxAge(15768000)" :class="{active:hsts.maxAge == 15768000}">[6个月/182.5天]</a> &nbsp;  &nbsp;
						<a href="" @click.prevent="setHSTSMaxAge(2592000)"  :class="{active:hsts.maxAge == 2592000}">[1个月/30天]</a>
					</p>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS包含子域名<em>（includeSubDomains）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsIncludeSubDomains" value="1" v-model="hsts.includeSubDomains"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS预加载<em>（preload）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsPreload" value="1" v-model="hsts.preload"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS生效的域名</td>
				<td colspan="2">
					<div class="names-box">
					<span class="ui label tiny basic" v-for="(domain, arrayIndex) in hsts.domains" :class="{blue:hstsDomainEditingIndex == arrayIndex}">{{domain}}
						<input type="hidden" name="hstsDomains" :value="domain"/> &nbsp;
						<a href="" @click.prevent="editHstsDomain(arrayIndex)" title="修改"><i class="icon pencil"></i></a>
						<a href="" @click.prevent="removeHstsDomain(arrayIndex)" title="删除"><i class="icon remove"></i></a>
					</span>
					</div>
					<div class="ui fields inline" v-if="hstsDomainAdding" style="margin-top:0.8em">
						<div class="ui field">
							<input type="text" name="addingHstsDomain" ref="addingHstsDomain" style="width:16em" maxlength="100" placeholder="域名，比如example.com" @keyup.enter="confirmAddHstsDomain()" @keypress.enter.prevent="1" v-model="addingHstsDomain" />
						</div>
						<div class="ui field">
							<button class="ui button tiny" type="button" @click="confirmAddHstsDomain()">确定</button>
							&nbsp; <a href="" @click.prevent="cancelHstsDomainAdding()">取消</a>
						</div>
					</div>
					<div class="ui field" style="margin-top: 1em">
						<button class="ui button tiny" type="button" @click="addHstsDomain()">+</button>
					</div>
					<p class="comment">如果没有设置域名的话，则默认支持所有的域名。</p>
				</td>
			</tr>
			
			<!-- OCSP -->
			<tr>
				<td>OCSP Stapling</td>
				<td><checkbox name="ocspIsOn" v-model="policy.ocspIsOn"></checkbox>
					<p class="comment">选中表示启用OCSP Stapling。</p>
				</td>
			</tr>
			
			<!-- 客户端认证 -->
			<tr>
				<td>客户端认证方式</td>
				<td>
					<select name="clientAuthType" v-model="policy.clientAuthType" class="ui dropdown auto-width">
						<option v-for="authType in allClientAuthTypes" :value="authType.type">{{authType.name}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>客户端认证CA证书</td>
				<td>
					<div v-if="policy.clientCACerts != null && policy.clientCACerts.length > 0">
						<div class="ui label small basic" v-for="(cert, index) in policy.clientCACerts">
							{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeClientCACert()"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="selectClientCACert()">选择已有证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadClientCACert()">上传新证书</button>
					<p class="comment">用来校验客户端证书以增强安全性，通常不需要设置。</p>
				</td>
			</tr>
		</tbody>	
	</table>
	<div class="ui margin"></div>
</div>`
})

// Action列表
Vue.component("http-firewall-actions-view", {
	props: ["v-actions"],
	template: `<div>
		<div v-for="action in vActions" style="margin-bottom: 0.3em">
			<span :class="{red: action.category == 'block', orange: action.category == 'verify', green: action.category == 'allow'}">{{action.name}} ({{action.code.toUpperCase()}})
			  	<div v-if="action.options != null">
			  		<span class="grey small" v-if="action.code.toLowerCase() == 'page'">[{{action.options.status}}]</span>
			  		<span class="grey small" v-if="action.code.toLowerCase() == 'allow' && action.options != null && action.options.scope != null && action.options.scope.length > 0">
			  			<span v-if="action.options.scope == 'group'">[分组]</span>
						<span v-if="action.options.scope == 'server'">[网站]</span>
						<span v-if="action.options.scope == 'global'">[网站和策略]</span>	
					</span>
					<span class="grey small" v-if="action.code.toLowerCase() == 'record_ip'">
						<span v-if="action.options.type == 'black'" class="red">黑名单</span>
						<span v-if="action.options.type == 'white'" class="green">白名单</span>
						<span v-if="action.options.type == 'grey'" class="grey">灰名单</span>
					</span>
				</div>	
			</span>
		</div>             
</div>`
})

// 显示WAF规则的标签
Vue.component("http-firewall-rule-label", {
	props: ["v-rule"],
	data: function () {
		return {
			rule: this.vRule
		}
	},
	methods: {
		showErr: function (err) {
			teaweb.popupTip("规则校验错误，请修正：<span class=\"red\">"  + teaweb.encodeHTML(err) + "</span>")
		},
		calculateParamName: function (param) {
			let paramName = ""
			if (param != null) {
				window.WAF_RULE_CHECKPOINTS.forEach(function (checkpoint) {
					if (param == "${" + checkpoint.prefix + "}" || param.startsWith("${" + checkpoint.prefix + ".")) {
						paramName = checkpoint.name
					}
				})
			}
			return paramName
		},
		calculateParamDescription: function (param) {
			let paramName = ""
			let paramDescription = ""
			if (param != null) {
				window.WAF_RULE_CHECKPOINTS.forEach(function (checkpoint) {
					if (param == "${" + checkpoint.prefix + "}" || param.startsWith("${" + checkpoint.prefix + ".")) {
						paramName = checkpoint.name
						paramDescription = checkpoint.description
					}
				})
			}
			return paramName + ": " + paramDescription
		},
		operatorName: function (operatorCode) {
			let operatorName = operatorCode
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorName = v.name
					}
				})
			}

			return operatorName
		},
		operatorDescription: function (operatorCode) {
			let operatorName = operatorCode
			let operatorDescription = ""
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorName = v.name
						operatorDescription = v.description
					}
				})
			}

			return operatorName + ": " + operatorDescription
		},
		operatorDataType: function (operatorCode) {
			let operatorDataType = "none"
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorDataType = v.dataType
					}
				})
			}

			return operatorDataType
		},
		isEmptyString: function (v) {
			return typeof v == "string" && v.length == 0
		}
	},
	template: `<div>
	<div class="ui label small basic" style="line-height: 1.5">
		{{rule.name}} <span :title="calculateParamDescription(rule.param)" class="hover">{{calculateParamName(rule.param)}}<span class="small grey"> {{rule.param}}</span></span>

		<!-- cc2 -->
		<span v-if="rule.param == '\${cc2}'">
			{{rule.checkpointOptions.period}}秒内请求数
		</span>

		<!-- refererBlock -->
		<span v-if="rule.param == '\${refererBlock}'">
			<span v-if="rule.checkpointOptions.allowDomains != null && rule.checkpointOptions.allowDomains.length > 0">允许{{rule.checkpointOptions.allowDomains}}</span>
			<span v-if="rule.checkpointOptions.denyDomains != null && rule.checkpointOptions.denyDomains.length > 0">禁止{{rule.checkpointOptions.denyDomains}}</span>
		</span>

		<span v-else>
			<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> 
		<span class="hover" :class="{dash:!rule.isComposed && rule.isCaseInsensitive}" :title="operatorDescription(rule.operator) + ((!rule.isComposed && rule.isCaseInsensitive) ? '\\n[大小写不敏感] ':'')">&lt;{{operatorName(rule.operator)}}&gt;</span> 
			<span class="hover" v-if="!isEmptyString(rule.value)">{{rule.value}}</span>
			<span v-else-if="operatorDataType(rule.operator) != 'none'" class="disabled" style="font-weight: normal" title="空字符串">[空]</span>
		</span>
		
		<!-- description -->
		<span v-if="rule.description != null && rule.description.length > 0" class="grey small">（{{rule.description}}）</span>
		
		<a href="" v-if="rule.err != null && rule.err.length > 0" @click.prevent="showErr(rule.err)" style="color: #db2828; opacity: 1; border-bottom: 1px #db2828 dashed; margin-left: 0.5em">规则错误</a>
	</div>
</div>`
})

// 缓存条件列表
Vue.component("http-cache-refs-box", {
	props: ["v-cache-refs"],
	data: function () {
		let refs = this.vCacheRefs
		if (refs == null) {
			refs = []
		}
		return {
			refs: refs
		}
	},
	methods: {
		timeUnitName: function (unit) {
			switch (unit) {
				case "ms":
					return "毫秒"
				case "second":
					return "秒"
				case "minute":
					return "分钟"
				case "hour":
					return "小时"
				case "day":
					return "天"
				case "week":
					return "周 "
			}
			return unit
		}
	},
	template: `<div>
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<p class="comment" v-if="refs.length == 0">暂时还没有缓存条件。</p>
	<div v-show="refs.length > 0">
		<table class="ui table selectable celled">
			<thead>
				<tr>
					<th>缓存条件</th>
					<th class="width6">缓存时间</th>
				</tr>
				<tr v-for="(cacheRef, index) in refs">
					<td :class="{'color-border': cacheRef.conds != null && cacheRef.conds.connector == 'and', disabled: !cacheRef.isOn}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds" :class="{disabled: !cacheRef.isOn}" v-if="cacheRef.conds != null && cacheRef.conds.groups != null"></http-request-conds-view>
							<http-request-cond-view :v-cond="cacheRef.simpleCond" v-if="cacheRef.simpleCond != null"></http-request-cond-view>
						
						<!-- 特殊参数 -->
						<grey-label v-if="cacheRef.key != null && cacheRef.key.indexOf('\${args}') < 0">忽略URI参数</grey-label>
						<grey-label v-if="cacheRef.minSize != null && cacheRef.minSize.count > 0">
							{{cacheRef.minSize.count}}{{cacheRef.minSize.unit}}
							<span v-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">- {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</span>
						</grey-label>
						<grey-label v-else-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">0 - {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</grey-label>
						<grey-label v-if="cacheRef.methods != null && cacheRef.methods.length > 0">{{cacheRef.methods.join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.expiresTime != null && cacheRef.expiresTime.isPrior && cacheRef.expiresTime.isOn">Expires</grey-label>
						<grey-label v-if="cacheRef.status != null && cacheRef.status.length > 0 && (cacheRef.status.length > 1 || cacheRef.status[0] != 200)">状态码：{{cacheRef.status.map(function(v) {return v.toString()}).join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.allowPartialContent">分片缓存</grey-label>
						<grey-label v-if="cacheRef.alwaysForwardRangeRequest">Range回源</grey-label>
						<grey-label v-if="cacheRef.enableIfNoneMatch">If-None-Match</grey-label>
						<grey-label v-if="cacheRef.enableIfModifiedSince">If-Modified-Since</grey-label>
						<grey-label v-if="cacheRef.enableReadingOriginAsync">支持异步</grey-label>
					</td>
					<td :class="{disabled: !cacheRef.isOn}">
						<span v-if="!cacheRef.isReverse">{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</span>
						<span v-else class="red">不缓存</span>
					</td>
				</tr>
			</thead>
		</table>
	</div>
	<div class="margin"></div>
</div>`
})

Vue.component("ssl-certs-box", {
	props: [
		"v-certs", // 证书列表
		"v-cert", // 单个证书
		"v-protocol", // 协议：https|tls
		"v-view-size", // 弹窗尺寸：normal, mini
		"v-single-mode", // 单证书模式
		"v-description", // 描述文字
		"v-domains", // 搜索的域名列表或者函数
		"v-user-id" // 用户ID
	],
	data: function () {
		let certs = this.vCerts
		if (certs == null) {
			certs = []
		}
		if (this.vCert != null) {
			certs.push(this.vCert)
		}

		let description = this.vDescription
		if (description == null || typeof (description) != "string") {
			description = ""
		}

		return {
			certs: certs,
			description: description
		}
	},
	methods: {
		certIds: function () {
			return this.certs.map(function (v) {
				return v.id
			})
		},
		// 删除证书
		removeCert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前网站不再使用此证书。", function () {
				that.certs.$remove(index)
			})
		},

		// 选择证书
		selectCert: function () {
			let that = this
			let width = "54em"
			let height = "32em"
			let viewSize = this.vViewSize
			if (viewSize == null) {
				viewSize = "normal"
			}
			if (viewSize == "mini") {
				width = "35em"
				height = "20em"
			}

			let searchingDomains = []
			if (this.vDomains != null) {
				if (typeof this.vDomains == "function") {
					let resultDomains = this.vDomains()
					if (resultDomains != null && typeof resultDomains == "object" && (resultDomains instanceof Array)) {
						searchingDomains = resultDomains
					}
				} else if (typeof this.vDomains == "object" && (this.vDomains instanceof Array)) {
					searchingDomains = this.vDomains
				}
				if (searchingDomains.length > 10000) {
					searchingDomains = searchingDomains.slice(0, 10000)
				}
			}

			let selectedCertIds = this.certs.map(function (cert) {
				return cert.id
			})
			let userId = this.vUserId
			if (userId == null) {
				userId = 0
			}

			teaweb.popup("/servers/certs/selectPopup?viewSize=" + viewSize + "&searchingDomains=" + window.encodeURIComponent(searchingDomains.join(",")) + "&selectedCertIds=" + selectedCertIds.join(",") + "&userId=" + userId, {
				width: width,
				height: height,
				callback: function (resp) {
					if (resp.data.cert != null) {
						that.certs.push(resp.data.cert)
					}
					if (resp.data.certs != null) {
						that.certs.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 上传证书
		uploadCert: function () {
			let that = this
			let userId = this.vUserId
			if (typeof userId != "number" && typeof userId != "string") {
				userId = 0
			}
			teaweb.popup("/servers/certs/uploadPopup?userId=" + userId, {
				height: "28em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						if (resp.data.cert != null) {
							that.certs.push(resp.data.cert)
						}
						if (resp.data.certs != null) {
							that.certs.$pushAll(resp.data.certs)
						}
						that.$forceUpdate()
					})
				}
			})
		},

		// 批量上传
		uploadBatch: function () {
			let that = this
			let userId = this.vUserId
			if (typeof userId != "number" && typeof userId != "string") {
				userId = 0
			}
			teaweb.popup("/servers/certs/uploadBatchPopup?userId=" + userId, {
				callback: function (resp) {
					if (resp.data.cert != null) {
						that.certs.push(resp.data.cert)
					}
					if (resp.data.certs != null) {
						that.certs.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 格式化时间
		formatTime: function (timestamp) {
			return new Date(timestamp * 1000).format("Y-m-d")
		},

		// 判断是否显示选择｜上传按钮
		buttonsVisible: function () {
			return this.vSingleMode == null || !this.vSingleMode || this.certs == null || this.certs.length == 0
		}
	},
	template: `<div>
	<input type="hidden" name="certIdsJSON" :value="JSON.stringify(certIds())"/>
	<div v-if="certs != null && certs.length > 0">
		<div class="ui label small basic" v-for="(cert, index) in certs">
			{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-else>
		<span class="red" v-if="description.length == 0">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
		<span class="grey" v-if="description.length > 0">{{description}}</span>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-if="buttonsVisible()">
		<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
		<span class="disabled">|</span> &nbsp;
		<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
		<button class="ui button tiny" type="button" @click.prevent="uploadBatch()">批量上传证书</button> &nbsp;
	</div>
</div>`
})

Vue.component("http-host-redirect-box", {
	props: ["v-redirects"],
	mounted: function () {
		let that = this
		sortTable(function (ids) {
			let newRedirects = []
			ids.forEach(function (id) {
				that.redirects.forEach(function (redirect) {
					if (redirect.id == id) {
						newRedirects.push(redirect)
					}
				})
			})
			that.updateRedirects(newRedirects)
		})
	},
	data: function () {
		let redirects = this.vRedirects
		if (redirects == null) {
			redirects = []
		}

		let id = 0
		redirects.forEach(function (v) {
			id++
			v.id = id
		})

		return {
			redirects: redirects,
			statusOptions: [
				{"code": 301, "text": "Moved Permanently"},
				{"code": 308, "text": "Permanent Redirect"},
				{"code": 302, "text": "Found"},
				{"code": 303, "text": "See Other"},
				{"code": 307, "text": "Temporary Redirect"}
			],
			id: id
		}
	},
	methods: {
		add: function () {
			let that = this
			window.UPDATING_REDIRECT = null

			teaweb.popup("/servers/server/settings/redirects/createPopup", {
				width: "50em",
				height: "36em",
				callback: function (resp) {
					that.id++
					resp.data.redirect.id = that.id
					that.redirects.push(resp.data.redirect)
					that.change()
				}
			})
		},
		update: function (index, redirect) {
			let that = this
			window.UPDATING_REDIRECT = redirect

			teaweb.popup("/servers/server/settings/redirects/createPopup", {
				width: "50em",
				height: "36em",
				callback: function (resp) {
					resp.data.redirect.id = redirect.id
					Vue.set(that.redirects, index, resp.data.redirect)
					that.change()
				}
			})
		},
		remove: function (index) {
			let that = this
			teaweb.confirm("确定要删除这条跳转规则吗？", function () {
				that.redirects.$remove(index)
				that.change()
			})
		},
		change: function () {
			let that = this
			setTimeout(function (){
				that.$emit("change", that.redirects)
			}, 100)
		},
		updateRedirects: function (newRedirects) {
			this.redirects = newRedirects
			this.change()
		}
	},
	template: `<div>
	<input type="hidden" name="hostRedirectsJSON" :value="JSON.stringify(redirects)"/>
	
	<first-menu>
		<menu-item @click.prevent="add">[创建]</menu-item>
	</first-menu>
	<div class="margin"></div>

	<p class="comment" v-if="redirects.length == 0">暂时还没有URL跳转规则。</p>
	<div v-show="redirects.length > 0">
		<table class="ui table celled selectable" id="sortable-table">
			<thead>
				<tr>
					<th style="width: 1em"></th>
					<th>跳转前</th>
					<th style="width: 1em"></th>
					<th>跳转后</th>
					<th>HTTP状态码</th>
					<th class="two wide">状态</th>
					<th class="two op">操作</th>
				</tr>
			</thead>
			<tbody v-for="(redirect, index) in redirects" :key="redirect.id" :v-id="redirect.id">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td>
						<div v-if="redirect.type == '' || redirect.type == 'url'">
							{{redirect.beforeURL}}
							<div style="margin-top: 0.4em">
								<grey-label><strong>URL跳转</strong></grey-label>
								<grey-label v-if="redirect.matchPrefix">匹配前缀</grey-label>
								<grey-label v-if="redirect.matchRegexp">正则匹配</grey-label>
								<grey-label v-if="!redirect.matchPrefix && !redirect.matchRegexp">精准匹配</grey-label>
								<grey-label v-if="redirect.exceptDomains != null && redirect.exceptDomains.length > 0" v-for="domain in redirect.exceptDomains">排除:{{domain}}</grey-label>
								<grey-label v-if="redirect.onlyDomains != null && redirect.onlyDomains.length > 0" v-for="domain in redirect.onlyDomains">仅限:{{domain}}</grey-label>
							</div>
						</div>
						<div v-if="redirect.type == 'domain'">
							<span v-if="redirect.domainsAll">所有域名</span>
							<span v-if="!redirect.domainsAll && redirect.domainsBefore != null">
								<span v-if="redirect.domainsBefore.length == 1">{{redirect.domainsBefore[0]}}</span>
								<span v-if="redirect.domainsBefore.length > 1">{{redirect.domainsBefore[0]}}等{{redirect.domainsBefore.length}}个域名</span>
							</span>
							<div style="margin-top: 0.4em">
								<grey-label><strong>域名跳转</strong></grey-label>
								<grey-label v-if="redirect.domainAfterScheme != null && redirect.domainAfterScheme.length > 0">{{redirect.domainAfterScheme}}</grey-label>
								<grey-label v-if="redirect.domainBeforeIgnorePorts">忽略端口</grey-label>
							</div>
						</div>
						<div v-if="redirect.type == 'port'">
							<span v-if="redirect.portsAll">所有端口</span>
							<span v-if="!redirect.portsAll && redirect.portsBefore != null">
								<span v-if="redirect.portsBefore.length <= 5">{{redirect.portsBefore.join(", ")}}</span>
								<span v-if="redirect.portsBefore.length > 5">{{redirect.portsBefore.slice(0, 5).join(", ")}}等{{redirect.portsBefore.length}}个端口</span>
							</span>
							<div style="margin-top: 0.4em">
								<grey-label><strong>端口跳转</strong></grey-label>
								<grey-label v-if="redirect.portAfterScheme != null && redirect.portAfterScheme.length > 0">{{redirect.portAfterScheme}}</grey-label>
							</div>
						</div>
						
						<div style="margin-top: 0.5em" v-if="redirect.conds != null && redirect.conds.groups != null && redirect.conds.groups.length > 0">
							<grey-label>匹配条件</grey-label>
						</div>
					</td>
					<td nowrap="">-&gt;</td>
					<td>
						<span v-if="redirect.type == '' || redirect.type == 'url'">{{redirect.afterURL}}</span>
						<span v-if="redirect.type == 'domain'">{{redirect.domainAfter}}</span>
						<span v-if="redirect.type == 'port'">{{redirect.portAfter}}</span>
					</td>
					<td>
						<span v-if="redirect.status > 0">{{redirect.status}}</span>
						<span v-else class="disabled">默认</span>
					</td>
					<td><label-on :v-is-on="redirect.isOn"></label-on></td>
					<td>
						<a href="" @click.prevent="update(index, redirect)">修改</a> &nbsp;
						<a href="" @click.prevent="remove(index)">删除</a>	
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="redirects.length > 1">所有规则匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。</p>
	</div>
	<div class="margin"></div>
</div>`
})

// 单个缓存条件设置
Vue.component("http-cache-ref-box", {
	props: ["v-cache-ref", "v-is-reverse"],
	mounted: function () {
		this.$refs.variablesDescriber.update(this.ref.key)
		if (this.ref.simpleCond != null) {
			this.condType = this.ref.simpleCond.type
			this.changeCondType(this.ref.simpleCond.type, true)
			this.condCategory = "simple"
		} else if (this.ref.conds != null && this.ref.conds.groups != null) {
			this.condCategory = "complex"
		}
		this.changeCondCategory(this.condCategory)
	},
	data: function () {
		let ref = this.vCacheRef
		if (ref == null) {
			ref = {
				isOn: true,
				cachePolicyId: 0,
				key: "${scheme}://${host}${requestPath}${isArgs}${args}",
				life: {count: 1, unit: "day"},
				status: [200],
				maxSize: {count: 128, unit: "mb"},
				minSize: {count: 0, unit: "kb"},
				skipCacheControlValues: ["private", "no-cache", "no-store"],
				skipSetCookie: true,
				enableRequestCachePragma: false,
				conds: null, // 复杂条件
				simpleCond: null, // 简单条件
				allowChunkedEncoding: true,
				allowPartialContent: true,
				forcePartialContent: false,
				enableIfNoneMatch: false,
				enableIfModifiedSince: false,
				enableReadingOriginAsync: false,
				isReverse: this.vIsReverse,
				methods: [],
				expiresTime: {
					isPrior: false,
					isOn: false,
					overwrite: true,
					autoCalculate: true,
					duration: {count: -1, "unit": "hour"}
				}
			}
		}
		if (ref.key == null) {
			ref.key = ""
		}
		if (ref.methods == null) {
			ref.methods = []
		}

		if (ref.life == null) {
			ref.life = {count: 2, unit: "hour"}
		}
		if (ref.maxSize == null) {
			ref.maxSize = {count: 32, unit: "mb"}
		}
		if (ref.minSize == null) {
			ref.minSize = {count: 0, unit: "kb"}
		}

		let condType = "url-extension"
		let condComponent = window.REQUEST_COND_COMPONENTS.$find(function (k, v) {
			return v.type == "url-extension"
		})

		return {
			ref: ref,

			keyIgnoreArgs: typeof ref.key == "string" && ref.key.indexOf("${args}") < 0,

			moreOptionsVisible: false,

			condCategory: "simple", // 条件分类：simple|complex
			condType: condType,
			condComponent: condComponent,
			condIsCaseInsensitive: (ref.simpleCond != null) ? ref.simpleCond.isCaseInsensitive : true,

			components: window.REQUEST_COND_COMPONENTS
		}
	},
	watch: {
		keyIgnoreArgs: function (b) {
			if (typeof this.ref.key != "string") {
				return
			}
			if (b) {
				this.ref.key = this.ref.key.replace("${isArgs}${args}", "")
				return;
			}
			if (this.ref.key.indexOf("${isArgs}") < 0) {
				this.ref.key = this.ref.key + "${isArgs}"
			}
			if (this.ref.key.indexOf("${args}") < 0) {
				this.ref.key = this.ref.key + "${args}"
			}
		}
	},
	methods: {
		changeOptionsVisible: function (v) {
			this.moreOptionsVisible = v
		},
		changeLife: function (v) {
			this.ref.life = v
		},
		changeMaxSize: function (v) {
			this.ref.maxSize = v
		},
		changeMinSize: function (v) {
			this.ref.minSize = v
		},
		changeConds: function (v) {
			this.ref.conds = v
			this.ref.simpleCond = null
		},
		changeStatusList: function (list) {
			let result = []
			list.forEach(function (status) {
				let statusNumber = parseInt(status)
				if (isNaN(statusNumber) || statusNumber < 100 || statusNumber > 999) {
					return
				}
				result.push(statusNumber)
			})
			this.ref.status = result
		},
		changeMethods: function (methods) {
			this.ref.methods = methods.map(function (v) {
				return v.toUpperCase()
			})
		},
		changeKey: function (key) {
			this.$refs.variablesDescriber.update(key)
		},
		changeExpiresTime: function (expiresTime) {
			this.ref.expiresTime = expiresTime
		},

		// 切换条件类型
		changeCondCategory: function (condCategory) {
			this.condCategory = condCategory

			// resize window
			let dialog = window.parent.document.querySelector("*[role='dialog']")
			if (dialog == null) {
				return
			}
			switch (condCategory) {
				case "simple":
					dialog.style.width = "45em"
					break
				case "complex":
					let width = window.parent.innerWidth
					if (width > 1024) {
						width = 1024
					}

					dialog.style.width = width + "px"
					if (this.ref.conds != null) {
						this.ref.conds.isOn = true
					}
					break
			}
		},
		changeCondType: function (condType, isInit) {
			if (!isInit && this.ref.simpleCond != null) {
				this.ref.simpleCond.value = null
			}
			let def = this.components.$find(function (k, component) {
				return component.type == condType
			})
			if (def != null) {
				this.condComponent = def
			}
		}
	},
	template: `<tbody>
	<tr v-if="condCategory == 'simple'">
		<td class="title">缓存对象 *</td>
		<td>
			<select class="ui dropdown auto-width" name="condType" v-model="condType" @change="changeCondType(condType, false)">
				<option value="url-extension">文件扩展名</option>
				<option value="url-eq-index">首页</option>
				<option value="url-all">全站</option>
				<option value="url-prefix">URL目录前缀</option>
				<option value="url-eq">URL完整路径</option>
				<option value="url-wildcard-match">URL通配符</option>
				<option value="url-regexp">URL正则匹配</option>
				<option value="params">参数匹配</option>
			</select>
			<p class="comment"><a href="" @click.prevent="changeCondCategory('complex')">切换到复杂条件 &raquo;</a></p>
		</td>
	</tr>
	<tr v-if="condCategory == 'simple'">
		<td>{{condComponent.paramsTitle}} *</td>
		<td>
			<component :is="condComponent.component" :v-cond="ref.simpleCond" v-if="condComponent.type != 'params'"></component>
			<table class="ui table" v-if="condComponent.type == 'params'">
				<component :is="condComponent.component" :v-cond="ref.simpleCond"></component>
			</table>
		</td>
	</tr>
	<tr v-if="condCategory == 'simple' && condComponent.caseInsensitive">
		<td>不区分大小写</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" name="condIsCaseInsensitive" value="1" v-model="condIsCaseInsensitive"/>
				<label></label>
			</div>
			<p class="comment">选中后表示对比时忽略参数值的大小写。</p>
		</td>
	</tr>
	<tr v-if="condCategory == 'complex'">
		<td class="title">匹配条件分组 *</td>
		<td>
			<http-request-conds-box :v-conds="ref.conds" @change="changeConds"></http-request-conds-box>
			<p class="comment"><a href="" @click.prevent="changeCondCategory('simple')">&laquo; 切换到简单条件</a></p>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td>缓存有效期 *</td>
		<td>
			<time-duration-box :v-value="ref.life" @change="changeLife" :v-min-unit="'minute'" maxlength="4"></time-duration-box>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td>忽略URI参数</td>
		<td>
			<checkbox v-model="keyIgnoreArgs"></checkbox>
			<p class="comment">选中后，表示缓存Key中不包含URI参数（即问号（?））后面的内容。</p>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td colspan="2"><more-options-indicator @change="changeOptionsVisible"></more-options-indicator></td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>缓存Key *</td>
		<td>
			<input type="text" v-model="ref.key" @input="changeKey(ref.key)"/>
			<p class="comment">用来区分不同缓存内容的唯一Key。<request-variables-describer ref="variablesDescriber"></request-variables-describer>。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>请求方法限制</td>
		<td>
			<values-box size="5" maxlength="10" :values="ref.methods" @change="changeMethods"></values-box>
			<p class="comment">允许请求的缓存方法，默认支持所有的请求方法。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>客户端过期时间<em>（Expires）</em></td>
		<td>
			<http-expires-time-config-box :v-expires-time="ref.expiresTime" @change="changeExpiresTime"></http-expires-time-config-box>		
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>可缓存的最大内容尺寸</td>
		<td>
			<size-capacity-box :v-value="ref.maxSize" @change="changeMaxSize"></size-capacity-box>
			<p class="comment">内容尺寸如果高于此值则不缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>可缓存的最小内容尺寸</td>
		<td>
			<size-capacity-box :v-value="ref.minSize" @change="changeMinSize"></size-capacity-box>
			<p class="comment">内容尺寸如果低于此值则不缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持缓存分片内容</td>
		<td>
			<checkbox name="allowPartialContent" value="1" v-model="ref.allowPartialContent"></checkbox>
			<p class="comment">选中后，支持缓存源站返回的某个分片的内容，该内容通过<code-label>206 Partial Content</code-label>状态码返回。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse && ref.allowPartialContent && !ref.alwaysForwardRangeReques">
		<td>强制返回分片内容</td>
		<td>
			<checkbox name="forcePartialContent" value="1" v-model="ref.forcePartialContent"></checkbox>
			<p class="comment">选中后，表示无论客户端是否发送<code-label>Range</code-label>报头，都会优先尝试返回已缓存的分片内容；如果你的应用有不支持分片内容的客户端（比如有些下载软件不支持<code-label>206 Partial Content</code-label>），请务必关闭此功能。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>强制Range回源</td>
		<td>
			<checkbox v-model="ref.alwaysForwardRangeRequest"></checkbox>
			<p class="comment">选中后，表示把所有包含Range报头的请求都转发到源站，而不是尝试从缓存中读取。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>状态码列表</td>
		<td>
			<values-box name="statusList" size="3" maxlength="3" :values="ref.status" @change="changeStatusList"></values-box>
			<p class="comment">允许缓存的HTTP状态码列表。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>跳过的Cache-Control值</td>
		<td>
			<values-box name="skipResponseCacheControlValues" size="10" maxlength="100" :values="ref.skipCacheControlValues"></values-box>
			<p class="comment">当响应的Cache-Control为这些值时不缓存响应内容，而且不区分大小写。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>跳过Set-Cookie</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" value="1" v-model="ref.skipSetCookie"/>
				<label></label>
			</div>
			<p class="comment">选中后，当响应的报头中有Set-Cookie时不缓存响应内容，防止动态内容被缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持请求no-cache刷新</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" name="enableRequestCachePragma" value="1" v-model="ref.enableRequestCachePragma"/>
				<label></label>
			</div>
			<p class="comment">选中后，当请求的报头中含有Pragma: no-cache或Cache-Control: no-cache时，会跳过缓存直接读取源内容，一般仅用于调试。</p>
		</td>
	</tr>	
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>允许If-None-Match回源</td>
		<td>
			<checkbox v-model="ref.enableIfNoneMatch"></checkbox>
			<p class="comment">特殊情况下才需要开启，可能会降低缓存命中率。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>允许If-Modified-Since回源</td>
		<td>
			<checkbox v-model="ref.enableIfModifiedSince"></checkbox>
			<p class="comment">特殊情况下才需要开启，可能会降低缓存命中率。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>允许异步读取源站</td>
		<td>
			<checkbox v-model="ref.enableReadingOriginAsync"></checkbox>
			<p class="comment">试验功能。允许客户端中断连接后，仍然继续尝试从源站读取内容并缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持分段内容</td>
		<td>
			<checkbox name="allowChunkedEncoding" value="1" v-model="ref.allowChunkedEncoding"></checkbox>
			<p class="comment">选中后，Gzip等压缩后的Chunked内容可以直接缓存，无需检查内容长度。</p>
		</td>
	</tr>
	<tr v-show="false">
		<td colspan="2"><input type="hidden" name="cacheRefJSON" :value="JSON.stringify(ref)"/></td>
	</tr>
</tbody>`
})

// 请求限制
Vue.component("http-request-limit-config-box", {
	props: ["v-request-limit-config", "v-is-group", "v-is-location"],
	data: function () {
		let config = this.vRequestLimitConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				maxConns: 0,
				maxConnsPerIP: 0,
				maxBodySize: {
					count: -1,
					unit: "kb"
				},
				outBandwidthPerConn: {
					count: -1,
					unit: "kb"
				}
			}
		}
		return {
			config: config,
			maxConns: config.maxConns,
			maxConnsPerIP: config.maxConnsPerIP
		}
	},
	watch: {
		maxConns: function (v) {
			let conns = parseInt(v, 10)
			if (isNaN(conns)) {
				this.config.maxConns = 0
				return
			}
			if (conns < 0) {
				this.config.maxConns = 0
			} else {
				this.config.maxConns = conns
			}
		},
		maxConnsPerIP: function (v) {
			let conns = parseInt(v, 10)
			if (isNaN(conns)) {
				this.config.maxConnsPerIP = 0
				return
			}
			if (conns < 0) {
				this.config.maxConnsPerIP = 0
			} else {
				this.config.maxConnsPerIP = conns
			}
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		}
	},
	template: `<div>
	<input type="hidden" name="requestLimitJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用请求限制</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>最大并发连接数</td>
				<td>
					<input type="text" maxlength="6" v-model="maxConns"/>
					<p class="comment">当前网站最大并发连接数，超出此限制则响应用户<code-label>429</code-label>代码。为0表示不限制。</p>
				</td>
			</tr>
			<tr>
				<td>单IP最大并发连接数</td>
				<td>
					<input type="text" maxlength="6" v-model="maxConnsPerIP"/>
					<p class="comment">单IP最大连接数，统计单个IP总连接数时不区分网站，超出此限制则响应用户<code-label>429</code-label>代码。为0表示不限制。<span v-if="maxConnsPerIP > 0 && maxConnsPerIP <= 3" class="red">当前设置的并发连接数过低，可能会影响正常用户访问，建议不小于3。</span></p>
				</td>
			</tr>
			<tr>
				<td>单连接带宽限制</td>
				<td>
					<size-capacity-box :v-value="config.outBandwidthPerConn" :v-supported-units="['byte', 'kb', 'mb']"></size-capacity-box>
					<p class="comment">客户端单个请求每秒可以读取的下行流量。</p>
				</td>
			</tr>
			<tr>
				<td>单请求最大尺寸</td>
				<td>
					<size-capacity-box :v-value="config.maxBodySize" :v-supported-units="['byte', 'kb', 'mb', 'gb']"></size-capacity-box>
					<p class="comment">单个请求能发送的最大内容尺寸。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-header-replace-values", {
	props: ["v-replace-values"],
	data: function () {
		let values = this.vReplaceValues
		if (values == null) {
			values = []
		}
		return {
			values: values,
			isAdding: false,
			addingValue: {"pattern": "", "replacement": "", "isCaseInsensitive": false, "isRegexp": false}
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.pattern.focus()
			})
		},
		remove: function (index) {
			this.values.$remove(index)
		},
		confirm: function () {
			let that = this
			if (this.addingValue.pattern.length == 0) {
				teaweb.warn("替换前内容不能为空", function () {
					that.$refs.pattern.focus()
				})
				return
			}

			this.values.push(this.addingValue)
			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.addingValue = {"pattern": "", "replacement": "", "isCaseInsensitive": false, "isRegexp": false}
		}
	},
	template: `<div>
	<input type="hidden" name="replaceValuesJSON" :value="JSON.stringify(values)"/>
	<div>
		<div v-for="(value, index) in values" class="ui label small" style="margin-bottom: 0.5em">
			<var>{{value.pattern}}</var><sup v-if="value.isCaseInsensitive" title="不区分大小写"><i class="icon info tiny"></i></sup> =&gt; <var v-if="value.replacement.length > 0">{{value.replacement}}</var><var v-else><span class="small grey">[空]</span></var>
			<a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="isAdding">
		<table class="ui table">
			<tr>
				<td class="title">替换前内容 *</td>
				<td><input type="text" v-model="addingValue.pattern" placeholder="替换前内容" ref="pattern" @keyup.enter="confirm()" @keypress.enter.prevent="1"/></td>
			</tr>	
			<tr>
				<td>替换后内容</td>
				<td><input type="text" v-model="addingValue.replacement" placeholder="替换后内容" @keyup.enter="confirm()" @keypress.enter.prevent="1"/></td>
			</tr>
			<tr>
				<td>是否忽略大小写</td>
				<td>
					<checkbox v-model="addingValue.isCaseInsensitive"></checkbox>
				</td>
			</tr>
		</table>

		<div>
			<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button> &nbsp;
			<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="!isAdding">
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`
})

// 浏览条件列表
Vue.component("http-request-conds-view", {
	props: ["v-conds"],
	data: function () {
		let conds = this.vConds
		if (conds == null) {
			conds = {
				isOn: true,
				connector: "or",
				groups: []
			}
		}
		if (conds.groups == null) {
			conds.groups = []
		}

		let that = this
		conds.groups.forEach(function (group) {
			group.conds.forEach(function (cond) {
				cond.typeName = that.typeName(cond)
			})
		})

		return {
			initConds: conds
		}
	},
	computed: {
		// 之所以使用computed，是因为需要动态更新
		conds: function () {
			return this.initConds
		}
	},
	methods: {
		typeName: function (cond) {
			let c = window.REQUEST_COND_COMPONENTS.$find(function (k, v) {
				return v.type == cond.type
			})
			if (c != null) {
				return c.name;
			}
			return cond.param + " " + cond.operator
		},
		updateConds: function (conds) {
			this.initConds = conds
		},
		notifyChange: function () {
			let that = this
			if (this.initConds.groups != null) {
				this.initConds.groups.forEach(function (group) {
					group.conds.forEach(function (cond) {
						cond.typeName = that.typeName(cond)
					})
				})
				this.$forceUpdate()
			}
		}
	},
	template: `<div>
		<div v-if="conds.groups.length > 0">
			<div v-for="(group, groupIndex) in conds.groups">
				<var v-for="(cond, index) in group.conds" style="font-style: normal;display: inline-block; margin-bottom:0.5em">
					<span class="ui label small basic" style="line-height: 1.5">
						<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
						<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{cond.typeName}}: </var>
						{{cond.value}}
						<sup v-if="cond.isCaseInsensitive" title="不区分大小写"><i class="icon info small"></i></sup>
					</span>
					
					<var v-if="index < group.conds.length - 1"> {{group.connector}} &nbsp;</var>
				</var>
				<div class="ui divider" v-if="groupIndex != conds.groups.length - 1" style="margin-top:0.3em;margin-bottom:0.5em"></div>
				<div>
					<span class="ui label tiny olive" v-if="group.description != null && group.description.length > 0">{{group.description}}</span>
				</div>
			</div>	
		</div>
	</div>	
</div>`
})

Vue.component("http-firewall-config-box", {
	props: ["v-firewall-config", "v-is-location", "v-is-group", "v-firewall-policy"],
	data: function () {
		let firewall = this.vFirewallConfig
		if (firewall == null) {
			firewall = {
				isPrior: false,
				isOn: false,
				firewallPolicyId: 0,
				ignoreGlobalRules: false,
				defaultCaptchaType: "none"
			}
		}

		if (firewall.defaultCaptchaType == null || firewall.defaultCaptchaType.length == 0) {
			firewall.defaultCaptchaType = "none"
		}

		let allCaptchaTypes = window.WAF_CAPTCHA_TYPES.$copy()

		// geetest
		let geeTestIsOn = false
		if (this.vFirewallPolicy != null && this.vFirewallPolicy.captchaAction != null && this.vFirewallPolicy.captchaAction.geeTestConfig != null) {
			geeTestIsOn = this.vFirewallPolicy.captchaAction.geeTestConfig.isOn
		}

		// 如果没有启用geetest，则还原
		if (!geeTestIsOn && firewall.defaultCaptchaType == "geetest") {
			firewall.defaultCaptchaType = "none"
		}

		return {
			firewall: firewall,
			moreOptionsVisible: false,
			execGlobalRules: !firewall.ignoreGlobalRules,
			captchaTypes: allCaptchaTypes,
			geeTestIsOn: geeTestIsOn
		}
	},
	watch: {
		execGlobalRules: function (v) {
			this.firewall.ignoreGlobalRules = !v
		}
	},
	methods: {
		changeOptionsVisible: function (v) {
			this.moreOptionsVisible = v
		}
	},
	template: `<div>
	<input type="hidden" name="firewallJSON" :value="JSON.stringify(firewall)"/>
	
	<table class="ui table selectable definition" v-show="!vIsGroup">
		<tr>
			<td class="title">全局WAF策略</td>
			<td>
				<div v-if="vFirewallPolicy != null">{{vFirewallPolicy.name}} <span v-if="vFirewallPolicy.modeInfo != null">&nbsp; <span :class="{green: vFirewallPolicy.modeInfo.code == 'defend', blue: vFirewallPolicy.modeInfo.code == 'observe', grey: vFirewallPolicy.modeInfo.code == 'bypass'}">[{{vFirewallPolicy.modeInfo.name}}]</span>&nbsp;</span> <link-icon :href="'/servers/components/waf/policy?firewallPolicyId=' + vFirewallPolicy.id"></link-icon>
					<p class="comment">当前网站所在集群的设置。</p>
				</div>
				<span v-else class="red">当前集群没有设置WAF策略，当前配置无法生效。</span>
			</td>
		</tr>
	</table>
	
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="firewall" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || firewall.isPrior">
			<tr>
				<td class="title">启用Web防火墙</td>
				<td>
					<checkbox v-model="firewall.isOn"></checkbox>
					<p class="comment">选中后，表示启用当前网站的WAF功能。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeOptionsVisible" v-show="firewall.isOn"></more-options-tbody>
		<tbody v-show="moreOptionsVisible">
			<tr>
				<td>人机识别验证方式</td>
				<td>
					<select class="ui dropdown auto-width" v-model="firewall.defaultCaptchaType">
						<option value="none">默认</option>
						<option v-for="captchaType in captchaTypes" v-if="captchaType.code != 'geetest' || geeTestIsOn" :value="captchaType.code">{{captchaType.name}}</option>
					</select>
					<p class="comment" v-if="firewall.defaultCaptchaType == 'none'">使用系统默认的设置。你需要在入站规则中添加规则集来决定哪些请求需要人机识别验证。</p>
					<p class="comment" v-for="captchaType in captchaTypes" v-if="captchaType.code == firewall.defaultCaptchaType">{{captchaType.description}}你需要在入站规则中添加规则集来决定哪些请求需要人机识别验证。</p>
				</td>
			</tr>
			<tr>
				<td>启用系统全局规则</td>
				<td>
					<checkbox v-model="execGlobalRules"></checkbox>
					<p class="comment">选中后，表示使用系统全局WAF策略中定义的规则。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

// 指标图表
Vue.component("metric-chart", {
	props: ["v-chart", "v-stats", "v-item", "v-column" /** in column? **/],
	mounted: function () {
		this.load()
	},
	data: function () {
		let stats = this.vStats
		if (stats == null) {
			stats = []
		}
		if (stats.length > 0) {
			let sum = stats.$sum(function (k, v) {
				return v.value
			})
			if (sum < stats[0].total) {
				if (this.vChart.type == "pie") {
					stats.push({
						keys: ["其他"],
						value: stats[0].total - sum,
						total: stats[0].total,
						time: stats[0].time
					})
				}
			}
		}
		if (this.vChart.maxItems > 0) {
			stats = stats.slice(0, this.vChart.maxItems)
		} else {
			stats = stats.slice(0, 10)
		}

		stats.$rsort(function (v1, v2) {
			return v1.value - v2.value
		})

		let widthPercent = 100
		if (this.vChart.widthDiv > 0) {
			widthPercent = 100 / this.vChart.widthDiv
		}

		return {
			chart: this.vChart,
			stats: stats,
			item: this.vItem,
			width: widthPercent + "%",
			chartId: "metric-chart-" + this.vChart.id,
			valueTypeName: (this.vItem != null && this.vItem.valueTypeName != null && this.vItem.valueTypeName.length > 0) ? this.vItem.valueTypeName : ""
		}
	},
	methods: {
		load: function () {
			var el = document.getElementById(this.chartId)
			if (el == null || el.offsetWidth == 0 || el.offsetHeight == 0) {
				setTimeout(this.load, 100)
			} else {
				this.render(el)
			}
		},
		render: function (el) {
			let chart = echarts.init(el)
			window.addEventListener("resize", function () {
				chart.resize()
			})
			switch (this.chart.type) {
				case "pie":
					this.renderPie(chart)
					break
				case "bar":
					this.renderBar(chart)
					break
				case "timeBar":
					this.renderTimeBar(chart)
					break
				case "timeLine":
					this.renderTimeLine(chart)
					break
				case "table":
					this.renderTable(chart)
					break
			}
		},
		renderPie: function (chart) {
			let values = this.stats.map(function (v) {
				return {
					name: v.keys[0],
					value: v.value
				}
			})
			let that = this
			chart.setOption({
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let percent = 0
						if (stat.total > 0) {
							percent = Math.round((stat.value * 100 / stat.total) * 100) / 100
						}
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
							case "count":
								value = teaweb.formatNumber(value)
								break
						}
						return stat.keys[0] + "<br/>" + that.valueTypeName + ": " + value + "<br/>占比：" + percent + "%"
					}
				},
				series: [
					{
						name: name,
						type: "pie",
						data: values,
						areaStyle: {},
						color: ["#9DD3E8", "#B2DB9E", "#F39494", "#FBD88A", "#879BD7"]
					}
				]
			})
		},
		renderTimeBar: function (chart) {
			this.stats.$sort(function (v1, v2) {
				return (v1.time < v2.time) ? -1 : 1
			})
			let values = this.stats.map(function (v) {
				return v.value
			})

			let axis = {unit: "", divider: 1}
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}

			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return that.formatTime(v.time)
					})
				},
				yAxis: {
					axisLabel: {
						formatter: function (value) {
							return value + axis.unit
						}
					}
				},
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
						}
						return that.formatTime(stat.time) + ": " + value
					}
				},
				grid: {
					left: 50,
					top: 10,
					right: 20,
					bottom: 25
				},
				series: [
					{
						name: name,
						type: "bar",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: teaweb.DefaultChartColor
						},
						areaStyle: {},
						barWidth: "10em"
					}
				]
			})
		},
		renderTimeLine: function (chart) {
			this.stats.$sort(function (v1, v2) {
				return (v1.time < v2.time) ? -1 : 1
			})
			let values = this.stats.map(function (v) {
				return v.value
			})

			let axis = {unit: "", divider: 1}
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}

			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return that.formatTime(v.time)
					})
				},
				yAxis: {
					axisLabel: {
						formatter: function (value) {
							return value + axis.unit
						}
					}
				},
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
						}
						return that.formatTime(stat.time) + ": " + value
					}
				},
				grid: {
					left: 50,
					top: 10,
					right: 20,
					bottom: 25
				},
				series: [
					{
						name: name,
						type: "line",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: teaweb.DefaultChartColor
						},
						areaStyle: {}
					}
				]
			})
		},
		renderBar: function (chart) {
			let values = this.stats.map(function (v) {
				return v.value
			})
			let axis = {unit: "", divider: 1}
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}
			let bottom = 24
			let rotate = 0
			let result = teaweb.xRotation(chart, this.stats.map(function (v) {
				return v.keys[0]
			}))
			if (result != null) {
				bottom = result[0]
				rotate = result[1]
			}
			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return v.keys[0]
					}),
					axisLabel: {
						interval: 0,
						rotate: rotate
					}
				},
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let percent = 0
						if (stat.total > 0) {
							percent = Math.round((stat.value * 100 / stat.total) * 100) / 100
						}
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
							case "count":
								value = teaweb.formatNumber(value)
								break
						}
						return stat.keys[0] + "<br/>" + that.valueTypeName + "：" + value + "<br/>占比：" + percent + "%"
					}
				},
				yAxis: {
					axisLabel: {
						formatter: function (value) {
							return value + axis.unit
						}
					}
				},
				grid: {
					left: 40,
					top: 10,
					right: 20,
					bottom: bottom
				},
				series: [
					{
						name: name,
						type: "bar",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: teaweb.DefaultChartColor
						},
						areaStyle: {},
						barWidth: "10em"
					}
				]
			})

			if (this.item.keys != null) {
				// IP相关操作
				if (this.item.keys.$contains("${remoteAddr}")) {
					let that = this
					chart.on("click", function (args) {
						let index = that.item.keys.$indexesOf("${remoteAddr}")[0]
						let value = that.stats[args.dataIndex].keys[index]
						teaweb.popup("/servers/ipbox?ip=" + value, {
							width: "50em",
							height: "30em"
						})
					})
				}
			}
		},
		renderTable: function (chart) {
			let table = `<table class="ui table celled">
	<thead>
		<tr>
			<th>对象</th>
			<th>数值</th>
			<th>占比</th>
		</tr>
	</thead>`
			let that = this
			this.stats.forEach(function (v) {
				let value = v.value
				switch (that.item.valueType) {
					case "byte":
						value = teaweb.formatBytes(value)
						break
				}
				table += "<tr><td>" + v.keys[0] + "</td><td>" + value + "</td>"
				let percent = 0
				if (v.total > 0) {
					percent = Math.round((v.value * 100 / v.total) * 100) / 100
				}
				table += "<td><div class=\"ui progress blue\"><div class=\"bar\" style=\"min-width: 0; height: 4px; width: " + percent + "%\"></div></div>" + percent + "%</td>"
				table += "</tr>"
			})

			table += `</table>`
			document.getElementById(this.chartId).innerHTML = table
		},
		formatTime: function (time) {
			if (time == null) {
				return ""
			}
			switch (this.item.periodUnit) {
				case "month":
					return time.substring(0, 4) + "-" + time.substring(4, 6)
				case "week":
					return time.substring(0, 4) + "-" + time.substring(4, 6)
				case "day":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8)
				case "hour":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8) + " " + time.substring(8, 10)
				case "minute":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8) + " " + time.substring(8, 10) + ":" + time.substring(10, 12)
			}
			return time
		}
	},
	template: `<div style="float: left" :style="{'width': this.vColumn ?  '' : width}" :class="{'ui column':this.vColumn}">
	<h4>{{chart.name}} <span>（{{valueTypeName}}）</span></h4>
	<div class="ui divider"></div>
	<div style="height: 14em; padding-bottom: 1em; " :id="chartId" :class="{'scroll-box': chart.type == 'table'}"></div>
</div>`
})

Vue.component("metric-board", {
	template: `<div><slot></slot></div>`
})

Vue.component("http-cache-config-box", {
	props: ["v-cache-config", "v-is-location", "v-is-group", "v-cache-policy", "v-web-id"],
	data: function () {
		let cacheConfig = this.vCacheConfig
		if (cacheConfig == null) {
			cacheConfig = {
				isPrior: false,
				isOn: false,
				addStatusHeader: true,
				addAgeHeader: false,
				enableCacheControlMaxAge: false,
				cacheRefs: [],
				purgeIsOn: false,
				purgeKey: "",
				disablePolicyRefs: false
			}
		}
		if (cacheConfig.cacheRefs == null) {
			cacheConfig.cacheRefs = []
		}

		let maxBytes = null
		if (this.vCachePolicy != null && this.vCachePolicy.maxBytes != null) {
			maxBytes = this.vCachePolicy.maxBytes
		}

		// key
		if (cacheConfig.key == null) {
			// use Vue.set to activate vue events
			Vue.set(cacheConfig, "key", {
				isOn: false,
				scheme: "https",
				host: ""
			})
		}

		return {
			cacheConfig: cacheConfig,
			moreOptionsVisible: false,
			enablePolicyRefs: !cacheConfig.disablePolicyRefs,
			maxBytes: maxBytes,

			searchBoxVisible: false,
			searchKeyword: "",

			keyOptionsVisible: false
		}
	},
	watch: {
		enablePolicyRefs: function (v) {
			this.cacheConfig.disablePolicyRefs = !v
		},
		searchKeyword: function (v) {
			this.$refs.cacheRefsConfigBoxRef.search(v)
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.cacheConfig.isPrior) && this.cacheConfig.isOn
		},
		isPlus: function () {
			return Tea.Vue.teaIsPlus
		},
		generatePurgeKey: function () {
			let r = Math.random().toString() + Math.random().toString()
			let s = r.replace(/0\./g, "")
				.replace(/\./g, "")
			let result = ""
			for (let i = 0; i < s.length; i++) {
				result += String.fromCharCode(parseInt(s.substring(i, i + 1)) + ((Math.random() < 0.5) ? "a" : "A").charCodeAt(0))
			}
			this.cacheConfig.purgeKey = result
		},
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},
		changeStale: function (stale) {
			this.cacheConfig.stale = stale
		},

		showSearchBox: function () {
			this.searchBoxVisible = !this.searchBoxVisible
			if (this.searchBoxVisible) {
				let that = this
				setTimeout(function () {
					that.$refs.searchBox.focus()
				})
			} else {
				this.searchKeyword = ""
			}
		}
	},
	template: `<div>
	<input type="hidden" name="cacheJSON" :value="JSON.stringify(cacheConfig)"/>
	
	<table class="ui table definition selectable" v-show="!vIsGroup">
		<tr>
			<td class="title">全局缓存策略</td>
			<td>
				<div v-if="vCachePolicy != null">{{vCachePolicy.name}} <link-icon :href="'/servers/components/cache/policy?cachePolicyId=' + vCachePolicy.id"></link-icon>
					<p class="comment">使用当前网站所在集群的设置。</p>
				</div>
				<span v-else class="red">当前集群没有设置缓存策略，当前配置无法生效。</span>
			</td>
		</tr>
	</table>
	
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="cacheConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || cacheConfig.isPrior">
			<tr>
				<td class="title">启用缓存</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="cacheConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && !vIsGroup">
			<tr>
				<td>缓存主域名</td>
				<td>
					<div v-show="!cacheConfig.key.isOn">默认 &nbsp; <a href="" @click.prevent="keyOptionsVisible = !keyOptionsVisible"><span class="small">[修改]</span></a></div>
					<div v-show="cacheConfig.key.isOn">使用主域名：{{cacheConfig.key.scheme}}://{{cacheConfig.key.host}} &nbsp;  <a href="" @click.prevent="keyOptionsVisible = !keyOptionsVisible"><span class="small">[修改]</span></a></div>
					<div v-show="keyOptionsVisible" style="margin-top: 1em">
						<div class="ui divider"></div>
						<table class="ui table definition">
							<tr>
								<td class="title">启用主域名</td>
								<td><checkbox v-model="cacheConfig.key.isOn"></checkbox>
									<p class="comment">启用主域名后，所有缓存键值中的协议和域名部分都会修改为主域名，用来实现缓存不区分域名。</p>
								</td>
							</tr>	
							<tr v-show="cacheConfig.key.isOn">
								<td>主域名 *</td>
								<td>
									<div class="ui fields inline">
										<div class="ui field">
											<select class="ui dropdown" v-model="cacheConfig.key.scheme">
												<option value="https">https://</option>
												<option value="http">http://</option>
											</select>
										</div>
										<div class="ui field">
											<input type="text" v-model="cacheConfig.key.host" placeholder="example.com" @keyup.enter="keyOptionsVisible = false" @keypress.enter.prevent="1"/>
										</div>
									</div>
									<p class="comment">此域名<strong>必须</strong>是当前网站已绑定域名，在刷新缓存时也需要使用此域名。</p>
								</td>
							</tr>
						</table>
						<button class="ui button tiny" type="button" @click.prevent="keyOptionsVisible = false">完成</button>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td colspan="2">
					<a href="" @click.prevent="showMoreOptions"><span v-if="moreOptionsVisible">收起选项</span><span v-else>更多选项</span><i class="icon angle" :class="{up: moreOptionsVisible, down:!moreOptionsVisible}"></i></a>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>使用默认缓存条件</td>
				<td>	
					<checkbox v-model="enablePolicyRefs"></checkbox>
					<p class="comment">选中后使用系统全局缓存策略中已经定义的默认缓存条件。</p>
				</td>
			</tr>
			<tr>
				<td>添加X-Cache报头</td>
				<td>
					<checkbox v-model="cacheConfig.addStatusHeader"></checkbox>
					<p class="comment">选中后自动在响应报头中增加<code-label>X-Cache: BYPASS|MISS|HIT|PURGE</code-label>；在浏览器端查看X-Cache值时请先禁用浏览器缓存，避免影响观察。</p>
				</td>
			</tr>
			<tr>
				<td>添加Age Header</td>
				<td>
					<checkbox v-model="cacheConfig.addAgeHeader"></checkbox>
					<p class="comment">选中后自动在响应Header中增加<code-label>Age: [存活时间秒数]</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>支持源站控制有效时间</td>
				<td>
					<checkbox v-model="cacheConfig.enableCacheControlMaxAge"></checkbox>
					<p class="comment">选中后表示支持源站在Header中设置的<code-label>Cache-Control: max-age=[有效时间秒数]</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td class="color-border">允许PURGE</td>
				<td>
					<checkbox v-model="cacheConfig.purgeIsOn"></checkbox>
					<p class="comment">允许使用PURGE方法清除某个URL缓存。</p>
				</td>
			</tr>
			<tr v-show="cacheConfig.purgeIsOn">
				<td class="color-border">PURGE Key *</td>
				<td>
					<input type="text" maxlength="200" v-model="cacheConfig.purgeKey"/>
					<p class="comment"><a href="" @click.prevent="generatePurgeKey">[随机生成]</a>。需要在PURGE方法调用时加入<code-label>X-Edge-Purge-Key: {{cacheConfig.purgeKey}}</code-label> Header。只能包含字符、数字、下划线。</p>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-if="isOn() && moreOptionsVisible && isPlus()">
		<h4>过时缓存策略</h4>
		<http-cache-stale-config :v-cache-stale-config="cacheConfig.stale" @change="changeStale"></http-cache-stale-config>
	</div>
	
	<div v-show="isOn()">
		<submit-btn></submit-btn>
		<div class="ui divider"></div>
	</div>
	
	<div v-show="isOn()" style="margin-top: 1em">
		<h4 style="position: relative">缓存条件 &nbsp; <a href="" style="font-size: 0.8em" @click.prevent="$refs.cacheRefsConfigBoxRef.addRef(false)">[添加]</a> &nbsp; <a href="" style="font-size: 0.8em" @click.prevent="showSearchBox" v-show="!searchBoxVisible">[搜索]</a> 
			<div class="ui input small right labeled" style="position: absolute; top: -0.4em; margin-left: 0.5em; zoom: 0.9" v-show="searchBoxVisible">
				<input type="text" placeholder="搜索..." ref="searchBox"  @keypress.enter.prevent="1" @keydown.esc="showSearchBox" v-model="searchKeyword" size="20"/>
				<a href="" class="ui label blue" @click.prevent="showSearchBox"><i class="icon remove small"></i></a>
			</div>
		</h4>
		<http-cache-refs-config-box ref="cacheRefsConfigBoxRef" :v-cache-config="cacheConfig" :v-cache-refs="cacheConfig.cacheRefs" :v-web-id="vWebId" :v-max-bytes="maxBytes"></http-cache-refs-config-box>
	</div>
	<div class="margin"></div>
</div>`
})

// 通用Header长度
let defaultGeneralHeaders = ["Cache-Control", "Connection", "Date", "Pragma", "Trailer", "Transfer-Encoding", "Upgrade", "Via", "Warning"]
Vue.component("http-cond-general-header-length", {
	props: ["v-checkpoint"],
	data: function () {
		let headers = null
		let length = null

		if (window.parent.UPDATING_RULE != null) {
			let options = window.parent.UPDATING_RULE.checkpointOptions
			if (options.headers != null && Array.$isArray(options.headers)) {
				headers = options.headers
			}
			if (options.length != null) {
				length = options.length
			}
		}


		if (headers == null) {
			headers = defaultGeneralHeaders
		}

		if (length == null) {
			length = 128
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			headers: headers,
			length: length
		}
	},
	watch: {
		length: function (v) {
			let len = parseInt(v)
			if (isNaN(len)) {
				len = 0
			}
			if (len < 0) {
				len = 0
			}
			this.length = len
			this.change()
		}
	},
	methods: {
		change: function () {
			this.vCheckpoint.options = [
				{
					code: "headers",
					value: this.headers
				},
				{
					code: "length",
					value: this.length
				}
			]
		}
	},
	template: `<div>
	<table class="ui table">
		<tr>
			<td class="title">通用Header列表</td>
			<td>
				<values-box :values="headers" :placeholder="'Header'" @change="change"></values-box>
				<p class="comment">需要检查的Header列表。</p>
			</td>
		</tr>
		<tr>
			<td>Header值超出长度</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" name="" style="width: 5em" v-model="length" maxlength="6"/>
					<span class="ui label">字节</span>
				</div>
				<p class="comment">超出此长度认为匹配成功，0表示不限制。</p>
			</td>
		</tr>
	</table>
</div>`
})

// CC
Vue.component("http-firewall-checkpoint-cc", {
	props: ["v-checkpoint"],
	data: function () {
		let keys = []
		let period = 60
		let threshold = 1000
		let ignoreCommonFiles = true
		let enableFingerprint = true

		let options = {}
		if (window.parent.UPDATING_RULE != null) {
			options = window.parent.UPDATING_RULE.checkpointOptions
		}

		if (options == null) {
			options = {}
		}
		if (options.keys != null) {
			keys = options.keys
		}
		if (keys.length == 0) {
			keys = ["${remoteAddr}", "${requestPath}"]
		}
		if (options.period != null) {
			period = options.period
		}
		if (options.threshold != null) {
			threshold = options.threshold
		}
		if (options.ignoreCommonFiles != null && typeof (options.ignoreCommonFiles) == "boolean") {
			ignoreCommonFiles = options.ignoreCommonFiles
		}
		if (options.enableFingerprint != null && typeof (options.enableFingerprint) == "boolean") {
			enableFingerprint = options.enableFingerprint
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			keys: keys,
			period: period,
			threshold: threshold,
			ignoreCommonFiles: ignoreCommonFiles,
			enableFingerprint: enableFingerprint,
			options: {},
			value: threshold
		}
	},
	watch: {
		period: function () {
			this.change()
		},
		threshold: function () {
			this.change()
		},
		ignoreCommonFiles: function () {
			this.change()
		},
		enableFingerprint: function () {
			this.change()
		}
	},
	methods: {
		changeKeys: function (keys) {
			this.keys = keys
			this.change()
		},
		change: function () {
			let period = parseInt(this.period.toString())
			if (isNaN(period) || period <= 0) {
				period = 60
			}

			let threshold = parseInt(this.threshold.toString())
			if (isNaN(threshold) || threshold <= 0) {
				threshold = 1000
			}
			this.value = threshold

			let ignoreCommonFiles = this.ignoreCommonFiles
			if (typeof ignoreCommonFiles != "boolean") {
				ignoreCommonFiles = false
			}

			let enableFingerprint = this.enableFingerprint
			if (typeof enableFingerprint != "boolean") {
				enableFingerprint = true
			}

			this.vCheckpoint.options = [
				{
					code: "keys",
					value: this.keys
				},
				{
					code: "period",
					value: period,
				},
				{
					code: "threshold",
					value: threshold
				},
				{
					code: "ignoreCommonFiles",
					value: ignoreCommonFiles
				},
				{
					code: "enableFingerprint",
					value: enableFingerprint
				}
			]
		},
		thresholdTooLow: function () {
			let threshold = parseInt(this.threshold.toString())
			if (isNaN(threshold) || threshold <= 0) {
				threshold = 1000
			}
			return threshold > 0 && threshold < 5
		}
	},
	template: `<div>
	<input type="hidden" name="operator" value="gt"/>
	<input type="hidden" name="value" :value="value"/>
	<table class="ui table">
		<tr>
			<td class="title">统计对象组合 *</td>
			<td>
				<metric-keys-config-box :v-keys="keys" @change="changeKeys"></metric-keys-config-box>
			</td>
		</tr>
		<tr>
			<td>统计周期 *</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="period" style="width: 6em" maxlength="8"/>
					<span class="ui label">秒</span>
				</div>
			</td>
		</tr>
		<tr>
			<td>阈值 *</td>
			<td>
				<input type="text" v-model="threshold" style="width: 6em" maxlength="8"/>
				<p class="comment" v-if="thresholdTooLow()"><span class="red">对于网站类应用来说，当前阈值设置的太低，有可能会影响用户正常访问。</span></p>
			</td>
		</tr>
		<tr>
			<td>检查请求来源指纹</td>
			<td>
				<checkbox v-model="enableFingerprint"></checkbox>
				<p class="comment">在接收到HTTPS请求时尝试检查请求来源的指纹，用来检测代理服务和爬虫攻击；如果你在网站前面放置了别的反向代理服务，请取消此选项。</p>
			</td>
		</tr>
		<tr>
			<td>忽略常用文件</td>
			<td>
				<checkbox v-model="ignoreCommonFiles"></checkbox>
				<p class="comment">忽略js、css、jpg等常在网页里被引用的文件名，即对这些文件的访问不加入计数，可以减少误判几率。</p>
			</td>
		</tr>
	</table>
</div>`
})

// 防盗链
Vue.component("http-firewall-checkpoint-referer-block", {
	props: ["v-checkpoint"],
	data: function () {
		let allowEmpty = true
		let allowSameDomain = true
		let allowDomains = []
		let denyDomains = []
		let checkOrigin = true

		let options = {}
		if (window.parent.UPDATING_RULE != null) {
			options = window.parent.UPDATING_RULE.checkpointOptions
		}

		if (options == null) {
			options = {}
		}
		if (typeof (options.allowEmpty) == "boolean") {
			allowEmpty = options.allowEmpty
		}
		if (typeof (options.allowSameDomain) == "boolean") {
			allowSameDomain = options.allowSameDomain
		}
		if (options.allowDomains != null && typeof (options.allowDomains) == "object") {
			allowDomains = options.allowDomains
		}
		if (options.denyDomains != null && typeof (options.denyDomains) == "object") {
			denyDomains = options.denyDomains
		}
		if (typeof options.checkOrigin == "boolean") {
			checkOrigin = options.checkOrigin
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			allowEmpty: allowEmpty,
			allowSameDomain: allowSameDomain,
			allowDomains: allowDomains,
			denyDomains: denyDomains,
			checkOrigin: checkOrigin,
			options: {},
			value: 0
		}
	},
	watch: {
		allowEmpty: function () {
			this.change()
		},
		allowSameDomain: function () {
			this.change()
		},
		checkOrigin: function () {
			this.change()
		}
	},
	methods: {
		changeAllowDomains: function (values) {
			this.allowDomains = values
			this.change()
		},
		changeDenyDomains: function (values) {
			this.denyDomains = values
			this.change()
		},
		change: function () {
			this.vCheckpoint.options = [
				{
					code: "allowEmpty",
					value: this.allowEmpty
				},
				{
					code: "allowSameDomain",
					value: this.allowSameDomain,
				},
				{
					code: "allowDomains",
					value: this.allowDomains
				},
				{
					code: "denyDomains",
					value: this.denyDomains
				},
				{
					code: "checkOrigin",
					value: this.checkOrigin
				}
			]
		}
	},
	template: `<div>
	<input type="hidden" name="operator" value="eq"/>
	<input type="hidden" name="value" :value="value"/>
	<table class="ui table">
		<tr>
			<td class="title">来源域名允许为空</td>
			<td>
				<checkbox v-model="allowEmpty"></checkbox>
				<p class="comment">允许不带来源的访问。</p>
			</td>
		</tr>
		<tr>
			<td>来源域名允许一致</td>
			<td>
				<checkbox v-model="allowSameDomain"></checkbox>
				<p class="comment">允许来源域名和当前访问的域名一致，相当于在站内访问。</p>
			</td>
		</tr>
		<tr>
			<td>允许的来源域名</td>
			<td>
				<values-box :values="allowDomains" @change="changeAllowDomains"></values-box>
				<p class="comment">允许的来源域名列表，比如<code-label>example.com</code-label>（顶级域名)、<code-label>*.example.com</code-label>（example.com的所有二级域名）。单个星号<code-label>*</code-label>表示允许所有域名。</p>
			</td>
		</tr>
		<tr>
			<td>禁止的来源域名</td>
			<td>
				<values-box :values="denyDomains" @change="changeDenyDomains"></values-box>
				<p class="comment">禁止的来源域名列表，比如<code-label>example.org</code-label>（顶级域名）、<code-label>*.example.org</code-label>（example.org的所有二级域名）；除了这些禁止的来源域名外，其他域名都会被允许，除非限定了允许的来源域名。</p>
			</td>
		</tr>
		<tr>
			<td>同时检查Origin</td>
			<td>
				<checkbox v-model="checkOrigin"></checkbox>
				<p class="comment">如果请求没有指定Referer Header，则尝试检查Origin Header，多用于跨站调用。</p>
			</td>
		</tr>
	</table>
</div>`
})

Vue.component("http-access-log-partitions-box", {
	props: ["v-partition", "v-day", "v-query"],
	mounted: function () {
		let that = this
		Tea.action("/servers/logs/partitionData")
			.params({
				day: this.vDay
			})
			.success(function (resp) {
				that.partitions = []
				resp.data.partitions.reverse().forEach(function (v) {
					that.partitions.push({
						code: v,
						isDisabled: false,
						hasLogs: false
					})
				})
				if (that.partitions.length > 0) {
					if (that.vPartition == null || that.vPartition < 0) {
						that.selectedPartition = that.partitions[0].code
					}

					if (that.partitions.length > 1) {
						that.checkLogs()
					}
				}
			})
			.post()
	},
	data: function () {
		return {
			partitions: [],
			selectedPartition: this.vPartition,
			checkingPartition: 0
		}
	},
	methods: {
		url: function (p) {
			let u = window.location.toString()
			u = u.replace(/\?partition=-?\d+/, "?")
			u = u.replace(/\?requestId=-?\d+/, "?")
			u = u.replace(/&partition=-?\d+/, "")
			u = u.replace(/&requestId=-?\d+/, "")
			if (u.indexOf("?") > 0) {
				u += "&partition=" + p
			} else {
				u += "?partition=" + p
			}
			return u
		},
		disable: function (partition) {
			this.partitions.forEach(function (p) {
				if (p.code == partition) {
					p.isDisabled = true
				}
			})
		},
		checkLogs: function () {
			let that = this
			let index = this.checkingPartition
			let params = {
				partition: index
			}
			let query = this.vQuery
			if (query == null || query.length == 0) {
				return
			}
			query.split("&").forEach(function (v) {
				let param = v.split("=")
				params[param[0]] = decodeURIComponent(param[1])
			})
			Tea.action("/servers/logs/hasLogs")
				.params(params)
				.post()
				.success(function (response) {
					if (response.data.hasLogs) {
						// 因为是倒序，所以这里需要使用总长度减去index
						that.partitions[that.partitions.length - 1 - index].hasLogs = true
					}

					index++
					if (index >= that.partitions.length) {
						return
					}
					that.checkingPartition = index
					that.checkLogs()
				})
		}
	},
	template: `<div v-if="partitions.length > 1">
	<div class="ui divider" style="margin-bottom: 0"></div>
	<div class="ui menu text small blue" style="margin-bottom: 0; margin-top: 0">
		<a v-for="(p, index) in partitions" :href="url(p.code)" class="item" :class="{active: selectedPartition == p.code, disabled: p.isDisabled}">分表{{p.code+1}} <span v-if="p.hasLogs">&nbsp; <dot></dot></span> &nbsp; &nbsp; <span class="disabled" v-if="index != partitions.length - 1">|</span></a>
	</div>
	<div class="ui divider" style="margin-top: 0"></div>
</div>`
})

Vue.component("http-cache-refs-config-box", {
	props: ["v-cache-refs", "v-cache-config", "v-cache-policy-id", "v-web-id", "v-max-bytes"],
	mounted: function () {
		let that = this
		sortTable(function (ids) {
			let newRefs = []
			ids.forEach(function (id) {
				that.refs.forEach(function (ref) {
					if (ref.id == id) {
						newRefs.push(ref)
					}
				})
			})
			that.updateRefs(newRefs)
			that.change()
		})
	},
	data: function () {
		let refs = this.vCacheRefs
		if (refs == null) {
			refs = []
		}

		let maxBytes = this.vMaxBytes

		let id = 0
		refs.forEach(function (ref) {
			// preset variables
			id++
			ref.id = id
			ref.visible = true

			// check max size
			if (ref.maxSize != null && maxBytes != null && maxBytes.count > 0 && teaweb.compareSizeCapacity(ref.maxSize, maxBytes) > 0) {
				ref.overMaxSize = maxBytes
			}
		})
		return {
			refs: refs,
			id: id // 用来对条件进行排序
		}
	},
	methods: {
		addRef: function (isReverse) {
			window.UPDATING_CACHE_REF = null

			let height = window.innerHeight
			if (height > 500) {
				height = 500
			}
			let that = this
			teaweb.popup("/servers/server/settings/cache/createPopup?isReverse=" + (isReverse ? 1 : 0), {
				height: height + "px",
				callback: function (resp) {
					let newRef = resp.data.cacheRef
					if (newRef.conds == null) {
						return
					}

					that.id++
					newRef.id = that.id

					if (newRef.isReverse) {
						let newRefs = []
						let isAdded = false
						that.refs.forEach(function (v) {
							if (!v.isReverse && !isAdded) {
								newRefs.push(newRef)
								isAdded = true
							}
							newRefs.push(v)
						})
						if (!isAdded) {
							newRefs.push(newRef)
						}

						that.updateRefs(newRefs)
					} else {
						that.refs.push(newRef)
					}

					// move to bottom
					var afterChangeCallback = function () {
						setTimeout(function () {
							let rightBox = document.querySelector(".right-box")
							if (rightBox != null) {
								rightBox.scrollTo(0, isReverse ? 0 : 100000)
							}
						}, 100)
					}

					that.change(afterChangeCallback)
				}
			})
		},
		updateRef: function (index, cacheRef) {
			window.UPDATING_CACHE_REF = teaweb.clone(cacheRef)

			let height = window.innerHeight
			if (height > 500) {
				height = 500
			}
			let that = this
			teaweb.popup("/servers/server/settings/cache/createPopup", {
				height: height + "px",
				callback: function (resp) {
					resp.data.cacheRef.id = that.refs[index].id
					Vue.set(that.refs, index, resp.data.cacheRef)
					that.change()
					that.$refs.cacheRef[index].updateConds(resp.data.cacheRef.conds, resp.data.cacheRef.simpleCond)
					that.$refs.cacheRef[index].notifyChange()
				}
			})
		},
		disableRef: function (ref) {
			ref.isOn = false
			this.change()
		},
		enableRef: function (ref) {
			ref.isOn = true
			this.change()
		},
		removeRef: function (index) {
			let that = this
			teaweb.confirm("确定要删除此缓存设置吗？", function () {
				that.refs.$remove(index)
				that.change()
			})
		},
		updateRefs: function (newRefs) {
			this.refs = newRefs
			if (this.vCacheConfig != null) {
				this.vCacheConfig.cacheRefs = newRefs
			}
		},
		timeUnitName: function (unit) {
			switch (unit) {
				case "ms":
					return "毫秒"
				case "second":
					return "秒"
				case "minute":
					return "分钟"
				case "hour":
					return "小时"
				case "day":
					return "天"
				case "week":
					return "周 "
			}
			return unit
		},
		change: function (callback) {
			this.$forceUpdate()

			// 自动保存
			if (this.vCachePolicyId != null && this.vCachePolicyId > 0) { // 缓存策略
				Tea.action("/servers/components/cache/updateRefs")
					.params({
						cachePolicyId: this.vCachePolicyId,
						refsJSON: JSON.stringify(this.refs)
					})
					.post()
			} else if (this.vWebId != null && this.vWebId > 0) { // Server Web or Group Web
				Tea.action("/servers/server/settings/cache/updateRefs")
					.params({
						webId: this.vWebId,
						refsJSON: JSON.stringify(this.refs)
					})
					.success(function (resp) {
						if (resp.data.isUpdated) {
							teaweb.successToast("保存成功", null, function () {
								if (typeof callback == "function") {
									callback()
								}
							})
						}
					})
					.post()
			}
		},
		search: function (keyword) {
			if (typeof keyword != "string") {
				keyword = ""
			}

			this.refs.forEach(function (ref) {
				if (keyword.length == 0) {
					ref.visible = true
					return
				}
				ref.visible = false

				// simple cond
				if (ref.simpleCond != null && typeof ref.simpleCond.value == "string" && teaweb.match(ref.simpleCond.value, keyword)) {
					ref.visible = true
					return
				}

				// composed conds
				if (ref.conds == null || ref.conds.groups == null || ref.conds.groups.length == 0) {
					return
				}

				ref.conds.groups.forEach(function (group) {
					if (group.conds != null) {
						group.conds.forEach(function (cond) {
							if (typeof cond.value == "string" && teaweb.match(cond.value, keyword)) {
								ref.visible = true
							}
						})
					}
				})
			})
			this.$forceUpdate()
		}
	},
	template: `<div>
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<div>
		<p class="comment" v-if="refs.length == 0">暂时还没有缓存条件。</p>
		<table class="ui table selectable celled" v-show="refs.length > 0" id="sortable-table">
			<thead>
				<tr>
					<th style="width:1em"></th>
					<th>缓存条件</th>
					<th style="width: 7em">缓存时间</th>
					<th class="three op">操作</th>
				</tr>
			</thead>	
			<tbody v-for="(cacheRef, index) in refs" :key="cacheRef.id" :v-id="cacheRef.id" v-show="cacheRef.visible !== false">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td :class="{'color-border': cacheRef.conds != null && cacheRef.conds.connector == 'and', disabled: !cacheRef.isOn}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds" ref="cacheRef" :class="{disabled: !cacheRef.isOn}" v-if="cacheRef.conds != null && cacheRef.conds.groups != null"></http-request-conds-view>
						<http-request-cond-view :v-cond="cacheRef.simpleCond" ref="cacheRef" v-if="cacheRef.simpleCond != null"></http-request-cond-view>
						
						<!-- 特殊参数 -->
						<grey-label v-if="cacheRef.key != null && cacheRef.key.indexOf('\${args}') < 0">忽略URI参数</grey-label>
						
						<grey-label v-if="cacheRef.minSize != null && cacheRef.minSize.count > 0">
							{{cacheRef.minSize.count}}{{cacheRef.minSize.unit}}
							<span v-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">- {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit.toUpperCase()}}</span>
						</grey-label>
						<grey-label v-else-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">0 - {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit.toUpperCase()}}</grey-label>
						
						<grey-label v-if="cacheRef.overMaxSize != null"><span class="red">系统限制{{cacheRef.overMaxSize.count}}{{cacheRef.overMaxSize.unit.toUpperCase()}}</span> </grey-label>
						
						<grey-label v-if="cacheRef.methods != null && cacheRef.methods.length > 0">{{cacheRef.methods.join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.expiresTime != null && cacheRef.expiresTime.isPrior && cacheRef.expiresTime.isOn">Expires</grey-label>
						<grey-label v-if="cacheRef.status != null && cacheRef.status.length > 0 && (cacheRef.status.length > 1 || cacheRef.status[0] != 200)">状态码：{{cacheRef.status.map(function(v) {return v.toString()}).join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.allowPartialContent">分片缓存</grey-label>
						<grey-label v-if="cacheRef.alwaysForwardRangeRequest">Range回源</grey-label>
						<grey-label v-if="cacheRef.enableIfNoneMatch">If-None-Match</grey-label>
						<grey-label v-if="cacheRef.enableIfModifiedSince">If-Modified-Since</grey-label>
						<grey-label v-if="cacheRef.enableReadingOriginAsync">支持异步</grey-label>
					</td>
					<td :class="{disabled: !cacheRef.isOn}">
						<span v-if="!cacheRef.isReverse">{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</span>
						<span v-else class="red">不缓存</span>
					</td>
					<td>
						<a href="" @click.prevent="updateRef(index, cacheRef)">修改</a> &nbsp;
						<a href="" v-if="cacheRef.isOn" @click.prevent="disableRef(cacheRef)">暂停</a><a href="" v-if="!cacheRef.isOn" @click.prevent="enableRef(cacheRef)"><span class="red">恢复</span></a> &nbsp;
						<a href="" @click.prevent="removeRef(index)">删除</a>
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="refs.length > 1">所有条件匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。服务设置的优先级比全局缓存策略设置的优先级要高。</p>
		
		<button class="ui button tiny" @click.prevent="addRef(false)" type="button">+添加缓存条件</button> &nbsp; &nbsp; <a href="" @click.prevent="addRef(true)" style="font-size: 0.8em">+添加不缓存条件</a>
	</div>
	<div class="margin"></div>
</div>`
})

Vue.component("origin-list-box", {
	props: ["v-primary-origins", "v-backup-origins", "v-server-type", "v-params"],
	data: function () {
		return {
			primaryOrigins: this.vPrimaryOrigins,
			backupOrigins: this.vBackupOrigins
		}
	},
	methods: {
		createPrimaryOrigin: function () {
			teaweb.popup("/servers/server/settings/origins/addPopup?originType=primary&" + this.vParams, {
				width: "45em",
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		createBackupOrigin: function () {
			teaweb.popup("/servers/server/settings/origins/addPopup?originType=backup&" + this.vParams, {
				width: "45em",
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		updateOrigin: function (originId, originType) {
			teaweb.popup("/servers/server/settings/origins/updatePopup?originType=" + originType + "&" + this.vParams + "&originId=" + originId, {
				width: "45em",
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		deleteOrigin: function (originId, originAddr, originType) {
			let that = this
			teaweb.confirm("确定要删除此源站（" + originAddr + "）吗？", function () {
				Tea.action("/servers/server/settings/origins/delete?" + that.vParams + "&originId=" + originId + "&originType=" + originType)
					.post()
					.success(function () {
						teaweb.success("删除成功", function () {
							window.location.reload()
						})
					})
			})
		},
		updateOriginIsOn: function (originId, originAddr, isOn) {
			let message
			let resultMessage
			if (isOn) {
				message = "确定要启用此源站（" + originAddr + "）吗？"
				resultMessage = "启用成功"
			} else {
				message = "确定要停用此源站（" + originAddr + "）吗？"
				resultMessage = "停用成功"
			}
			let that = this
			teaweb.confirm(message, function () {
				Tea.action("/servers/server/settings/origins/updateIsOn?" + that.vParams + "&originId=" + originId + "&isOn=" + (isOn ? 1 : 0))
					.post()
					.success(function () {
						teaweb.success(resultMessage, function () {
							window.location.reload()
						})
					})
			})
		}
	},
	template: `<div>
	<h3>主要源站 <a href="" @click.prevent="createPrimaryOrigin()">[添加主要源站]</a> </h3>
	<p class="comment" v-if="primaryOrigins.length == 0">暂时还没有主要源站。</p>
	<origin-list-table v-if="primaryOrigins.length > 0" :v-origins="vPrimaryOrigins" :v-origin-type="'primary'" @delete-origin="deleteOrigin" @update-origin="updateOrigin" @update-origin-is-on="updateOriginIsOn"></origin-list-table>

	<h3>备用源站 <a href="" @click.prevent="createBackupOrigin()">[添加备用源站]</a></h3>
	<p class="comment" v-if="backupOrigins.length == 0">暂时还没有备用源站。</p>
	<origin-list-table v-if="backupOrigins.length > 0" :v-origins="backupOrigins" :v-origin-type="'backup'" @delete-origin="deleteOrigin" @update-origin="updateOrigin" @update-origin-is-on="updateOriginIsOn"></origin-list-table>
</div>`
})

Vue.component("origin-list-table", {
	props: ["v-origins", "v-origin-type"],
	data: function () {
		let hasMatchedDomains = false
		let origins = this.vOrigins
		if (origins != null && origins.length > 0) {
			origins.forEach(function (origin) {
				if (origin.domains != null && origin.domains.length > 0) {
					hasMatchedDomains = true
				}
			})
		}

		return {
			hasMatchedDomains: hasMatchedDomains
		}
	},
	methods: {
		deleteOrigin: function (originId, originAddr) {
			this.$emit("delete-origin", originId, originAddr, this.vOriginType)
		},
		updateOrigin: function (originId) {
			this.$emit("update-origin", originId, this.vOriginType)
		},
		updateOriginIsOn: function (originId, originAddr, isOn) {
			this.$emit("update-origin-is-on", originId, originAddr, isOn)
		}
	},
	template: `
<table class="ui table selectable">
	<thead>
		<tr>
			<th>源站地址</th>
			<th class="width5">权重</th>
			<th class="width6">状态</th>
			<th class="three op">操作</th>
		</tr>	
	</thead>
	<tbody>
		<tr v-for="origin in vOrigins">
			<td :class="{disabled:!origin.isOn}">
				<a href="" @click.prevent="updateOrigin(origin.id)" :class="{disabled:!origin.isOn}">{{origin.addr}} &nbsp;<i class="icon expand small"></i></a>
				<div style="margin-top: 0.3em">
					<tiny-basic-label class="grey border-grey" v-if="origin.isOSS"><i class="icon hdd outline"></i>对象存储</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.name.length > 0">{{origin.name}}</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.hasCert">证书</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.host != null && origin.host.length > 0">主机名: {{origin.host}}</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.followPort">端口跟随</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.addr != null && origin.addr.startsWith('https://') && origin.http2Enabled">HTTP/2</tiny-basic-label>
	
					<span v-if="origin.domains != null && origin.domains.length > 0"><tiny-basic-label class="grey border-grey" v-for="domain in origin.domains">匹配: {{domain}}</tiny-basic-label></span>
					<span v-else-if="hasMatchedDomains"><tiny-basic-label class="grey  border-grey">匹配: 所有域名</tiny-basic-label></span>
				</div>
			</td>
			<td :class="{disabled:!origin.isOn}">{{origin.weight}}</td>
			<td>
				<label-on :v-is-on="origin.isOn"></label-on>
			</td>
			<td>
				<a href="" @click.prevent="updateOrigin(origin.id)">修改</a> &nbsp;
				<a href="" v-if="origin.isOn" @click.prevent="updateOriginIsOn(origin.id, origin.addr, false)">停用</a><a href=""  v-if="!origin.isOn" @click.prevent="updateOriginIsOn(origin.id, origin.addr, true)"><span class="red">启用</span></a> &nbsp;
				<a href="" @click.prevent="deleteOrigin(origin.id, origin.addr)">删除</a>
			</td>
		</tr>
	</tbody>
</table>`
})

Vue.component("http-cors-header-config-box", {
	props: ["value"],
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isOn: false,
				allowMethods: [],
				allowOrigin: "",
				allowCredentials: true,
				exposeHeaders: [],
				maxAge: 0,
				requestHeaders: [],
				requestMethod: "",
				optionsMethodOnly: false
			}
		}
		if (config.allowMethods == null) {
			config.allowMethods = []
		}
		if (config.exposeHeaders == null) {
			config.exposeHeaders = []
		}

		let maxAgeSecondsString = config.maxAge.toString()
		if (maxAgeSecondsString == "0") {
			maxAgeSecondsString = ""
		}

		return {
			config: config,

			maxAgeSecondsString: maxAgeSecondsString,

			moreOptionsVisible: false
		}
	},
	watch: {
		maxAgeSecondsString: function (v) {
			let seconds = parseInt(v)
			if (isNaN(seconds)) {
				seconds = 0
			}
			this.config.maxAge = seconds
		}
	},
	methods: {
		changeMoreOptions: function (visible) {
			this.moreOptionsVisible = visible
		},
		addDefaultAllowMethods: function () {
			let that = this
			let defaultMethods = ["PUT", "GET", "POST", "DELETE", "HEAD", "OPTIONS", "PATCH"]
			defaultMethods.forEach(function (method) {
				if (!that.config.allowMethods.$contains(method)) {
					that.config.allowMethods.push(method)
				}
			})
		}
	},
	template: `<div>
	<input type="hidden" name="corsJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">启用CORS自适应跨域</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
					<p class="comment">启用后，自动在响应报头中增加对应的<code-label>Access-Control-*</code-label>相关内容。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td colspan="2"><more-options-indicator @change="changeMoreOptions"></more-options-indicator></td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn && moreOptionsVisible">
			<tr>
				<td>允许的请求方法列表</td>
				<td>
					<http-methods-box :v-methods="config.allowMethods"></http-methods-box>
					<p class="comment"><a href="" @click.prevent="addDefaultAllowMethods">[添加默认]</a>。<code-label>Access-Control-Allow-Methods</code-label>值设置。所访问资源允许使用的方法列表，不设置则表示默认为<code-label>PUT</code-label>、<code-label>GET</code-label>、<code-label>POST</code-label>、<code-label>DELETE</code-label>、<code-label>HEAD</code-label>、<code-label>OPTIONS</code-label>、<code-label>PATCH</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>预检结果缓存时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 6em" maxlength="6" v-model="maxAgeSecondsString"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment"><code-label>Access-Control-Max-Age</code-label>值设置。预检结果缓存时间，0或者不填表示使用浏览器默认设置。注意每个浏览器有不同的缓存时间上限。</p>
				</td>
			</tr>
			<tr>
				<td>允许服务器暴露的报头</td>
				<td>
					<values-box :v-values="config.exposeHeaders"></values-box>
					<p class="comment"><code-label>Access-Control-Expose-Headers</code-label>值设置。允许服务器暴露的报头，请注意报头的大小写。</p>
				</td>
			</tr>
			<tr>
				<td>实际请求方法</td>
				<td>
					<input type="text" v-model="config.requestMethod"/>
					<p class="comment"><code-label>Access-Control-Request-Method</code-label>值设置。实际请求服务器时使用的方法，比如<code-label>POST</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>仅OPTIONS有效</td>
				<td>
					<checkbox v-model="config.optionsMethodOnly"></checkbox>
					<p class="comment">选中后，表示当前CORS设置仅在OPTIONS方法请求时有效。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-firewall-policy-selector", {
	props: ["v-http-firewall-policy"],
	mounted: function () {
		let that = this
		Tea.action("/servers/components/waf/count")
			.post()
			.success(function (resp) {
				that.count = resp.data.count
			})
	},
	data: function () {
		let firewallPolicy = this.vHttpFirewallPolicy
		return {
			count: 0,
			firewallPolicy: firewallPolicy
		}
	},
	methods: {
		remove: function () {
			this.firewallPolicy = null
		},
		select: function () {
			let that = this
			teaweb.popup("/servers/components/waf/selectPopup", {
				height: "26em",
				callback: function (resp) {
					that.firewallPolicy = resp.data.firewallPolicy
				}
			})
		},
		create: function () {
			let that = this
			teaweb.popup("/servers/components/waf/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.firewallPolicy = resp.data.firewallPolicy
				}
			})
		}
	},
	template: `<div>
	<div v-if="firewallPolicy != null" class="ui label basic">
		<input type="hidden" name="httpFirewallPolicyId" :value="firewallPolicy.id"/>
		{{firewallPolicy.name}} &nbsp; <a :href="'/servers/components/waf/policy?firewallPolicyId=' + firewallPolicy.id" target="_blank" title="修改"><i class="icon pencil small"></i></a>&nbsp; <a href="" @click.prevent="remove()" title="删除"><i class="icon remove small"></i></a>
	</div>
	<div v-if="firewallPolicy == null">
		<span v-if="count > 0"><a href="" @click.prevent="select">[选择已有策略]</a> &nbsp; &nbsp; </span><a href="" @click.prevent="create">[创建新策略]</a>
	</div>
</div>`
})

// 压缩配置
Vue.component("http-optimization-config-box", {
	props: ["v-optimization-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vOptimizationConfig

		return {
			config: config,
			htmlMoreOptions: false,
			javascriptMoreOptions: false,
			cssMoreOptions: false
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		}
	},
	template: `<div>
	<input type="hidden" name="optimizationJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable" v-if="vIsLocation || vIsGroup">
		<prior-checkbox :v-config="config"></prior-checkbox>
	</table>
	
	<div v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
		<div class="margin"></div>
		<table class="ui table definition selectable">
			<tbody>
				<tr>
					<td class="title">HTML优化</td>
					<td>
						<div class="ui checkbox">
							<input type="checkbox" value="1" v-model="config.html.isOn"/>
							<label></label>
						</div>
						<p class="comment">可以自动优化HTML中包含的空白、注释、空标签等。只有文件可以缓存时才会被优化。</p>
					</td>
				</tr>
				<tr v-show="config.html.isOn">
					<td colspan="2"><more-options-indicator v-model="htmlMoreOptions"></more-options-indicator></td>
				</tr>
			</tbody>
			<tbody v-show="htmlMoreOptions">
				<tr>
					<td>HTML例外URL</td>
					<td>
						<url-patterns-box v-model="config.html.exceptURLPatterns"></url-patterns-box>
						<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
					</td>
				</tr>
				<tr>
					<td>HTML限制URL</td>
					<td>
						<url-patterns-box v-model="config.html.onlyURLPatterns"></url-patterns-box>
						<p class="comment">如果填写了限制URL，表示只对这些URL进行优化处理；如果不填则表示支持所有的URL。</p>
					</td>
				</tr>	
			</tbody>
		</table>
		
		<table class="ui table definition selectable">
			<tbody>
				<tr>
					<td class="title">Javascript优化</td>
					<td>
						<div class="ui checkbox">
							<input type="checkbox" value="1" v-model="config.javascript.isOn"/>
							<label></label>
						</div>
						<p class="comment">可以自动缩短Javascript中变量、函数名称等。只有文件可以缓存时才会被优化。</p>
					</td>
				</tr>
				<tr v-show="config.javascript.isOn">
					<td colspan="2"><more-options-indicator v-model="javascriptMoreOptions"></more-options-indicator></td>
				</tr>
			</tbody>
			<tbody v-show="javascriptMoreOptions">
				<tr>
					<td>Javascript例外URL</td>
					<td>
						<url-patterns-box v-model="config.javascript.exceptURLPatterns"></url-patterns-box>
						<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
					</td>
				</tr>
				<tr>
					<td>Javascript限制URL</td>
					<td>
						<url-patterns-box v-model="config.javascript.onlyURLPatterns"></url-patterns-box>
						<p class="comment">如果填写了限制URL，表示只对这些URL进行优化处理；如果不填则表示支持所有的URL。</p>
					</td>
				</tr>	
			</tbody>
		</table>
		
		<table class="ui table definition selectable">
			<tbody>
				<tr>
					<td class="title">CSS优化</td>
					<td>
						<div class="ui checkbox">
							<input type="checkbox" value="1" v-model="config.css.isOn"/>
							<label></label>
						</div>
						<p class="comment">可以自动去除CSS中包含的空白。只有文件可以缓存时才会被优化。</p>
					</td>
				</tr>
				<tr v-show="config.css.isOn">
					<td colspan="2"><more-options-indicator v-model="cssMoreOptions"></more-options-indicator></td>
				</tr>
			</tbody>
			<tbody v-show="cssMoreOptions">
				<tr>
					<td>CSS例外URL</td>
					<td>
						<url-patterns-box v-model="config.css.exceptURLPatterns"></url-patterns-box>
						<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
					</td>
				</tr>
				<tr>
					<td>CSS限制URL</td>
					<td>
						<url-patterns-box v-model="config.css.onlyURLPatterns"></url-patterns-box>
						<p class="comment">如果填写了限制URL，表示只对这些URL进行优化处理；如果不填则表示支持所有的URL。</p>
					</td>
				</tr>	
			</tbody>
		</table>
	</div>
	
	<div class="margin"></div>
</div>`
})

Vue.component("http-websocket-box", {
	props: ["v-websocket-ref", "v-websocket-config", "v-is-location", "v-is-group"],
	data: function () {
		let websocketRef = this.vWebsocketRef
		if (websocketRef == null) {
			websocketRef = {
				isPrior: false,
				isOn: false,
				websocketId: 0
			}
		}

		let websocketConfig = this.vWebsocketConfig
		if (websocketConfig == null) {
			websocketConfig = {
				id: 0,
				isOn: false,
				handshakeTimeout: {
					count: 30,
					unit: "second"
				},
				allowAllOrigins: true,
				allowedOrigins: [],
				requestSameOrigin: true,
				requestOrigin: ""
			}
		} else {
			if (websocketConfig.handshakeTimeout == null) {
				websocketConfig.handshakeTimeout = {
					count: 30,
					unit: "second",
				}
			}
			if (websocketConfig.allowedOrigins == null) {
				websocketConfig.allowedOrigins = []
			}
		}

		return {
			websocketRef: websocketRef,
			websocketConfig: websocketConfig,
			handshakeTimeoutCountString: websocketConfig.handshakeTimeout.count.toString(),
			advancedVisible: false
		}
	},
	watch: {
		handshakeTimeoutCountString: function (v) {
			let count = parseInt(v)
			if (!isNaN(count) && count >= 0) {
				this.websocketConfig.handshakeTimeout.count = count
			} else {
				this.websocketConfig.handshakeTimeout.count = 0
			}
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.websocketRef.isPrior) && this.websocketRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		},
		createOrigin: function () {
			let that = this
			teaweb.popup("/servers/server/settings/websocket/createOrigin", {
				height: "12.5em",
				callback: function (resp) {
					that.websocketConfig.allowedOrigins.push(resp.data.origin)
				}
			})
		},
		removeOrigin: function (index) {
			this.websocketConfig.allowedOrigins.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="websocketRefJSON" :value="JSON.stringify(websocketRef)"/>
	<input type="hidden" name="websocketJSON" :value="JSON.stringify(websocketConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="websocketRef" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="((!vIsLocation && !vIsGroup) || websocketRef.isPrior)">
			<tr>
				<td class="title">启用Websocket</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="color-border">允许所有来源域<em>（Origin）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketConfig.allowAllOrigins"/>
						<label></label>
					</div>
					<p class="comment">选中表示允许所有的来源域。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && !websocketConfig.allowAllOrigins">
			<tr>
				<td class="color-border">允许的来源域列表<em>（Origin）</em></td>
				<td>
					<div v-if="websocketConfig.allowedOrigins.length > 0">
						<div class="ui label small basic" v-for="(origin, index) in websocketConfig.allowedOrigins">
							{{origin}} <a href="" title="删除" @click.prevent="removeOrigin(index)"><i class="icon remove small"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="createOrigin()">+</button>
					<p class="comment">只允许在列表中的来源域名访问Websocket服务。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-show="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td class="color-border">传递请求来源域</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketConfig.requestSameOrigin"/>
						<label></label>
					</div>
					<p class="comment">选中后，表示把接收到的请求中的<code-label>Origin</code-label>字段传递到源站。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && advancedVisible && !websocketConfig.requestSameOrigin">
			<tr>
				<td class="color-border">指定传递的来源域</td>
				<td>
					<input type="text" v-model="websocketConfig.requestOrigin" maxlength="200"/>
					<p class="comment">指定向源站传递的<span class="ui label tiny">Origin</span>字段值。</p>
				</td>
			</tr>
		</tbody>
		<!-- TODO 这个选项暂时保留 -->
		<tbody v-show="isOn() && false">
			<tr>
				<td>握手超时时间<em>（Handshake）</em></td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" maxlength="10" v-model="handshakeTimeoutCountString" style="width:6em"/>
						</div>
						<div class="ui field">
							秒
						</div>
					</div>
					<p class="comment">0表示使用默认的时间设置。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-rewrite-rule-list", {
	props: ["v-web-id", "v-rewrite-rules"],
	mounted: function () {
		setTimeout(this.sort, 1000)
	},
	data: function () {
		let rewriteRules = this.vRewriteRules
		if (rewriteRules == null) {
			rewriteRules = []
		}
		return {
			rewriteRules: rewriteRules
		}
	},
	methods: {
		updateRewriteRule: function (rewriteRuleId) {
			teaweb.popup("/servers/server/settings/rewrite/updatePopup?webId=" + this.vWebId + "&rewriteRuleId=" + rewriteRuleId, {
				height: "26em",
				callback: function () {
					window.location.reload()
				}
			})
		},
		deleteRewriteRule: function (rewriteRuleId) {
			let that = this
			teaweb.confirm("确定要删除此重写规则吗？", function () {
				Tea.action("/servers/server/settings/rewrite/delete")
					.params({
						webId: that.vWebId,
						rewriteRuleId: rewriteRuleId
					})
					.post()
					.refresh()
			})
		},
		// 排序
		sort: function () {
			if (this.rewriteRules.length == 0) {
				return
			}
			let that = this
			sortTable(function (rowIds) {
				Tea.action("/servers/server/settings/rewrite/sort")
					.post()
					.params({
						webId: that.vWebId,
						rewriteRuleIds: rowIds
					})
					.success(function () {
						teaweb.success("保存成功")
					})
			})
		}
	},
	template: `<div>
	<div class="margin"></div>
	<p class="comment" v-if="rewriteRules.length == 0">暂时还没有重写规则。</p>
	<table class="ui table selectable" v-if="rewriteRules.length > 0" id="sortable-table">
		<thead>
			<tr>
				<th style="width:1em"></th>
				<th>匹配规则</th>
				<th>转发目标</th>
				<th>转发方式</th>
				<th class="two wide">状态</th>
				<th class="two op">操作</th>
			</tr>
		</thead>
		<tbody v-for="rule in rewriteRules" :v-id="rule.id">
			<tr>
				<td><i class="icon bars grey handle"></i></td>
				<td>{{rule.pattern}}
				<br/>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.isBreak">BREAK</http-rewrite-labels-label>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.mode == 'redirect' && rule.redirectStatus != 307">{{rule.redirectStatus}}</http-rewrite-labels-label>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.proxyHost.length > 0">Host: {{rule.proxyHost}}</http-rewrite-labels-label>
				</td>
				<td>{{rule.replace}}</td>
				<td>
					<span v-if="rule.mode == 'proxy'">隐式</span>
					<span v-if="rule.mode == 'redirect'">显示</span>
				</td>
				<td>
					<label-on :v-is-on="rule.isOn"></label-on>
				</td>
				<td>
					<a href="" @click.prevent="updateRewriteRule(rule.id)">修改</a> &nbsp;
					<a href="" @click.prevent="deleteRewriteRule(rule.id)">删除</a>
				</td>
			</tr>
		</tbody>
	</table>
	<p class="comment" v-if="rewriteRules.length > 0">拖动左侧的<i class="icon bars grey"></i>图标可以对重写规则进行排序。</p>

</div>`
})

Vue.component("http-rewrite-labels-label", {
	props: ["v-class"],
	template: `<span class="ui label tiny" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></span>`
})

Vue.component("server-name-box", {
	props: ["v-server-names"],
	data: function () {
		let serverNames = this.vServerNames;
		if (serverNames == null) {
			serverNames = []
		}
		return {
			serverNames: serverNames,
			isSearching: false,
			keyword: ""
		}
	},
	methods: {
		addServerName: function () {
			window.UPDATING_SERVER_NAME = null
			let that = this
			teaweb.popup("/servers/addServerNamePopup", {
				callback: function (resp) {
					var serverName = resp.data.serverName
					that.serverNames.push(serverName)
					setTimeout(that.submitForm, 100)
				}
			});
		},

		removeServerName: function (index) {
			this.serverNames.$remove(index)
		},

		updateServerName: function (index, serverName) {
			window.UPDATING_SERVER_NAME = teaweb.clone(serverName)
			let that = this
			teaweb.popup("/servers/addServerNamePopup", {
				callback: function (resp) {
					var serverName = resp.data.serverName
					Vue.set(that.serverNames, index, serverName)
					setTimeout(that.submitForm, 100)
				}
			});
		},
		showSearchBox: function () {
			this.isSearching = !this.isSearching
			if (this.isSearching) {
				let that = this
				setTimeout(function () {
					that.$refs.keywordRef.focus()
				}, 200)
			} else {
				this.keyword = ""
			}
		},
		allServerNames: function () {
			if (this.serverNames == null) {
				return []
			}
			let result = []
			this.serverNames.forEach(function (serverName) {
				if (serverName.subNames != null && serverName.subNames.length > 0) {
					serverName.subNames.forEach(function (subName) {
						if (subName != null && subName.length > 0) {
							if (!result.$contains(subName)) {
								result.push(subName)
							}
						}
					})
				} else if (serverName.name != null && serverName.name.length > 0) {
					if (!result.$contains(serverName.name)) {
						result.push(serverName.name)
					}
				}
			})
			return result
		},
		submitForm: function () {
			Tea.runActionOn(this.$refs.serverNamesRef.form)
		}
	},
	watch: {
		keyword: function (v) {
			this.serverNames.forEach(function (serverName) {
				if (v.length == 0) {
					serverName.isShowing = true
					return
				}
				if (serverName.subNames == null || serverName.subNames.length == 0) {
					if (!teaweb.match(serverName.name, v)) {
						serverName.isShowing = false
					}
				} else {
					let found = false
					serverName.subNames.forEach(function (subName) {
						if (teaweb.match(subName, v)) {
							found = true
						}
					})
					serverName.isShowing = found
				}
			})
		}
	},
	template: `<div>
	<input type="hidden" name="serverNames" :value="JSON.stringify(serverNames)" ref="serverNamesRef"/>
	<div v-if="serverNames.length > 0">
		<div v-for="(serverName, index) in serverNames" class="ui label small basic" :class="{hidden: serverName.isShowing === false}">
			<em v-if="serverName.type != 'full'">{{serverName.type}}</em>  
			<span v-if="serverName.subNames == null || serverName.subNames.length == 0" :class="{disabled: serverName.isShowing === false}">{{serverName.name}}</span>
			<span v-else :class="{disabled: serverName.isShowing === false}">{{serverName.subNames[0]}}等{{serverName.subNames.length}}个域名</span>
			<a href="" title="修改" @click.prevent="updateServerName(index, serverName)"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="removeServerName(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline">
	    <div class="ui field"><a href="" @click.prevent="addServerName()">[添加域名绑定]</a></div>
	    <div class="ui field" v-if="serverNames.length > 0"><span class="grey">|</span> </div>
	    <div class="ui field" v-if="serverNames.length > 0">
	        <a href="" @click.prevent="showSearchBox()" v-if="!isSearching"><i class="icon search small"></i></a>
	        <a href="" @click.prevent="showSearchBox()" v-if="isSearching"><i class="icon close small"></i></a>
        </div>
        <div class="ui field" v-if="isSearching">
            <input type="text" placeholder="搜索域名" ref="keywordRef" class="ui input tiny" v-model="keyword"/>
        </div>
    </div>
</div>`
})

// UAM模式配置
Vue.component("uam-config-box", {
	props: ["v-uam-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vUamConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				addToWhiteList: true,
				onlyURLPatterns: [],
				exceptURLPatterns: [],
				minQPSPerIP: 0,
				keyLife: 0
			}
		}
		if (config.onlyURLPatterns == null) {
			config.onlyURLPatterns = []
		}
		if (config.exceptURLPatterns == null) {
			config.exceptURLPatterns = []
		}
		return {
			config: config,
			moreOptionsVisible: false,
			minQPSPerIP: config.minQPSPerIP,
			keyLife: config.keyLife
		}
	},
	watch: {
		minQPSPerIP: function (v) {
			let qps = parseInt(v.toString())
			if (isNaN(qps) || qps < 0) {
				qps = 0
			}
			this.config.minQPSPerIP = qps
		},
		keyLife: function (v) {
			let keyLife = parseInt(v)
			if (isNaN(keyLife) || keyLife <= 0) {
				keyLife = 0
			}
			this.config.keyLife = keyLife
		}
	},
	methods: {
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},
		changeConds: function (conds) {
			this.config.conds = conds
		}
	},
	template: `<div>
<input type="hidden" name="uamJSON" :value="JSON.stringify(config)"/>
<table class="ui table definition selectable">
	<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
	<tbody v-show="((!vIsLocation && !vIsGroup) || config.isPrior)">
		<tr>
			<td class="title">启用5秒盾</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment"><plus-label></plus-label>启用后，访问网站时，自动检查浏览器环境，阻止非正常访问。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="config.isOn">
		<tr>
			<td colspan="2"><more-options-indicator @change="showMoreOptions"></more-options-indicator></td>
		</tr>
	</tbody>
	<tbody v-show="moreOptionsVisible && config.isOn">
		<tr>
			<td>验证有效期</td>
			<td>
				<div class="ui input right labeled">
                     <input type="text" name="keyLife" v-model="keyLife" maxlength="6" size="6" style="width: 6em"/>
                     <span class="ui label">秒</span>
                </div>
                <p class="comment">单个客户端验证通过后，在这个有效期内不再重复验证；如果为0则表示系统默认。</p>
			</td>
		</tr>
		<tr>
			<td>单IP最低QPS</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" name="minQPSPerIP" maxlength="6" style="width: 6em" v-model="minQPSPerIP"/>
					<span class="ui label">请求数/秒</span>
				</div>
				<p class="comment">当某个IP在1分钟内平均QPS达到此值时，才会触发5秒盾；如果设置为0，表示任何访问都会触发。</p>
			</td>
		</tr>
		<tr>
			<td>加入IP白名单</td>
			<td>
				<checkbox v-model="config.addToWhiteList"></checkbox>
				<p class="comment">选中后，表示验证通过后，将访问者IP加入到临时白名单中，此IP下次访问时不再校验5秒盾；此白名单只对5秒盾有效，不影响其他规则。此选项主要用于可能无法正常使用Cookie的网站。</p>
			</td>
		</tr>
		<tr>
			<td>例外URL</td>
			<td>
				<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
				<p class="comment">如果填写了例外URL，表示这些URL跳过5秒盾不做处理。</p>
			</td>
		</tr>
		<tr>
			<td>限制URL</td>
			<td>
				<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
				<p class="comment">如果填写了限制URL，表示只对这些URL进行5秒盾处理；如果不填则表示支持所有的URL。</p>
			</td>
		</tr>
		<tr>
			<td>匹配条件</td>
			<td>
				<http-request-conds-box :v-conds="config.conds" @change="changeConds"></http-request-conds-box>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`
})

Vue.component("http-cache-stale-config", {
	props: ["v-cache-stale-config"],
	data: function () {
		let config = this.vCacheStaleConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				status: [],
				supportStaleIfErrorHeader: true,
				life: {
					count: 1,
					unit: "day"
				}
			}
		}
		return {
			config: config
		}
	},
	watch: {
		config: {
			deep: true,
			handler: function () {
				this.$emit("change", this.config)
			}
		}
	},
	methods: {},
	template: `<table class="ui table definition selectable">
	<tbody>
		<tr>
			<td class="title">启用过时缓存</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment"><plus-label></plus-label>选中后，在更新缓存失败后会尝试读取过时的缓存。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>有效期</td>
			<td>
				<time-duration-box :v-value="config.life"></time-duration-box>
				<p class="comment">缓存在过期之后，仍然保留的时间。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>状态码</td>
			<td><http-status-box :v-status-list="config.status"></http-status-box>
				<p class="comment">在这些状态码出现时使用过时缓存，默认支持<code-label>50x</code-label>状态码。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>支持stale-if-error</td>
			<td>
				<checkbox v-model="config.supportStaleIfErrorHeader"></checkbox>
				<p class="comment">选中后，支持在Cache-Control中通过<code-label>stale-if-error</code-label>指定过时缓存有效期。</p>
			</td>
		</tr>
	</tbody>
</table>`
})

Vue.component("firewall-syn-flood-config-viewer", {
	props: ["v-syn-flood-config"],
	data: function () {
		let config = this.vSynFloodConfig
		if (config == null) {
			config = {
				isOn: false,
				minAttempts: 10,
				timeoutSeconds: 600,
				ignoreLocal: true
			}
		}
		return {
			config: config
		}
	},
	template: `<div>
	<span v-if="config.isOn">
		已启用 / <span>空连接次数：{{config.minAttempts}}次/分钟</span> / 封禁时长：{{config.timeoutSeconds}}秒 <span v-if="config.ignoreLocal">/ 忽略局域网访问</span>
	</span>
	<span v-else>未启用</span>
</div>`
})

// 域名列表
Vue.component("domains-box", {
	props: ["v-domains", "name", "v-support-wildcard"],
	data: function () {
		let domains = this.vDomains
		if (domains == null) {
			domains = []
		}

		let realName = "domainsJSON"
		if (this.name != null && typeof this.name == "string") {
			realName = this.name
		}

		let supportWildcard = true
		if (typeof this.vSupportWildcard == "boolean") {
			supportWildcard = this.vSupportWildcard
		}

		return {
			domains: domains,

			mode: "single", // single | batch
			batchDomains: "",

			isAdding: false,
			addingDomain: "",

			isEditing: false,
			editingIndex: -1,

			realName: realName,
			supportWildcard: supportWildcard
		}
	},
	watch: {
		vSupportWildcard: function (v) {
			if (typeof v == "boolean") {
				this.supportWildcard = v
			}
		},
		mode: function (mode) {
			let that = this
			setTimeout(function () {
				if (mode == "single") {
					if (that.$refs.addingDomain != null) {
						that.$refs.addingDomain.focus()
					}
				} else if (mode == "batch") {
					if (that.$refs.batchDomains != null) {
						that.$refs.batchDomains.focus()
					}
				}
			}, 100)
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingDomain.focus()
			}, 100)
		},
		confirm: function () {
			if (this.mode == "batch") {
				this.confirmBatch()
				return
			}

			let that = this

			// 删除其中的空格
			this.addingDomain = this.addingDomain.replace(/\s/g, "")

			if (this.addingDomain.length == 0) {
				teaweb.warn("请输入要添加的域名", function () {
					that.$refs.addingDomain.focus()
				})
				return
			}

			// 基本校验
			if (this.supportWildcard) {
				if (this.addingDomain[0] == "~") {
					let expr = this.addingDomain.substring(1)
					try {
						new RegExp(expr)
					} catch (e) {
						teaweb.warn("正则表达式错误：" + e.message, function () {
							that.$refs.addingDomain.focus()
						})
						return
					}
				}
			} else {
				if (/[*~^]/.test(this.addingDomain)) {
					teaweb.warn("当前只支持添加普通域名，域名中不能含有特殊符号", function () {
						that.$refs.addingDomain.focus()
					})
					return
				}
			}

			if (this.isEditing && this.editingIndex >= 0) {
				this.domains[this.editingIndex] = this.addingDomain
			} else {
				// 分割逗号（，）、顿号（、）
				if (this.addingDomain.match("[，、,;]")) {
					let domainList = this.addingDomain.split(new RegExp("[，、,;]"))
					domainList.forEach(function (v) {
						if (v.length > 0) {
							that.domains.push(v)
						}
					})
				} else {
					this.domains.push(this.addingDomain)
				}
			}
			this.cancel()
			this.change()
		},
		confirmBatch: function () {
			let domains = this.batchDomains.split("\n")
			let realDomains = []
			let that = this
			let hasProblems = false
			domains.forEach(function (domain) {
				if (hasProblems) {
					return
				}
				if (domain.length == 0) {
					return
				}
				if (that.supportWildcard) {
					if (domain == "~") {
						let expr = domain.substring(1)
						try {
							new RegExp(expr)
						} catch (e) {
							hasProblems = true
							teaweb.warn("正则表达式错误：" + e.message, function () {
								that.$refs.batchDomains.focus()
							})
							return
						}
					}
				} else {
					if (/[*~^]/.test(domain)) {
						hasProblems = true
						teaweb.warn("当前只支持添加普通域名，域名中不能含有特殊符号", function () {
							that.$refs.batchDomains.focus()
						})
						return
					}
				}
				realDomains.push(domain)
			})
			if (hasProblems) {
				return
			}
			if (realDomains.length == 0) {
				teaweb.warn("请输入要添加的域名", function () {
					that.$refs.batchDomains.focus()
				})
				return
			}

			realDomains.forEach(function (domain) {
				that.domains.push(domain)
			})
			this.cancel()
			this.change()
		},
		edit: function (index) {
			this.addingDomain = this.domains[index]
			this.isEditing = true
			this.editingIndex = index
			let that = this
			setTimeout(function () {
				that.$refs.addingDomain.focus()
			}, 50)
		},
		remove: function (index) {
			this.domains.$remove(index)
			this.change()
		},
		cancel: function () {
			this.isAdding = false
			this.mode = "single"
			this.batchDomains = ""
			this.isEditing = false
			this.editingIndex = -1
			this.addingDomain = ""
		},
		change: function () {
			this.$emit("change", this.domains)
		}
	},
	template: `<div>
	<input type="hidden" :name="realName" :value="JSON.stringify(domains)"/>
	<div v-if="domains.length > 0">
		<span class="ui label small basic" v-for="(domain, index) in domains" :class="{blue: index == editingIndex}">
			<span v-if="domain.length > 0 && domain[0] == '~'" class="grey" style="font-style: normal">[正则]</span>
			<span v-if="domain.length > 0 && domain[0] == '.'" class="grey" style="font-style: normal">[后缀]</span>
			<span v-if="domain.length > 0 && domain[0] == '*'" class="grey" style="font-style: normal">[泛域名]</span>
			{{domain}}
			<span v-if="!isAdding && !isEditing">
				&nbsp; <a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil small"></i></a>
				&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
			</span>
			<span v-if="isAdding || isEditing">
				&nbsp; <a class="disabled"><i class="icon pencil small"></i></a>
				&nbsp; <a class="disabled"><i class="icon remove small"></i></a>
			</span>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding || isEditing">
		<div class="ui fields">
			<div class="ui field" v-if="isAdding">
				<select class="ui dropdown" v-model="mode">
					<option value="single">单个</option>
					<option value="batch">批量</option>
				</select>
			</div>
			<div class="ui field">
				<div v-show="mode == 'single'">
					<input type="text" v-model="addingDomain" @keyup.enter="confirm()" @keypress.enter.prevent="1" @keydown.esc="cancel()" ref="addingDomain" :placeholder="supportWildcard ? 'example.com、*.example.com' : 'example.com、www.example.com'" size="30" maxlength="100"/>
				</div>
				<div v-show="mode == 'batch'">
					<textarea cols="30" v-model="batchDomains" placeholder="example1.com\nexample2.com\n每行一个域名" ref="batchDomains"></textarea>
				</div>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment" v-if="supportWildcard">支持普通域名（<code-label>example.com</code-label>）、泛域名（<code-label>*.example.com</code-label>）<span v-if="vSupportWildcard == undefined">、域名后缀（以点号开头，如<code-label>.example.com</code-label>）和正则表达式（以波浪号开头，如<code-label>~.*.example.com</code-label>）</span>；如果域名后有端口，请加上端口号。</p>
		<p class="comment" v-if="!supportWildcard">只支持普通域名（<code-label>example.com</code-label>、<code-label>www.example.com</code-label>）。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("http-firewall-province-selector", {
	props: ["v-type", "v-provinces"],
	data: function () {
		let provinces = this.vProvinces
		if (provinces == null) {
			provinces = []
		}

		return {
			listType: this.vType,
			provinces: provinces
		}
	},
	methods: {
		addProvince: function () {
			let selectedProvinceIds = this.provinces.map(function (province) {
				return province.id
			})
			let that = this
			teaweb.popup("/servers/server/settings/waf/ipadmin/selectProvincesPopup?type=" + this.listType + "&selectedProvinceIds=" + selectedProvinceIds.join(","), {
				width: "50em",
				height: "26em",
				callback: function (resp) {
					that.provinces = resp.data.selectedProvinces
					that.$forceUpdate()
					that.notifyChange()
				}
			})
		},
		removeProvince: function (index) {
			this.provinces.$remove(index)
			this.notifyChange()
		},
		resetProvinces: function () {
			this.provinces = []
			this.notifyChange()
		},
		notifyChange: function () {
			this.$emit("change", {
				"provinces": this.provinces
			})
		}
	},
	template: `<div>
	<span v-if="provinces.length == 0" class="disabled">暂时没有选择<span v-if="listType =='allow'">允许</span><span v-else>封禁</span>的省份。</span>
	<div v-show="provinces.length > 0">
		<div class="ui label tiny basic" v-for="(province, index) in provinces" style="margin-bottom: 0.5em">
			<input type="hidden" :name="listType + 'ProvinceIds'" :value="province.id"/>
			{{province.name}} <a href="" @click.prevent="removeProvince(index)" title="删除"><i class="icon remove"></i></a>
		</div>
	</div>
	<div class="ui divider"></div>
	<button type="button" class="ui button tiny" @click.prevent="addProvince">修改</button> &nbsp; <button type="button" class="ui button tiny" v-show="provinces.length > 0" @click.prevent="resetProvinces">清空</button>
</div>`
})

Vue.component("http-referers-config-box", {
	props: ["v-referers-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vReferersConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				allowEmpty: true,
				allowSameDomain: true,
				allowDomains: [],
				denyDomains: [],
				checkOrigin: true
			}
		}
		if (config.allowDomains == null) {
			config.allowDomains = []
		}
		if (config.denyDomains == null) {
			config.denyDomains = []
		}
		return {
			config: config,
			moreOptionsVisible: false
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeAllowDomains: function (domains) {
			if (typeof (domains) == "object") {
				this.config.allowDomains = domains
				this.$forceUpdate()
			}
		},
		changeDenyDomains: function (domains) {
			if (typeof (domains) == "object") {
				this.config.denyDomains = domains
				this.$forceUpdate()
			}
		},
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		}
	},
	template: `<div>
<input type="hidden" name="referersJSON" :value="JSON.stringify(config)"/>
<table class="ui table selectable definition">
	<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
	<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
		<tr>
			<td class="title">启用防盗链</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="config.isOn"/>
					<label></label>
				</div>
				<p class="comment">选中后表示开启防盗链。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="isOn()">
		<tr>
			<td class="title">允许直接访问网站</td>
			<td>
				<checkbox v-model="config.allowEmpty"></checkbox>
				<p class="comment">允许用户直接访问网站，用户第一次访问网站时来源域名通常为空。</p>
			</td>
		</tr>
		<tr>
			<td>来源域名允许一致</td>
			<td>
				<checkbox v-model="config.allowSameDomain"></checkbox>
				<p class="comment">允许来源域名和当前访问的域名一致，相当于在站内访问。</p>
			</td>
		</tr>
		<tr>
			<td>允许的来源域名</td>
			<td>
				<domains-box :v-domains="config.allowDomains" @change="changeAllowDomains">></domains-box>
				<p class="comment">允许的其他来源域名列表，比如<code-label>example.com</code-label>、<code-label>*.example.com</code-label>。单个星号<code-label>*</code-label>表示允许所有域名。</p>
			</td>
		</tr>
		<tr>
			<td>禁止的来源域名</td>
			<td>
				<domains-box :v-domains="config.denyDomains" @change="changeDenyDomains"></domains-box>
				<p class="comment">禁止的来源域名列表，比如<code-label>example.org</code-label>、<code-label>*.example.org</code-label>；除了这些禁止的来源域名外，其他域名都会被允许，除非限定了允许的来源域名。</p>
			</td>
		</tr>
		<tr>
			<td colspan="2"><more-options-indicator @change="showMoreOptions"></more-options-indicator></td>
		</tr>
	</tbody>
	<tbody v-show="moreOptionsVisible && isOn()">
		<tr>
			<td>同时检查Origin</td>
			<td>
				<checkbox v-model="config.checkOrigin"></checkbox>
				<p class="comment">如果请求没有指定Referer Header，则尝试检查Origin Header，多用于跨站调用。</p>
			</td>
		</tr>
		<tr>
			<td>例外URL</td>
			<td>
				<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
				<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
			</td>
		</tr>
		<tr>
			<td>限制URL</td>
			<td>
				<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
				<p class="comment">如果填写了限制URL，表示只对这些URL进行处理；如果不填则表示支持所有的URL。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="ui margin"></div>
</div>`
})

Vue.component("server-traffic-limit-status-viewer", {
	props: ["value"],
	data: function () {
		let targetTypeName = "流量"
		if (this.value != null) {
			targetTypeName = this.targetTypeToName(this.value.targetType)
		}

		return {
			status: this.value,
			targetTypeName: targetTypeName
		}
	},
	methods: {
		targetTypeToName: function (targetType) {
			switch (targetType) {
				case "traffic":
					return "流量"
				case "request":
					return "请求数"
				case "websocketConnections":
					return "Websocket连接数"
			}
			return "流量"
		}
	},
	template: `<span v-if="status != null">
	<span v-if="status.dateType == 'day'" class="small red">已达到<span v-if="status.planId > 0">套餐</span>当日{{targetTypeName}}限制</span>
	<span v-if="status.dateType == 'month'" class="small red">已达到<span v-if="status.planId > 0">套餐</span>当月{{targetTypeName}}限制</span>
	<span v-if="status.dateType == 'total'" class="small red">已达到<span v-if="status.planId > 0">套餐</span>总体{{targetTypeName}}限制</span>
</span>`
})

Vue.component("http-redirect-to-https-box", {
	props: ["v-redirect-to-https-config", "v-is-location"],
	data: function () {
		let redirectToHttpsConfig = this.vRedirectToHttpsConfig
		if (redirectToHttpsConfig == null) {
			redirectToHttpsConfig = {
				isPrior: false,
				isOn: false,
				host: "",
				port: 0,
				status: 0,
				onlyDomains: [],
				exceptDomains: []
			}
		} else {
			if (redirectToHttpsConfig.onlyDomains == null) {
				redirectToHttpsConfig.onlyDomains = []
			}
			if (redirectToHttpsConfig.exceptDomains == null) {
				redirectToHttpsConfig.exceptDomains = []
			}
		}
		return {
			redirectToHttpsConfig: redirectToHttpsConfig,
			portString: (redirectToHttpsConfig.port > 0) ? redirectToHttpsConfig.port.toString() : "",
			moreOptionsVisible: false,
			statusOptions: [
				{"code": 301, "text": "Moved Permanently"},
				{"code": 308, "text": "Permanent Redirect"},
				{"code": 302, "text": "Found"},
				{"code": 303, "text": "See Other"},
				{"code": 307, "text": "Temporary Redirect"}
			]
		}
	},
	watch: {
		"redirectToHttpsConfig.status": function () {
			this.redirectToHttpsConfig.status = parseInt(this.redirectToHttpsConfig.status)
		},
		portString: function (v) {
			let port = parseInt(v)
			if (!isNaN(port)) {
				this.redirectToHttpsConfig.port = port
			} else {
				this.redirectToHttpsConfig.port = 0
			}
		}
	},
	methods: {
		changeMoreOptions: function (isVisible) {
			this.moreOptionsVisible = isVisible
		},
		changeOnlyDomains: function (values) {
			this.redirectToHttpsConfig.onlyDomains = values
			this.$forceUpdate()
		},
		changeExceptDomains: function (values) {
			this.redirectToHttpsConfig.exceptDomains = values
			this.$forceUpdate()
		}
	},
	template: `<div>
	<input type="hidden" name="redirectToHTTPSJSON" :value="JSON.stringify(redirectToHttpsConfig)"/>
	
	<!-- Location -->
	<table class="ui table selectable definition" v-if="vIsLocation">
		<prior-checkbox :v-config="redirectToHttpsConfig"></prior-checkbox>
		<tbody v-show="redirectToHttpsConfig.isPrior">
			<tr>
				<td class="title">自动跳转到HTTPS</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="redirectToHttpsConfig.isOn"/>
						<label></label>
					</div>
					<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上，<more-options-angle @change="changeMoreOptions"></more-options-angle></p>
					
					<!--  TODO 如果已经设置了特殊设置，需要在界面上显示 -->
					<table class="ui table" v-show="moreOptionsVisible">
						<tr>
							<td class="title">状态码</td>
							<td>
								<select class="ui dropdown auto-width" v-model="redirectToHttpsConfig.status">
									<option value="0">[使用默认]</option>
									<option v-for="option in statusOptions" :value="option.code">{{option.code}} {{option.text}}</option>
								</select>
							</td>
						</tr>
						<tr>
							<td>域名或IP地址</td>
							<td>
								<input type="text" name="host" v-model="redirectToHttpsConfig.host"/>
								<p class="comment">默认和用户正在访问的域名或IP地址一致。</p>
							</td>
						</tr>
						<tr>
							<td>端口</td>
							<td>
								<input type="text" name="port" v-model="portString" maxlength="5" style="width:6em"/>
								<p class="comment">默认端口为443。</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>	
		</tbody>
	</table>
	
	<!-- 非Location -->
	<div v-if="!vIsLocation">
		<div class="ui checkbox">
			<input type="checkbox" v-model="redirectToHttpsConfig.isOn"/>
			<label></label>
		</div>
		<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上，<more-options-angle @change="changeMoreOptions"></more-options-angle></p>
		
		<!--  TODO 如果已经设置了特殊设置，需要在界面上显示 -->
		<table class="ui table" v-show="moreOptionsVisible">
			<tr>
				<td class="title">状态码</td>
				<td>
					<select class="ui dropdown auto-width" v-model="redirectToHttpsConfig.status">
						<option value="0">[使用默认]</option>
						<option v-for="option in statusOptions" :value="option.code">{{option.code}} {{option.text}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>跳转后域名或IP地址</td>
				<td>
					<input type="text" name="host" v-model="redirectToHttpsConfig.host"/>
					<p class="comment">默认和用户正在访问的域名或IP地址一致，不填写就表示使用当前的域名。</p>
				</td>
			</tr>
			<tr>
				<td>端口</td>
				<td>
					<input type="text" name="port" v-model="portString" maxlength="5" style="width:6em"/>
					<p class="comment">默认端口为443。</p>
				</td>
			</tr>
			<tr>
				<td>允许的域名</td>
				<td>
					<domains-box :v-domains="redirectToHttpsConfig.onlyDomains" @change="changeOnlyDomains"></domains-box>
					<p class="comment">如果填写了允许的域名，那么只有这些域名可以自动跳转。</p>
				</td>
			</tr>
			<tr>
				<td>排除的域名</td>
				<td>
					<domains-box :v-domains="redirectToHttpsConfig.exceptDomains" @change="changeExceptDomains"></domains-box>
					<p class="comment">如果填写了排除的域名，那么这些域名将不跳转。</p>
				</td>
			</tr>
		</table>
	</div>
	<div class="margin"></div>
</div>`
})

// 动作选择
Vue.component("http-firewall-actions-box", {
	props: ["v-actions", "v-firewall-policy", "v-action-configs", "v-group-type"],
	mounted: function () {
		let that = this
		Tea.action("/servers/iplists/levelOptions")
			.success(function (resp) {
				that.ipListLevels = resp.data.levels
			})
			.post()

		this.loadJS(function () {
			let box = document.getElementById("actions-box")
			Sortable.create(box, {
				draggable: ".label",
				handle: ".icon.handle",
				onStart: function () {
					that.cancel()
				},
				onUpdate: function (event) {
					let labels = box.getElementsByClassName("label")
					let newConfigs = []
					for (let i = 0; i < labels.length; i++) {
						let index = parseInt(labels[i].getAttribute("data-index"))
						newConfigs.push(that.configs[index])
					}
					that.configs = newConfigs
				}
			})
		})
	},
	data: function () {
		if (this.vFirewallPolicy.inbound == null) {
			this.vFirewallPolicy.inbound = {}
		}
		if (this.vFirewallPolicy.inbound.groups == null) {
			this.vFirewallPolicy.inbound.groups = []
		}

		if (this.vFirewallPolicy.outbound == null) {
			this.vFirewallPolicy.outbound = {}
		}
		if (this.vFirewallPolicy.outbound.groups == null) {
			this.vFirewallPolicy.outbound.groups = []
		}

		let id = 0
		let configs = []
		if (this.vActionConfigs != null) {
			configs = this.vActionConfigs
			configs.forEach(function (v) {
				v.id = (id++)
			})
		}

		var defaultPageBody = `<!DOCTYPE html>
<html lang="en">
<head>
\t<title>403 Forbidden</title>
\t<style>
\t\taddress { line-height: 1.8; }
\t</style>
</head>
<body>
<h1>403 Forbidden</h1>
<address>Connection: \${remoteAddr} (Client) -&gt; \${serverAddr} (Server)</address>
<address>Request ID: \${requestId}</address>
</body>
</html>`


		return {
			id: id,

			actions: this.vActions,
			configs: configs,
			isAdding: false,
			editingIndex: -1,

			action: null,
			actionCode: "",
			actionOptions: {},

			// IPList相关
			ipListLevels: [],

			// 动作参数
			allowScope: "global",

			blockTimeout: "",
			blockTimeoutMax: "",
			blockScope: "global",

			captchaLife: "",
			captchaMaxFails: "",
			captchaFailBlockTimeout: "",

			get302Life: "",

			post307Life: "",

			recordIPType: "black",
			recordIPLevel: "critical",
			recordIPTimeout: "",
			recordIPListId: 0,
			recordIPListName: "",

			tagTags: [],

			pageUseDefault: true,
			pageStatus: 403,
			pageBody: defaultPageBody,
			defaultPageBody: defaultPageBody,

			redirectStatus: 307,
			redirectURL: "",

			goGroupName: "",
			goGroupId: 0,
			goGroup: null,

			goSetId: 0,
			goSetName: "",

			jsCookieLife: "",
			jsCookieMaxFails: "",
			jsCookieFailBlockTimeout: "",

			statusOptions: [
				{"code": 301, "text": "Moved Permanently"},
				{"code": 308, "text": "Permanent Redirect"},
				{"code": 302, "text": "Found"},
				{"code": 303, "text": "See Other"},
				{"code": 307, "text": "Temporary Redirect"}
			]
		}
	},
	watch: {
		actionCode: function (code) {
			this.action = this.actions.$find(function (k, v) {
				return v.code == code
			})
			this.actionOptions = {}
		},
		allowScope: function (v) {
			this.actionOptions["scope"] = v
		},
		blockTimeout: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["timeout"] = 0
			} else {
				this.actionOptions["timeout"] = v
			}
		},
		blockTimeoutMax: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["timeoutMax"] = 0
			} else {
				this.actionOptions["timeoutMax"] = v
			}
		},
		blockScope: function (v) {
			this.actionOptions["scope"] = v
		},
		captchaLife: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["life"] = 0
			} else {
				this.actionOptions["life"] = v
			}
		},
		captchaMaxFails: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["maxFails"] = 0
			} else {
				this.actionOptions["maxFails"] = v
			}
		},
		captchaFailBlockTimeout: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["failBlockTimeout"] = 0
			} else {
				this.actionOptions["failBlockTimeout"] = v
			}
		},
		get302Life: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["life"] = 0
			} else {
				this.actionOptions["life"] = v
			}
		},
		post307Life: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["life"] = 0
			} else {
				this.actionOptions["life"] = v
			}
		},
		recordIPType: function (v) {
			this.recordIPListId = 0
		},
		recordIPTimeout: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["timeout"] = 0
			} else {
				this.actionOptions["timeout"] = v
			}
		},
		goGroupId: function (groupId) {
			let group = this.vFirewallPolicy.inbound.groups.$find(function (k, v) {
				return v.id == groupId
			})
			this.goGroup = group
			if (group == null) {
				// search outbound groups
				group = this.vFirewallPolicy.outbound.groups.$find(function (k, v) {
					return v.id == groupId
				})
				if (group == null) {
					this.goGroupName = ""
				} else {
					this.goGroup = group
					this.goGroupName = group.name
				}
			} else {
				this.goGroupName = group.name
			}
			this.goSetId = 0
			this.goSetName = ""
		},
		goSetId: function (setId) {
			if (this.goGroup == null) {
				return
			}
			let set = this.goGroup.sets.$find(function (k, v) {
				return v.id == setId
			})
			if (set == null) {
				this.goSetId = 0
				this.goSetName = ""
			} else {
				this.goSetName = set.name
			}
		},
		jsCookieLife: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["life"] = 0
			} else {
				this.actionOptions["life"] = v
			}
		},
		jsCookieMaxFails: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["maxFails"] = 0
			} else {
				this.actionOptions["maxFails"] = v
			}
		},
		jsCookieFailBlockTimeout: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["failBlockTimeout"] = 0
			} else {
				this.actionOptions["failBlockTimeout"] = v
			}
		},
	},
	methods: {
		add: function () {
			this.action = null
			this.actionCode = "page"
			this.isAdding = true
			this.actionOptions = {}

			// 动作参数
			this.allowScope = "global"

			this.blockTimeout = ""
			this.blockTimeoutMax = ""
			this.blockScope = "global"

			this.captchaLife = ""
			this.captchaMaxFails = ""
			this.captchaFailBlockTimeout = ""

			this.jsCookieLife = ""
			this.jsCookieMaxFails = ""
			this.jsCookieFailBlockTimeout = ""

			this.get302Life = ""

			this.post307Life = ""

			this.recordIPLevel = "critical"
			this.recordIPType = "black"
			this.recordIPTimeout = ""
			this.recordIPListId = 0
			this.recordIPListName = ""

			this.tagTags = []

			this.pageUseDefault = true
			this.pageStatus = 403
			this.pageBody = this.defaultPageBody

			this.redirectStatus = 307
			this.redirectURL = ""

			this.goGroupName = ""
			this.goGroupId = 0
			this.goGroup = null

			this.goSetId = 0
			this.goSetName = ""

			let that = this
			this.action = this.vActions.$find(function (k, v) {
				return v.code == that.actionCode
			})

			// 滚到界面底部
			this.scroll()
		},
		remove: function (index) {
			this.isAdding = false
			this.editingIndex = -1
			this.configs.$remove(index)
		},
		update: function (index, config) {
			if (this.isAdding && this.editingIndex == index) {
				this.cancel()
				return
			}

			this.add()

			this.isAdding = true
			this.editingIndex = index

			this.actionCode = config.code
			this.action = this.actions.$find(function (k, v) {
				return v.code == config.code
			})

			switch (config.code) {
				case "block":
					this.blockTimeout = ""
					this.blockTimeoutMax = ""
					if (config.options.timeout != null || config.options.timeout > 0) {
						this.blockTimeout = config.options.timeout.toString()
					}
					if (config.options.timeoutMax != null || config.options.timeoutMax > 0) {
						this.blockTimeoutMax = config.options.timeoutMax.toString()
					}
					if (config.options.scope != null && config.options.scope.length > 0) {
						this.blockScope = config.options.scope
					} else {
						this.blockScope = "global" // 兼容先前版本遗留的默认值
					}
					break
				case "allow":
					if (config.options != null && config.options.scope != null && config.options.scope.length > 0) {
						this.allowScope = config.options.scope
					} else {
						this.allowScope = "global"
					}
					break
				case "log":
					break
				case "captcha":
					this.captchaLife = ""
					if (config.options.life != null || config.options.life > 0) {
						this.captchaLife = config.options.life.toString()
					}
					this.captchaMaxFails = ""
					if (config.options.maxFails != null || config.options.maxFails > 0) {
						this.captchaMaxFails = config.options.maxFails.toString()
					}
					this.captchaFailBlockTimeout = ""
					if (config.options.failBlockTimeout != null || config.options.failBlockTimeout > 0) {
						this.captchaFailBlockTimeout = config.options.failBlockTimeout.toString()
					}
					break
				case "js_cookie":
					this.jsCookieLife = ""
					if (config.options.life != null || config.options.life > 0) {
						this.jsCookieLife = config.options.life.toString()
					}
					this.jsCookieMaxFails = ""
					if (config.options.maxFails != null || config.options.maxFails > 0) {
						this.jsCookieMaxFails = config.options.maxFails.toString()
					}
					this.jsCookieFailBlockTimeout = ""
					if (config.options.failBlockTimeout != null || config.options.failBlockTimeout > 0) {
						this.jsCookieFailBlockTimeout = config.options.failBlockTimeout.toString()
					}
					break
				case "notify":
					break
				case "get_302":
					this.get302Life = ""
					if (config.options.life != null || config.options.life > 0) {
						this.get302Life = config.options.life.toString()
					}
					break
				case "post_307":
					this.post307Life = ""
					if (config.options.life != null || config.options.life > 0) {
						this.post307Life = config.options.life.toString()
					}
					break;
				case "record_ip":
					if (config.options != null) {
						this.recordIPLevel = config.options.level
						this.recordIPType = config.options.type
						if (config.options.timeout > 0) {
							this.recordIPTimeout = config.options.timeout.toString()
						}
						let that = this

						// VUE需要在函数执行完之后才会调用watch函数，这样会导致设置的值被覆盖，所以这里使用setTimeout
						setTimeout(function () {
							that.recordIPListId = config.options.ipListId
							that.recordIPListName = config.options.ipListName
						})
					}
					break
				case "tag":
					this.tagTags = []
					if (config.options.tags != null) {
						this.tagTags = config.options.tags
					}
					break
				case "page":
					this.pageUseDefault = true
					this.pageStatus = 403
					this.pageBody = this.defaultPageBody
					if (typeof config.options.useDefault === "boolean") {
						this.pageUseDefault = config.options.useDefault
					} else {
						this.pageUseDefault = false
					}
					if (config.options.status != null) {
						this.pageStatus = config.options.status
					}
					if (config.options.body != null) {
						this.pageBody = config.options.body
					}
					break
				case "redirect":
					this.redirectStatus = 307
					this.redirectURL = ""
					if (config.options.status != null) {
						this.redirectStatus = config.options.status
					}
					if (config.options.url != null) {
						this.redirectURL = config.options.url
					}
					break
				case "go_group":
					if (config.options != null) {
						this.goGroupName = config.options.groupName
						this.goGroupId = config.options.groupId
						this.goGroup = this.vFirewallPolicy.inbound.groups.$find(function (k, v) {
							return v.id == config.options.groupId
						})
					}
					break
				case "go_set":
					if (config.options != null) {
						this.goGroupName = config.options.groupName
						this.goGroupId = config.options.groupId
						this.goGroup = this.vFirewallPolicy.inbound.groups.$find(function (k, v) {
							return v.id == config.options.groupId
						})

						// VUE需要在函数执行完之后才会调用watch函数，这样会导致设置的值被覆盖，所以这里使用setTimeout
						let that = this
						setTimeout(function () {
							that.goSetId = config.options.setId
							if (that.goGroup != null) {
								let set = that.goGroup.sets.$find(function (k, v) {
									return v.id == config.options.setId
								})
								if (set != null) {
									that.goSetName = set.name
								}
							}
						})
					}
					break
			}

			// 滚到界面底部
			this.scroll()
		},
		cancel: function () {
			this.isAdding = false
			this.editingIndex = -1
		},
		confirm: function () {
			if (this.action == null) {
				return
			}

			if (this.actionOptions == null) {
				this.actionOptions = {}
			}

			// record_ip
			if (this.actionCode == "record_ip") {
				let timeout = parseInt(this.recordIPTimeout)
				if (isNaN(timeout)) {
					timeout = 0
				}
				if (this.recordIPListId < 0) {
					return
				}

				// recordIPListId can be 0

				this.actionOptions = {
					type: this.recordIPType,
					level: this.recordIPLevel,
					timeout: timeout,
					ipListId: this.recordIPListId,
					ipListName: this.recordIPListName
				}
			} else if (this.actionCode == "tag") { // tag
				if (this.tagTags == null || this.tagTags.length == 0) {
					return
				}
				this.actionOptions = {
					tags: this.tagTags
				}
			} else if (this.actionCode == "page") {
				let pageStatus = this.pageStatus.toString()
				if (!pageStatus.match(/^\d{3}$/)) {
					pageStatus = 403
				} else {
					pageStatus = parseInt(pageStatus)
				}

				this.actionOptions = {
					useDefault: this.pageUseDefault,
					status: pageStatus,
					body: this.pageBody
				}
			} else if (this.actionCode == "redirect") {
				let redirectStatus = this.redirectStatus.toString()
				if (!redirectStatus.match(/^\d{3}$/)) {
					redirectStatus = 307
				} else {
					redirectStatus = parseInt(redirectStatus)
				}

				if (this.redirectURL.length == 0) {
					teaweb.warn("请输入跳转到URL")
					return
				}

				this.actionOptions = {
					status: redirectStatus,
					url: this.redirectURL
				}
			} else if (this.actionCode == "go_group") { // go_group
				let groupId = this.goGroupId
				if (typeof (groupId) == "string") {
					groupId = parseInt(groupId)
					if (isNaN(groupId)) {
						groupId = 0
					}
				}
				if (groupId <= 0) {
					return
				}
				this.actionOptions = {
					groupId: groupId.toString(),
					groupName: this.goGroupName
				}
			} else if (this.actionCode == "go_set") { // go_set
				let groupId = this.goGroupId
				if (typeof (groupId) == "string") {
					groupId = parseInt(groupId)
					if (isNaN(groupId)) {
						groupId = 0
					}
				}

				let setId = this.goSetId
				if (typeof (setId) == "string") {
					setId = parseInt(setId)
					if (isNaN(setId)) {
						setId = 0
					}
				}
				if (setId <= 0) {
					return
				}
				this.actionOptions = {
					groupId: groupId.toString(),
					groupName: this.goGroupName,
					setId: setId.toString(),
					setName: this.goSetName
				}
			}

			let options = {}
			for (let k in this.actionOptions) {
				if (this.actionOptions.hasOwnProperty(k)) {
					options[k] = this.actionOptions[k]
				}
			}
			if (this.editingIndex > -1) {
				this.configs[this.editingIndex] = {
					id: this.configs[this.editingIndex].id,
					code: this.actionCode,
					name: this.action.name,
					options: options
				}
			} else {
				this.configs.push({
					id: (this.id++),
					code: this.actionCode,
					name: this.action.name,
					options: options
				})
			}

			this.cancel()
		},
		removeRecordIPList: function () {
			this.recordIPListId = 0
		},
		selectRecordIPList: function () {
			let that = this
			teaweb.popup("/servers/iplists/selectPopup?type=" + this.recordIPType, {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					that.recordIPListId = resp.data.list.id
					that.recordIPListName = resp.data.list.name
				}
			})
		},
		changeTags: function (tags) {
			this.tagTags = tags
		},
		loadJS: function (callback) {
			if (typeof Sortable != "undefined") {
				callback()
				return
			}

			// 引入js
			let jsFile = document.createElement("script")
			jsFile.setAttribute("src", "/js/sortable.min.js")
			jsFile.addEventListener("load", function () {
				callback()
			})
			document.head.appendChild(jsFile)
		},
		scroll: function () {
			setTimeout(function () {
				let mainDiv = document.getElementsByClassName("main")
				if (mainDiv.length > 0) {
					mainDiv[0].scrollTo(0, 1000)
				}
			}, 10)
		}
	},
	template: `<div>
	<input type="hidden" name="actionsJSON" :value="JSON.stringify(configs)"/>
	<div v-show="configs.length > 0" style="margin-bottom: 0.5em" id="actions-box"> 
		<div v-for="(config, index) in configs" :data-index="index" :key="config.id" class="ui label small basic" :class="{blue: index == editingIndex}" style="margin-bottom: 0.4em">
			{{config.name}} <span class="small">({{config.code.toUpperCase()}})</span> 
			
			<!-- allow -->
			<span class="small" v-if="config.code == 'allow' && config.options != null && config.options.scope != null && config.options.scope.length > 0">
				<span v-if="config.options.scope == 'group'">[分组]</span>
				<span v-if="config.options.scope == 'server'">[网站]</span>
				<span v-if="config.options.scope == 'global'">[网站和策略]</span>
			</span>
			
			<!-- block -->
			<span v-if="config.code == 'block' && config.options.timeout > 0">：封禁时长{{config.options.timeout}}<span v-if="config.options.timeoutMax > config.options.timeout">-{{config.options.timeoutMax}}</span>秒</span>
			
			<!-- captcha -->
			<span v-if="config.code == 'captcha' && config.options.life > 0">：有效期{{config.options.life}}秒
				<span v-if="config.options.maxFails > 0"> / 最多失败{{config.options.maxFails}}次</span>
			</span>
			
			<!-- js cookie -->
			<span v-if="config.code == 'js_cookie' && config.options.life > 0">：有效期{{config.options.life}}秒
				<span v-if="config.options.maxFails > 0"> / 最多失败{{config.options.maxFails}}次</span>
			</span>
			
			<!-- get 302 -->
			<span v-if="config.code == 'get_302' && config.options.life > 0">：有效期{{config.options.life}}秒</span>
			
			<!-- post 307 -->
			<span v-if="config.code == 'post_307' && config.options.life > 0">：有效期{{config.options.life}}秒</span>
			
			<!-- record_ip -->
			<span v-if="config.code == 'record_ip'">：<span :class="{red: config.options.ipListIsDeleted}">{{config.options.ipListName}}</span></span>
			
			<!-- tag -->
			<span v-if="config.code == 'tag'">：{{config.options.tags.join(", ")}}</span>
			
			<!-- page -->
			<span v-if="config.code == 'page'">：[{{config.options.status}}]<span v-if="config.options.useDefault">&nbsp; [默认页面]</span></span>
			
			<!-- redirect -->
			<span v-if="config.code == 'redirect'">：{{config.options.url}}</span>
			
			<!-- go_group -->
			<span v-if="config.code == 'go_group'">：{{config.options.groupName}}</span>
			
			<!-- go_set -->
			<span v-if="config.code == 'go_set'">：{{config.options.groupName}} / {{config.options.setName}}</span>
			
			<!-- 范围 -->
			<span v-if="config.code != 'allow' && config.options.scope != null && config.options.scope.length > 0" class="small grey">
				&nbsp; 
				<span v-if="config.options.scope == 'global'">[所有网站]</span>
				<span v-if="config.options.scope == 'service'">[当前网站]</span>
			</span>
			
			<!-- 操作按钮 -->
			 &nbsp; <a href="" title="修改" @click.prevent="update(index, config)"><i class="icon pencil small"></i></a> &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a> &nbsp; <a href="" title="拖动改变顺序"><i class="icon bars handle"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div style="margin-bottom: 0.5em" v-if="isAdding">
		<table class="ui table" :class="{blue: editingIndex > -1}">
			<tr>
				<td class="title">动作类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="actionCode">
						<option v-for="action in actions" :value="action.code">{{action.name}} ({{action.code.toUpperCase()}})</option>
					</select>
					<p class="comment" v-if="action != null && action.description.length > 0">{{action.description}}</p>
				</td>
			</tr>
			
			<!-- allow -->
			<tr v-if="actionCode == 'allow'">
				<td>有效范围</td>
				<td>
					<select class="ui dropdown auto-width" v-model="allowScope">
						<option value="group">分组</option>
						<option value="server">网站</option>
						<option value="global">网站和策略</option>
					</select>
					<p class="comment" v-if="allowScope == 'group'">跳过当前分组其他规则集，继续执行其他分组的规则集。</p>
					<p class="comment" v-if="allowScope == 'server'">跳过当前网站所有的规则集。</p>
					<p class="comment" v-if="allowScope =='global'">跳过当前网站和网站对应WAF策略所有的规则集。</p>
				</td>
			</tr>
			
			<!-- block -->
			<tr v-if="actionCode == 'block'">
				<td>封禁范围</td>
				<td>
					<select class="ui dropdown auto-width" v-model="blockScope">
						<option value="service">当前网站</option>
						<option value="global">所有网站</option>
					</select>
					<p class="comment" v-if="blockScope == 'service'">只封禁用户对当前网站的访问，其他网站不受影响。</p>
					<p class="comment" v-if="blockScope =='global'">封禁用户对所有网站的访问。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'block'">
				<td>封禁时长</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="blockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
				</td>
			</tr>
			<tr v-if="actionCode == 'block'">
				<td>最大封禁时长</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="blockTimeoutMax" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">选填项。如果同时填写了封禁时长和最大封禁时长，则会在两者之间随机选择一个数字作为最终的封禁时长。</p>
				</td>
			</tr>
			
			<!-- captcha -->
			<tr v-if="actionCode == 'captcha'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="captchaLife" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'captcha'">
				<td>最多失败次数</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="captchaMaxFails" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">次</span>
					</div>
					<p class="comment"><span v-if="captchaMaxFails > 0 && captchaMaxFails < 5" class="red">建议填入一个不小于5的数字，以减少误判几率。</span>允许用户失败尝试的最多次数，超过这个次数将被自动加入黑名单；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'captcha'">
				<td>失败拦截时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="captchaFailBlockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">在达到最多失败次数（大于0）时，自动拦截的时间；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			
			<!-- js cookie -->
			<tr v-if="actionCode == 'js_cookie'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="jsCookieLife" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'js_cookie'">
				<td>最多失败次数</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="jsCookieMaxFails" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">次</span>
					</div>
					<p class="comment">允许用户失败尝试的最多次数，超过这个次数将被自动加入黑名单；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'js_cookie'">
				<td>失败拦截时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="jsCookieFailBlockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">在达到最多失败次数（大于0）时，自动拦截的时间；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			
			<!-- get_302 -->
			<tr v-if="actionCode == 'get_302'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="get302Life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证。</p>
				</td>
			</tr>
			
			<!-- post_307 -->
			<tr v-if="actionCode == 'post_307'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="post307Life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证。</p>
				</td>
			</tr>
			
			<!-- record_ip -->
			<tr v-if="actionCode == 'record_ip'">
				<td>IP名单类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="recordIPType">
					<option value="black">黑名单</option>
					<option value="white">白名单</option>
					<option value="grey">灰名单</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>选择IP名单</td>
				<td>
					<div v-if="recordIPListId > 0" class="ui label basic small">{{recordIPListName}} <a href="" @click.prevent="removeRecordIPList"><i class="icon remove small"></i></a></div>
					<button type="button" class="ui button tiny" @click.prevent="selectRecordIPList">+</button>
					<p class="comment">如不选择，则自动添加到当前策略的IP名单中。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>级别</td>
				<td>
					<select class="ui dropdown auto-width" v-model="recordIPLevel">
						<option v-for="level in ipListLevels" :value="level.code">{{level.name}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>超时时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 6em" maxlength="9" v-model="recordIPTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">0表示不超时。</p>
				</td>
			</tr>
			
			<!-- tag -->
			<tr v-if="actionCode == 'tag'">
				<td>标签 *</td>
				<td>
					<values-box @change="changeTags" :values="tagTags"></values-box>
				</td>
			</tr>
			
			<!-- page -->
			<tr v-if="actionCode == 'page'">
				<td>使用默认提示</td>
				<td>
					<checkbox v-model="pageUseDefault"></checkbox>
				</td>
			</tr>
			<tr v-if="actionCode == 'page' && !pageUseDefault">
				<td class="color-border">状态码 *</td>
				<td><input type="text" style="width: 4em" maxlength="3" v-model="pageStatus"/></td>
			</tr>
			<tr v-if="actionCode == 'page' && !pageUseDefault">
				<td class="color-border">网页内容</td>
				<td>
					<textarea v-model="pageBody"></textarea>
				</td>
			</tr>
			
			<!-- redirect -->
			<tr v-if="actionCode == 'redirect'">
				<td>状态码 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="redirectStatus">
						<option v-for="status in statusOptions" :value="status.code">{{status.code}} {{status.text}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'redirect'">
				<td>跳转到URL</td>
				<td>
					<input type="text" v-model="redirectURL"/>
				</td>
			</tr>
			
			<!-- 规则分组 -->
			<tr v-if="actionCode == 'go_group'">
				<td>下一个分组 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goGroupId">
						<option value="0">[选择分组]</option>
						<option v-if="vFirewallPolicy.inbound != null && vFirewallPolicy.inbound.groups != null" v-for="group in vFirewallPolicy.inbound.groups" :value="group.id">入站：{{group.name}}</option>
						<option v-if="vGroupType == 'outbound' && vFirewallPolicy.outbound != null && vFirewallPolicy.outbound.groups != null" v-for="group in vFirewallPolicy.outbound.groups" :value="group.id">出站：{{group.name}}</option>
					</select>
				</td>
			</tr>
			
			<!-- 规则集 -->
			<tr v-if="actionCode == 'go_set'">
				<td>下一个分组 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goGroupId">
						<option value="0">[选择分组]</option>
						<option v-if="vFirewallPolicy.inbound != null && vFirewallPolicy.inbound.groups != null" v-for="group in vFirewallPolicy.inbound.groups" :value="group.id">入站：{{group.name}}</option>
						<option v-if="vGroupType == 'outbound' && vFirewallPolicy.outbound != null && vFirewallPolicy.outbound.groups != null" v-for="group in vFirewallPolicy.outbound.groups" :value="group.id">出站：{{group.name}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'go_set' && goGroup != null">
				<td>下一个规则集 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goSetId">
						<option value="0">[选择规则集]</option>
						<option v-for="set in goGroup.sets" :value="set.id">{{set.name}}</option>
					</select>
				</td>
			</tr>
		</table>
		<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button> &nbsp;
		<a href="" @click.prevent="cancel" title="取消"><i class="icon remove small"></i></a>
	</div>
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
	<p class="comment">系统总是会先执行记录日志、标签等不会修改请求的动作，再执行阻止、验证码等可能改变请求的动作。</p>
</div>`
})

// 认证设置
Vue.component("http-auth-config-box", {
	props: ["v-auth-config", "v-is-location"],
	data: function () {
		let authConfig = this.vAuthConfig
		if (authConfig == null) {
			authConfig = {
				isPrior: false,
				isOn: false
			}
		}
		if (authConfig.policyRefs == null) {
			authConfig.policyRefs = []
		}
		return {
			authConfig: authConfig
		}
	},
	methods: {
		isOn: function () {
			return (!this.vIsLocation || this.authConfig.isPrior) && this.authConfig.isOn
		},
		add: function () {
			let that = this
			teaweb.popup("/servers/server/settings/access/createPopup", {
				callback: function (resp) {
					that.authConfig.policyRefs.push(resp.data.policyRef)
					that.change()
				},
				height: "28em"
			})
		},
		update: function (index, policyId) {
			let that = this
			teaweb.popup("/servers/server/settings/access/updatePopup?policyId=" + policyId, {
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						teaweb.reload()
					})
				},
				height: "28em"
			})
		},
		remove: function (index) {
			this.authConfig.policyRefs.$remove(index)
			this.change()
		},
		methodName: function (methodType) {
			switch (methodType) {
				case "basicAuth":
					return "BasicAuth"
				case "subRequest":
					return "子请求"
				case "typeA":
					return "URL鉴权A"
				case "typeB":
					return "URL鉴权B"
				case "typeC":
					return "URL鉴权C"
				case "typeD":
					return "URL鉴权D"
			}
			return ""
		},
		change: function () {
			let that = this
			setTimeout(function () {
				// 延时通知，是为了让表单有机会变更数据
				that.$emit("change", this.authConfig)
			}, 100)
		}
	},
	template: `<div>
<input type="hidden" name="authJSON" :value="JSON.stringify(authConfig)"/> 
<table class="ui table selectable definition">
	<prior-checkbox :v-config="authConfig" v-if="vIsLocation"></prior-checkbox>
	<tbody v-show="!vIsLocation || authConfig.isPrior">
		<tr>
			<td class="title">启用鉴权</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" v-model="authConfig.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
<!-- 鉴权方式 -->
<div v-show="isOn()">
	<h4>鉴权方式</h4>
	<table class="ui table selectable celled" v-show="authConfig.policyRefs.length > 0">
		<thead>
			<tr>
				<th class="three wide">名称</th>
				<th class="three wide">鉴权方法</th>
				<th>参数</th>
				<th class="two wide">状态</th>
				<th class="two op">操作</th>
			</tr>
		</thead>
		<tbody v-for="(ref, index) in authConfig.policyRefs" :key="ref.authPolicyId">
			<tr>
				<td>{{ref.authPolicy.name}}</td>
				<td>
					{{methodName(ref.authPolicy.type)}}
				</td>
				<td>
					<span v-if="ref.authPolicy.type == 'basicAuth'">{{ref.authPolicy.params.users.length}}个用户</span>
					<span v-if="ref.authPolicy.type == 'subRequest'">
						<span v-if="ref.authPolicy.params.method.length > 0" class="grey">[{{ref.authPolicy.params.method}}]</span>
						{{ref.authPolicy.params.url}}
					</span>
					<span v-if="ref.authPolicy.type == 'typeA'">{{ref.authPolicy.params.signParamName}}/有效期{{ref.authPolicy.params.life}}秒</span>
					<span v-if="ref.authPolicy.type == 'typeB'">有效期{{ref.authPolicy.params.life}}秒</span>
					<span v-if="ref.authPolicy.type == 'typeC'">有效期{{ref.authPolicy.params.life}}秒</span>
					<span v-if="ref.authPolicy.type == 'typeD'">{{ref.authPolicy.params.signParamName}}/{{ref.authPolicy.params.timestampParamName}}/有效期{{ref.authPolicy.params.life}}秒</span>
					
					<div v-if="(ref.authPolicy.params.exts != null && ref.authPolicy.params.exts.length > 0) || (ref.authPolicy.params.domains != null && ref.authPolicy.params.domains.length > 0)">
						<grey-label v-if="ref.authPolicy.params.exts != null" v-for="ext in ref.authPolicy.params.exts">扩展名：{{ext}}</grey-label>
						<grey-label v-if="ref.authPolicy.params.domains != null" v-for="domain in ref.authPolicy.params.domains">域名：{{domain}}</grey-label>
					</div>
				</td>
				<td>
					<label-on :v-is-on="ref.authPolicy.isOn"></label-on>
				</td>
				<td>
					<a href="" @click.prevent="update(index, ref.authPolicyId)">修改</a> &nbsp;
					<a href="" @click.prevent="remove(index)">删除</a>
				</td>
			</tr>
		</tbody>
	</table>
	<button class="ui button small" type="button" @click.prevent="add">+添加鉴权方式</button>
</div>
<div class="margin"></div>
</div>`
})

Vue.component("user-selector", {
	props: ["v-user-id", "data-url"],
	data: function () {
		let userId = this.vUserId
		if (userId == null) {
			userId = 0
		}

		let dataURL = this.dataUrl
		if (dataURL == null || dataURL.length == 0) {
			dataURL = "/servers/users/options"
		}

		return {
			users: [],
			userId: userId,
			dataURL: dataURL
		}
	},
	methods: {
		change: function(item) {
			if (item != null) {
				this.$emit("change", item.id)
			} else {
				this.$emit("change", 0)
			}
		},
		clear: function () {
			this.$refs.comboBox.clear()
		}
	},
	template: `<div>
	<combo-box placeholder="选择用户" :data-url="dataURL" :data-key="'users'" data-search="on" name="userId" :v-value="userId" @change="change" ref="comboBox"></combo-box>
</div>`
})

Vue.component("http-header-policy-box", {
	props: ["v-request-header-policy", "v-request-header-ref", "v-response-header-policy", "v-response-header-ref", "v-params", "v-is-location", "v-is-group", "v-has-group-request-config", "v-has-group-response-config", "v-group-setting-url"],
	data: function () {
		let type = "response"
		let hash = window.location.hash
		if (hash == "#request") {
			type = "request"
		}

		// ref
		let requestHeaderRef = this.vRequestHeaderRef
		if (requestHeaderRef == null) {
			requestHeaderRef = {
				isPrior: false,
				isOn: true,
				headerPolicyId: 0
			}
		}

		let responseHeaderRef = this.vResponseHeaderRef
		if (responseHeaderRef == null) {
			responseHeaderRef = {
				isPrior: false,
				isOn: true,
				headerPolicyId: 0
			}
		}

		// 请求相关
		let requestSettingHeaders = []
		let requestDeletingHeaders = []
		let requestNonStandardHeaders = []

		let requestPolicy = this.vRequestHeaderPolicy
		if (requestPolicy != null) {
			if (requestPolicy.setHeaders != null) {
				requestSettingHeaders = requestPolicy.setHeaders
			}
			if (requestPolicy.deleteHeaders != null) {
				requestDeletingHeaders = requestPolicy.deleteHeaders
			}
			if (requestPolicy.nonStandardHeaders != null) {
				requestNonStandardHeaders = requestPolicy.nonStandardHeaders
			}
		}

		// 响应相关
		let responseSettingHeaders = []
		let responseDeletingHeaders = []
		let responseNonStandardHeaders = []

		let responsePolicy = this.vResponseHeaderPolicy
		if (responsePolicy != null) {
			if (responsePolicy.setHeaders != null) {
				responseSettingHeaders = responsePolicy.setHeaders
			}
			if (responsePolicy.deleteHeaders != null) {
				responseDeletingHeaders = responsePolicy.deleteHeaders
			}
			if (responsePolicy.nonStandardHeaders != null) {
				responseNonStandardHeaders = responsePolicy.nonStandardHeaders
			}
		}

		let responseCORS = {
			isOn: false
		}
		if (responsePolicy.cors != null) {
			responseCORS = responsePolicy.cors
		}

		return {
			type: type,
			typeName: (type == "request") ? "请求" : "响应",

			requestHeaderRef: requestHeaderRef,
			responseHeaderRef: responseHeaderRef,
			requestSettingHeaders: requestSettingHeaders,
			requestDeletingHeaders: requestDeletingHeaders,
			requestNonStandardHeaders: requestNonStandardHeaders,

			responseSettingHeaders: responseSettingHeaders,
			responseDeletingHeaders: responseDeletingHeaders,
			responseNonStandardHeaders: responseNonStandardHeaders,
			responseCORS: responseCORS
		}
	},
	methods: {
		selectType: function (type) {
			this.type = type
			window.location.hash = "#" + type
			window.location.reload()
		},
		addSettingHeader: function (policyId) {
			teaweb.popup("/servers/server/settings/headers/createSetPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + this.type, {
				height: "22em",
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		},
		addDeletingHeader: function (policyId, type) {
			teaweb.popup("/servers/server/settings/headers/createDeletePopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + type, {
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		},
		addNonStandardHeader: function (policyId, type) {
			teaweb.popup("/servers/server/settings/headers/createNonStandardPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + type, {
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		},
		updateSettingPopup: function (policyId, headerId) {
			teaweb.popup("/servers/server/settings/headers/updateSetPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&headerId=" + headerId + "&type=" + this.type, {
				height: "22em",
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		},
		deleteDeletingHeader: function (policyId, headerName) {
			teaweb.confirm("确定要删除'" + headerName + "'吗？", function () {
				Tea.action("/servers/server/settings/headers/deleteDeletingHeader")
					.params({
						headerPolicyId: policyId,
						headerName: headerName
					})
					.post()
					.refresh()
			})
		},
		deleteNonStandardHeader: function (policyId, headerName) {
			teaweb.confirm("确定要删除'" + headerName + "'吗？", function () {
				Tea.action("/servers/server/settings/headers/deleteNonStandardHeader")
					.params({
						headerPolicyId: policyId,
						headerName: headerName
					})
					.post()
					.refresh()
			})
		},
		deleteHeader: function (policyId, type, headerId) {
			teaweb.confirm("确定要删除此报头吗？", function () {
					this.$post("/servers/server/settings/headers/delete")
						.params({
							headerPolicyId: policyId,
							type: type,
							headerId: headerId
						})
						.refresh()
				}
			)
		},
		updateCORS: function (policyId) {
			teaweb.popup("/servers/server/settings/headers/updateCORSPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + this.type, {
				height: "30em",
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		}
	},
	template: `<div>
	<div class="ui menu tabular small">
		<a class="item" :class="{active:type == 'response'}" @click.prevent="selectType('response')">响应报头<span v-if="responseSettingHeaders.length > 0">({{responseSettingHeaders.length}})</span></a>
		<a class="item" :class="{active:type == 'request'}" @click.prevent="selectType('request')">请求报头<span v-if="requestSettingHeaders.length > 0">({{requestSettingHeaders.length}})</span></a>
	</div>
	
	<div class="margin"></div>
	
	<input type="hidden" name="type" :value="type"/>
	
	<!-- 请求 -->
	<div v-if="(vIsLocation || vIsGroup) && type == 'request'">
		<input type="hidden" name="requestHeaderJSON" :value="JSON.stringify(requestHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="requestHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="((!vIsLocation && !vIsGroup) || requestHeaderRef.isPrior) && type == 'request'">
		<div v-if="vHasGroupRequestConfig">
        	<div class="margin"></div>
        	<warning-message>由于已经在当前<a :href="vGroupSettingUrl + '#request'">网站分组</a>中进行了对应的配置，在这里的配置将不会生效。</warning-message>
    	</div>
    	<div :class="{'opacity-mask': vHasGroupRequestConfig}">
		<h4>设置请求报头 &nbsp; <a href="" @click.prevent="addSettingHeader(vRequestHeaderPolicy.id)" style="font-size: 0.8em">[添加新报头]</a></h4>
			<p class="comment" v-if="requestSettingHeaders.length == 0">暂时还没有自定义报头。</p>
			<table class="ui table selectable celled" v-if="requestSettingHeaders.length > 0">
				<thead>
					<tr>
						<th>名称</th>
						<th>值</th>
						<th class="two op">操作</th>
					</tr>
				</thead>
				<tbody v-for="header in requestSettingHeaders">
					<tr>
						<td class="five wide">
							<a href="" @click.prevent="updateSettingPopup(vRequestHeaderPolicy.id, header.id)">{{header.name}} <i class="icon expand small"></i></a>
							<div>
								<span v-if="header.status != null && header.status.codes != null && !header.status.always"><grey-label v-for="code in header.status.codes" :key="code">{{code}}</grey-label></span>
								<span v-if="header.methods != null && header.methods.length > 0"><grey-label v-for="method in header.methods" :key="method">{{method}}</grey-label></span>
								<span v-if="header.domains != null && header.domains.length > 0"><grey-label v-for="domain in header.domains" :key="domain">{{domain}}</grey-label></span>
								<grey-label v-if="header.shouldAppend">附加</grey-label>
								<grey-label v-if="header.disableRedirect">跳转禁用</grey-label>
								<grey-label v-if="header.shouldReplace && header.replaceValues != null && header.replaceValues.length > 0">替换</grey-label>
							</div>
						</td>
						<td>{{header.value}}</td>
						<td><a href="" @click.prevent="updateSettingPopup(vRequestHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vRequestHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
					</tr>
				</tbody>
			</table>
			
			<h4>其他设置</h4>
			
			<table class="ui table definition selectable">
				<tbody>
					<tr>
						<td class="title">删除报头 <tip-icon content="可以通过此功能删除转发到源站的请求报文中不需要的报头"></tip-icon></td>
						<td>
							<div v-if="requestDeletingHeaders.length > 0">
								<div class="ui label small basic" v-for="headerName in requestDeletingHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteDeletingHeader(vRequestHeaderPolicy.id, headerName)"></i></a> </div>
								<div class="ui divider" ></div>
							</div>
							<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vRequestHeaderPolicy.id, 'request')">+</button>
						</td>
					</tr>
					<tr>
						<td class="title">非标报头 <tip-icon content="可以通过此功能设置转发到源站的请求报文中非标准的报头，比如hello_world"></tip-icon></td>
						<td>
							<div v-if="requestNonStandardHeaders.length > 0">
								<div class="ui label small basic" v-for="headerName in requestNonStandardHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteNonStandardHeader(vRequestHeaderPolicy.id, headerName)"></i></a> </div>
								<div class="ui divider" ></div>
							</div>
							<button class="ui button small" type="button" @click.prevent="addNonStandardHeader(vRequestHeaderPolicy.id, 'request')">+</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>			
	</div>
	
	<!-- 响应 -->
	<div v-if="(vIsLocation || vIsGroup) && type == 'response'">
		<input type="hidden" name="responseHeaderJSON" :value="JSON.stringify(responseHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="responseHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="((!vIsLocation && !vIsGroup) || responseHeaderRef.isPrior) && type == 'response'">
		<div v-if="vHasGroupResponseConfig">
        	<div class="margin"></div>
        	<warning-message>由于已经在当前<a :href="vGroupSettingUrl + '#response'">网站分组</a>中进行了对应的配置，在这里的配置将不会生效。</warning-message>
    	</div>
    	<div :class="{'opacity-mask': vHasGroupResponseConfig}">
			<h4>设置响应报头 &nbsp; <a href="" @click.prevent="addSettingHeader(vResponseHeaderPolicy.id)" style="font-size: 0.8em">[添加新报头]</a></h4>
			<p class="comment" style="margin-top: 0; padding-top: 0">将会覆盖已有的同名报头。</p>
			<p class="comment" v-if="responseSettingHeaders.length == 0">暂时还没有自定义报头。</p>
			<table class="ui table selectable celled" v-if="responseSettingHeaders.length > 0">
				<thead>
					<tr>
						<th>名称</th>
						<th>值</th>
						<th class="two op">操作</th>
					</tr>
				</thead>
				<tbody v-for="header in responseSettingHeaders">
					<tr>
						<td class="five wide">
							<a href="" @click.prevent="updateSettingPopup(vResponseHeaderPolicy.id, header.id)">{{header.name}} <i class="icon expand small"></i></a>
							<div>
								<span v-if="header.status != null && header.status.codes != null && !header.status.always"><grey-label v-for="code in header.status.codes" :key="code">{{code}}</grey-label></span>
								<span v-if="header.methods != null && header.methods.length > 0"><grey-label v-for="method in header.methods" :key="method">{{method}}</grey-label></span>
								<span v-if="header.domains != null && header.domains.length > 0"><grey-label v-for="domain in header.domains" :key="domain">{{domain}}</grey-label></span>
								<grey-label v-if="header.shouldAppend">附加</grey-label>
								<grey-label v-if="header.disableRedirect">跳转禁用</grey-label>
								<grey-label v-if="header.shouldReplace && header.replaceValues != null && header.replaceValues.length > 0">替换</grey-label>
							</div>
							
							<!-- CORS -->
							<div v-if="header.name == 'Access-Control-Allow-Origin' && header.value == '*'">
								<span class="red small">建议使用当前页面下方的"CORS自适应跨域"功能代替Access-Control-*-*相关报头。</span>
							</div>
						</td>
						<td>{{header.value}}</td>
						<td><a href="" @click.prevent="updateSettingPopup(vResponseHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vResponseHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
					</tr>
				</tbody>
			</table>
			
			<h4>其他设置</h4>
			
			<table class="ui table definition selectable">
				<tbody>
					<tr>
						<td class="title">删除报头 <tip-icon content="可以通过此功能删除响应报文中不需要的报头"></tip-icon></td>
						<td>
							<div v-if="responseDeletingHeaders.length > 0">
								<div class="ui label small basic" v-for="headerName in responseDeletingHeaders">{{headerName}} &nbsp; <a href=""><i class="icon remove small" title="删除" @click.prevent="deleteDeletingHeader(vResponseHeaderPolicy.id, headerName)"></i></a></div>
								<div class="ui divider" ></div>
							</div>
							<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vResponseHeaderPolicy.id, 'response')">+</button>
						</td>
					</tr>
					<tr>
						<td>非标报头 <tip-icon content="可以通过此功能设置响应报文中非标准的报头，比如hello_world"></tip-icon></td>
						<td>
							<div v-if="responseNonStandardHeaders.length > 0">
								<div class="ui label small basic" v-for="headerName in responseNonStandardHeaders">{{headerName}} &nbsp; <a href=""><i class="icon remove small" title="删除" @click.prevent="deleteNonStandardHeader(vResponseHeaderPolicy.id, headerName)"></i></a></div>
								<div class="ui divider" ></div>
							</div>
							<button class="ui button small" type="button" @click.prevent="addNonStandardHeader(vResponseHeaderPolicy.id, 'response')">+</button>
						</td>
					</tr>
					<tr>
						<td class="title">CORS自适应跨域</td>
						<td>
							<span v-if="responseCORS.isOn" class="green">已启用</span><span class="disabled" v-else="">未启用</span> &nbsp; <a href="" @click.prevent="updateCORS(vResponseHeaderPolicy.id)">[修改]</a>
							<p class="comment"><span v-if="!responseCORS.isOn">启用后，服务器可以</span><span v-else>服务器会</span>自动生成<code-label>Access-Control-*-*</code-label>相关的报头。</p>
						</td>
					</tr>
				</tbody>
			</table>
		</div>		
	</div>
	<div class="margin"></div>
</div>`
})

// 通用设置
Vue.component("http-common-config-box", {
	props: ["v-common-config"],
	data: function () {
		let config = this.vCommonConfig
		if (config == null) {
			config = {
				mergeSlashes: false
			}
		}
		return {
			config: config
		}
	},
	template: `<div>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">合并重复的路径分隔符</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" name="mergeSlashes" value="1" v-model="config.mergeSlashes"/>
					<label></label>
				</div>
				<p class="comment">合并URL中重复的路径分隔符为一个，比如<code-label>//hello/world</code-label>中的<code-label>//</code-label>。</p>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-cache-policy-selector", {
	props: ["v-cache-policy"],
	mounted: function () {
		let that = this
		Tea.action("/servers/components/cache/count")
			.post()
			.success(function (resp) {
				that.count = resp.data.count
			})
	},
	data: function () {
		let cachePolicy = this.vCachePolicy
		return {
			count: 0,
			cachePolicy: cachePolicy
		}
	},
	methods: {
		remove: function () {
			this.cachePolicy = null
		},
		select: function () {
			let that = this
			teaweb.popup("/servers/components/cache/selectPopup", {
				width: "42em",
				height: "26em",
				callback: function (resp) {
					that.cachePolicy = resp.data.cachePolicy
				}
			})
		},
		create: function () {
			let that = this
			teaweb.popup("/servers/components/cache/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.cachePolicy = resp.data.cachePolicy
				}
			})
		}
	},
	template: `<div>
	<div v-if="cachePolicy != null" class="ui label basic">
		<input type="hidden" name="cachePolicyId" :value="cachePolicy.id"/>
		{{cachePolicy.name}} &nbsp; <a :href="'/servers/components/cache/update?cachePolicyId=' + cachePolicy.id" target="_blank" title="修改"><i class="icon pencil small"></i></a>&nbsp; <a href="" @click.prevent="remove()" title="删除"><i class="icon remove small"></i></a>
	</div>
	<div v-if="cachePolicy == null">
		<span v-if="count > 0"><a href="" @click.prevent="select">[选择已有策略]</a> &nbsp; &nbsp; </span><a href="" @click.prevent="create">[创建新策略]</a>
	</div>
</div>`
})

Vue.component("http-pages-and-shutdown-box", {
	props: ["v-enable-global-pages", "v-pages", "v-shutdown-config", "v-is-location"],
	data: function () {
		let pages = []
		if (this.vPages != null) {
			pages = this.vPages
		}
		let shutdownConfig = {
			isPrior: false,
			isOn: false,
			bodyType: "html",
			url: "",
			body: "",
			status: 0
		}
		if (this.vShutdownConfig != null) {
			if (this.vShutdownConfig.body == null) {
				this.vShutdownConfig.body = ""
			}
			if (this.vShutdownConfig.bodyType == null) {
				this.vShutdownConfig.bodyType = "html"
			}
			shutdownConfig = this.vShutdownConfig
		}

		let shutdownStatus = ""
		if (shutdownConfig.status > 0) {
			shutdownStatus = shutdownConfig.status.toString()
		}

		return {
			pages: pages,
			shutdownConfig: shutdownConfig,
			shutdownStatus: shutdownStatus,
			enableGlobalPages: this.vEnableGlobalPages
		}
	},
	watch: {
		shutdownStatus: function (status) {
			let statusInt = parseInt(status)
			if (!isNaN(statusInt) && statusInt > 0 && statusInt < 1000) {
				this.shutdownConfig.status = statusInt
			} else {
				this.shutdownConfig.status = 0
			}
		}
	},
	methods: {
		addPage: function () {
			let that = this
			teaweb.popup("/servers/server/settings/pages/createPopup", {
				height: "30em",
				callback: function (resp) {
					that.pages.push(resp.data.page)
					that.notifyChange()
				}
			})
		},
		updatePage: function (pageIndex, pageId) {
			let that = this
			teaweb.popup("/servers/server/settings/pages/updatePopup?pageId=" + pageId, {
				height: "30em",
				callback: function (resp) {
					Vue.set(that.pages, pageIndex, resp.data.page)
					that.notifyChange()
				}
			})
		},
		removePage: function (pageIndex) {
			let that = this
			teaweb.confirm("确定要删除此自定义页面吗？", function () {
				that.pages.$remove(pageIndex)
				that.notifyChange()
			})
		},
		addShutdownHTMLTemplate: function () {
			this.shutdownConfig.body = `<!DOCTYPE html>
<html lang="en">
<head>
\t<title>升级中</title>
\t<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
\t<style>
\t\taddress { line-height: 1.8; }
\t</style>
</head>
<body>

<h1>网站升级中</h1>
<p>为了给您提供更好的服务，我们正在升级网站，请稍后重新访问。</p>

<address>Connection: \${remoteAddr} (Client) -&gt; \${serverAddr} (Server)</address>
<address>Request ID: \${requestId}</address>

</body>
</html>`
		},
		notifyChange: function () {
			let parent = this.$el.parentNode
			while (true) {
				if (parent == null) {
					break
				}
				if (parent.tagName == "FORM") {
					break
				}
				parent = parent.parentNode
			}
			if (parent != null) {
				setTimeout(function () {
					Tea.runActionOn(parent)
				}, 100)
			}
		}
	},
	template: `<div>
<input type="hidden" name="pagesJSON" :value="JSON.stringify(pages)"/>
<input type="hidden" name="shutdownJSON" :value="JSON.stringify(shutdownConfig)"/>

<h4 style="margin-bottom: 0.5em">自定义页面</h4>

<p class="comment" style="padding-top: 0; margin-top: 0">根据响应状态码返回一些自定义页面，比如404，500等错误页面。</p>

<div v-if="pages.length > 0">
	<table class="ui table selectable celled">
		<thead>
			<tr>
				<th class="two wide">响应状态码</th>
				<th>页面类型</th>
				<th class="two wide">新状态码</th>
				<th>例外URL</th>
				<th>限制URL</th>
				<th class="two op">操作</th>
			</tr>	
		</thead>
		<tr v-for="(page,index) in pages">
			<td>
				<a href="" @click.prevent="updatePage(index, page.id)">
					<span v-if="page.status != null && page.status.length == 1">{{page.status[0]}}</span>
					<span v-else>{{page.status}}</span>
					
					<i class="icon expand small"></i>
				</a>
			</td>
			<td style="word-break: break-all">
				<div v-if="page.bodyType == 'url'">
					{{page.url}}
					<div>
						<grey-label>读取URL</grey-label>
					</div>
				</div>
				<div v-if="page.bodyType == 'redirectURL'">
					{{page.url}}
					<div>
						<grey-label>跳转URL</grey-label>	
						<grey-label v-if="page.newStatus > 0">{{page.newStatus}}</grey-label>
					</div>
				</div>
				<div v-if="page.bodyType == 'html'">
					[HTML内容]
					<div>
						<grey-label v-if="page.newStatus > 0">{{page.newStatus}}</grey-label>
					</div>
				</div>
			</td>
			<td>
				<span v-if="page.newStatus > 0">{{page.newStatus}}</span>
				<span v-else class="disabled">保持</span>	
			</td>
			<td>
				<div v-if="page.exceptURLPatterns != null && page.exceptURLPatterns">
					<span v-for="urlPattern in page.exceptURLPatterns" class="ui basic label small">{{urlPattern.pattern}}</span>
				</div>
				<span v-else class="disabled">-</span>
			</td>
			<td>
				<div v-if="page.onlyURLPatterns != null && page.onlyURLPatterns">
					<span v-for="urlPattern in page.onlyURLPatterns" class="ui basic label small">{{urlPattern.pattern}}</span>
				</div>
				<span v-else class="disabled">-</span>
			</td>
			<td>
				<a href="" title="修改" @click.prevent="updatePage(index, page.id)">修改</a> &nbsp; 
				<a href="" title="删除" @click.prevent="removePage(index)">删除</a>
			</td>
		</tr>
	</table>
</div>
<div style="margin-top: 1em">
	<button class="ui button small" type="button" @click.prevent="addPage()">+添加自定义页面</button>
</div>

<h4 style="margin-top: 2em;">临时关闭页面</h4>
<p class="comment" style="margin-top: 0; padding-top: 0">开启临时关闭页面时，所有请求都会直接显示此页面。可用于临时升级网站或者禁止用户访问某个网页。</p>	
<div>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="shutdownConfig" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="!vIsLocation || shutdownConfig.isPrior">
			<tr>
				<td class="title">启用临时关闭网站</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="shutdownConfig.isOn" />
						<label></label>
					</div>
					<p class="comment">选中后，表示临时关闭当前网站，并显示自定义内容。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="(!vIsLocation || shutdownConfig.isPrior) && shutdownConfig.isOn">
			<tr>
				<td>显示内容类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="shutdownConfig.bodyType">
						<option value="html">HTML</option>
						<option value="url">读取URL</option>
						<option value="redirectURL">跳转URL</option>
					</select>
				</td>
			</tr>
			<tr v-if="shutdownConfig.bodyType == 'url'">
				<td class="title">显示页面URL *</td>
				<td>
					<input type="text" v-model="shutdownConfig.url" placeholder="类似于 https://example.com/page.html"/>
					<p class="comment">将从此URL中读取内容。</p>
				</td>
			</tr>
			<tr v-if="shutdownConfig.bodyType == 'redirectURL'">
				<td class="title">跳转到URL *</td>
				<td>
					<input type="text" v-model="shutdownConfig.url" placeholder="类似于 https://example.com/page.html"/>
					 <p class="comment">将会跳转到此URL。</p>
				</td>
			</tr>
			<tr v-show="shutdownConfig.bodyType == 'html'">
				<td>显示页面HTML *</td>
				<td>
					<textarea name="body" ref="shutdownHTMLBody" v-model="shutdownConfig.body"></textarea>
					<p class="comment"><a href="" @click.prevent="addShutdownHTMLTemplate">[使用模板]</a>。填写页面的HTML内容，支持请求变量。</p>
				</td>
			</tr>
			<tr>
				<td>状态码</td>
				<td><input type="text" size="3" maxlength="3" name="shutdownStatus" style="width:5.2em" placeholder="状态码" v-model="shutdownStatus"/></td>
			</tr>
		</tbody>
	</table>
</div>

<h4 style="margin-top: 2em;">其他设置</h4>
<table class="ui table definition selectable">
	<tr>
		<td class="title">启用系统自定义页面</td>
		<td>
			<checkbox name="enableGlobalPages" v-model="enableGlobalPages"></checkbox>
			<p class="comment">选中后，表示如果当前网站没有自定义页面，则尝试使用系统对应的自定义页面。</p>
		</td>
	</tr>
</table>
<div class="ui margin"></div>

</div>`
})

Vue.component("http-firewall-page-options", {
	props: ["v-page-options"],
	data: function () {
		var defaultPageBody = `<!DOCTYPE html>
<html lang="en">
<head>
	<title>403 Forbidden</title>
	<style>
		address { line-height: 1.8; }
	</style>
</head>
<body>
<h1>403 Forbidden By WAF</h1>
<address>Connection: \${remoteAddr} (Client) -&gt; \${serverAddr} (Server)</address>
<address>Request ID: \${requestId}</address>
</body>
</html>`

		return {
			pageOptions: this.vPageOptions,
			status: this.vPageOptions.status,
			body: this.vPageOptions.body,
			defaultPageBody: defaultPageBody,
			isEditing: false
		}
	},
	watch: {
		status: function (v) {
			if (typeof v === "string" && v.length != 3) {
				return
			}
			let statusCode = parseInt(v)
			if (isNaN(statusCode)) {
				this.pageOptions.status = 403
			} else {
				this.pageOptions.status = statusCode
			}
		},
		body: function (v) {
			this.pageOptions.body = v
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	<input type="hidden" name="pageOptionsJSON" :value="JSON.stringify(pageOptions)"/>
	<a href="" @click.prevent="edit">状态码：{{status}} / 提示内容：<span v-if="pageOptions.body != null && pageOptions.body.length > 0">[{{pageOptions.body.length}}字符]</span><span v-else class="disabled">[无]</span>
	 <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<table class="ui table" v-show="isEditing">
		<tr>
			<td class="title">状态码 *</td>
			<td><input type="text" style="width: 4em" maxlength="3" v-model="status"/></td>
		</tr>
		<tr>
			<td>网页内容</td>
			<td>
				<textarea v-model="body"></textarea>
				<p class="comment"><a href="" @click.prevent="body = defaultPageBody">[使用模板]</a> </p>
			</td>
		</tr>
	</table>
</div>	
`
})

Vue.component("http-firewall-js-cookie-options", {
	props: ["v-js-cookie-options"],
	mounted: function () {
		this.updateSummary()
	},
	data: function () {
		let options = this.vJsCookieOptions
		if (options == null) {
			options = {
				life: 0,
				maxFails: 0,
				failBlockTimeout: 0,
				failBlockScopeAll: false,
				scope: "service"
			}
		}

		return {
			options: options,
			isEditing: false,
			summary: ""
		}
	},
	watch: {
		"options.life": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.life = i
			this.updateSummary()
		},
		"options.maxFails": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.maxFails = i
			this.updateSummary()
		},
		"options.failBlockTimeout": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.failBlockTimeout = i
			this.updateSummary()
		},
		"options.failBlockScopeAll": function (v) {
			this.updateSummary()
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		},
		updateSummary: function () {
			let summaryList = []
			if (this.options.life > 0) {
				summaryList.push("有效时间" + this.options.life + "秒")
			}
			if (this.options.maxFails > 0) {
				summaryList.push("最多失败" + this.options.maxFails + "次")
			}
			if (this.options.failBlockTimeout > 0) {
				summaryList.push("失败拦截" + this.options.failBlockTimeout + "秒")
			}
			if (this.options.failBlockScopeAll) {
				summaryList.push("尝试全局封禁")
			}

			if (summaryList.length == 0) {
				this.summary = "默认配置"
			} else {
				this.summary = summaryList.join(" / ")
			}
		},
		confirm: function () {
			this.isEditing = false
		}
	},
	template: `<div>
	<input type="hidden" name="jsCookieOptionsJSON" :value="JSON.stringify(options)"/>
	<a href="" @click.prevent="edit">{{summary}} <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<div v-show="isEditing" style="margin-top: 0.5em">
		<table class="ui table definition selectable">
			<tbody>
				<tr>
					<td class="title">有效时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">验证通过后在这个时间内不再验证，默认3600秒。</p>
					</td>
				</tr>
				<tr>
					<td>最多失败次数</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.maxFails" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">次</span>
						</div>
						<p class="comment"><span v-if="options.maxFails > 0 && options.maxFails < 5" class="red">建议填入一个不小于5的数字，以减少误判几率。</span>允许用户失败尝试的最多次数，超过这个次数将被自动加入黑名单。如果为空或者为0，表示不限制。</p>
					</td>
				</tr>
				<tr>
					<td>失败拦截时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.failBlockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">在达到最多失败次数（大于0）时，自动拦截的时长；如果为0表示不自动拦截。</p>
					</td>
				</tr>
				<tr>
					<td>失败全局封禁</td>
					<td>
						<checkbox v-model="options.failBlockScopeAll"></checkbox>
						<p class="comment">选中后，表示允许系统尝试全局封禁某个IP，以提升封禁性能。</p>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
`
})

// 压缩配置
Vue.component("http-compression-config-box", {
	props: ["v-compression-config", "v-is-location", "v-is-group"],
	mounted: function () {
		let that = this
		sortLoad(function () {
			that.initSortableTypes()
		})
	},
	data: function () {
		let config = this.vCompressionConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				useDefaultTypes: true,
				types: ["brotli", "gzip", "zstd", "deflate"],
				level: 0,
				decompressData: false,
				gzipRef: null,
				deflateRef: null,
				brotliRef: null,
				minLength: {count: 1, "unit": "kb"},
				maxLength: {count: 32, "unit": "mb"},
				mimeTypes: ["text/*", "application/javascript", "application/json", "application/atom+xml", "application/rss+xml", "application/xhtml+xml", "font/*", "image/svg+xml"],
				extensions: [".js", ".json", ".html", ".htm", ".xml", ".css", ".woff2", ".txt"],
				exceptExtensions: [".apk", ".ipa"],
				conds: null,
				enablePartialContent: false,
				onlyURLPatterns: [],
				exceptURLPatterns: []
			}
		}

		if (config.types == null) {
			config.types = []
		}
		if (config.mimeTypes == null) {
			config.mimeTypes = []
		}
		if (config.extensions == null) {
			config.extensions = []
		}

		let allTypes = [
			{
				name: "Gzip",
				code: "gzip",
				isOn: true
			},
			{
				name: "Deflate",
				code: "deflate",
				isOn: true
			},
			{
				name: "Brotli",
				code: "brotli",
				isOn: true
			},
			{
				name: "ZSTD",
				code: "zstd",
				isOn: true
			}
		]

		let configTypes = []
		config.types.forEach(function (typeCode) {
			allTypes.forEach(function (t) {
				if (typeCode == t.code) {
					t.isOn = true
					configTypes.push(t)
				}
			})
		})
		allTypes.forEach(function (t) {
			if (!config.types.$contains(t.code)) {
				t.isOn = false
				configTypes.push(t)
			}
		})

		return {
			config: config,
			moreOptionsVisible: false,
			allTypes: configTypes
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeExtensions: function (values) {
			values.forEach(function (v, k) {
				if (v.length > 0 && v[0] != ".") {
					values[k] = "." + v
				}
			})
			this.config.extensions = values
		},
		changeExceptExtensions: function (values) {
			values.forEach(function (v, k) {
				if (v.length > 0 && v[0] != ".") {
					values[k] = "." + v
				}
			})
			this.config.exceptExtensions = values
		},
		changeMimeTypes: function (values) {
			this.config.mimeTypes = values
		},
		changeAdvancedVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},
		changeConds: function (conds) {
			this.config.conds = conds
		},
		changeType: function () {
			this.config.types = []
			let that = this
			this.allTypes.forEach(function (v) {
				if (v.isOn) {
					that.config.types.push(v.code)
				}
			})
		},
		initSortableTypes: function () {
			let box = document.querySelector("#compression-types-box")
			let that = this
			Sortable.create(box, {
				draggable: ".checkbox",
				handle: ".icon.handle",
				onStart: function () {

				},
				onUpdate: function (event) {
					let checkboxes = box.querySelectorAll(".checkbox")
					let codes = []
					checkboxes.forEach(function (checkbox) {
						let code = checkbox.getAttribute("data-code")
						codes.push(code)
					})
					that.config.types = codes
				}
			})
		}
	},
	template: `<div>
	<input type="hidden" name="compressionJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用内容压缩</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被压缩，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>例外扩展名</td>
				<td>
					<values-box :values="config.exceptExtensions" @change="changeExceptExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将<strong>不会</strong>被压缩，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>支持的MimeType</td>
				<td>
					<values-box :values="config.mimeTypes" @change="changeMimeTypes" placeholder="比如 text/*"></values-box>
					<p class="comment">响应的Content-Type里包含这些MimeType的内容将会被压缩。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>压缩算法</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.useDefaultTypes" id="compression-use-default"/>
						<label v-if="config.useDefaultTypes" for="compression-use-default">使用默认顺序<span class="grey small">（brotli、gzip、 zstd、deflate）</span></label>
						<label v-if="!config.useDefaultTypes" for="compression-use-default">使用默认顺序</label>
					</div>
					<div v-show="!config.useDefaultTypes">
						<div class="ui divider"></div>
						<div id="compression-types-box">
							<div class="ui checkbox" v-for="t in allTypes" style="margin-right: 2em" :data-code="t.code">
								<input type="checkbox" v-model="t.isOn" :id="'compression-type-' + t.code" @change="changeType"/>
								<label :for="'compression-type-' + t.code">{{t.name}} &nbsp; <i class="icon list small grey handle"></i></label>
							</div>
						</div>
					</div>
					
					<p class="comment" v-show="!config.useDefaultTypes">选择支持的压缩算法和优先顺序，拖动<i class="icon list small grey"></i>图表排序。</p>
				</td>
			</tr>
			<tr>
				<td>支持已压缩内容</td>
				<td>
					<checkbox v-model="config.decompressData"></checkbox>
					<p class="comment">支持对已压缩内容尝试重新使用新的算法压缩；不选中表示保留当前的压缩格式。</p>
				</td>
			</tr>
			<tr>
				<td>内容最小长度</td>
				<td>
					<size-capacity-box :v-name="'minLength'" :v-value="config.minLength" :v-unit="'kb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>内容最大长度</td>
				<td>
					<size-capacity-box :v-name="'maxLength'" :v-value="config.maxLength" :v-unit="'mb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>支持Partial<br/>Content</td>
				<td>
					<checkbox v-model="config.enablePartialContent"></checkbox>
					<p class="comment">支持对分片内容（PartialContent）的压缩；除非客户端有特殊要求，一般不需要启用。</p>
				</td>
			</tr>
				<tr>
				<td>例外URL</td>
				<td>
					<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
				</td>
			</tr>
			<tr>
				<td>限制URL</td>
				<td>
					<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了限制URL，表示只对这些URL进行压缩处理；如果不填则表示支持所有的URL。</p>
				</td>
			</tr>
			<tr>
				<td>匹配条件</td>
				<td>
					<http-request-conds-box :v-conds="config.conds" @change="changeConds"></http-request-conds-box>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

// HTTP CC防护配置
Vue.component("http-cc-config-box", {
	props: ["v-cc-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vCcConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				enableFingerprint: true,
				enableGET302: true,
				onlyURLPatterns: [],
				exceptURLPatterns: [],
				useDefaultThresholds: true,
				ignoreCommonFiles: true
			}
		}

		if (config.thresholds == null || config.thresholds.length == 0) {
			config.thresholds = [
				{
					maxRequests: 0
				},
				{
					maxRequests: 0
				},
				{
					maxRequests: 0
				}
			]
		}

		if (typeof config.enableFingerprint != "boolean") {
			config.enableFingerprint = true
		}
		if (typeof config.enableGET302 != "boolean") {
			config.enableGET302 = true
		}

		if (config.onlyURLPatterns == null) {
			config.onlyURLPatterns = []
		}
		if (config.exceptURLPatterns == null) {
			config.exceptURLPatterns = []
		}
		return {
			config: config,
			moreOptionsVisible: false,
			minQPSPerIP: config.minQPSPerIP,
			useCustomThresholds: !config.useDefaultThresholds,

			thresholdMaxRequests0: this.maxRequestsStringAtThresholdIndex(config, 0),
			thresholdMaxRequests1: this.maxRequestsStringAtThresholdIndex(config, 1),
			thresholdMaxRequests2: this.maxRequestsStringAtThresholdIndex(config, 2)
		}
	},
	watch: {
		minQPSPerIP: function (v) {
			let qps = parseInt(v.toString())
			if (isNaN(qps) || qps < 0) {
				qps = 0
			}
			this.config.minQPSPerIP = qps
		},
		thresholdMaxRequests0: function (v) {
			this.setThresholdMaxRequests(0, v)
		},
		thresholdMaxRequests1: function (v) {
			this.setThresholdMaxRequests(1, v)
		},
		thresholdMaxRequests2: function (v) {
			this.setThresholdMaxRequests(2, v)
		},
		useCustomThresholds: function (b) {
			this.config.useDefaultThresholds = !b
		}
	},
	methods: {
		maxRequestsStringAtThresholdIndex: function (config, index) {
			if (config.thresholds == null) {
				return ""
			}
			if (index < config.thresholds.length) {
				let s = config.thresholds[index].maxRequests.toString()
				if (s == "0") {
					s = ""
				}
				return s
			}
			return ""
		},
		setThresholdMaxRequests: function (index, v) {
			let maxRequests = parseInt(v)
			if (isNaN(maxRequests) || maxRequests < 0) {
				maxRequests = 0
			}
			if (index < this.config.thresholds.length) {
				this.config.thresholds[index].maxRequests = maxRequests
			}
		},
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		}
	},
	template: `<div>
<input type="hidden" name="ccJSON" :value="JSON.stringify(config)"/>
<table class="ui table definition selectable">
	<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
	<tbody v-show="((!vIsLocation && !vIsGroup) || config.isPrior)">
		<tr>
			<td class="title">启用CC无感防护</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment"><plus-label></plus-label>启用后，自动检测并拦截CC攻击。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="config.isOn">
		<tr>
			<td colspan="2"><more-options-indicator @change="showMoreOptions"></more-options-indicator></td>
		</tr>
	</tbody>
	<tbody v-show="moreOptionsVisible && config.isOn">
		<tr>
			<td>例外URL</td>
			<td>
				<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
				<p class="comment">如果填写了例外URL，表示这些URL跳过CC防护不做处理。</p>
			</td>
		</tr>
		<tr>
			<td>限制URL</td>
			<td>
				<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
				<p class="comment">如果填写了限制URL，表示只对这些URL进行CC防护处理；如果不填则表示支持所有的URL。</p>
			</td>
		</tr>
		<tr>
			<td>忽略常用文件</td>
			<td>
				<checkbox v-model="config.ignoreCommonFiles"></checkbox>
				<p class="comment">忽略js、css、jpg等常在网页里被引用的文件名，即对这些文件的访问不加入计数，可以减少误判几率。</p>
			</td>
		</tr>
		<tr>
			<td>检查请求来源指纹</td>
			<td>
				<checkbox v-model="config.enableFingerprint"></checkbox>
				<p class="comment">在接收到HTTPS请求时尝试检查请求来源的指纹，用来检测代理服务和爬虫攻击；如果你在网站前面放置了别的反向代理服务，请取消此选项。</p>
			</td>
		</tr>
		<tr>
			<td>启用GET302校验</td>
			<td>
				<checkbox v-model="config.enableGET302"></checkbox>
				<p class="comment">选中后，表示自动通过GET302方法来校验客户端。</p>
			</td>
		</tr>
		<tr>
			<td>单IP最低QPS</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" name="minQPSPerIP" maxlength="6" style="width: 6em" v-model="minQPSPerIP"/>
					<span class="ui label">请求数/秒</span>
				</div>
				<p class="comment">当某个IP在1分钟内平均QPS达到此值时，才会开始检测；如果设置为0，表示任何访问都会检测。（注意这里设置的是检测开启阈值，不是拦截阈值，拦截阈值在当前表单下方可以设置）</p>
			</td>
		</tr>
		<tr>
			<td class="color-border">使用自定义拦截阈值</td>
			<td>
				<checkbox v-model="useCustomThresholds"></checkbox>
			</td>
		</tr>
		<tr v-show="!config.useDefaultThresholds">
			<td class="color-border">自定义拦截阈值设置</td>
			<td>
				<div>
					<div class="ui input left right labeled">
						<span class="ui label basic">单IP每5秒最多</span>
						<input type="text" style="width: 6em" maxlength="6" v-model="thresholdMaxRequests0"/>
						<span class="ui label basic">请求</span>
					</div>
				</div>
					
				<div style="margin-top: 1em">
					<div class="ui input left right labeled">
						<span class="ui label basic">单IP每60秒</span>
						<input type="text" style="width: 6em" maxlength="6" v-model="thresholdMaxRequests1"/>
						<span class="ui label basic">请求</span>
					</div>
				</div>
				<div style="margin-top: 1em">
					<div class="ui input left right labeled">
						<span class="ui label basic">单IP每300秒</span>
						<input type="text" style="width: 6em" maxlength="6" v-model="thresholdMaxRequests2"/>
						<span class="ui label basic">请求</span>
					</div>
				</div>
			</td>
		</tr>
	</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`
})

Vue.component("firewall-event-level-options", {
    props: ["v-value"],
    mounted: function () {
        let that = this
        Tea.action("/ui/eventLevelOptions")
            .post()
            .success(function (resp) {
                that.levels = resp.data.eventLevels
                that.change()
            })
    },
    data: function () {
        let value = this.vValue
        if (value == null || value.length == 0) {
            value = "" // 不要给默认值，因为黑白名单等默认值均有不同
        }

        return {
            levels: [],
            description: "",
            level: value
        }
    },
    methods: {
        change: function () {
            this.$emit("change")

            let that = this
            let l = this.levels.$find(function (k, v) {
                return v.code == that.level
            })
            if (l != null) {
                this.description = l.description
            } else {
                this.description = ""
            }
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="eventLevel" v-model="level" @change="change">
        <option v-for="level in levels" :value="level.code">{{level.name}}</option>
    </select>
    <p class="comment">{{description}}</p>
</div>`
})

Vue.component("prior-checkbox", {
	props: ["v-config", "description"],
	data: function () {
		let description = this.description
		if (description == null) {
			description = "打开后可以覆盖父级或子级配置"
		}
		return {
			isPrior: this.vConfig.isPrior,
			realDescription: description
		}
	},
	watch: {
		isPrior: function (v) {
			this.vConfig.isPrior = v
		}
	},
	template: `<tbody>
	<tr :class="{active:isPrior}">
		<td class="title">打开独立配置</td>
		<td>
			<div class="ui toggle checkbox">
				<input type="checkbox" v-model="isPrior"/>
				<label class="red"></label>
			</div>
			<p class="comment"><strong v-if="isPrior">[已打开]</strong> {{realDescription}}。</p>
		</td>
	</tr>
</tbody>`
})

Vue.component("http-charsets-box", {
	props: ["v-usual-charsets", "v-all-charsets", "v-charset-config", "v-is-location", "v-is-group"],
	data: function () {
		let charsetConfig = this.vCharsetConfig
		if (charsetConfig == null) {
			charsetConfig = {
				isPrior: false,
				isOn: false,
				charset: "",
				isUpper: false,
				force: false
			}
		}
		return {
			charsetConfig: charsetConfig,
			advancedVisible: false
		}
	},
	methods: {
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		}
	},
	template: `<div>
	<input type="hidden" name="charsetJSON" :value="JSON.stringify(charsetConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="charsetConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || charsetConfig.isPrior">
			<tr>
				<td class="title">启用字符编码</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="charsetConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn">	
			<tr>
				<td class="title">选择字符编码</td>
				<td><select class="ui dropdown" style="width:20em" name="charset" v-model="charsetConfig.charset">
						<option value="">[未选择]</option>
						<optgroup label="常用字符编码"></optgroup>
						<option v-for="charset in vUsualCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
						<optgroup label="全部字符编码"></optgroup>
						<option v-for="charset in vAllCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
					</select>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn"></more-options-tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn && advancedVisible">
			<tr>
				<td>强制替换</td>
				<td>
					<checkbox v-model="charsetConfig.force"></checkbox>
					<p class="comment">选中后，表示强制覆盖已经设置的字符集；不选中，表示如果源站已经设置了字符集，则保留不修改。</p>
				</td>
			</tr>
			<tr>
				<td>字符编码大写</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="charsetConfig.isUpper"/>
						<label></label>
					</div>
					<p class="comment">选中后将指定的字符编码转换为大写，比如默认为<code-label>utf-8</code-label>，选中后将改为<code-label>UTF-8</code-label>。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-expires-time-config-box", {
	props: ["v-expires-time"],
	data: function () {
		let expiresTime = this.vExpiresTime
		if (expiresTime == null) {
			expiresTime = {
				isPrior: false,
				isOn: false,
				overwrite: true,
				autoCalculate: true,
				duration: {count: -1, "unit": "hour"}
			}
		}
		return {
			expiresTime: expiresTime
		}
	},
	watch: {
		"expiresTime.isPrior": function () {
			this.notifyChange()
		},
		"expiresTime.isOn": function () {
			this.notifyChange()
		},
		"expiresTime.overwrite": function () {
			this.notifyChange()
		},
		"expiresTime.autoCalculate": function () {
			this.notifyChange()
		}
	},
	methods: {
		notifyChange: function () {
			this.$emit("change", this.expiresTime)
		}
	},
	template: `<div>
	<table class="ui table">
		<prior-checkbox :v-config="expiresTime"></prior-checkbox>
		<tbody v-show="expiresTime.isPrior">
			<tr>
				<td class="title">启用</td>
				<td><checkbox v-model="expiresTime.isOn"></checkbox>
					<p class="comment">启用后，将会在响应的Header中添加<code-label>Expires</code-label>字段，浏览器据此会将内容缓存在客户端；同时，在管理后台执行清理缓存时，也将无法清理客户端已有的缓存。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn">
				<td>覆盖源站设置</td>
				<td>
					<checkbox v-model="expiresTime.overwrite"></checkbox>
					<p class="comment">选中后，会覆盖源站Header中已有的<code-label>Expires</code-label>字段。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn">
				<td>自动计算时间</td>
				<td><checkbox v-model="expiresTime.autoCalculate"></checkbox>
					<p class="comment">根据已设置的缓存有效期进行计算。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn && !expiresTime.autoCalculate">
				<td>强制缓存时间</td>
				<td>
					<time-duration-box :v-value="expiresTime.duration" @change="notifyChange"></time-duration-box>
					<p class="comment">从客户端访问的时间开始要缓存的时长。</p>
				</td>
			</tr>
		</tbody>
	</table>
</div>`
})

Vue.component("http-access-log-box", {
	props: ["v-access-log", "v-keyword", "v-show-server-link"],
	data: function () {
		let accessLog = this.vAccessLog
		if (accessLog.header != null && accessLog.header.Upgrade != null && accessLog.header.Upgrade.values != null && accessLog.header.Upgrade.values.$contains("websocket")) {
			if (accessLog.scheme == "http") {
				accessLog.scheme = "ws"
			} else if (accessLog.scheme == "https") {
				accessLog.scheme = "wss"
			}
		}

		// 对TAG去重
		if (accessLog.tags != null && accessLog.tags.length > 0) {
			let tagMap = {}
			accessLog.tags = accessLog.tags.$filter(function (k, tag) {
				let b = (typeof (tagMap[tag]) == "undefined")
				tagMap[tag] = true
				return b
			})
		}

		// 域名
		accessLog.unicodeHost = ""
		if (accessLog.host != null && accessLog.host.startsWith("xn--")) {
			// port
			let portIndex = accessLog.host.indexOf(":")
			if (portIndex > 0) {
				accessLog.unicodeHost = punycode.ToUnicode(accessLog.host.substring(0, portIndex))
			} else {
				accessLog.unicodeHost = punycode.ToUnicode(accessLog.host)
			}
		}

		return {
			accessLog: accessLog
		}
	},
	methods: {
		formatCost: function (seconds) {
			if (seconds == null) {
				return "0"
			}
			let s = (seconds * 1000).toString();
			let pieces = s.split(".");
			if (pieces.length < 2) {
				return s;
			}

			return pieces[0] + "." + pieces[1].substring(0, 3);
		},
		showLog: function () {
			let that = this
			let requestId = this.accessLog.requestId
			this.$parent.$children.forEach(function (v) {
				if (v.deselect != null) {
					v.deselect()
				}
			})
			this.select()
			teaweb.popup("/servers/server/log/viewPopup?requestId=" + requestId, {
				width: "50em",
				height: "28em",
				onClose: function () {
					that.deselect()
				}
			})
		},
		select: function () {
			this.$refs.box.parentNode.style.cssText = "background: rgba(0, 0, 0, 0.1)"
		},
		deselect: function () {
			this.$refs.box.parentNode.style.cssText = ""
		},
		mismatch: function () {
			teaweb.warn("当前访问没有匹配到任何网站")
		}
	},
	template: `<div style="word-break: break-all" :style="{'color': (accessLog.status >= 400) ? '#dc143c' : ''}" ref="box">
	<div>
		<a v-if="accessLog.node != null && accessLog.node.nodeCluster != null" :href="'/clusters/cluster/node?nodeId=' + accessLog.node.id + '&clusterId=' + accessLog.node.nodeCluster.id" title="点击查看节点详情" target="_top"><span class="grey">[{{accessLog.node.name}}<span v-if="!accessLog.node.name.endsWith('节点')">节点</span>]</span></a>
		
		<!-- 网站 -->
		<a :href="'/servers/server/log?serverId=' + accessLog.serverId" title="点击到网站" v-if="vShowServerLink && accessLog.serverId > 0"><span class="grey">[网站]</span></a>
		<span v-if="vShowServerLink && (accessLog.serverId == null || accessLog.serverId == 0)" @click.prevent="mismatch()"><span class="disabled">[网站]</span></span>
		
		<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey"><ip-box :v-ip="accessLog.remoteAddr">[{{accessLog.region}}]</ip-box></span> 
		<ip-box><keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword></ip-box> [{{accessLog.timeLocal}}] <em>&quot;<keyword :v-word="vKeyword">{{accessLog.requestMethod}}</keyword> {{accessLog.scheme}}://<keyword :v-word="vKeyword">{{accessLog.host}}</keyword><keyword :v-word="vKeyword">{{accessLog.requestURI}}</keyword> <a :href="accessLog.scheme + '://' + accessLog.host + accessLog.requestURI" target="_blank" title="新窗口打开" class="disabled"><i class="external icon tiny"></i> </a> {{accessLog.proto}}&quot; </em> <keyword :v-word="vKeyword">{{accessLog.status}}</keyword> 
		
		<code-label v-if="accessLog.unicodeHost != null && accessLog.unicodeHost.length > 0">{{accessLog.unicodeHost}}</code-label>
		
		<!-- attrs -->
		<code-label v-if="accessLog.attrs != null && (accessLog.attrs['cache.status'] == 'HIT' || accessLog.attrs['cache.status'] == 'STALE')">cache {{accessLog.attrs['cache.status'].toLowerCase()}}</code-label> 
		<!-- waf -->
		<code-label v-if="accessLog.firewallActions != null && accessLog.firewallActions.length > 0">waf {{accessLog.firewallActions}}</code-label> 
		
		<!-- tags -->
		<span v-if="accessLog.tags != null && accessLog.tags.length > 0">- <code-label v-for="tag in accessLog.tags" :key="tag">{{tag}}</code-label>
		</span>
		<span  v-if="accessLog.wafInfo != null">
			<a :href="(accessLog.wafInfo.policy.serverId == 0) ? '/servers/components/waf/group?firewallPolicyId=' +  accessLog.firewallPolicyId + '&type=inbound&groupId=' + accessLog.firewallRuleGroupId+ '#set' + accessLog.firewallRuleSetId : '/servers/server/settings/waf/group?serverId=' + accessLog.serverId + '&firewallPolicyId=' + accessLog.firewallPolicyId + '&type=inbound&groupId=' + accessLog.firewallRuleGroupId + '#set' + accessLog.firewallRuleSetId" target="_blank">
				<code-label-plain>
					<span>
						WAF -
						<span v-if="accessLog.wafInfo.group != null">{{accessLog.wafInfo.group.name}} -</span>
						<span v-if="accessLog.wafInfo.set != null">{{accessLog.wafInfo.set.name}}</span>
					</span>
				</code-label-plain>
			</a>
		</span>
			
		<span v-if="accessLog.requestTime != null"> - 耗时:{{formatCost(accessLog.requestTime)}} ms </span><span v-if="accessLog.humanTime != null && accessLog.humanTime.length > 0" class="grey small">&nbsp; ({{accessLog.humanTime}})</span>
		&nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>
	</div>
</div>`
})

// Javascript Punycode converter derived from example in RFC3492.
// This implementation is created by some@domain.name and released into public domain
// 代码来自：https://stackoverflow.com/questions/183485/converting-punycode-with-dash-character-to-unicode
var punycode = new function Punycode() {
	// This object converts to and from puny-code used in IDN
	//
	// punycode.ToASCII ( domain )
	//
	// Returns a puny coded representation of "domain".
	// It only converts the part of the domain name that
	// has non ASCII characters. I.e. it dosent matter if
	// you call it with a domain that already is in ASCII.
	//
	// punycode.ToUnicode (domain)
	//
	// Converts a puny-coded domain name to unicode.
	// It only converts the puny-coded parts of the domain name.
	// I.e. it dosent matter if you call it on a string
	// that already has been converted to unicode.
	//
	//
	this.utf16 = {
		// The utf16-class is necessary to convert from javascripts internal character representation to unicode and back.
		decode: function (input) {
			var output = [], i = 0, len = input.length, value, extra;
			while (i < len) {
				value = input.charCodeAt(i++);
				if ((value & 0xF800) === 0xD800) {
					extra = input.charCodeAt(i++);
					if (((value & 0xFC00) !== 0xD800) || ((extra & 0xFC00) !== 0xDC00)) {
						throw new RangeError("UTF-16(decode): Illegal UTF-16 sequence");
					}
					value = ((value & 0x3FF) << 10) + (extra & 0x3FF) + 0x10000;
				}
				output.push(value);
			}
			return output;
		},
		encode: function (input) {
			var output = [], i = 0, len = input.length, value;
			while (i < len) {
				value = input[i++];
				if ((value & 0xF800) === 0xD800) {
					throw new RangeError("UTF-16(encode): Illegal UTF-16 value");
				}
				if (value > 0xFFFF) {
					value -= 0x10000;
					output.push(String.fromCharCode(((value >>> 10) & 0x3FF) | 0xD800));
					value = 0xDC00 | (value & 0x3FF);
				}
				output.push(String.fromCharCode(value));
			}
			return output.join("");
		}
	}

	//Default parameters
	var initial_n = 0x80;
	var initial_bias = 72;
	var delimiter = "\x2D";
	var base = 36;
	var damp = 700;
	var tmin = 1;
	var tmax = 26;
	var skew = 38;
	var maxint = 0x7FFFFFFF;

	// decode_digit(cp) returns the numeric value of a basic code
	// point (for use in representing integers) in the range 0 to
	// base-1, or base if cp is does not represent a value.

	function decode_digit(cp) {
		return cp - 48 < 10 ? cp - 22 : cp - 65 < 26 ? cp - 65 : cp - 97 < 26 ? cp - 97 : base;
	}

	// encode_digit(d,flag) returns the basic code point whose value
	// (when used for representing integers) is d, which needs to be in
	// the range 0 to base-1. The lowercase form is used unless flag is
	// nonzero, in which case the uppercase form is used. The behavior
	// is undefined if flag is nonzero and digit d has no uppercase form.

	function encode_digit(d, flag) {
		return d + 22 + 75 * (d < 26) - ((flag != 0) << 5);
		//  0..25 map to ASCII a..z or A..Z
		// 26..35 map to ASCII 0..9
	}

	//** Bias adaptation function **
	function adapt(delta, numpoints, firsttime) {
		var k;
		delta = firsttime ? Math.floor(delta / damp) : (delta >> 1);
		delta += Math.floor(delta / numpoints);

		for (k = 0; delta > (((base - tmin) * tmax) >> 1); k += base) {
			delta = Math.floor(delta / (base - tmin));
		}
		return Math.floor(k + (base - tmin + 1) * delta / (delta + skew));
	}

	// encode_basic(bcp,flag) forces a basic code point to lowercase if flag is zero,
	// uppercase if flag is nonzero, and returns the resulting code point.
	// The code point is unchanged if it is caseless.
	// The behavior is undefined if bcp is not a basic code point.

	function encode_basic(bcp, flag) {
		bcp -= (bcp - 97 < 26) << 5;
		return bcp + ((!flag && (bcp - 65 < 26)) << 5);
	}

	// Main decode
	this.decode = function (input, preserveCase) {
		// Dont use utf16
		var output = [];
		var case_flags = [];
		var input_length = input.length;

		var n, out, i, bias, basic, j, ic, oldi, w, k, digit, t, len;

		// Initialize the state:

		n = initial_n;
		i = 0;
		bias = initial_bias;

		// Handle the basic code points: Let basic be the number of input code
		// points before the last delimiter, or 0 if there is none, then
		// copy the first basic code points to the output.

		basic = input.lastIndexOf(delimiter);
		if (basic < 0) basic = 0;

		for (j = 0; j < basic; ++j) {
			if (preserveCase) case_flags[output.length] = (input.charCodeAt(j) - 65 < 26);
			if (input.charCodeAt(j) >= 0x80) {
				throw new RangeError("Illegal input >= 0x80");
			}
			output.push(input.charCodeAt(j));
		}

		// Main decoding loop: Start just after the last delimiter if any
		// basic code points were copied; start at the beginning otherwise.

		for (ic = basic > 0 ? basic + 1 : 0; ic < input_length;) {

			// ic is the index of the next character to be consumed,

			// Decode a generalized variable-length integer into delta,
			// which gets added to i. The overflow checking is easier
			// if we increase i as we go, then subtract off its starting
			// value at the end to obtain delta.
			for (oldi = i, w = 1, k = base; ; k += base) {
				if (ic >= input_length) {
					throw RangeError("punycode_bad_input(1)");
				}
				digit = decode_digit(input.charCodeAt(ic++));

				if (digit >= base) {
					throw RangeError("punycode_bad_input(2)");
				}
				if (digit > Math.floor((maxint - i) / w)) {
					throw RangeError("punycode_overflow(1)");
				}
				i += digit * w;
				t = k <= bias ? tmin : k >= bias + tmax ? tmax : k - bias;
				if (digit < t) {
					break;
				}
				if (w > Math.floor(maxint / (base - t))) {
					throw RangeError("punycode_overflow(2)");
				}
				w *= (base - t);
			}

			out = output.length + 1;
			bias = adapt(i - oldi, out, oldi === 0);

			// i was supposed to wrap around from out to 0,
			// incrementing n each time, so we'll fix that now:
			if (Math.floor(i / out) > maxint - n) {
				throw RangeError("punycode_overflow(3)");
			}
			n += Math.floor(i / out);
			i %= out;

			// Insert n at position i of the output:
			// Case of last character determines uppercase flag:
			if (preserveCase) {
				case_flags.splice(i, 0, input.charCodeAt(ic - 1) - 65 < 26);
			}

			output.splice(i, 0, n);
			i++;
		}
		if (preserveCase) {
			for (i = 0, len = output.length; i < len; i++) {
				if (case_flags[i]) {
					output[i] = (String.fromCharCode(output[i]).toUpperCase()).charCodeAt(0);
				}
			}
		}
		return this.utf16.encode(output);
	};

	//** Main encode function **

	this.encode = function (input, preserveCase) {
		//** Bias adaptation function **

		var n, delta, h, b, bias, j, m, q, k, t, ijv, case_flags;

		if (preserveCase) {
			// Preserve case, step1 of 2: Get a list of the unaltered string
			case_flags = this.utf16.decode(input);
		}
		// Converts the input in UTF-16 to Unicode
		input = this.utf16.decode(input.toLowerCase());

		var input_length = input.length; // Cache the length

		if (preserveCase) {
			// Preserve case, step2 of 2: Modify the list to true/false
			for (j = 0; j < input_length; j++) {
				case_flags[j] = input[j] != case_flags[j];
			}
		}

		var output = [];


		// Initialize the state:
		n = initial_n;
		delta = 0;
		bias = initial_bias;

		// Handle the basic code points:
		for (j = 0; j < input_length; ++j) {
			if (input[j] < 0x80) {
				output.push(
					String.fromCharCode(
						case_flags ? encode_basic(input[j], case_flags[j]) : input[j]
					)
				);
			}
		}

		h = b = output.length;

		// h is the number of code points that have been handled, b is the
		// number of basic code points

		if (b > 0) output.push(delimiter);

		// Main encoding loop:
		//
		while (h < input_length) {
			// All non-basic code points < n have been
			// handled already. Find the next larger one:

			for (m = maxint, j = 0; j < input_length; ++j) {
				ijv = input[j];
				if (ijv >= n && ijv < m) m = ijv;
			}

			// Increase delta enough to advance the decoder's
			// <n,i> state to <m,0>, but guard against overflow:

			if (m - n > Math.floor((maxint - delta) / (h + 1))) {
				throw RangeError("punycode_overflow (1)");
			}
			delta += (m - n) * (h + 1);
			n = m;

			for (j = 0; j < input_length; ++j) {
				ijv = input[j];

				if (ijv < n) {
					if (++delta > maxint) return Error("punycode_overflow(2)");
				}

				if (ijv == n) {
					// Represent delta as a generalized variable-length integer:
					for (q = delta, k = base; ; k += base) {
						t = k <= bias ? tmin : k >= bias + tmax ? tmax : k - bias;
						if (q < t) break;
						output.push(String.fromCharCode(encode_digit(t + (q - t) % (base - t), 0)));
						q = Math.floor((q - t) / (base - t));
					}
					output.push(String.fromCharCode(encode_digit(q, preserveCase && case_flags[j] ? 1 : 0)));
					bias = adapt(delta, h + 1, h == b);
					delta = 0;
					++h;
				}
			}

			++delta, ++n;
		}
		return output.join("");
	}

	this.ToASCII = function (domain) {
		var domain_array = domain.split(".");
		var out = [];
		for (var i = 0; i < domain_array.length; ++i) {
			var s = domain_array[i];
			out.push(
				s.match(/[^A-Za-z0-9-]/) ?
					"xn--" + punycode.encode(s) :
					s
			);
		}
		return out.join(".");
	}
	this.ToUnicode = function (domain) {
		var domain_array = domain.split(".");
		var out = [];
		for (var i = 0; i < domain_array.length; ++i) {
			var s = domain_array[i];
			out.push(
				s.match(/^xn--/) ?
					punycode.decode(s.slice(4)) :
					s
			);
		}
		return out.join(".");
	}
}();

Vue.component("http-firewall-block-options-viewer", {
	props: ["v-block-options"],
	data: function () {
		return {
			options: this.vBlockOptions
		}
	},
	template: `<div>
	<span v-if="options == null">默认设置</span>
	<div v-else>
		状态码：{{options.statusCode}} / 提示内容：<span v-if="options.body != null && options.body.length > 0">[{{options.body.length}}字符]</span><span v-else class="disabled">[无]</span>  / 超时时间：{{options.timeout}}秒 <span v-if="options.timeoutMax > options.timeout">/ 最大封禁时长：{{options.timeoutMax}}秒</span>
		<span v-if="options.failBlockScopeAll"> / 尝试全局封禁</span>
	</div>
</div>	
`
})

Vue.component("http-access-log-config-box", {
	props: ["v-access-log-config", "v-fields", "v-default-field-codes", "v-is-location", "v-is-group"],
	data: function () {
		let that = this

		// 初始化
		setTimeout(function () {
			that.changeFields()
		}, 100)

		let accessLog = {
			isPrior: false,
			isOn: false,
			fields: [1, 2, 6, 7],
			status1: true,
			status2: true,
			status3: true,
			status4: true,
			status5: true,

            firewallOnly: false,
			enableClientClosed: false
		}
		if (this.vAccessLogConfig != null) {
			accessLog = this.vAccessLogConfig
		}

		this.vFields.forEach(function (v) {
			if (that.vAccessLogConfig == null) { // 初始化默认值
				v.isChecked = that.vDefaultFieldCodes.$contains(v.code)
			} else {
				v.isChecked = accessLog.fields.$contains(v.code)
			}
		})

		return {
			accessLog: accessLog,
			hasRequestBodyField: this.vFields.$contains(8),
			showAdvancedOptions: false
		}
	},
	methods: {
		changeFields: function () {
			this.accessLog.fields = this.vFields.filter(function (v) {
				return v.isChecked
			}).map(function (v) {
				return v.code
			})
			this.hasRequestBodyField = this.accessLog.fields.$contains(8)
		},
		changeAdvanced: function (v) {
			this.showAdvancedOptions = v
		}
	},
	template: `<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(accessLog)"/>
	<table class="ui table definition selectable" :class="{'opacity-mask': this.accessLog.firewallOnly}">
		<prior-checkbox :v-config="accessLog" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || accessLog.isPrior">
			<tr>
				<td class="title">启用访问日志</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn">
			<tr>
				<td colspan="2"><more-options-indicator @change="changeAdvanced"></more-options-indicator></td>
			</tr>
		</tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn && showAdvancedOptions">
			<tr>
				<td>基础信息</td>
				<td><p class="comment" style="padding-top: 0">默认记录客户端IP、请求URL等基础信息。</p></td>
			</tr>
			<tr>
				<td>高级信息</td>
				<td>
					<div class="ui checkbox" v-for="(field, index) in vFields" style="width:10em;margin-bottom:0.8em">
						<input type="checkbox" v-model="field.isChecked" @change="changeFields" :id="'access-log-field-' + index"/>
						<label :for="'access-log-field-' + index">{{field.name}}</label>
					</div>
					<p class="comment">在基础信息之外要存储的信息。
						<span class="red" v-if="hasRequestBodyField">记录"请求Body"将会显著消耗更多的系统资源，建议仅在调试时启用，最大记录尺寸为2MiB。</span>
					</p>
				</td>
			</tr>
			<tr>
				<td>要存储的访问日志状态码</td>
				<td>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status1"/>
						<label>1xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status2"/>
						<label>2xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status3"/>
						<label>3xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status4"/>
						<label>4xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status5"/>
						<label>5xx</label>
					</div>
				</td>
			</tr>
			<tr>
				<td>记录客户端中断日志</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.enableClientClosed"/>
						<label></label>
					</div>
					<p class="comment">以<code-label>499</code-label>的状态码记录客户端主动中断日志。</p>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn">
        <h4>WAF相关</h4>
        <table class="ui table definition selectable">
            <tr>
                <td class="title">只记录WAF相关日志</td>
                <td>
                    <checkbox v-model="accessLog.firewallOnly"></checkbox>
                    <p class="comment">选中后只记录WAF相关的日志。通过此选项可有效减少访问日志数量，降低网络带宽和存储压力。</p>
                </td>
            </tr>
        </table>
    </div>
	<div class="margin"></div>
</div>`
})

// 基本认证用户配置
Vue.component("http-auth-basic-auth-user-box", {
	props: ["v-users"],
	data: function () {
		let users = this.vUsers
		if (users == null) {
			users = []
		}
		return {
			users: users,
			isAdding: false,
			updatingIndex: -1,

			username: "",
			password: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.username = ""
			this.password = ""

			let that = this
			setTimeout(function () {
				that.$refs.username.focus()
			}, 100)
		},
		cancel: function () {
			this.isAdding = false
			this.updatingIndex = -1
		},
		confirm: function () {
			let that = this
			if (this.username.length == 0) {
				teaweb.warn("请输入用户名", function () {
					that.$refs.username.focus()
				})
				return
			}
			if (this.password.length == 0) {
				teaweb.warn("请输入密码", function () {
					that.$refs.password.focus()
				})
				return
			}
			if (this.updatingIndex < 0) {
				this.users.push({
					username: this.username,
					password: this.password
				})
			} else {
				this.users[this.updatingIndex].username = this.username
				this.users[this.updatingIndex].password = this.password
			}
			this.cancel()
		},
		update: function (index, user) {
			this.updatingIndex = index

			this.isAdding = true
			this.username = user.username
			this.password = user.password

			let that = this
			setTimeout(function () {
				that.$refs.username.focus()
			}, 100)
		},
		remove: function (index) {
			this.users.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="httpAuthBasicAuthUsersJSON" :value="JSON.stringify(users)"/>
	<div v-if="users.length > 0">
		<div class="ui label small basic" v-for="(user, index) in users">
			{{user.username}} <a href="" title="修改" @click.prevent="update(index, user)"><i class="icon pencil tiny"></i></a>
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-show="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="用户名" v-model="username" size="15" ref="username"/>
			</div>
			<div class="ui field">
				<input type="password" placeholder="密码" v-model="password" size="15" ref="password"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>&nbsp;
				<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding" style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("http-location-labels", {
	props: ["v-location-config", "v-server-id"],
	data: function () {
		return {
			location: this.vLocationConfig
		}
	},
	methods: {
		// 判断是否已启用某配置
		configIsOn: function (config) {
			return config != null && config.isPrior && config.isOn
		},

		refIsOn: function (ref, config) {
			return this.configIsOn(ref) && config != null && config.isOn
		},

		len: function (arr) {
			return (arr == null) ? 0 : arr.length
		},
		url: function (path) {
			return "/servers/server/settings/locations" + path + "?serverId=" + this.vServerId + "&locationId=" + this.location.id
		}
	},
	template: `	<div class="labels-box">
	<!-- 基本信息 -->
	<http-location-labels-label v-if="location.name != null && location.name.length > 0" :class="'olive'" :href="url('/location')">{{location.name}}</http-location-labels-label>
	
	<!-- domains -->
	<div v-if="location.domains != null && location.domains.length > 0">
		<grey-label v-for="domain in location.domains">{{domain}}</grey-label>
	</div>
	
	<!-- break -->
	<http-location-labels-label v-if="location.isBreak" :href="url('/location')">BREAK</http-location-labels-label>
	
	<!-- redirectToHTTPS -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.redirectToHTTPS)" :href="url('/http')">自动跳转HTTPS</http-location-labels-label>
	
	<!-- Web -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.root)" :href="url('/web')">文档根目录</http-location-labels-label>
	
	<!-- 反向代理 -->
	<http-location-labels-label v-if="refIsOn(location.reverseProxyRef, location.reverseProxy)" :v-href="url('/reverseProxy')">源站</http-location-labels-label>
	
	<!-- UAM -->
	<http-location-labels-label v-if="location.web != null && location.web.uam != null && location.web.uam.isPrior"><span :class="{disabled: !location.web.uam.isOn, red:location.web.uam.isOn}">5秒盾</span></http-location-labels-label>
	
	<!-- CC -->
	<http-location-labels-label v-if="location.web != null && location.web.cc != null && location.web.cc.isPrior"><span :class="{disabled: !location.web.cc.isOn, red:location.web.cc.isOn}">CC防护</span></http-location-labels-label>
	
	<!-- WAF -->
	<!-- TODO -->
	
	<!-- Cache -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.cache)" :v-href="url('/cache')">CACHE</http-location-labels-label>
	
	<!-- Charset -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.charset) && location.web.charset.charset.length > 0" :href="url('/charset')">{{location.web.charset.charset}}</http-location-labels-label>
	
	<!-- 访问日志 -->
	<!-- TODO -->
	
	<!-- 统计 -->
	<!-- TODO -->
	
	<!-- Gzip -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.gzipRef, location.web.gzip) && location.web.gzip.level > 0" :href="url('/gzip')">Gzip:{{location.web.gzip.level}}</http-location-labels-label>
	
	<!-- HTTP Header -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.requestHeaderPolicyRef, location.web.requestHeaderPolicy) && (len(location.web.requestHeaderPolicy.addHeaders) > 0 || len(location.web.requestHeaderPolicy.setHeaders) > 0 || len(location.web.requestHeaderPolicy.replaceHeaders) > 0 || len(location.web.requestHeaderPolicy.deleteHeaders) > 0)" :href="url('/headers')">请求Header</http-location-labels-label>
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.responseHeaderPolicyRef, location.web.responseHeaderPolicy) && (len(location.web.responseHeaderPolicy.addHeaders) > 0 || len(location.web.responseHeaderPolicy.setHeaders) > 0 || len(location.web.responseHeaderPolicy.replaceHeaders) > 0 || len(location.web.responseHeaderPolicy.deleteHeaders) > 0)" :href="url('/headers')">响应Header</http-location-labels-label>
	
	<!-- Websocket -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.websocketRef, location.web.websocket)" :href="url('/websocket')">Websocket</http-location-labels-label>
	
	<!-- 请求脚本 -->
	<http-location-labels-label v-if="location.web != null && location.web.requestScripts != null && ((location.web.requestScripts.initGroup != null && location.web.requestScripts.initGroup.isPrior) || (location.web.requestScripts.requestGroup != null && location.web.requestScripts.requestGroup.isPrior))" :href="url('/requestScripts')">请求脚本</http-location-labels-label>
	
	<!-- 自定义访客IP地址 -->
	<http-location-labels-label v-if="location.web != null && location.web.remoteAddr != null && location.web.remoteAddr.isPrior" :href="url('/remoteAddr')" :class="{disabled: !location.web.remoteAddr.isOn}">访客IP地址</http-location-labels-label>
	
	<!-- 请求限制 -->
	<http-location-labels-label v-if="location.web != null && location.web.requestLimit != null && location.web.requestLimit.isPrior" :href="url('/requestLimit')" :class="{disabled: !location.web.requestLimit.isOn}">请求限制</http-location-labels-label>
		
	<!-- 自定义页面 -->
	<div v-if="location.web != null && location.web.pages != null && location.web.pages.length > 0">
		<div v-for="page in location.web.pages" :key="page.id"><http-location-labels-label :href="url('/pages')">PAGE [状态码{{page.status[0]}}] -&gt; {{page.url}}</http-location-labels-label></div>
	</div>
	<div v-if="location.web != null && configIsOn(location.web.shutdown)">
		<http-location-labels-label :v-class="'red'" :href="url('/pages')">临时关闭</http-location-labels-label>
	</div>
	
	<!-- 重写规则 -->
	<div v-if="location.web != null && location.web.rewriteRules != null && location.web.rewriteRules.length > 0">
		<div v-for="rewriteRule in location.web.rewriteRules">
			<http-location-labels-label :href="url('/rewrite')">REWRITE {{rewriteRule.pattern}} -&gt; {{rewriteRule.replace}}</http-location-labels-label>
		</div>
	</div>
</div>`
})

Vue.component("http-location-labels-label", {
	props: ["v-class", "v-href"],
	template: `<a :href="vHref" class="ui label tiny basic" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></a>`
})

Vue.component("http-gzip-box", {
	props: ["v-gzip-config", "v-gzip-ref", "v-is-location"],
	data: function () {
		let gzip = this.vGzipConfig
		if (gzip == null) {
			gzip = {
				isOn: true,
				level: 0,
				minLength: null,
				maxLength: null,
				conds: null
			}
		}

		return {
			gzip: gzip,
			advancedVisible: false
		}
	},
	methods: {
		isOn: function () {
			return (!this.vIsLocation || this.vGzipRef.isPrior) && this.vGzipRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		}
	},
	template: `<div>
<input type="hidden" name="gzipRefJSON" :value="JSON.stringify(vGzipRef)"/> 
<table class="ui table selectable definition">
	<prior-checkbox :v-config="vGzipRef" v-if="vIsLocation"></prior-checkbox>
	<tbody v-show="!vIsLocation || vGzipRef.isPrior">
		<tr>
			<td class="title">启用Gzip压缩</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" v-model="vGzipRef.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
	<tbody v-show="isOn()">
		<tr>
			<td class="title">压缩级别</td>
			<td>
				<select class="dropdown auto-width" name="level" v-model="gzip.level">
					<option value="0">不压缩</option>
					<option v-for="i in 9" :value="i">{{i}}</option>
				</select>
				<p class="comment">级别越高，压缩比例越大。</p>
			</td>
		</tr>
	</tbody>
	<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
	<tbody v-show="isOn() && advancedVisible">
		<tr>
			<td>Gzip内容最小长度</td>
			<td>
				<size-capacity-box :v-name="'minLength'" :v-value="gzip.minLength" :v-unit="'kb'"></size-capacity-box>
				<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
			</td>
		</tr>
		<tr>
			<td>Gzip内容最大长度</td>
			<td>
				<size-capacity-box :v-name="'maxLength'" :v-value="gzip.maxLength" :v-unit="'mb'"></size-capacity-box>
				<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
			</td>
		</tr>
		<tr>
			<td>匹配条件</td>
			<td>
				<http-request-conds-box :v-conds="gzip.conds"></http-request-conds-box>
</td>
		</tr>
	</tbody>
</table>
</div>`
})

Vue.component("script-config-box", {
	props: ["id", "v-script-config", "comment", "v-auditing-status"],
	mounted: function () {
		let that = this
		setTimeout(function () {
			that.$forceUpdate()
		}, 100)
	},
	data: function () {
		let config = this.vScriptConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				code: "",
				auditingCode: ""
			}
		}

		let auditingStatus = null
		if (config.auditingCodeMD5 != null && config.auditingCodeMD5.length > 0 && config.auditingCode != null && config.auditingCode.length > 0) {
			config.code = config.auditingCode

			if (this.vAuditingStatus != null) {
				for (let i = 0; i < this.vAuditingStatus.length; i++) {
					let status = this.vAuditingStatus[i]
					if (status.md5 == config.auditingCodeMD5) {
						auditingStatus = status
						break
					}
				}
			}
		}

		if (config.code.length == 0) {
			config.code = "\n\n\n\n"
		}

		return {
			config: config,
			auditingStatus: auditingStatus
		}
	},
	watch: {
		"config.isOn": function () {
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.config)
		},
		changeCode: function (code) {
			this.config.code = code
			this.change()
		},
		isPlus: function () {
			if (Tea == null || Tea.Vue == null) {
				return false
			}
			return Tea.Vue.teaIsPlus
		}
	},
	template: `<div>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">启用脚本设置</td>
				<td><checkbox v-model="config.isOn"></checkbox></td>
			</tr>
		</tbody>
		<tbody>
			<tr :style="{opacity: !config.isOn ? 0.5 : 1}">
				<td>脚本代码</td>	
				<td>
					<p class="comment" v-if="auditingStatus != null">
						<span class="green" v-if="auditingStatus.isPassed">管理员审核结果：审核通过。</span>
						<span class="red" v-else-if="auditingStatus.isRejected">管理员审核结果：驳回 &nbsp; &nbsp; 驳回理由：{{auditingStatus.rejectedReason}}</span>
						<span class="red" v-else>当前脚本将在审核后生效，请耐心等待审核结果。 <a href="/servers/user-scripts" target="_blank" v-if="isPlus()">去审核 &raquo;</a></span>
					</p>
					<p class="comment" v-if="auditingStatus == null"><span class="green">管理员审核结果：审核通过。</span></p>
					<source-code-box :id="id" type="text/javascript" :read-only="false" @change="changeCode">{{config.code}}</source-code-box>
					<p class="comment">{{comment}}</p>
				</td>
			</tr>
		</tbody>
	</table>
</div>`
})

Vue.component("http-firewall-js-cookie-options-viewer", {
	props: ["v-js-cookie-options"],
	mounted: function () {
		this.updateSummary()
	},
	data: function () {
		let options = this.vJsCookieOptions
		if (options == null) {
			options = {
				life: 0,
				maxFails: 0,
				failBlockTimeout: 0,
				failBlockScopeAll: false,
				scope: ""
			}
		}
		return {
			options: options,
			summary: ""
		}
	},
	methods: {
		updateSummary: function () {
			let summaryList = []
			if (this.options.life > 0) {
				summaryList.push("有效时间" + this.options.life + "秒")
			}
			if (this.options.maxFails > 0) {
				summaryList.push("最多失败" + this.options.maxFails + "次")
			}
			if (this.options.failBlockTimeout > 0) {
				summaryList.push("失败拦截" + this.options.failBlockTimeout + "秒")
			}
			if (this.options.failBlockScopeAll) {
				summaryList.push("尝试全局封禁")
			}

			if (summaryList.length == 0) {
				this.summary = "默认配置"
			} else {
				this.summary = summaryList.join(" / ")
			}
		}
	},
	template: `<div>{{summary}}</div>
`
})

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

Vue.component("http-firewall-captcha-options-viewer", {
	props: ["v-captcha-options"],
	mounted: function () {
		this.updateSummary()
	},
	data: function () {
		let options = this.vCaptchaOptions
		if (options == null) {
			options = {
				life: 0,
				maxFails: 0,
				failBlockTimeout: 0,
				failBlockScopeAll: false,
				uiIsOn: false,
				uiTitle: "",
				uiPrompt: "",
				uiButtonTitle: "",
				uiShowRequestId: false,
				uiCss: "",
				uiFooter: "",
				uiBody: "",
				cookieId: "",
				lang: ""
			}
		}
		return {
			options: options,
			summary: "",
			captchaTypes: window.WAF_CAPTCHA_TYPES
		}
	},
	methods: {
		updateSummary: function () {
			let summaryList = []
			if (this.options.life > 0) {
				summaryList.push("有效时间" + this.options.life + "秒")
			}
			if (this.options.maxFails > 0) {
				summaryList.push("最多失败" + this.options.maxFails + "次")
			}
			if (this.options.failBlockTimeout > 0) {
				summaryList.push("失败拦截" + this.options.failBlockTimeout + "秒")
			}
			if (this.options.failBlockScopeAll) {
				summaryList.push("全局封禁")
			}
			let that = this
			let typeDef = this.captchaTypes.$find(function (k, v) {
				return v.code == that.options.captchaType
			})
			if (typeDef != null) {
				summaryList.push("默认验证方式：" + typeDef.name)
			}

			if (this.options.captchaType == "default") {
				if (this.options.uiIsOn) {
					summaryList.push("定制UI")
				}
			}

			if (this.options.geeTestConfig != null && this.options.geeTestConfig.isOn) {
				summaryList.push("已配置极验")
			}

			if (summaryList.length == 0) {
				this.summary = "默认配置"
			} else {
				this.summary = summaryList.join(" / ")
			}
		}
	},
	template: `<div>{{summary}}</div>
`
})

Vue.component("reverse-proxy-box", {
	props: ["v-reverse-proxy-ref", "v-reverse-proxy-config", "v-is-location", "v-is-group", "v-family"],
	data: function () {
		let reverseProxyRef = this.vReverseProxyRef
		if (reverseProxyRef == null) {
			reverseProxyRef = {
				isPrior: false,
				isOn: false,
				reverseProxyId: 0
			}
		}

		let reverseProxyConfig = this.vReverseProxyConfig
		if (reverseProxyConfig == null) {
			reverseProxyConfig = {
				requestPath: "",
				stripPrefix: "",
				requestURI: "",
				requestHost: "",
				requestHostType: 0,
				requestHostExcludingPort: false,
				addHeaders: [],
				connTimeout: {count: 0, unit: "second"},
				readTimeout: {count: 0, unit: "second"},
				idleTimeout: {count: 0, unit: "second"},
				maxConns: 0,
				maxIdleConns: 0,
				followRedirects: false,
				retry50X: false,
				retry40X: false
			}
		}
		if (reverseProxyConfig.addHeaders == null) {
			reverseProxyConfig.addHeaders = []
		}
		if (reverseProxyConfig.connTimeout == null) {
			reverseProxyConfig.connTimeout = {count: 0, unit: "second"}
		}
		if (reverseProxyConfig.readTimeout == null) {
			reverseProxyConfig.readTimeout = {count: 0, unit: "second"}
		}
		if (reverseProxyConfig.idleTimeout == null) {
			reverseProxyConfig.idleTimeout = {count: 0, unit: "second"}
		}

		if (reverseProxyConfig.proxyProtocol == null) {
			// 如果直接赋值Vue将不会触发变更通知
			Vue.set(reverseProxyConfig, "proxyProtocol", {
				isOn: false,
				version: 1
			})
		}

		let forwardHeaders = [
			{
				name: "X-Real-IP",
				isChecked: false
			},
			{
				name: "X-Forwarded-For",
				isChecked: false
			},
			{
				name: "X-Forwarded-By",
				isChecked: false
			},
			{
				name: "X-Forwarded-Host",
				isChecked: false
			},
			{
				name: "X-Forwarded-Proto",
				isChecked: false
			}
		]
		forwardHeaders.forEach(function (v) {
			v.isChecked = reverseProxyConfig.addHeaders.$contains(v.name)
		})

		return {
			reverseProxyRef: reverseProxyRef,
			reverseProxyConfig: reverseProxyConfig,
			advancedVisible: false,
			family: this.vFamily,
			forwardHeaders: forwardHeaders
		}
	},
	watch: {
		"reverseProxyConfig.requestHostType": function (v) {
			let requestHostType = parseInt(v)
			if (isNaN(requestHostType)) {
				requestHostType = 0
			}
			this.reverseProxyConfig.requestHostType = requestHostType
		},
		"reverseProxyConfig.connTimeout.count": function (v) {
			let count = parseInt(v)
			if (isNaN(count) || count < 0) {
				count = 0
			}
			this.reverseProxyConfig.connTimeout.count = count
		},
		"reverseProxyConfig.readTimeout.count": function (v) {
			let count = parseInt(v)
			if (isNaN(count) || count < 0) {
				count = 0
			}
			this.reverseProxyConfig.readTimeout.count = count
		},
		"reverseProxyConfig.idleTimeout.count": function (v) {
			let count = parseInt(v)
			if (isNaN(count) || count < 0) {
				count = 0
			}
			this.reverseProxyConfig.idleTimeout.count = count
		},
		"reverseProxyConfig.maxConns": function (v) {
			let maxConns = parseInt(v)
			if (isNaN(maxConns) || maxConns < 0) {
				maxConns = 0
			}
			this.reverseProxyConfig.maxConns = maxConns
		},
		"reverseProxyConfig.maxIdleConns": function (v) {
			let maxIdleConns = parseInt(v)
			if (isNaN(maxIdleConns) || maxIdleConns < 0) {
				maxIdleConns = 0
			}
			this.reverseProxyConfig.maxIdleConns = maxIdleConns
		},
		"reverseProxyConfig.proxyProtocol.version": function (v) {
			let version = parseInt(v)
			if (isNaN(version)) {
				version = 1
			}
			this.reverseProxyConfig.proxyProtocol.version = version
		}
	},
	methods: {
		isOn: function () {
			if (this.vIsLocation || this.vIsGroup) {
				return this.reverseProxyRef.isPrior && this.reverseProxyRef.isOn
			}
			return this.reverseProxyRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		},
		changeAddHeader: function () {
			this.reverseProxyConfig.addHeaders = this.forwardHeaders.filter(function (v) {
				return v.isChecked
			}).map(function (v) {
				return v.name
			})
		}
	},
	template: `<div>
	<input type="hidden" name="reverseProxyRefJSON" :value="JSON.stringify(reverseProxyRef)"/>
	<input type="hidden" name="reverseProxyJSON" :value="JSON.stringify(reverseProxyConfig)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="reverseProxyRef" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || reverseProxyRef.isPrior">
			<tr>
				<td class="title">启用源站</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyRef.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后，所有源站设置才会生效。</p>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
				<td>回源主机名<em>（Host）</em></td>
				<td>	
					<radio :v-value="0" v-model="reverseProxyConfig.requestHostType">跟随CDN域名</radio> &nbsp;
					<radio :v-value="1" v-model="reverseProxyConfig.requestHostType">跟随源站</radio> &nbsp;
					<radio :v-value="2" v-model="reverseProxyConfig.requestHostType">自定义</radio>
					<div v-show="reverseProxyConfig.requestHostType == 2" style="margin-top: 0.8em">
						<input type="text" placeholder="比如example.com" v-model="reverseProxyConfig.requestHost"/>
					</div>
					<p class="comment">请求源站时的主机名（Host），用于修改源站接收到的域名
					<span v-if="reverseProxyConfig.requestHostType == 0">，"跟随CDN域名"是指源站接收到的域名和当前CDN访问域名保持一致</span>
					<span v-if="reverseProxyConfig.requestHostType == 1">，"跟随源站"是指源站接收到的域名仍然是填写的源站地址中的信息，不随代理服务域名改变而改变</span>					
					<span v-if="reverseProxyConfig.requestHostType == 2">，自定义Host内容中支持请求变量</span>。</p>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
				<td>回源主机名移除端口</td>
				<td><checkbox v-model="reverseProxyConfig.requestHostExcludingPort"></checkbox>
					<p class="comment">选中后表示移除回源主机名中的端口部分。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && advancedVisible">
			<tr v-show="family == null || family == 'http'">
				<td>回源跟随</td>
				<td>
					<checkbox v-model="reverseProxyConfig.followRedirects"></checkbox>
					<p class="comment">选中后，自动读取源站跳转后的网页内容。</p>
				</td>
			</tr>
		    <tr v-show="family == null || family == 'http'">
		        <td>自动添加报头</td>
		        <td>
		            <div>
		                <div style="width: 14em; float: left; margin-bottom: 1em" v-for="header in forwardHeaders" :key="header.name">
		                    <checkbox v-model="header.isChecked" @input="changeAddHeader">{{header.name}}</checkbox>
                        </div>
                        <div style="clear: both"></div>
                    </div>
                    <p class="comment">选中后，会自动向源站请求添加这些报头，以便于源站获取客户端信息。</p>
                </td> 
            </tr>
			<tr v-show="family == null || family == 'http'">
				<td>请求URI<em>（RequestURI）</em></td>
				<td>
					<input type="text" placeholder="\${requestURI}" v-model="reverseProxyConfig.requestURI"/>
					<p class="comment">\${requestURI}为完整的请求URI，可以使用类似于"\${requestURI}?arg1=value1&arg2=value2"的形式添加你的参数。</p>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
				<td>去除URL前缀<em>（StripPrefix）</em></td>
				<td>
					<input type="text" v-model="reverseProxyConfig.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
			<tr v-if="family == null || family == 'http'">
				<td>自动刷新缓存区<em>（AutoFlush）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyConfig.autoFlush"/>
						<label></label>
					</div>
					<p class="comment">开启后将自动刷新缓冲区数据到客户端，在类似于SSE（server-sent events）等场景下很有用。</p>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
            	<td>自动重试50X</td>
            	<td>
            		<checkbox v-model="reverseProxyConfig.retry50X"></checkbox>
            		<p class="comment">选中后，表示当源站返回状态码为50X（比如502、504等）时，自动重试其他源站。</p>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
            	<td>自动重试40X</td>
            	<td>
            		<checkbox v-model="reverseProxyConfig.retry40X"></checkbox>
            		<p class="comment">选中后，表示当源站返回状态码为40X（403或404）时，自动重试其他源站。</p>
				</td>
			</tr>
            <tr v-show="family != 'unix'">
            	<td>PROXY Protocol</td>
            	<td>
            		<checkbox name="proxyProtocolIsOn" v-model="reverseProxyConfig.proxyProtocol.isOn"></checkbox>
            		<p class="comment">选中后表示启用PROXY Protocol，每次连接源站时都会在头部写入客户端地址信息。</p>
				</td>
			</tr>
			<tr v-show="family != 'unix' && reverseProxyConfig.proxyProtocol.isOn">
				<td>PROXY Protocol版本</td>
				<td>
					<select class="ui dropdown auto-width" name="proxyProtocolVersion" v-model="reverseProxyConfig.proxyProtocol.version">
						<option value="1">1</option>
						<option value="2">2</option>
					</select>
					<p class="comment" v-if="reverseProxyConfig.proxyProtocol.version == 1">发送类似于<code-label>PROXY TCP4 192.168.1.1 192.168.1.10 32567 443</code-label>的头部信息。</p>
					<p class="comment" v-if="reverseProxyConfig.proxyProtocol.version == 2">发送二进制格式的头部信息。</p>
				</td>
			</tr>
			<tr v-if="family == null || family == 'http'">
                <td class="color-border">源站连接失败超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="connTimeout" value="10" size="6" v-model="reverseProxyConfig.connTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">连接源站失败的最大超时时间，0表示不限制。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站读取超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="readTimeout" value="0" size="6" v-model="reverseProxyConfig.readTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">读取内容时的最大超时时间，0表示不限制。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大并发连接数</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="maxConns" value="0" size="6" maxlength="10" v-model="reverseProxyConfig.maxConns"/>
                        </div>
                    </div>
                    <p class="comment">源站可以接受到的最大并发连接数，0表示使用系统默认。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大空闲连接数</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="maxIdleConns" value="0" size="6" maxlength="10" v-model="reverseProxyConfig.maxIdleConns"/>
                        </div>
                    </div>
                    <p class="comment">当没有请求时，源站保持等待的最大空闲连接数量，0表示使用系统默认。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大空闲超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="idleTimeout" value="0" size="6" v-model="reverseProxyConfig.idleTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">源站保持等待的空闲超时时间，0表示使用默认时间。</p>
                </td>
            </tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-firewall-param-filters-box", {
	props: ["v-filters"],
	data: function () {
		let filters = this.vFilters
		if (filters == null) {
			filters = []
		}

		return {
			filters: filters,
			isAdding: false,
			options: [
				{name: "MD5", code: "md5"},
				{name: "URLEncode", code: "urlEncode"},
				{name: "URLDecode", code: "urlDecode"},
				{name: "BASE64Encode", code: "base64Encode"},
				{name: "BASE64Decode", code: "base64Decode"},
				{name: "UNICODE编码", code: "unicodeEncode"},
				{name: "UNICODE解码", code: "unicodeDecode"},
				{name: "HTML实体编码", code: "htmlEscape"},
				{name: "HTML实体解码", code: "htmlUnescape"},
				{name: "计算长度", code: "length"},
				{name: "十六进制->十进制", "code": "hex2dec"},
				{name: "十进制->十六进制", "code": "dec2hex"},
				{name: "SHA1", "code": "sha1"},
				{name: "SHA256", "code": "sha256"}
			],
			addingCode: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.addingCode = ""
		},
		confirm: function () {
			if (this.addingCode.length == 0) {
				return
			}
			let that = this
			this.filters.push(this.options.$find(function (k, v) {
				return (v.code == that.addingCode)
			}))
			this.isAdding = false
		},
		cancel: function () {
			this.isAdding = false
		},
		remove: function (index) {
			this.filters.$remove(index)
		}
	},
	template: `<div>
		<input type="hidden" name="paramFiltersJSON" :value="JSON.stringify(filters)" />
		<div v-if="filters.length > 0">
			<div v-for="(filter, index) in filters" class="ui label small basic">
				{{filter.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a>
			</div>
			<div class="ui divider"></div>
		</div>
		<div v-if="isAdding">
			<div class="ui fields inline">
				<div class="ui field">
					<select class="ui dropdown auto-width" v-model="addingCode">
						<option value="">[请选择]</option>
						<option v-for="option in options" :value="option.code">{{option.name}}</option>
					</select>
				</div>
				<div class="ui field">
					<button class="ui button tiny" type="button" @click.prevent="confirm()">确定</button>
					&nbsp; <a href="" @click.prevent="cancel()" title="取消"><i class="icon remove"></i></a>
				</div>
			</div>
		</div>
		<div v-if="!isAdding">
			<button class="ui button tiny" type="button" @click.prevent="add">+</button>
		</div>
		<p class="comment">可以对参数值进行特定的编解码处理。</p>
</div>`
})

Vue.component("http-remote-addr-config-box", {
	props: ["v-remote-addr-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vRemoteAddrConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				value: "${rawRemoteAddr}",
				type: "default",

				requestHeaderName: ""
			}
		}

		// type
		if (config.type == null || config.type.length == 0) {
			config.type = "default"
			switch (config.value) {
				case "${rawRemoteAddr}":
					config.type = "default"
					break
				case "${remoteAddrValue}":
					config.type = "default"
					break
				case "${remoteAddr}":
					config.type = "proxy"
					break
				default:
					if (config.value != null && config.value.length > 0) {
						config.type = "variable"
					}
			}
		}

		// value
		if (config.value == null || config.value.length == 0) {
			config.value = "${rawRemoteAddr}"
		}

		return {
			config: config,
			options: [
				{
					name: "直接获取",
					description: "用户直接访问边缘节点，即 \"用户 --> 边缘节点\" 模式，这时候系统会试图从直接的连接中读取到客户端IP地址。",
					value: "${rawRemoteAddr}",
					type: "default"
				},
				{
					name: "从上级代理中获取",
					description: "用户和边缘节点之间有别的代理服务转发，即 \"用户 --> [第三方代理服务] --> 边缘节点\"，这时候只能从上级代理中获取传递的IP地址；上级代理传递的请求报头中必须包含 X-Forwarded-For 或 X-Real-IP 信息。",
					value: "${remoteAddr}",
					type: "proxy"
				},
				{
					name: "从请求报头中读取",
					description: "从自定义请求报头读取客户端IP。",
					value: "",
					type: "requestHeader"
				},
				{
					name: "[自定义变量]",
					description: "通过自定义变量来获取客户端真实的IP地址。",
					value: "",
					type: "variable"
				}
			]
		}
	},
	watch: {
		"config.requestHeaderName": function (value) {
			if (this.config.type == "requestHeader"){
				this.config.value = "${header." + value.trim() + "}"
			}
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeOptionType: function () {
			let that = this

			switch(this.config.type) {
				case "default":
					this.config.value = "${rawRemoteAddr}"
					break
				case "proxy":
					this.config.value = "${remoteAddr}"
					break
				case "requestHeader":
					this.config.value = ""
					if (this.requestHeaderName != null && this.requestHeaderName.length > 0) {
						this.config.value = "${header." + this.requestHeaderName + "}"
					}

					setTimeout(function () {
						that.$refs.requestHeaderInput.focus()
					})
					break
				case "variable":
					this.config.value = "${rawRemoteAddr}"

					setTimeout(function () {
						that.$refs.variableInput.focus()
					})

					break
			}
		}
	},
	template: `<div>
	<input type="hidden" name="remoteAddrJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用访客IP设置</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后，表示使用自定义的请求变量获取客户端IP。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>获取IP方式 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="config.type" @change="changeOptionType">
						<option v-for="option in options" :value="option.type">{{option.name}}</option>
					</select>
					<p class="comment" v-for="option in options" v-if="option.type == config.type && option.description.length > 0">{{option.description}}</p>
				</td>
			</tr>
			
			<!-- read from request header -->
			<tr v-show="config.type == 'requestHeader'">
				<td>请求报头 *</td>
				<td>
					<input type="text" name="requestHeaderName" v-model="config.requestHeaderName" maxlength="100" ref="requestHeaderInput"/>
					<p class="comment">请输入包含有客户端IP的请求报头，需要注意大小写，常见的有<code-label>X-Forwarded-For</code-label>、<code-label>X-Real-IP</code-label>、<code-label>X-Client-IP</code-label>等。</p>
				</td>
			</tr>
			
			<!-- read from variable -->
			<tr v-show="config.type == 'variable'">
				<td>读取IP变量值 *</td>
				<td>
					<input type="text" name="value" v-model="config.value" maxlength="100" ref="variableInput"/>
					<p class="comment">通过此变量获取用户的IP地址。具体可用的请求变量列表可参考官方网站文档；比如通过报头传递IP的情形，可以使用<code-label>\${header.你的自定义报头}</code-label>（类似于<code-label>\${header.X-Forwarded-For}</code-label>，需要注意大小写规范）。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>		
</div>`
})

// 访问日志搜索框
Vue.component("http-access-log-search-box", {
	props: ["v-ip", "v-domain", "v-keyword", "v-cluster-id", "v-node-id"],
	data: function () {
		let ip = this.vIp
		if (ip == null) {
			ip = ""
		}

		let domain = this.vDomain
		if (domain == null) {
			domain = ""
		}

		let keyword = this.vKeyword
		if (keyword == null) {
			keyword = ""
		}

		return {
			ip: ip,
			domain: domain,
			keyword: keyword,
			clusterId: this.vClusterId
		}
	},
	methods: {
		cleanIP: function () {
			this.ip = ""
			this.submit()
		},
		cleanDomain: function () {
			this.domain = ""
			this.submit()
		},
		cleanKeyword: function () {
			this.keyword = ""
			this.submit()
		},
		submit: function () {
			let parent = this.$el.parentNode
			while (true) {
				if (parent == null) {
					break
				}
				if (parent.tagName == "FORM") {
					break
				}
				parent = parent.parentNode
			}
			if (parent != null) {
				setTimeout(function () {
					parent.submit()
				}, 500)
			}
		},
		changeCluster: function (clusterId) {
			this.clusterId = clusterId
		}
	},
	template: `<div style="z-index: 10">
	<div class="margin"></div>
	<div class="ui fields inline">
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">IP</span>
				<input type="text" name="ip" placeholder="x.x.x.x" size="15" v-model="ip"/>
				<a class="ui label basic" :class="{disabled: ip.length == 0}" @click.prevent="cleanIP"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small" >
				<span class="ui label basic" style="font-weight: normal">域名</span>
				<input type="text" name="domain" placeholder="example.com" size="15" v-model="domain"/>
				<a class="ui label basic" :class="{disabled: domain.length == 0}" @click.prevent="cleanDomain"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">关键词</span>
				<input type="text" name="keyword" v-model="keyword" placeholder="路径、UserAgent、请求ID等..." size="30"/>
				<a class="ui label basic" :class="{disabled: keyword.length == 0}" @click.prevent="cleanKeyword"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field"><tip-icon content="一些特殊的关键词：<br/>单个状态码：status:200<br/>状态码范围：status:500-504<br/>查询IP：ip:192.168.1.100<br/>查询URL：https://goedge.cloud/docs<br/>查询路径部分：requestPath:/hello/world<br/>查询协议版本：proto:HTTP/1.1<br/>协议：scheme:http<br/>请求方法：method:POST<br/>请求来源：referer:example.com"></tip-icon></div>
	</div>
	<div class="ui fields inline" style="margin-top: 0.5em">
		<div class="ui field">
			<node-cluster-combo-box :v-cluster-id="clusterId" @change="changeCluster"></node-cluster-combo-box>
		</div>
		<div class="ui field" v-if="clusterId > 0">
			<node-combo-box :v-cluster-id="clusterId" :v-node-id="vNodeId"></node-combo-box>
		</div>
		<slot></slot>
		<div class="ui field">
			<button class="ui button small" type="submit">搜索日志</button>
		</div>
	</div>
</div>`
})

Vue.component("server-config-copy-link", {
	props: ["v-server-id", "v-config-code"],
	data: function () {
		return {
			serverId: this.vServerId,
			configCode: this.vConfigCode
		}
	},
	methods: {
		copy: function () {
			teaweb.popup("/servers/server/settings/copy?serverId=" + this.serverId + "&configCode=" + this.configCode, {
				height: "25em",
				callback: function () {
					teaweb.success("批量复制成功")
				}
			})
		}
	},
	template: `<a href=\"" class="item" @click.prevent="copy" style="padding-right:0"><span style="font-size: 0.8em">批量</span>&nbsp;<i class="icon copy small"></i></a>`
})

// 显示指标对象名
Vue.component("metric-key-label", {
	props: ["v-key"],
	data: function () {
		return {
			keyDefs: window.METRIC_HTTP_KEYS
		}
	},
	methods: {
		keyName: function (key) {
			let that = this
			let subKey = ""
			let def = this.keyDefs.$find(function (k, v) {
				if (v.code == key) {
					return true
				}
				if (key.startsWith("${arg.") && v.code.startsWith("${arg.")) {
					subKey = that.getSubKey("arg.", key)
					return true
				}
				if (key.startsWith("${header.") && v.code.startsWith("${header.")) {
					subKey = that.getSubKey("header.", key)
					return true
				}
				if (key.startsWith("${cookie.") && v.code.startsWith("${cookie.")) {
					subKey = that.getSubKey("cookie.", key)
					return true
				}
				return false
			})
			if (def != null) {
				if (subKey.length > 0) {
					return def.name + ": " + subKey
				}
				return def.name
			}
			return key
		},
		getSubKey: function (prefix, key) {
			prefix = "${" + prefix
			let index = key.indexOf(prefix)
			if (index >= 0) {
				key = key.substring(index + prefix.length)
				key = key.substring(0, key.length - 1)
				return key
			}
			return ""
		}
	},
	template: `<div class="ui label basic small">
	{{keyName(this.vKey)}}
</div>`
})

// 指标对象
Vue.component("metric-keys-config-box", {
	props: ["v-keys"],
	data: function () {
		let keys = this.vKeys
		if (keys == null) {
			keys = []
		}
		return {
			keys: keys,
			isAdding: false,
			key: "",
			subKey: "",
			keyDescription: "",

			keyDefs: window.METRIC_HTTP_KEYS
		}
	},
	watch: {
		keys: function () {
			this.$emit("change", this.keys)
		}
	},
	methods: {
		cancel: function () {
			this.key = ""
			this.subKey = ""
			this.keyDescription = ""
			this.isAdding = false
		},
		confirm: function () {
			if (this.key.length == 0) {
				return
			}

			if (this.key.indexOf(".NAME") > 0) {
				if (this.subKey.length == 0) {
					teaweb.warn("请输入参数值")
					return
				}
				this.key = this.key.replace(".NAME", "." + this.subKey)
			}
			this.keys.push(this.key)
			this.cancel()
		},
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				if (that.$refs.key != null) {
					that.$refs.key.focus()
				}
			}, 100)
		},
		remove: function (index) {
			this.keys.$remove(index)
		},
		changeKey: function () {
			if (this.key.length == 0) {
				return
			}
			let that = this
			let def = this.keyDefs.$find(function (k, v) {
				return v.code == that.key
			})
			if (def != null) {
				this.keyDescription = def.description
			}
		},
		keyName: function (key) {
			let that = this
			let subKey = ""
			let def = this.keyDefs.$find(function (k, v) {
				if (v.code == key) {
					return true
				}
				if (key.startsWith("${arg.") && v.code.startsWith("${arg.")) {
					subKey = that.getSubKey("arg.", key)
					return true
				}
				if (key.startsWith("${header.") && v.code.startsWith("${header.")) {
					subKey = that.getSubKey("header.", key)
					return true
				}
				if (key.startsWith("${cookie.") && v.code.startsWith("${cookie.")) {
					subKey = that.getSubKey("cookie.", key)
					return true
				}
				return false
			})
			if (def != null) {
				if (subKey.length > 0) {
					return def.name + ": " + subKey
				}
				return def.name
			}
			return key
		},
		getSubKey: function (prefix, key) {
			prefix = "${" + prefix
			let index = key.indexOf(prefix)
			if (index >= 0) {
				key = key.substring(index + prefix.length)
				key = key.substring(0, key.length - 1)
				return key
			}
			return ""
		}
	},
	template: `<div>
	<input type="hidden" name="keysJSON" :value="JSON.stringify(keys)"/>
	<div>
		<div v-for="(key, index) in keys" class="ui label small basic">
			{{keyName(key)}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="isAdding" style="margin-top: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" v-model="key" @change="changeKey">
					<option value="">[选择对象]</option>
					<option v-for="def in keyDefs" :value="def.code">{{def.name}}</option>
				</select>
			</div>
			<div class="ui field" v-if="key == '\${arg.NAME}'">
				<input type="text" v-model="subKey" placeholder="参数名" size="15"/>
			</div>
			<div class="ui field" v-if="key == '\${header.NAME}'">
				<input type="text" v-model="subKey" placeholder="Header名" size="15">
			</div>
			<div class="ui field" v-if="key == '\${cookie.NAME}'">
				<input type="text" v-model="subKey" placeholder="Cookie名" size="15">
			</div>
			<div class="ui field">
				<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
				<a href="" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment" v-if="keyDescription.length > 0">{{keyDescription}}</p>
	</div>
	<div style="margin-top: 1em" v-if="!isAdding">
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("http-web-root-box", {
	props: ["v-root-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vRootConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				dir: "",
				indexes: [],
				stripPrefix: "",
				decodePath: false,
				isBreak: false,
				exceptHiddenFiles: true,
				onlyURLPatterns: [],
				exceptURLPatterns: []
			}
		}
		if (config.indexes == null) {
			config.indexes = []
		}

		if (config.onlyURLPatterns == null) {
			config.onlyURLPatterns = []
		}
		if (config.exceptURLPatterns == null) {
			config.exceptURLPatterns = []
		}

		return {
			config: config,
			advancedVisible: false
		}
	},
	methods: {
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		},
		addIndex: function () {
			let that = this
			teaweb.popup("/servers/server/settings/web/createIndex", {
				height: "10em",
				callback: function (resp) {
					that.config.indexes.push(resp.data.index)
				}
			})
		},
		removeIndex: function (i) {
			this.config.indexes.$remove(i)
		},
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		}
	},
	template: `<div>
	<input type="hidden" name="rootJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用静态资源分发</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="title">静态资源根目录</td>
				<td>
					<input type="text" name="root" v-model="config.dir" ref="focus" placeholder="类似于 /home/www"/>
					<p class="comment">可以访问此根目录下的静态资源。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>

		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td>首页文件</td>
				<td>
					<!-- TODO 支持排序 -->
					<div v-if="config.indexes.length > 0">
						<div v-for="(index, i) in config.indexes" class="ui label small basic">
							{{index}} <a href="" title="删除" @click.prevent="removeIndex(i)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="addIndex()">+</button>
					<p class="comment">在URL中只有目录没有文件名时默认查找的首页文件。</p>
				</td>
			</tr>
			<tr>
				<td>例外URL</td>
				<td>
					<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了例外URL，表示不支持通过这些URL访问。</p>
				</td>
			</tr>
			<tr>
				<td>限制URL</td>
				<td>
					<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了限制URL，表示仅支持通过这些URL访问。</p>
				</td>
			</tr>	
			<tr>
				<td>排除隐藏文件</td>
				<td>
					<checkbox v-model="config.exceptHiddenFiles"></checkbox>
					<p class="comment">排除以点（.）符号开头的隐藏目录或文件，比如<code-label>/.git/logs/HEAD</code-label></p>
				</td>
			</tr>
			<tr>
				<td>去除URL前缀</td>
				<td>
					<input type="text" v-model="config.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
			<tr>
				<td>路径解码</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.decodePath"/>
						<label></label>	
					</div>
					<p class="comment">是否对请求路径进行URL解码，比如把 <span class="ui label tiny">/Web+App+Browser.html</span> 解码成 <span class="ui label tiny">/Web App Browser.html</span> 再查找文件。</p>
				</td>
			</tr>
			<tr>
				<td>终止请求</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.isBreak"/>
						<label></label>	
					</div>
					<p class="comment">在找不到要访问的文件的情况下是否终止请求并返回404，如果选择终止请求，则不再尝试反向代理等设置。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-webp-config-box", {
	props: ["v-webp-config", "v-is-location", "v-is-group", "v-require-cache"],
	data: function () {
		let config = this.vWebpConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				minLength: {count: 0, "unit": "kb"},
				maxLength: {count: 0, "unit": "kb"},
				mimeTypes: ["image/png", "image/jpeg", "image/bmp", "image/x-ico"],
				extensions: [".png", ".jpeg", ".jpg", ".bmp", ".ico"],
				conds: null
			}
		}

		if (config.mimeTypes == null) {
			config.mimeTypes = []
		}
		if (config.extensions == null) {
			config.extensions = []
		}

		return {
			config: config,
			moreOptionsVisible: false
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeExtensions: function (values) {
			values.forEach(function (v, k) {
				if (v.length > 0 && v[0] != ".") {
					values[k] = "." + v
				}
			})
			this.config.extensions = values
		},
		changeMimeTypes: function (values) {
			this.config.mimeTypes = values
		},
		changeAdvancedVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},
		changeConds: function (conds) {
			this.config.conds = conds
		}
	},
	template: `<div>
	<input type="hidden" name="webpJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用WebP压缩</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示开启自动WebP压缩；图片的宽和高均不能超过16383像素<span v-if="vRequireCache">；只有满足缓存条件的图片内容才会被转换</span>。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被转成WebP，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>支持的MimeType</td>
				<td>
					<values-box :values="config.mimeTypes" @change="changeMimeTypes" placeholder="比如 text/*"></values-box>
					<p class="comment">响应的Content-Type里包含这些MimeType的内容将会被转成WebP。</p>
				</td>
			</tr>
			<tr>
				<td>内容最小长度</td>
				<td>
					<size-capacity-box :v-name="'minLength'" :v-value="config.minLength" :v-unit="'kb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>内容最大长度</td>
				<td>
					<size-capacity-box :v-name="'maxLength'" :v-value="config.maxLength" :v-unit="'mb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>匹配条件</td>
				<td>
					<http-request-conds-box :v-conds="config.conds" @change="changeConds"></http-request-conds-box>
	</td>
			</tr>
		</tbody>
	</table>			
	<div class="ui margin"></div>
</div>`
})

Vue.component("origin-scheduling-view-box", {
	props: ["v-scheduling", "v-params"],
	data: function () {
		let scheduling = this.vScheduling
		if (scheduling == null) {
			scheduling = {}
		}
		return {
			scheduling: scheduling
		}
	},
	methods: {
		update: function () {
			teaweb.popup("/servers/server/settings/reverseProxy/updateSchedulingPopup?" + this.vParams, {
				height: "21em",
				callback: function () {
					window.location.reload()
				},
			})
		}
	},
	template: `<div>
	<div class="margin"></div>
	<table class="ui table selectable definition">
		<tr>
			<td class="title">当前正在使用的算法</td>
			<td>
				{{scheduling.name}} &nbsp; <a href="" @click.prevent="update()"><span>[修改]</span></a>
				<p class="comment">{{scheduling.description}}</p>
			</td>
		</tr>
	</table>
</div>`
})

Vue.component("http-firewall-block-options", {
	props: ["v-block-options"],
	data: function () {
		return {
			options: this.vBlockOptions,
			statusCode: this.vBlockOptions.statusCode,
			timeout: this.vBlockOptions.timeout,
			timeoutMax: this.vBlockOptions.timeoutMax,
			isEditing: false
		}
	},
	watch: {
		statusCode: function (v) {
			let statusCode = parseInt(v)
			if (isNaN(statusCode)) {
				this.options.statusCode = 403
			} else {
				this.options.statusCode = statusCode
			}
		},
		timeout: function (v) {
			let timeout = parseInt(v)
			if (isNaN(timeout)) {
				this.options.timeout = 0
			} else {
				this.options.timeout = timeout
			}
		},
		timeoutMax: function (v) {
			let timeoutMax = parseInt(v)
			if (isNaN(timeoutMax)) {
				this.options.timeoutMax = 0
			} else {
				this.options.timeoutMax = timeoutMax
			}
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	<input type="hidden" name="blockOptionsJSON" :value="JSON.stringify(options)"/>
	<a href="" @click.prevent="edit">状态码：{{statusCode}} / 提示内容：<span v-if="options.body != null && options.body.length > 0">[{{options.body.length}}字符]</span><span v-else class="disabled">[无]</span> <span v-if="timeout > 0"> / 封禁时长：{{timeout}}秒</span>
	 <span v-if="timeoutMax > timeout"> / 最大封禁时长：{{timeoutMax}}秒</span>
	 <span v-if="options.failBlockScopeAll"> / 尝试全局封禁</span>
	 <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<table class="ui table" v-show="isEditing">
		<tr>
			<td class="title">状态码</td>
			<td>
				<input type="text" v-model="statusCode" style="width:4.5em" maxlength="3"/>
			</td>
		</tr>
		<tr>
			<td>提示内容</td>
			<td>
				<textarea rows="3" v-model="options.body"></textarea>
			</td>
		</tr>
		<tr>
			<td>封禁时长</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeout" style="width: 5em" maxlength="6"/>
					<span class="ui label">秒</span>
				</div>
				<p class="comment">触发阻止动作时，封禁客户端IP的时间。</p>
			</td>
		</tr>
		<tr>
			<td>最大封禁时长</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeoutMax" style="width: 5em" maxlength="6"/>
					<span class="ui label">秒</span>
				</div>
				<p class="comment">如果最大封禁时长大于封禁时长（{{timeout}}秒），那么表示每次封禁的时候，将会在这两个时长数字之间随机选取一个数字作为最终的封禁时长。</p>
			</td>
		</tr>
		<tr>
			<td>失败全局封禁</td>
			<td>
				<checkbox v-model="options.failBlockScopeAll"></checkbox>
				<p class="comment">选中后，表示允许系统尝试全局封禁某个IP，以提升封禁性能。</p>
			</td>
		</tr>
	</table>
</div>	
`
})

Vue.component("http-hls-config-box", {
	props: ["value", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isPrior: false
			}
		}

		let encryptingConfig = config.encrypting
		if (encryptingConfig == null) {
			encryptingConfig = {
				isOn: false,
				onlyURLPatterns: [],
				exceptURLPatterns: []
			}
			config.encrypting = encryptingConfig
		}

		return {
			config: config,

			encryptingConfig: encryptingConfig,
			encryptingMoreOptionsVisible: false
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior)
		},

		showEncryptingMoreOptions: function () {
			this.encryptingMoreOptionsVisible = !this.encryptingMoreOptionsVisible
		}
	},
	template: `<div>
	<input type="hidden" name="hlsJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable" v-show="vIsLocation || vIsGroup">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
	</table>
	
	<table class="ui table definition selectable" v-show="isOn()">
		<tbody>
			<tr>
				<td class="title">启用HLS加密</td>
				<td>
					<checkbox v-model="encryptingConfig.isOn"></checkbox>
					<p class="comment">启用后，系统会自动在<code-label>.m3u8</code-label>文件中加入<code-label>#EXT-X-KEY:METHOD=AES-128...</code-label>，并将其中的<code-label>.ts</code-label>文件内容进行加密。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="encryptingConfig.isOn">
			<tr>
				<td colspan="2"><more-options-indicator @change="showEncryptingMoreOptions"></more-options-indicator></td>
			</tr>
		</tbody>
		<tbody v-show="encryptingConfig.isOn && encryptingMoreOptionsVisible">
			<tr>
				<td>例外URL</td>
				<td>
					<url-patterns-box v-model="encryptingConfig.exceptURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
				</td>
			</tr>
			<tr>
				<td>限制URL</td>
				<td>
					<url-patterns-box v-model="encryptingConfig.onlyURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了限制URL，表示只对这些URL进行加密处理；如果不填则表示支持所有的URL。</p>
				</td>
			</tr>	
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-oss-bucket-params", {
	props: ["v-oss-config", "v-params", "name"],
	data: function () {
		let params = this.vParams
		if (params == null) {
			params = []
		}

		let ossConfig = this.vOssConfig
		if (ossConfig == null) {
			ossConfig = {
				bucketParam: "input",
				bucketName: "",
				bucketArgName: ""
			}
		} else {
			// 兼容以往
			if (ossConfig.bucketParam != null && ossConfig.bucketParam.length == 0) {
				ossConfig.bucketParam = "input"
			}
			if (ossConfig.options != null && ossConfig.options.bucketName != null && ossConfig.options.bucketName.length > 0) {
				ossConfig.bucketName = ossConfig.options.bucketName
			}
		}

		return {
			params: params,
			ossConfig: ossConfig
		}
	},
	template: `<tbody>
	<tr>
		<td>{{name}}名称获取方式 *</td>
		<td>
			<select class="ui dropdown auto-width" name="bucketParam" v-model="ossConfig.bucketParam">
				<option v-for="param in params" :value="param.code" v-if="param.example.length == 0">{{param.name.replace("\${optionName}", name)}}</option>
				<option v-for="param in params" :value="param.code" v-if="param.example.length > 0">{{param.name}} - {{param.example}}</option>
			</select>
			<p class="comment" v-for="param in params" v-if="param.code == ossConfig.bucketParam">{{param.description.replace("\${optionName}", name)}}</p>
		</td>
	</tr>
    <tr v-if="ossConfig.bucketParam == 'input'">
        <td>{{name}}名称 *</td>
        <td>
            <input type="text" name="bucketName" maxlength="100" v-model="ossConfig.bucketName"/>
            <p class="comment">{{name}}名称，类似于<code-label>bucket-12345678</code-label>。</p>
        </td>
    </tr>
    <tr v-if="ossConfig.bucketParam == 'arg'">
    	<td>{{name}}参数名称 *</td>
        <td>
            <input type="text" name="bucketArgName" maxlength="100" v-model="ossConfig.bucketArgName"/>
            <p class="comment">{{name}}参数名称，比如<code-label>?myBucketName=BUCKET-NAME</code-label>中的<code-label>myBucketName</code-label>。</p>
        </td>
	</tr>
</tbody>`
})

Vue.component("http-request-scripts-config-box", {
	props: ["vRequestScriptsConfig", "v-auditing-status", "v-is-location"],
	data: function () {
		let config = this.vRequestScriptsConfig
		if (config == null) {
			config = {}
		}

		return {
			config: config
		}
	},
	methods: {
		changeInitGroup: function (group) {
			this.config.initGroup = group
			this.$forceUpdate()
		},
		changeRequestGroup: function (group) {
			this.config.requestGroup = group
			this.$forceUpdate()
		}
	},
	template: `<div>
	<input type="hidden" name="requestScriptsJSON" :value="JSON.stringify(config)"/>
	<div class="margin"></div>
	<h4 style="margin-bottom: 0">请求初始化</h4>
	<p class="comment">在请求刚初始化时调用，此时自定义报头等尚未生效。</p>
	<div>
		<script-group-config-box :v-group="config.initGroup" :v-auditing-status="vAuditingStatus" @change="changeInitGroup" :v-is-location="vIsLocation"></script-group-config-box>
	</div>
	<h4 style="margin-bottom: 0">准备发送请求</h4>
	<p class="comment">在准备执行请求或者转发请求之前调用，此时自定义报头、源站等已准备好。</p>
	<div>
		<script-group-config-box :v-group="config.requestGroup" :v-auditing-status="vAuditingStatus" @change="changeRequestGroup" :v-is-location="vIsLocation"></script-group-config-box>
	</div>
	<div class="margin"></div>
</div>`
})

Vue.component("http-request-cond-view", {
	props: ["v-cond"],
	data: function () {
		return {
			cond: this.vCond,
			components: window.REQUEST_COND_COMPONENTS
		}
	},
	methods: {
		typeName: function (cond) {
			let c = this.components.$find(function (k, v) {
				return v.type == cond.type
			})
			if (c != null) {
				return c.name;
			}
			return cond.param + " " + cond.operator
		},
		updateConds: function (conds, simpleCond) {
			for (let k in simpleCond) {
				if (simpleCond.hasOwnProperty(k)) {
					this.cond[k] = simpleCond[k]
				}
			}
		},
		notifyChange: function () {

		}
	},
	template: `<div style="margin-bottom: 0.5em">
	<span class="ui label small basic">
		<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
		<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{typeName(cond)}}: </var>
		{{cond.value}}
		<sup v-if="cond.isCaseInsensitive" title="不区分大小写"><i class="icon info small"></i></sup>
	</span>
</div>`
})

Vue.component("http-header-assistant", {
	props: ["v-type", "v-value"],
	mounted: function () {
		let that = this
		Tea.action("/servers/headers/options?type=" + this.vType)
			.post()
			.success(function (resp) {
				that.allHeaders = resp.data.headers
			})
	},
	data: function () {
		return {
			allHeaders: [],
			matchedHeaders: [],

			selectedHeaderName: ""
		}
	},
	watch: {
		vValue: function (v) {
			if (v != this.selectedHeaderName) {
				this.selectedHeaderName = ""
			}

			if (v.length == 0) {
				this.matchedHeaders = []
				return
			}
			this.matchedHeaders = this.allHeaders.filter(function (header) {
				return teaweb.match(header, v)
			}).slice(0, 10)
		}
	},
	methods: {
		select: function (header) {
			this.$emit("select", header)
			this.selectedHeaderName = header
		}
	},
	template: `<span v-if="selectedHeaderName.length == 0">
	<a href="" v-for="header in matchedHeaders" class="ui label basic tiny blue" style="font-weight: normal; margin-bottom: 0.3em" @click.prevent="select(header)">{{header}}</a>
	<span v-if="matchedHeaders.length > 0">&nbsp; &nbsp;</span>
</span>`
})

Vue.component("http-firewall-rules-box", {
	props: ["v-rules", "v-type"],
	data: function () {
		let rules = this.vRules
		if (rules == null) {
			rules = []
		}
		return {
			rules: rules
		}
	},
	methods: {
		addRule: function () {
			window.UPDATING_RULE = null
			let that = this
			teaweb.popup("/servers/components/waf/createRulePopup?type=" + this.vType, {
				height: "30em",
				callback: function (resp) {
					that.rules.push(resp.data.rule)
				}
			})
		},
		updateRule: function (index, rule) {
			window.UPDATING_RULE = teaweb.clone(rule)
			let that = this
			teaweb.popup("/servers/components/waf/createRulePopup?type=" + this.vType, {
				height: "30em",
				callback: function (resp) {
					Vue.set(that.rules, index, resp.data.rule)
				}
			})
		},
		removeRule: function (index) {
			let that = this
			teaweb.confirm("确定要删除此规则吗？", function () {
				that.rules.$remove(index)
			})
		},
		operatorName: function (operatorCode) {
			let operatorName = operatorCode
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorName = v.name
					}
				})
			}

			return operatorName
		},
		operatorDescription: function (operatorCode) {
			let operatorName = operatorCode
			let operatorDescription = ""
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorName = v.name
						operatorDescription = v.description
					}
				})
			}

			return operatorName + ": " + operatorDescription
		},
		operatorDataType: function (operatorCode) {
			let operatorDataType = "none"
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorDataType = v.dataType
					}
				})
			}

			return operatorDataType
		},
		calculateParamName: function (param) {
			let paramName = ""
			if (param != null) {
				window.WAF_RULE_CHECKPOINTS.forEach(function (checkpoint) {
					if (param == "${" + checkpoint.prefix + "}" || param.startsWith("${" + checkpoint.prefix + ".")) {
						paramName = checkpoint.name
					}
				})
			}
			return paramName
		},
		calculateParamDescription: function (param) {
			let paramName = ""
			let paramDescription = ""
			if (param != null) {
				window.WAF_RULE_CHECKPOINTS.forEach(function (checkpoint) {
					if (param == "${" + checkpoint.prefix + "}" || param.startsWith("${" + checkpoint.prefix + ".")) {
						paramName = checkpoint.name
						paramDescription = checkpoint.description
					}
				})
			}
			return paramName + ": " + paramDescription
		},
		isEmptyString: function (v) {
			return typeof v == "string" && v.length == 0
		}
	},
	template: `<div>
		<input type="hidden" name="rulesJSON" :value="JSON.stringify(rules)"/>
		<div v-if="rules.length > 0">
			<div v-for="(rule, index) in rules" class="ui label small basic" style="margin-bottom: 0.5em; line-height: 1.5">
				{{rule.name}} <span :title="calculateParamDescription(rule.param)" class="hover">{{calculateParamName(rule.param)}}<span class="small grey"> {{rule.param}}</span></span>
				
				<!-- cc2 -->
				<span v-if="rule.param == '\${cc2}'">
					{{rule.checkpointOptions.period}}秒内请求数
				</span>	
				
				<!-- refererBlock -->
				<span v-if="rule.param == '\${refererBlock}'">
					<span v-if="rule.checkpointOptions.allowDomains != null && rule.checkpointOptions.allowDomains.length > 0">允许{{rule.checkpointOptions.allowDomains}}</span>
					<span v-if="rule.checkpointOptions.denyDomains != null && rule.checkpointOptions.denyDomains.length > 0">禁止{{rule.checkpointOptions.denyDomains}}</span>
				</span>
				
				<span v-else>
					<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> <span class="hover" :title="operatorDescription(rule.operator) + ((!rule.isComposed && rule.isCaseInsensitive) ? '\\n[大小写不敏感] ':'')">&lt;{{operatorName(rule.operator)}}&gt;</span> 
						<span v-if="!isEmptyString(rule.value)" class="hover">{{rule.value}}</span>
						<span v-else-if="operatorDataType(rule.operator) != 'none'" class="disabled" style="font-weight: normal" title="空字符串">[空]</span>
				</span>
				
				<!-- description -->
				<span v-if="rule.description != null && rule.description.length > 0" class="grey small">（{{rule.description}}）</span>
				
				<a href="" title="修改" @click.prevent="updateRule(index, rule)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeRule(index)"><i class="icon remove"></i></a>
			</div>
			<div class="ui divider"></div>
		</div>
		<button class="ui button tiny" type="button" @click.prevent="addRule()">+</button>
</div>`
})

Vue.component("http-fastcgi-box", {
	props: ["v-fastcgi-ref", "v-fastcgi-configs", "v-is-location"],
	data: function () {
		let fastcgiRef = this.vFastcgiRef
		if (fastcgiRef == null) {
			fastcgiRef = {
				isPrior: false,
				isOn: false,
				fastcgiIds: []
			}
		}
		let fastcgiConfigs = this.vFastcgiConfigs
		if (fastcgiConfigs == null) {
			fastcgiConfigs = []
		} else {
			fastcgiRef.fastcgiIds = fastcgiConfigs.map(function (v) {
				return v.id
			})
		}

		return {
			fastcgiRef: fastcgiRef,
			fastcgiConfigs: fastcgiConfigs,
			advancedVisible: false
		}
	},
	methods: {
		isOn: function () {
			return (!this.vIsLocation || this.fastcgiRef.isPrior) && this.fastcgiRef.isOn
		},
		createFastcgi: function () {
			let that = this
			teaweb.popup("/servers/server/settings/fastcgi/createPopup", {
				height: "26em",
				callback: function (resp) {
					teaweb.success("添加成功", function () {
						that.fastcgiConfigs.push(resp.data.fastcgi)
						that.fastcgiRef.fastcgiIds.push(resp.data.fastcgi.id)
					})
				}
			})
		},
		updateFastcgi: function (fastcgiId, index) {
			let that = this
			teaweb.popup("/servers/server/settings/fastcgi/updatePopup?fastcgiId=" + fastcgiId, {
				callback: function (resp) {
					teaweb.success("修改成功", function () {
						Vue.set(that.fastcgiConfigs, index, resp.data.fastcgi)
					})
				}
			})
		},
		removeFastcgi: function (index) {
			this.fastcgiRef.fastcgiIds.$remove(index)
			this.fastcgiConfigs.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="fastcgiRefJSON" :value="JSON.stringify(fastcgiRef)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="fastcgiRef" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="(!this.vIsLocation || this.fastcgiRef.isPrior)">
			<tr>
				<td class="title">启用配置</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="fastcgiRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-if="isOn()">
			<tr>
				<td>Fastcgi服务</td>
				<td>
					<div v-show="fastcgiConfigs.length > 0" style="margin-bottom: 0.5em">
						<div class="ui label basic small" :class="{disabled: !fastcgi.isOn}" v-for="(fastcgi, index) in fastcgiConfigs">
							{{fastcgi.address}} &nbsp; <a href="" title="修改" @click.prevent="updateFastcgi(fastcgi.id, index)"><i class="ui icon pencil small"></i></a> &nbsp; <a href="" title="删除" @click.prevent="removeFastcgi(index)"><i class="ui icon remove"></i></a>
						</div>
						<div class="ui divided"></div>
					</div>
					<button type="button" class="ui button tiny" @click.prevent="createFastcgi()">+</button>
				</td>
			</tr>
		</tbody>
	</table>	
	<div class="margin"></div>
</div>`
})

// 请求方法列表
Vue.component("http-methods-box", {
	props: ["v-methods"],
	data: function () {
		let methods = this.vMethods
		if (methods == null) {
			methods = []
		}
		return {
			methods: methods,
			isAdding: false,
			addingMethod: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingMethod.focus()
			}, 100)
		},
		confirm: function () {
			let that = this

			// 删除其中的空格
			this.addingMethod = this.addingMethod.replace(/\s/g, "").toUpperCase()

			if (this.addingMethod.length == 0) {
				teaweb.warn("请输入要添加的请求方法", function () {
					that.$refs.addingMethod.focus()
				})
				return
			}

			// 是否已经存在
			if (this.methods.$contains(this.addingMethod)) {
				teaweb.warn("此请求方法已经存在，无需重复添加", function () {
					that.$refs.addingMethod.focus()
				})
				return
			}

			this.methods.push(this.addingMethod)
			this.cancel()
		},
		remove: function (index) {
			this.methods.$remove(index)
		},
		cancel: function () {
			this.isAdding = false
			this.addingMethod = ""
		}
	},
	template: `<div>
	<input type="hidden" name="methodsJSON" :value="JSON.stringify(methods)"/>
	<div v-if="methods.length > 0">
		<span class="ui label small basic" v-for="(method, index) in methods">
			{{method}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingMethod" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingMethod" placeholder="如GET" size="10"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">格式为大写，比如<code-label>GET</code-label>、<code-label>POST</code-label>等。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

// URL扩展名条件
Vue.component("http-cond-url-extension", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPathLowerExtension}",
			operator: "in",
			value: "[]"
		}
		if (this.vCond != null && this.vCond.param == cond.param) {
			cond.value = this.vCond.value
		}

		let extensions = []
		try {
			extensions = JSON.parse(cond.value)
		} catch (e) {

		}

		return {
			cond: cond,
			extensions: extensions, // TODO 可以拖动排序

			isAdding: false,
			addingExt: ""
		}
	},
	watch: {
		extensions: function () {
			this.cond.value = JSON.stringify(this.extensions)
		}
	},
	methods: {
		addExt: function () {
			this.isAdding = !this.isAdding

			if (this.isAdding) {
				let that = this
				setTimeout(function () {
					that.$refs.addingExt.focus()
				}, 100)
			}
		},
		cancelAdding: function () {
			this.isAdding = false
			this.addingExt = ""
		},
		confirmAdding: function () {
			// TODO 做更详细的校验
			// TODO 如果有重复的则提示之

			if (this.addingExt.length == 0) {
				return
			}

			let that = this
			this.addingExt.split(/[,;，；|]/).forEach(function (ext) {
				ext = ext.trim()
				if (ext.length > 0) {
					if (ext[0] != ".") {
						ext = "." + ext
					}
					ext = ext.replace(/\s+/g, "").toLowerCase()
					that.extensions.push(ext)
				}
			})

			// 清除状态
			this.cancelAdding()
		},
		removeExt: function (index) {
			this.extensions.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="extensions.length > 0">
		<div class="ui label small basic" v-for="(ext, index) in extensions">{{ext}} <a href="" title="删除" @click.prevent="removeExt(index)"><i class="icon remove small"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="20" maxlength="100" v-model="addingExt" ref="addingExt" placeholder=".xxx, .yyy" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny basic" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em" v-show="!isAdding">
		<button class="ui button tiny basic" type="button" @click.prevent="addExt()">+添加扩展名</button>
	</div>
	<p class="comment">扩展名需要包含点（.）符号，例如<code-label>.jpg</code-label>、<code-label>.png</code-label>之类；多个扩展名用逗号分割。</p>
</div>`
})

// 排除URL扩展名条件
Vue.component("http-cond-url-not-extension", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPathLowerExtension}",
			operator: "not in",
			value: "[]"
		}
		if (this.vCond != null && this.vCond.param == cond.param) {
			cond.value = this.vCond.value
		}

		let extensions = []
		try {
			extensions = JSON.parse(cond.value)
		} catch (e) {

		}

		return {
			cond: cond,
			extensions: extensions, // TODO 可以拖动排序

			isAdding: false,
			addingExt: ""
		}
	},
	watch: {
		extensions: function () {
			this.cond.value = JSON.stringify(this.extensions)
		}
	},
	methods: {
		addExt: function () {
			this.isAdding = !this.isAdding

			if (this.isAdding) {
				let that = this
				setTimeout(function () {
					that.$refs.addingExt.focus()
				}, 100)
			}
		},
		cancelAdding: function () {
			this.isAdding = false
			this.addingExt = ""
		},
		confirmAdding: function () {
			// TODO 做更详细的校验
			// TODO 如果有重复的则提示之

			if (this.addingExt.length == 0) {
				return
			}
			if (this.addingExt[0] != ".") {
				this.addingExt = "." + this.addingExt
			}
			this.addingExt = this.addingExt.replace(/\s+/g, "").toLowerCase()
			this.extensions.push(this.addingExt)

			// 清除状态
			this.cancelAdding()
		},
		removeExt: function (index) {
			this.extensions.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="extensions.length > 0">
		<div class="ui label small basic" v-for="(ext, index) in extensions">{{ext}} <a href="" title="删除" @click.prevent="removeExt(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="6" maxlength="100" v-model="addingExt" ref="addingExt" placeholder=".xxx" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny basic" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em" v-show="!isAdding">
		<button class="ui button tiny basic" type="button" @click.prevent="addExt()">+添加扩展名</button>
	</div>
	<p class="comment">扩展名需要包含点（.）符号，例如<code-label>.jpg</code-label>、<code-label>.png</code-label>之类。</p>
</div>`
})

// 根据URL前缀
Vue.component("http-cond-url-prefix", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "prefix",
			value: "",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof (this.vCond.value) == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment">URL前缀，有此前缀的URL都将会被匹配，通常以<code-label>/</code-label>开头，比如<code-label>/static</code-label>、<code-label>/images</code-label>，不需要带域名。</p>
</div>`
})

Vue.component("http-cond-url-not-prefix", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "prefix",
			value: "",
			isReverse: true,
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment">要排除的URL前缀，有此前缀的URL都将会被匹配，通常以<code-label>/</code-label>开头，比如<code-label>/static</code-label>、<code-label>/images</code-label>，不需要带域名。</p>
</div>`
})

// 首页
Vue.component("http-cond-url-eq-index", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "eq",
			value: "/",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" disabled="disabled" style="background: #eee"/>
	<p class="comment">检查URL路径是为<code-label>/</code-label>，不需要带域名。</p>
</div>`
})

// 全站
Vue.component("http-cond-url-all", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "prefix",
			value: "/",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" disabled="disabled" style="background: #eee"/>
	<p class="comment">支持全站所有URL。</p>
</div>`
})

// URL精准匹配
Vue.component("http-cond-url-eq", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "eq",
			value: "",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment">完整的URL路径，通常以<code-label>/</code-label>开头，比如<code-label>/static/ui.js</code-label>，不需要带域名。</p>
</div>`
})

Vue.component("http-cond-url-not-eq", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "eq",
			value: "",
			isReverse: true,
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment">要排除的完整的URL路径，通常以<code-label>/</code-label>开头，比如<code-label>/static/ui.js</code-label>，不需要带域名。</p>
</div>`
})

// URL正则匹配
Vue.component("http-cond-url-regexp", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "regexp",
			value: "",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment">匹配URL的正则表达式，比如<code-label>^/static/(.*).js$</code-label>，不需要带域名。</p>
</div>`
})

// 排除URL正则匹配
Vue.component("http-cond-url-not-regexp", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "not regexp",
			value: "",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment"><strong>不要</strong>匹配URL的正则表达式，意即只要匹配成功则排除此条件，比如<code-label>^/static/(.*).js$</code-label>，不需要带域名。</p>
</div>`
})

// URL通配符
Vue.component("http-cond-url-wildcard-match", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "wildcard match",
			value: "",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment">匹配URL的通配符，用星号（<code-label>*</code-label>）表示任意字符，比如（<code-label>/images/*.png</code-label>、<code-label>/static/*</code-label>，不需要带域名。</p>
</div>`
})

// User-Agent正则匹配
Vue.component("http-cond-user-agent-regexp", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${userAgent}",
			operator: "regexp",
			value: "",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment">匹配User-Agent的正则表达式，比如<code-label>Android|iPhone</code-label>。</p>
</div>`
})

// User-Agent正则不匹配
Vue.component("http-cond-user-agent-not-regexp", {
	props: ["v-cond"],
	mounted: function () {
		this.$refs.valueInput.focus()
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "${userAgent}",
			operator: "not regexp",
			value: "",
			isCaseInsensitive: false
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	methods: {
		changeCaseInsensitive: function (isCaseInsensitive) {
			this.cond.isCaseInsensitive = isCaseInsensitive
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value" ref="valueInput"/>
	<p class="comment">匹配User-Agent的正则表达式，比如<code-label>Android|iPhone</code-label>，如果匹配，则排除此条件。</p>
</div>`
})

// 根据MimeType
Vue.component("http-cond-mime-type", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: false,
			param: "${response.contentType}",
			operator: "mime type",
			value: "[]"
		}
		if (this.vCond != null && this.vCond.param == cond.param) {
			cond.value = this.vCond.value
		}
		return {
			cond: cond,
			mimeTypes: JSON.parse(cond.value), // TODO 可以拖动排序

			isAdding: false,
			addingMimeType: ""
		}
	},
	watch: {
		mimeTypes: function () {
			this.cond.value = JSON.stringify(this.mimeTypes)
		}
	},
	methods: {
		addMimeType: function () {
			this.isAdding = !this.isAdding

			if (this.isAdding) {
				let that = this
				setTimeout(function () {
					that.$refs.addingMimeType.focus()
				}, 100)
			}
		},
		cancelAdding: function () {
			this.isAdding = false
			this.addingMimeType = ""
		},
		confirmAdding: function () {
			// TODO 做更详细的校验
			// TODO 如果有重复的则提示之

			if (this.addingMimeType.length == 0) {
				return
			}
			this.addingMimeType = this.addingMimeType.replace(/\s+/g, "")
			this.mimeTypes.push(this.addingMimeType)

			// 清除状态
			this.cancelAdding()
		},
		removeMimeType: function (index) {
			this.mimeTypes.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="mimeTypes.length > 0">
		<div class="ui label small" v-for="(mimeType, index) in mimeTypes">{{mimeType}} <a href="" title="删除" @click.prevent="removeMimeType(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="16" maxlength="100" v-model="addingMimeType" ref="addingMimeType" placeholder="类似于image/png" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny basic" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em">
		<button class="ui button tiny basic" type="button" @click.prevent="addMimeType()">+添加MimeType</button>
	</div>
	<p class="comment">服务器返回的内容的MimeType，比如<span class="ui label tiny">text/html</span>、<span class="ui label tiny">image/*</span>等。</p>
</div>`
})

// 参数匹配
Vue.component("http-cond-params", {
	props: ["v-cond"],
	mounted: function () {
		let cond = this.vCond
		if (cond == null) {
			return
		}
		this.operator = cond.operator

		// stringValue
		if (["regexp", "not regexp", "eq", "not", "prefix", "suffix", "contains", "not contains", "eq ip", "gt ip", "gte ip", "lt ip", "lte ip", "ip range"].$contains(cond.operator)) {
			this.stringValue = cond.value
			return
		}

		// numberValue
		if (["eq int", "eq float", "gt", "gte", "lt", "lte", "mod 10", "ip mod 10", "mod 100", "ip mod 100"].$contains(cond.operator)) {
			this.numberValue = cond.value
			return
		}

		// modValue
		if (["mod", "ip mod"].$contains(cond.operator)) {
			let pieces = cond.value.split(",")
			this.modDivValue = pieces[0]
			if (pieces.length > 1) {
				this.modRemValue = pieces[1]
			}
			return
		}

		// stringValues
		let that = this
		if (["in", "not in", "file ext", "mime type"].$contains(cond.operator)) {
			try {
				let arr = JSON.parse(cond.value)
				if (arr != null && (arr instanceof Array)) {
					arr.forEach(function (v) {
						that.stringValues.push(v)
					})
				}
			} catch (e) {

			}
			return
		}

		// versionValue
		if (["version range"].$contains(cond.operator)) {
			let pieces = cond.value.split(",")
			this.versionRangeMinValue = pieces[0]
			if (pieces.length > 1) {
				this.versionRangeMaxValue = pieces[1]
			}
			return
		}
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "",
			operator: window.REQUEST_COND_OPERATORS[0].op,
			value: "",
			isCaseInsensitive: false
		}
		if (this.vCond != null) {
			cond = this.vCond
		}
		return {
			cond: cond,
			operators: window.REQUEST_COND_OPERATORS,
			operator: window.REQUEST_COND_OPERATORS[0].op,
			operatorDescription: window.REQUEST_COND_OPERATORS[0].description,
			variables: window.REQUEST_VARIABLES,
			variable: "",

			// 各种类型的值
			stringValue: "",
			numberValue: "",

			modDivValue: "",
			modRemValue: "",

			stringValues: [],

			versionRangeMinValue: "",
			versionRangeMaxValue: ""
		}
	},
	methods: {
		changeVariable: function () {
			let v = this.cond.param
			if (v == null) {
				v = ""
			}
			this.cond.param = v + this.variable
		},
		changeOperator: function () {
			let that = this
			this.operators.forEach(function (v) {
				if (v.op == that.operator) {
					that.operatorDescription = v.description
				}
			})

			this.cond.operator = this.operator

			// 移动光标
			let box = document.getElementById("variables-value-box")
			if (box != null) {
				setTimeout(function () {
					let input = box.getElementsByTagName("INPUT")
					if (input.length > 0) {
						input[0].focus()
					}
				}, 100)
			}
		},
		changeStringValues: function (v) {
			this.stringValues = v
			this.cond.value = JSON.stringify(v)
		}
	},
	watch: {
		stringValue: function (v) {
			this.cond.value = v
		},
		numberValue: function (v) {
			// TODO 校验数字
			this.cond.value = v
		},
		modDivValue: function (v) {
			if (v.length == 0) {
				return
			}
			let div = parseInt(v)
			if (isNaN(div)) {
				div = 1
			}
			this.modDivValue = div
			this.cond.value = div + "," + this.modRemValue
		},
		modRemValue: function (v) {
			if (v.length == 0) {
				return
			}
			let rem = parseInt(v)
			if (isNaN(rem)) {
				rem = 0
			}
			this.modRemValue = rem
			this.cond.value = this.modDivValue + "," + rem
		},
		versionRangeMinValue: function (v) {
			this.cond.value = this.versionRangeMinValue + "," + this.versionRangeMaxValue
		},
		versionRangeMaxValue: function (v) {
			this.cond.value = this.versionRangeMinValue + "," + this.versionRangeMaxValue
		}
	},
	template: `<tbody>
	<tr>
		<td style="width: 8em">参数值</td>
		<td>
			<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
			<div>
				<div class="ui field">
					<input type="text" placeholder="\${xxx}" v-model="cond.param"/>
				</div>
				<div class="ui field">
					<select class="ui dropdown" style="width: 16em; color: grey" v-model="variable" @change="changeVariable">
						<option value="">[常用参数]</option>
						<option v-for="v in variables" :value="v.code">{{v.code}} - {{v.name}}</option>
					</select>
				</div>
			</div>
			<p class="comment">其中可以使用变量，类似于<code-label>\${requestPath}</code-label>，也可以是多个变量的组合。</p>
		</td>
	</tr>
	<tr>
		<td>操作符</td>
		<td>
			<div>
				<select class="ui dropdown auto-width" v-model="operator" @change="changeOperator">
					<option v-for="operator in operators" :value="operator.op">{{operator.name}}</option>
				</select>
				<p class="comment" v-html="operatorDescription"></p>
			</div>
		</td>
	</tr>
	<tr v-show="!['file exist', 'file not exist'].$contains(cond.operator)">
		<td>对比值</td>
		<td id="variables-value-box">
			<!-- 正则表达式 -->
			<div v-if="['regexp', 'not regexp'].$contains(cond.operator)">
				<input type="text" v-model="stringValue"/>
				<p class="comment">要匹配的正则表达式，比如<code-label>^/static/(.+).js</code-label>。</p>
			</div>
			
			<!-- 数字相关 -->
			<div v-if="['eq int', 'eq float', 'gt', 'gte', 'lt', 'lte'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">要对比的数字。</p>
			</div>
			
			<!-- 取模 -->
			<div v-if="['mod 10'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数值除以10的余数，在0-9之间。</p>
			</div>
			<div v-if="['mod 100'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数值除以100的余数，在0-99之间。</p>
			</div>
			<div v-if="['mod', 'ip mod'].$contains(cond.operator)">
				<div class="ui fields inline">
					<div class="ui field">除：</div>
					<div class="ui field">
						<input type="text" maxlength="11" size="11" style="width: 5em" v-model="modDivValue" placeholder="除数"/>
					</div>
					<div class="ui field">余：</div>
					<div class="ui field">
						<input type="text" maxlength="11" size="11" style="width: 5em" v-model="modRemValue" placeholder="余数"/>
					</div>
				</div>
			</div>
			
			<!-- 字符串相关 -->
			<div v-if="['eq', 'not', 'prefix', 'suffix', 'contains', 'not contains'].$contains(cond.operator)">
				<input type="text" v-model="stringValue"/>
				<p class="comment" v-if="cond.operator == 'eq'">和参数值一致的字符串。</p>
				<p class="comment" v-if="cond.operator == 'not'">和参数值不一致的字符串。</p>
				<p class="comment" v-if="cond.operator == 'prefix'">参数值的前缀。</p>
				<p class="comment" v-if="cond.operator == 'suffix'">参数值的后缀为此字符串。</p>
				<p class="comment" v-if="cond.operator == 'contains'">参数值包含此字符串。</p>
				<p class="comment" v-if="cond.operator == 'not contains'">参数值不包含此字符串。</p>
			</div>
			<div v-if="['in', 'not in', 'file ext', 'mime type'].$contains(cond.operator)">
				<values-box @change="changeStringValues" :values="stringValues" size="15"></values-box>
				<p class="comment" v-if="cond.operator == 'in'">添加参数值列表。</p>
				<p class="comment" v-if="cond.operator == 'not in'">添加参数值列表。</p>
				<p class="comment" v-if="cond.operator == 'file ext'">添加扩展名列表，比如<code-label>png</code-label>、<code-label>html</code-label>，不包括点。</p>
				<p class="comment" v-if="cond.operator == 'mime type'">添加MimeType列表，类似于<code-label>text/html</code-label>、<code-label>image/*</code-label>。</p>
			</div>
			<div v-if="['version range'].$contains(cond.operator)">
				<div class="ui fields inline">
					<div class="ui field"><input type="text" v-model="versionRangeMinValue" maxlength="200" placeholder="最小版本" style="width: 10em"/></div>
					<div class="ui field">-</div>
					<div class="ui field"><input type="text" v-model="versionRangeMaxValue" maxlength="200" placeholder="最大版本" style="width: 10em"/></div>
				</div>
			</div>
			
			<!-- IP相关 -->
			<div v-if="['eq ip', 'gt ip', 'gte ip', 'lt ip', 'lte ip', 'ip range'].$contains(cond.operator)">
				<input type="text" style="width: 10em" v-model="stringValue" placeholder="x.x.x.x"/>
				<p class="comment">要对比的IP。</p>
			</div>
			<div v-if="['ip mod 10'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数中IP转换成整数后除以10的余数，在0-9之间。</p>
			</div>
			<div v-if="['ip mod 100'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数中IP转换成整数后除以100的余数，在0-99之间。</p>
			</div>
		</td>
	</tr>
	<tr v-if="['regexp', 'not regexp', 'eq', 'not', 'prefix', 'suffix', 'contains', 'not contains', 'in', 'not in'].$contains(cond.operator)">
		<td>不区分大小写</td>
		<td>
		   <div class="ui checkbox">
				<input type="checkbox" name="condIsCaseInsensitive" v-model="cond.isCaseInsensitive"/>
				<label></label>
			</div>
			<p class="comment">选中后表示对比时忽略参数值的大小写。</p>
		</td>
	</tr>
</tbody>
`
})

// 请求方法列表
Vue.component("http-status-box", {
	props: ["v-status-list"],
	data: function () {
		let statusList = this.vStatusList
		if (statusList == null) {
			statusList = []
		}
		return {
			statusList: statusList,
			isAdding: false,
			addingStatus: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingStatus.focus()
			}, 100)
		},
		confirm: function () {
			let that = this

			// 删除其中的空格
			this.addingStatus = this.addingStatus.replace(/\s/g, "").toUpperCase()

			if (this.addingStatus.length == 0) {
				teaweb.warn("请输入要添加的状态码", function () {
					that.$refs.addingStatus.focus()
				})
				return
			}

			// 是否已经存在
			if (this.statusList.$contains(this.addingStatus)) {
				teaweb.warn("此状态码已经存在，无需重复添加", function () {
					that.$refs.addingStatus.focus()
				})
				return
			}

			// 格式
			if (!this.addingStatus.match(/^\d{3}$/)) {
				teaweb.warn("请输入正确的状态码", function () {
					that.$refs.addingStatus.focus()
				})
				return
			}

			this.statusList.push(parseInt(this.addingStatus, 10))
			this.cancel()
		},
		remove: function (index) {
			this.statusList.$remove(index)
		},
		cancel: function () {
			this.isAdding = false
			this.addingStatus = ""
		}
	},
	template: `<div>
	<input type="hidden" name="statusListJSON" :value="JSON.stringify(statusList)"/>
	<div v-if="statusList.length > 0">
		<span class="ui label small basic" v-for="(status, index) in statusList">
			{{status}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingStatus" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingStatus" placeholder="如200" size="3" maxlength="3" style="width: 5em"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">格式为三位数字，比如<code-label>200</code-label>、<code-label>404</code-label>等。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("server-group-selector", {
	props: ["v-groups"],
	data: function () {
		let groups = this.vGroups
		if (groups == null) {
			groups = []
		}
		return {
			groups: groups
		}
	},
	methods: {
		selectGroup: function () {
			let that = this
			let groupIds = this.groups.map(function (v) {
				return v.id.toString()
			}).join(",")
			teaweb.popup("/servers/groups/selectPopup?selectedGroupIds=" + groupIds, {
				callback: function (resp) {
					that.groups.push(resp.data.group)
				}
			})
		},
		addGroup: function () {
			let that = this
			teaweb.popup("/servers/groups/createPopup", {
				callback: function (resp) {
					that.groups.push(resp.data.group)
				}
			})
		},
		removeGroup: function (index) {
			this.groups.$remove(index)
		},
		groupIds: function () {
			return this.groups.map(function (v) {
				return v.id
			})
		}
	},
	template: `<div>
	<div v-if="groups.length > 0">
		<div class="ui label small basic" v-if="groups.length > 0" v-for="(group, index) in groups">
			<input type="hidden" name="groupIds" :value="group.id"/>
			{{group.name}} &nbsp;<a href="" title="删除" @click.prevent="removeGroup(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<a href="" @click.prevent="selectGroup()">[选择分组]</a> &nbsp; <a href="" @click.prevent="addGroup()">[添加分组]</a>
	</div>
</div>`
})

Vue.component("script-group-config-box", {
	props: ["v-group", "v-auditing-status", "v-is-location"],
	data: function () {
		let group = this.vGroup
		if (group == null) {
			group = {
				isPrior: false,
				isOn: true,
				scripts: []
			}
		}
		if (group.scripts == null) {
			group.scripts = []
		}

		let script = null
		if (group.scripts.length > 0) {
			script = group.scripts[group.scripts.length - 1]
		}

		return {
			group: group,
			script: script
		}
	},
	methods: {
		changeScript: function (script) {
			this.group.scripts = [script] // 目前只支持单个脚本
			this.change()
		},
		change: function () {
			this.$emit("change", this.group)
		}
	},
	template: `<div>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="group" v-if="vIsLocation"></prior-checkbox>
		</table>
		<div :style="{opacity: (!vIsLocation || group.isPrior) ? 1 : 0.5}">
			<script-config-box :v-script-config="script" :v-auditing-status="vAuditingStatus" comment="在接收到客户端请求之后立即调用。预置req、resp变量。" @change="changeScript" :v-is-location="vIsLocation"></script-config-box>
		</div>
</div>`
})

// 指标周期设置
Vue.component("metric-period-config-box", {
	props: ["v-period", "v-period-unit"],
	data: function () {
		let period = this.vPeriod
		let periodUnit = this.vPeriodUnit
		if (period == null || period.toString().length == 0) {
			period = 1
		}
		if (periodUnit == null || periodUnit.length == 0) {
			periodUnit = "day"
		}
		return {
			periodConfig: {
				period: period,
				unit: periodUnit
			}
		}
	},
	watch: {
		"periodConfig.period": function (v) {
			v = parseInt(v)
			if (isNaN(v) || v <= 0) {
				v = 1
			}
			this.periodConfig.period = v
		}
	},
	template: `<div>
	<input type="hidden" name="periodJSON" :value="JSON.stringify(periodConfig)"/>
	<div class="ui fields inline">
		<div class="ui field">
			<input type="text" v-model="periodConfig.period" maxlength="4" size="4"/>
		</div>
		<div class="ui field">
			<select class="ui dropdown" v-model="periodConfig.unit">
				<option value="minute">分钟</option>
				<option value="hour">小时</option>
				<option value="day">天</option>
				<option value="week">周</option>
				<option value="month">月</option>
			</select>
		</div>
	</div>
	<p class="comment">在此周期内同一对象累积为同一数据。</p>
</div>`
})

Vue.component("traffic-limit-config-box", {
	props: ["v-traffic-limit"],
	data: function () {
		let config = this.vTrafficLimit
		if (config == null) {
			config = {
				isOn: false,
				dailySize: {
					count: -1,
					unit: "gb"
				},
				monthlySize: {
					count: -1,
					unit: "gb"
				},
				totalSize: {
					count: -1,
					unit: "gb"
				},
				noticePageBody: ""
			}
		}
		if (config.dailySize == null) {
			config.dailySize = {
				count: -1,
				unit: "gb"
			}
		}
		if (config.monthlySize == null) {
			config.monthlySize = {
				count: -1,
				unit: "gb"
			}
		}
		if (config.totalSize == null) {
			config.totalSize = {
				count: -1,
				unit: "gb"
			}
		}
		return {
			config: config
		}
	},
	methods: {
		showBodyTemplate: function () {
			this.config.noticePageBody = `<!DOCTYPE html>
<html>
<head>
<title>Traffic Limit Exceeded Warning</title>
<body>

<h1>Traffic Limit Exceeded Warning</h1>
<p>The site traffic has exceeded the limit. Please contact with the site administrator.</p>
<address>Request ID: \${requestId}.</address>

</body>
</html>`
		}
	},
	template: `<div>
	<input type="hidden" name="trafficLimitJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<tbody>
			<tr>
				<td class="title">启用流量限制</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
					<p class="comment">注意：由于流量统计是每5分钟统计一次，所以超出流量限制后，对用户的提醒也会有所延迟。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>日流量限制</td>
				<td>
					<size-capacity-box :v-value="config.dailySize"></size-capacity-box>
				</td>
			</tr>
			<tr>
				<td>月流量限制</td>
				<td>
					<size-capacity-box :v-value="config.monthlySize"></size-capacity-box>
				</td>
			</tr>
			<!--<tr>
				<td>总体限制</td>
				<td>
					<size-capacity-box :v-value="config.totalSize"></size-capacity-box>
					<p class="comment"></p>
				</td>
			</tr>-->
			<tr>
				<td>网页提示内容</td>
				<td>
					<textarea v-model="config.noticePageBody"></textarea>
					<p class="comment"><a href="" @click.prevent="showBodyTemplate">[使用模板]</a>。当达到流量限制时网页显示的HTML内容，不填写则显示默认的提示内容，适用于网站类服务。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-firewall-captcha-options", {
	props: ["v-captcha-options"],
	mounted: function () {
		this.updateSummary()
	},
	data: function () {
		let options = this.vCaptchaOptions
		if (options == null) {
			options = {
				captchaType: "default",
				countLetters: 0,
				life: 0,
				maxFails: 0,
				failBlockTimeout: 0,
				failBlockScopeAll: false,
				uiIsOn: false,
				uiTitle: "",
				uiPrompt: "",
				uiButtonTitle: "",
				uiShowRequestId: true,
				uiCss: "",
				uiFooter: "",
				uiBody: "",
				cookieId: "",
				lang: "",
				geeTestConfig: {
					isOn: false,
					captchaId: "",
					captchaKey: ""
				}
			}
		}
		if (options.countLetters <= 0) {
			options.countLetters = 6
		}

		if (options.captchaType == null || options.captchaType.length == 0) {
			options.captchaType = "default"
		}


		return {
			options: options,
			isEditing: false,
			summary: "",
			uiBodyWarning: "",
			captchaTypes: window.WAF_CAPTCHA_TYPES
		}
	},
	watch: {
		"options.countLetters": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			} else if (i < 0) {
				i = 0
			} else if (i > 10) {
				i = 10
			}
			this.options.countLetters = i
		},
		"options.life": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.life = i
			this.updateSummary()
		},
		"options.maxFails": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.maxFails = i
			this.updateSummary()
		},
		"options.failBlockTimeout": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.failBlockTimeout = i
			this.updateSummary()
		},
		"options.failBlockScopeAll": function (v) {
			this.updateSummary()
		},
		"options.captchaType": function (v) {
			this.updateSummary()
		},
		"options.uiIsOn": function (v) {
			this.updateSummary()
		},
		"options.uiBody": function (v) {
			if (/<form(>|\s).+\$\{body}.*<\/form>/s.test(v)) {
				this.uiBodyWarning = "页面模板中不能使用<form></form>标签包裹\${body}变量，否则将导致验证码表单无法提交。"
			} else {
				this.uiBodyWarning = ""
			}
		},
		"options.geeTestConfig.isOn": function (v) {
			this.updateSummary()
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		},
		updateSummary: function () {
			let summaryList = []
			if (this.options.life > 0) {
				summaryList.push("有效时间" + this.options.life + "秒")
			}
			if (this.options.maxFails > 0) {
				summaryList.push("最多失败" + this.options.maxFails + "次")
			}
			if (this.options.failBlockTimeout > 0) {
				summaryList.push("失败拦截" + this.options.failBlockTimeout + "秒")
			}
			if (this.options.failBlockScopeAll) {
				summaryList.push("尝试全局封禁")
			}

			let that = this
			let typeDef = this.captchaTypes.$find(function (k, v) {
				return v.code == that.options.captchaType
			})
			if (typeDef != null) {
				summaryList.push("默认验证方式：" + typeDef.name)
			}

			if (this.options.captchaType == "default") {
				if (this.options.uiIsOn) {
					summaryList.push("定制UI")
				}
			}

			if (this.options.geeTestConfig != null && this.options.geeTestConfig.isOn) {
				summaryList.push("已配置极验")
			}

			if (summaryList.length == 0) {
				this.summary = "默认配置"
			} else {
				this.summary = summaryList.join(" / ")
			}
		},
		confirm: function () {
			this.isEditing = false
		}
	},
	template: `<div>
	<input type="hidden" name="captchaOptionsJSON" :value="JSON.stringify(options)"/>
	<a href="" @click.prevent="edit">{{summary}} <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<div v-show="isEditing" style="margin-top: 0.5em">
		<table class="ui table definition selectable">
			<tbody>
				<tr>
					<td>默认验证方式</td>
					<td>
						<select class="ui dropdown auto-width" v-model="options.captchaType">
							<option v-for="captchaDef in captchaTypes" :value="captchaDef.code">{{captchaDef.name}}</option>
						</select>
						<p class="comment" v-for="captchaDef in captchaTypes" v-if="captchaDef.code == options.captchaType">{{captchaDef.description}}</p>
					</td>
				</tr>
				<tr>
					<td class="title">有效时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">验证通过后在这个时间内不再验证，默认600秒。</p>
					</td>
				</tr>
				<tr>
					<td>最多失败次数</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.maxFails" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">次</span>
						</div>
						<p class="comment"><span v-if="options.maxFails > 0 && options.maxFails < 5" class="red">建议填入一个不小于5的数字，以减少误判几率。</span>允许用户失败尝试的最多次数，超过这个次数将被自动加入黑名单。如果为空或者为0，表示不限制。</p>
					</td>
				</tr>
				<tr>
					<td>失败拦截时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.failBlockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">在达到最多失败次数（大于0）时，自动拦截的时长；如果为0表示不自动拦截。</p>
					</td>
				</tr>
				<tr>
					<td>失败全局封禁</td>
					<td>
						<checkbox v-model="options.failBlockScopeAll"></checkbox>
						<p class="comment">选中后，表示允许系统尝试全局封禁某个IP，以提升封禁性能。</p>
					</td>
				</tr>
				
				<tr v-show="options.captchaType == 'default'">
					<td>验证码中数字个数</td>
					<td>
						<select class="ui dropdown auto-width" v-model="options.countLetters">
							<option v-for="i in 10" :value="i">{{i}}</option>
						</select>
					</td>
				</tr>
				<tr v-show="options.captchaType == 'default'">
					<td class="color-border">定制UI</td>
					<td><checkbox v-model="options.uiIsOn"></checkbox></td>
				</tr>
			</tbody>
			<tbody v-show="options.uiIsOn && options.captchaType == 'default'">
				<tr>
					<td class="color-border">页面标题</td>
					<td>
						<input type="text" v-model="options.uiTitle" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
					</td>
				</tr>
				<tr>
					<td class="color-border">按钮标题</td>
					<td>
						<input type="text" v-model="options.uiButtonTitle" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<p class="comment">类似于<code-label>提交验证</code-label>。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">显示请求ID</td>
					<td>
						<checkbox v-model="options.uiShowRequestId"></checkbox>
						<p class="comment">在界面上显示请求ID，方便用户报告问题。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">CSS样式</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiCss" rows="2"></textarea>
					</td>
				</tr>
				<tr>
					<td class="color-border">页头提示</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiPrompt" rows="2"></textarea>
						<p class="comment">类似于<code-label>请输入上面的验证码</code-label>，支持HTML。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">页尾提示</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiFooter" rows="2"></textarea>
						<p class="comment">支持HTML。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">页面模板</td>
					<td>
						<textarea spellcheck="false" rows="2" v-model="options.uiBody"></textarea>
						<p class="comment"><span v-if="uiBodyWarning.length > 0" class="red">警告：{{uiBodyWarning}}</span><span v-if="options.uiBody.length > 0 && options.uiBody.indexOf('\${body}') < 0 " class="red">模板中必须包含\${body}表示验证码表单！</span>整个页面的模板，支持HTML，其中必须使用<code-label>\${body}</code-label>变量代表验证码表单，否则将无法正常显示验证码。</p>
					</td>
				</tr>
			</tbody>
		</table>
		
		<table class="ui table definition selectable">
			<tr>
				<td class="title">允许用户使用极验</td>
				<td><checkbox v-model="options.geeTestConfig.isOn"></checkbox>
					<p class="comment">选中后，表示允许用户在WAF设置中选择极验。</p>
				</td>
			</tr>
			<tbody v-show="options.geeTestConfig.isOn">
				<tr>
					<td class="color-border">极验-验证ID *</td>
					<td>
						<input type="text" maxlength="100" name="geetestCaptchaId" v-model="options.geeTestConfig.captchaId" spellcheck="false"/>
						<p class="comment">在极验控制台--业务管理中获取。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">极验-验证Key *</td>
					<td>
						<input type="text" maxlength="100" name="geetestCaptchaKey" v-model="options.geeTestConfig.captchaKey" spellcheck="false"/>
						<p class="comment">在极验控制台--业务管理中获取。</p>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
`
})

Vue.component("user-agent-config-box", {
	props: ["v-is-location", "v-is-group", "value"],
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				filters: []
			}
		}
		if (config.filters == null) {
			config.filters = []
		}
		return {
			config: config,
			isAdding: false,
			addingFilter: {
				keywords: [],
				action: "deny"
			},
			moreOptionsVisible: false,
			batchKeywords: ""
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		remove: function (index) {
			let that = this
			teaweb.confirm("确定要删除此名单吗？", function () {
				that.config.filters.$remove(index)
			})
		},
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.batchKeywords.focus()
			})
		},
		confirm: function () {
			if (this.addingFilter.action == "deny") {
				this.config.filters.push(this.addingFilter)
			} else {
				let index = -1
				this.config.filters.forEach(function (filter, filterIndex) {
					if (filter.action == "allow") {
						index = filterIndex
					}
				})

				if (index < 0) {
					this.config.filters.unshift(this.addingFilter)
				} else {
					this.config.filters.$insert(index + 1, this.addingFilter)
				}
			}

			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.addingFilter = {
				keywords: [],
				action: "deny"
			}
			this.batchKeywords = ""
		},
		changeKeywords: function (keywords) {
			let arr = keywords.split(/\n/)
			let resultKeywords = []
			arr.forEach(function (keyword){
				keyword = keyword.trim()
				if (!resultKeywords.$contains(keyword)) {
					resultKeywords.push(keyword)
				}
			})
			this.addingFilter.keywords = resultKeywords
		},
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		}
	},
	template: `<div>
	<input type="hidden" name="userAgentJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用UA名单</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示开启UserAgent名单。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>UA名单</td>
				<td>
					<div v-if="config.filters.length > 0">
						<table class="ui table celled">
							<thead class="full-width">
								<tr>
									<th>UA关键词</th>
									<th class="two wide">动作</th>
									<th class="one op">操作</th>
								</tr>
							</thead>
							<tbody v-for="(filter, index) in config.filters">
								<tr>
									<td style="background: white">
										<span v-for="keyword in filter.keywords" class="ui label basic tiny">
											<span v-if="keyword.length > 0">{{keyword}}</span>
											<span v-if="keyword.length == 0" class="disabled">[空]</span>
										</span>
									</td>
									<td>
										<span v-if="filter.action == 'allow'" class="green">允许</span><span v-if="filter.action == 'deny'" class="red">不允许</span>
									</td>
									<td><a href="" @click.prevent="remove(index)">删除</a></td>
								</tr>
							</tbody>
						</table>
					</div>
					<div v-if="isAdding" style="margin-top: 0.5em">
						<table class="ui table definition">
							<tr>
								<td class="title">UA关键词</td>
								<td>
									<textarea v-model="batchKeywords" @input="changeKeywords(batchKeywords)" ref="batchKeywords" style="width: 20em" placeholder="*浏览器标识*"></textarea>
									<p class="comment">每行一个关键词；不区分大小写，比如<code-label>Chrome</code-label>；支持<code-label>*</code-label>通配符，比如<code-label>*Firefox*</code-label>；也支持空行关键词，表示空UserAgent。</p>
								</td>
							</tr>
							<tr>
								<td>动作</td>
								<td><select class="ui dropdown auto-width" v-model="addingFilter.action">
										<option value="deny">不允许</option>
										<option value="allow">允许</option>
									</select>
								</td>
							</tr>
						</table>
						<button type="button" class="ui button tiny" @click.prevent="confirm">保存</button> &nbsp; <a href="" @click.prevent="cancel" title="取消"><i class="icon remove small"></i></a>
					</div>
					<div v-show="!isAdding" style="margin-top: 0.5em">
						<button class="ui button tiny" type="button" @click.prevent="add">+</button>
					</div>
				</td>
			</tr>
			<tr>
				<td colspan="2"><more-options-indicator @change="showMoreOptions"></more-options-indicator></td>
			</tr>
		</tbody>
		<tbody v-show="moreOptionsVisible && isOn()">
			<tr>
				<td>例外URL</td>
				<td>
					<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
				</td>
			</tr>
			<tr>
				<td>限制URL</td>
				<td>
					<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了限制URL，表示只对这些URL进行处理；如果不填则表示支持所有的URL。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-pages-box", {
	props: ["v-pages"],
	data: function () {
		let pages = []
		if (this.vPages != null) {
			pages = this.vPages
		}

		return {
			pages: pages
		}
	},
	methods: {
		addPage: function () {
			let that = this
			teaweb.popup("/servers/server/settings/pages/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.pages.push(resp.data.page)
					that.notifyChange()
				}
			})
		},
		updatePage: function (pageIndex, pageId) {
			let that = this
			teaweb.popup("/servers/server/settings/pages/updatePopup?pageId=" + pageId, {
				height: "26em",
				callback: function (resp) {
					Vue.set(that.pages, pageIndex, resp.data.page)
					that.notifyChange()
				}
			})
		},
		removePage: function (pageIndex) {
			let that = this
			teaweb.confirm("确定要移除此页面吗？", function () {
				that.pages.$remove(pageIndex)
				that.notifyChange()
			})
		},
		notifyChange: function () {
			let parent = this.$el.parentNode
			while (true) {
				if (parent == null) {
					break
				}
				if (parent.tagName == "FORM") {
					break
				}
				parent = parent.parentNode
			}
			if (parent != null) {
				setTimeout(function () {
					Tea.runActionOn(parent)
				}, 100)
			}
		}
	},
	template: `<div>
<input type="hidden" name="pagesJSON" :value="JSON.stringify(pages)"/>

<div v-if="pages.length > 0">
	<table class="ui table selectable celled">
		<thead>
			<tr>
				<th class="two wide">响应状态码</th>
				<th>页面类型</th>
				<th class="two wide">新状态码</th>
				<th>例外URL</th>
				<th>限制URL</th>
				<th class="two op">操作</th>
			</tr>	
		</thead>
		<tr v-for="(page,index) in pages">
			<td>
				<a href="" @click.prevent="updatePage(index, page.id)">
					<span v-if="page.status != null && page.status.length == 1">{{page.status[0]}}</span>
					<span v-else>{{page.status}}</span>
					
					<i class="icon expand small"></i>
				</a>
			</td>
			<td style="word-break: break-all">
				<div v-if="page.bodyType == 'url'">
					{{page.url}}
					<div>
						<grey-label>读取URL</grey-label>
					</div>
				</div>
				<div v-if="page.bodyType == 'redirectURL'">
					{{page.url}}
					<div>
						<grey-label>跳转URL</grey-label>	
						<grey-label v-if="page.newStatus > 0">{{page.newStatus}}</grey-label>
					</div>
				</div>
				<div v-if="page.bodyType == 'html'">
					[HTML内容]
					<div>
						<grey-label v-if="page.newStatus > 0">{{page.newStatus}}</grey-label>
					</div>
				</div>
			</td>
			<td>
				<span v-if="page.newStatus > 0">{{page.newStatus}}</span>
				<span v-else class="disabled">保持</span>	
			</td>
			<td>
				<div v-if="page.exceptURLPatterns != null && page.exceptURLPatterns">
					<span v-for="urlPattern in page.exceptURLPatterns" class="ui basic label small">{{urlPattern.pattern}}</span>
				</div>
				<span v-else class="disabled">-</span>
			</td>
			<td>
				<div v-if="page.onlyURLPatterns != null && page.onlyURLPatterns">
					<span v-for="urlPattern in page.onlyURLPatterns" class="ui basic label small">{{urlPattern.pattern}}</span>
				</div>
				<span v-else class="disabled">-</span>
			</td>
			<td>
				<a href="" title="修改" @click.prevent="updatePage(index, page.id)">修改</a> &nbsp; 
				<a href="" title="删除" @click.prevent="removePage(index)">删除</a>
			</td>
		</tr>
	</table>
</div>
<div style="margin-top: 1em">
	<button class="ui button small" type="button" @click.prevent="addPage()">+添加自定义页面</button>
</div>
<div class="ui margin"></div>
</div>`
})

Vue.component("firewall-syn-flood-config-box", {
	props: ["v-syn-flood-config"],
	data: function () {
		let config = this.vSynFloodConfig
		if (config == null) {
			config = {
				isOn: false,
				minAttempts: 10,
				timeoutSeconds: 600,
				ignoreLocal: true
			}
		}
		return {
			config: config,
			isEditing: false,
			minAttempts: config.minAttempts,
			timeoutSeconds: config.timeoutSeconds
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	watch: {
		minAttempts: function (v) {
			let count = parseInt(v)
			if (isNaN(count)) {
				count = 10
			}
			if (count < 5) {
				count = 5
			}
			this.config.minAttempts = count
		},
		timeoutSeconds: function (v) {
			let seconds = parseInt(v)
			if (isNaN(seconds)) {
				seconds = 10
			}
			if (seconds < 60) {
				seconds = 60
			}
			this.config.timeoutSeconds = seconds
		}
	},
	template: `<div>
	<input type="hidden" name="synFloodJSON" :value="JSON.stringify(config)"/>
	<a href="" @click.prevent="edit">
		<span v-if="config.isOn">
			已启用 / <span>空连接次数：{{config.minAttempts}}次/分钟</span> / 封禁时长：{{config.timeoutSeconds}}秒 <span v-if="config.ignoreLocal">/ 忽略局域网访问</span>
		</span>
		<span v-else>未启用</span>
		<i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i>
	</a>
	
	<table class="ui table selectable" v-show="isEditing">
		<tr>
			<td class="title">启用</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment">启用后，WAF将会尝试自动检测并阻止SYN Flood攻击。此功能需要节点已安装并启用nftables或Firewalld。</p>
			</td>
		</tr>
		<tr>
			<td>空连接次数</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="minAttempts" style="width: 5em" maxlength="4"/>
					<span class="ui label">次/分钟</span>
				</div>
				<p class="comment">超过此数字的"空连接"将被视为SYN Flood攻击，为了防止误判，此数值默认不小于5。</p>
			</td>
		</tr>
		<tr>
			<td>封禁时长</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeoutSeconds" style="width: 5em" maxlength="8"/>
					<span class="ui label">秒</span>
				</div>
			</td>
		</tr>
		<tr>
			<td>忽略局域网访问</td>
			<td>
				<checkbox v-model="config.ignoreLocal"></checkbox>
			</td>
		</tr>
	</table>
</div>`
})

Vue.component("http-firewall-region-selector", {
	props: ["v-type", "v-countries"],
	data: function () {
		let countries = this.vCountries
		if (countries == null) {
			countries = []
		}

		return {
			listType: this.vType,
			countries: countries
		}
	},
	methods: {
		addCountry: function () {
			let selectedCountryIds = this.countries.map(function (country) {
				return country.id
			})
			let that = this
			teaweb.popup("/servers/server/settings/waf/ipadmin/selectCountriesPopup?type=" + this.listType + "&selectedCountryIds=" + selectedCountryIds.join(","), {
				width: "52em",
				height: "30em",
				callback: function (resp) {
					that.countries = resp.data.selectedCountries
					that.$forceUpdate()
					that.notifyChange()
				}
			})
		},
		removeCountry: function (index) {
			this.countries.$remove(index)
			this.notifyChange()
		},
		resetCountries: function () {
			this.countries = []
			this.notifyChange()
		},
		notifyChange: function () {
			this.$emit("change", {
				"countries": this.countries
			})
		}
	},
	template: `<div>
	<span v-if="countries.length == 0" class="disabled">暂时没有选择<span v-if="listType =='allow'">允许</span><span v-else>封禁</span>的区域。</span>
	<div v-show="countries.length > 0">
		<div class="ui label tiny basic" v-for="(country, index) in countries" style="margin-bottom: 0.5em">
			<input type="hidden" :name="listType + 'CountryIds'" :value="country.id"/>
			({{country.letter}}){{country.name}} <a href="" @click.prevent="removeCountry(index)" title="删除"><i class="icon remove"></i></a>
		</div>
	</div>
	<div class="ui divider"></div>
	<button type="button" class="ui button tiny" @click.prevent="addCountry">修改</button> &nbsp; <button type="button" class="ui button tiny" v-show="countries.length > 0" @click.prevent="resetCountries">清空</button>
</div>`
})

// TODO 支持关键词搜索
// TODO 改成弹窗选择
Vue.component("admin-selector", {
    props: ["v-admin-id"],
    mounted: function () {
        let that = this
        Tea.action("/admins/options")
            .post()
            .success(function (resp) {
                that.admins = resp.data.admins
            })
    },
    data: function () {
        let adminId = this.vAdminId
        if (adminId == null) {
            adminId = 0
        }
        return {
            admins: [],
            adminId: adminId
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="adminId" v-model="adminId">
        <option value="0">[选择系统用户]</option>
        <option v-for="admin in admins" :value="admin.id">{{admin.name}}（{{admin.username}}）</option>
    </select>
</div>`
})

// 绑定IP列表
Vue.component("ip-list-bind-box", {
	props: ["v-http-firewall-policy-id", "v-type"],
	mounted: function () {
		this.refresh()
	},
	data: function () {
		return {
			policyId: this.vHttpFirewallPolicyId,
			type: this.vType,
			lists: []
		}
	},
	methods: {
		bind: function () {
			let that = this
			teaweb.popup("/servers/iplists/bindHTTPFirewallPopup?httpFirewallPolicyId=" + this.policyId + "&type=" + this.type, {
				width: "50em",
				height: "34em",
				callback: function () {

				},
				onClose: function () {
					that.refresh()
				}
			})
		},
		remove: function (index, listId) {
			let that = this
			teaweb.confirm("确定要删除这个绑定的IP名单吗？", function () {
				Tea.action("/servers/iplists/unbindHTTPFirewall")
					.params({
						httpFirewallPolicyId: that.policyId,
						listId: listId
					})
					.post()
					.success(function (resp) {
						that.lists.$remove(index)
					})
			})
		},
		refresh: function () {
			let that = this
			Tea.action("/servers/iplists/httpFirewall")
				.params({
					httpFirewallPolicyId: this.policyId,
					type: this.vType
				})
				.post()
				.success(function (resp) {
					that.lists = resp.data.lists
				})
		}
	},
	template: `<div>
	<a href="" @click.prevent="bind()" style="color: rgba(0,0,0,.6)">绑定+</a> &nbsp; <span v-if="lists.length > 0"><span class="disabled small">|&nbsp;</span> 已绑定：</span>
	<div class="ui label basic small" v-for="(list, index) in lists">
		<a :href="'/servers/iplists/list?listId=' + list.id" title="点击查看详情" style="opacity: 1">{{list.name}}</a>
		<a href="" title="删除" @click.prevent="remove(index, list.id)"><i class="icon remove small"></i></a>
	</div>
</div>`
})

Vue.component("ip-list-table", {
	props: ["v-items", "v-keyword", "v-show-search-button", "v-total"/** total items >= items length **/],
	data: function () {
		let maxDeletes = 10000
		if (this.vTotal != null && this.vTotal > 0 && this.vTotal < maxDeletes) {
			maxDeletes = this.vTotal
		}

		return {
			items: this.vItems,
			keyword: (this.vKeyword != null) ? this.vKeyword : "",
			selectedAll: false,
			hasSelectedItems: false,

			MaxDeletes: maxDeletes
		}
	},
	methods: {
		updateItem: function (itemId) {
			this.$emit("update-item", itemId)
		},
		deleteItem: function (itemId) {
			this.$emit("delete-item", itemId)
		},
		viewLogs: function (itemId) {
			teaweb.popup("/servers/iplists/accessLogsPopup?itemId=" + itemId, {
				width: "50em",
				height: "30em"
			})
		},
		changeSelectedAll: function () {
			let boxes = this.$refs.itemCheckBox
			if (boxes == null) {
				return
			}

			let that = this
			boxes.forEach(function (box) {
				box.checked = that.selectedAll
			})

			this.hasSelectedItems = this.selectedAll
		},
		changeSelected: function (e) {
			let that = this
			that.hasSelectedItems = false
			let boxes = that.$refs.itemCheckBox
			if (boxes == null) {
				return
			}
			boxes.forEach(function (box) {
				if (box.checked) {
					that.hasSelectedItems = true
				}
			})
		},
		deleteAll: function () {
			let boxes = this.$refs.itemCheckBox
			if (boxes == null) {
				return
			}
			let itemIds = []
			boxes.forEach(function (box) {
				if (box.checked) {
					itemIds.push(box.value)
				}
			})
			if (itemIds.length == 0) {
				return
			}

			Tea.action("/servers/iplists/deleteItems")
				.post()
				.params({
					itemIds: itemIds
				})
				.success(function () {
					teaweb.successToast("批量删除成功", 1200, teaweb.reload)
				})
		},
		deleteCount: function () {
			let that = this
			teaweb.confirm("确定要批量删除当前列表中的" + this.MaxDeletes + "个IP吗？", function () {
				let query = window.location.search
				if (query.startsWith("?")) {
					query = query.substring(1)
				}
				Tea.action("/servers/iplists/deleteCount?" + query)
					.post()
					.params({count: that.MaxDeletes})
					.success(function () {
						teaweb.successToast("批量删除成功", 1200, teaweb.reload)
					})
			})
		},
		formatSeconds: function (seconds) {
			if (seconds < 60) {
				return seconds + "秒"
			}
			if (seconds < 3600) {
				return Math.ceil(seconds / 60) + "分钟"
			}
			if (seconds < 86400) {
				return Math.ceil(seconds / 3600) + "小时"
			}
			return Math.ceil(seconds / 86400) + "天"
		},
		cancelChecked: function () {
			this.hasSelectedItems = false
			this.selectedAll = false

			let boxes = this.$refs.itemCheckBox
			if (boxes == null) {
				return
			}
			boxes.forEach(function (box) {
				box.checked = false
			})
		}
	},
	template: `<div>
 <div v-show="hasSelectedItems">
 	<div class="ui divider"></div>
 	<button class="ui button basic" type="button" @click.prevent="deleteAll">批量删除所选</button>
 	&nbsp; &nbsp; 
 	<button class="ui button basic" type="button" @click.prevent="deleteCount" v-if="vTotal != null && vTotal >= MaxDeletes">批量删除所有搜索结果（{{MaxDeletes}}个）</button>
 	
 	&nbsp; &nbsp; 
 	<button class="ui button basic" type="button" @click.prevent="cancelChecked">取消选中</button>
</div>
 <table class="ui table selectable celled" v-if="items.length > 0">
        <thead>
            <tr>
            	<th style="width: 1em">
            		<div class="ui checkbox">
						<input type="checkbox" v-model="selectedAll" @change="changeSelectedAll"/>
						<label></label>
					</div>
				</th>
                <th style="width:18em">IP</th>
                <th style="width: 6em">类型</th>
                <th style="width: 6em">级别</th>
                <th style="width: 12em">过期时间</th>
                <th>备注</th>
                <th class="three op">操作</th>
            </tr>
        </thead>
		<tbody v-for="item in items">
			<tr>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" :value="item.id" @change="changeSelected" ref="itemCheckBox"/>
						<label></label>
					</div>
				</td>
				<td>
					<span v-if="item.type != 'all'" :class="{green: item.list != null && item.list.type == 'white', grey: item.list != null && item.list.type == 'grey'}">
						<span v-if="item.value != null && item.value.length > 0"><keyword :v-word="keyword">{{item.value}}</keyword></span>
						<span v-else>
							<keyword :v-word="keyword">{{item.ipFrom}}</keyword> <span> <span class="small red" v-if="item.isRead != null && !item.isRead">&nbsp;New&nbsp;</span>&nbsp;<a :href="'/servers/iplists?ip=' + item.ipFrom" v-if="vShowSearchButton" title="搜索此IP"><span><i class="icon search small" style="color: #ccc"></i></span></a></span>
							<span v-if="item.ipTo.length > 0"> - <keyword :v-word="keyword">{{item.ipTo}}</keyword></span>
						</span>
					</span>
					<span v-else class="disabled">*</span>
					
					<div v-if="item.region != null && item.region.length > 0">
						<span class="grey small">{{item.region}}</span>
						<span v-if="item.isp != null && item.isp.length > 0 && item.isp != '内网IP'" class="grey small"><span class="disabled">|</span> {{item.isp}}</span>
					</div>
					<div v-else-if="item.isp != null && item.isp.length > 0 && item.isp != '内网IP'"><span class="grey small">{{item.isp}}</span></div>
				
					<div v-if="item.createdTime != null">
						<span class="small grey">添加于 {{item.createdTime}}
							<span v-if="item.list != null && item.list.id > 0">
								@ 
								
								<a :href="'/servers/iplists/list?listId=' + item.list.id" v-if="item.policy.id == 0"><span>[<span v-if="item.list.type == 'black'">黑</span><span v-if="item.list.type == 'white'">白</span><span v-if="item.list.type == 'grey'">灰</span>名单：{{item.list.name}}]</span></a>
								<span v-else>[<span v-if="item.list.type == 'black'">黑</span><span v-if="item.list.type == 'white'">白</span><span v-if="item.list.type == 'grey'">灰</span>名单：{{item.list.name}}</span>
								
								<span v-if="item.policy.id > 0">
									<span v-if="item.policy.server != null">
										<a :href="'/servers/server/settings/waf/ipadmin/allowList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'white'">[网站：{{item.policy.server.name}}]</a>
										<a :href="'/servers/server/settings/waf/ipadmin/denyList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'black'">[网站：{{item.policy.server.name}}]</a>
										<a :href="'/servers/server/settings/waf/ipadmin/greyList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'grey'">[网站：{{item.policy.server.name}}]</a>
									</span>
									<span v-else>
										<a :href="'/servers/components/waf/ipadmin/lists?firewallPolicyId=' + item.policy.id +  '&type=' + item.list.type">[策略：{{item.policy.name}}]</a>
									</span>
								</span>
							</span>
						</span>
					</div>
				</td>
				<td>
					<span v-if="item.type.length == 0">IPv4</span>
					<span v-else-if="item.type == 'ipv4'">IPv4</span>
					<span v-else-if="item.type == 'ipv6'">IPv6</span>
					<span v-else-if="item.type == 'all'"><strong>所有IP</strong></span>
				</td>
				<td>
					<span v-if="item.eventLevelName != null && item.eventLevelName.length > 0">{{item.eventLevelName}}</span>
					<span v-else class="disabled">-</span>
				</td>
				<td>
					<div v-if="item.expiredTime.length > 0">
						{{item.expiredTime}}
						<div v-if="item.isExpired && item.lifeSeconds == null" style="margin-top: 0.5em">
							<span class="ui label tiny basic red">已过期</span>
						</div>
						<div  v-if="item.lifeSeconds != null">
							<span class="small grey" v-if="item.lifeSeconds > 0">{{formatSeconds(item.lifeSeconds)}}</span>
							<span class="small red" v-if="item.lifeSeconds < 0">已过期</span>
						</div>
					</div>
					<span v-else class="disabled">不过期</span>
				</td>
				<td>
					<span v-if="item.reason.length > 0">{{item.reason}}</span>
					<span v-else class="disabled">-</span>
					
					<div v-if="item.sourceNode != null && item.sourceNode.id > 0" style="margin-top: 0.4em">
						<a :href="'/clusters/cluster/node?clusterId=' + item.sourceNode.clusterId + '&nodeId=' + item.sourceNode.id"><span class="small"><i class="icon cloud"></i>{{item.sourceNode.name}}</span></a>
					</div>
					<div style="margin-top: 0.4em" v-if="item.sourceServer != null && item.sourceServer.id > 0">
						<a :href="'/servers/server?serverId=' + item.sourceServer.id" style="border: 0"><span class="small "><i class="icon clone outline"></i>{{item.sourceServer.name}}</span></a>
					</div>
					<div v-if="item.sourcePolicy != null && item.sourcePolicy.id > 0" style="margin-top: 0.4em">
						<a :href="'/servers/components/waf/group?firewallPolicyId=' +  item.sourcePolicy.id + '&type=inbound&groupId=' + item.sourceGroup.id + '#set' + item.sourceSet.id" v-if="item.sourcePolicy.serverId == 0"><span class="small "><i class="icon shield"></i>{{item.sourcePolicy.name}} &raquo; {{item.sourceGroup.name}} &raquo; {{item.sourceSet.name}}</span></a>
						<a :href="'/servers/server/settings/waf/group?serverId=' + item.sourcePolicy.serverId + '&firewallPolicyId=' + item.sourcePolicy.id + '&type=inbound&groupId=' + item.sourceGroup.id + '#set' + item.sourceSet.id" v-if="item.sourcePolicy.serverId > 0"><span class="small "><i class="icon shield"></i> {{item.sourcePolicy.name}} &raquo; {{item.sourceGroup.name}} &raquo; {{item.sourceSet.name}}</span></a>
					</div>
				</td>
				<td>
					<a href="" @click.prevent="viewLogs(item.id)">日志</a> &nbsp;
					<a href="" @click.prevent="updateItem(item.id)">修改</a> &nbsp;
					<a href="" @click.prevent="deleteItem(item.id)">删除</a>
				</td>
			</tr>
        </tbody>
    </table>
</div>`
})

Vue.component("ip-item-text", {
    props: ["v-item"],
    template: `<span>
    <span v-if="vItem.type == 'all'">*</span>
    <span v-else>
    	<span v-if="vItem.value != null && vItem.value.length > 0">{{vItem.value}}</span>
    	<span v-else>
			{{vItem.ipFrom}}
			<span v-if="vItem.ipTo != null &&vItem.ipTo.length > 0">- {{vItem.ipTo}}</span>
		</span>
	</span>
    <span v-if="vItem.eventLevelName != null && vItem.eventLevelName.length > 0">&nbsp; 级别：{{vItem.eventLevelName}}</span>
</span>`
})

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

Vue.component("sms-sender", {
	props: ["value", "name"],
	mounted: function () {
		this.initType(this.config.type)
	},
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isOn: false,
				type: "webHook",
				webHookParams: {
					url: "",
					method: "POST"
				},
				aliyunSMSParams: {
					sign: "",
					templateCode: "",
					codeVarName: "code",
					accessKeyId: "",
					accessKeySecret: ""
				},
				tencentSMSParams: {
					sdkAppId: "",
					sign: "",
					templateId: "",
					accessKeyId: "",
					accessKeySecret: ""
				}
			}
		}

		if (config.aliyunSMSParams == null) {
			Vue.set(config, "aliyunSMSParams", {
				sign: "",
				templateCode: "",
				codeVarName: "code",
				accessKeyId: "",
				accessKeySecret: ""
			})
		}
		if (config.tencentSMSParams == null) {
			Vue.set(config, "tencentSMSParams",  {
				sdkAppId: "",
				sign: "",
				templateId: "",
				accessKeyId: "",
				accessKeySecret: ""
			})
		}

		return {
			config: config
		}
	},
	watch: {
		"config.type": function (v) {
			this.initType(v)
		}
	},
	methods: {
		initType: function (v) {
			// initialize params
			switch (v) {
				case "webHook":
					if (this.config.webHookParams == null) {
						this.config.webHookParams = {
							url: "",
							method: "POST"
						}
					}
					break
			}
		},
		test: function () {
			window.TESTING_SMS_CONFIG = this.config
			teaweb.popup("/users/setting/smsTest", {
				height: "22em"
			})
		}
	},
	template: `<div>
	<input type="hidden" :name="name" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<tbody>
			<tr>
				<td class="title">启用</td>
				<td><checkbox v-model="config.isOn"></checkbox></td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>发送渠道</td>
				<td>
					<select class="ui dropdown auto-width" v-model="config.type">
						<option value="webHook">自定义HTTP接口</option>
						<option value="aliyunSMS">阿里云短信</option>
						<option value="tencentSMS">腾讯云短信</option>
					</select>
					<p class="comment" v-if="config.type == 'webHook'">通过HTTP接口的方式调用你的自定义发送短信接口。</p>
					<p class="comment" v-if="config.type == 'aliyunSMS'">通过阿里云短信服务发送短信接口；<strong>目前仅支持发送验证码</strong>。</p>
					<p class="comment" v-if="config.type == 'tencentSMS'">通过腾讯云短信服务发送短信接口；<strong>目前仅支持发送验证码</strong>。</p>
				</td>				
			</tr>
			
			<!-- webhook -->
			<tr v-if="config.type == 'webHook' && config.webHookParams != null">
				<td class="color-border">HTTP接口的URL地址 *</td>
				<td>
					<input type="text" maxlength="100" placeholder="https://..." v-model="config.webHookParams.url"/>
					<p class="comment">接收发送短信请求的URL，必须以<code-label>http://</code-label>或<code-label>https://</code-label>开头。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'webHook' && config.webHookParams != null">
				<td class="color-border">HTTP接口的请求方法</td>
				<td>
					<select class="ui dropdown auto-width" v-model="config.webHookParams.method">
						<option value="GET">GET</option>
						<option value="POST">POST</option>
					</select>
					<p class="comment" v-if="config.webHookParams.method == 'GET'">以在URL参数中加入mobile、body和code三个参数（<code-label>YOUR_API_URL?mobile=手机号&amp;body=短信内容&code=验证码</code-label>)的方式调用你的HTTP接口的URL地址；状态码返回200表示成功。</p>
					<p class="comment" v-if="config.webHookParams.method == 'POST'">通过POST表单发送mobile、body和code三个参数（<code-label>mobile=手机号&amp;body=短信内容&code=验证码</code-label>）的方式调用你的HTTP接口URL地址；状态码返回200表示成功。</p>
				</td>
			</tr>
			
			<!-- aliyun sms -->
			<tr v-if="config.type == 'aliyunSMS'">
				<td class="color-border">签名名称 *</td>
				<td><input type="text" v-model="config.aliyunSMSParams.sign" maxlength="12"/>
					<p class="comment">在阿里云短信服务 “签名管理” 中添加并通过审核后才能使用。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'aliyunSMS'">
				<td class="color-border">模板CODE *</td>
				<td>
					<input type="text" v-model="config.aliyunSMSParams.templateCode" maxlength="30"/>
					<p class="comment">在阿里云短信服务 “模板管理” 中添加并通过审核后才能使用。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'aliyunSMS'">
				<td class="color-border">模板中验证码变量名称 *</td>
				<td>
					<input type="text" v-model="config.aliyunSMSParams.codeVarName" maxlength="30"/>
					<p class="comment">默认为<code-label>code</code-label>，不需要带\${}等符号，即表示在模板中使用<code-label>\${<span>{{ config.aliyunSMSParams.codeVarName }}</span>}</code-label>代表要发送的验证码。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'aliyunSMS'">
				<td class="color-border">AccessKey ID *</td>
				<td>
					<input type="text" v-model="config.aliyunSMSParams.accessKeyId" maxlength="100"/>
					<p class="comment">在阿里云 -- RAM访问控制 -- AccessKey中可以创建和获取。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'aliyunSMS'">
				<td class="color-border">AccessKey Secret *</td>
				<td>
					<input type="text" v-model="config.aliyunSMSParams.accessKeySecret" maxlength="100"/>
					<p class="comment">和表单中的AccessKey ID对应，在阿里云 -- RAM访问控制 -- AccessKey中可以创建和获取。</p>
				</td>
			</tr>
			
			<!-- tencent sms -->
			<tr v-if="config.type == 'tencentSMS'">
				<td>SDK应用ID *</td>
				<td>
					<input type="text" v-model="config.tencentSMSParams.sdkAppId" maxlength="30"/>
					<p class="comment">在腾讯云 -- 短信 -- 应用管理 -- 应用列表中可以查看。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'tencentSMS'">
				<td>签名内容 *</td>
				<td>
					<input type="text" v-model="config.tencentSMSParams.sign" maxlength="12"/>
					<p class="comment">比如“腾讯云”，在腾讯云 -- 短信 -- 签名管理中可以查看。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'tencentSMS'">
				<td>正文模板ID *</td>
				<td>
					<input type="text" v-model="config.tencentSMSParams.templateId" maxlength="50"/>
					<p class="comment">在腾讯云 -- 短信 -- 正文模板管理中可以查看。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'tencentSMS'">
				<td>密钥SecretId *</td>
				<td>
					<input type="text" v-model="config.tencentSMSParams.accessKeyId"/>
					<p class="comment">同SecretKey一同在腾讯云 -- 访问管理 -- API密钥管理中获取。</p>
				</td>
			</tr>
			<tr v-if="config.type == 'tencentSMS'">
				<td>密钥SecretKey *</td>
				<td>
					<input type="text" v-model="config.tencentSMSParams.accessKeySecret"/>
					<p class="comment">同SecretId一同在腾讯云 -- 访问管理 -- API密钥管理中获取。</p>
				</td>
			</tr>
			
			<tr>
				<td>发送测试</td>
				<td><a href="" @click.prevent="test">[点此测试]</a></td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("email-sender", {
	props: ["value", "name"],
	data: function () {
		let value = this.value
		if (value == null) {
			value = {
				isOn: false,
				smtpHost: "",
				smtpPort: 0,
				username: "",
				password: "",
				fromEmail: "",
				fromName: ""
			}
		}
		let smtpPortString = value.smtpPort.toString()
		if (smtpPortString == "0") {
			smtpPortString = ""
		}

		return {
			config: value,
			smtpPortString: smtpPortString
		}
	},
	watch: {
		smtpPortString: function (v) {
			let port = parseInt(v)
			if (!isNaN(port)) {
				this.config.smtpPort = port
			}
		}
	},
	methods: {
		test: function () {
			window.TESTING_EMAIL_CONFIG = this.config
			teaweb.popup("/users/setting/emailTest", {
				height: "36em"
			})
		}
	},
	template: `<div>
	<input type="hidden" :name="name" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<tbody>
			<tr>
				<td class="title">启用</td>
				<td><checkbox v-model="config.isOn"></checkbox></td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>SMTP地址 *</td>
				<td>
					<input type="text" :name="name + 'SmtpHost'" v-model="config.smtpHost"/>
					<p class="comment">SMTP主机地址，比如<code-label>smtp.qq.com</code-label>，目前仅支持TLS协议，如不清楚，请查询对应邮件服务商文档。</p>
				</td>
			</tr>
			<tr>
				<td>SMTP端口 *</td>
				<td>
					<input type="text" :name="name + 'SmtpPort'" v-model="smtpPortString" style="width: 5em" maxlength="5"/>
					<p class="comment">SMTP主机端口，比如<code-label>587</code-label>、<code-label>465</code-label>，如不清楚，请查询对应邮件服务商文档。</p>
				</td>
			</tr>
			<tr>
				<td>用户名 *</td>
				<td>
					<input type="text" :name="name + 'Username'" v-model="config.username"/>
					<p class="comment">通常为发件人邮箱地址。</p>
				</td>
			</tr>
			<tr>
				<td>密码 *</td>
				<td>
					<input type="password" :name="name + 'Password'" v-model="config.password"/>
					<p class="comment">邮箱登录密码或授权码，如不清楚，请查询对应邮件服务商文档。。</p>
				</td>
			</tr>
			<tr>
				<td>发件人Email *</td>
				<td>
					<input type="text" :name="name + 'FromEmail'" v-model="config.fromEmail" maxlength="128"/>
					<p class="comment">使用的发件人邮箱地址，通常和发件用户名一致。</p>
				</td>
			</tr>
			<tr>
				<td>发件人名称</td>
				<td>
					<input type="text" :name="name + 'FromName'" v-model="config.fromName" maxlength="30"/>
					<p class="comment">使用的发件人名称，默认使用系统设置的<a href="/settings/ui" target="_blank">产品名称</a>。</p>
				</td>
			</tr>
			<tr>
				<td>发送测试</td>
				<td><a href="" @click.prevent="test">[点此测试]</a></td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("api-node-selector", {
	props: [],
	data: function () {
		return {}
	},
	template: `<div>
	暂未实现
</div>`
})

Vue.component("api-node-addresses-box", {
	props: ["v-addrs", "v-name"],
	data: function () {
		let addrs = this.vAddrs
		if (addrs == null) {
			addrs = []
		}
		return {
			addrs: addrs
		}
	},
	methods: {
		// 添加IP地址
		addAddr: function () {
			let that = this;
			teaweb.popup("/settings/api/node/createAddrPopup", {
				height: "16em",
				callback: function (resp) {
					that.addrs.push(resp.data.addr);
				}
			})
		},

		// 修改地址
		updateAddr: function (index, addr) {
			let that = this;
			window.UPDATING_ADDR = addr
			teaweb.popup("/settings/api/node/updateAddrPopup?addressId=", {
				callback: function (resp) {
					Vue.set(that.addrs, index, resp.data.addr);
				}
			})
		},

		// 删除IP地址
		removeAddr: function (index) {
			this.addrs.$remove(index);
		}
	},
	template: `<div>
	<input type="hidden" :name="vName" :value="JSON.stringify(addrs)"/>
	<div v-if="addrs.length > 0">
		<div>
			<div v-for="(addr, index) in addrs" class="ui label small basic">
				{{addr.protocol}}://{{addr.host.quoteIP()}}:{{addr.portRange}}</span>
				<a href="" title="修改" @click.prevent="updateAddr(index, addr)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeAddr(index)"><i class="icon remove"></i></a>
			</div>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button small" type="button" @click.prevent="addAddr()">+</button>
	</div>
</div>`
})

// 给Table增加排序功能
function sortTable(callback) {
	// 引入js
	let jsFile = document.createElement("script")
	jsFile.setAttribute("src", "/js/sortable.min.js")
	jsFile.addEventListener("load", function () {
		// 初始化
		let box = document.querySelector("#sortable-table")
		if (box == null) {
			return
		}
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {
			},
			onUpdate: function (event) {
				let rows = box.querySelectorAll("tbody")
				let rowIds = []
				rows.forEach(function (row) {
					rowIds.push(parseInt(row.getAttribute("v-id")))
				})
				callback(rowIds)
			}
		})
	})
	document.head.appendChild(jsFile)
}

function sortLoad(callback) {
	let jsFile = document.createElement("script")
	jsFile.setAttribute("src", "/js/sortable.min.js")
	jsFile.addEventListener("load", function () {
		if (typeof (callback) == "function") {
			callback()
		}
	})
	document.head.appendChild(jsFile)
}


Vue.component("page-box", {
	data: function () {
		return {
			page: ""
		}
	},
	created: function () {
		let that = this;
		setTimeout(function () {
			that.page = Tea.Vue.page;
		})
	},
	template: `<div>
	<div class="page" v-html="page"></div>
</div>`
})

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

/**
 * 保存按钮
 */
Vue.component("submit-btn", {
	template: '<button class="ui button primary" type="submit"><slot>保存</slot></button>'
});

// 可以展示更多条目的角图表
Vue.component("more-items-angle", {
	props: ["v-data-url", "v-url"],
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		show: function () {
			this.visible = !this.visible
			if (this.visible) {
				this.showBox()
			} else {
				this.hideBox()
			}
		},
		showBox: function () {
			let that = this

			this.visible = true

			Tea.action(this.vDataUrl)
				.params({
					url: this.vUrl
				})
				.post()
				.success(function (resp) {
					let groups = resp.data.groups

					let boxLeft = that.$el.offsetLeft + 120;
					let boxTop = that.$el.offsetTop + 70;

					let box = document.createElement("div")
					box.setAttribute("id", "more-items-box")
					box.style.cssText = "z-index: 100; position: absolute; left: " + boxLeft + "px; top: " + boxTop + "px; max-height: 30em; overflow: auto; border-bottom: 1px solid rgba(34,36,38,.15)"
					document.body.append(box)

					let menuHTML = "<ul class=\"ui labeled menu vertical borderless\" style=\"padding: 0\">"
					groups.forEach(function (group) {
						menuHTML += "<div class=\"item header\">" + teaweb.encodeHTML(group.name) + "</div>"
						group.items.forEach(function (item) {
							menuHTML += "<a href=\"" + item.url + "\" class=\"item " + (item.isActive ? "active" : "") + "\" style=\"font-size: 0.9em;\">" + teaweb.encodeHTML(item.name) + "<i class=\"icon right angle\"></i></a>"
						})
					})
					menuHTML += "</ul>"
					box.innerHTML = menuHTML

					let listener = function (e) {
						if (e.target.tagName == "I") {
							return
						}

						if (!that.isInBox(box, e.target)) {
							document.removeEventListener("click", listener)
							that.hideBox()
						}
					}
					document.addEventListener("click", listener)
				})
		},
		hideBox: function () {
			let box = document.getElementById("more-items-box")
			if (box != null) {
				box.parentNode.removeChild(box)
			}
			this.visible = false
		},
		isInBox: function (parent, child) {
			while (true) {
				if (child == null) {
					break
				}
				if (child.parentNode == parent) {
					return true
				}
				child = child.parentNode
			}
			return false
		}
	},
	template: `<a href="" class="item" @click.prevent="show" style="padding-right: 0"><span style="font-size: 0.8em">切换</span><i class="icon angle" :class="{down: !visible, up: visible}"></i></a>`
})

/**
 * 菜单项
 */
Vue.component("menu-item", {
	props: ["href", "active", "code"],
	data: function () {
		let active = this.active
		if (typeof (active) == "undefined") {
			var itemCode = ""
			if (typeof (window.TEA.ACTION.data.firstMenuItem) != "undefined") {
				itemCode = window.TEA.ACTION.data.firstMenuItem
			}
			if (itemCode != null && itemCode.length > 0 && this.code != null && this.code.length > 0) {
				if (itemCode.indexOf(",") > 0) {
					active = itemCode.split(",").$contains(this.code)
				} else {
					active = (itemCode == this.code)
				}
			}
		}

		let href = (this.href == null) ? "" : this.href
		if (typeof (href) == "string" && href.length > 0 && href.startsWith(".")) {
			let qIndex = href.indexOf("?")
			if (qIndex >= 0) {
				href = Tea.url(href.substring(0, qIndex)) + href.substring(qIndex)
			} else {
				href = Tea.url(href)
			}
		}

		return {
			vHref: href,
			vActive: active
		}
	},
	methods: {
		click: function (e) {
			this.$emit("click", e)
		}
	},
	template: '\
		<a :href="vHref" class="item" :class="{active:vActive}" @click="click"><slot></slot></a> \
		'
});

// 使用Icon的链接方式
Vue.component("link-icon", {
	props: ["href", "title", "target", "size"],
	data: function () {
		let realSize = this.size
		if (realSize == null || realSize.length == 0) {
			realSize = "small"
		}

		return {
			vTitle: (this.title == null) ? "打开链接" : this.title,
			realSize: realSize
		}
	},
	template: `<span><slot></slot>&nbsp;<a :href="href" :title="vTitle" class="link grey" :target="target"><i class="icon linkify" :class="realSize"></i></a></span>`
})

// 带有下划虚线的连接
Vue.component("link-red", {
	props: ["href", "title"],
	data: function () {
		let href = this.href
		if (href == null) {
			href = ""
		}
		return {
			vHref: href
		}
	},
	methods: {
		clickPrevent: function () {
			emitClick(this, arguments)

			if (this.vHref.length > 0) {
				window.location = this.vHref
			}
		}
	},
	template: `<a :href="vHref" :title="title" style="border-bottom: 1px #db2828 dashed" @click.prevent="clickPrevent"><span class="red"><slot></slot></span></a>`
})

// 会弹出窗口的链接
Vue.component("link-popup", {
	props: ["title"],
	methods: {
		clickPrevent: function () {
			emitClick(this, arguments)
		}
	},
	template: `<a href="" :title="title" @click.prevent="clickPrevent"><slot></slot></a>`
})

Vue.component("popup-icon", {
	props: ["title", "href", "height"],
	methods: {
		clickPrevent: function () {
			if (this.href != null && this.href.length > 0) {
				teaweb.popup(this.href, {
					height: this.height
				})
			}
		}
	},
	template: `<span><slot></slot>&nbsp;<a href="" :title="title" @click.prevent="clickPrevent"><i class="icon expand small"></i></a></span>`
})

// 小提示
Vue.component("tip-icon", {
	props: ["content"],
	methods: {
		showTip: function () {
			teaweb.popupTip(this.content)
		}
	},
	template: `<a href="" title="查看帮助" @click.prevent="showTip"><i class="icon question circle grey"></i></a>`
})

// 提交点击事件
function emitClick(obj, arguments) {
	let event = "click"
	let newArgs = [event]
	for (let i = 0; i < arguments.length; i++) {
		newArgs.push(arguments[i])
	}
	obj.$emit.apply(obj, newArgs)
}

Vue.component("countries-selector", {
	props: ["v-countries"],
	data: function () {
		let countries = this.vCountries
		if (countries == null) {
			countries = []
		}
		let countryIds = countries.$map(function (k, v) {
			return v.id
		})
		return {
			countries: countries,
			countryIds: countryIds
		}
	},
	methods: {
		add: function () {
			let countryStringIds = this.countryIds.map(function (v) {
				return v.toString()
			})
			let that = this
			teaweb.popup("/ui/selectCountriesPopup?countryIds=" + countryStringIds.join(","), {
				width: "48em",
				height: "23em",
				callback: function (resp) {
					that.countries = resp.data.countries
					that.change()
				}
			})
		},
		remove: function (index) {
			this.countries.$remove(index)
			this.change()
		},
		change: function () {
			this.countryIds = this.countries.$map(function (k, v) {
				return v.id
			})
		}
	},
	template: `<div>
	<input type="hidden" name="countryIdsJSON" :value="JSON.stringify(countryIds)"/>
	<div v-if="countries.length > 0" style="margin-bottom: 0.5em">
		<div v-for="(country, index) in countries" class="ui label tiny basic">{{country.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("raquo-item", {
	template: `<span class="item disabled" style="padding: 0">&raquo;</span>`
})

Vue.component("bandwidth-size-capacity-view", {
	props: ["v-value"],
	data: function () {
		let capacity = this.vValue
		if (capacity != null && capacity.count > 0 && typeof capacity.unit === "string") {
			capacity.unit = capacity.unit[0].toUpperCase() + capacity.unit.substring(1) + "ps"
		}
		return {
			capacity: capacity
		}
	},
	template: `<span>
	<span v-if="capacity != null && capacity.count > 0">{{capacity.count}}{{capacity.unit}}</span>
</span>`
})

Vue.component("more-options-tbody", {
	data: function () {
		return {
			isVisible: false
		}
	},
	methods: {
		show: function () {
			this.isVisible = !this.isVisible
			this.$emit("change", this.isVisible)
		}
	},
	template: `<tbody>
	<tr>
		<td colspan="2"><a href="" @click.prevent="show()"><span v-if="!isVisible">更多选项</span><span v-if="isVisible">收起选项</span><i class="icon angle" :class="{down:!isVisible, up:isVisible}"></i></a></td>
	</tr>
</tbody>`
})

Vue.component("download-link", {
	props: ["v-element", "v-file", "v-value"],
	created: function () {
		let that = this
		setTimeout(function () {
			that.url = that.composeURL()
		}, 1000)
	},
	data: function () {
		let filename = this.vFile
		if (filename == null || filename.length == 0) {
			filename = "unknown-file"
		}
		return {
			file: filename,
			url: this.composeURL()
		}
	},
	methods: {
		composeURL: function () {
			let text = ""
			if (this.vValue != null) {
				text = this.vValue
			} else {
				let e = document.getElementById(this.vElement)
				if (e == null) {
					// 不提示错误，因为此时可能页面未加载完整
					return
				}
				text = e.innerText
				if (text == null) {
					text = e.textContent
				}
			}
			return Tea.url("/ui/download", {
				file: this.file,
				text: text
			})
		}
	},
	template: `<a :href="url" target="_blank" style="font-weight: normal"><slot></slot></a>`,
})

Vue.component("values-box", {
	props: ["values", "v-values", "size", "maxlength", "name", "placeholder", "v-allow-empty", "validator"],
	data: function () {
		let values = this.values;
		if (values == null) {
			values = [];
		}

		if (this.vValues != null && typeof this.vValues == "object") {
			values = this.vValues
		}

		return {
			"realValues": values,
			"isUpdating": false,
			"isAdding": false,
			"index": 0,
			"value": "",
			isEditing: false
		}
	},
	methods: {
		create: function () {
			this.isAdding = true;
			var that = this;
			setTimeout(function () {
				that.$refs.value.focus();
			}, 200);
		},
		update: function (index) {
			this.cancel()
			this.isUpdating = true;
			this.index = index;
			this.value = this.realValues[index];
			var that = this;
			setTimeout(function () {
				that.$refs.value.focus();
			}, 200);
		},
		confirm: function () {
			if (this.value.length == 0) {
				if (typeof(this.vAllowEmpty) != "boolean" || !this.vAllowEmpty) {
					return
				}
			}

			// validate
			if (typeof(this.validator) == "function") {
				let resp = this.validator.call(this, this.value)
				if (typeof resp == "object") {
					if (typeof resp.isOk == "boolean" && !resp.isOk) {
						if (typeof resp.message == "string") {
							let that = this
							teaweb.warn(resp.message, function () {
								that.$refs.value.focus();
							})
						}
						return
					}
				}
			}

			if (this.isUpdating) {
				Vue.set(this.realValues, this.index, this.value);
			} else {
				this.realValues.push(this.value);
			}
			this.cancel()
			this.$emit("change", this.realValues)
		},
		remove: function (index) {
			this.realValues.$remove(index)
			this.$emit("change", this.realValues)
		},
		cancel: function () {
			this.isUpdating = false;
			this.isAdding = false;
			this.value = "";
		},
		updateAll: function (values) {
			this.realValues = values
		},
		addValue: function (v) {
			this.realValues.push(v)
		},

		startEditing: function () {
			this.isEditing = !this.isEditing
		},
		allValues: function () {
			return this.realValues
		}
	},
	template: `<div>
	<div v-show="!isEditing && realValues.length > 0">
		<div class="ui label tiny basic" v-for="(value, index) in realValues" style="margin-top:0.4em;margin-bottom:0.4em">
			<span v-if="value.toString().length > 0">{{value}}</span>
			<span v-if="value.toString().length == 0" class="disabled">[空]</span>
		</div>
		<a href="" @click.prevent="startEditing" style="font-size: 0.8em; margin-left: 0.2em">[修改]</a>
	</div>
	<div v-show="isEditing || realValues.length == 0">
		<div style="margin-bottom: 1em" v-if="realValues.length > 0">
			<div class="ui label tiny basic" v-for="(value, index) in realValues" style="margin-top:0.4em;margin-bottom:0.4em">
				<span v-if="value.toString().length > 0">{{value}}</span>
				<span v-if="value.toString().length == 0" class="disabled">[空]</span>
				<input type="hidden" :name="name" :value="value"/>
				&nbsp; <a href="" @click.prevent="update(index)" title="修改"><i class="icon pencil small" ></i></a> 
				<a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a> 
			</div> 
			<div class="ui divider"></div>
		</div> 
		<!-- 添加|修改 -->
		<div v-if="isAdding || isUpdating">
			<div class="ui fields inline">
				<div class="ui field">
					<input type="text" :size="size" :maxlength="maxlength" :placeholder="placeholder" v-model="value" ref="value" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
				</div> 
				<div class="ui field">
					<button class="ui button small" type="button" @click.prevent="confirm()">确定</button> 
				</div>
				<div class="ui field">
					<a href="" @click.prevent="cancel()" title="取消"><i class="icon remove small"></i></a> 
				</div> 
			</div> 
		</div> 
		<div v-if="!isAdding && !isUpdating">
			<button class="ui button tiny" type="button" @click.prevent="create()">+</button> 
		</div>
	</div>	
</div>`
});

Vue.component("datetime-input", {
	props: ["v-name", "v-timestamp"],
	mounted: function () {
		let that = this
		teaweb.datepicker(this.$refs.dayInput, function (v) {
			that.day = v
			that.hour = "23"
			that.minute = "59"
			that.second = "59"
			that.change()
		})
	},
	data: function () {
		let timestamp = this.vTimestamp
		if (timestamp != null) {
			timestamp = parseInt(timestamp)
			if (isNaN(timestamp)) {
				timestamp = 0
			}
		} else {
			timestamp = 0
		}

		let day = ""
		let hour = ""
		let minute = ""
		let second = ""

		if (timestamp > 0) {
			let date = new Date()
			date.setTime(timestamp * 1000)

			let year = date.getFullYear().toString()
			let month = this.leadingZero((date.getMonth() + 1).toString(), 2)
			day = year + "-" + month + "-" + this.leadingZero(date.getDate().toString(), 2)

			hour = this.leadingZero(date.getHours().toString(), 2)
			minute = this.leadingZero(date.getMinutes().toString(), 2)
			second = this.leadingZero(date.getSeconds().toString(), 2)
		}

		return {
			timestamp: timestamp,
			day: day,
			hour: hour,
			minute: minute,
			second: second,

			hasDayError: false,
			hasHourError: false,
			hasMinuteError: false,
			hasSecondError: false
		}
	},
	methods: {
		change: function () {
			// day
			if (!/^\d{4}-\d{1,2}-\d{1,2}$/.test(this.day)) {
				this.hasDayError = true
				return
			}
			let pieces = this.day.split("-")
			let year = parseInt(pieces[0])

			let month = parseInt(pieces[1])
			if (month < 1 || month > 12) {
				this.hasDayError = true
				return
			}

			let day = parseInt(pieces[2])
			if (day < 1 || day > 32) {
				this.hasDayError = true
				return
			}

			this.hasDayError = false

			// hour
			if (!/^\d+$/.test(this.hour)) {
				this.hasHourError = true
				return
			}
			let hour = parseInt(this.hour)
			if (isNaN(hour)) {
				this.hasHourError = true
				return
			}
			if (hour < 0 || hour >= 24) {
				this.hasHourError = true
				return
			}
			this.hasHourError = false

			// minute
			if (!/^\d+$/.test(this.minute)) {
				this.hasMinuteError = true
				return
			}
			let minute = parseInt(this.minute)
			if (isNaN(minute)) {
				this.hasMinuteError = true
				return
			}
			if (minute < 0 || minute >= 60) {
				this.hasMinuteError = true
				return
			}
			this.hasMinuteError = false

			// second
			if (!/^\d+$/.test(this.second)) {
				this.hasSecondError = true
				return
			}
			let second = parseInt(this.second)
			if (isNaN(second)) {
				this.hasSecondError = true
				return
			}
			if (second < 0 || second >= 60) {
				this.hasSecondError = true
				return
			}
			this.hasSecondError = false

			let date = new Date(year, month - 1, day, hour, minute, second)
			this.timestamp = Math.floor(date.getTime() / 1000)
		},
		leadingZero: function (s, l) {
			s = s.toString()
			if (l <= s.length) {
				return s
			}
			for (let i = 0; i < l - s.length; i++) {
				s = "0" + s
			}
			return s
		},
		resultTimestamp: function () {
			return this.timestamp
		},
		nextYear: function () {
			let date = new Date()
			date.setFullYear(date.getFullYear()+1)
			this.day = date.getFullYear() + "-" + this.leadingZero(date.getMonth() + 1, 2) + "-" + this.leadingZero(date.getDate(), 2)
			this.hour = this.leadingZero(date.getHours(), 2)
			this.minute = this.leadingZero(date.getMinutes(), 2)
			this.second = this.leadingZero(date.getSeconds(), 2)
			this.change()
		},
		nextDays: function (days) {
			let date = new Date()
			date.setTime(date.getTime() + days * 86400 * 1000)
			this.day = date.getFullYear() + "-" + this.leadingZero(date.getMonth() + 1, 2) + "-" + this.leadingZero(date.getDate(), 2)
			this.hour = this.leadingZero(date.getHours(), 2)
			this.minute = this.leadingZero(date.getMinutes(), 2)
			this.second = this.leadingZero(date.getSeconds(), 2)
			this.change()
		},
		nextHours: function (hours) {
			let date = new Date()
			date.setTime(date.getTime() + hours * 3600 * 1000)
			this.day = date.getFullYear() + "-" + this.leadingZero(date.getMonth() + 1, 2) + "-" + this.leadingZero(date.getDate(), 2)
			this.hour = this.leadingZero(date.getHours(), 2)
			this.minute = this.leadingZero(date.getMinutes(), 2)
			this.second = this.leadingZero(date.getSeconds(), 2)
			this.change()
		}
	},
	template: `<div>
	<input type="hidden" :name="vName" :value="timestamp"/>
	<div class="ui fields inline" style="padding: 0; margin:0">
		<div class="ui field" :class="{error: hasDayError}">
			<input type="text" v-model="day" placeholder="YYYY-MM-DD" style="width:8.6em" maxlength="10" @input="change" ref="dayInput"/>
		</div>
		<div class="ui field" :class="{error: hasHourError}"><input type="text" v-model="hour" maxlength="2" style="width:4em" placeholder="时" @input="change"/></div>
		<div class="ui field">:</div>
		<div class="ui field" :class="{error: hasMinuteError}"><input type="text" v-model="minute" maxlength="2" style="width:4em" placeholder="分" @input="change"/></div>
		<div class="ui field">:</div>
		<div class="ui field" :class="{error: hasSecondError}"><input type="text" v-model="second" maxlength="2" style="width:4em" placeholder="秒" @input="change"/></div>
	</div>
	<p class="comment">常用时间：<a href="" @click.prevent="nextHours(1)"> &nbsp;1小时&nbsp; </a> <span class="disabled">|</span> <a href="" @click.prevent="nextDays(1)"> &nbsp;1天&nbsp; </a> <span class="disabled">|</span> <a href="" @click.prevent="nextDays(3)"> &nbsp;3天&nbsp; </a> <span class="disabled">|</span> <a href="" @click.prevent="nextDays(7)"> &nbsp;1周&nbsp; </a> <span class="disabled">|</span> <a href="" @click.prevent="nextDays(30)"> &nbsp;30天&nbsp; </a> <span class="disabled">|</span> <a href="" @click.prevent="nextYear()"> &nbsp;1年&nbsp; </a> </p>
</div>`
})

// 启用状态标签
Vue.component("label-on", {
	props: ["v-is-on"],
	template: '<div><span v-if="vIsOn" class="green">已启用</span><span v-if="!vIsOn" class="red">已停用</span></div>'
})

// 文字代码标签
Vue.component("code-label", {
	methods: {
		click: function (args) {
			this.$emit("click", args)
		}
	},
	template: `<span class="ui label basic small" style="padding: 3px;margin-left:2px;margin-right:2px" @click.prevent="click"><slot></slot></span>`
})

Vue.component("code-label-plain", {
	template: `<span class="ui label basic tiny" style="padding: 3px;margin-left:2px;margin-right:2px"><slot></slot></span>`
})


// tiny标签
Vue.component("tiny-label", {
	template: `<span class="ui label tiny" style="margin-bottom: 0.5em"><slot></slot></span>`
})

Vue.component("tiny-basic-label", {
	template: `<span class="ui label tiny basic" style="margin-bottom: 0.5em"><slot></slot></span>`
})

// 更小的标签
Vue.component("micro-basic-label", {
	template: `<span class="ui label tiny basic" style="margin-bottom: 0.5em; font-size: 0.7em; padding: 4px"><slot></slot></span>`
})


// 灰色的Label
Vue.component("grey-label", {
	props: ["color"],
	data: function () {
		let color = "grey"
		if (this.color != null && this.color.length > 0) {
			color = "red"
		}
		return {
			labelColor: color
		}
	},
	template: `<span class="ui label basic tiny" :class="labelColor" style="margin-top: 0.4em; font-size: 0.7em; border: 1px solid #ddd!important; font-weight: normal;"><slot></slot></span>`
})

// 可选标签
Vue.component("optional-label", {
	template: `<em><span class="grey">（可选）</span></em>`
})

// Plus专属
Vue.component("plus-label", {
	template: `<span style="color: #B18701;">Plus专属功能。</span>`
})

// 提醒设置项为专业设置
Vue.component("pro-warning-label", {
	template: `<span><i class="icon warning circle yellow"></i>注意：通常不需要修改；如要修改，请在专家指导下进行。</span>`
})


Vue.component("js-page", {
	props: ["v-max"],
	data: function () {
		let max = this.vMax
		if (max == null) {
			max = 0
		}
		return {
			max: max,
			page: 1
		}
	},
	methods: {
		updateMax: function (max) {
			this.max = max
		},
		selectPage: function(page) {
			this.page = page
			this.$emit("change", page)
		}
	},
	template:`<div>
	<div class="page" v-if="max > 1">
		<a href="" v-for="i in max" :class="{active: i == page}" @click.prevent="selectPage(i)">{{i}}</a>
	</div>
</div>`
})

/**
 * 一级菜单
 */
Vue.component("first-menu", {
	props: [],
	template: ' \
		<div class="first-menu"> \
			<div class="ui menu text blue small">\
				<slot></slot>\
			</div> \
			<div class="ui divider"></div> \
		</div>'
});

/**
 * 更多选项
 */
Vue.component("more-options-indicator", {
	props:[],
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		changeVisible: function () {
			this.visible = !this.visible
			if (Tea.Vue != null) {
				Tea.Vue.moreOptionsVisible = this.visible
			}
			this.$emit("change", this.visible)
			this.$emit("input", this.visible)
		}
	},
	template: '<a href="" style="font-weight: normal" @click.prevent="changeVisible()"><slot><span v-if="!visible">更多选项</span><span v-if="visible">收起选项</span></slot> <i class="icon angle" :class="{down:!visible, up:visible}"></i> </a>'
});

Vue.component("page-size-selector", {
	data: function () {
		let query = window.location.search
		let pageSize = 10
		if (query.length > 0) {
			query = query.substr(1)
			let params = query.split("&")
			params.forEach(function (v) {
				let pieces = v.split("=")
				if (pieces.length == 2 && pieces[0] == "pageSize") {
					let pageSizeString = pieces[1]
					if (pageSizeString.match(/^\d+$/)) {
						pageSize = parseInt(pageSizeString, 10)
						if (isNaN(pageSize) || pageSize < 1) {
							pageSize = 10
						}
					}
				}
			})
		}
		return {
			pageSize: pageSize
		}
	},
	watch: {
		pageSize: function () {
			window.ChangePageSize(this.pageSize)
		}
	},
	template: `<select class="ui dropdown" style="height:34px;padding-top:0;padding-bottom:0;margin-left:1em;color:#666" v-model="pageSize">
\t<option value="10">[每页]</option><option value="10" selected="selected">10条</option><option value="20">20条</option><option value="30">30条</option><option value="40">40条</option><option value="50">50条</option><option value="60">60条</option><option value="70">70条</option><option value="80">80条</option><option value="90">90条</option><option value="100">100条</option>
</select>`
})

/**
 * 二级菜单
 */
Vue.component("second-menu", {
	template: ' \
		<div class="second-menu"> \
			<div class="ui menu text blue small">\
				<slot></slot>\
			</div> \
			<div class="ui divider"></div> \
		</div>'
});

Vue.component("loading-message", {
	template: `<div class="ui message loading">
        <div class="ui active inline loader small"></div>  &nbsp; <slot></slot>
    </div>`
})

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

Vue.component("more-options-angle", {
	data: function () {
		return {
			isVisible: false
		}
	},
	methods: {
		show: function () {
			this.isVisible = !this.isVisible
			this.$emit("change", this.isVisible)
		}
	},
	template: `<a href="" @click.prevent="show()"><span v-if="!isVisible">更多选项</span><span v-if="isVisible">收起选项</span><i class="icon angle" :class="{down:!isVisible, up:isVisible}"></i></a>`
})

Vue.component("columns-grid", {
	props: [],
	mounted: function () {
		this.columns = this.calculateColumns()

		let that = this
		window.addEventListener("resize", function () {
			that.columns = that.calculateColumns()
		})
	},
	data: function () {
		return {
			columns: "four"
		}
	},
	methods: {
		calculateColumns: function () {
			let w = window.innerWidth
			let columns = Math.floor(w / 250)
			if (columns == 0) {
				columns = 1
			}

			let columnElements = this.$el.getElementsByClassName("column")
			if (columnElements.length == 0) {
				return
			}
			let maxColumns = columnElements.length
			if (columns > maxColumns) {
				columns = maxColumns
			}

			// 添加右侧边框
			for (let index = 0; index < columnElements.length; index++) {
				let el = columnElements[index]
				el.className = el.className.replace("with-border", "")
				if (index % columns == columns - 1 || index == columnElements.length - 1 /** 最后一个 **/) {
					el.className += " with-border"
				}
			}

			switch (columns) {
				case 1:
					return "one"
				case 2:
					return "two"
				case 3:
					return "three"
				case 4:
					return "four"
				case 5:
					return "five"
				case 6:
					return "six"
				case 7:
					return "seven"
				case 8:
					return "eight"
				case 9:
					return "nine"
				case 10:
					return "ten"
				default:
					return "ten"
			}
		}
	},
	template: `<div :class="'ui ' + columns + ' columns grid counter-chart'">
	<slot></slot>
</div>`
})

/**
 * 菜单项
 */
Vue.component("inner-menu-item", {
	props: ["href", "active", "code"],
	data: function () {
		var active = this.active;
		if (typeof(active) =="undefined") {
			var itemCode = "";
			if (typeof (window.TEA.ACTION.data.firstMenuItem) != "undefined") {
				itemCode = window.TEA.ACTION.data.firstMenuItem;
			}
			active = (itemCode == this.code);
		}
		return {
			vHref: (this.href == null) ? "" : this.href,
			vActive: active
		};
	},
	template: '\
		<a :href="vHref" class="item right" style="color:#4183c4" :class="{active:vActive}">[<slot></slot>]</a> \
		'
});

Vue.component("bandwidth-size-capacity-box", {
	props: ["v-name", "v-value", "v-count", "v-unit", "size", "maxlength", "v-supported-units"],
	data: function () {
		let v = this.vValue
		if (v == null) {
			v = {
				count: this.vCount,
				unit: this.vUnit
			}
		}
		if (v.unit == null || v.unit.length == 0) {
			v.unit = "mb"
		}

		if (typeof (v["count"]) != "number") {
			v["count"] = -1
		}

		let vSize = this.size
		if (vSize == null) {
			vSize = 6
		}

		let vMaxlength = this.maxlength
		if (vMaxlength == null) {
			vMaxlength = 10
		}

		let supportedUnits = this.vSupportedUnits
		if (supportedUnits == null) {
			supportedUnits = []
		}

		return {
			capacity: v,
			countString: (v.count >= 0) ? v.count.toString() : "",
			vSize: vSize,
			vMaxlength: vMaxlength,
			supportedUnits: supportedUnits
		}
	},
	watch: {
		"countString": function (newValue) {
			let value = newValue.trim()
			if (value.length == 0) {
				this.capacity.count = -1
				this.change()
				return
			}
			let count = parseInt(value)
			if (!isNaN(count)) {
				this.capacity.count = count
			}
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.capacity)
		}
	},
	template: `<div class="ui fields inline">
	<input type="hidden" :name="vName" :value="JSON.stringify(capacity)"/>
	<div class="ui field">
		<input type="text" v-model="countString" :maxlength="vMaxlength" :size="vSize"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="capacity.unit" @change="change">
			<option value="b" v-if="supportedUnits.length == 0 || supportedUnits.$contains('b')">Bps</option>
			<option value="kb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('kb')">Kbps</option>
			<option value="mb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('mb')">Mbps</option>
			<option value="gb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('gb')">Gbps</option>
			<option value="tb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('tb')">Tbps</option>
			<option value="pb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('pb')">Pbps</option>
			<option value="eb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('eb')">Ebps</option>
		</select>
	</div>
</div>`
})

Vue.component("health-check-config-box", {
	props: ["v-health-check-config", "v-check-domain-url", "v-is-plus"],
	data: function () {
		let healthCheckConfig = this.vHealthCheckConfig
		let urlProtocol = "http"
		let urlPort = ""
		let urlRequestURI = "/"
		let urlHost = ""

		if (healthCheckConfig == null) {
			healthCheckConfig = {
				isOn: false,
				url: "",
				interval: {count: 60, unit: "second"},
				statusCodes: [200],
				timeout: {count: 10, unit: "second"},
				countTries: 3,
				tryDelay: {count: 100, unit: "ms"},
				autoDown: true,
				countUp: 1,
				countDown: 3,
				userAgent: "",
				onlyBasicRequest: true,
				accessLogIsOn: true
			}
			let that = this
			setTimeout(function () {
				that.changeURL()
			}, 500)
		} else {
			try {
				let url = new URL(healthCheckConfig.url)
				urlProtocol = url.protocol.substring(0, url.protocol.length - 1)

				// 域名
				urlHost = url.host
				if (urlHost == "%24%7Bhost%7D") {
					urlHost = "${host}"
				}
				let colonIndex = urlHost.indexOf(":")
				if (colonIndex > 0) {
					urlHost = urlHost.substring(0, colonIndex)
				}

				urlPort = url.port
				urlRequestURI = url.pathname
				if (url.search.length > 0) {
					urlRequestURI += url.search
				}
			} catch (e) {
			}

			if (healthCheckConfig.statusCodes == null) {
				healthCheckConfig.statusCodes = [200]
			}
			if (healthCheckConfig.interval == null) {
				healthCheckConfig.interval = {count: 60, unit: "second"}
			}
			if (healthCheckConfig.timeout == null) {
				healthCheckConfig.timeout = {count: 10, unit: "second"}
			}
			if (healthCheckConfig.tryDelay == null) {
				healthCheckConfig.tryDelay = {count: 100, unit: "ms"}
			}
			if (healthCheckConfig.countUp == null || healthCheckConfig.countUp < 1) {
				healthCheckConfig.countUp = 1
			}
			if (healthCheckConfig.countDown == null || healthCheckConfig.countDown < 1) {
				healthCheckConfig.countDown = 3
			}
		}

		return {
			healthCheck: healthCheckConfig,
			advancedVisible: false,
			urlProtocol: urlProtocol,
			urlHost: urlHost,
			urlPort: urlPort,
			urlRequestURI: urlRequestURI,
			urlIsEditing: healthCheckConfig.url.length == 0,

			hostErr: ""
		}
	},
	watch: {
		urlRequestURI: function () {
			if (this.urlRequestURI.length > 0 && this.urlRequestURI[0] != "/") {
				this.urlRequestURI = "/" + this.urlRequestURI
			}
			this.changeURL()
		},
		urlPort: function (v) {
			let port = parseInt(v)
			if (!isNaN(port)) {
				this.urlPort = port.toString()
			} else {
				this.urlPort = ""
			}
			this.changeURL()
		},
		urlProtocol: function () {
			this.changeURL()
		},
		urlHost: function () {
			this.changeURL()
			this.hostErr = ""
		},
		"healthCheck.countTries": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countTries = count
			} else {
				this.healthCheck.countTries = 0
			}
		},
		"healthCheck.countUp": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countUp = count
			} else {
				this.healthCheck.countUp = 0
			}
		},
		"healthCheck.countDown": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countDown = count
			} else {
				this.healthCheck.countDown = 0
			}
		}
	},
	methods: {
		showAdvanced: function () {
			this.advancedVisible = !this.advancedVisible
		},
		changeURL: function () {
			let urlHost = this.urlHost
			if (urlHost.length == 0) {
				urlHost = "${host}"
			}
			this.healthCheck.url = this.urlProtocol + "://" + urlHost + ((this.urlPort.length > 0) ? ":" + this.urlPort : "") + this.urlRequestURI
		},
		changeStatus: function (values) {
			this.healthCheck.statusCodes = values.$map(function (k, v) {
				let status = parseInt(v)
				if (isNaN(status)) {
					return 0
				} else {
					return status
				}
			})
		},
		onChangeURLHost: function () {
			let checkDomainURL = this.vCheckDomainUrl
			if (checkDomainURL == null || checkDomainURL.length == 0) {
				return
			}

			let that = this
			Tea.action(checkDomainURL)
				.params({host: this.urlHost})
				.success(function (resp) {
					if (!resp.data.isOk) {
						that.hostErr = "在当前集群中找不到此域名，可能会影响健康检查结果。"
					} else {
						that.hostErr = ""
					}
				})
				.post()
		},
		editURL: function () {
			this.urlIsEditing = !this.urlIsEditing
		}
	},
	template: `<div>
<input type="hidden" name="healthCheckJSON" :value="JSON.stringify(healthCheck)"/>
<table class="ui table definition selectable">
	<tbody>
		<tr>
			<td class="title">启用健康检查</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.isOn"/>
					<label></label>
				</div>
				<p class="comment">通过访问节点上的网站URL来确定节点是否健康。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td>检测URL *</td>
			<td>
				<div v-if="healthCheck.url.length > 0" style="margin-bottom: 1em"><code-label>{{healthCheck.url}}</code-label> &nbsp; <a href="" @click.prevent="editURL"><span class="small">修改 <i class="icon angle" :class="{down: !urlIsEditing, up: urlIsEditing}"></i></span></a> </div>
				<div v-show="urlIsEditing">
					<table class="ui table">
						 <tr>
							<td class="title">协议</td> 
							<td>
								<select class="ui dropdown auto-width" v-model="urlProtocol">
								<option value="http">http://</option>
								<option value="https">https://</option>
								</select>
							</td>
						</tr>
						<tr>
							<td>域名</td>
							<td>
								<input type="text" v-model="urlHost" @change="onChangeURLHost"/>
								<p class="comment"><span v-if="hostErr.length > 0" class="red">{{hostErr}}</span>已经部署到当前集群的一个域名；如果为空则使用节点IP作为域名。<span class="red" v-if="urlProtocol == 'https' && urlHost.length == 0">如果协议是https，这里必须填写一个已经设置了SSL证书的域名。</span></p>
							</td>
						</tr>
						<tr>
							<td>端口</td>
							<td>
								<input type="text" maxlength="5" style="width:5.4em" placeholder="端口" v-model="urlPort"/>
								<p class="comment">域名或者IP的端口，可选项，默认为80/443。</p>
							</td>
						</tr>
						<tr>
							<td>RequestURI</td>
							<td><input type="text" v-model="urlRequestURI" placeholder="/" style="width:20em"/>
								<p class="comment">请求的路径，可以带参数，可选项。</p>
							</td>
						</tr>
					</table>
					<div class="ui divider"></div>
					<p class="comment" v-if="healthCheck.url.length > 0">拼接后的检测URL：<code-label>{{healthCheck.url}}</code-label>，其中\${host}指的是域名。</p>
				</div>
			</td>
		</tr>
		<tr>
			<td>检测时间间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.interval"></time-duration-box>
				<p class="comment">两次检查之间的间隔。</p>
			</td>
		</tr>
		<tr>
			<td>自动上/下线<span v-if="vIsPlus">IP</span></td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.autoDown"/>
					<label></label>
				</div>
				<p class="comment">选中后系统会根据健康检查的结果自动标记<span v-if="vIsPlus">节点IP</span><span v-else>节点</span>的上线/下线状态，并可能自动同步DNS设置。<span v-if="!vIsPlus">注意：免费版的只能整体上下线整个节点，商业版的可以下线单个IP。</span></p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续上线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countUp" style="width:5em" maxlength="6"/>
				<p class="comment">连续{{healthCheck.countUp}}次检查成功后自动恢复上线。</p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续下线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countDown" style="width:5em" maxlength="6"/>
				<p class="comment">连续{{healthCheck.countDown}}次检查失败后自动下线。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td colspan="2"><more-options-angle @change="showAdvanced"></more-options-angle></td>
		</tr>
	</tbody>
	<tbody v-show="advancedVisible && healthCheck.isOn">
		<tr>
			<td>允许的状态码</td>
			<td>
				<values-box :values="healthCheck.statusCodes" maxlength="3" @change="changeStatus"></values-box>
				<p class="comment">允许检测URL返回的状态码列表。</p>
			</td>
		</tr>
		<tr>
			<td>超时时间</td>
			<td>
				<time-duration-box :v-value="healthCheck.timeout"></time-duration-box>
				<p class="comment">读取检测URL超时时间。</p>
			</td>	
		</tr>
		<tr>
			<td>连续尝试次数</td>
			<td>
				<input type="text" v-model="healthCheck.countTries" style="width: 5em" maxlength="2"/>
				<p class="comment">如果读取检测URL失败后需要再次尝试的次数。</p>
			</td>
		</tr>
		<tr>
			<td>每次尝试间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.tryDelay"></time-duration-box>
				<p class="comment">如果读取检测URL失败后再次尝试时的间隔时间。</p>
			</td>
		</tr>
		<tr>
			<td>终端信息<em>（User-Agent）</em></td>
			<td>
				<input type="text" v-model="healthCheck.userAgent" maxlength="200"/>
				<p class="comment">发送到服务器的User-Agent值，不填写表示使用默认值。</p>
			</td>
		</tr>
		<tr>
			<td>只基础请求</td>
			<td>
				<checkbox v-model="healthCheck.onlyBasicRequest"></checkbox>
				<p class="comment">只做基础的请求，不处理反向代理（不检查源站）、WAF等。</p>
			</td>
		</tr>
		<tr>
			<td>记录访问日志</td>
			<td>
				<checkbox v-model="healthCheck.accessLogIsOn"></checkbox>
				<p class="comment">记录健康检查的访问日志。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`
})

// 将变量转换为中文
Vue.component("request-variables-describer", {
	data: function () {
		return {
			vars:[]
		}
	},
	methods: {
		update: function (variablesString) {
			this.vars = []
			let that = this
			variablesString.replace(/\${.+?}/g, function (v) {
				let def = that.findVar(v)
				if (def == null) {
					return v
				}
				that.vars.push(def)
			})
		},
		findVar: function (name) {
			let def = null
			window.REQUEST_VARIABLES.forEach(function (v) {
				if (v.code == name) {
					def = v
				}
			})
			return def
		}
	},
	template: `<span>
	<span v-for="(v, index) in vars"><code-label :title="v.description">{{v.code}}</code-label> - {{v.name}}<span v-if="index < vars.length-1">；</span></span>
</span>`
})


Vue.component("combo-box", {
	// data-url 和 data-key 成对出现
	props: [
		"name", "title", "placeholder", "size", "v-items", "v-value",
		"data-url", // 数据源URL
		"data-key", // 数据源中数据的键名
		"data-search", // 是否启用动态搜索，如果值为on或true，则表示启用
		"width"
	],
	mounted: function () {
		if (this.dataURL.length > 0) {
			this.search("")
		}

		// 设定菜单宽度
		let searchBox = this.$refs.searchBox
		if (searchBox != null) {
			let inputWidth = searchBox.offsetWidth
			if (inputWidth != null && inputWidth > 0) {
				this.$refs.menu.style.width = inputWidth + "px"
			} else if (this.styleWidth.length > 0) {
				this.$refs.menu.style.width = this.styleWidth
			}
		}
	},
	data: function () {
		let items = this.vItems
		if (items == null || !(items instanceof Array)) {
			items = []
		}
		items = this.formatItems(items)

		// 当前选中项
		let selectedItem = null
		if (this.vValue != null) {
			let that = this
			items.forEach(function (v) {
				if (v.value == that.vValue) {
					selectedItem = v
				}
			})
		}

		let width = this.width
		if (width == null || width.length == 0) {
			width = "11em"
		} else {
			if (/\d+$/.test(width)) {
				width += "em"
			}
		}

		// data url
		let dataURL = ""
		if (typeof this.dataUrl == "string" && this.dataUrl.length > 0) {
			dataURL = this.dataUrl
		}

		return {
			allItems: items, // 原始的所有的items
			items: items.$copy(), // 候选的items
			selectedItem: selectedItem, // 选中的item
			keyword: "",
			visible: false,
			hideTimer: null,
			hoverIndex: 0,
			styleWidth: width,

			isInitial: true,
			dataURL: dataURL,
			urlRequestId: 0 // 记录URL请求ID，防止并行冲突
		}
	},
	methods: {
		search: function (keyword) {
			// 从URL中获取选项数据
			let dataUrl = this.dataURL
			let dataKey = this.dataKey
			let that = this

			let requestId = Math.random()
			this.urlRequestId = requestId

			Tea.action(dataUrl)
				.params({
					keyword: (keyword == null) ? "" : keyword
				})
				.post()
				.success(function (resp) {
					if (requestId != that.urlRequestId) {
						return
					}

					if (resp.data != null) {
						if (typeof (resp.data[dataKey]) == "object") {
							let items = that.formatItems(resp.data[dataKey])
							that.allItems = items
							that.items = items.$copy()

							if (that.isInitial) {
								that.isInitial = false
								if (that.vValue != null) {
									items.forEach(function (v) {
										if (v.value == that.vValue) {
											that.selectedItem = v
										}
									})
								}
							}
						}
					}
				})
		},
		formatItems: function (items) {
			items.forEach(function (v) {
				if (v.value == null) {
					v.value = v.id
				}
			})
			return items
		},
		reset: function () {
			this.selectedItem = null
			this.change()
			this.hoverIndex = 0

			let that = this
			setTimeout(function () {
				if (that.$refs.searchBox) {
					that.$refs.searchBox.focus()
				}
			})
		},
		clear: function () {
			this.selectedItem = null
			this.change()
			this.hoverIndex = 0
		},
		changeKeyword: function () {
			let shouldSearch = this.dataURL.length > 0 && (this.dataSearch == "on" || this.dataSearch == "true")

			this.hoverIndex = 0
			let keyword = this.keyword
			if (keyword.length == 0) {
				if (shouldSearch) {
					this.search(keyword)
				} else {
					this.items = this.allItems.$copy()
				}
				return
			}


			if (shouldSearch) {
				this.search(keyword)
			} else {
				this.items = this.allItems.$copy().filter(function (v) {
					if (v.fullname != null && v.fullname.length > 0 && teaweb.match(v.fullname, keyword)) {
						return true
					}
					return teaweb.match(v.name, keyword)
				})
			}
		},
		selectItem: function (item) {
			this.selectedItem = item
			this.change()
			this.hoverIndex = 0
			this.keyword = ""
			this.changeKeyword()
		},
		confirm: function () {
			if (this.items.length > this.hoverIndex) {
				this.selectItem(this.items[this.hoverIndex])
			}
		},
		show: function () {
			this.visible = true

			// 不要重置hoverIndex，以便焦点可以在输入框和可选项之间切换
		},
		hide: function () {
			let that = this
			this.hideTimer = setTimeout(function () {
				that.visible = false
			}, 500)
		},
		downItem: function () {
			this.hoverIndex++
			if (this.hoverIndex > this.items.length - 1) {
				this.hoverIndex = 0
			}
			this.focusItem()
		},
		upItem: function () {
			this.hoverIndex--
			if (this.hoverIndex < 0) {
				this.hoverIndex = 0
			}
			this.focusItem()
		},
		focusItem: function () {
			if (this.hoverIndex < this.items.length) {
				this.$refs.itemRef[this.hoverIndex].focus()
				let that = this
				setTimeout(function () {
					that.$refs.searchBox.focus()
					if (that.hideTimer != null) {
						clearTimeout(that.hideTimer)
						that.hideTimer = null
					}
				})
			}
		},
		change: function () {
			this.$emit("change", this.selectedItem)

			let that = this
			setTimeout(function () {
				if (that.$refs.selectedLabel != null) {
					that.$refs.selectedLabel.focus()
				}
			})
		},
		submitForm: function (event) {
			if (event.target.tagName != "A") {
				return
			}
			let parentBox = this.$refs.selectedLabel.parentNode
			while (true) {
				parentBox = parentBox.parentNode
				if (parentBox == null || parentBox.tagName == "BODY") {
					return
				}
				if (parentBox.tagName == "FORM") {
					parentBox.submit()
					break
				}
			}
		},

		setDataURL: function (dataURL) {
			this.dataURL = dataURL
		},
		reloadData: function () {
			this.search("")
		}
	},
	template: `<div style="display: inline; z-index: 10; background: white" class="combo-box">
	<!-- 搜索框 -->
	<div v-if="selectedItem == null">
		<input type="text" v-model="keyword" :placeholder="placeholder" :size="size" :style="{'width': styleWidth}"  @input="changeKeyword" @focus="show" @blur="hide" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="searchBox" @keydown.down.prevent="downItem" @keydown.up.prevent="upItem"/>
	</div>
	
	<!-- 当前选中 -->
	<div v-if="selectedItem != null">
		<input type="hidden" :name="name" :value="selectedItem.value"/>
		<span class="ui label basic" style="line-height: 1.4; font-weight: normal; font-size: 1em" ref="selectedLabel"><span><span v-if="title != null && title.length > 0">{{title}}：</span>{{selectedItem.name}}</span>
			<a href="" title="清除" @click.prevent="reset"><i class="icon remove small"></i></a>
		</span>
	</div>
	
	<!-- 菜单 -->
	<div v-show="selectedItem == null && items.length > 0 && visible">
		<div class="ui menu vertical small narrow-scrollbar" ref="menu">
			<a href="" v-for="(item, index) in items" ref="itemRef" class="item" :class="{active: index == hoverIndex, blue: index == hoverIndex}" @click.prevent="selectItem(item)" style="line-height: 1.4">
				<span v-if="item.fullname != null && item.fullname.length > 0">{{item.fullname}}</span>
				<span v-else>{{item.name}}</span>
			</a>
		</div>
	</div>
</div>`
})

Vue.component("search-box", {
	props: ["placeholder", "width"],
	data: function () {
		let width = this.width
		if (width == null) {
			width = "10em"
		}
		return {
			realWidth: width,
			realValue: ""
		}
	},
	methods: {
		onInput: function () {
			this.$emit("input", { value: this.realValue})
			this.$emit("change", { value: this.realValue})
		},
		clearValue: function () {
			this.realValue = ""
			this.focus()
			this.onInput()
		},
		focus: function () {
			this.$refs.valueRef.focus()
		}
	},
	template: `<div>
	<div class="ui input small" :class="{'right labeled': realValue.length > 0}">
		<input type="text" :placeholder="placeholder" :style="{width: realWidth}" @input="onInput" v-model="realValue" ref="valueRef"/>
		<a href="" class="ui label blue" v-if="realValue.length > 0" @click.prevent="clearValue" style="padding-right: 0"><i class="icon remove"></i></a>
	</div>
</div>`
})

Vue.component("dot", {
	template: '<span style="display: inline-block; padding-bottom: 3px"><i class="icon circle tiny"></i></span>'
})

Vue.component("time-duration-box", {
	props: ["v-name", "v-value", "v-count", "v-unit", "placeholder", "v-min-unit", "maxlength"],
	mounted: function () {
		this.change()
	},
	data: function () {
		let v = this.vValue
		if (v == null) {
			v = {
				count: this.vCount,
				unit: this.vUnit
			}
		}
		if (typeof (v["count"]) != "number") {
			v["count"] = -1
		}

		let minUnit = this.vMinUnit
		let units = [
			{
				code: "ms",
				name: "毫秒"
			},
			{
				code: "second",
				name: "秒"
			},
			{
				code: "minute",
				name: "分钟"
			},
			{
				code: "hour",
				name: "小时"
			},
			{
				code: "day",
				name: "天"
			}
		]
		let minUnitIndex = -1
		if (minUnit != null && typeof minUnit == "string" && minUnit.length > 0) {
			for (let i = 0; i < units.length; i++) {
				if (units[i].code == minUnit) {
					minUnitIndex = i
					break
				}
			}
		}
		if (minUnitIndex > -1) {
			units = units.slice(minUnitIndex)
		}

		let maxLength = parseInt(this.maxlength)
		if (typeof maxLength != "number") {
			maxLength = 10
		}

		return {
			duration: v,
			countString: (v.count >= 0) ? v.count.toString() : "",
			units: units,
			realMaxLength: maxLength
		}
	},
	watch: {
		"countString": function (newValue) {
			let value = newValue.trim()
			if (value.length == 0) {
				this.duration.count = -1
				return
			}
			let count = parseInt(value)
			if (!isNaN(count)) {
				this.duration.count = count
			}
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.duration)
		}
	},
	template: `<div class="ui fields inline" style="padding-bottom: 0; margin-bottom: 0">
	<input type="hidden" :name="vName" :value="JSON.stringify(duration)"/>
	<div class="ui field">
		<input type="text" v-model="countString" :maxlength="realMaxLength" :size="realMaxLength" :placeholder="placeholder" @keypress.enter.prevent="1"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="duration.unit" @change="change">
			<option v-for="unit in units" :value="unit.code">{{unit.name}}</option>
		</select>
	</div>
</div>`
})

Vue.component("time-duration-text", {
	props: ["v-value"],
	methods: {
		unitName: function (unit) {
			switch (unit) {
				case "ms":
					return "毫秒"
				case "second":
					return "秒"
				case "minute":
					return "分钟"
				case "hour":
					return "小时"
				case "day":
					return "天"
			}
		}
	},
	template: `<span>
	{{vValue.count}} {{unitName(vValue.unit)}}
</span>`
})

Vue.component("not-found-box", {
	props: ["message"],
	template: `<div style="text-align: center; margin-top: 5em;">
	<div style="font-size: 2em; margin-bottom: 1em"><i class="icon exclamation triangle large grey"></i></div>
	<p class="comment">{{message}}<slot></slot></p>
</div>`
})

// 警告消息
Vue.component("warning-message", {
	template: `<div class="ui icon message warning"><i class="icon warning circle"></i><div class="content"><slot></slot></div></div>`
})

let checkboxId = 0
Vue.component("checkbox", {
	props: ["name", "value", "v-value", "id", "checked"],
	data: function () {
		checkboxId++
		let elementId = this.id
		if (elementId == null) {
			elementId = "checkbox" + checkboxId
		}

		let elementValue = this.vValue
		if (elementValue == null) {
			elementValue = "1"
		}

		let checkedValue = this.value
		if (checkedValue == null && this.checked == "checked") {
			checkedValue = elementValue
		}

		return {
			elementId: elementId,
			elementValue: elementValue,
			newValue: checkedValue
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.newValue)
		},
		check: function () {
			this.newValue = this.elementValue
		},
		uncheck: function () {
			this.newValue = ""
		},
		isChecked: function () {
			return (typeof (this.newValue) == "boolean" && this.newValue) || this.newValue == this.elementValue
		}
	},
	watch: {
		value: function (v) {
			if (typeof v == "boolean") {
				this.newValue = v
			}
		}
	},
	template: `<div class="ui checkbox">
	<input type="checkbox" :name="name" :value="elementValue" :id="elementId" @change="change" v-model="newValue"/>
	<label :for="elementId"><slot></slot></label>
</div>`
})

Vue.component("network-addresses-view", {
	props: ["v-addresses"],
	template: `<div>
	<div class="ui label tiny basic" v-if="vAddresses != null" v-for="addr in vAddresses">
		{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host.quoteIP()}}</span><span v-else>*</span>:{{addr.portRange}}
	</div>
</div>`
})

Vue.component("url-patterns-box", {
	props: ["value"],
	data: function () {
		let patterns = []
		if (this.value != null) {
			patterns = this.value
		}

		return {
			patterns: patterns,
			isAdding: false,

			addingPattern: {"type": "wildcard", "pattern": ""},
			editingIndex: -1,

			patternIsInvalid: false,

			windowIsSmall: window.innerWidth < 600
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.patternInput.focus()
			})
		},
		edit: function (index) {
			this.isAdding = true
			this.editingIndex = index
			this.addingPattern = {
				type: this.patterns[index].type,
				pattern: this.patterns[index].pattern
			}
		},
		confirm: function () {
			if (this.requireURL(this.addingPattern.type)) {
				let pattern = this.addingPattern.pattern.trim()
				if (pattern.length == 0) {
					let that = this
					teaweb.warn("请输入URL", function () {
						that.$refs.patternInput.focus()
					})
					return
				}
			}
			if (this.editingIndex < 0) {
				this.patterns.push({
					type: this.addingPattern.type,
					pattern: this.addingPattern.pattern
				})
			} else {
				this.patterns[this.editingIndex].type = this.addingPattern.type
				this.patterns[this.editingIndex].pattern = this.addingPattern.pattern
			}
			this.notifyChange()
			this.cancel()
		},
		remove: function (index) {
			this.patterns.$remove(index)
			this.cancel()
			this.notifyChange()
		},
		cancel: function () {
			this.isAdding = false
			this.addingPattern = {"type": "wildcard", "pattern": ""}
			this.editingIndex = -1
		},
		patternTypeName: function (patternType) {
			switch (patternType) {
				case "wildcard":
					return "通配符"
				case "regexp":
					return "正则"
				case "images":
					return "常见图片文件"
				case "audios":
					return "常见音频文件"
				case "videos":
					return "常见视频文件"
			}
			return ""
		},
		notifyChange: function () {
			this.$emit("input", this.patterns)
		},
		changePattern: function () {
			this.patternIsInvalid = false
			let pattern = this.addingPattern.pattern
			switch (this.addingPattern.type) {
				case "wildcard":
					if (pattern.indexOf("?") >= 0) {
						this.patternIsInvalid = true
					}
					break
				case "regexp":
					if (pattern.indexOf("?") >= 0) {
						let pieces = pattern.split("?")
						for (let i = 0; i < pieces.length - 1; i++) {
							if (pieces[i].length == 0 || pieces[i][pieces[i].length - 1] != "\\") {
								this.patternIsInvalid = true
							}
						}
					}
					break
			}
		},
		requireURL: function (patternType) {
			return patternType == "wildcard" || patternType == "regexp"
		}
	},
	template: `<div>
	<div v-show="patterns.length > 0">
		<div v-for="(pattern, index) in patterns" class="ui label basic small" :class="{blue: index == editingIndex, disabled: isAdding && index != editingIndex}" style="margin-bottom: 0.8em">
			<span class="grey" style="font-weight: normal">[{{patternTypeName(pattern.type)}}]</span> <span >{{pattern.pattern}}</span> &nbsp; 
			<a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil tiny"></i></a> 
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-show="isAdding" style="margin-top: 0.5em">
		<div :class="{'ui fields inline': !windowIsSmall}">
			<div class="ui field">
				<select class="ui dropdown auto-width" v-model="addingPattern.type">
					<option value="wildcard">通配符</option>
					<option value="regexp">正则表达式</option>
					<option value="images">常见图片</option>
					<option value="audios">常见音频</option>
					<option value="videos">常见视频</option>
				</select>
			</div>
			<div class="ui field" v-show="addingPattern.type == 'wildcard' || addingPattern.type ==  'regexp'">
				<input type="text" :placeholder="(addingPattern.type == 'wildcard') ? '可以使用星号（*）通配符，不区分大小写' : '可以使用正则表达式，不区分大小写'" v-model="addingPattern.pattern" @input="changePattern" size="36" ref="patternInput" @keyup.enter="confirm()" @keypress.enter.prevent="1" spellcheck="false"/>
				<p class="comment" v-if="patternIsInvalid"><span class="red" style="font-weight: normal"><span v-if="addingPattern.type == 'wildcard'">通配符</span><span v-if="addingPattern.type == 'regexp'">正则表达式</span>中不能包含问号（?）及问号以后的内容。</span></p>
			</div>
			<div class="ui field" style="padding-left: 0"  v-show="addingPattern.type == 'wildcard' || addingPattern.type ==  'regexp'">
				<tip-icon content="通配符示例：<br/>单个路径开头：/hello/world/*<br/>单个路径结尾：*/hello/world<br/>包含某个路径：*/article/*<br/>某个域名下的所有URL：*example.com/*<br/>忽略某个扩展名：*.js" v-if="addingPattern.type == 'wildcard'"></tip-icon>
				<tip-icon content="正则表达式示例：<br/>单个路径开头：^/hello/world<br/>单个路径结尾：/hello/world$<br/>包含某个路径：/article/<br/>匹配某个数字路径：/article/(\\d+)<br/>某个域名下的所有URL：^(http|https)://example.com/" v-if="addingPattern.type == 'regexp'"></tip-icon>
			</div>
			<div class="ui field">
				<button class="ui button tiny" :class="{disabled:this.patternIsInvalid}" type="button" @click.prevent="confirm">确定</button><a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<div v-if=!isAdding style="margin-top: 0.5em">
		<button class="ui button tiny basic" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("size-capacity-view", {
	props:["v-default-text", "v-value"],
	methods: {
		composeCapacity: function (capacity) {
			return teaweb.convertSizeCapacityToString(capacity)
		}
	},
	template: `<div>
	<span v-if="vValue != null && vValue.count > 0">{{composeCapacity(vValue)}}</span>
	<span v-else>{{vDefaultText}}</span>
</div>`
})

// 信息提示窗口
Vue.component("tip-message-box", {
	props: ["code"],
	mounted: function () {
		let that = this
		Tea.action("/ui/showTip")
			.params({
				code: this.code
			})
			.success(function (resp) {
				that.visible = resp.data.visible
			})
			.post()
	},
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		close: function () {
			this.visible = false
			Tea.action("/ui/hideTip")
				.params({
					code: this.code
				})
				.post()
		}
	},
	template: `<div class="ui icon message" v-if="visible">
	<i class="icon info circle"></i>
	<i class="close icon" title="取消" @click.prevent="close" style="margin-top: 1em"></i>
	<div class="content">
		<slot></slot>
	</div>
</div>`
})

Vue.component("digit-input", {
	props: ["value", "maxlength", "size", "min", "max", "required", "placeholder"],
	mounted: function () {
		let that = this
		setTimeout(function () {
			that.check()
		})
	},
	data: function () {
		let realMaxLength = this.maxlength
		if (realMaxLength == null) {
			realMaxLength = 20
		}

		let realSize = this.size
		if (realSize == null) {
			realSize = 6
		}

		return {
			realValue: this.value,
			realMaxLength: realMaxLength,
			realSize: realSize,
			isValid: true
		}
	},
	watch: {
		realValue: function (v) {
			this.notifyChange()
		}
	},
	methods: {
		notifyChange: function () {
			let v = parseInt(this.realValue.toString(), 10)
			if (isNaN(v)) {
				v = 0
			}
			this.check()
			this.$emit("input", v)
		},
		check: function () {
			if (this.realValue == null) {
				return
			}
			let s = this.realValue.toString()
			if (!/^\d+$/.test(s)) {
				this.isValid = false
				return
			}
			let v = parseInt(s, 10)
			if (isNaN(v)) {
				this.isValid = false
			} else {
				if (this.required) {
					this.isValid = (this.min == null || this.min <= v) && (this.max == null || this.max >= v)
				} else {
					this.isValid = (v == 0 || (this.min == null || this.min <= v) && (this.max == null || this.max >= v))
				}
			}
		}
	},
	template: `<input type="text" v-model="realValue" :maxlength="realMaxLength" :size="realSize" :class="{error: !this.isValid}" :placeholder="placeholder" autocomplete="off"/>`
})

Vue.component("keyword", {
	props: ["v-word"],
	data: function () {
		let word = this.vWord
		if (word == null) {
			word = ""
		} else {
			word = word.replace(/\)/g, "\\)")
			word = word.replace(/\(/g, "\\(")
			word = word.replace(/\+/g, "\\+")
			word = word.replace(/\^/g, "\\^")
			word = word.replace(/\$/g, "\\$")
			word = word.replace(/\?/g, "\\?")
			word = word.replace(/\*/g, "\\*")
			word = word.replace(/\[/g, "\\[")
			word = word.replace(/{/g, "\\{")
			word = word.replace(/\./g, "\\.")
		}

		let slot = this.$slots["default"][0]
		let text = slot.text
		if (word.length > 0) {
			let that = this
			let m = []  // replacement => tmp
			let tmpIndex = 0
			text = text.replaceAll(new RegExp("(" + word + ")", "ig"), function (replacement) {
				tmpIndex++
				let s = "<span style=\"border: 1px #ccc dashed; color: #ef4d58\">" + that.encodeHTML(replacement) + "</span>"
				let tmpKey = "$TMP__KEY__" + tmpIndex.toString() + "$"
				m.push([tmpKey, s])
				return tmpKey
			})
			text = this.encodeHTML(text)

			m.forEach(function (r) {
				text = text.replace(r[0], r[1])
			})

		} else {
			text = this.encodeHTML(text)
		}

		return {
			word: word,
			text: text
		}
	},
	methods: {
		encodeHTML: function (s) {
			s = s.replace(/&/g, "&amp;")
			s = s.replace(/</g, "&lt;")
			s = s.replace(/>/g, "&gt;")
			s = s.replace(/"/g, "&quot;")
			return s
		}
	},
	template: `<span><span style="display: none"><slot></slot></span><span v-html="text"></span></span>`
})

Vue.component("bits-var", {
	props: ["v-bits"],
	data: function () {
		let bits = this.vBits
		if (typeof bits != "number") {
			bits = 0
		}
		let format = teaweb.splitFormat(teaweb.formatBits(bits))
		return {
			format: format
		}
	},
	template:`<var class="normal">
	<span>{{format[0]}}</span>{{format[1]}}
</var>`
})

Vue.component("mask-warning", {
	template: `<span class="red">为了安全起见，此项数据保存后将不允许在界面查看完整明文，为避免忘记，请自行记录原始数据。</span>`
})

Vue.component("chart-columns-grid", {
	props: [],
	mounted: function () {
		this.columns = this.calculateColumns()

		let that = this
		window.addEventListener("resize", function () {
			that.columns = that.calculateColumns()
		})
	},
	updated: function () {
		let totalElements = this.$el.getElementsByClassName("column").length
		if (totalElements == this.totalElements) {
			return
		}
		this.totalElements = totalElements
		this.calculateColumns()
	},
	data: function () {
		return {
			columns: "four",
			totalElements: 0
		}
	},
	methods: {
		calculateColumns: function () {
			let w = window.innerWidth
			let columns = Math.floor(w / 500)
			if (columns == 0) {
				columns = 1
			}

			let columnElements = this.$el.getElementsByClassName("column")
			if (columnElements.length == 0) {
				return "one"
			}
			let maxColumns = columnElements.length
			if (columns > maxColumns) {
				columns = maxColumns
			}

			// 添加右侧边框
			for (let index = 0; index < columnElements.length; index++) {
				let el = columnElements[index]
				el.className = el.className.replace("with-border", "")
				if (index % columns == columns - 1 || index == columnElements.length - 1 /** 最后一个 **/) {
					el.className += " with-border"
				}
			}

			switch (columns) {
				case 1:
					return "one"
				case 2:
					return "two"
				case 3:
					return "three"
				case 4:
					return "four"
				case 5:
					return "five"
				case 6:
					return "six"
				case 7:
					return "seven"
				case 8:
					return "eight"
				case 9:
					return "nine"
				case 10:
					return "ten"
				default:
					return "ten"
			}
		}
	},
	template: `<div :class="'ui ' + columns + ' columns grid chart-grid'">
	<slot></slot>
</div>`
})

Vue.component("bytes-var", {
	props: ["v-bytes"],
	data: function () {
		let bytes = this.vBytes
		if (typeof bytes != "number") {
			bytes = 0
		}
		let format = teaweb.splitFormat(teaweb.formatBytes(bytes))
		return {
			format: format
		}
	},
	template:`<var class="normal">
	<span>{{format[0]}}</span>{{format[1]}}
</var>`
})

Vue.component("node-log-row", {
	props: ["v-log", "v-keyword"],
	data: function () {
		return {
			log: this.vLog,
			keyword: this.vKeyword
		}
	},
	template: `<div>
	<pre class="log-box" style="margin: 0; padding: 0"><span :class="{red:log.level == 'error', orange:log.level == 'warning', green: log.level == 'success'}"><span v-if="!log.isToday">[{{log.createdTime}}]</span><strong v-if="log.isToday">[{{log.createdTime}}]</strong><keyword :v-word="keyword">[{{log.tag}}]{{log.description}}</keyword></span> &nbsp; <span v-if="log.count > 1" class="ui label tiny" :class="{red:log.level == 'error', orange:log.level == 'warning'}">共{{log.count}}条</span> <span v-if="log.server != null && log.server.id > 0"><a :href="'/servers/server?serverId=' + log.server.id" class="ui label tiny basic">{{log.server.name}}</a></span></pre>
</div>`
})

Vue.component("provinces-selector", {
	props: ["v-provinces"],
	data: function () {
		let provinces = this.vProvinces
		if (provinces == null) {
			provinces = []
		}
		let provinceIds = provinces.$map(function (k, v) {
			return v.id
		})
		return {
			provinces: provinces,
			provinceIds: provinceIds
		}
	},
	methods: {
		add: function () {
			let provinceStringIds = this.provinceIds.map(function (v) {
				return v.toString()
			})
			let that = this
			teaweb.popup("/ui/selectProvincesPopup?provinceIds=" + provinceStringIds.join(","), {
				width: "48em",
				height: "23em",
				callback: function (resp) {
					that.provinces = resp.data.provinces
					that.change()
				}
			})
		},
		remove: function (index) {
			this.provinces.$remove(index)
			this.change()
		},
		change: function () {
			this.provinceIds = this.provinces.$map(function (k, v) {
				return v.id
			})
		}
	},
	template: `<div>
	<input type="hidden" name="provinceIdsJSON" :value="JSON.stringify(provinceIds)"/>
	<div v-if="provinces.length > 0" style="margin-bottom: 0.5em">
		<div v-for="(province, index) in provinces" class="ui label tiny basic">{{province.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("csrf-token", {
	created: function () {
		this.refreshToken()
	},
	mounted: function () {
		let that = this
		this.$refs.token.form.addEventListener("submit", function () {
			that.refreshToken()
		})

		// 自动刷新
		setInterval(function () {
			that.refreshToken()
		}, 10 * 60 * 1000)
	},
	data: function () {
		return {
			token: ""
		}
	},
	methods: {
		refreshToken: function () {
			let that = this
			Tea.action("/csrf/token")
				.get()
				.success(function (resp) {
					that.token = resp.data.token
				})
		}
	},
	template: `<input type="hidden" name="csrfToken" :value="token" ref="token"/>`
})


Vue.component("labeled-input", {
	props: ["name", "size", "maxlength", "label", "value"],
	template: '<div class="ui input right labeled"> \
	<input type="text" :name="name" :size="size" :maxlength="maxlength" :value="value"/>\
	<span class="ui label">{{label}}</span>\
</div>'
});

let radioId = 0
Vue.component("radio", {
	props: ["name", "value", "v-value", "id"],
	data: function () {
		radioId++
		let elementId = this.id
		if (elementId == null) {
			elementId = "radio" + radioId
		}
		return {
			"elementId": elementId
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.vValue)
		}
	},
	template: `<div class="ui checkbox radio">
	<input type="radio" :name="name" :value="vValue" :id="elementId" @change="change" :checked="(vValue == value)"/>
	<label :for="elementId"><slot></slot></label>
</div>`
})

Vue.component("copy-to-clipboard", {
	props: ["v-target"],
	created: function () {
		if (typeof ClipboardJS == "undefined") {
			let jsFile = document.createElement("script")
			jsFile.setAttribute("src", "/js/clipboard.min.js")
			document.head.appendChild(jsFile)
		}
	},
	methods: {
		copy: function () {
			new ClipboardJS('[data-clipboard-target]');
			teaweb.successToast("已复制到剪切板")
		}
	},
	template: `<a href="" title="拷贝到剪切板" :data-clipboard-target="'#' + vTarget" @click.prevent="copy"><i class="ui icon copy small"></i></em></a>`
})

// 节点角色名称
Vue.component("node-role-name", {
	props: ["v-role"],
	data: function () {
		let roleName = ""
		switch (this.vRole) {
			case "node":
				roleName = "边缘节点"
				break
			case "monitor":
				roleName = "监控节点"
				break
			case "api":
				roleName = "API节点"
				break
			case "user":
				roleName = "用户平台"
				break
			case "admin":
				roleName = "管理平台"
				break
			case "database":
				roleName = "数据库节点"
				break
			case "dns":
				roleName = "DNS节点"
				break
			case "report":
				roleName = "区域监控终端"
				break
		}
		return {
			roleName: roleName
		}
	},
	template: `<span>{{roleName}}</span>`
})

let sourceCodeBoxIndex = 0

Vue.component("source-code-box", {
	props: ["name", "type", "id", "read-only", "width", "height", "focus"],
	mounted: function () {
		let readOnly = this.readOnly
		if (typeof readOnly != "boolean") {
			readOnly = true
		}
		let box = document.getElementById("source-code-box-" + this.index)
		let valueBox = document.getElementById(this.valueBoxId)
		let value = ""
		if (valueBox.textContent != null) {
			value = valueBox.textContent
		} else if (valueBox.innerText != null) {
			value = valueBox.innerText
		}

		this.createEditor(box, value, readOnly)
	},
	data: function () {
		let index = sourceCodeBoxIndex++

		let valueBoxId = 'source-code-box-value-' + sourceCodeBoxIndex
		if (this.id != null) {
			valueBoxId = this.id
		}

		return {
			index: index,
			valueBoxId: valueBoxId
		}
	},
	methods: {
		createEditor: function (box, value, readOnly) {
			let boxEditor = CodeMirror.fromTextArea(box, {
				theme: "idea",
				lineNumbers: true,
				value: "",
				readOnly: readOnly,
				showCursorWhenSelecting: true,
				height: "auto",
				//scrollbarStyle: null,
				viewportMargin: Infinity,
				lineWrapping: true,
				highlightFormatting: false,
				indentUnit: 4,
				indentWithTabs: true,
			})
			let that = this
			boxEditor.on("change", function () {
				that.change(boxEditor.getValue())
			})
			boxEditor.setValue(value)

			if (this.focus) {
				boxEditor.focus()
			}

			let width = this.width
			let height = this.height
			if (width != null && height != null) {
				width = parseInt(width)
				height = parseInt(height)
				if (!isNaN(width) && !isNaN(height)) {
					if (width <= 0) {
						width = box.parentNode.offsetWidth
					}
					boxEditor.setSize(width, height)
				}
			} else if (height != null) {
				height = parseInt(height)
				if (!isNaN(height)) {
					boxEditor.setSize("100%", height)
				}
			}

			let info = CodeMirror.findModeByMIME(this.type)
			if (info != null) {
				boxEditor.setOption("mode", info.mode)
				CodeMirror.modeURL = "/codemirror/mode/%N/%N.js"
				CodeMirror.autoLoadMode(boxEditor, info.mode)
			}
		},
		change: function (code) {
			this.$emit("change", code)
		}
	},
	template: `<div class="source-code-box">
	<div style="display: none" :id="valueBoxId"><slot></slot></div>
	<textarea :id="'source-code-box-' + index" :name="name"></textarea>
</div>`
})

Vue.component("size-capacity-box", {
	props: ["v-name", "v-value", "v-count", "v-unit", "size", "maxlength", "v-supported-units"],
	data: function () {
		let v = this.vValue
		if (v == null) {
			v = {
				count: this.vCount,
				unit: this.vUnit
			}
		}
		if (typeof (v["count"]) != "number") {
			v["count"] = -1
		}

		let vSize = this.size
		if (vSize == null) {
			vSize = 6
		}

		let vMaxlength = this.maxlength
		if (vMaxlength == null) {
			vMaxlength = 10
		}

		let supportedUnits = this.vSupportedUnits
		if (supportedUnits == null) {
			supportedUnits = []
		}

		return {
			capacity: v,
			countString: (v.count >= 0) ? v.count.toString() : "",
			vSize: vSize,
			vMaxlength: vMaxlength,
			supportedUnits: supportedUnits
		}
	},
	watch: {
		"countString": function (newValue) {
			let value = newValue.trim()
			if (value.length == 0) {
				this.capacity.count = -1
				this.change()
				return
			}
			let count = parseInt(value)
			if (!isNaN(count)) {
				this.capacity.count = count
			}
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.capacity)
		}
	},
	template: `<div class="ui fields inline">
	<input type="hidden" :name="vName" :value="JSON.stringify(capacity)"/>
	<div class="ui field">
		<input type="text" v-model="countString" :maxlength="vMaxlength" :size="vSize"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="capacity.unit" @change="change">
			<option value="byte" v-if="supportedUnits.length == 0 || supportedUnits.$contains('byte')">字节</option>
			<option value="kb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('kb')">KiB</option>
			<option value="mb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('mb')">MiB</option>
			<option value="gb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('gb')">GiB</option>
			<option value="tb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('tb')">TiB</option>
			<option value="pb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('pb')">PiB</option>
			<option value="eb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('eb')">EiB</option>
		</select>
	</div>
</div>`
})

/**
 * 二级菜单
 */
Vue.component("inner-menu", {
	template: `
		<div class="second-menu" style="width:80%;position: absolute;top:-8px;right:1em"> 
			<div class="ui menu text blue small">
				<slot></slot>
			</div> 
		</div>`
});

Vue.component("datepicker", {
	props: ["value", "v-name", "name", "v-value", "v-bottom-left", "placeholder"],
	mounted: function () {
		let that = this
		teaweb.datepicker(this.$refs.dayInput, function (v) {
			that.day = v
			that.change()
		}, !!this.vBottomLeft)
	},
	data: function () {
		let name = this.vName
		if (name == null) {
			name = this.name
		}
		if (name == null) {
			name = "day"
		}

		let day = this.vValue
		if (day == null) {
			day = this.value
			if (day == null) {
				day = ""
			}
		}

		let placeholder = "YYYY-MM-DD"
		if (this.placeholder != null) {
			placeholder = this.placeholder
		}

		return {
			realName: name,
			realPlaceholder: placeholder,
			day: day
		}
	},
	watch: {
		value: function (v) {
			this.day = v

			let picker = this.$refs.dayInput.picker
			if (picker != null) {
				if (v != null && /^\d+-\d+-\d+$/.test(v)) {
					picker.setDate(v)
				}
			}
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.day) // support v-model，事件触发需要在 change 之前
			this.$emit("change", this.day)
		}
	},
	template: `<div style="display: inline-block">
	<input type="text" :name="realName" v-model="day" :placeholder="realPlaceholder" style="width:8.6em" maxlength="10" @input="change" ref="dayInput" autocomplete="off"/>
</div>`
})

// 排序使用的箭头
Vue.component("sort-arrow", {
	props: ["name"],
	data: function () {
		let url = window.location.toString()
		let order = ""
		let iconTitle = ""
		let newArgs = []
		if (window.location.search != null && window.location.search.length > 0) {
			let queryString = window.location.search.substring(1)
			let pieces = queryString.split("&")
			let that = this
			pieces.forEach(function (v) {
				let eqIndex = v.indexOf("=")
				if (eqIndex > 0) {
					let argName = v.substring(0, eqIndex)
					let argValue = v.substring(eqIndex + 1)
					if (argName == that.name) {
						order = argValue
					} else if (argName != "page" && argValue != "asc" && argValue != "desc") {
						newArgs.push(v)
					}
				} else {
					newArgs.push(v)
				}
			})
		}
		if (order == "asc") {
			newArgs.push(this.name + "=desc")
			iconTitle = "当前正序排列"
		} else if (order == "desc") {
			newArgs.push(this.name + "=asc")
			iconTitle = "当前倒序排列"
		} else {
			newArgs.push(this.name + "=desc")
			iconTitle = "当前正序排列"
		}

		let qIndex = url.indexOf("?")
		if (qIndex > 0) {
			url = url.substring(0, qIndex) + "?" + newArgs.join("&")
		} else {
			url = url + "?" + newArgs.join("&")
		}

		return {
			order: order,
			url: url,
			iconTitle: iconTitle
		}
	},
	template: `<a :href="url" :title="iconTitle">&nbsp; <i class="ui icon long arrow small" :class="{down: order == 'asc', up: order == 'desc', 'down grey': order == '' || order == null}"></i></a>`
})

Vue.component("user-link", {
	props: ["v-user", "v-keyword"],
	data: function () {
		let user = this.vUser
		if (user == null) {
			user = {id: 0, "username": "", "fullname": ""}
		}
		return {
			user: user
		}
	},
	template: `<div style="display: inline-block">
	<span v-if="user.id > 0"><keyword :v-word="vKeyword">{{user.fullname}}</keyword><span class="small grey">（<keyword :v-word="vKeyword">{{user.username}}</keyword>）</span></span>
	<span v-else class="disabled">[已删除]</span>
</div>`
})

// 监控节点分组选择
Vue.component("report-node-groups-selector", {
	props: ["v-group-ids"],
	mounted: function () {
		let that = this
		Tea.action("/clusters/monitors/groups/options")
			.post()
			.success(function (resp) {
				that.groups = resp.data.groups.map(function (group) {
					group.isChecked = that.groupIds.$contains(group.id)
					return group
				})
				that.isLoaded = true
			})
	},
	data: function () {
		var groupIds = this.vGroupIds
		if (groupIds == null) {
			groupIds = []
		}

		return {
			groups: [],
			groupIds: groupIds,
			isLoaded: false,
			allGroups: groupIds.length == 0
		}
	},
	methods: {
		check: function (group) {
			group.isChecked = !group.isChecked
			this.groupIds = []
			let that = this
			this.groups.forEach(function (v) {
				if (v.isChecked) {
					that.groupIds.push(v.id)
				}
			})
			this.change()
		},
		change: function () {
			let that = this
			let groups = []
			this.groupIds.forEach(function (groupId) {
				let group = that.groups.$find(function (k, v) {
					return v.id == groupId
				})
				if (group == null) {
					return
				}
				groups.push({
					id: group.id,
					name: group.name
				})
			})
			this.$emit("change", groups)
		}
	},
	watch: {
		allGroups: function (b) {
			if (b) {
				this.groupIds = []
				this.groups.forEach(function (v) {
					v.isChecked = false
				})
			}

			this.change()
		}
	},
	template: `<div>
	<input type="hidden" name="reportNodeGroupIdsJSON" :value="JSON.stringify(groupIds)"/>
	<span class="disabled" v-if="isLoaded && groups.length == 0">还没有分组。</span>
	<div v-if="groups.length > 0">
		<div>
			<div class="ui checkbox">
				<input type="checkbox" v-model="allGroups" id="all-group"/>
				<label for="all-group">全部分组</label>
			</div>
			<div class="ui divider" v-if="!allGroups"></div>
		</div>
		<div v-show="!allGroups">
			<div v-for="group in groups" :key="group.id" style="float: left; width: 7.6em; margin-bottom: 0.5em">
				<div class="ui checkbox">
					<input type="checkbox" v-model="group.isChecked" value="1" :id="'report-node-group-' + group.id" @click.prevent="check(group)"/>
					<label :for="'report-node-group-' + group.id">{{group.name}}</label>
				</div>
			</div>
		</div>
	</div>
</div>`
})

Vue.component("finance-user-selector", {
	props: ["v-user-id"],
	data: function () {
		return {}
	},
	methods: {
		change: function (userId) {
			this.$emit("change", userId)
		}
	},
	template: `<div>
	<user-selector :v-user-id="vUserId" data-url="/finance/users/options" @change="change"></user-selector>
</div>`
})

Vue.component("node-cache-disk-dirs-box", {
	props: ["value", "name"],
	data: function () {
		let dirs = this.value
		if (dirs == null) {
			dirs = []
		}
		return {
			dirs: dirs,

			isEditing: false,
			isAdding: false,

			addingPath: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingPath.focus()
			}, 100)
		},
		confirm: function () {
			let addingPath = this.addingPath.trim()
			if (addingPath.length == 0) {
				let that = this
				teaweb.warn("请输入要添加的缓存目录", function () {
					that.$refs.addingPath.focus()
				})
				return
			}
			if (addingPath[0] != "/") {
				addingPath = "/" + addingPath
			}
			this.dirs.push({
				path: addingPath
			})
			this.cancel()
		},
		cancel: function () {
			this.addingPath = ""
			this.isAdding = false
			this.isEditing = false
		},
		remove: function (index) {
			let that = this
			teaweb.confirm("确定要删除此目录吗？", function () {
				that.dirs.$remove(index)
			})
		}
	},
	template: `<div>
	<input type="hidden" :name="name" :value="JSON.stringify(dirs)"/>
	<div style="margin-bottom: 0.3em">
		<span class="ui label small basic" v-for="(dir, index) in dirs">
			<i class="icon folder"></i>{{dir.path}}  &nbsp;  <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
	</div>
	
	<!-- 添加 -->
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" style="width: 30em" v-model="addingPath" @keyup.enter="confirm()" @keypress.enter.prevent="1" @keydown.esc="cancel()" ref="addingPath" placeholder="新的缓存目录，比如 /mnt/cache"/>
			</div>
			<div class="ui field">
				<button class="ui button small" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("node-ip-address-clusters-selector", {
	props: ["vClusters"],
	mounted: function () {
		this.checkClusters()
	},
	data: function () {
		let clusters = this.vClusters
		if (clusters == null) {
			clusters = []
		}
		return {
			clusters: clusters,
			hasCheckedCluster: false,
			clustersVisible: false
		}
	},
	methods: {
		checkClusters: function () {
			let that = this

			let b = false
			this.clusters.forEach(function (cluster) {
				if (cluster.isChecked) {
					b = true
				}
			})

			this.hasCheckedCluster = b

			return b
		},
		changeCluster: function (cluster) {
			cluster.isChecked = !cluster.isChecked
			this.checkClusters()
		},
		showClusters: function () {
			this.clustersVisible = !this.clustersVisible
		}
	},
	template: `<div>
  <span v-if="!hasCheckedCluster">默认用于所有集群 &nbsp; <a href="" @click.prevent="showClusters">修改 <i class="icon angle" :class="{down: !clustersVisible, up:clustersVisible}"></i></a></span>
	<div v-if="hasCheckedCluster">
		<span v-for="cluster in clusters" class="ui label basic small" v-if="cluster.isChecked">{{cluster.name}}</span> &nbsp; <a href="" @click.prevent="showClusters">修改 <i class="icon angle" :class="{down: !clustersVisible, up:clustersVisible}"></i></a>
		<p class="comment">当前IP仅在所选集群中有效。</p>
	</div>
	<div v-show="clustersVisible">
		<div class="ui divider"></div>
		<checkbox v-for="cluster in clusters" :v-value="cluster.id" :value="cluster.isChecked ? cluster.id : 0" style="margin-right: 1em" @input="changeCluster(cluster)" name="clusterIds">
			{{cluster.name}}
		</checkbox>
	</div>
</div>`
})

// 节点登录推荐端口
Vue.component("node-login-suggest-ports", {
	data: function () {
		return {
			ports: [],
			availablePorts: [],
			autoSelected: false,
			isLoading: false
		}
	},
	methods: {
		reload: function (host) {
			let that = this
			this.autoSelected = false
			this.isLoading = true
			Tea.action("/clusters/cluster/suggestLoginPorts")
				.params({
					host: host
				})
				.success(function (resp) {
					if (resp.data.availablePorts != null) {
						that.availablePorts = resp.data.availablePorts
						if (that.availablePorts.length > 0) {
							that.autoSelectPort(that.availablePorts[0])
							that.autoSelected = true
						}
					}
					if (resp.data.ports != null) {
						that.ports = resp.data.ports
						if (that.ports.length > 0 && !that.autoSelected) {
							that.autoSelectPort(that.ports[0])
							that.autoSelected = true
						}
					}
				})
				.done(function () {
					that.isLoading = false
				})
				.post()
		},
		selectPort: function (port) {
			this.$emit("select", port)
		},
		autoSelectPort: function (port) {
			this.$emit("auto-select", port)
		}
	},
	template: `<span>
	<span v-if="isLoading">正在检查端口...</span>
	<span v-if="availablePorts.length > 0">
		可能端口：<a href="" v-for="port in availablePorts" @click.prevent="selectPort(port)" class="ui label tiny basic blue" style="border: 1px #2185d0 dashed; font-weight: normal">{{port}}</a>
		&nbsp; &nbsp;
	</span>
	<span v-if="ports.length > 0">
		常用端口：<a href="" v-for="port in ports" @click.prevent="selectPort(port)" class="ui label tiny basic blue" style="border: 1px #2185d0 dashed;  font-weight: normal">{{port}}</a>
	</span>
	<span v-if="ports.length == 0">常用端口有22等。</span>
	<span v-if="ports.length > 0" class="grey small">（可以点击要使用的端口）</span>
</span>`
})

Vue.component("node-group-selector", {
	props: ["v-cluster-id", "v-group"],
	data: function () {
		return {
			selectedGroup: this.vGroup
		}
	},
	methods: {
		selectGroup: function () {
			let that = this
			teaweb.popup("/clusters/cluster/groups/selectPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedGroup = resp.data.group
				}
			})
		},
		addGroup: function () {
			let that = this
			teaweb.popup("/clusters/cluster/groups/createPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedGroup = resp.data.group
				}
			})
		},
		removeGroup: function () {
			this.selectedGroup = null
		}
	},
	template: `<div>
	<div class="ui label small basic" v-if="selectedGroup != null">
		<input type="hidden" name="groupId" :value="selectedGroup.id"/>
		{{selectedGroup.name}} &nbsp;<a href="" title="删除" @click.prevent="removeGroup()"><i class="icon remove"></i></a>
	</div>
	<div v-if="selectedGroup == null">
		<a href="" @click.prevent="selectGroup()">[选择分组]</a> &nbsp; <a href="" @click.prevent="addGroup()">[添加分组]</a>
	</div>
</div>`
})

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

Vue.component("node-schedule-conds-box", {
	props: ["value", "v-params", "v-operators"],
	mounted: function () {
		this.formatConds(this.condsConfig.conds)
		this.$forceUpdate()
	},
	data: function () {
		let condsConfig = this.value
		if (condsConfig == null) {
			condsConfig = {
				conds: [],
				connector: "and"
			}
		}
		if (condsConfig.conds == null) {
			condsConfig.conds = []
		}

		let paramMap = {}
		this.vParams.forEach(function (param) {
			paramMap[param.code] = param
		})

		let operatorMap = {}
		this.vOperators.forEach(function (operator) {
			operatorMap[operator.code] = operator.name
		})

		return {
			condsConfig: condsConfig,
			params: this.vParams,
			paramMap: paramMap,
			operatorMap: operatorMap,
			operator: "",

			isAdding: false,

			paramCode: "",
			param: null,

			valueBandwidth: {
				count: 100,
				unit: "mb"
			},
			valueTraffic: {
				count: 1,
				unit: "gb"
			},
			valueCPU: 80,
			valueMemory: 90,
			valueLoad: 20,
			valueRate: 0
		}
	},
	watch: {
		paramCode: function (code) {
			if (code.length == 0) {
				this.param = null
			} else {
				this.param = this.params.$find(function (k, v) {
					return v.code == code
				})
			}
			this.$emit("changeparam", this.param)
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
		},
		confirm: function () {
			if (this.param == null) {
				teaweb.warn("请选择参数")
				return
			}
			if (this.param.operators != null && this.param.operators.length > 0 && this.operator.length == 0) {
				teaweb.warn("请选择操作符")
				return
			}
			if (this.param.operators == null || this.param.operators.length == 0) {
				this.operator = ""
			}

			let value = null
			switch (this.param.valueType) {
				case "bandwidth": {
					if (this.valueBandwidth.unit.length == 0) {
						teaweb.warn("请选择带宽单位")
						return
					}
					let count = parseInt(this.valueBandwidth.count.toString())
					if (isNaN(count)) {
						count = 0
					}
					if (count < 0) {
						count = 0
					}
					value = {
						count: count,
						unit: this.valueBandwidth.unit
					}
				}
					break
				case "traffic": {
					if (this.valueTraffic.unit.length == 0) {
						teaweb.warn("请选择带宽单位")
						return
					}
					let count = parseInt(this.valueTraffic.count.toString())
					if (isNaN(count)) {
						count = 0
					}
					if (count < 0) {
						count = 0
					}
					value = {
						count: count,
						unit: this.valueTraffic.unit
					}
				}
					break
				case "cpu":
					let cpu = parseInt(this.valueCPU.toString())
					if (isNaN(cpu)) {
						cpu = 0
					}
					if (cpu < 0) {
						cpu = 0
					}
					if (cpu > 100) {
						cpu = 100
					}
					value = cpu
					break
				case "memory":
					let memory = parseInt(this.valueMemory.toString())
					if (isNaN(memory)) {
						memory = 0
					}
					if (memory < 0) {
						memory = 0
					}
					if (memory > 100) {
						memory = 100
					}
					value = memory
					break
				case "load":
					let load = parseInt(this.valueLoad.toString())
					if (isNaN(load)) {
						load = 0
					}
					if (load < 0) {
						load = 0
					}
					value = load
					break
				case "rate":
					let rate = parseInt(this.valueRate.toString())
					if (isNaN(rate)) {
						rate = 0
					}
					if (rate < 0) {
						rate = 0
					}
					value = rate
					break
			}

			this.condsConfig.conds.push({
				param: this.param.code,
				operator: this.operator,
				value: value
			})
			this.formatConds(this.condsConfig.conds)

			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.paramCode = ""
			this.param = null
		},
		remove: function (index) {
			this.condsConfig.conds.$remove(index)
		},
		formatConds: function (conds) {
			let that = this
			conds.forEach(function (cond) {
				switch (that.paramMap[cond.param].valueType) {
					case "bandwidth":
						cond.valueFormat = cond.value.count + cond.value.unit[0].toUpperCase() + cond.value.unit.substring(1) + "ps"
						return
					case "traffic":
						cond.valueFormat = cond.value.count + cond.value.unit.toUpperCase()
						return
					case "cpu":
						cond.valueFormat = cond.value + "%"
						return
					case "memory":
						cond.valueFormat = cond.value + "%"
						return
					case "load":
						cond.valueFormat = cond.value
						return
					case "rate":
						cond.valueFormat = cond.value + "/秒"
						return
				}
			})
		}
	},
	template: `<div>
	<input type="hidden" name="condsJSON" :value="JSON.stringify(this.condsConfig)"/>
	
	<!-- 已有条件 -->
	<div v-if="condsConfig.conds.length > 0" style="margin-bottom: 1em">
		<span v-for="(cond, index) in condsConfig.conds">
			<span class="ui label basic small">
				<span>{{paramMap[cond.param].name}} 
					<span v-if="paramMap[cond.param].operators != null && paramMap[cond.param].operators.length > 0"><span class="grey">{{operatorMap[cond.operator]}}</span> {{cond.valueFormat}}</span> 
					&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
				</span>
			</span>
			<span v-if="index < condsConfig.conds.length - 1"> &nbsp;<span v-if="condsConfig.connector == 'and'">且</span><span v-else>或</span>&nbsp; </span>
		</span>
	</div>
	
	<div v-if="isAdding">
		<table class="ui table">
			<tbody>
				<tr>
					<td class="title">参数</td>
					<td>
						<select class="ui dropdown auto-width" v-model="paramCode">
							<option value="">[选择参数]</option>
							<option v-for="paramOption in params" :value="paramOption.code">{{paramOption.name}}</option>
						</select>
						<p class="comment" v-if="param != null">{{param.description}}</p>
					</td>
				</tr>
				<tr v-if="param != null && param.operators != null && param.operators.length > 0">
					<td>操作符</td>
					<td>
						<select class="ui dropdown auto-width" v-if="param != null" v-model="operator">
							<option value="">[选择操作符]</option>
							<option v-for="operator in param.operators" :value="operator">{{operatorMap[operator]}}</option>
						</select>
					</td>
				</tr>
				<tr v-if="param != null && param.operators != null && param.operators.length > 0">
					<td>{{param.valueName}}</td>
					<td>
						<!-- 带宽 -->
						<div v-if="param.valueType == 'bandwidth'">
							<div class="ui fields inline">
								<div class="ui field">
									<input type="text" maxlength="10" size="6" v-model="valueBandwidth.count" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
								</div>
								<div class="ui field">
									<select class="ui dropdown auto-width" v-model="valueBandwidth.unit">
										<option value="gb">Gbps</option>
										<option value="mb">Mbps</option>
									</select>
								</div>
							</div>
						</div>
						
						<!-- 流量 -->
						<div v-if="param.valueType == 'traffic'">
							<div class="ui fields inline">
								<div class="ui field">
									<input type="text" maxlength="10" size="6" v-model="valueTraffic.count" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
								</div>
								<div class="ui field">
									<select class="ui dropdown auto-width" v-model="valueTraffic.unit">
										<option value="mb">MiB</option>
										<option value="gb">GiB</option>
										<option value="tb">TiB</option>
										<option value="pb">PiB</option>
										<option value="eb">EiB</option>
									</select>
								</div>
							</div>
						</div>
						
						<!-- cpu -->
						<div v-if="param.valueType == 'cpu'">
							<div class="ui input right labeled">
								<input type="text" v-model="valueCPU" maxlength="3" size="3" style="width: 4em" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
								<span class="ui label">%</span>
							</div>
						</div>
						
						<!-- memory -->
						<div v-if="param.valueType == 'memory'">
							<div class="ui input right labeled">
								<input type="text" v-model="valueMemory" maxlength="3" size="3" style="width: 4em" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
								<span class="ui label">%</span>
							</div>
						</div>
						
						<!-- load -->
						<div v-if="param.valueType == 'load'">
							<input type="text" v-model="valueLoad" maxlength="3" size="3" style="width: 4em" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
						</div>
						
						<!-- rate -->
						<div v-if="param.valueType == 'rate'">
							<div class="ui input right labeled">
								<input type="text" v-model="valueRate" maxlength="8" size="8" style="width: 8em" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
								<span class="ui label">/秒</span>
							</div>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
		<button class="ui button small" type="button" @click.prevent="confirm">确定</button> &nbsp; <a href="" @click.prevent="cancel">取消</a>
	</div>
	
	<div v-if="!isAdding">
		<button class="ui button small" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("node-schedule-action-box", {
	props: ["value", "v-actions"],
	data: function () {
		let actionConfig = this.value
		if (actionConfig == null) {
			actionConfig = {
				code: "",
				params: {}
			}
		}

		return {
			actions: this.vActions,
			currentAction: null,
			actionConfig: actionConfig
		}
	},
	watch: {
		"actionConfig.code": function (actionCode) {
			if (actionCode.length == 0) {
				this.currentAction = null
			} else {
				this.currentAction = this.actions.$find(function (k, v) {
					return v.code == actionCode
				})
			}
			this.actionConfig.params = {}
		}
	},
	template: `<div>
	<input type="hidden" name="actionJSON" :value="JSON.stringify(actionConfig)"/>
	<div>
		<div>
			<select class="ui dropdown auto-width" v-model="actionConfig.code">
				<option value="">[选择动作]</option>
				<option v-for="action in actions" :value="action.code">{{action.name}}</option>
			</select>
		</div>
		<p class="comment" v-if="currentAction != null">{{currentAction.description}}</p>
		
		<div v-if="actionConfig.code == 'webHook'">
			<input type="text" placeholder="https://..." v-model="actionConfig.params.url"/>
			<p class="comment">接收通知的URL。</p>
		</div>
	</div>
</div>`
})

// 节点IP阈值
Vue.component("node-ip-address-thresholds-view", {
	props: ["v-thresholds"],
	data: function () {
		let thresholds = this.vThresholds
		if (thresholds == null) {
			thresholds = []
		} else {
			thresholds.forEach(function (v) {
				if (v.items == null) {
					v.items = []
				}
				if (v.actions == null) {
					v.actions = []
				}
			})
		}

		return {
			thresholds: thresholds,
			allItems: window.IP_ADDR_THRESHOLD_ITEMS,
			allOperators: [
				{
					"name": "小于等于",
					"code": "lte"
				},
				{
					"name": "大于",
					"code": "gt"
				},
				{
					"name": "不等于",
					"code": "neq"
				},
				{
					"name": "小于",
					"code": "lt"
				},
				{
					"name": "大于等于",
					"code": "gte"
				}
			],
			allActions: window.IP_ADDR_THRESHOLD_ACTIONS
		}
	},
	methods: {
		itemName: function (item) {
			let result = ""
			this.allItems.forEach(function (v) {
				if (v.code == item) {
					result = v.name
				}
			})
			return result
		},
		itemUnitName: function (itemCode) {
			let result = ""
			this.allItems.forEach(function (v) {
				if (v.code == itemCode) {
					result = v.unit
				}
			})
			return result
		},
		itemDurationUnitName: function (unit) {
			switch (unit) {
				case "minute":
					return "分钟"
				case "second":
					return "秒"
				case "hour":
					return "小时"
				case "day":
					return "天"
			}
			return unit
		},
		itemOperatorName: function (operator) {
			let result = ""
			this.allOperators.forEach(function (v) {
				if (v.code == operator) {
					result = v.name
				}
			})
			return result
		},
		actionName: function (actionCode) {
			let result = ""
			this.allActions.forEach(function (v) {
				if (v.code == actionCode) {
					result = v.name
				}
			})
			return result
		}
	},
	template: `<div>
	<!-- 已有条件 -->
	<div v-if="thresholds.length > 0">
		<div class="ui label basic small" v-for="(threshold, index) in thresholds" style="margin-bottom: 0.8em">
			<span v-for="(item, itemIndex) in threshold.items">
				<span>
					<span v-if="item.item != 'nodeHealthCheck'">
						[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
					</span>	 
					{{itemName(item.item)}}
					
					<span v-if="item.item == 'nodeHealthCheck'">
						<!-- 健康检查 -->
						<span v-if="item.value == 1">成功</span>
						<span v-if="item.value == 0">失败</span>
					</span>
					<span v-else>
						<!-- 连通性 -->
						<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
						
						 <span class="grey">[{{itemOperatorName(item.operator)}}]</span> {{item.value}}{{itemUnitName(item.item)}} &nbsp;
					 </span>
				 </span>
				 <span v-if="itemIndex != threshold.items.length - 1" style="font-style: italic">AND &nbsp;</span></span>
				-&gt;
				<span v-for="(action, actionIndex) in threshold.actions">{{actionName(action.action)}}
				<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
				<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
				 &nbsp;					 
				 <span v-if="actionIndex != threshold.actions.length - 1" style="font-style: italic">AND &nbsp;</span>
			 </span>
		</div>
	</div>
</div>`
})

// 节点IP阈值
Vue.component("node-ip-address-thresholds-box", {
	props: ["v-thresholds"],
	data: function () {
		let thresholds = this.vThresholds
		if (thresholds == null) {
			thresholds = []
		} else {
			thresholds.forEach(function (v) {
				if (v.items == null) {
					v.items = []
				}
				if (v.actions == null) {
					v.actions = []
				}
			})
		}

		return {
			editingIndex: -1,
			thresholds: thresholds,
			addingThreshold: {
				items: [],
				actions: []
			},
			isAdding: false,
			isAddingItem: false,
			isAddingAction: false,

			itemCode: "nodeAvgRequests",
			itemReportGroups: [],
			itemOperator: "lte",
			itemValue: "",
			itemDuration: "5",
			allItems: window.IP_ADDR_THRESHOLD_ITEMS,
			allOperators: [
				{
					"name": "小于等于",
					"code": "lte"
				},
				{
					"name": "大于",
					"code": "gt"
				},
				{
					"name": "不等于",
					"code": "neq"
				},
				{
					"name": "小于",
					"code": "lt"
				},
				{
					"name": "大于等于",
					"code": "gte"
				}
			],
			allActions: window.IP_ADDR_THRESHOLD_ACTIONS,

			actionCode: "up",
			actionBackupIPs: "",
			actionWebHookURL: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = !this.isAdding
		},
		cancel: function () {
			this.isAdding = false
			this.editingIndex = -1
			this.addingThreshold = {
				items: [],
				actions: []
			}
		},
		confirm: function () {
			if (this.addingThreshold.items.length == 0) {
				teaweb.warn("需要至少添加一个阈值")
				return
			}
			if (this.addingThreshold.actions.length == 0) {
				teaweb.warn("需要至少添加一个动作")
				return
			}

			if (this.editingIndex >= 0) {
				this.thresholds[this.editingIndex].items = this.addingThreshold.items
				this.thresholds[this.editingIndex].actions = this.addingThreshold.actions
			} else {
				this.thresholds.push({
					items: this.addingThreshold.items,
					actions: this.addingThreshold.actions
				})
			}

			// 还原
			this.cancel()
		},
		remove: function (index) {
			this.cancel()
			this.thresholds.$remove(index)
		},
		update: function (index) {
			this.editingIndex = index
			this.addingThreshold = {
				items: this.thresholds[index].items.$copy(),
				actions: this.thresholds[index].actions.$copy()
			}
			this.isAdding = true
		},

		addItem: function () {
			this.isAddingItem = !this.isAddingItem
			let that = this
			setTimeout(function () {
				that.$refs.itemValue.focus()
			}, 100)
		},
		cancelItem: function () {
			this.isAddingItem = false

			this.itemCode = "nodeAvgRequests"
			this.itmeOperator = "lte"
			this.itemValue = ""
			this.itemDuration = "5"
			this.itemReportGroups = []
		},
		confirmItem: function () {
			// 特殊阈值快速添加
			if (["nodeHealthCheck"].$contains(this.itemCode)) {
				if (this.itemValue.toString().length == 0) {
					teaweb.warn("请选择检查结果")
					return
				}

				let value = parseInt(this.itemValue)
				if (isNaN(value)) {
					value = 0
				} else if (value < 0) {
					value = 0
				} else if (value > 1) {
					value = 1
				}

				// 添加
				this.addingThreshold.items.push({
					item: this.itemCode,
					operator: this.itemOperator,
					value: value,
					duration: 0,
					durationUnit: "minute",
					options: {}
				})
				this.cancelItem()
				return
			}

			if (this.itemDuration.length == 0) {
				let that = this
				teaweb.warn("请输入统计周期", function () {
					that.$refs.itemDuration.focus()
				})
				return
			}
			let itemDuration = parseInt(this.itemDuration)
			if (isNaN(itemDuration) || itemDuration <= 0) {
				teaweb.warn("请输入正确的统计周期", function () {
					that.$refs.itemDuration.focus()
				})
				return
			}

			if (this.itemValue.length == 0) {
				let that = this
				teaweb.warn("请输入对比值", function () {
					that.$refs.itemValue.focus()
				})
				return
			}
			let itemValue = parseFloat(this.itemValue)
			if (isNaN(itemValue)) {
				teaweb.warn("请输入正确的对比值", function () {
					that.$refs.itemValue.focus()
				})
				return
			}


			let options = {}

			switch (this.itemCode) {
				case "connectivity": // 连通性校验
					if (itemValue > 100) {
						let that = this
						teaweb.warn("连通性对比值不能超过100", function () {
							that.$refs.itemValue.focus()
						})
						return
					}

					options["groups"] = this.itemReportGroups
					break
			}

			// 添加
			this.addingThreshold.items.push({
				item: this.itemCode,
				operator: this.itemOperator,
				value: itemValue,
				duration: itemDuration,
				durationUnit: "minute",
				options: options
			})

			// 还原
			this.cancelItem()
		},
		removeItem: function (index) {
			this.cancelItem()
			this.addingThreshold.items.$remove(index)
		},
		changeReportGroups: function (groups) {
			this.itemReportGroups = groups
		},
		itemName: function (item) {
			let result = ""
			this.allItems.forEach(function (v) {
				if (v.code == item) {
					result = v.name
				}
			})
			return result
		},
		itemUnitName: function (itemCode) {
			let result = ""
			this.allItems.forEach(function (v) {
				if (v.code == itemCode) {
					result = v.unit
				}
			})
			return result
		},
		itemDurationUnitName: function (unit) {
			switch (unit) {
				case "minute":
					return "分钟"
				case "second":
					return "秒"
				case "hour":
					return "小时"
				case "day":
					return "天"
			}
			return unit
		},
		itemOperatorName: function (operator) {
			let result = ""
			this.allOperators.forEach(function (v) {
				if (v.code == operator) {
					result = v.name
				}
			})
			return result
		},

		addAction: function () {
			this.isAddingAction = !this.isAddingAction
		},
		cancelAction: function () {
			this.isAddingAction = false
			this.actionCode = "up"
			this.actionBackupIPs = ""
			this.actionWebHookURL = ""
		},
		confirmAction: function () {
			this.doConfirmAction(false)
		},
		doConfirmAction: function (validated, options) {
			// 是否已存在
			let exists = false
			let that = this
			this.addingThreshold.actions.forEach(function (v) {
				if (v.action == that.actionCode) {
					exists = true
				}
			})
			if (exists) {
				teaweb.warn("此动作已经添加过了，无需重复添加")
				return
			}

			if (options == null) {
				options = {}
			}

			switch (this.actionCode) {
				case "switch":
					if (!validated) {
						Tea.action("/ui/validateIPs")
							.params({
								"ips": this.actionBackupIPs
							})
							.success(function (resp) {
								if (resp.data.ips.length == 0) {
									teaweb.warn("请输入备用IP", function () {
										that.$refs.actionBackupIPs.focus()
									})
									return
								}
								options["ips"] = resp.data.ips
								that.doConfirmAction(true, options)
							})
							.fail(function (resp) {
								teaweb.warn("输入的IP '" + resp.data.failIP + "' 格式不正确，请改正后提交", function () {
									that.$refs.actionBackupIPs.focus()
								})
							})
							.post()
						return
					}
					break
				case "webHook":
					if (this.actionWebHookURL.length == 0) {
						teaweb.warn("请输入WebHook URL", function () {
							that.$refs.webHookURL.focus()
						})
						return
					}
					if (!this.actionWebHookURL.match(/^(http|https):\/\//i)) {
						teaweb.warn("URL开头必须是http://或者https://", function () {
							that.$refs.webHookURL.focus()
						})
						return
					}
					options["url"] = this.actionWebHookURL
			}

			this.addingThreshold.actions.push({
				action: this.actionCode,
				options: options
			})

			// 还原
			this.cancelAction()
		},
		removeAction: function (index) {
			this.cancelAction()
			this.addingThreshold.actions.$remove(index)
		},
		actionName: function (actionCode) {
			let result = ""
			this.allActions.forEach(function (v) {
				if (v.code == actionCode) {
					result = v.name
				}
			})
			return result
		}
	},
	template: `<div>
	<input type="hidden" name="thresholdsJSON" :value="JSON.stringify(thresholds)"/>
		
	<!-- 已有条件 -->
	<div v-if="thresholds.length > 0">
		<div class="ui label basic small" v-for="(threshold, index) in thresholds">
			<span v-for="(item, itemIndex) in threshold.items">
				<span v-if="item.item != 'nodeHealthCheck'">
					[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
				</span> 
				{{itemName(item.item)}}
				
				<span v-if="item.item == 'nodeHealthCheck'">
					<!-- 健康检查 -->
					<span v-if="item.value == 1">成功</span>
					<span v-if="item.value == 0">失败</span>
				</span>
				<span v-else>
					<!-- 连通性 -->
					<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
				
					<span  class="grey">[{{itemOperatorName(item.operator)}}]</span> &nbsp;{{item.value}}{{itemUnitName(item.item)}} 
			 	</span>
			 	&nbsp;<span v-if="itemIndex != threshold.items.length - 1" style="font-style: italic">AND &nbsp;</span>
			</span>
			-&gt;
			<span v-for="(action, actionIndex) in threshold.actions">{{actionName(action.action)}}
			<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
			<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
			 &nbsp;<span v-if="actionIndex != threshold.actions.length - 1" style="font-style: italic">AND &nbsp;</span></span>
			&nbsp;
			<a href="" title="修改" @click.prevent="update(index)"><i class="icon pencil small"></i></a> 
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon small remove"></i></a>
		</div>
	</div>
	
	<!-- 新阈值 -->
	<div v-if="isAdding" style="margin-top: 0.5em">
		<table class="ui table celled">
			<thead>
				<tr>
					<td style="width: 50%; background: #f9fafb; border-bottom: 1px solid rgba(34,36,38,.1)">阈值</td>
					<th>动作</th>
				</tr>
			</thead>
			<tr>
				<td style="background: white">
					<!-- 已经添加的项目 -->
					<div>
						<div v-for="(item, index) in addingThreshold.items" class="ui label basic small" style="margin-bottom: 0.5em;">
							<span v-if="item.item != 'nodeHealthCheck'">
								[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
							</span> 
							{{itemName(item.item)}}
							
							<span v-if="item.item == 'nodeHealthCheck'">
								<!-- 健康检查 -->
								<span v-if="item.value == 1">成功</span>
								<span v-if="item.value == 0">失败</span>
							</span>
							<span v-else>
								<!-- 连通性 -->
								<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
								 <span class="grey">[{{itemOperatorName(item.operator)}}]</span> {{item.value}}{{itemUnitName(item.item)}}
							 </span> 
							 &nbsp;
							<a href="" title="删除" @click.prevent="removeItem(index)"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<!-- 正在添加的项目 -->
					<div v-if="isAddingItem" style="margin-top: 0.8em">
						<table class="ui table">
							<tr>
								<td style="width: 6em">统计项目</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemCode">
									<option v-for="item in allItems" :value="item.code">{{item.name}}</option>
									</select>
									<p class="comment" style="font-weight: normal" v-for="item in allItems" v-if="item.code == itemCode">{{item.description}}</p>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>统计周期</td>
								<td>
									<div class="ui input right labeled">
										<input type="text" v-model="itemDuration" style="width: 4em" maxlength="4" ref="itemDuration" @keyup.enter="confirmItem()" @keypress.enter.prevent="1"/>
										<span class="ui label">分钟</span>
									</div>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>操作符</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemOperator">
										<option v-for="operator in allOperators" :value="operator.code">{{operator.name}}</option>
									</select>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>对比值</td>
								<td>
									<div class="ui input right labeled">
										<input type="text" maxlength="20" style="width: 5em" v-model="itemValue" ref="itemValue" @keyup.enter="confirmItem()" @keypress.enter.prevent="1"/>
										<span class="ui label" v-for="item in allItems" v-if="item.code == itemCode">{{item.unit}}</span>
									</div>
								</td>
							</tr>
							<tr v-show="itemCode == 'nodeHealthCheck'">
								<td>检查结果</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemValue">
										<option value="1">成功</option>
										<option value="0">失败</option>
									</select>
									<p class="comment" style="font-weight: normal">只有状态发生改变的时候才会触发。</p>
								</td>
							</tr>
							
							<!-- 连通性 -->
							<tr v-if="itemCode == 'connectivity'">
								<td>终端分组</td>
								<td style="font-weight: normal">
									<div style="zoom: 0.8"><report-node-groups-selector @change="changeReportGroups"></report-node-groups-selector></div>
								</td>
							</tr>
						</table>
						<div style="margin-top: 0.8em">
							<button class="ui button tiny" type="button" @click.prevent="confirmItem">确定</button>							 &nbsp;
							<a href="" title="取消" @click.prevent="cancelItem"><i class="icon remove small"></i></a>
						</div>
					</div>
					<div style="margin-top: 0.8em" v-if="!isAddingItem">
						<button class="ui button tiny" type="button" @click.prevent="addItem">+</button>
					</div>
				</td>
				<td style="background: white">
					<!-- 已经添加的动作 -->
					<div>
						<div v-for="(action, index) in addingThreshold.actions" class="ui label basic small" style="margin-bottom: 0.5em">
							{{actionName(action.action)}} &nbsp;
							<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
							<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
							<a href="" title="删除" @click.prevent="removeAction(index)"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<!-- 正在添加的动作 -->
					<div v-if="isAddingAction" style="margin-top: 0.8em">
						<table class="ui table">
							<tr>
								<td style="width: 6em">动作类型</td>
								<td>
									<select class="ui dropdown auto-width" v-model="actionCode">
										<option v-for="action in allActions" :value="action.code">{{action.name}}</option>
									</select>
									<p class="comment" v-for="action in allActions" v-if="action.code == actionCode">{{action.description}}</p>
								</td>
							</tr>
							
							<!-- 切换 -->
							<tr v-if="actionCode == 'switch'">
								<td>备用IP *</td>
								<td>
									<textarea rows="2" v-model="actionBackupIPs" ref="actionBackupIPs"></textarea>
									<p class="comment">每行一个备用IP。</p>
								</td>
							</tr>
							
							<!-- WebHook -->
							<tr v-if="actionCode == 'webHook'">
								<td>URL *</td>
								<td>
									<input type="text" maxlength="1000" placeholder="https://..." v-model="actionWebHookURL" ref="webHookURL" @keyup.enter="confirmAction()" @keypress.enter.prevent="1"/>
									<p class="comment">完整的URL，比如<code-label>https://example.com/webhook/api</code-label>，系统会在触发阈值的时候通过GET调用此URL。</p>
								</td>
							</tr>
						</table>
						<div style="margin-top: 0.8em">
							<button class="ui button tiny" type="button" @click.prevent="confirmAction">确定</button>	 &nbsp;
							<a href="" title="取消" @click.prevent="cancelAction"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<div style="margin-top: 0.8em" v-if="!isAddingAction">
						<button class="ui button tiny" type="button" @click.prevent="addAction">+</button>
					</div>	
				</td>
			</tr>
		</table>
		
		<!-- 添加阈值 -->
		<div>
			<button class="ui button tiny" :class="{disabled: (isAddingItem || isAddingAction)}" type="button" @click.prevent="confirm">确定</button> &nbsp;
			<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
		</div>
	</div>
	
	<div v-if="!isAdding" style="margin-top: 0.5em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("node-region-selector", {
	props: ["v-region"],
	data: function () {
		return {
			selectedRegion: this.vRegion
		}
	},
	methods: {
		selectRegion: function () {
			let that = this
			teaweb.popup("/clusters/regions/selectPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedRegion = resp.data.region
				}
			})
		},
		addRegion: function () {
			let that = this
			teaweb.popup("/clusters/regions/createPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedRegion = resp.data.region
				}
			})
		},
		removeRegion: function () {
			this.selectedRegion = null
		}
	},
	template: `<div>
	<div class="ui label small basic" v-if="selectedRegion != null">
		<input type="hidden" name="regionId" :value="selectedRegion.id"/>
		{{selectedRegion.name}} &nbsp;<a href="" title="删除" @click.prevent="removeRegion()"><i class="icon remove"></i></a>
	</div>
	<div v-if="selectedRegion == null">
		<a href="" @click.prevent="selectRegion()">[选择区域]</a> &nbsp; <a href="" @click.prevent="addRegion()">[添加区域]</a>
	</div>
</div>`
})

Vue.component("node-combo-box", {
	props: ["v-cluster-id", "v-node-id"],
	data: function () {
		let that = this
		Tea.action("/clusters/nodeOptions")
			.params({
				clusterId: this.vClusterId
			})
			.post()
			.success(function (resp) {
				that.nodes = resp.data.nodes
			})
		return {
			nodes: []
		}
	},
	template: `<div v-if="nodes.length > 0">
	<combo-box title="节点" placeholder="节点名称" :v-items="nodes" name="nodeId" :v-value="vNodeId"></combo-box>
</div>`
})

// 节点级别选择器
Vue.component("node-level-selector", {
	props: ["v-node-level"],
	data: function () {
		let levelCode = this.vNodeLevel
		if (levelCode == null || levelCode < 1) {
			levelCode = 1
		}
		return {
			levels: [
				{
					name: "边缘节点",
					code: 1,
					description: "普通的边缘节点。"
				},
				{
					name: "L2节点",
					code: 2,
					description: "特殊的边缘节点，同时负责同组上一级节点的回源。"
				}
			],
			levelCode: levelCode
		}
	},
	watch: {
		levelCode: function (code) {
			this.$emit("change", code)
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="level" v-model="levelCode">
	<option v-for="level in levels" :value="level.code">{{level.name}}</option>
</select>
<p class="comment" v-if="typeof(levels[levelCode - 1]) != null"><plus-label
></plus-label>{{levels[levelCode - 1].description}}</p>
</div>`
})

Vue.component("node-schedule-conds-viewer", {
	props: ["value", "v-params", "v-operators"],
	mounted: function () {
		this.formatConds(this.condsConfig.conds)
		this.$forceUpdate()
	},
	data: function () {
		let paramMap = {}
		this.vParams.forEach(function (param) {
			paramMap[param.code] = param
		})

		let operatorMap = {}
		this.vOperators.forEach(function (operator) {
			operatorMap[operator.code] = operator.name
		})

		return {
			condsConfig: this.value,
			paramMap: paramMap,
			operatorMap: operatorMap
		}
	},
	methods: {
		formatConds: function (conds) {
			let that = this
			conds.forEach(function (cond) {
				switch (that.paramMap[cond.param].valueType) {
					case "bandwidth":
						cond.valueFormat = cond.value.count + cond.value.unit[0].toUpperCase() + cond.value.unit.substring(1) + "ps"
						return
					case "traffic":
						cond.valueFormat = cond.value.count + cond.value.unit.toUpperCase()
						return
					case "cpu":
						cond.valueFormat = cond.value + "%"
						return
					case "memory":
						cond.valueFormat = cond.value + "%"
						return
					case "load":
						cond.valueFormat = cond.value
						return
					case "rate":
						cond.valueFormat = cond.value + "/秒"
						return
				}
			})
		}
	},
	template: `<div>
	<span v-for="(cond, index) in condsConfig.conds">
		<span class="ui label basic small">
			<span>{{paramMap[cond.param].name}} 
				<span v-if="paramMap[cond.param].operators != null && paramMap[cond.param].operators.length > 0"><span class="grey">{{operatorMap[cond.operator]}}</span> {{cond.valueFormat}}</span> 
			</span>
		</span>
		<span v-if="index < condsConfig.conds.length - 1"> &nbsp;<span v-if="condsConfig.connector == 'and'">且</span><span v-else>或</span>&nbsp; </span>
	</span>
</div>`
})

Vue.component("dns-route-selector", {
	props: ["v-all-routes", "v-routes"],
	data: function () {
		let routes = this.vRoutes
		if (routes == null) {
			routes = []
		}
		routes.$sort(function (v1, v2) {
			if (v1.domainId == v2.domainId) {
				return v1.code < v2.code
			}
			return (v1.domainId < v2.domainId) ? 1 : -1
		})
		return {
			routes: routes,
			routeCodes: routes.$map(function (k, v) {
				return v.code + "@" + v.domainId
			}),
			isAdding: false,
			routeCode: "",
			keyword: "",
			searchingRoutes: this.vAllRoutes.$copy()
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.keyword = ""
			this.routeCode = ""

			let that = this
			setTimeout(function () {
				that.$refs.keywordRef.focus()
			}, 200)
		},
		cancel: function () {
			this.isAdding = false
		},
		confirm: function () {
			if (this.routeCode.length == 0) {
				return
			}
			if (this.routeCodes.$contains(this.routeCode)) {
				teaweb.warn("已经添加过此线路，不能重复添加")
				return
			}
			let that = this
			let route = this.vAllRoutes.$find(function (k, v) {
				return v.code + "@" + v.domainId == that.routeCode
			})
			if (route == null) {
				return
			}

			this.routeCodes.push(this.routeCode)
			this.routes.push(route)

			this.routes.$sort(function (v1, v2) {
				if (v1.domainId == v2.domainId) {
					return v1.code < v2.code
				}
				return (v1.domainId < v2.domainId) ? 1 : -1
			})

			this.routeCode = ""
			this.isAdding = false
		},
		remove: function (route) {
			this.routeCodes.$removeValue(route.code + "@" + route.domainId)
			this.routes.$removeIf(function (k, v) {
				return v.code + "@" + v.domainId == route.code + "@" + route.domainId
			})
		},
		clearKeyword: function () {
			this.keyword = ""
		}
	},
	watch: {
		keyword: function (keyword) {
			if (keyword.length == 0) {
				this.searchingRoutes = this.vAllRoutes.$copy()
				this.routeCode = ""
				return
			}
			this.searchingRoutes = this.vAllRoutes.filter(function (route) {
				return teaweb.match(route.name, keyword) || teaweb.match(route.code, keyword) || teaweb.match(route.domainName, keyword)
			})
			if (this.searchingRoutes.length > 0) {
				this.routeCode = this.searchingRoutes[0].code + "@" + this.searchingRoutes[0].domainId
			} else {
				this.routeCode = ""
			}
		}
	},
	template: `<div>
	<input type="hidden" name="dnsRoutesJSON" :value="JSON.stringify(routeCodes)"/>
	<div v-if="routes.length > 0">
		<tiny-basic-label v-for="route in routes" :key="route.code + '@' + route.domainId">
			{{route.name}} <span class="grey small">（{{route.domainName}}）</span><a href="" @click.prevent="remove(route)"><i class="icon remove"></i></a>
		</tiny-basic-label>
		<div class="ui divider"></div>
	</div>
	<button type="button" class="ui button small" @click.prevent="add" v-if="!isAdding">+</button>
	<div v-if="isAdding">
		<table class="ui table">
			<tr>
				<td class="title">所有线路</td>
				<td>
					<span v-if="keyword.length > 0 && searchingRoutes.length == 0">没有和关键词“{{keyword}}”匹配的线路</span>
					<span v-show="keyword.length == 0 || searchingRoutes.length > 0">
						<select class="ui dropdown" v-model="routeCode">
							<option value="" v-if="keyword.length == 0">[请选择]</option>
							<option v-for="route in searchingRoutes" :value="route.code + '@' + route.domainId">{{route.name}}（{{route.code}}/{{route.domainName}}）</option>
						</select>
					</span>
				</td>
			</tr>
			<tr>
				<td>搜索线路</td>
				<td>
					<div class="ui input" :class="{'right labeled':keyword.length > 0}">
						<input type="text" placeholder="线路名称或代号..." size="10" style="width: 10em" v-model="keyword" ref="keywordRef" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
						<a class="ui label" v-if="keyword.length > 0" @click.prevent="clearKeyword" href=""><i class="icon remove small blue"></i></a>
					</div>
				</td>
			</tr>
		</table>
		
		<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button> &nbsp; <a href="" @click.prevent="cancel()"><i class="icon remove small"></i></a>
	</div>
</div>`
})

Vue.component("dns-domain-selector", {
	props: ["v-domain-id", "v-domain-name", "v-provider-name"],
	data: function () {
		let domainId = this.vDomainId
		if (domainId == null) {
			domainId = 0
		}
		let domainName = this.vDomainName
		if (domainName == null) {
			domainName = ""
		}

		let providerName = this.vProviderName
		if (providerName == null) {
			providerName = ""
		}

		return {
			domainId: domainId,
			domainName: domainName,
			providerName: providerName
		}
	},
	methods: {
		select: function () {
			let that = this
			teaweb.popup("/dns/domains/selectPopup", {
				callback: function (resp) {
					that.domainId = resp.data.domainId
					that.domainName = resp.data.domainName
					that.providerName = resp.data.providerName
					that.change()
				}
			})
		},
		remove: function() {
			this.domainId = 0
			this.domainName = ""
			this.change()
		},
		update: function () {
			let that = this
			teaweb.popup("/dns/domains/selectPopup?domainId=" + this.domainId, {
				callback: function (resp) {
					that.domainId = resp.data.domainId
					that.domainName = resp.data.domainName
					that.providerName = resp.data.providerName
					that.change()
				}
			})
		},
		change: function () {
			this.$emit("change", {
				id: this.domainId,
				name: this.domainName
			})
		}
	},
	template: `<div>
	<input type="hidden" name="dnsDomainId" :value="domainId"/>
	<div v-if="domainName.length > 0">
		<span class="ui label small basic">
			<span v-if="providerName != null && providerName.length > 0">{{providerName}} &raquo; </span> {{domainName}}
			<a href="" @click.prevent="update"><i class="icon pencil small"></i></a>
			<a href="" @click.prevent="remove()"><i class="icon remove"></i></a>
		</span>
	</div>
	<div v-if="domainName.length == 0">
		<a href="" @click.prevent="select()">[选择域名]</a>
	</div>
</div>`
})

Vue.component("dns-resolver-config-box", {
	props:["v-dns-resolver-config"],
	data: function () {
		let config = this.vDnsResolverConfig
		if (config == null) {
			config = {
				type: "default"
			}
		}
		return {
			config: config,
			types: [
				{
					name: "默认",
					code: "default"
				},
				{
					name: "CGO",
					code: "cgo"
				},
				{
					name: "Go原生",
					code: "goNative"
				},
			]
		}
	},
	template: `<div>
	<input type="hidden" name="dnsResolverJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">使用的DNS解析库</td>
			<td>
				<select class="ui dropdown auto-width" v-model="config.type">
					<option v-for="t in types" :value="t.code">{{t.name}}</option>
				</select>
				<p class="comment">边缘节点使用的DNS解析库。修改此项配置后，需要重启节点进程才会生效。<pro-warning-label></pro-warning-label></p>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("dns-resolvers-config-box", {
	props: ["value", "name"],
	data: function () {
		let resolvers = this.value
		if (resolvers == null) {
			resolvers = []
		}

		let name = this.name
		if (name == null || name.length == 0) {
			name = "dnsResolversJSON"
		}

		return {
			formName: name,
			resolvers: resolvers,

			host: "",

			isAdding: false
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.hostRef.focus()
			})
		},
		confirm: function () {
			let host = this.host.trim()
			if (host.length == 0) {
				let that = this
				setTimeout(function () {
					that.$refs.hostRef.focus()
				})
				return
			}
			this.resolvers.push({
				host: host,
				port: 0, // TODO
				protocol: "" // TODO
			})
			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.host = ""
			this.port = 0
			this.protocol = ""
		},
		remove: function (index) {
			this.resolvers.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" :name="formName" :value="JSON.stringify(resolvers)"/>
	<div v-if="resolvers.length > 0">
		<div v-for="(resolver, index) in resolvers" class="ui label basic small">
			<span v-if="resolver.protocol.length > 0">{{resolver.protocol}}</span>{{resolver.host}}<span v-if="resolver.port > 0">:{{resolver.port}}</span>
			&nbsp;
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
	</div>
	
	<div v-if="isAdding" style="margin-top: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="x.x.x.x" @keyup.enter="confirm" @keypress.enter.prevent="1" ref="hostRef" v-model="host"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确认</button>
				&nbsp; <a href="" @click.prevent="cancel" title="取消"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	
	<div v-if="!isAdding" style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("ad-instance-objects-box", {
	props: ["v-objects", "v-user-id"],
	mounted: function () {
		this.getUserServers(1)
	},
	data: function () {
		let objects = this.vObjects
		if (objects == null) {
			objects = []
		}

		let objectCodes = []
		objects.forEach(function (v) {
			objectCodes.push(v.code)
		})

		return {
			userId: this.vUserId,
			objects: objects,
			objectCodes: objectCodes,
			isAdding: true,

			servers: [],
			serversIsLoading: false
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
		},
		cancel: function () {
			this.isAdding = false
		},
		remove: function (index) {
			let that = this
			teaweb.confirm("确定要删除此防护对象吗？", function () {
				that.objects.$remove(index)
				that.notifyChange()
			})
		},
		removeObjectCode: function (objectCode) {
			let index = -1
			this.objectCodes.forEach(function (v, k) {
				if (objectCode == v) {
					index = k
				}
			})
			if (index >= 0) {
				this.objects.$remove(index)
				this.notifyChange()
			}
		},
		getUserServers: function (page) {
			if (Tea.Vue == null) {
				let that = this
				setTimeout(function () {
					that.getUserServers(page)
				}, 100)
				return
			}

			let that = this
			this.serversIsLoading = true
			Tea.Vue.$post(".userServers")
				.params({
					userId: this.userId,
					page: page,
					pageSize: 5
				})
				.success(function (resp) {
					that.servers = resp.data.servers

					that.$refs.serverPage.updateMax(resp.data.page.max)
					that.serversIsLoading = false
				})
				.error(function () {
					that.serversIsLoading = false
				})
		},
		changeServerPage: function (page) {
			this.getUserServers(page)
		},
		selectServerObject: function (server) {
			if (this.existObjectCode("server:" + server.id)) {
				return
			}

			this.objects.push({
				"type": "server",
				"code": "server:" + server.id,

				"id": server.id,
				"name": server.name
			})
			this.notifyChange()
		},
		notifyChange: function () {
			let objectCodes = []
			this.objects.forEach(function (v) {
				objectCodes.push(v.code)
			})
			this.objectCodes = objectCodes
		},
		existObjectCode: function (objectCode) {
			let found = false
			this.objects.forEach(function (v) {
				if (v.code == objectCode) {
					found = true
				}
			})
			return found
		}
	},
	template: `<div>
	<input type="hidden" name="objectCodesJSON" :value="JSON.stringify(objectCodes)"/>
	
	<!-- 已有对象 -->
	<div>
		<div v-if="objects.length == 0"><span class="grey">暂时还没有设置任何防护对象。</span></div>
		<div v-if="objects.length > 0">
			<table class="ui table">
				<tr>
					<td class="title">已选中防护对象</td>
					<td>
						<div v-for="(object, index) in objects" class="ui label basic small" style="margin-bottom: 0.5em">
							<span v-if="object.type == 'server'">网站：{{object.name}}</span>
							&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
						</div>
					</td>
				</tr>
			</table>
		</div>
	</div>
	<div class="margin"></div>
	
	<!-- 添加表单 -->
	<div v-if="isAdding">
		<table class="ui table celled">
			<tr>
				<td class="title">对象类型</td>
				<td>网站</td>
			</tr>
			<!-- 网站列表 -->
			<tr>
				<td>网站列表</td>
				<td>
					<span v-if="serversIsLoading">加载中...</span>
					<div v-if="!serversIsLoading && servers.length == 0">暂时还没有可选的网站。</div>
					<table class="ui table" v-show="!serversIsLoading && servers.length > 0">
						<thead class="full-width">
							<tr>
								<th>网站名称</th>
								<th class="one op">操作</th>
							</tr>	
						</thead>
						<tr v-for="server in servers">
							<td style="background: white">{{server.name}}</td>
							<td>
								<a href="" @click.prevent="selectServerObject(server)" v-if="!existObjectCode('server:' + server.id)">选中</a>
								<a href="" @click.prevent="removeObjectCode('server:' + server.id)" v-else><span class="red">取消</span></a>
							</td>
						</tr>
					</table>
					
					<js-page ref="serverPage" @change="changeServerPage"></js-page>
				</td>
			</tr>
		</table>
	</div>
	
	<!-- 添加按钮 -->
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("grant-selector", {
	props: ["v-grant", "v-node-cluster-id", "v-ns-cluster-id"],
	data: function () {
		return {
			grantId: (this.vGrant == null) ? 0 : this.vGrant.id,
			grant: this.vGrant,
			nodeClusterId: (this.vNodeClusterId != null) ? this.vNodeClusterId : 0,
			nsClusterId: (this.vNsClusterId != null) ? this.vNsClusterId : 0
		}
	},
	methods: {
		// 选择授权
		select: function () {
			let that = this;
			teaweb.popup("/clusters/grants/selectPopup?nodeClusterId=" + this.nodeClusterId + "&nsClusterId=" + this.nsClusterId, {
				callback: (resp) => {
					that.grantId = resp.data.grant.id;
					if (that.grantId > 0) {
						that.grant = resp.data.grant;
					}
					that.notifyUpdate()
				},
				height: "26em"
			})
		},

		// 创建授权
		create: function () {
			let that = this
			teaweb.popup("/clusters/grants/createPopup", {
				height: "26em",
				callback: (resp) => {
					that.grantId = resp.data.grant.id;
					if (that.grantId > 0) {
						that.grant = resp.data.grant;
					}
					that.notifyUpdate()
				}
			})
		},

		// 修改授权
		update: function () {
			if (this.grant == null) {
				window.location.reload();
				return;
			}
			let that = this
			teaweb.popup("/clusters/grants/updatePopup?grantId=" + this.grant.id, {
				height: "26em",
				callback: (resp) => {
					that.grant = resp.data.grant
					that.notifyUpdate()
				}
			})
		},

		// 删除已选择授权
		remove: function () {
			this.grant = null
			this.grantId = 0
			this.notifyUpdate()
		},
		notifyUpdate: function () {
			this.$emit("change", this.grant)
		}
	},
	template: `<div>
	<input type="hidden" name="grantId" :value="grantId"/>
	<div class="ui label small basic" v-if="grant != null">{{grant.name}}<span class="small grey">（{{grant.methodName}}）</span><span class="small grey" v-if="grant.username != null && grant.username.length > 0">（{{grant.username}}）</span> <a href="" title="修改" @click.prevent="update()"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="remove()"><i class="icon remove"></i></a> </div>
	<div v-if="grant == null">
		<a href="" @click.prevent="select()">[选择已有认证]</a> &nbsp; &nbsp; <a href="" @click.prevent="create()">[添加新认证]</a>
	</div>
</div>`
})

window.REQUEST_COND_COMPONENTS = [{"type":"url-extension","name":"文件扩展名","description":"根据URL中的文件路径扩展名进行过滤","component":"http-cond-url-extension","paramsTitle":"扩展名列表","isRequest":true,"caseInsensitive":false},{"type":"url-eq-index","name":"首页","description":"检查URL路径是为\"/\"","component":"http-cond-url-eq-index","paramsTitle":"URL完整路径","isRequest":true,"caseInsensitive":false},{"type":"url-all","name":"全站","description":"全站所有URL","component":"http-cond-url-all","paramsTitle":"URL完整路径","isRequest":true,"caseInsensitive":false},{"type":"url-prefix","name":"URL目录前缀","description":"根据URL中的文件路径前缀进行过滤","component":"http-cond-url-prefix","paramsTitle":"URL目录前缀","isRequest":true,"caseInsensitive":true},{"type":"url-eq","name":"URL完整路径","description":"检查URL中的文件路径是否一致","component":"http-cond-url-eq","paramsTitle":"URL完整路径","isRequest":true,"caseInsensitive":true},{"type":"url-regexp","name":"URL正则匹配","description":"使用正则表达式检查URL中的文件路径是否一致","component":"http-cond-url-regexp","paramsTitle":"正则表达式","isRequest":true,"caseInsensitive":true},{"type":"url-wildcard-match","name":"URL通配符","description":"使用通配符检查URL中的文件路径是否一致","component":"http-cond-url-wildcard-match","paramsTitle":"通配符","isRequest":true,"caseInsensitive":true},{"type":"user-agent-regexp","name":"User-Agent正则匹配","description":"使用正则表达式检查User-Agent中是否含有某些浏览器和系统标识","component":"http-cond-user-agent-regexp","paramsTitle":"正则表达式","isRequest":true,"caseInsensitive":true},{"type":"params","name":"参数匹配","description":"根据参数值进行匹配","component":"http-cond-params","paramsTitle":"参数配置","isRequest":true,"caseInsensitive":false},{"type":"url-not-extension","name":"排除：URL扩展名","description":"根据URL中的文件路径扩展名进行过滤","component":"http-cond-url-not-extension","paramsTitle":"扩展名列表","isRequest":true,"caseInsensitive":false},{"type":"url-not-prefix","name":"排除：URL前缀","description":"根据URL中的文件路径前缀进行过滤","component":"http-cond-url-not-prefix","paramsTitle":"URL前缀","isRequest":true,"caseInsensitive":true},{"type":"url-not-eq","name":"排除：URL完整路径","description":"检查URL中的文件路径是否一致","component":"http-cond-url-not-eq","paramsTitle":"URL完整路径","isRequest":true,"caseInsensitive":true},{"type":"url-not-regexp","name":"排除：URL正则匹配","description":"使用正则表达式检查URL中的文件路径是否一致，如果一致，则不匹配","component":"http-cond-url-not-regexp","paramsTitle":"正则表达式","isRequest":true,"caseInsensitive":true},{"type":"user-agent-not-regexp","name":"排除：User-Agent正则匹配","description":"使用正则表达式检查User-Agent中是否含有某些浏览器和系统标识，如果含有，则不匹配","component":"http-cond-user-agent-not-regexp","paramsTitle":"正则表达式","isRequest":true,"caseInsensitive":true},{"type":"mime-type","name":"内容MimeType","description":"根据服务器返回的内容的MimeType进行过滤。注意：当用于缓存条件时，此条件需要结合别的请求条件使用。","component":"http-cond-mime-type","paramsTitle":"MimeType列表","isRequest":false,"caseInsensitive":false}];

window.REQUEST_COND_OPERATORS = [{"description":"判断是否正则表达式匹配","name":"正则表达式匹配","op":"regexp"},{"description":"判断是否正则表达式不匹配","name":"正则表达式不匹配","op":"not regexp"},{"description":"判断是否和指定的通配符匹配","name":"通配符匹配","op":"wildcard match"},{"description":"判断是否和指定的通配符不匹配","name":"通配符不匹配","op":"wildcard not match"},{"description":"使用字符串对比参数值是否相等于某个值","name":"字符串等于","op":"eq"},{"description":"参数值包含某个前缀","name":"字符串前缀","op":"prefix"},{"description":"参数值包含某个后缀","name":"字符串后缀","op":"suffix"},{"description":"参数值包含另外一个字符串","name":"字符串包含","op":"contains"},{"description":"参数值不包含另外一个字符串","name":"字符串不包含","op":"not contains"},{"description":"使用字符串对比参数值是否不相等于某个值","name":"字符串不等于","op":"not"},{"description":"判断参数值在某个列表中","name":"在列表中","op":"in"},{"description":"判断参数值不在某个列表中","name":"不在列表中","op":"not in"},{"description":"判断小写的扩展名（不带点）在某个列表中","name":"扩展名","op":"file ext"},{"description":"判断MimeType在某个列表中，支持类似于image/*的语法","name":"MimeType","op":"mime type"},{"description":"判断版本号在某个范围内，格式为version1,version2","name":"版本号范围","op":"version range"},{"description":"将参数转换为整数数字后进行对比","name":"整数等于","op":"eq int"},{"description":"将参数转换为可以有小数的浮点数字进行对比","name":"浮点数等于","op":"eq float"},{"description":"将参数转换为数字进行对比","name":"数字大于","op":"gt"},{"description":"将参数转换为数字进行对比","name":"数字大于等于","op":"gte"},{"description":"将参数转换为数字进行对比","name":"数字小于","op":"lt"},{"description":"将参数转换为数字进行对比","name":"数字小于等于","op":"lte"},{"description":"对整数参数值取模，除数为10，对比值为余数","name":"整数取模10","op":"mod 10"},{"description":"对整数参数值取模，除数为100，对比值为余数","name":"整数取模100","op":"mod 100"},{"description":"对整数参数值取模，对比值格式为：除数,余数，比如10,1","name":"整数取模","op":"mod"},{"description":"将参数转换为IP进行对比","name":"IP等于","op":"eq ip"},{"description":"将参数转换为IP进行对比","name":"IP大于","op":"gt ip"},{"description":"将参数转换为IP进行对比","name":"IP大于等于","op":"gte ip"},{"description":"将参数转换为IP进行对比","name":"IP小于","op":"lt ip"},{"description":"将参数转换为IP进行对比","name":"IP小于等于","op":"lte ip"},{"description":"IP在某个范围之内，范围格式可以是英文逗号分隔的\u003ccode-label\u003e开始IP,结束IP\u003c/code-label\u003e，比如\u003ccode-label\u003e192.168.1.100,192.168.2.200\u003c/code-label\u003e，或者CIDR格式的ip/bits，比如\u003ccode-label\u003e192.168.2.1/24\u003c/code-label\u003e","name":"IP范围","op":"ip range"},{"description":"对IP参数值取模，除数为10，对比值为余数","name":"IP取模10","op":"ip mod 10"},{"description":"对IP参数值取模，除数为100，对比值为余数","name":"IP取模100","op":"ip mod 100"},{"description":"对IP参数值取模，对比值格式为：除数,余数，比如10,1","name":"IP取模","op":"ip mod"}];

window.REQUEST_VARIABLES = [{"code":"${edgeVersion}","description":"","name":"边缘节点版本"},{"code":"${remoteAddr}","description":"会依次根据X-Forwarded-For、X-Real-IP、RemoteAddr获取，适合前端有别的反向代理服务时使用，存在伪造的风险","name":"客户端地址（IP）"},{"code":"${rawRemoteAddr}","description":"返回直接连接服务的客户端原始IP地址","name":"客户端地址（IP）"},{"code":"${remotePort}","description":"","name":"客户端端口"},{"code":"${remoteUser}","description":"","name":"客户端用户名"},{"code":"${requestURI}","description":"比如/hello?name=lily","name":"请求URI"},{"code":"${requestPath}","description":"比如/hello","name":"请求路径（不包括参数）"},{"code":"${requestURL}","description":"比如https://example.com/hello?name=lily","name":"完整的请求URL"},{"code":"${requestLength}","description":"","name":"请求内容长度"},{"code":"${requestMethod}","description":"比如GET、POST","name":"请求方法"},{"code":"${requestFilename}","description":"","name":"请求文件路径"},{"code":"${requestPathExtension}","description":"请求路径中的文件扩展名，包括点符号，比如.html、.png","name":"请求文件扩展名"},{"code":"${requestPathLowerExtension}","description":"请求路径中的文件扩展名，其中大写字母会被自动转换为小写，包括点符号，比如.html、.png","name":"请求文件小写扩展名"},{"code":"${scheme}","description":"","name":"请求协议，http或https"},{"code":"${proto}","description:":"类似于HTTP/1.0","name":"包含版本的HTTP请求协议"},{"code":"${timeISO8601}","description":"比如2018-07-16T23:52:24.839+08:00","name":"ISO 8601格式的时间"},{"code":"${timeLocal}","description":"比如17/Jul/2018:09:52:24 +0800","name":"本地时间"},{"code":"${msec}","description":"比如1531756823.054","name":"带有毫秒的时间"},{"code":"${timestamp}","description":"","name":"unix时间戳，单位为秒"},{"code":"${host}","description":"","name":"主机名"},{"code":"${cname}","description":"比如38b48e4f.goedge.cn","name":"当前网站的CNAME"},{"code":"${serverName}","description":"","name":"接收请求的服务器名"},{"code":"${serverPort}","description":"","name":"接收请求的服务器端口"},{"code":"${referer}","description":"","name":"请求来源URL"},{"code":"${referer.host}","description":"","name":"请求来源URL域名"},{"code":"${userAgent}","description":"","name":"客户端信息"},{"code":"${contentType}","description":"","name":"请求头部的Content-Type"},{"code":"${cookies}","description":"","name":"所有cookie组合字符串"},{"code":"${cookie.NAME}","description":"","name":"单个cookie值"},{"code":"${isArgs}","description":"如果URL有参数，则值为`?`；否则，则值为空","name":"问号（?）标记"},{"code":"${args}","description":"","name":"所有参数组合字符串"},{"code":"${arg.NAME}","description":"","name":"单个参数值"},{"code":"${headers}","description":"","name":"所有Header信息组合字符串"},{"code":"${header.NAME}","description":"","name":"单个Header值"},{"code":"${geo.country.name}","description":"","name":"国家/地区名称"},{"code":"${geo.country.id}","description":"","name":"国家/地区ID"},{"code":"${geo.province.name}","description":"目前只包含中国省份","name":"省份名称"},{"code":"${geo.province.id}","description":"目前只包含中国省份","name":"省份ID"},{"code":"${geo.city.name}","description":"目前只包含中国城市","name":"城市名称"},{"code":"${geo.city.id}","description":"目前只包含中国城市","name":"城市名称"},{"code":"${isp.name}","description":"","name":"ISP服务商名称"},{"code":"${isp.id}","description":"","name":"ISP服务商ID"},{"code":"${browser.os.name}","description":"客户端所在操作系统名称","name":"操作系统名称"},{"code":"${browser.os.version}","description":"客户端所在操作系统版本","name":"操作系统版本"},{"code":"${browser.name}","description":"客户端浏览器名称","name":"浏览器名称"},{"code":"${browser.version}","description":"客户端浏览器版本","name":"浏览器版本"},{"code":"${browser.isMobile}","description":"如果客户端是手机，则值为1，否则为0","name":"手机标识"}];

window.METRIC_HTTP_KEYS = [{"name":"客户端地址（IP）","code":"${remoteAddr}","description":"会依次根据X-Forwarded-For、X-Real-IP、RemoteAddr获取，适用于前端可能有别的反向代理的情形，存在被伪造的可能","icon":""},{"name":"直接客户端地址（IP）","code":"${rawRemoteAddr}","description":"返回直接连接服务的客户端原始IP地址","icon":""},{"name":"客户端用户名","code":"${remoteUser}","description":"通过基本认证填入的用户名","icon":""},{"name":"请求URI","code":"${requestURI}","description":"包含参数，比如/hello?name=lily","icon":""},{"name":"请求路径","code":"${requestPath}","description":"不包含参数，比如/hello","icon":""},{"name":"完整URL","code":"${requestURL}","description":"比如https://example.com/hello?name=lily","icon":""},{"name":"请求方法","code":"${requestMethod}","description":"比如GET、POST等","icon":""},{"name":"请求协议Scheme","code":"${scheme}","description":"http或https","icon":""},{"name":"文件扩展名","code":"${requestPathExtension}","description":"请求路径中的文件扩展名，包括点符号，比如.html、.png","icon":""},{"name":"小写文件扩展名","code":"${requestPathLowerExtension}","description":"请求路径中的文件扩展名小写形式，包括点符号，比如.html、.png","icon":""},{"name":"主机名","code":"${host}","description":"通常是请求的域名","icon":""},{"name":"HTTP协议","code":"${proto}","description":"包含版本的HTTP请求协议，类似于HTTP/1.0","icon":""},{"name":"URL参数值","code":"${arg.NAME}","description":"单个URL参数值","icon":""},{"name":"请求来源URL","code":"${referer}","description":"请求来源Referer URL","icon":""},{"name":"请求来源URL域名","code":"${referer.host}","description":"请求来源Referer URL域名","icon":""},{"name":"Header值","code":"${header.NAME}","description":"单个Header值，比如${header.User-Agent}","icon":""},{"name":"Cookie值","code":"${cookie.NAME}","description":"单个cookie值，比如${cookie.sid}","icon":""},{"name":"状态码","code":"${status}","description":"","icon":""},{"name":"响应的Content-Type值","code":"${response.contentType}","description":"","icon":""}];

window.IP_ADDR_THRESHOLD_ITEMS = [{"code":"nodeAvgRequests","description":"当前节点在单位时间内接收到的平均请求数。","name":"节点平均请求数","unit":"个"},{"code":"nodeAvgTrafficOut","description":"当前节点在单位时间内发送的下行流量。","name":"节点平均下行流量","unit":"M"},{"code":"nodeAvgTrafficIn","description":"当前节点在单位时间内接收的上行流量。","name":"节点平均上行流量","unit":"M"},{"code":"nodeHealthCheck","description":"当前节点健康检查结果。","name":"节点健康检查结果","unit":""},{"code":"connectivity","description":"通过区域监控得到的当前IP地址的连通性数值，取值在0和100之间。","name":"IP连通性","unit":"%"},{"code":"groupAvgRequests","description":"当前节点所在分组在单位时间内接收到的平均请求数。","name":"分组平均请求数","unit":"个"},{"code":"groupAvgTrafficOut","description":"当前节点所在分组在单位时间内发送的下行流量。","name":"分组平均下行流量","unit":"M"},{"code":"groupAvgTrafficIn","description":"当前节点所在分组在单位时间内接收的上行流量。","name":"分组平均上行流量","unit":"M"},{"code":"clusterAvgRequests","description":"当前节点所在集群在单位时间内接收到的平均请求数。","name":"集群平均请求数","unit":"个"},{"code":"clusterAvgTrafficOut","description":"当前节点所在集群在单位时间内发送的下行流量。","name":"集群平均下行流量","unit":"M"},{"code":"clusterAvgTrafficIn","description":"当前节点所在集群在单位时间内接收的上行流量。","name":"集群平均上行流量","unit":"M"}];

window.IP_ADDR_THRESHOLD_ACTIONS = [{"code":"up","description":"上线当前IP。","name":"上线"},{"code":"down","description":"下线当前IP。","name":"下线"},{"code":"notify","description":"发送已达到阈值通知。","name":"通知"},{"code":"switch","description":"在DNS中记录中将IP切换到指定的备用IP。","name":"切换"},{"code":"webHook","description":"调用外部的WebHook。","name":"WebHook"}];

window.WAF_RULE_CHECKPOINTS = [{"description":"通用报头比如Cache-Control、Accept之类的长度限制，防止缓冲区溢出攻击。","name":"通用请求报头长度限制","prefix":"requestGeneralHeaderLength"},{"description":"通用报头比如Cache-Control、Date之类的长度限制，防止缓冲区溢出攻击。","name":"通用响应报头长度限制","prefix":"responseGeneralHeaderLength"},{"description":"试图通过分析X-Forwarded-For等报头获取的客户端地址，比如192.168.1.100，存在伪造的可能。","name":"客户端地址（IP）","prefix":"remoteAddr"},{"description":"直接连接的客户端地址，比如192.168.1.100。","name":"客户端源地址（IP）","prefix":"rawRemoteAddr"},{"description":"直接连接的客户端地址端口。","name":"客户端端口","prefix":"remotePort"},{"description":"通过BasicAuth登录的客户端用户名。","name":"客户端用户名","prefix":"remoteUser"},{"description":"包含URL参数的请求URI，类似于 /hello/world?lang=go，不包含域名部分。","name":"请求URI","prefix":"requestURI"},{"description":"不包含URL参数的请求路径，类似于 /hello/world，不包含域名部分。","name":"请求路径","prefix":"requestPath"},{"description":"完整的请求URL，包含协议、域名、请求路径、参数等，类似于 https://example.com/hello?name=lily 。","name":"请求完整URL","prefix":"requestURL"},{"description":"请求报头中的Content-Length。","name":"请求内容长度","prefix":"requestLength"},{"description":"通常在POST或者PUT等操作时会附带请求体，最大限制32M。","name":"请求体内容","prefix":"requestBody"},{"description":"${requestURI}和${requestBody}组合。","name":"请求URI和请求体组合","prefix":"requestAll"},{"description":"获取POST或者其他方法发送的表单参数，最大请求体限制32M。","name":"请求表单参数","prefix":"requestForm"},{"description":"获取POST上传的文件信息，最大请求体限制32M。","name":"上传文件","prefix":"requestUpload"},{"description":"获取POST或者其他方法发送的JSON，最大请求体限制32M，使用点（.）符号表示多级数据。","name":"请求JSON参数","prefix":"requestJSON"},{"description":"比如GET、POST。","name":"请求方法","prefix":"requestMethod"},{"description":"比如http或https。","name":"请求协议","prefix":"scheme"},{"description":"比如HTTP/1.1。","name":"HTTP协议版本","prefix":"proto"},{"description":"比如example.com。","name":"主机名","prefix":"host"},{"description":"当前网站服务CNAME，比如38b48e4f.example.com。","name":"CNAME","prefix":"cname"},{"description":"是否为CNAME，值为1（是）或0（否）。","name":"是否为CNAME","prefix":"isCNAME"},{"description":"请求报头中的Referer和Origin值。","name":"请求来源","prefix":"refererOrigin"},{"description":"请求报头中的Referer值。","name":"请求来源Referer","prefix":"referer"},{"description":"比如Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103。","name":"客户端信息","prefix":"userAgent"},{"description":"请求报头的Content-Type。","name":"内容类型","prefix":"contentType"},{"description":"比如sid=IxZVPFhE\u0026city=beijing\u0026uid=18237。","name":"所有cookie组合字符串","prefix":"cookies"},{"description":"单个cookie值。","name":"单个cookie值","prefix":"cookie"},{"description":"比如name=lu\u0026age=20。","name":"所有URL参数组合","prefix":"args"},{"description":"单个URL参数值。","name":"单个URL参数值","prefix":"arg"},{"description":"使用换行符（\\n）隔开的报头内容字符串，每行均为\"NAME: VALUE格式\"。","name":"所有请求报头内容","prefix":"headers"},{"description":"使用换行符（\\n）隔开的报头名称字符串，每行一个名称。","name":"所有请求报头名称","prefix":"headerNames"},{"description":"单个报头值。","name":"单个请求报头值","prefix":"header"},{"description":"最长的请求报头的长度。","name":"请求报头最大长度","prefix":"headerMaxLength"},{"description":"当前客户端所处国家/地区名称。","name":"国家/地区名称","prefix":"geoCountryName"},{"description":"当前客户端所处中国省份名称。","name":"省份名称","prefix":"geoProvinceName"},{"description":"当前客户端所处中国城市名称。","name":"城市名称","prefix":"geoCityName"},{"description":"当前客户端所处ISP名称。","name":"ISP名称","prefix":"ispName"},{"description":"对统计对象进行统计。","name":"CC统计","prefix":"cc2"},{"description":"对统计对象进行统计。","name":"防盗链","prefix":"refererBlock"},{"description":"统计某段时间段内的请求信息（不推荐再使用，请使用新的CC2统计代替）。","name":"CC统计（旧）","prefix":"cc"},{"description":"响应状态码，比如200、404、500。","name":"响应状态码","prefix":"status"},{"description":"响应报头值。","name":"响应报头","prefix":"responseHeader"},{"description":"响应内容字符串。","name":"响应内容","prefix":"responseBody"},{"description":"响应内容长度，通过响应的报头Content-Length获取。","name":"响应内容长度","prefix":"bytesSent"}];

window.WAF_RULE_OPERATORS = [{"name":"正则匹配","code":"match","description":"使用正则表达式匹配，在头部使用(?i)表示不区分大小写，\u003ca href=\"https://goedge.cloud/docs/Appendix/Regexp/Index.md\" target=\"_blank\"\u003e正则表达式语法 \u0026raquo;\u003c/a\u003e。","caseInsensitive":"yes","dataType":"regexp"},{"name":"正则不匹配","code":"not match","description":"使用正则表达式不匹配，在头部使用(?i)表示不区分大小写，\u003ca href=\"https://goedge.cloud/docs/Appendix/Regexp/Index.md\" target=\"_blank\"\u003e正则表达式语法 \u0026raquo;\u003c/a\u003e。","caseInsensitive":"yes","dataType":"regexp"},{"name":"通配符匹配","code":"wildcard match","description":"判断是否和指定的通配符匹配，可以在对比值中使用星号通配符（*）表示任意字符。","caseInsensitive":"yes","dataType":"wildcard"},{"name":"通配符不匹配","code":"wildcard not match","description":"判断是否和指定的通配符不匹配，可以在对比值中使用星号通配符（*）表示任意字符。","caseInsensitive":"yes","dataType":"wildcard"},{"name":"字符串等于","code":"eq string","description":"使用字符串对比等于。","caseInsensitive":"no","dataType":"string"},{"name":"字符串不等于","code":"neq string","description":"使用字符串对比不等于。","caseInsensitive":"no","dataType":"string"},{"name":"包含字符串","code":"contains","description":"包含某个字符串，比如Hello World包含了World。","caseInsensitive":"no","dataType":"string"},{"name":"不包含字符串","code":"not contains","description":"不包含某个字符串，比如Hello字符串中不包含Hi。","caseInsensitive":"no","dataType":"string"},{"name":"包含任一字符串","code":"contains any","description":"包含字符串列表中的任意一个，比如/hello/world包含/hello和/hi中的/hello，对比值中每行一个字符串。","caseInsensitive":"no","dataType":"strings"},{"name":"包含所有字符串","code":"contains all","description":"包含字符串列表中的所有字符串，比如/hello/world必须包含/hello和/world，对比值中每行一个字符串。","caseInsensitive":"no","dataType":"strings"},{"name":"包含前缀","code":"prefix","description":"包含字符串前缀部分，比如/hello前缀会匹配/hello, /hello/world等。","caseInsensitive":"no","dataType":"string"},{"name":"包含后缀","code":"suffix","description":"包含字符串后缀部分，比如/hello后缀会匹配/hello, /hi/hello等。","caseInsensitive":"no","dataType":"string"},{"name":"包含任一单词","code":"contains any word","description":"包含某个独立单词，对比值中每行一个单词，比如mozilla firefox里包含了mozilla和firefox两个单词，但是不包含fire和fox这两个单词。","caseInsensitive":"no","dataType":"strings"},{"name":"包含所有单词","code":"contains all words","description":"包含所有的独立单词，对比值中每行一个单词，比如mozilla firefox里包含了mozilla和firefox两个单词，但是不包含fire和fox这两个单词。","caseInsensitive":"no","dataType":"strings"},{"name":"不包含任一单词","code":"not contains any word","description":"不包含某个独立单词，对比值中每行一个单词，比如mozilla firefox里包含了mozilla和firefox两个单词，但是不包含fire和fox这两个单词。","caseInsensitive":"no","dataType":"strings"},{"name":"包含SQL注入","code":"contains sql injection","description":"检测字符串内容是否包含SQL注入。","caseInsensitive":"none","dataType":"none"},{"name":"包含SQL注入-严格模式","code":"contains sql injection strictly","description":"更加严格地检测字符串内容是否包含SQL注入，相对于非严格模式，有一定的误报几率。","caseInsensitive":"none","dataType":"none"},{"name":"包含XSS注入","code":"contains xss","description":"检测字符串内容是否包含XSS注入。","caseInsensitive":"none","dataType":"none"},{"name":"包含XSS注入-严格模式","code":"contains xss strictly","description":"更加严格地检测字符串内容是否包含XSS注入，相对于非严格模式，此时xml、audio、video等标签也会被匹配。","caseInsensitive":"none","dataType":"none"},{"name":"包含二进制数据","code":"contains binary","description":"包含一组二进制数据。","caseInsensitive":"no","dataType":"string"},{"name":"不包含二进制数据","code":"not contains binary","description":"不包含一组二进制数据。","caseInsensitive":"no","dataType":"string"},{"name":"数值大于","code":"gt","description":"使用数值对比大于，对比值需要是一个数字。","caseInsensitive":"none","dataType":"number"},{"name":"数值大于等于","code":"gte","description":"使用数值对比大于等于，对比值需要是一个数字。","caseInsensitive":"none","dataType":"number"},{"name":"数值小于","code":"lt","description":"使用数值对比小于，对比值需要是一个数字。","caseInsensitive":"none","dataType":"number"},{"name":"数值小于等于","code":"lte","description":"使用数值对比小于等于，对比值需要是一个数字。","caseInsensitive":"none","dataType":"number"},{"name":"数值等于","code":"eq","description":"使用数值对比等于，对比值需要是一个数字。","caseInsensitive":"none","dataType":"number"},{"name":"数值不等于","code":"neq","description":"使用数值对比不等于，对比值需要是一个数字。","caseInsensitive":"none","dataType":"number"},{"name":"包含索引","code":"has key","description":"对于一组数据拥有某个键值或者索引。","caseInsensitive":"no","dataType":"string|number"},{"name":"版本号大于","code":"version gt","description":"对比版本号大于。","caseInsensitive":"none","dataType":"version"},{"name":"版本号小于","code":"version lt","description":"对比版本号小于。","caseInsensitive":"none","dataType":"version"},{"name":"版本号范围","code":"version range","description":"判断版本号在某个范围内，格式为 起始version1,结束version2。","caseInsensitive":"none","dataType":"versionRange"},{"name":"IP等于","code":"eq ip","description":"将参数转换为IP进行对比，只能对比单个IP。","caseInsensitive":"none","dataType":"ip"},{"name":"在一组IP中","code":"in ip list","description":"判断参数IP在一组IP内，对比值中每行一个IP。","caseInsensitive":"none","dataType":"ips"},{"name":"IP大于","code":"gt ip","description":"将参数转换为IP进行对比。","caseInsensitive":"none","dataType":"ip"},{"name":"IP大于等于","code":"gte ip","description":"将参数转换为IP进行对比。","caseInsensitive":"none","dataType":"ip"},{"name":"IP小于","code":"lt ip","description":"将参数转换为IP进行对比。","caseInsensitive":"none","dataType":"ip"},{"name":"IP小于等于","code":"lte ip","description":"将参数转换为IP进行对比。","caseInsensitive":"none","dataType":"ip"},{"name":"IP范围","code":"ip range","description":"IP在某个范围之内，范围格式可以是英文逗号分隔的\u003ccode-label\u003e开始IP,结束IP\u003c/code-label\u003e，比如\u003ccode-label\u003e192.168.1.100,192.168.2.200\u003c/code-label\u003e；或者CIDR格式的ip/bits，比如\u003ccode-label\u003e192.168.2.1/24\u003c/code-label\u003e；或者单个IP。可以填写多行，每行一个IP范围。","caseInsensitive":"none","dataType":"ips"},{"name":"不在IP范围","code":"not ip range","description":"IP不在某个范围之内，范围格式可以是英文逗号分隔的\u003ccode-label\u003e开始IP,结束IP\u003c/code-label\u003e，比如\u003ccode-label\u003e192.168.1.100,192.168.2.200\u003c/code-label\u003e；或者CIDR格式的ip/bits，比如\u003ccode-label\u003e192.168.2.1/24\u003c/code-label\u003e；或者单个IP。可以填写多行，每行一个IP范围。","caseInsensitive":"none","dataType":"ips"},{"name":"IP取模10","code":"ip mod 10","description":"对IP参数值取模，除数为10，对比值为余数。","caseInsensitive":"none","dataType":"number"},{"name":"IP取模100","code":"ip mod 100","description":"对IP参数值取模，除数为100，对比值为余数。","caseInsensitive":"none","dataType":"number"},{"name":"IP取模","code":"ip mod","description":"对IP参数值取模，对比值格式为：除数,余数，比如10,1。","caseInsensitive":"none","dataType":"number"}];

window.WAF_CAPTCHA_TYPES = [{"name":"验证码","code":"default","description":"通过输入验证码来验证人机。","icon":""},{"name":"点击验证","code":"oneClick","description":"通过点击界面元素来验证人机。","icon":""},{"name":"滑动解锁","code":"slide","description":"通过滑动方块解锁来验证人机。","icon":""},{"name":"极验-行为验","code":"geetest","description":"使用极验-行为验提供的人机验证方式。","icon":""}];

