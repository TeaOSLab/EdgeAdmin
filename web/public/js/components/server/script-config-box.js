Vue.component("script-config-box", {
	props: ["id", "v-script-config", "comment"],
	data: function () {
		let config = this.vScriptConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				code: ""
			}
		}

		if (config.code.length == 0) {
			config.code = "\n\n\n\n"
		}

		return {
			config: config
		}
	},
	watch: {
		"config.isOn": function () {
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.config)
		},
		changeCode: function (code) {
			this.config.code = code
			this.change()
		}
	},
	template: `<div>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">启用脚本设置</td>
				<td><checkbox v-model="config.isOn"></checkbox></td>
			</tr>
		</tbody>
		<tbody>
			<tr :style="{opacity: !config.isOn ? 0.5 : 1}">
				<td>脚本代码</td>	
				<td><source-code-box :id="id" type="text/javascript" :read-only="false" @change="changeCode">{{config.code}}</source-code-box>
					<p class="comment">{{comment}}</p>
				</td>
			</tr>
		</tbody>
	</table>
</div>`
})