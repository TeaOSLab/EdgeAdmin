Vue.component("ns-domain-group-selector", {
	props: ["v-domain-group-id"],
	data: function () {
		let groupId = this.vDomainGroupId
		if (groupId == null) {
			groupId = 0
		}
		return {
			userId: 0,
			groupId: groupId
		}
	},
	methods: {
		change: function (group) {
			if (group != null) {
				this.$emit("change", group.id)
			} else {
				this.$emit("change", 0)
			}
		},
		reload: function (userId) {
			this.userId = userId
			this.$refs.comboBox.clear()
			this.$refs.comboBox.setDataURL("/ns/domains/groups/options?userId=" + userId)
			this.$refs.comboBox.reloadData()
		}
	},
	template: `<div>
	<combo-box 
		data-url="/ns/domains/groups/options" 
		placeholder="选择分组" 
		data-key="groups" 
		name="groupId"
		:v-value="groupId" 
		@change="change"
		ref="comboBox">	
	</combo-box>
</div>`
})