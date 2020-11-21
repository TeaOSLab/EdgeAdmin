Vue.component("countries-selector", {
	props: ["v-countries"],
	data: function () {
		let countries = this.vCountries
		if (countries == null) {
			countries = []
		}
		let countryIds = countries.$map(function (k, v) {
			return v.id
		})
		return {
			countries: countries,
			countryIds: countryIds
		}
	},
	methods: {
		add: function () {
			let countryStringIds = this.countryIds.map(function (v) {
				return v.toString()
			})
			let that = this
			teaweb.popup("/ui/selectCountriesPopup?countryIds=" + countryStringIds.join(","), {
				width: "48em",
				height: "23em",
				callback: function (resp) {
					that.countries = resp.data.countries
					that.change()
				}
			})
		},
		remove: function (index) {
			this.countries.$remove(index)
			this.change()
		},
		change: function () {
			this.countryIds = this.countries.$map(function (k, v) {
				return v.id
			})
		}
	},
	template: `<div>
	<input type="hidden" name="countryIdsJSON" :value="JSON.stringify(countryIds)"/>
	<div v-if="countries.length > 0" style="margin-bottom: 0.5em">
		<div v-for="(country, index) in countries" class="ui label tiny basic">{{country.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})