Vue.component("firewall-syn-flood-config-viewer", {
	props: ["v-syn-flood-config"],
	data: function () {
		let config = this.vSynFloodConfig
		if (config == null) {
			config = {
				isOn: false,
				minAttempts: 10,
				timeoutSeconds: 600,
				ignoreLocal: true
			}
		}
		return {
			config: config
		}
	},
	template: `<div>
	<span v-if="config.isOn">
		已启用 / <span>空连接次数：{{config.minAttempts}}次/分钟</span> / 封禁时长：{{config.timeoutSeconds}}秒 <span v-if="config.ignoreLocal">/ 忽略局域网访问</span>
	</span>
	<span v-else>未启用</span>
</div>`
})