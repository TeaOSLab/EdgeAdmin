Vue.component("http-access-log-config-box", {
	props: ["v-access-log-config", "v-fields", "v-default-field-codes", "v-is-location", "v-is-group"],
	data: function () {
		let that = this

		// 初始化
		setTimeout(function () {
			that.changeFields()
		}, 100)

		let accessLog = {
			isPrior: false,
			isOn: false,
			fields: [1, 2, 6, 7],
			status1: true,
			status2: true,
			status3: true,
			status4: true,
			status5: true,

            firewallOnly: false,
			enableClientClosed: false
		}
		if (this.vAccessLogConfig != null) {
			accessLog = this.vAccessLogConfig
		}

		this.vFields.forEach(function (v) {
			if (that.vAccessLogConfig == null) { // 初始化默认值
				v.isChecked = that.vDefaultFieldCodes.$contains(v.code)
			} else {
				v.isChecked = accessLog.fields.$contains(v.code)
			}
		})

		return {
			accessLog: accessLog,
			hasRequestBodyField: this.vFields.$contains(8),
			showAdvancedOptions: false
		}
	},
	methods: {
		changeFields: function () {
			this.accessLog.fields = this.vFields.filter(function (v) {
				return v.isChecked
			}).map(function (v) {
				return v.code
			})
			this.hasRequestBodyField = this.accessLog.fields.$contains(8)
		},
		changeAdvanced: function (v) {
			this.showAdvancedOptions = v
		}
	},
	template: `<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(accessLog)"/>
	<table class="ui table definition selectable" :class="{'opacity-mask': this.accessLog.firewallOnly}">
		<prior-checkbox :v-config="accessLog" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || accessLog.isPrior">
			<tr>
				<td class="title">启用访问日志</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn">
			<tr>
				<td colspan="2"><more-options-indicator @change="changeAdvanced"></more-options-indicator></td>
			</tr>
		</tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn && showAdvancedOptions">
			<tr>
				<td>基础信息</td>
				<td><p class="comment" style="padding-top: 0">默认记录客户端IP、请求URL等基础信息。</p></td>
			</tr>
			<tr>
				<td>高级信息</td>
				<td>
					<div class="ui checkbox" v-for="(field, index) in vFields" style="width:10em;margin-bottom:0.8em">
						<input type="checkbox" v-model="field.isChecked" @change="changeFields" :id="'access-log-field-' + index"/>
						<label :for="'access-log-field-' + index">{{field.name}}</label>
					</div>
					<p class="comment">在基础信息之外要存储的信息。
						<span class="red" v-if="hasRequestBodyField">记录"请求Body"将会显著消耗更多的系统资源，建议仅在调试时启用，最大记录尺寸为2MiB。</span>
					</p>
				</td>
			</tr>
			<tr>
				<td>要存储的访问日志状态码</td>
				<td>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status1"/>
						<label>1xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status2"/>
						<label>2xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status3"/>
						<label>3xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status4"/>
						<label>4xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status5"/>
						<label>5xx</label>
					</div>
				</td>
			</tr>
			<tr>
				<td>记录客户端中断日志</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.enableClientClosed"/>
						<label></label>
					</div>
					<p class="comment">以<code-label>499</code-label>的状态码记录客户端主动中断日志。</p>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn">
        <h4>WAF相关</h4>
        <table class="ui table definition selectable">
            <tr>
                <td class="title">只记录WAF相关日志</td>
                <td>
                    <checkbox v-model="accessLog.firewallOnly"></checkbox>
                    <p class="comment">选中后只记录WAF相关的日志。通过此选项可有效减少访问日志数量，降低网络带宽和存储压力。</p>
                </td>
            </tr>
        </table>
    </div>
	<div class="margin"></div>
</div>`
})