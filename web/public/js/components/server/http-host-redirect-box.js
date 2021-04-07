Vue.component("http-host-redirect-box", {
	props: ["v-redirects"],
	data: function () {
		let redirects = this.vRedirects
		if (redirects == null) {
			redirects = []
		}

		return {
			redirects: redirects,
			statusOptions: [
				{"code": 301, "text": "Moved Permanently"},
				{"code": 308, "text": "Permanent Redirect"},
				{"code": 302, "text": "Found"},
				{"code": 303, "text": "See Other"},
				{"code": 307, "text": "Temporary Redirect"}
			]
		}
	},
	methods: {
		add: function () {
			let that = this
			window.UPDATING_REDIRECT = null

			teaweb.popup("/servers/server/settings/redirects/createPopup", {
				height: "22em",
				callback: function (resp) {
					that.redirects.push(resp.data.redirect)
				}
			})
		},
		update: function (index, redirect) {
			let that = this
			window.UPDATING_REDIRECT = redirect

			teaweb.popup("/servers/server/settings/redirects/createPopup", {
				height: "22em",
				callback: function (resp) {
					Vue.set(that.redirects, index, resp.data.redirect)
				}
			})
		},
		remove: function (index) {
			this.redirects.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="hostRedirectsJSON" :value="JSON.stringify(redirects)"/>

	<!-- TODO 将来支持排序，并支持isOn切换 -->
	<div v-if="redirects.length > 0">
		<div v-for="(redirect, index) in redirects" class="ui label basic small" style="margin-bottom: 0.5em;margin-top: 0.5em">
			<span v-if="redirect.status > 0">[{{redirect.status}}]</span><span v-if="redirect.matchPrefix">[prefix]</span> {{redirect.beforeURL}} -&gt; {{redirect.afterURL}} <a href="" @click.prevent="update(index, redirect)" title="修改"><i class="icon pencil small"></i></a> &nbsp; <a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>	
	</div>
	<div>
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`
})