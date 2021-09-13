Vue.component("dns-route-selector", {
	props: ["v-all-routes", "v-routes"],
	data: function () {
		let routes = this.vRoutes
		if (routes == null) {
			routes = []
		}
		routes.$sort(function (v1, v2) {
			if (v1.domainId == v2.domainId) {
				return v1.code < v2.code
			}
			return (v1.domainId < v2.domainId) ? 1 : -1
		})
		return {
			routes: routes,
			routeCodes: routes.$map(function (k, v) {
				return v.code + "@" + v.domainId
			}),
			isAdding: false,
			routeCode: "",
			keyword: "",
			searchingRoutes: this.vAllRoutes.$copy()
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.keyword = ""
			this.routeCode = ""

			let that = this
			setTimeout(function () {
				that.$refs.keywordRef.focus()
			}, 200)
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

			this.routes.$sort(function (v1, v2) {
				if (v1.domainId == v2.domainId) {
					return v1.code < v2.code
				}
				return (v1.domainId < v2.domainId) ? 1 : -1
			})

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
	watch: {
		keyword: function (keyword) {
			if (keyword.length == 0) {
				this.searchingRoutes = this.vAllRoutes.$copy()
				this.routeCode = ""
				return
			}
			this.searchingRoutes = this.vAllRoutes.filter(function (route) {
				return teaweb.match(route.name, keyword)
			})
			if (this.searchingRoutes.length > 0) {
				this.routeCode = this.searchingRoutes[0].code + "@" + this.searchingRoutes[0].domainId
			} else {
				this.routeCode = ""
			}
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
				<select class="ui dropdown" style="width: 18em" v-model="routeCode">
					<option value="" v-if="keyword.length == 0">[请选择]</option>
					<option v-for="route in searchingRoutes" :value="route.code + '@' + route.domainId">{{route.name}}（{{route.domainName}}）</option>
				</select>
			</div>
			<div class="ui field">
				<input type="text" placeholder="搜索..." size="10" v-model="keyword" ref="keywordRef" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
			</div>
			
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
			</div>
			<div class="ui field">
				<a href="" @click.prevent="cancel()"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
</div>`
})