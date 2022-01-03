Vue.component("http-request-scripts-config-box", {
	props: ["vRequestScriptsConfig"],
	data: function () {
		let config = this.vRequestScriptsConfig
		if (config == null) {
			config = {}
		}
		return {
			config: config
		}
	},
	methods: {
		changeInitGroup: function (group) {
			this.config.initGroup = group
			this.$forceUpdate()
		},
		changeRequestGroup: function (group) {
			this.config.requestGroup = group
			this.$forceUpdate()
		}
	},
	template: `<div>
	<input type="hidden" name="requestScriptsJSON" :value="JSON.stringify(config)"/>
	<div class="margin"></div>
	<h4 style="margin-bottom: 0">请求初始化</h4>
	<p class="comment">在请求刚初始化时调用，此时自定义Header等尚未生效。</p>
	<div>
		<script-group-config-box :v-group="config.initGroup" @change="changeInitGroup"></script-group-config-box>
	</div>
	<h4 style="margin-bottom: 0">准备发送请求</h4>
	<p class="comment">在准备执行请求或者转发请求之前调用，此时自定义Header、源站等已准备好。</p>
	<div>
		<script-group-config-box :v-group="config.requestGroup" @change="changeRequestGroup"></script-group-config-box>
	</div>
	<div class="margin"></div>
</div>`
})