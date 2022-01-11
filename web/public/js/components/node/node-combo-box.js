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