Vue.component("plan-price-view", {
	props: ["v-plan"],
	data: function () {
		return {
			plan: this.vPlan
		}
	},
	template: `<div>
	 <span v-if="plan.priceType == 'period'">
		<span v-if="plan.monthlyPrice > 0">月度：￥{{plan.monthlyPrice}}元<br/></span>
		<span v-if="plan.seasonallyPrice > 0">季度：￥{{plan.seasonallyPrice}}元<br/></span>
		<span v-if="plan.yearlyPrice > 0">年度：￥{{plan.yearlyPrice}}元</span>
	</span>
	<span v-if="plan.priceType == 'traffic'">
		基础价格：￥{{plan.trafficPrice.base}}元/GB
	</span>
</div>`
})