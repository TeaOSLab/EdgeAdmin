Vue.component("page-size-selector", {
	data: function () {
		let query = window.location.search
		let pageSize = 10
		if (query.length > 0) {
			query = query.substr(1)
			let params = query.split("&")
			params.forEach(function (v) {
				let pieces = v.split("=")
				if (pieces.length == 2 && pieces[0] == "pageSize") {
					let pageSizeString = pieces[1]
					if (pageSizeString.match(/^\d+$/)) {
						pageSize = parseInt(pageSizeString, 10)
						if (isNaN(pageSize) || pageSize < 1) {
							pageSize = 10
						}
					}
				}
			})
		}
		return {
			pageSize: pageSize
		}
	},
	watch: {
		pageSize: function () {
			window.ChangePageSize(this.pageSize)
		}
	},
	template: `<select class="ui dropdown" style="height:34px;padding-top:0;padding-bottom:0;margin-left:1em;color:#666" v-model="pageSize">
\t<option value="10">[每页]</option><option value="10" selected="selected">10条</option><option value="20">20条</option><option value="30">30条</option><option value="40">40条</option><option value="50">50条</option><option value="60">60条</option><option value="70">70条</option><option value="80">80条</option><option value="90">90条</option><option value="100">100条</option>
</select>`
})