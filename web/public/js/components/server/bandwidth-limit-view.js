// 显示带宽限制说明
Vue.component("bandwidth-limit-view", {
	props: ["v-bandwidth-limit"],
	data: function () {
		return {
			config: this.vBandwidthLimit
		}
	},
	template: `<div>
	<div v-if="config.isOn">
		<span v-if="config.dailySize != null && config.dailySize.count > 0">日带宽限制：{{config.dailySize.count}}{{config.dailySize.unit.toUpperCase()}}<br/></span>
		<span v-if="config.monthlySize != null && config.monthlySize.count > 0">月带宽限制：{{config.monthlySize.count}}{{config.monthlySize.unit.toUpperCase()}}<br/></span>
	</div>
	<span v-else class="disabled">没有限制。</span>
</div>`
})