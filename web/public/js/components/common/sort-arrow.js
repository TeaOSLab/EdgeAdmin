// 排序使用的箭头
Vue.component("sort-arrow", {
	props: ["name"],
	data: function () {
		let url = window.location.toString()
		let order = ""
		let newArgs = []
		if (window.location.search != null && window.location.search.length > 0) {
			let queryString = window.location.search.substring(1)
			let pieces = queryString.split("&")
			let that = this
			pieces.forEach(function (v) {
				let eqIndex = v.indexOf("=")
				if (eqIndex > 0) {
					let argName = v.substring(0, eqIndex)
					let argValue = v.substring(eqIndex + 1)
					if (argName == that.name) {
						order = argValue
					} else if (argValue != "asc" && argValue != "desc") {
						newArgs.push(v)
					}
				} else {
					newArgs.push(v)
				}
			})
		}
		if (order == "asc") {
			newArgs.push(this.name + "=desc")
		} else if (order == "desc") {
			newArgs.push(this.name + "=asc")
		} else {
			newArgs.push(this.name + "=desc")
		}

		let qIndex = url.indexOf("?")
		if (qIndex > 0) {
			url = url.substring(0, qIndex) + "?" + newArgs.join("&")
		} else {
			url = url + "?" + newArgs.join("&")
		}

		return {
			order: order,
			url: url
		}
	},
	template: `<a :href="url" title="排序">&nbsp; <i class="ui icon long arrow small" :class="{down: order == 'asc', up: order == 'desc', 'down grey': order == '' || order == null}"></i></a>`
})