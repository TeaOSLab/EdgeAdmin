Vue.component("http-stat-config-box", {
	props: ["v-stat-config", "v-is-location"],
	data: function () {
		let stat = this.vStatConfig
		if (stat == null) {
			stat = {
				isPrior: false,
				isOn: true
			}
		}
		return {
			stat: stat
		}
	},
	template: `<div>
	<input type="hidden" name="statJSON" :value="JSON.stringify(stat)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="stat" v-if="vIsLocation" ></prior-checkbox>
		<tbody v-show="!vIsLocation || stat.isPrior">
			<tr>
				<td class="title">是否开启统计</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="stat.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
<div class="margin"></div>
</div>`
})