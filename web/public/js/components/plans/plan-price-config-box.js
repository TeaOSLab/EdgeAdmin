// 套餐价格配置
Vue.component("plan-price-config-box", {
	props: ["v-price-type", "v-monthly-price", "v-seasonally-price", "v-yearly-price", "v-traffic-price"],
	data: function () {
		let priceType = this.vPriceType
		if (priceType == null) {
			priceType = "period"
		}

		let monthlyPriceNumber = 0
		let monthlyPrice = this.vMonthlyPrice
		if (monthlyPrice == null || monthlyPrice <= 0) {
			monthlyPrice = ""
		} else {
			monthlyPrice = monthlyPrice.toString()
			monthlyPriceNumber = parseFloat(monthlyPrice)
			if (isNaN(monthlyPriceNumber)) {
				monthlyPriceNumber = 0
			}
		}

		let seasonallyPriceNumber = 0
		let seasonallyPrice = this.vSeasonallyPrice
		if (seasonallyPrice == null || seasonallyPrice <= 0) {
			seasonallyPrice = ""
		} else {
			seasonallyPrice = seasonallyPrice.toString()
			seasonallyPriceNumber = parseFloat(seasonallyPrice)
			if (isNaN(seasonallyPriceNumber)) {
				seasonallyPriceNumber = 0
			}
		}

		let yearlyPriceNumber = 0
		let yearlyPrice = this.vYearlyPrice
		if (yearlyPrice == null || yearlyPrice <= 0) {
			yearlyPrice = ""
		} else {
			yearlyPrice = yearlyPrice.toString()
			yearlyPriceNumber = parseFloat(yearlyPrice)
			if (isNaN(yearlyPriceNumber)) {
				yearlyPriceNumber = 0
			}
		}

		let trafficPrice = this.vTrafficPrice
		let trafficPriceBaseNumber = 0
		if (trafficPrice != null) {
			trafficPriceBaseNumber = trafficPrice.base
		} else {
			trafficPrice = {
				base: 0
			}
		}
		let trafficPriceBase = ""
		if (trafficPriceBaseNumber > 0) {
			trafficPriceBase = trafficPriceBaseNumber.toString()
		}

		return {
			priceType: priceType,
			monthlyPrice: monthlyPrice,
			seasonallyPrice: seasonallyPrice,
			yearlyPrice: yearlyPrice,

			monthlyPriceNumber: monthlyPriceNumber,
			seasonallyPriceNumber: seasonallyPriceNumber,
			yearlyPriceNumber: yearlyPriceNumber,

			trafficPriceBase: trafficPriceBase,
			trafficPrice: trafficPrice
		}
	},
	watch: {
		monthlyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.monthlyPriceNumber = price
		},
		seasonallyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.seasonallyPriceNumber = price
		},
		yearlyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.yearlyPriceNumber = price
		},
		trafficPriceBase: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.trafficPrice.base = price
		}
	},
	template: `<div>
	<input type="hidden" name="priceType" :value="priceType"/>
	<input type="hidden" name="monthlyPrice" :value="monthlyPriceNumber"/>
	<input type="hidden" name="seasonallyPrice" :value="seasonallyPriceNumber"/>
	<input type="hidden" name="yearlyPrice" :value="yearlyPriceNumber"/>
	<input type="hidden" name="trafficPriceJSON" :value="JSON.stringify(trafficPrice)"/>
	
	<div>
		<radio :v-value="'period'" :value="priceType" v-model="priceType">&nbsp;按时间周期</radio> &nbsp; &nbsp;
		<radio :v-value="'traffic'" :value="priceType" v-model="priceType">&nbsp;按流量</radio>
	</div>
	
	<!-- 按时间周期 -->
	<div v-show="priceType == 'period'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">月度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="monthlyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">季度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="seasonallyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">年度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="yearlyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
		</table>
	</div>
	
	<!-- 按流量 -->
	<div v-show="priceType =='traffic'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">基础流量费用</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" v-model="trafficPriceBase" maxlength="10" style="width: 7em"/>
						<span class="ui label">元/GB</span>
					</div>
				</td>
			</tr>
		</table>
	</div>
</div>`
})