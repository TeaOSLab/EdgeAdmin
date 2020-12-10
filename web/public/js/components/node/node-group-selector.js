Vue.component("node-group-selector", {
	props: ["v-cluster-id", "v-group"],
	data: function () {
		return {
			selectedGroup: this.vGroup
		}
	},
	methods: {
		selectGroup: function () {
			let that = this
			teaweb.popup("/clusters/cluster/groups/selectPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedGroup = resp.data.group
				}
			})
		},
		addGroup: function () {
			let that = this
			teaweb.popup("/clusters/cluster/groups/createPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedGroup = resp.data.group
				}
			})
		},
		removeGroup: function () {
			this.selectedGroup = null
		}
	},
	template: `<div>
	<div class="ui label small basic" v-if="selectedGroup != null">
		<input type="hidden" name="groupId" :value="selectedGroup.id"/>
		{{selectedGroup.name}} &nbsp;<a href="" title="删除" @click.prevent="removeGroup()"><i class="icon remove"></i></a>
	</div>
	<div v-if="selectedGroup == null">
		<a href="" @click.prevent="selectGroup()">[选择分组]</a> &nbsp; <a href="" @click.prevent="addGroup()">[添加分组]</a>
	</div>
</div>`
})