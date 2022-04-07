Vue.component("node-cluster-combo-box", {
	props: ["v-cluster-id"],
	data: function () {
		let that = this
		Tea.action("/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})
		return {
			clusters: []
		}
	},
	methods: {
		change: function (item) {
			if (item == null) {
				this.$emit("change", 0)
			} else {
				this.$emit("change", item.value)
			}
		}
	},
	template: `<div v-if="clusters.length > 0" style="min-width: 10.4em">
	<combo-box title="集群" placeholder="集群名称" :v-items="clusters" name="clusterId" :v-value="vClusterId" @change="change"></combo-box>
</div>`
})