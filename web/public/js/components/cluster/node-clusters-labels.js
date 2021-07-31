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
		if (labelSize == "tiny") {
			labelSize += " olive"
		}
		return {
			cluster: cluster,
			secondaryClusters: secondaryClusters,
			labelSize: labelSize
		}
	},
	template: `<div>
	<a v-if="cluster != null" :href="'/clusters/cluster?clusterId=' + cluster.id" class="ui label basic" :class="labelSize" title="主集群" style="margin-bottom: 0.3em;">{{cluster.name}}</a>
	<a v-for="c in secondaryClusters" :href="'/clusters/cluster?clusterId=' + c.id" class="ui label basic" :class="labelSize" title="从集群" style="margin-bottom: 0.3em;"><span class="grey" style="text-decoration: none">{{c.name}}</span></a>
</div>`
})