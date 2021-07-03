// 显示指标对象名
Vue.component("metric-key-label", {
	props: ["v-key"],
	data: function () {
		return {
			keyDefs: window.METRIC_HTTP_KEYS
		}
	},
	methods: {
		keyName: function (key) {
			let that = this
			let subKey = ""
			let def = this.keyDefs.$find(function (k, v) {
				if (v.code == key) {
					return true
				}
				if (key.startsWith("${arg.") && v.code.startsWith("${arg.")) {
					subKey = that.getSubKey("arg.", key)
					return true
				}
				if (key.startsWith("${header.") && v.code.startsWith("${header.")) {
					subKey = that.getSubKey("header.", key)
					return true
				}
				if (key.startsWith("${cookie.") && v.code.startsWith("${cookie.")) {
					subKey = that.getSubKey("cookie.", key)
					return true
				}
				return false
			})
			if (def != null) {
				if (subKey.length > 0) {
					return def.name + ": " + subKey
				}
				return def.name
			}
			return key
		},
		getSubKey: function (prefix, key) {
			prefix = "${" + prefix
			let index = key.indexOf(prefix)
			if (index >= 0) {
				key = key.substring(index + prefix.length)
				key = key.substring(0, key.length - 1)
				return key
			}
			return ""
		}
	},
	template: `<div class="ui label basic small">
	{{keyName(this.vKey)}}
</div>`
})