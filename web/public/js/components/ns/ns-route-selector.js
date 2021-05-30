// 选择单一线路
Vue.component("ns-route-selector", {
	props: ["v-route-id"],
	mounted: function () {
		let that = this
		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.routes = resp.data.routes
			})
	},
	data: function () {
		let routeId = this.vRouteId
		if (routeId == null) {
			routeId = 0
		}
		return {
			routeId: routeId,
			routes: []
		}
	},
	template: `<div>
	<div v-if="routes.length > 0">
		<select class="ui dropdown" name="routeId" v-model="routeId">
			<option value="0">[线路]</option>
			<option v-for="route in routes" :value="route.id">{{route.name}}</option>
		</select>
	</div>
</div>`
})