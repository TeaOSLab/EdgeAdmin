// 将变量转换为中文
Vue.component("request-variables-describer", {
	data: function () {
		return {
			vars:[]
		}
	},
	methods: {
		update: function (variablesString) {
			this.vars = []
			let that = this
			variablesString.replace(/\${.+?}/g, function (v) {
				let def = that.findVar(v)
				if (def == null) {
					return v
				}
				that.vars.push(def)
			})
		},
		findVar: function (name) {
			let def = null
			window.REQUEST_VARIABLES.forEach(function (v) {
				if (v.code == name) {
					def = v
				}
			})
			return def
		}
	},
	template: `<span>
	<span v-for="(v, index) in vars"><code-label :title="v.description">{{v.code}}</code-label> - {{v.name}}<span v-if="index < vars.length-1">；</span></span>
</span>`
})
