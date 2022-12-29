Vue.component("http-header-assistant", {
	props: ["v-type", "v-value"],
	mounted: function () {
		let that = this
		Tea.action("/servers/headers/options?type=" + this.vType)
			.post()
			.success(function (resp) {
				that.allHeaders = resp.data.headers
			})
	},
	data: function () {
		return {
			allHeaders: [],
			matchedHeaders: [],

			selectedHeaderName: ""
		}
	},
	watch: {
		vValue: function (v) {
			if (v != this.selectedHeaderName) {
				this.selectedHeaderName = ""
			}

			if (v.length == 0) {
				this.matchedHeaders = []
				return
			}
			this.matchedHeaders = this.allHeaders.filter(function (header) {
				return teaweb.match(header, v)
			}).slice(0, 10)
		}
	},
	methods: {
		select: function (header) {
			this.$emit("select", header)
			this.selectedHeaderName = header
		}
	},
	template: `<span v-if="selectedHeaderName.length == 0">
	<a href="" v-for="header in matchedHeaders" class="ui label basic tiny blue" style="font-weight: normal; margin-bottom: 0.3em" @click.prevent="select(header)">{{header}}</a>
	<span v-if="matchedHeaders.length > 0">&nbsp; &nbsp;</span>
</span>`
})