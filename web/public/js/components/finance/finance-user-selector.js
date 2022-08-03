Vue.component("finance-user-selector", {
	props: ["v-user-id"],
	data: function () {
		return {}
	},
	methods: {
		change: function (userId) {
			this.$emit("change", userId)
		}
	},
	template: `<div>
	<user-selector :v-user-id="vUserId" data-url="/finance/users/options" @change="change"></user-selector>
</div>`
})