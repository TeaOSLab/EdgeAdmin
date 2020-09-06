Vue.component("page-box", {
	data: function () {
		return {
			page: ""
		}
	},
	created: function () {
		let that = this;
		setTimeout(function () {
			that.page = Tea.Vue.page;
		})
	},
	template: `<div>
	<div class="page" v-html="page"></div>
</div>`
})