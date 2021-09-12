// 监控节点分组选择
Vue.component("report-node-groups-selector", {
	props: ["v-group-ids"],
	mounted: function () {
		let that = this
		Tea.action("/clusters/monitors/groups/options")
			.post()
			.success(function (resp) {
				that.groups = resp.data.groups.map(function (group) {
					group.isChecked = that.groupIds.$contains(group.id)
					return group
				})
				that.isLoaded = true
			})
	},
	data: function () {
		var groupIds = this.vGroupIds
		if (groupIds == null) {
			groupIds = []
		}

		return {
			groups: [],
			groupIds: groupIds,
			isLoaded: false,
			allGroups: groupIds.length == 0
		}
	},
	methods: {
		check: function (group) {
			group.isChecked = !group.isChecked
			this.groupIds = []
			let that = this
			this.groups.forEach(function (v) {
				if (v.isChecked) {
					that.groupIds.push(v.id)
				}
			})
			this.change()
		},
		change: function () {
			let that = this
			let groups = []
			this.groupIds.forEach(function (groupId) {
				let group = that.groups.$find(function (k, v) {
					return v.id == groupId
				})
				if (group == null) {
					return
				}
				groups.push({
					id: group.id,
					name: group.name
				})
			})
			this.$emit("change", groups)
		}
	},
	watch: {
		allGroups: function (b) {
			if (b) {
				this.groupIds = []
				this.groups.forEach(function (v) {
					v.isChecked = false
				})
			}

			this.change()
		}
	},
	template: `<div>
	<input type="hidden" name="reportNodeGroupIdsJSON" :value="JSON.stringify(groupIds)"/>
	<span class="disabled" v-if="isLoaded && groups.length == 0">还没有分组。</span>
	<div v-if="groups.length > 0">
		<div>
			<div class="ui checkbox">
				<input type="checkbox" v-model="allGroups" id="all-group"/>
				<label for="all-group">全部分组</label>
			</div>
			<div class="ui divider" v-if="!allGroups"></div>
		</div>
		<div v-show="!allGroups">
			<div v-for="group in groups" :key="group.id" style="float: left; width: 7.6em; margin-bottom: 0.5em">
				<div class="ui checkbox">
					<input type="checkbox" v-model="group.isChecked" value="1" :id="'report-node-group-' + group.id" @click.prevent="check(group)"/>
					<label :for="'report-node-group-' + group.id">{{group.name}}</label>
				</div>
			</div>
		</div>
	</div>
</div>`
})