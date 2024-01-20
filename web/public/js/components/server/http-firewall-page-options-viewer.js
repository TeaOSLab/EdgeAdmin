Vue.component("http-firewall-page-options-viewer", {
	props: ["v-page-options"],
	data: function () {
		return {
			options: this.vPageOptions
		}
	},
	template: `<div>
	<span v-if="options == null">默认设置</span>
	<div v-else>
		状态码：{{options.status}} / 提示内容：<span v-if="options.body != null && options.body.length > 0">[{{options.body.length}}字符]</span>
	</div>
</div>	
`
})