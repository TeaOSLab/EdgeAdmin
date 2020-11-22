Vue.component("http-firewall-block-options", {
	props: ["v-block-options"],
	data: function () {
		return {
			blockOptions: this.vBlockOptions,
			statusCode: this.vBlockOptions.statusCode
		}
	},
	watch: {
		statusCode: function (v) {
			let statusCode = parseInt(v)
			if (isNaN(statusCode)) {
				this.blockOptions.statusCode = 403
			} else {
				this.blockOptions.statusCode = statusCode
			}
		}
	},
	template: `<div>
<input type="hidden" name="blockOptionsJSON" :value="JSON.stringify(blockOptions)"/>
	<table class="ui table">
		<tr>
			<td class="title">状态码</td>
			<td>
				<input type="text" v-model="statusCode" style="width:4.5em" maxlength="3"/>
			</td>
		</tr>
		<tr>
			<td>提示内容</td>
			<td>
				<textarea rows="3" v-model="blockOptions.body"></textarea>
			</td>
		</tr>
	</table>
</div>	
`
})