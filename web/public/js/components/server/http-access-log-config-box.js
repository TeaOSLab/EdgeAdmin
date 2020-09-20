Vue.component("http-access-log-config-box", {
	props: ["v-access-log-config", "v-fields", "v-default-field-codes", "v-access-log-policies"],
	data: function () {
		let that = this

		// 初始化
		setTimeout(function () {
			that.changeFields()
			that.changePolicy()
		}, 100)

		let accessLog = {
			isOn: true,
			fields: [],
			status1: true,
			status2: true,
			status3: true,
			status4: true,
			status5: true,

			storageOnly: false,
			storagePolicies: []
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
		this.vAccessLogPolicies.forEach(function (v) {
			v.isChecked = accessLog.storagePolicies.$contains(v.id)
		})

		return {
			accessLog: accessLog
		}
	},
	methods: {
		changeFields: function () {
			this.accessLog.fields = this.vFields.filter(function (v) {
				return v.isChecked
			}).map(function (v) {
				return v.code
			})
		},
		changePolicy: function () {
			this.accessLog.storagePolicies = this.vAccessLogPolicies.filter(function (v) {
				return v.isChecked
			}).map(function (v) {
				return v.id
			})
		}
	},
	template: `<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(accessLog)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">是否开启访问日志存储</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.isOn"/>
						<label></label>
					</div>
					<p class="comment">关闭访问日志，并不影响统计的运行。</p>
				</td>
			</tr>
		</tbody>
		<tbody  v-show="accessLog.isOn">
			<tr>
				<td>要存储的访问日志字段</td>
				<td>
					<div class="ui checkbox" v-for="field in vFields" style="width:10em;margin-bottom:0.8em">
						<input type="checkbox" v-model="field.isChecked" @change="changeFields"/>
						<label>{{field.name}}</label>
					</div>
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
				<td>选择输出的日志策略</td>
				<td>
					<span class="disabled" v-if="vAccessLogPolicies.length == 0">暂时还没有缓存策略。</span>
					<div v-if="vAccessLogPolicies.length > 0">
						<div class="ui checkbox" v-for="policy in vAccessLogPolicies" style="width:10em;margin-bottom:0.8em">
							<input type="checkbox" v-model="policy.isChecked" @change="changePolicy" />
							<label>{{policy.name}}</label>
						</div>
					</div>
				</td>
			</tr>
			<tr>
				<td>是否只输出到日志策略</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.storageOnly"/>
						<label></label>
					</div>
					<p class="comment">选中表示只输出日志到日志策略，而停止默认的日志存储。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})