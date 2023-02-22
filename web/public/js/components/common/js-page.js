Vue.component("js-page", {
	props: ["v-max"],
	data: function () {
		let max = this.vMax
		if (max == null) {
			max = 0
		}
		return {
			max: max,
			page: 1
		}
	},
	methods: {
		updateMax: function (max) {
			this.max = max
		},
		selectPage: function(page) {
			this.page = page
			this.$emit("change", page)
		}
	},
	template:`<div>
	<div class="page" v-if="max > 1">
		<a href="" v-for="i in max" :class="{active: i == page}" @click.prevent="selectPage(i)">{{i}}</a>
	</div>
</div>`
})