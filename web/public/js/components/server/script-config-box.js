Vue.component("script-config-box", {
	props: ["id", "v-script-config", "comment", "v-auditing-status"],
	mounted: function () {
		let that = this
		setTimeout(function () {
			that.$forceUpdate()
		}, 100)
	},
	data: function () {
		let config = this.vScriptConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				code: "",
				auditingCode: ""
			}
		}

		let auditingStatus = null
		if (config.auditingCodeMD5 != null && config.auditingCodeMD5.length > 0 && config.auditingCode != null && config.auditingCode.length > 0) {
			config.code = config.auditingCode

			if (this.vAuditingStatus != null) {
				for (let i = 0; i < this.vAuditingStatus.length; i++) {
					let status = this.vAuditingStatus[i]
					if (status.md5 == config.auditingCodeMD5) {
						auditingStatus = status
						break
					}
				}
			}
		}

		if (config.code.length == 0) {
			config.code = "\n\n\n\n"
		}

		return {
			config: config,
			auditingStatus: auditingStatus
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
		},
		isPlus: function () {
			if (Tea == null || Tea.Vue == null) {
				return false
			}
			return Tea.Vue.teaIsPlus
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
				<td>
					<p class="comment" v-if="auditingStatus != null">
						<span class="green" v-if="auditingStatus.isPassed">管理员审核结果：审核通过。</span>
						<span class="red" v-else-if="auditingStatus.isRejected">管理员审核结果：驳回 &nbsp; &nbsp; 驳回理由：{{auditingStatus.rejectedReason}}</span>
						<span class="red" v-else>当前脚本将在审核后生效，请耐心等待审核结果。 <a href="/servers/user-scripts" target="_blank" v-if="isPlus()">去审核 &raquo;</a></span>
					</p>
					<p class="comment" v-if="auditingStatus == null"><span class="green">管理员审核结果：审核通过。</span></p>
					<source-code-box :id="id" type="text/javascript" :read-only="false" @change="changeCode">{{config.code}}</source-code-box>
					<p class="comment">{{comment}}</p>
				</td>
			</tr>
		</tbody>
	</table>
</div>`
})