Vue.component("ns-access-log-ref-box", {
	props:["v-access-log-ref"],
	data: function () {
		let config = this.vAccessLogRef
		if (config == null) {
			config = {
				isOn: false,
				isPrior: false
			}
		}
		return {
			config: config
		}
	},
	template: `<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">是否启用</td>
			<td>
				<checkbox name="isOn" value="1" v-model="config.isOn"></checkbox>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`
})