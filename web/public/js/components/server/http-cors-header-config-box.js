Vue.component("http-cors-header-config-box", {
	props: ["value"],
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isOn: false,
				allowMethods: [],
				allowOrigin: "",
				allowCredentials: false,
				exposeHeaders: [],
				maxAge: 0,
				requestHeaders: [],
				requestMethod: ""
			}
		}

		return {
			config: config
		}
	},
	template: `<div>
	<input type="hidden" name="corsJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">启用CORS自适应跨域</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`
})