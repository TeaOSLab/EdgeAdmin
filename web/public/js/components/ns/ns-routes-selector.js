// 选择多个线路
Vue.component("ns-routes-selector", {
	props: ["v-routes"],
	mounted: function () {
		let that = this
		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.routes = resp.data.routes
			})
	},
	data: function () {
		let selectedRoutes = this.vRoutes
		if (selectedRoutes == null) {
			selectedRoutes = []
		}

		return {
			routeCode: "default",
			routes: [],
			isAdding: false,
			routeType: "default",
			selectedRoutes: selectedRoutes
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
		},
		cancel: function () {
			this.isAdding = false
		},
		confirm: function () {
			if (this.routeCode.length == 0) {
				return
			}

			let that = this
			this.routes.forEach(function (v) {
				if (v.code == that.routeCode) {
					that.selectedRoutes.push(v)
				}
			})
			this.cancel()
		},
		remove: function (index) {
			this.selectedRoutes.$remove(index)
		}
	}
	,
	template: `<div>
	<div>
		<div class="ui label basic text small" v-for="(route, index) in selectedRoutes" style="margin-bottom: 0.3em">
			<input type="hidden" name="routeCodes" :value="route.code"/>
			{{route.name}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding" style="margin-bottom: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" v-model="routeType">
					<option value="default">[默认线路]</option>
					<option value="user">自定义线路</option>
					<option value="isp">运营商</option>
					<option value="china">中国省市</option>
					<option value="world">全球国家地区</option>
				</select>
			</div>
			
			<div class="ui field">
				<select class="ui dropdown" v-model="routeCode" style="width: 10em">
					<option v-for="route in routes" :value="route.code" v-if="route.type == routeType">{{route.name}}</option>
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