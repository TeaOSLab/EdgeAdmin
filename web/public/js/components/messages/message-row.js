Vue.component("message-row", {
	props: ["v-message"],
	data: function () {
		let paramsJSON = this.vMessage.params
		let params = null
		if (paramsJSON != null && paramsJSON.length > 0) {
			params = JSON.parse(paramsJSON)
		}

		return {
			message: this.vMessage,
			params: params
		}
	},
	methods: {
		viewCert: function (certId) {
			teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
				height: "28em",
				width: "48em"
			})
		}
	},
	template: `<div>
<table class="ui table selectable">
	<tr :class="{error: message.level == 'error', positive: message.level == 'success', warning: message.level == 'warning'}">
		<td>
			<strong>{{message.datetime}}</strong>
			<span v-if="message.cluster != null && message.cluster.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster?clusterId=' + message.cluster.id">集群：{{message.cluster.name}}</a>
			</span>
			<span v-if="message.node != null && message.node.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id">节点：{{message.node.name}}</a>
			</span>
		</td>
	</tr>
	<tr :class="{error: message.level == 'error', positive: message.level == 'success', warning: message.level == 'warning'}">
		<td>
			{{message.body}}
			
			<!-- 健康检查 -->
			<div v-if="message.type == 'HealthCheckFailed'" style="margin-top: 0.8em">
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + param.node.id" v-for="param in params" class="ui label small basic" style="margin-bottom: 0.5em">{{param.node.name}}: {{param.error}}</a>
			</div>
			
			<!-- 集群DNS设置 -->
			<div v-if="message.type == 'ClusterDNSSyncFailed'" style="margin-top: 0.8em">
				<a :href="'/dns/clusters/cluster?clusterId=' + message.cluster.id">查看问题 &raquo;</a>
			</div>
			
			<!-- 证书即将过期 -->
			<div v-if="message.type == 'SSLCertExpiring'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0">查看任务&raquo;</a>
			</div>
			
			<!-- 证书续期成功 -->
			<div v-if="message.type == 'SSLCertACMETaskSuccess'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0">查看任务&raquo;</a>
			</div>
			<div v-if="message.type == 'SSLCertACMETaskFailed'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0">查看任务&raquo;</a>
			</div>
		</td>
	</tr>
</table>
<div class="margin"></div>
</div>`
})