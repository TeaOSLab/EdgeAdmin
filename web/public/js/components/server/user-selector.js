Vue.component("user-selector", {
	props: ["v-user-id", "data-url"],
	data: function () {
		let userId = this.vUserId
		if (userId == null) {
			userId = 0
		}

		let dataURL = this.dataUrl
		if (dataURL == null || dataURL.length == 0) {
			dataURL = "/servers/users/options"
		}

		return {
			users: [],
			userId: userId,
			dataURL: dataURL
		}
	},
	methods: {
		change: function(item) {
			if (item != null) {
				this.$emit("change", item.id)
			} else {
				this.$emit("change", 0)
			}
		},
		clear: function () {
			this.$refs.comboBox.clear()
		}
	},
	template: `<div>
	<combo-box placeholder="选择用户" :data-url="dataURL" :data-key="'users'" data-search="on" name="userId" :v-value="userId" @change="change" ref="comboBox"></combo-box>
</div>`
})