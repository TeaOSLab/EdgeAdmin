Vue.component("user-selector", {
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
	<combo-box placeholder="选择用户" :data-url="'/servers/users/options'" :data-key="'users'" name="userId" :v-value="userId"></combo-box>
</div>`
})