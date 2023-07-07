Vue.component("http-firewall-region-selector", {
	props: ["v-type", "v-countries"],
	data: function () {
		let countries = this.vCountries
		if (countries == null) {
			countries = []
		}

		return {
			listType: this.vType,
			countries: countries
		}
	},
	methods: {
		addCountry: function () {
			let selectedCountryIds = this.countries.map(function (country) {
				return country.id
			})
			let that = this
			teaweb.popup("/servers/server/settings/waf/ipadmin/selectCountriesPopup?type=" + this.listType + "&selectedCountryIds=" + selectedCountryIds.join(","), {
				width: "52em",
				height: "30em",
				callback: function (resp) {
					that.countries = resp.data.selectedCountries
					that.$forceUpdate()
					that.notifyChange()
				}
			})
		},
		removeCountry: function (index) {
			this.countries.$remove(index)
			this.notifyChange()
		},
		resetCountries: function () {
			this.countries = []
			this.notifyChange()
		},
		notifyChange: function () {
			this.$emit("change", {
				"countries": this.countries
			})
		}
	},
	template: `<div>
	<span v-if="countries.length == 0" class="disabled">暂时没有选择<span v-if="listType =='allow'">允许</span><span v-else>封禁</span>的区域。</span>
	<div v-show="countries.length > 0">
		<div class="ui label tiny basic" v-for="(country, index) in countries" style="margin-bottom: 0.5em">
			<input type="hidden" :name="listType + 'CountryIds'" :value="country.id"/>
			({{country.letter}}){{country.name}} <a href="" @click.prevent="removeCountry(index)" title="删除"><i class="icon remove"></i></a>
		</div>
	</div>
	<div class="ui divider"></div>
	<button type="button" class="ui button tiny" @click.prevent="addCountry">修改</button> &nbsp; <button type="button" class="ui button tiny" v-show="countries.length > 0" @click.prevent="resetCountries">清空</button>
</div>`
})