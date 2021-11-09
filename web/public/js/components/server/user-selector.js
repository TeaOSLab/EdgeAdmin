Vue.component("user-selector", {
	mounted: function () {
		let that = this

		Tea.action("/servers/users/options")
			.post()
			.success(function (resp) {
				that.users = resp.data.users
			})
	},
	props: ["v-user-id"],
	data: function () {
		let userId = this.vUserId
		if (userId == null) {
			userId = 0
		}
		return {
			users: [],
			userId: userId
		}
	},
	watch: {
		userId: function (v) {
			this.$emit("change", v)
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`
})