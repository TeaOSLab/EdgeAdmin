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