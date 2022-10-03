Vue.component("bits-var", {
	props: ["v-bits"],
	data: function () {
		let bits = this.vBits
		if (typeof bits != "number") {
			bits = 0
		}
		let format = teaweb.splitFormat(teaweb.formatBits(bits))
		return {
			format: format
		}
	},
	template:`<var class="normal">
	<span>{{format[0]}}</span>{{format[1]}}
</var>`
})