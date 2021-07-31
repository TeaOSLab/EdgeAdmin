// 单个集群选择
Vue.component("cluster-selector", {
	mounted: function () {
		let that = this

		Tea.action("/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})
	},
	props: ["v-cluster-id"],
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