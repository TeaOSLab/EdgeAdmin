Vue.component("ns-cluster-combo-box", {
	props: ["v-cluster-id", "name"],
	data: function () {
		let that = this
		Tea.action("/ns/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})


		let inputName = "clusterId"
		if (this.name != null && this.name.length > 0) {
			inputName = this.name
		}

		return {
			clusters: [],
			inputName: inputName
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
	<combo-box title="集群" placeholder="集群名称" :v-items="clusters" :name="inputName" :v-value="vClusterId" @change="change"></combo-box>
</div>`
})