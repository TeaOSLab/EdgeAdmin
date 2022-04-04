Vue.component("http-cache-stale-config", {
	props: ["v-cache-stale-config"],
	data: function () {
		let config = this.vCacheStaleConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				status: [],
				supportStaleIfErrorHeader: true,
				life: {
					count: 1,
					unit: "day"
				}
			}
		}
		return {
			config: config
		}
	},
	watch: {
		config: {
			deep: true,
			handler: function () {
				this.$emit("change", this.config)
			}
		}
	},
	methods: {},
	template: `<table class="ui table definition selectable">
	<tbody>
		<tr>
			<td class="title">启用过时缓存</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment"><plus-label></plus-label>选中后，在更新缓存失败后会尝试读取过时的缓存。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>有效期</td>
			<td>
				<time-duration-box :v-value="config.life"></time-duration-box>
				<p class="comment">缓存在过期之后，仍然保留的时间。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>状态码</td>
			<td><http-status-box :v-status-list="config.status"></http-status-box>
				<p class="comment">在这些状态码出现时使用过时缓存，默认支持<code-label>50x</code-label>状态码。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>支持stale-if-error</td>
			<td>
				<checkbox v-model="config.supportStaleIfErrorHeader"></checkbox>
				<p class="comment">选中后，支持在Cache-Control中通过<code-label>stale-if-error</code-label>指定过时缓存有效期。</p>
			</td>
		</tr>
	</tbody>
</table>`
})