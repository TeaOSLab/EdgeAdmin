Vue.component("http-firewall-province-selector", {
	props: ["v-type", "v-provinces"],
	data: function () {
		let provinces = this.vProvinces
		if (provinces == null) {
			provinces = []
		}

		return {
			listType: this.vType,
			provinces: provinces
		}
	},
	methods: {
		addProvince: function () {
			let selectedProvinceIds = this.provinces.map(function (province) {
				return province.id
			})
			let that = this
			teaweb.popup("/servers/server/settings/waf/ipadmin/selectProvincesPopup?type=" + this.listType + "&selectedProvinceIds=" + selectedProvinceIds.join(","), {
				width: "50em",
				height: "26em",
				callback: function (resp) {
					that.provinces = resp.data.selectedProvinces
					that.$forceUpdate()
					that.notifyChange()
				}
			})
		},
		removeProvince: function (index) {
			this.provinces.$remove(index)
			this.notifyChange()
		},
		resetProvinces: function () {
			this.provinces = []
			this.notifyChange()
		},
		notifyChange: function () {
			this.$emit("change", {
				"provinces": this.provinces
			})
		}
	},
	template: `<div>
	<span v-if="provinces.length == 0" class="disabled">暂时没有选择<span v-if="listType =='allow'">允许</span><span v-else>封禁</span>的省份。</span>
	<div v-show="provinces.length > 0">
		<div class="ui label tiny basic" v-for="(province, index) in provinces" style="margin-bottom: 0.5em">
			<input type="hidden" :name="listType + 'ProvinceIds'" :value="province.id"/>
			{{province.name}} <a href="" @click.prevent="removeProvince(index)" title="删除"><i class="icon remove"></i></a>
		</div>
	</div>
	<div class="ui divider"></div>
	<button type="button" class="ui button tiny" @click.prevent="addProvince">修改</button> &nbsp; <button type="button" class="ui button tiny" v-show="provinces.length > 0" @click.prevent="resetProvinces">清空</button>
</div>`
})