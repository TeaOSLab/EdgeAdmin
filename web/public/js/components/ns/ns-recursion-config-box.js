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