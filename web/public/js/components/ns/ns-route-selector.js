// 选择单一线路
Vue.component("ns-route-selector", {
	props: ["v-route-code"],
	mounted: function () {
		let that = this
		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.routes = resp.data.routes
			})
	},
	data: function () {
		let routeCode = this.vRouteCode
		if (routeCode == null) {
			routeCode = ""
		}
		return {
			routeCode: routeCode,
			routes: []
		}
	},
	template: `<div>
	<div v-if="routes.length > 0">
		<select class="ui dropdown" name="routeCode" v-model="routeCode">
			<option value="">[线路]</option>
			<option v-for="route in routes" :value="route.code">{{route.name}}</option>
		</select>
	</div>
</div>`
})