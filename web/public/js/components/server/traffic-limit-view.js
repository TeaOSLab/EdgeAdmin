// 显示流量限制说明
Vue.component("traffic-limit-view", {
	props: ["v-traffic-limit"],
	data: function () {
		return {
			config: this.vTrafficLimit
		}
	},
	template: `<div>
	<div v-if="config.isOn">
		<span v-if="config.dailySize != null && config.dailySize.count > 0">日流量限制：{{config.dailySize.count}}{{config.dailySize.unit.toUpperCase()}}<br/></span>
		<span v-if="config.monthlySize != null && config.monthlySize.count > 0">月流量限制：{{config.monthlySize.count}}{{config.monthlySize.unit.toUpperCase()}}<br/></span>
	</div>
	<span v-else class="disabled">没有限制。</span>
</div>`
})