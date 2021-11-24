// 通用设置
Vue.component("http-common-config-box", {
	props: ["v-common-config"],
	data: function () {
		let config = this.vCommonConfig
		if (config == null) {
			config = {
				mergeSlashes: false
			}
		}
		return {
			config: config
		}
	},
	template: `<div>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">合并重复的路径分隔符</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" name="mergeSlashes" value="1" v-model="config.mergeSlashes"/>
					<label></label>
				</div>
				<p class="comment">合并URL中重复的路径分隔符为一个，比如<code-label>//hello/world</code-label>中的<code-label>//</code-label>。</p>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`
})