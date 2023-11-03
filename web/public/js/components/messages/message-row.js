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