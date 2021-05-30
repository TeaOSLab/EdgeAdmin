// 选择多个线路
Vue.component("ns-routes-selector", {
	props: ["v-route-ids"],
	mounted: function () {
		let that = this

		let routeIds = this.vRouteIds
		if (routeIds == null) {
			routeIds = []
		}

		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.allRoutes = resp.data.routes
				that.allRoutes.forEach(function (v) {
					v.isChecked = (routeIds.$contains(v.id))
				})
			})
	},
	data: function () {
		return {
			routeId: 0,
			allRoutes: [],
			routes: [],
			isAdding: false
		}
	}
	,
	methods: {
		add: function () {
			this.isAdding = true
			this.routes = this.allRoutes.$findAll(function (k, v) {
				return !v.isChecked
			})
			this.routeId = 0
		},
		cancel: function () {
			this.isAdding = false
		},
		confirm: function () {
			if (this.routeId == 0) {
				return
			}

			let that = this
			this.routes.forEach(function (v) {
				if (v.id == that.routeId) {
					v.isChecked = true
				}
			})
			this.cancel()
		},
		remove: function (index) {
			this.allRoutes[index].isChecked = false
			Vue.set(this.allRoutes, index, this.allRoutes[index])
		}
	}
	,
	template: `<div>
	<div>
		<div class="ui label basic text small" v-for="(route, index) in allRoutes" v-if="route.isChecked">
			<input type="hidden" name="routeIds" :value="route.id"/>
			{{route.name}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding" style="margin-bottom: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" name="routeId" v-model="routeId">
					<option value="0">[线路]</option>
					<option v-for="route in routes" :value="route.id">{{route.name}}</option>
				</select>
			</div>
			<div class="ui field">
				<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<button class="ui button tiny" type="button" @click.prevent="add">+</button>
</div>`
})