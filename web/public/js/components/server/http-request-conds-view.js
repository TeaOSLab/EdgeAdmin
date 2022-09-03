// 浏览条件列表
Vue.component("http-request-conds-view", {
	props: ["v-conds"],
	data: function () {
		let conds = this.vConds
		if (conds == null) {
			conds = {
				isOn: true,
				connector: "or",
				groups: []
			}
		}
		if (conds.groups == null) {
			conds.groups = []
		}

		let that = this
		conds.groups.forEach(function (group) {
			group.conds.forEach(function (cond) {
				cond.typeName = that.typeName(cond)
			})
		})

		return {
			initConds: conds
		}
	},
	computed: {
		// 之所以使用computed，是因为需要动态更新
		conds: function () {
			return this.initConds
		}
	},
	methods: {
		typeName: function (cond) {
			let c = window.REQUEST_COND_COMPONENTS.$find(function (k, v) {
				return v.type == cond.type
			})
			if (c != null) {
				return c.name;
			}
			return cond.param + " " + cond.operator
		},
		updateConds: function (conds) {
			this.initConds = conds
		},
		notifyChange: function () {
			let that = this
			if (this.initConds.groups != null) {
				this.initConds.groups.forEach(function (group) {
					group.conds.forEach(function (cond) {
						cond.typeName = that.typeName(cond)
					})
				})
				this.$forceUpdate()
			}
		}
	},
	template: `<div>
		<div v-if="conds.groups.length > 0">
			<div v-for="(group, groupIndex) in conds.groups">
				<var v-for="(cond, index) in group.conds" style="font-style: normal;display: inline-block; margin-bottom:0.5em">
					<span class="ui label small basic" style="line-height: 1.5">
						<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
						<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{cond.typeName}}: </var>
						{{cond.value}}
						<sup v-if="cond.isCaseInsensitive" title="不区分大小写"><i class="icon info small"></i></sup>
					</span>
					
					<var v-if="index < group.conds.length - 1"> {{group.connector}} &nbsp;</var>
				</var>
				<div class="ui divider" v-if="groupIndex != conds.groups.length - 1" style="margin-top:0.3em;margin-bottom:0.5em"></div>
				<div>
					<span class="ui label tiny olive" v-if="group.description != null && group.description.length > 0">{{group.description}}</span>
				</div>
			</div>	
		</div>
	</div>	
</div>`
})