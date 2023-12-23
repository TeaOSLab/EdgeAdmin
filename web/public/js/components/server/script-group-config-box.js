Vue.component("script-group-config-box", {
	props: ["v-group", "v-auditing-status", "v-is-location"],
	data: function () {
		let group = this.vGroup
		if (group == null) {
			group = {
				isPrior: false,
				isOn: true,
				scripts: []
			}
		}
		if (group.scripts == null) {
			group.scripts = []
		}

		let script = null
		if (group.scripts.length > 0) {
			script = group.scripts[group.scripts.length - 1]
		}

		return {
			group: group,
			script: script
		}
	},
	methods: {
		changeScript: function (script) {
			this.group.scripts = [script] // 目前只支持单个脚本
			this.change()
		},
		change: function () {
			this.$emit("change", this.group)
		}
	},
	template: `<div>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="group" v-if="vIsLocation"></prior-checkbox>
		</table>
		<div :style="{opacity: (!vIsLocation || group.isPrior) ? 1 : 0.5}">
			<script-config-box :v-script-config="script" :v-auditing-status="vAuditingStatus" comment="在接收到客户端请求之后立即调用。预置req、resp变量。" @change="changeScript" :v-is-location="vIsLocation"></script-config-box>
		</div>
</div>`
})