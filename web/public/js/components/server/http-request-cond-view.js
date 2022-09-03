Vue.component("http-request-cond-view", {
	props: ["v-cond"],
	data: function () {
		return {
			cond: this.vCond,
			components: window.REQUEST_COND_COMPONENTS
		}
	},
	methods: {
		typeName: function (cond) {
			let c = this.components.$find(function (k, v) {
				return v.type == cond.type
			})
			if (c != null) {
				return c.name;
			}
			return cond.param + " " + cond.operator
		},
		updateConds: function (conds, simpleCond) {
			for (let k in simpleCond) {
				if (simpleCond.hasOwnProperty(k)) {
					this.cond[k] = simpleCond[k]
				}
			}
		},
		notifyChange: function () {

		}
	},
	template: `<div style="margin-bottom: 0.5em">
	<span class="ui label small basic">
		<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
		<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{typeName(cond)}}: </var>
		{{cond.value}}
		<sup v-if="cond.isCaseInsensitive" title="不区分大小写"><i class="icon info small"></i></sup>
	</span>
</div>`
})