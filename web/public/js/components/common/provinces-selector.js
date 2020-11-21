Vue.component("provinces-selector", {
	props: ["v-provinces"],
	data: function () {
		let provinces = this.vProvinces
		if (provinces == null) {
			provinces = []
		}
		let provinceIds = provinces.$map(function (k, v) {
			return v.id
		})
		return {
			provinces: provinces,
			provinceIds: provinceIds
		}
	},
	methods: {
		add: function () {
			let provinceStringIds = this.provinceIds.map(function (v) {
				return v.toString()
			})
			let that = this
			teaweb.popup("/ui/selectProvincesPopup?provinceIds=" + provinceStringIds.join(","), {
				width: "48em",
				height: "23em",
				callback: function (resp) {
					that.provinces = resp.data.provinces
					that.change()
				}
			})
		},
		remove: function (index) {
			this.provinces.$remove(index)
			this.change()
		},
		change: function () {
			this.provinceIds = this.provinces.$map(function (k, v) {
				return v.id
			})
		}
	},
	template: `<div>
	<input type="hidden" name="provinceIdsJSON" :value="JSON.stringify(provinceIds)"/>
	<div v-if="provinces.length > 0" style="margin-bottom: 0.5em">
		<div v-for="(province, index) in provinces" class="ui label tiny basic">{{province.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})