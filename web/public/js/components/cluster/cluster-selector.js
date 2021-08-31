// 单个集群选择
Vue.component("cluster-selector", {
	props: ["v-cluster-id"],
	mounted: function () {
		let that = this

		Tea.action("/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})
	},
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
	<select class="ui dropdown" style="max-width: 10em" name="clusterId" v-model="clusterId">
		<option value="0">[选择集群]</option>
		<option v-for="cluster in clusters" :value="cluster.id">{{cluster.name}}</option>
	</select>
</div>`
})