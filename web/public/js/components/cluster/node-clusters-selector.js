// 一个节点的多个集群选择器
Vue.component("node-clusters-selector", {
	props: ["v-primary-cluster", "v-secondary-clusters"],
	data: function () {
		let primaryCluster = this.vPrimaryCluster

		let secondaryClusters = this.vSecondaryClusters
		if (secondaryClusters == null) {
			secondaryClusters = []
		}

		return {
			primaryClusterId: (primaryCluster == null) ? 0 : primaryCluster.id,
			secondaryClusterIds: secondaryClusters.map(function (v) {
				return v.id
			}),

			primaryCluster: primaryCluster,
			secondaryClusters: secondaryClusters
		}
	},
	methods: {
		addPrimary: function () {
			let that = this
			let selectedClusterIds = [this.primaryClusterId].concat(this.secondaryClusterIds)
			teaweb.popup("/clusters/selectPopup?selectedClusterIds=" + selectedClusterIds.join(",") + "&mode=single", {
				height: "30em",
				width: "50em",
				callback: function (resp) {
					if (resp.data.cluster != null) {
						that.primaryCluster = resp.data.cluster
						that.primaryClusterId = that.primaryCluster.id
						that.notifyChange()
					}
				}
			})
		},
		removePrimary: function () {
			this.primaryClusterId = 0
			this.primaryCluster = null
			this.notifyChange()
		},
		addSecondary: function () {
			let that = this
			let selectedClusterIds = [this.primaryClusterId].concat(this.secondaryClusterIds)
			teaweb.popup("/clusters/selectPopup?selectedClusterIds=" + selectedClusterIds.join(",") + "&mode=multiple", {
				height: "30em",
				width: "50em",
				callback: function (resp) {
					if (resp.data.cluster != null) {
						that.secondaryClusterIds.push(resp.data.cluster.id)
						that.secondaryClusters.push(resp.data.cluster)
						that.notifyChange()
					}
				}
			})
		},
		removeSecondary: function (index) {
			this.secondaryClusterIds.$remove(index)
			this.secondaryClusters.$remove(index)
			this.notifyChange()
		},
		notifyChange: function () {
			this.$emit("change", {
				clusterId: this.primaryClusterId
			})
		}
	},
	template: `<div>
	<input type="hidden" name="primaryClusterId" :value="primaryClusterId"/>
	<input type="hidden" name="secondaryClusterIds" :value="JSON.stringify(secondaryClusterIds)"/>
	<table class="ui table">
		<tr>
			<td class="title">主集群</td>
			<td>
				<div v-if="primaryCluster != null">
					<div class="ui label basic small">{{primaryCluster.name}} &nbsp; <a href="" title="删除" @click.prevent="removePrimary"><i class="icon remove small"></i></a> </div>
				</div>
				<div style="margin-top: 0.6em" v-if="primaryClusterId == 0">
					<button class="ui button tiny" type="button" @click.prevent="addPrimary">+</button>
				</div>
				<p class="comment">多个集群配置有冲突时，优先使用主集群配置。</p>
			</td>
		</tr>
		<tr>
			<td>从集群</td>
			<td>
				<div v-if="secondaryClusters.length > 0">
					<div class="ui label basic small" v-for="(cluster, index) in secondaryClusters"><span class="grey">{{cluster.name}}</span> &nbsp; <a href="" title="删除" @click.prevent="removeSecondary(index)"><i class="icon remove small"></i></a> </div>
				</div>
				<div style="margin-top: 0.6em">
					<button class="ui button tiny" type="button" @click.prevent="addSecondary">+</button>
				</div>
			</td>
		</tr>
	</table>
</div>`
})