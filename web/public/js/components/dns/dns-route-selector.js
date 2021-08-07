Vue.component("dns-route-selector", {
	props: ["v-all-routes", "v-routes"],
	data: function () {
		let routes = this.vRoutes
		if (routes == null) {
			routes = []
		}
		return {
			routes: routes,
			routeCodes: routes.$map(function (k, v) {
				return v.code + "@" + v.domainId
			}),
			isAdding: false,
			routeCode: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
		},
		cancel: function () {
			this.isAdding = false
		},
		confirm: function () {
			if (this.routeCode.length == 0) {
				return
			}
			if (this.routeCodes.$contains(this.routeCode)) {
				teaweb.warn("已经添加过此线路，不能重复添加")
				return
			}
			let that = this
			let route = this.vAllRoutes.$find(function (k, v) {
				return v.code + "@" + v.domainId == that.routeCode
			})
			if (route == null) {
				return
			}

			this.routeCodes.push(this.routeCode)
			this.routes.push(route)

			this.routeCode = ""
			this.isAdding = false
		},
		remove: function (route) {
			this.routeCodes.$removeValue(route.code + "@" + route.domainId)
			this.routes.$removeIf(function (k, v) {
				return v.code + "@" + v.domainId == route.code + "@" + route.domainId
			})
		}
	},
	template: `<div>
	<input type="hidden" name="dnsRoutesJSON" :value="JSON.stringify(routeCodes)"/>
	<div v-if="routes.length > 0">
		<tiny-basic-label v-for="route in routes" :key="route.code + '@' + route.domainId">
			{{route.name}} <span class="grey small">（{{route.domainName}}）</span><a href="" @click.prevent="remove(route)"><i class="icon remove"></i></a>
		</tiny-basic-label>
		<div class="ui divider"></div>
	</div>
	<button type="button" class="ui button small" @click.prevent="add" v-if="!isAdding">+</button>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown auto-width" v-model="routeCode">
					<option value="">[请选择]</option>
					<option v-for="route in vAllRoutes" :value="route.code + '@' + route.domainId">{{route.name}}（{{route.domainName}}）</option>
				</select>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
			</div>
			<div class="ui field">
				<a href="" @click.prevent="cancel()"><i class="icon remove"></i></a>
			</div>
		</div>
	</div>
</div>`
})