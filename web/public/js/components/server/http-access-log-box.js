Vue.component("http-access-log-box", {
	props: ["v-access-log", "v-keyword"],
	data: function () {
		let accessLog = this.vAccessLog
		if (accessLog.header != null && accessLog.header.Upgrade != null && accessLog.header.Upgrade.values != null && accessLog.header.Upgrade.values.$contains("websocket")) {
			if (accessLog.scheme == "http") {
				accessLog.scheme = "ws"
			} else if (accessLog.scheme == "https") {
				accessLog.scheme = "wss"
			}
		}

		return {
			accessLog: accessLog
		}
	},
	methods: {
		formatCost: function (seconds) {
			var s = (seconds * 1000).toString();
			var pieces = s.split(".");
			if (pieces.length < 2) {
				return s;
			}

			return pieces[0] + "." + pieces[1].substr(0, 3);
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
		}
	},
	template: `<div style="word-break: break-all" :style="{'color': (accessLog.status >= 400) ? '#dc143c' : ''}" ref="box">
	<a v-if="accessLog.node != null && accessLog.node.nodeCluster != null" :href="'/clusters/cluster/node?nodeId=' + accessLog.node.id + '&clusterId=' + accessLog.node.nodeCluster.id" title="点击查看节点详情" target="_top"><span class="grey">[{{accessLog.node.name}}<span v-if="!accessLog.node.name.endsWith('节点')">节点</span>]</span></a><span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey">[{{accessLog.region}}]</span> <ip-box><keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword></ip-box> [{{accessLog.timeLocal}}] <em>&quot;<keyword :v-word="vKeyword">{{accessLog.requestMethod}}</keyword> {{accessLog.scheme}}://<keyword :v-word="vKeyword">{{accessLog.host}}</keyword><keyword :v-word="vKeyword">{{accessLog.requestURI}}</keyword> <a :href="accessLog.scheme + '://' + accessLog.host + accessLog.requestURI" target="_blank" title="新窗口打开" class="disabled"><i class="external icon tiny"></i> </a> {{accessLog.proto}}&quot; </em> <keyword :v-word="vKeyword">{{accessLog.status}}</keyword> <code-label v-if="accessLog.attrs != null && accessLog.attrs['cache.status'] == 'HIT'">cache hit</code-label> <code-label v-if="accessLog.firewallActions != null && accessLog.firewallActions.length > 0">waf {{accessLog.firewallActions}}</code-label> <span v-if="accessLog.tags != null && accessLog.tags.length > 0">- <code-label v-for="tag in accessLog.tags" :key="tag">{{tag}}</code-label></span> - 耗时:{{formatCost(accessLog.requestTime)}} ms <span v-if="accessLog.humanTime != null && accessLog.humanTime.length > 0" class="grey small">&nbsp; ({{accessLog.humanTime}})</span>
	&nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>
</div>`
})