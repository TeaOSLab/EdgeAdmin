Vue.component("plan-price-view", {
	props: ["v-plan"],
	data: function () {
		return {
			plan: this.vPlan
		}
	},
	template: `<div>
	 <span v-if="plan.priceType == 'period'">
	 	按时间周期计费
	 	<div>
	 		<span class="grey small">
				<span v-if="plan.monthlyPrice > 0">月度：￥{{plan.monthlyPrice}}元<br/></span>
				<span v-if="plan.seasonallyPrice > 0">季度：￥{{plan.seasonallyPrice}}元<br/></span>
				<span v-if="plan.yearlyPrice > 0">年度：￥{{plan.yearlyPrice}}元</span>
			</span>
		</div>
	</span>
	<span v-if="plan.priceType == 'traffic'">
		按流量计费
		<div>
			<span class="grey small">基础价格：￥{{plan.trafficPrice.base}}元/GiB</span>
		</div>
	</span>
	<div v-if="plan.priceType == 'bandwidth' && plan.bandwidthPrice != null && plan.bandwidthPrice.percentile > 0">
		按{{plan.bandwidthPrice.percentile}}th带宽计费 
		<div>
			<div v-for="range in plan.bandwidthPrice.ranges">
				<span class="small grey">{{range.minMB}} - <span v-if="range.maxMB > 0">{{range.maxMB}}MiB</span><span v-else>&infin;</span>： <span v-if="range.totalPrice > 0">{{range.totalPrice}}元</span><span v-else="">{{range.pricePerMB}}元/MiB</span></span>
			</div>
		</div>
	</div>
</div>`
})