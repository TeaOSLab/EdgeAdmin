Vue.component("bytes-var", {
	props: ["v-bytes"],
	data: function () {
		let bytes = this.vBytes
		if (typeof bytes != "number") {
			bytes = 0
		}
		let format = teaweb.splitFormat(teaweb.formatBytes(bytes))
		return {
			format: format
		}
	},
	template:`<var class="normal">
	<span>{{format[0]}}</span>{{format[1]}}
</var>`
})