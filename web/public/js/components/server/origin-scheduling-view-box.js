Vue.component("origin-scheduling-view-box", {
	props: ["v-scheduling", "v-params"],
	data: function () {
		let scheduling = this.vScheduling
		if (scheduling == null) {
			scheduling = {}
		}
		return {
			scheduling: scheduling
		}
	},
	methods: {
		update: function () {
			teaweb.popup("/servers/server/settings/reverseProxy/updateSchedulingPopup?" + this.vParams, {
				height: "21em",
				callback: function () {
					window.location.reload()
				},
			})
		}
	},
	template: `<div>
	<div class="margin"></div>
	<table class="ui table selectable definition">
		<tr>
			<td class="title">当前正在使用的算法</td>
			<td>
				{{scheduling.name}} &nbsp; <a href="" @click.prevent="update()"><span>[修改]</span></a>
				<p class="comment">{{scheduling.description}}</p>
			</td>
		</tr>
	</table>
</div>`
})