// 指标周期设置
Vue.component("metric-period-config-box", {
	props: ["v-period", "v-period-unit"],
	data: function () {
		let period = this.vPeriod
		let periodUnit = this.vPeriodUnit
		if (period == null || period.toString().length == 0) {
			period = 1
		}
		if (periodUnit == null || periodUnit.length == 0) {
			periodUnit = "day"
		}
		return {
			periodConfig: {
				period: period,
				unit: periodUnit
			}
		}
	},
	watch: {
		"periodConfig.period": function (v) {
			v = parseInt(v)
			if (isNaN(v) || v <= 0) {
				v = 1
			}
			this.periodConfig.period = v
		}
	},
	template: `<div>
	<input type="hidden" name="periodJSON" :value="JSON.stringify(periodConfig)"/>
	<div class="ui fields inline">
		<div class="ui field">
			<input type="text" v-model="periodConfig.period" maxlength="4" size="4"/>
		</div>
		<div class="ui field">
			<select class="ui dropdown" v-model="periodConfig.unit">
				<option value="minute">分钟</option>
				<option value="hour">小时</option>
				<option value="day">天</option>
				<option value="week">周</option>
				<option value="month">月</option>
			</select>
		</div>
	</div>
	<p class="comment">在此周期内同一对象累积为同一数据。</p>
</div>`
})