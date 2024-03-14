// 选择多个线路
Vue.component("ns-routes-selector", {
	props: ["v-routes", "name"],
	mounted: function () {
		let that = this
		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.routes = resp.data.routes

				// provinces
				let provinces = {}
				if (resp.data.provinces != null && resp.data.provinces.length > 0) {
					for (const province of resp.data.provinces) {
						let countryCode = province.countryCode
						if (typeof provinces[countryCode] == "undefined") {
							provinces[countryCode] = []
						}
						provinces[countryCode].push({
							name: province.name,
							code: province.code
						})
					}
				}
				that.provinces = provinces
			})
	},
	data: function () {
		let selectedRoutes = this.vRoutes
		if (selectedRoutes == null) {
			selectedRoutes = []
		}

		let inputName = this.name
		if (typeof inputName != "string" || inputName.length == 0) {
			inputName = "routeCodes"
		}

		return {
			routeCode: "default",
			inputName: inputName,
			routes: [],

			provinces: {}, // country code => [ province1, province2, ... ]
			provinceRouteCode: "",

			isAdding: false,
			routeType: "default",
			selectedRoutes: selectedRoutes,
		}
	},
	watch: {
		routeType: function (v) {
			this.routeCode = ""
			let that = this
			this.routes.forEach(function (route) {
				if (route.type == v && that.routeCode.length == 0) {
					that.routeCode = route.code
				}
			})
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.routeType = "default"
			this.routeCode = "default"
			this.provinceRouteCode = ""
			this.$emit("add")
		},
		cancel: function () {
			this.isAdding = false
			this.$emit("cancel")
		},
		confirm: function () {
			if (this.routeCode.length == 0) {
				return
			}

			let that = this

			// route
			let selectedRoute = null
			for (const route of this.routes) {
				if (route.code == this.routeCode) {
					selectedRoute = route
					break
				}
			}

			if (selectedRoute != null) {
				// province route
				if (this.provinceRouteCode.length > 0 && this.provinces[this.routeCode] != null) {
					for (const province of this.provinces[this.routeCode]) {
						if (province.code == this.provinceRouteCode) {
							selectedRoute = {
								name: selectedRoute.name + "-" + province.name,
								code: province.code
							}
							break
						}
					}
				}

				that.selectedRoutes.push(selectedRoute)
			}

			this.$emit("change", this.selectedRoutes)
			this.cancel()
		},
		remove: function (index) {
			this.selectedRoutes.$remove(index)
			this.$emit("change", this.selectedRoutes)
		}
	}
	,
	template: `<div>
	<div v-show="selectedRoutes.length > 0">
		<div class="ui label basic text small" v-for="(route, index) in selectedRoutes" style="margin-bottom: 0.3em">
			<input type="hidden" :name="inputName" :value="route.code"/>
			{{route.name}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding" style="margin-bottom: 1em">
		<table class="ui table">
			<tr>
				<td class="title">选择类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="routeType">
						<option value="default">[默认线路]</option>
						<option value="user">自定义线路</option>
						<option value="isp">运营商</option>
						<option value="china">中国省市</option>
						<option value="world">全球国家地区</option>
						<option value="agent">搜索引擎</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>选择线路 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="routeCode">
						<option v-for="route in routes" :value="route.code" v-if="route.type == routeType">{{route.name}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="routeCode.length > 0 && provinces[routeCode] != null">
				<td>选择省/州</td>
				<td>
					<select class="ui dropdown auto-width" v-model="provinceRouteCode">
						<option value="">[全域]</option>
						<option v-for="province in provinces[routeCode]" :value="province.code">{{province.name}}</option>
					</select>
				</td>
			</tr>
		</table>
		<div>
			<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
			&nbsp; <a href="" title="取消" @click.prevent="cancel">取消</a>
		</div>	
	</div>
	<button class="ui button tiny" type="button" @click.prevent="add" v-if="!isAdding">+</button>
</div>`
})