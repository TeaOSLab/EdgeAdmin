Vue.component("server-group-selector", {
	props: ["v-groups"],
	data: function () {
		let groups = this.vGroups
		if (groups == null) {
			groups = []
		}
		return {
			groups: groups
		}
	},
	methods: {
		selectGroup: function () {
			let that = this
			let groupIds = this.groups.map(function (v) {
				return v.id.toString()
			}).join(",")
			teaweb.popup("/servers/groups/selectPopup?selectedGroupIds=" + groupIds, {
				callback: function (resp) {
					that.groups.push(resp.data.group)
				}
			})
		},
		addGroup: function () {
			let that = this
			teaweb.popup("/servers/groups/createPopup", {
				callback: function (resp) {
					that.groups.push(resp.data.group)
				}
			})
		},
		removeGroup: function (index) {
			this.groups.$remove(index)
		},
		groupIds: function () {
			return this.groups.map(function (v) {
				return v.id
			})
		}
	},
	template: `<div>
	<div v-if="groups.length > 0">
		<div class="ui label small basic" v-if="groups.length > 0" v-for="(group, index) in groups">
			<input type="hidden" name="groupIds" :value="group.id"/>
			{{group.name}} &nbsp;<a href="" title="删除" @click.prevent="removeGroup(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<a href="" @click.prevent="selectGroup()">[选择分组]</a> &nbsp; <a href="" @click.prevent="addGroup()">[添加分组]</a>
	</div>
</div>`
})