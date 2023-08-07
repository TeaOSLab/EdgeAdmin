Vue.component("http-cache-policy-selector", {
	props: ["v-cache-policy"],
	mounted: function () {
		let that = this
		Tea.action("/servers/components/cache/count")
			.post()
			.success(function (resp) {
				that.count = resp.data.count
			})
	},
	data: function () {
		let cachePolicy = this.vCachePolicy
		return {
			count: 0,
			cachePolicy: cachePolicy
		}
	},
	methods: {
		remove: function () {
			this.cachePolicy = null
		},
		select: function () {
			let that = this
			teaweb.popup("/servers/components/cache/selectPopup", {
				width: "42em",
				height: "26em",
				callback: function (resp) {
					that.cachePolicy = resp.data.cachePolicy
				}
			})
		},
		create: function () {
			let that = this
			teaweb.popup("/servers/components/cache/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.cachePolicy = resp.data.cachePolicy
				}
			})
		}
	},
	template: `<div>
	<div v-if="cachePolicy != null" class="ui label basic">
		<input type="hidden" name="cachePolicyId" :value="cachePolicy.id"/>
		{{cachePolicy.name}} &nbsp; <a :href="'/servers/components/cache/update?cachePolicyId=' + cachePolicy.id" target="_blank" title="修改"><i class="icon pencil small"></i></a>&nbsp; <a href="" @click.prevent="remove()" title="删除"><i class="icon remove small"></i></a>
	</div>
	<div v-if="cachePolicy == null">
		<span v-if="count > 0"><a href="" @click.prevent="select">[选择已有策略]</a> &nbsp; &nbsp; </span><a href="" @click.prevent="create">[创建新策略]</a>
	</div>
</div>`
})