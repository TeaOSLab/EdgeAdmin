Vue.component("http-firewall-block-options", {
	props: ["v-block-options"],
	data: function () {
		return {
			blockOptions: this.vBlockOptions,
			statusCode: this.vBlockOptions.statusCode,
			timeout: this.vBlockOptions.timeout,
			isEditing: false
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
		},
		timeout: function (v) {
			let timeout = parseInt(v)
			if (isNaN(timeout)) {
				this.blockOptions.timeout = 0
			} else {
				this.blockOptions.timeout = timeout
			}
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	<input type="hidden" name="blockOptionsJSON" :value="JSON.stringify(blockOptions)"/>
	<a href="" @click.prevent="edit">状态码：{{statusCode}} / 提示内容：<span v-if="blockOptions.body != null && blockOptions.body.length > 0">[{{blockOptions.body.length}}字符]</span><span v-else class="disabled">[无]</span>  / 超时时间：{{timeout}}秒 <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
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
				<textarea rows="3" v-model="blockOptions.body"></textarea>
			</td>
		</tr>
		<tr>
			<td>超时时间</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeout" style="width: 5em" maxlength="6"/>
					<span class="ui label">秒</span>
				</div>
				<p class="comment">触发阻止动作时，封锁客户端IP的时间。</p>
			</td>
		</tr>
	</table>
</div>	
`
})