Vue.component("ns-access-log-ref-box", {
	props: ["v-access-log-ref", "v-is-parent"],
	data: function () {
		let config = this.vAccessLogRef
		if (config == null) {
			config = {
				isOn: false,
				isPrior: false,
				logMissingDomains: false
			}
		}
		if (typeof (config.logMissingDomains) == "undefined") {
			config.logMissingDomains = false
		}
		return {
			config: config
		}
	},
	template: `<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="!vIsParent"></prior-checkbox>
		<tbody v-show="vIsParent || config.isPrior">
			<tr>
				<td class="title">启用</td>
				<td>
					<checkbox name="isOn" value="1" v-model="config.isOn"></checkbox>
				</td>
			</tr>
			<tr>
				<td>只记录失败查询</td>
				<td>
					<checkbox v-model="config.missingRecordsOnly"></checkbox>
					<p class="comment">选中后，表示只记录查询失败的日志。</p>
				</td>
			</tr>
			<tr>
				<td>包含未添加的域名</td>
				<td>
					<checkbox name="logMissingDomains" value="1" v-model="config.logMissingDomains"></checkbox>
					<p class="comment">选中后，表示日志中包含对没有在系统里创建的域名访问。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})