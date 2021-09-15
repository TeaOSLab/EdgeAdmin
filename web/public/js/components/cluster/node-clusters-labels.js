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