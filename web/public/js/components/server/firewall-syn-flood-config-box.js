Vue.component("firewall-syn-flood-config-box", {
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
			config: config,
			isEditing: false,
			minAttempts: config.minAttempts,
			timeoutSeconds: config.timeoutSeconds
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	watch: {
		minAttempts: function (v) {
			let count = parseInt(v)
			if (isNaN(count)) {
				count = 10
			}
			if (count < 5) {
				count = 5
			}
			this.config.minAttempts = count
		},
		timeoutSeconds: function (v) {
			let seconds = parseInt(v)
			if (isNaN(seconds)) {
				seconds = 10
			}
			if (seconds < 60) {
				seconds = 60
			}
			this.config.timeoutSeconds = seconds
		}
	},
	template: `<div>
	<input type="hidden" name="synFloodJSON" :value="JSON.stringify(config)"/>
	<a href="" @click.prevent="edit">
		<span v-if="config.isOn">
			已启用 / <span>空连接次数：{{config.minAttempts}}次/分钟</span> / 封禁时长：{{config.timeoutSeconds}}秒 <span v-if="config.ignoreLocal">/ 忽略局域网访问</span>
		</span>
		<span v-else>未启用</span>
		<i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i>
	</a>
	
	<table class="ui table selectable" v-show="isEditing">
		<tr>
			<td class="title">启用</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment">启用后，WAF将会尝试自动检测并阻止SYN Flood攻击。此功能需要节点已安装并启用nftables或Firewalld。</p>
			</td>
		</tr>
		<tr>
			<td>空连接次数</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="minAttempts" style="width: 5em" maxlength="4"/>
					<span class="ui label">次/分钟</span>
				</div>
				<p class="comment">超过此数字的"空连接"将被视为SYN Flood攻击，为了防止误判，此数值默认不小于5。</p>
			</td>
		</tr>
		<tr>
			<td>封禁时长</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeoutSeconds" style="width: 5em" maxlength="8"/>
					<span class="ui label">秒</span>
				</div>
			</td>
		</tr>
		<tr>
			<td>忽略局域网访问</td>
			<td>
				<checkbox v-model="config.ignoreLocal"></checkbox>
			</td>
		</tr>
	</table>
</div>`
})