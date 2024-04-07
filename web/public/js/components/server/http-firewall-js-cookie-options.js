Vue.component("http-firewall-js-cookie-options", {
	props: ["v-js-cookie-options"],
	mounted: function () {
		this.updateSummary()
	},
	data: function () {
		let options = this.vJsCookieOptions
		if (options == null) {
			options = {
				life: 0,
				maxFails: 0,
				failBlockTimeout: 0,
				failBlockScopeAll: false,
				scope: "service"
			}
		}

		return {
			options: options,
			isEditing: false,
			summary: ""
		}
	},
	watch: {
		"options.life": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.life = i
			this.updateSummary()
		},
		"options.maxFails": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.maxFails = i
			this.updateSummary()
		},
		"options.failBlockTimeout": function (v) {
			let i = parseInt(v, 10)
			if (isNaN(i)) {
				i = 0
			}
			this.options.failBlockTimeout = i
			this.updateSummary()
		},
		"options.failBlockScopeAll": function (v) {
			this.updateSummary()
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		},
		updateSummary: function () {
			let summaryList = []
			if (this.options.life > 0) {
				summaryList.push("有效时间" + this.options.life + "秒")
			}
			if (this.options.maxFails > 0) {
				summaryList.push("最多失败" + this.options.maxFails + "次")
			}
			if (this.options.failBlockTimeout > 0) {
				summaryList.push("失败拦截" + this.options.failBlockTimeout + "秒")
			}
			if (this.options.failBlockScopeAll) {
				summaryList.push("尝试全局封禁")
			}

			if (summaryList.length == 0) {
				this.summary = "默认配置"
			} else {
				this.summary = summaryList.join(" / ")
			}
		},
		confirm: function () {
			this.isEditing = false
		}
	},
	template: `<div>
	<input type="hidden" name="jsCookieOptionsJSON" :value="JSON.stringify(options)"/>
	<a href="" @click.prevent="edit">{{summary}} <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<div v-show="isEditing" style="margin-top: 0.5em">
		<table class="ui table definition selectable">
			<tbody>
				<tr>
					<td class="title">有效时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">验证通过后在这个时间内不再验证，默认3600秒。</p>
					</td>
				</tr>
				<tr>
					<td>最多失败次数</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.maxFails" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">次</span>
						</div>
						<p class="comment"><span v-if="options.maxFails > 0 && options.maxFails < 5" class="red">建议填入一个不小于5的数字，以减少误判几率。</span>允许用户失败尝试的最多次数，超过这个次数将被自动加入黑名单。如果为空或者为0，表示不限制。</p>
					</td>
				</tr>
				<tr>
					<td>失败拦截时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.failBlockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">在达到最多失败次数（大于0）时，自动拦截的时长；如果为0表示不自动拦截。</p>
					</td>
				</tr>
				<tr>
					<td>失败全局封禁</td>
					<td>
						<checkbox v-model="options.failBlockScopeAll"></checkbox>
						<p class="comment">选中后，表示允许系统尝试全局封禁某个IP，以提升封禁性能。</p>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
`
})