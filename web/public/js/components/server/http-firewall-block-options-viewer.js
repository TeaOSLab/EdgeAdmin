Vue.component("http-firewall-block-options-viewer", {
	props: ["v-block-options"],
	data: function () {
		return {
			options: this.vBlockOptions
		}
	},
	template: `<div>
	<span v-if="options == null">默认设置</span>
	<div v-else>
		状态码：{{options.statusCode}} / 提示内容：<span v-if="options.body != null && options.body.length > 0">[{{options.body.length}}字符]</span><span v-else class="disabled">[无]</span>  / 超时时间：{{options.timeout}}秒 <span v-if="options.timeoutMax > options.timeout">/ 最大封禁时长：{{options.timeoutMax}}秒</span>
		<span v-if="options.failBlockScopeAll"> / 尝试全局封禁</span>
	</div>
</div>	
`
})