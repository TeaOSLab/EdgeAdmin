// 套餐价格配置
Vue.component("plan-price-config-box", {
	props: ["v-price-type", "v-monthly-price", "v-seasonally-price", "v-yearly-price", "v-traffic-price", "v-bandwidth-price", "v-disable-period"],
	data: function () {
		let priceType = this.vPriceType
		if (priceType == null) {
			priceType = "bandwidth"
		}

		// 按时间周期计费
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

		// 按流量计费
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

		// 按带宽计费
		let bandwidthPrice = this.vBandwidthPrice
		if (bandwidthPrice == null) {
			bandwidthPrice = {
				percentile: 95,
				ranges: []
			}
		} else if (bandwidthPrice.ranges == null) {
			bandwidthPrice.ranges = []
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
			trafficPrice: trafficPrice,

			bandwidthPrice: bandwidthPrice,
			bandwidthPercentile: bandwidthPrice.percentile
		}
	},
	methods: {
		changeBandwidthPriceRanges: function (ranges) {
			this.bandwidthPrice.ranges = ranges
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
		},
		bandwidthPercentile: function (v) {
			let percentile = parseInt(v)
			if (isNaN(percentile) || percentile <= 0) {
				percentile = 95
			} else if (percentile > 100) {
				percentile = 100
			}
			this.bandwidthPrice.percentile = percentile
		}
	},
	template: `<div>
	<input type="hidden" name="priceType" :value="priceType"/>
	<input type="hidden" name="monthlyPrice" :value="monthlyPriceNumber"/>
	<input type="hidden" name="seasonallyPrice" :value="seasonallyPriceNumber"/>
	<input type="hidden" name="yearlyPrice" :value="yearlyPriceNumber"/>
	<input type="hidden" name="trafficPriceJSON" :value="JSON.stringify(trafficPrice)"/>
	<input type="hidden" name="bandwidthPriceJSON" :value="JSON.stringify(bandwidthPrice)"/>
	
	<div>
		<radio :v-value="'bandwidth'" :value="priceType" v-model="priceType">&nbsp;按带宽</radio> &nbsp;
		<radio :v-value="'traffic'" :value="priceType" v-model="priceType">&nbsp;按流量</radio> &nbsp;
		<radio :v-value="'period'" :value="priceType" v-model="priceType" v-show="typeof(vDisablePeriod) != 'boolean' || !vDisablePeriod">&nbsp;按时间周期</radio>
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
				<td class="title">基础流量费用 *</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" v-model="trafficPriceBase" maxlength="10" style="width: 7em"/>
						<span class="ui label">元/GB</span>
					</div>
				</td>
			</tr>
		</table>
	</div>
	
	<!-- 按带宽 -->
	<div v-show="priceType == 'bandwidth'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">带宽百分位 *</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 4em" maxlength="3" v-model="bandwidthPercentile"/>
						<span class="ui label">th</span>
					</div>
				</td>
			</tr>
			<tr>
				<td>带宽价格</td>
				<td>
					<plan-bandwidth-ranges :v-ranges="bandwidthPrice.ranges" @change="changeBandwidthPriceRanges"></plan-bandwidth-ranges>
				</td>
			</tr>
		</table>
	</div>
</div>`
})