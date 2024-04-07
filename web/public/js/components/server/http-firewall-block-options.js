Vue.component("http-firewall-block-options", {
	props: ["v-block-options"],
	data: function () {
		return {
			options: this.vBlockOptions,
			statusCode: this.vBlockOptions.statusCode,
			timeout: this.vBlockOptions.timeout,
			timeoutMax: this.vBlockOptions.timeoutMax,
			isEditing: false
		}
	},
	watch: {
		statusCode: function (v) {
			let statusCode = parseInt(v)
			if (isNaN(statusCode)) {
				this.options.statusCode = 403
			} else {
				this.options.statusCode = statusCode
			}
		},
		timeout: function (v) {
			let timeout = parseInt(v)
			if (isNaN(timeout)) {
				this.options.timeout = 0
			} else {
				this.options.timeout = timeout
			}
		},
		timeoutMax: function (v) {
			let timeoutMax = parseInt(v)
			if (isNaN(timeoutMax)) {
				this.options.timeoutMax = 0
			} else {
				this.options.timeoutMax = timeoutMax
			}
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	<input type="hidden" name="blockOptionsJSON" :value="JSON.stringify(options)"/>
	<a href="" @click.prevent="edit">状态码：{{statusCode}} / 提示内容：<span v-if="options.body != null && options.body.length > 0">[{{options.body.length}}字符]</span><span v-else class="disabled">[无]</span> <span v-if="timeout > 0"> / 封禁时长：{{timeout}}秒</span>
	 <span v-if="timeoutMax > timeout"> / 最大封禁时长：{{timeoutMax}}秒</span>
	 <span v-if="options.failBlockScopeAll"> / 尝试全局封禁</span>
	 <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<table class="ui table" v-show="isEditing">
		<tr>
			<td class="title">状态码</td>
			<td>
				<input type="text" v-model="statusCode" style="width:4.5em" maxlength="3"/>
			</td>
		</tr>
		<tr>
			<td>提示内容</td>
			<td>
				<textarea rows="3" v-model="options.body"></textarea>
			</td>
		</tr>
		<tr>
			<td>封禁时长</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeout" style="width: 5em" maxlength="6"/>
					<span class="ui label">秒</span>
				</div>
				<p class="comment">触发阻止动作时，封禁客户端IP的时间。</p>
			</td>
		</tr>
		<tr>
			<td>最大封禁时长</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeoutMax" style="width: 5em" maxlength="6"/>
					<span class="ui label">秒</span>
				</div>
				<p class="comment">如果最大封禁时长大于封禁时长（{{timeout}}秒），那么表示每次封禁的时候，将会在这两个时长数字之间随机选取一个数字作为最终的封禁时长。</p>
			</td>
		</tr>
		<tr>
			<td>失败全局封禁</td>
			<td>
				<checkbox v-model="options.failBlockScopeAll"></checkbox>
				<p class="comment">选中后，表示允许系统尝试全局封禁某个IP，以提升封禁性能。</p>
			</td>
		</tr>
	</table>
</div>	
`
})