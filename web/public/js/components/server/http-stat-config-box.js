Vue.component("http-stat-config-box", {
	props: ["v-stat-config", "v-is-location", "v-is-group"],
	data: function () {
		let stat = this.vStatConfig
		if (stat == null) {
			stat = {
				isPrior: false,
				isOn: false
			}
		}
		return {
			stat: stat
		}
	},
	template: `<div>
	<input type="hidden" name="statJSON" :value="JSON.stringify(stat)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="stat" v-if="vIsLocation || vIsGroup" ></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || stat.isPrior">
			<tr>
				<td class="title">启用统计</td>
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