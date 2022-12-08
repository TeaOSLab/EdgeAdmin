Vue.component("email-sender", {
	props: ["value", "name"],
	data: function () {
		let value = this.value
		if (value == null) {
			value = {
				isOn: false,
				smtpHost: "",
				smtpPort: 0,
				username: "",
				password: "",
				fromEmail: "",
				fromName: ""
			}
		}
		let smtpPortString = value.smtpPort.toString()
		if (smtpPortString == "0") {
			smtpPortString = ""
		}

		return {
			config: value,
			smtpPortString: smtpPortString
		}
	},
	watch: {
		smtpPortString: function (v) {
			let port = parseInt(v)
			if (!isNaN(port)) {
				this.config.smtpPort = port
			}
		}
	},
	methods: {
		test: function () {
			window.TESTING_EMAIL_CONFIG = this.config
			teaweb.popup("/users/setting/emailTest", {
				height: "36em"
			})
		}
	},
	template: `<div>
	<input type="hidden" :name="name" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<tbody>
			<tr>
				<td class="title">启用</td>
				<td><checkbox v-model="config.isOn"></checkbox></td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>SMTP地址 *</td>
				<td>
					<input type="text" :name="name + 'SmtpHost'" v-model="config.smtpHost"/>
					<p class="comment">SMTP主机地址，比如<code-label>smtp.qq.com</code-label>，目前仅支持TLS协议，如不清楚，请查询对应邮件服务商文档。</p>
				</td>
			</tr>
			<tr>
				<td>SMTP端口 *</td>
				<td>
					<input type="text" :name="name + 'SmtpPort'" v-model="smtpPortString" style="width: 5em" maxlength="5"/>
					<p class="comment">SMTP主机端口，比如<code-label>587</code-label>、<code-label>465</code-label>，如不清楚，请查询对应邮件服务商文档。</p>
				</td>
			</tr>
			<tr>
				<td>用户名 *</td>
				<td>
					<input type="text" :name="name + 'Username'" v-model="config.username"/>
					<p class="comment">通常为发件人邮箱地址。</p>
				</td>
			</tr>
			<tr>
				<td>密码 *</td>
				<td>
					<input type="password" :name="name + 'Password'" v-model="config.password"/>
					<p class="comment">邮箱登录密码或授权码，如不清楚，请查询对应邮件服务商文档。。</p>
				</td>
			</tr>
			<tr>
				<td>发件人Email *</td>
				<td>
					<input type="text" :name="name + 'FromEmail'" v-model="config.fromEmail" maxlength="128"/>
					<p class="comment">使用的发件人邮箱地址，通常和发件用户名一致。</p>
				</td>
			</tr>
			<tr>
				<td>发件人名称</td>
				<td>
					<input type="text" :name="name + 'FromName'" v-model="config.fromName" maxlength="30"/>
					<p class="comment">使用的发件人名称，默认使用系统设置的<a href="/settings/ui" target="_blank">产品名称</a>。</p>
				</td>
			</tr>
			<tr>
				<td>发送测试</td>
				<td><a href="" @click.prevent="test">[点此测试]</a></td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})