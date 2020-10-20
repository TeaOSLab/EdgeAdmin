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
	template: `<div>
<table class="ui table selectable">
	<tr :class="{error: message.level == 'error'}">
		<td>
			<strong>{{message.datetime}}</strong>
			<span v-if="message.cluster != null && message.cluster.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster?clusterId=' + message.cluster.id">集群：{{message.cluster.name}}</a>
			</span>
		</td>
	</tr>
	<tr :class="{error: message.level == 'error'}">
		<td>
			{{message.body}}
			
			<div v-if="message.type == 'HealthCheckFail'" style="margin-top: 0.8em">
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + param.node.id" v-for="param in params" class="ui label tiny" style="margin-bottom: 0.5em">{{param.node.name}}: {{param.error}}</a>
			</div>
		</td>
	</tr>
</table>
<div class="margin"></div>
</div>`
})